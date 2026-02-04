package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/RayUI/RayUI/internal/autostart"
	"github.com/RayUI/RayUI/internal/core"
	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/parser"
	"github.com/RayUI/RayUI/internal/service"
	"github.com/RayUI/RayUI/internal/store"
	"github.com/RayUI/RayUI/internal/sysproxy"
	"github.com/RayUI/RayUI/internal/util"
)

// App is the main Wails binding struct that connects all backend services.
type App struct {
	ctx context.Context

	// Stores
	configStore  *store.ConfigStore
	profileStore *store.ProfileStore
	subStore     *store.SubStore
	routingStore *store.RoutingStore
	dnsStore     *store.DNSStore
	statsStore   *store.StatsStore

	// Services
	coreManager          core.CoreManager
	subscriptionService  *service.SubscriptionService
	statsService         *service.StatsService
	speedTestService     *service.SpeedTestService
	updaterService       *service.UpdaterService
	connectionsService   *service.ConnectionsService
	geoDataService       *service.GeoDataService

	// System integration
	sysProxy     sysproxy.SysProxy
	autoStart    autostart.AutoStart
	powerMonitor *PowerMonitor

	// Logging
	logWriter *core.LogWriter

	// Quit flag - when true, allow app to close instead of hiding to tray
	quitRequested bool
}

// NewApp creates a new App instance with all dependencies.
func NewApp() *App {
	ps := store.NewProfileStore()
	ss := store.NewSubStore()
	sts := store.NewStatsStore()

	return &App{
		configStore:  store.NewConfigStore(),
		profileStore: ps,
		subStore:     ss,
		routingStore: store.NewRoutingStore(),
		dnsStore:     store.NewDNSStore(),
		statsStore:   sts,

		subscriptionService: &service.SubscriptionService{
			SubStore:     ss,
			ProfileStore: ps,
			StatsStore:   sts,
		},
		statsService:     &service.StatsService{StatsStore: sts},
		speedTestService:   &service.SpeedTestService{},
		updaterService:     &service.UpdaterService{DataDir: util.AppDataDir()},
		connectionsService: service.NewConnectionsService(),
		geoDataService:     service.NewGeoDataService(util.AppDataDir()),

		sysProxy:  sysproxy.NewSysProxy(),
		autoStart: autostart.NewAutoStart(),
	}
}

// Startup is called when the Wails app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Ensure directory structure.
	dataDir := util.AppDataDir()
	for _, sub := range []string{"cores", "logs", "data"} {
		_ = util.EnsureDir(filepath.Join(dataDir, sub))
	}

	// Wire stats callback to emit events.
	a.statsService.SetCallback(func(stats model.TrafficStats) {
		a.emitTraffic(stats)
	})

	// Wire updater progress callback.
	a.updaterService.OnProgress = func(p service.UpdateProgress) {
		a.emitUpdateProgress(p)
	}

	// Wire auto-sync notification callback and start background updater.
	a.subscriptionService.OnAutoSync = func(subID string, count int, err error) {
		if err != nil {
			a.emitNotification("error", "Auto-Sync Failed", err.Error())
		} else {
			a.emitNotification("info", "Auto-Sync", fmt.Sprintf("Updated %d profiles", count))
		}
	}
	a.subscriptionService.StartAutoUpdate()

	// Start power monitor to handle sleep/wake events.
	a.powerMonitor = NewPowerMonitor(a.handlePowerEvent)
	a.powerMonitor.Start()
}

// Shutdown is called when the Wails app is closing.
func (a *App) Shutdown(ctx context.Context) {
	if a.powerMonitor != nil {
		a.powerMonitor.Stop()
		a.powerMonitor.StopPlatform()
	}
	a.subscriptionService.StopAutoUpdate()
	a.statsService.StopCollecting()
	if a.coreManager != nil && a.coreManager.IsRunning() {
		_ = a.coreManager.Stop()
	}
	if a.logWriter != nil {
		_ = a.logWriter.Close()
		a.logWriter = nil
	}
	_ = a.sysProxy.Clear()
}

// handlePowerEvent is called when system sleep/wake events occur.
func (a *App) handlePowerEvent(isSleeping bool) {
	if isSleeping {
		// System is going to sleep - pause stats collection
		a.statsService.StopCollecting()
	} else {
		// System woke up - verify core is alive and resume stats
		if a.coreManager != nil && a.coreManager.IsRunning() {
			// Check if core process is still alive
			status := a.coreManager.GetStatus()
			if !status.Running {
				// Core died during sleep, try to restart
				_ = a.StartCore()
			} else {
				// Core is alive, resume stats collection
				cfg, err := a.configStore.Load()
				if err == nil && cfg.ActiveProfileID != "" {
					coreType := status.CoreType
					a.statsService.StartCollecting(cfg.ActiveProfileID, coreType)
				}
			}
		}
	}
}

// --- Core bindings ---

// StartCore starts the proxy core with the active profile.
func (a *App) StartCore() error {
	cfg, err := a.configStore.Load()
	if err != nil {
		return err
	}
	if cfg.ActiveProfileID == "" {
		return fmt.Errorf("no active profile selected")
	}

	profile, err := a.profileStore.GetByID(cfg.ActiveProfileID)
	if err != nil {
		return err
	}
	if profile == nil {
		return fmt.Errorf("active profile not found")
	}

	routing, err := a.getActiveRouting(cfg)
	if err != nil {
		return err
	}
	dns, err := a.dnsStore.Load()
	if err != nil {
		return err
	}

	// Determine core type and create manager.
	coreType := core.SelectCore(*profile)
	a.coreManager = core.NewCoreManager(coreType, util.AppDataDir())

	// Ensure geo data files exist before starting core.
	_ = a.geoDataService.EnsureGeoData(coreType)

	// Set up log writer to file + event streaming.
	logDir := filepath.Join(util.AppDataDir(), "logs")
	lw, err := core.NewLogWriter(logDir)
	if err == nil {
		lw.SetCallback(func(line string) {
			a.emitCoreLog(line)
		})
		a.logWriter = lw
		a.coreManager.SetLogWriter(lw)
	}

	if err := a.coreManager.Start(*profile, *routing, dns, cfg); err != nil {
		if a.logWriter != nil {
			a.logWriter.Close()
			a.logWriter = nil
		}
		return err
	}

	// Apply proxy mode.
	if cfg.ProxyMode == model.ProxyModeSystem {
		httpPort, socksPort := 10809, 10808
		for _, ib := range cfg.Inbounds {
			if ib.Protocol == "http" {
				httpPort = ib.Port
			}
			if ib.Protocol == "socks" {
				socksPort = ib.Port
			}
		}
		_ = a.sysProxy.Set("127.0.0.1", httpPort, "127.0.0.1", socksPort)
	}

	// Start stats collection.
	a.statsService.StartCollecting(profile.ID, coreType)
	a.emitCoreStatus(a.coreManager.GetStatus())

	return nil
}

// StopCore stops the running proxy core.
func (a *App) StopCore() error {
	a.statsService.StopCollecting()
	_ = a.sysProxy.Clear()

	if a.coreManager != nil {
		if err := a.coreManager.Stop(); err != nil {
			return err
		}
		a.emitCoreStatus(a.coreManager.GetStatus())
	}
	if a.logWriter != nil {
		_ = a.logWriter.Close()
		a.logWriter = nil
	}
	return nil
}

// RestartCore restarts the proxy core.
func (a *App) RestartCore() error {
	if err := a.StopCore(); err != nil {
		return err
	}
	return a.StartCore()
}

// GetCoreStatus returns the current core status.
func (a *App) GetCoreStatus() model.CoreStatus {
	if a.coreManager == nil {
		return model.CoreStatus{}
	}
	return a.coreManager.GetStatus()
}

// --- Profile bindings ---

// GetProfiles returns all profiles, optionally filtered by subscription ID.
func (a *App) GetProfiles(subID string) ([]model.ProfileItem, error) {
	if subID != "" {
		return a.profileStore.GetBySubID(subID)
	}
	return a.profileStore.GetAll()
}

// AddProfile adds a new profile.
func (a *App) AddProfile(profile model.ProfileItem) error {
	if profile.ID == "" {
		profile.ID = util.GenerateUUID()
	}
	return a.profileStore.Add(profile)
}

// UpdateProfile updates an existing profile.
func (a *App) UpdateProfile(profile model.ProfileItem) error {
	return a.profileStore.Update(profile)
}

// DeleteProfiles deletes profiles by IDs.
func (a *App) DeleteProfiles(ids []string) error {
	return a.profileStore.Delete(ids)
}

// SetActiveProfile sets the active profile and restarts core if running.
func (a *App) SetActiveProfile(id string) error {
	cfg, err := a.configStore.Load()
	if err != nil {
		return err
	}
	cfg.ActiveProfileID = id
	if err := a.configStore.Save(cfg); err != nil {
		return err
	}
	if a.coreManager != nil && a.coreManager.IsRunning() {
		return a.RestartCore()
	}
	return nil
}

// TestProfiles runs latency tests on the given profile IDs.
func (a *App) TestProfiles(ids []string) ([]model.SpeedTestResult, error) {
	allProfiles, err := a.profileStore.GetAll()
	if err != nil {
		return nil, err
	}
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}
	var toTest []model.ProfileItem
	for _, p := range allProfiles {
		if _, ok := idSet[p.ID]; ok {
			toTest = append(toTest, p)
		}
	}
	cfg, _ := a.configStore.Load()
	timeout := time.Duration(cfg.SpeedTest.Timeout) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return a.speedTestService.TestProfiles(toTest, cfg.SpeedTest.Concurrent, timeout), nil
}

// TestAllProfiles tests all profiles.
func (a *App) TestAllProfiles() ([]model.SpeedTestResult, error) {
	profiles, err := a.profileStore.GetAll()
	if err != nil {
		return nil, err
	}
	cfg, _ := a.configStore.Load()
	timeout := time.Duration(cfg.SpeedTest.Timeout) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return a.speedTestService.TestProfiles(profiles, cfg.SpeedTest.Concurrent, timeout), nil
}

// --- Subscription bindings ---

// GetSubscriptions returns all subscriptions.
func (a *App) GetSubscriptions() ([]model.SubItem, error) {
	return a.subStore.GetAll()
}

// AddSubscription adds a new subscription.
func (a *App) AddSubscription(sub model.SubItem) error {
	if sub.ID == "" {
		sub.ID = util.GenerateUUID()
	}
	return a.subStore.Add(sub)
}

// UpdateSubscription updates a subscription.
func (a *App) UpdateSubscription(sub model.SubItem) error {
	return a.subStore.Update(sub)
}

// DeleteSubscription deletes a subscription and its associated profiles.
func (a *App) DeleteSubscription(id string) error {
	profiles, _ := a.profileStore.GetBySubID(id)
	if len(profiles) > 0 {
		ids := make([]string, len(profiles))
		for i, p := range profiles {
			ids[i] = p.ID
		}
		_ = a.profileStore.Delete(ids)
	}
	return a.subStore.Delete(id)
}

// SyncSubscription syncs a single subscription.
func (a *App) SyncSubscription(id string) (int, error) {
	count, err := a.subscriptionService.Sync(id)
	if err != nil {
		a.emitNotification("error", "Sync Failed", err.Error())
		return 0, err
	}
	a.emitNotification("success", "Sync Complete", fmt.Sprintf("Imported %d profiles", count))
	return count, nil
}

// SyncAllSubscriptions syncs all subscriptions.
func (a *App) SyncAllSubscriptions() (map[string]int, error) {
	return a.subscriptionService.SyncAll()
}

// --- Routing bindings ---

// GetRoutings returns all routing items.
func (a *App) GetRoutings() ([]model.RoutingItem, error) {
	return a.routingStore.GetAll()
}

// AddRouting adds a routing item.
func (a *App) AddRouting(item model.RoutingItem) error {
	if item.ID == "" {
		item.ID = util.GenerateUUID()
	}
	return a.routingStore.Add(item)
}

// UpdateRouting updates a routing item.
func (a *App) UpdateRouting(item model.RoutingItem) error {
	return a.routingStore.Update(item)
}

// DeleteRouting deletes a routing item.
func (a *App) DeleteRouting(id string) error {
	return a.routingStore.Delete(id)
}

// SetActiveRouting sets the active routing and restarts core if running.
func (a *App) SetActiveRouting(id string) error {
	cfg, err := a.configStore.Load()
	if err != nil {
		return err
	}
	cfg.ActiveRoutingID = id
	if err := a.configStore.Save(cfg); err != nil {
		return err
	}
	if a.coreManager != nil && a.coreManager.IsRunning() {
		return a.RestartCore()
	}
	return nil
}

// --- DNS bindings ---

// GetDNSConfig returns the DNS configuration.
func (a *App) GetDNSConfig() (model.DNSItem, error) {
	return a.dnsStore.Load()
}

// UpdateDNSConfig updates the DNS configuration.
func (a *App) UpdateDNSConfig(dns model.DNSItem) error {
	if err := a.dnsStore.Save(dns); err != nil {
		return err
	}
	if a.coreManager != nil && a.coreManager.IsRunning() {
		return a.RestartCore()
	}
	return nil
}

// --- Config bindings ---

// GetConfig returns the global application configuration.
func (a *App) GetConfig() (model.Config, error) {
	return a.configStore.Load()
}

// UpdateConfig updates the global configuration.
func (a *App) UpdateConfig(cfg model.Config) error {
	if err := a.configStore.Save(cfg); err != nil {
		return err
	}
	if a.coreManager != nil && a.coreManager.IsRunning() {
		return a.RestartCore()
	}
	return nil
}

// --- System bindings ---

// SystemInfo holds platform information.
type SystemInfo struct {
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	AppVersion string `json:"appVersion"`
}

// appVersion is set from main.go via SetAppVersion.
var appVersion string

// SetAppVersion is called from main.go to pass the build-time version.
func SetAppVersion(v string) {
	appVersion = v
}

// GetSystemInfo returns platform information.
func (a *App) GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:         util.GetOS(),
		Arch:       util.GetArch(),
		AppVersion: appVersion,
	}
}

// SetProxyMode changes the proxy mode.
func (a *App) SetProxyMode(mode model.EProxyMode) error {
	cfg, err := a.configStore.Load()
	if err != nil {
		return err
	}
	cfg.ProxyMode = mode
	if err := a.configStore.Save(cfg); err != nil {
		return err
	}
	if a.coreManager != nil && a.coreManager.IsRunning() {
		return a.RestartCore()
	}
	return nil
}

// SetAutoStart enables or disables launch at login.
func (a *App) SetAutoStart(enabled bool) error {
	if enabled {
		// Use the current executable path.
		return a.autoStart.Enable("")
	}
	return a.autoStart.Disable()
}

// IsAutoStartEnabled checks if auto-start is enabled.
func (a *App) IsAutoStartEnabled() bool {
	enabled, _ := a.autoStart.IsEnabled()
	return enabled
}

// --- Import / Export bindings ---

// ImportFromText parses share links from text and imports them.
func (a *App) ImportFromText(text string) (int, error) {
	items, err := parser.ParseBatch(text)
	if err != nil {
		return 0, err
	}
	for _, item := range items {
		if item.ID == "" {
			item.ID = util.GenerateUUID()
		}
		if err := a.profileStore.Add(item); err != nil {
			continue
		}
	}
	return len(items), nil
}

// ExportShareLink converts a profile to a share URI.
func (a *App) ExportShareLink(id string) (string, error) {
	profile, err := a.profileStore.GetByID(id)
	if err != nil {
		return "", err
	}
	if profile == nil {
		return "", fmt.Errorf("profile not found")
	}
	return parser.ToShareURI(*profile)
}

// --- Stats bindings ---

// GetStats returns all server statistics.
func (a *App) GetStats() ([]model.ServerStatItem, error) {
	return a.statsStore.GetAll()
}

// ResetStats resets statistics for a profile.
func (a *App) ResetStats(profileID string) error {
	return a.statsStore.DeleteByProfileID(profileID)
}

// ResetAllStats resets all statistics.
func (a *App) ResetAllStats() error {
	return a.statsStore.Clear()
}

// --- Logs bindings ---

// GetLogs returns the last N lines from the core log file.
func (a *App) GetLogs(lines int) []string {
	logPath := filepath.Join(util.AppDataDir(), "logs", "core.log")
	f, err := os.Open(logPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var all []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		all = append(all, scanner.Text())
	}

	if lines <= 0 || lines >= len(all) {
		return all
	}
	return all[len(all)-lines:]
}

// ClearLogs truncates the core log file.
func (a *App) ClearLogs() error {
	logPath := filepath.Join(util.AppDataDir(), "logs", "core.log")
	return os.Truncate(logPath, 0)
}

// --- Core availability bindings ---

// IsCoreInstalled checks if a specific core binary is installed.
func (a *App) IsCoreInstalled(coreType model.ECoreType) bool {
	var binaryPath string
	dataDir := util.AppDataDir()
	switch coreType {
	case model.CoreSingbox:
		binaryPath = filepath.Join(dataDir, "cores", "sing-box")
	case model.CoreXray:
		binaryPath = filepath.Join(dataDir, "cores", "xray")
	default:
		return false
	}
	_, err := os.Stat(binaryPath)
	return err == nil
}

// --- Core Update bindings ---

// CheckCoreUpdate checks for a new version of the given core.
func (a *App) CheckCoreUpdate(coreType model.ECoreType) (*service.UpdateInfo, error) {
	return a.updaterService.CheckUpdate(coreType)
}

// DownloadCoreUpdate downloads and installs a core update.
func (a *App) DownloadCoreUpdate(info service.UpdateInfo) error {
	return a.updaterService.DownloadUpdate(info)
}

// --- GeoData bindings ---

// GetGeoDataInfo returns information about installed geo data files.
func (a *App) GetGeoDataInfo(coreType model.ECoreType) service.GeoDataInfo {
	return a.geoDataService.GetGeoDataInfo(coreType)
}

// EnsureGeoData ensures geo data files exist, downloading if necessary.
func (a *App) EnsureGeoData(coreType model.ECoreType) error {
	return a.geoDataService.EnsureGeoData(coreType)
}

// UpdateGeoData downloads the latest geo data files.
func (a *App) UpdateGeoData(coreType model.ECoreType) error {
	return a.geoDataService.UpdateGeoData(coreType)
}

// CheckGeoDataUpdate checks if newer geo data is available.
func (a *App) CheckGeoDataUpdate(coreType model.ECoreType) (bool, string, error) {
	return a.geoDataService.CheckGeoDataUpdate(coreType)
}

// Context returns the Wails runtime context (for use by system tray).
func (a *App) Context() context.Context {
	return a.ctx
}

// ToggleCore starts or stops the core depending on its current state.
func (a *App) ToggleCore() {
	if a.coreManager != nil && a.coreManager.IsRunning() {
		_ = a.StopCore()
	} else {
		_ = a.StartCore()
	}
}

// ShouldCloseToTray returns true if the window should hide instead of quit.
// Returns false if a quit was explicitly requested (e.g., from tray menu).
func (a *App) ShouldCloseToTray() bool {
	// If quit was explicitly requested, don't hide to tray
	if a.quitRequested {
		return false
	}
	cfg, err := a.configStore.Load()
	if err != nil {
		return false
	}
	return cfg.UI.CloseToTray
}

// RequestQuit sets the quit flag to allow the app to close.
// This should be called before triggering wailsruntime.Quit().
func (a *App) RequestQuit() {
	a.quitRequested = true
}

// --- Connections bindings ---

// GetConnections returns all active connections via Clash API.
func (a *App) GetConnections() (*service.ConnectionsResponse, error) {
	return a.connectionsService.GetConnections()
}

// CloseConnection closes a specific connection by ID.
func (a *App) CloseConnection(id string) error {
	return a.connectionsService.CloseConnection(id)
}

// CloseAllConnections closes all active connections.
func (a *App) CloseAllConnections() error {
	return a.connectionsService.CloseAllConnections()
}

// --- helpers ---

func (a *App) getActiveRouting(cfg model.Config) (*model.RoutingItem, error) {
	if cfg.ActiveRoutingID != "" {
		r, err := a.routingStore.GetByID(cfg.ActiveRoutingID)
		if err == nil && r != nil {
			return r, nil
		}
	}
	// Fallback: return the first enabled routing.
	routings, err := a.routingStore.GetAll()
	if err != nil {
		return nil, err
	}
	for i := range routings {
		if routings[i].Enabled {
			return &routings[i], nil
		}
	}
	if len(routings) > 0 {
		return &routings[0], nil
	}
	return nil, fmt.Errorf("no routing available")
}

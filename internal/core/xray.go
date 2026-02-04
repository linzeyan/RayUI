package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/RayUI/RayUI/internal/config"
	"github.com/RayUI/RayUI/internal/model"
)

// XrayCore manages the xray-core binary lifecycle.
type XrayCore struct {
	dataDir   string
	cmd       *exec.Cmd
	cancel    context.CancelFunc
	mu        sync.Mutex
	status    model.CoreStatus
	generator config.XrayConfigGenerator
	logWriter io.Writer
}

func NewXrayCore(dataDir string) *XrayCore {
	return &XrayCore{
		dataDir: dataDir,
		status: model.CoreStatus{
			CoreType: model.CoreXray,
		},
	}
}

func (c *XrayCore) Start(profile model.ProfileItem, routing model.RoutingItem, dns model.DNSItem, cfg model.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd != nil {
		return fmt.Errorf("xray already running")
	}

	cfgBytes, err := c.generator.Generate(profile, routing, dns, cfg)
	if err != nil {
		return fmt.Errorf("generate config: %w", err)
	}

	cfgPath := filepath.Join(c.dataDir, "cores", "running-config.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(cfgPath, cfgBytes, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	bin := c.BinaryPath()
	cmd := exec.CommandContext(ctx, bin, "run", "-c", cfgPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// Tell Xray where to find geoip.dat / geosite.dat.
	cmd.Env = append(os.Environ(), "XRAY_LOCATION_ASSET="+filepath.Join(c.dataDir, "data"))
	if c.logWriter != nil {
		cmd.Stdout = c.logWriter
		cmd.Stderr = c.logWriter
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("start xray: %w", err)
	}

	c.cmd = cmd
	now := time.Now().Unix()
	c.status = model.CoreStatus{
		Running:   true,
		CoreType:  model.CoreXray,
		StartTime: &now,
		PID:       cmd.Process.Pid,
		Profile:   profile.Remarks,
	}

	go func() {
		_ = cmd.Wait()
		c.mu.Lock()
		defer c.mu.Unlock()
		c.cmd = nil
		c.status.Running = false
		c.status.PID = 0
	}()

	return nil
}

func (c *XrayCore) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}

	_ = c.cmd.Process.Signal(syscall.SIGTERM)
	c.cancel()

	done := make(chan struct{})
	go func() {
		_ = c.cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		_ = c.cmd.Process.Kill()
		<-done
	}

	c.cmd = nil
	c.status.Running = false
	c.status.PID = 0
	return nil
}

func (c *XrayCore) Restart() error {
	return c.Stop()
}

func (c *XrayCore) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cmd != nil
}

func (c *XrayCore) GetStatus() model.CoreStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := c.status
	if v, err := c.versionLocked(); err == nil {
		s.Version = v
	}
	return s
}

func (c *XrayCore) GenerateConfig(profile model.ProfileItem, routing model.RoutingItem, dns model.DNSItem, cfg model.Config) ([]byte, error) {
	return c.generator.Generate(profile, routing, dns, cfg)
}

func (c *XrayCore) CoreType() model.ECoreType {
	return model.CoreXray
}

func (c *XrayCore) Version() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.versionLocked()
}

func (c *XrayCore) versionLocked() (string, error) {
	vFile := filepath.Join(c.dataDir, "cores", "xray.version")
	data, err := os.ReadFile(vFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (c *XrayCore) SetLogWriter(w io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logWriter = w
}

func (c *XrayCore) BinaryPath() string {
	return filepath.Join(c.dataDir, "cores", "xray")
}

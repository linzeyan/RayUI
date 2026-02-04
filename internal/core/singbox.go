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

// SingboxCore manages the sing-box binary lifecycle.
type SingboxCore struct {
	dataDir   string
	cmd       *exec.Cmd
	cancel    context.CancelFunc
	mu        sync.Mutex
	status    model.CoreStatus
	generator config.SingboxConfigGenerator
	logWriter io.Writer
}

func NewSingboxCore(dataDir string) *SingboxCore {
	return &SingboxCore{
		dataDir: dataDir,
		status: model.CoreStatus{
			CoreType: model.CoreSingbox,
		},
	}
}

func (c *SingboxCore) Start(profile model.ProfileItem, routing model.RoutingItem, dns model.DNSItem, cfg model.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd != nil {
		return fmt.Errorf("sing-box already running")
	}

	// Generate config.
	cfgBytes, err := c.generator.Generate(profile, routing, dns, cfg)
	if err != nil {
		return fmt.Errorf("generate config: %w", err)
	}

	// Write config file.
	cfgPath := filepath.Join(c.dataDir, "cores", "running-config.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(cfgPath, cfgBytes, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	// Start process.
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	bin := c.BinaryPath()
	cmd := exec.CommandContext(ctx, bin, "run", "-c", cfgPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if c.logWriter != nil {
		cmd.Stdout = c.logWriter
		cmd.Stderr = c.logWriter
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("start sing-box: %w", err)
	}

	c.cmd = cmd
	now := time.Now().Unix()
	c.status = model.CoreStatus{
		Running:   true,
		CoreType:  model.CoreSingbox,
		StartTime: &now,
		PID:       cmd.Process.Pid,
		Profile:   profile.Remarks,
	}

	// Wait in background to update status on exit.
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

func (c *SingboxCore) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}

	// Graceful shutdown.
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

func (c *SingboxCore) Restart() error {
	status := c.GetStatus()
	if err := c.Stop(); err != nil {
		return err
	}
	_ = status // Restart requires the caller to re-supply params.
	return nil
}

func (c *SingboxCore) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cmd != nil
}

func (c *SingboxCore) GetStatus() model.CoreStatus {
	c.mu.Lock()
	defer c.mu.Unlock()
	s := c.status
	if v, err := c.versionLocked(); err == nil {
		s.Version = v
	}
	return s
}

func (c *SingboxCore) GenerateConfig(profile model.ProfileItem, routing model.RoutingItem, dns model.DNSItem, cfg model.Config) ([]byte, error) {
	return c.generator.Generate(profile, routing, dns, cfg)
}

func (c *SingboxCore) CoreType() model.ECoreType {
	return model.CoreSingbox
}

func (c *SingboxCore) Version() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.versionLocked()
}

func (c *SingboxCore) versionLocked() (string, error) {
	vFile := filepath.Join(c.dataDir, "cores", "sing-box.version")
	data, err := os.ReadFile(vFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (c *SingboxCore) SetLogWriter(w io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logWriter = w
}

func (c *SingboxCore) BinaryPath() string {
	return filepath.Join(c.dataDir, "cores", "sing-box")
}

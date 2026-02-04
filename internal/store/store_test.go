package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/util"
)

// withTempAppDir overrides AppDataDir to use a temp directory for testing.
func withTempAppDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir) // AppDataDir reads $HOME
	return filepath.Join(dir, ".RayUI")
}

func TestStoreGenericLoadSave(t *testing.T) {
	withTempAppDir(t)

	type item struct {
		Name string `json:"name"`
	}

	s := NewStore[[]item]("test_generic.json", []item{{Name: "default"}})

	// First load should write and return default.
	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got) != 1 || got[0].Name != "default" {
		t.Fatalf("got %+v, want [{Name: default}]", got)
	}

	// Save new data.
	if err := s.Save([]item{{Name: "a"}, {Name: "b"}}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err = s.Load()
	if err != nil {
		t.Fatalf("Load after save: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}

func TestStoreDefaultOnMissing(t *testing.T) {
	withTempAppDir(t)

	s := NewStore[model.Config]("config_test.json", model.DefaultConfig())
	cfg, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.UI.Language != "en" {
		t.Errorf("expected default language 'en', got %q", cfg.UI.Language)
	}

	// File should now exist.
	if _, err := os.Stat(s.GetPath()); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestProfileStoreCRUD(t *testing.T) {
	withTempAppDir(t)

	ps := NewProfileStore()

	// Empty initially.
	all, err := ps.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Fatalf("expected 0, got %d", len(all))
	}

	// Add.
	p := model.NewProfileItem()
	p.Remarks = "Test Server"
	p.Address = "1.2.3.4"
	p.Port = 443
	if err := ps.Add(p); err != nil {
		t.Fatal(err)
	}

	// Get by ID.
	got, err := ps.GetByID(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Remarks != "Test Server" {
		t.Fatalf("GetByID returned %+v", got)
	}

	// Update.
	p.Remarks = "Updated"
	if err := ps.Update(p); err != nil {
		t.Fatal(err)
	}
	got, _ = ps.GetByID(p.ID)
	if got.Remarks != "Updated" {
		t.Fatalf("expected Updated, got %q", got.Remarks)
	}

	// Delete.
	if err := ps.Delete([]string{p.ID}); err != nil {
		t.Fatal(err)
	}
	all, _ = ps.GetAll()
	if len(all) != 0 {
		t.Fatalf("expected 0 after delete, got %d", len(all))
	}
}

func TestStatsStoreUpdateAndReset(t *testing.T) {
	withTempAppDir(t)

	ss := NewStatsStore()

	// Update creates new record.
	if err := ss.UpdateStats("p1", 100, 200); err != nil {
		t.Fatal(err)
	}

	stat, err := ss.GetByProfileID("p1")
	if err != nil {
		t.Fatal(err)
	}
	if stat == nil {
		t.Fatal("expected stat record")
	}
	if stat.TotalUp != 100 || stat.TotalDown != 200 {
		t.Fatalf("unexpected totals: up=%d down=%d", stat.TotalUp, stat.TotalDown)
	}
	if stat.TodayUp != 100 || stat.TodayDown != 200 {
		t.Fatalf("unexpected today: up=%d down=%d", stat.TodayUp, stat.TodayDown)
	}

	// Accumulate.
	if err := ss.UpdateStats("p1", 50, 75); err != nil {
		t.Fatal(err)
	}
	stat, _ = ss.GetByProfileID("p1")
	if stat.TotalUp != 150 || stat.TotalDown != 275 {
		t.Fatalf("unexpected accumulated totals: up=%d down=%d", stat.TotalUp, stat.TotalDown)
	}

	// ResetDaily should be a no-op today.
	if err := ss.ResetDaily(); err != nil {
		t.Fatal(err)
	}
	stat, _ = ss.GetByProfileID("p1")
	if stat.TodayUp != 150 {
		t.Fatalf("unexpected todayUp after reset: %d", stat.TodayUp)
	}

	// Simulate date change by directly editing the record.
	items, _ := ss.store.Load()
	items[0].DateNow = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_ = ss.store.Save(items)

	if err := ss.ResetDaily(); err != nil {
		t.Fatal(err)
	}
	stat, _ = ss.GetByProfileID("p1")
	if stat.TodayUp != 0 || stat.TodayDown != 0 {
		t.Fatalf("expected daily reset, got up=%d down=%d", stat.TodayUp, stat.TodayDown)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}
	for _, tt := range tests {
		got := util.FormatBytes(tt.input)
		if got != tt.want {
			t.Errorf("FormatBytes(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

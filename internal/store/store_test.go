package store

import (
	"os"
	"path/filepath"
	"sync"
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

func TestSubStoreCRUD(t *testing.T) {
	withTempAppDir(t)
	ss := NewSubStore()

	all, err := ss.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Fatalf("expected 0, got %d", len(all))
	}

	s1 := model.NewSubItem()
	s1.Remarks = "Sub1"
	s1.URL = "https://example.com/1"
	if err := ss.Add(s1); err != nil {
		t.Fatal(err)
	}

	got, err := ss.GetByID(s1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Remarks != "Sub1" {
		t.Fatalf("GetByID got %+v", got)
	}

	// Non-existent ID returns nil.
	got, err = ss.GetByID("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatal("expected nil for nonexistent ID")
	}

	// Update.
	s1.Remarks = "Updated"
	if err := ss.Update(s1); err != nil {
		t.Fatal(err)
	}
	got, _ = ss.GetByID(s1.ID)
	if got.Remarks != "Updated" {
		t.Fatalf("expected Updated, got %q", got.Remarks)
	}

	// Update non-existent returns error.
	fake := model.NewSubItem()
	fake.Remarks = "Fake"
	if err := ss.Update(fake); err == nil {
		t.Error("expected error for non-existent update")
	}

	// Delete.
	if err := ss.Delete(s1.ID); err != nil {
		t.Fatal(err)
	}
	all, _ = ss.GetAll()
	if len(all) != 0 {
		t.Fatalf("expected 0, got %d", len(all))
	}
}

func TestRoutingStoreLockedDelete(t *testing.T) {
	withTempAppDir(t)
	rs := NewRoutingStore()

	// Default routing items are locked.
	all, err := rs.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 4 {
		t.Fatalf("expected 4 defaults, got %d", len(all))
	}

	// Delete a locked item should fail.
	if err := rs.Delete(all[0].ID); err == nil {
		t.Error("expected error deleting locked routing item")
	}

	// Add a custom non-locked item and delete it.
	custom := model.RoutingItem{
		ID: "custom-1", Remarks: "Custom", Locked: false, Enabled: true,
	}
	if err := rs.Add(custom); err != nil {
		t.Fatal(err)
	}
	if err := rs.Delete("custom-1"); err != nil {
		t.Fatalf("Delete non-locked: %v", err)
	}
	all, _ = rs.GetAll()
	if len(all) != 4 {
		t.Fatalf("expected 4 after delete, got %d", len(all))
	}
}

func TestDNSStoreLoadSave(t *testing.T) {
	withTempAppDir(t)
	ds := NewDNSStore()

	dns, err := ds.Load()
	if err != nil {
		t.Fatal(err)
	}
	if dns.RemoteDNS == "" {
		t.Error("default RemoteDNS should not be empty")
	}

	dns.RemoteDNS = "https://custom.dns/query"
	dns.FakeIP = true
	if err := ds.Save(dns); err != nil {
		t.Fatal(err)
	}

	loaded, err := ds.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.RemoteDNS != "https://custom.dns/query" {
		t.Errorf("RemoteDNS = %q", loaded.RemoteDNS)
	}
	if !loaded.FakeIP {
		t.Error("FakeIP should be true")
	}
}

func TestConfigStoreLoadSave(t *testing.T) {
	withTempAppDir(t)
	cs := NewConfigStore()

	cfg, err := cs.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.UI.Language != "en" {
		t.Errorf("default language = %q", cfg.UI.Language)
	}

	cfg.UI.Language = "zh-TW"
	cfg.UI.Theme = "dark"
	if err := cs.Save(cfg); err != nil {
		t.Fatal(err)
	}

	loaded, err := cs.Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.UI.Language != "zh-TW" {
		t.Errorf("Language = %q", loaded.UI.Language)
	}
	if loaded.UI.Theme != "dark" {
		t.Errorf("Theme = %q", loaded.UI.Theme)
	}
}

func TestProfileStoreReplaceBySubID(t *testing.T) {
	withTempAppDir(t)
	ps := NewProfileStore()

	p1 := model.NewProfileItem()
	p1.SubID = "sub-A"
	p1.Remarks = "A1"
	p2 := model.NewProfileItem()
	p2.SubID = "sub-A"
	p2.Remarks = "A2"
	p3 := model.NewProfileItem()
	p3.SubID = "sub-B"
	p3.Remarks = "B1"
	_ = ps.Add(p1)
	_ = ps.Add(p2)
	_ = ps.Add(p3)

	// Replace sub-A with a single new profile.
	newP := model.NewProfileItem()
	newP.SubID = "sub-A"
	newP.Remarks = "A-New"
	if err := ps.ReplaceBySubID("sub-A", []model.ProfileItem{newP}); err != nil {
		t.Fatal(err)
	}

	all, _ := ps.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}

	subA, _ := ps.GetBySubID("sub-A")
	if len(subA) != 1 || subA[0].Remarks != "A-New" {
		t.Errorf("sub-A: %+v", subA)
	}

	subB, _ := ps.GetBySubID("sub-B")
	if len(subB) != 1 || subB[0].Remarks != "B1" {
		t.Errorf("sub-B should be untouched: %+v", subB)
	}
}

func TestStatsStoreDeleteAndClear(t *testing.T) {
	withTempAppDir(t)
	ss := NewStatsStore()

	_ = ss.UpdateStats("p1", 100, 200)
	_ = ss.UpdateStats("p2", 300, 400)

	if err := ss.DeleteByProfileID("p1"); err != nil {
		t.Fatal(err)
	}
	s, _ := ss.GetByProfileID("p1")
	if s != nil {
		t.Error("p1 should be deleted")
	}
	s, _ = ss.GetByProfileID("p2")
	if s == nil {
		t.Error("p2 should still exist")
	}

	_ = ss.UpdateStats("p3", 50, 60)
	if err := ss.Clear(); err != nil {
		t.Fatal(err)
	}
	all, _ := ss.GetAll()
	if len(all) != 0 {
		t.Errorf("expected 0 after clear, got %d", len(all))
	}
}

func TestStoreConcurrentAccess(t *testing.T) {
	withTempAppDir(t)

	type item struct{ V int }
	s := NewStore[[]item]("concurrent.json", []item{{V: 0}})
	_, _ = s.Load() // seed file

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = s.Save([]item{{V: val}})
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.Load()
		}()
	}
	wg.Wait()

	// Final load should not error.
	got, err := s.Load()
	if err != nil {
		t.Fatalf("final Load: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 item, got %d", len(got))
	}
}

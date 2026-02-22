package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRoleTypeFor(t *testing.T) {
	tests := []struct {
		role   string
		expect RoleType
	}{
		{"polecat", Autonomous},
		{"witness", Autonomous},
		{"refinery", Autonomous},
		{"deacon", Autonomous},
		{"boot", Autonomous},
		{"mayor", Interactive},
		{"crew", Interactive},
		{"unknown", Interactive},
		{"", Interactive},
	}
	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			got := RoleTypeFor(tt.role)
			if got != tt.expect {
				t.Errorf("RoleTypeFor(%q) = %q, want %q", tt.role, got, tt.expect)
			}
		})
	}
}

func TestEnsureSettingsAt_CreatesFile(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettingsAt(dir, Interactive, ".claude", "settings.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	info, err := os.Stat(settingsPath)
	if err != nil {
		t.Fatalf("settings file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("settings file is empty")
	}
}

func TestEnsureSettingsAt_DoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	claudeDir := filepath.Join(dir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte("custom"), 0600); err != nil {
		t.Fatal(err)
	}

	err := EnsureSettingsAt(dir, Interactive, ".claude", "settings.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}

	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "custom" {
		t.Errorf("settings file was overwritten; got %q, want %q", string(content), "custom")
	}
}

func TestEnsureSettingsAt_Autonomous(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettingsAt(dir, Autonomous, ".claude", "settings.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("failed to read settings: %v", err)
	}
	if len(content) == 0 {
		t.Error("autonomous settings file is empty")
	}
}

func TestEnsureSettingsAt_CustomDir(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettingsAt(dir, Interactive, "my-settings", "config.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}

	settingsPath := filepath.Join(dir, "my-settings", "config.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("settings file not created at custom path: %v", err)
	}
}

func TestEnsureSettings(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettings(dir, Interactive)
	if err != nil {
		t.Fatalf("EnsureSettings failed: %v", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("settings file not created: %v", err)
	}
}

func TestEnsureSettingsForRole(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettingsForRole(dir, "polecat")
	if err != nil {
		t.Fatalf("EnsureSettingsForRole failed: %v", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("settings file not created: %v", err)
	}
}

func TestEnsureSettingsForRoleAt(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettingsForRoleAt(dir, "witness", "custom-dir", "custom.json")
	if err != nil {
		t.Fatalf("EnsureSettingsForRoleAt failed: %v", err)
	}

	settingsPath := filepath.Join(dir, "custom-dir", "custom.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("settings file not created at custom path: %v", err)
	}

	// witness is autonomous — verify content is valid JSON
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("failed to read settings: %v", err)
	}
	if !json.Valid(content) {
		t.Errorf("settings file is not valid JSON")
	}
}

func TestEnsureSettingsAt_MkdirError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission test not reliable on Windows")
	}
	// Place a regular file where the settings directory would go,
	// so MkdirAll fails.
	dir := t.TempDir()
	blockingFile := filepath.Join(dir, ".claude")
	if err := os.WriteFile(blockingFile, []byte("not a dir"), 0600); err != nil {
		t.Fatal(err)
	}

	err := EnsureSettingsAt(dir, Interactive, ".claude", "settings.json")
	if err == nil {
		t.Fatal("expected error when directory creation blocked by file, got nil")
	}
}

func TestEnsureSettingsAt_WriteError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission test not reliable on Windows")
	}
	dir := t.TempDir()
	claudeDir := filepath.Join(dir, ".claude")
	// Create settings dir as read-only so WriteFile fails
	if err := os.MkdirAll(claudeDir, 0555); err != nil {
		t.Fatal(err)
	}

	err := EnsureSettingsAt(dir, Interactive, ".claude", "settings.json")
	if err == nil {
		t.Fatal("expected error when write is not permitted, got nil")
	}
}

func TestEnsureSettingsAt_ContentIsValidJSON(t *testing.T) {
	dir := t.TempDir()

	// Test interactive template
	err := EnsureSettingsAt(dir, Interactive, ".claude", "settings.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(dir, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("failed to read settings: %v", err)
	}
	if !json.Valid(content) {
		t.Error("interactive settings is not valid JSON")
	}

	// Test autonomous template
	dir2 := t.TempDir()
	err = EnsureSettingsAt(dir2, Autonomous, ".claude", "settings.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}
	content2, err := os.ReadFile(filepath.Join(dir2, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("failed to read settings: %v", err)
	}
	if !json.Valid(content2) {
		t.Error("autonomous settings is not valid JSON")
	}
}

func TestEnsureSettingsAt_FilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permission test")
	}
	dir := t.TempDir()

	err := EnsureSettingsAt(dir, Interactive, ".claude", "settings.json")
	if err != nil {
		t.Fatalf("EnsureSettingsAt failed: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("failed to stat settings: %v", err)
	}
	// File should be created with 0600 permissions
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("settings file permissions = %o, want 0600", perm)
	}
}

func TestEnsureSettingsAt_AutonomousVsInteractiveDiffer(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	if err := EnsureSettingsAt(dir1, Interactive, ".claude", "settings.json"); err != nil {
		t.Fatal(err)
	}
	if err := EnsureSettingsAt(dir2, Autonomous, ".claude", "settings.json"); err != nil {
		t.Fatal(err)
	}

	c1, _ := os.ReadFile(filepath.Join(dir1, ".claude", "settings.json"))
	c2, _ := os.ReadFile(filepath.Join(dir2, ".claude", "settings.json"))

	// Both should have hooks, but autonomous has SessionStart mail injection
	var m1, m2 map[string]interface{}
	if err := json.Unmarshal(c1, &m1); err != nil {
		t.Fatalf("interactive settings not valid JSON: %v", err)
	}
	if err := json.Unmarshal(c2, &m2); err != nil {
		t.Fatalf("autonomous settings not valid JSON: %v", err)
	}

	// Both should contain hooks configuration
	if _, ok := m1["hooks"]; !ok {
		t.Error("interactive settings missing hooks")
	}
	if _, ok := m2["hooks"]; !ok {
		t.Error("autonomous settings missing hooks")
	}
}

func TestEnsureSettingsForRole_InteractiveRole(t *testing.T) {
	dir := t.TempDir()

	err := EnsureSettingsForRole(dir, "mayor")
	if err != nil {
		t.Fatalf("EnsureSettingsForRole failed: %v", err)
	}

	settingsPath := filepath.Join(dir, ".claude", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("settings file not created: %v", err)
	}
	if !json.Valid(content) {
		t.Error("settings file is not valid JSON")
	}
}

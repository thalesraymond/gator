package config

import (
	"path/filepath"
	"strings"
	"testing"

	"os"
)

func TestReadConfigTable(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		writeFile bool
		wantErr   bool
		wantUser  string
		wantDBURL string
	}{
		{
			name:      "reads valid config",
			content:   `{"db_url":"postgres://localhost","current_user_name":"thales"}`,
			writeFile: true,
			wantUser:  "thales",
			wantDBURL: "postgres://localhost",
		},
		{
			name:      "returns error when config is missing",
			writeFile: false,
			wantErr:   true,
		},
		{
			name:      "returns error for invalid json",
			content:   `{"db_url":`,
			writeFile: true,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tmpHome := t.TempDir()
			t.Setenv("HOME", tmpHome)

			configPath := filepath.Join(tmpHome, ".gatorconfig.json")

			if tc.writeFile {
				if err := os.WriteFile(configPath, []byte(tc.content), 0600); err != nil {
					t.Fatalf("failed to write config fixture: %v", err)
				}
			}

			cfg, err := ReadConfig()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if cfg.CurrentUserName != tc.wantUser {
				t.Fatalf("expected user %q, got %q", tc.wantUser, cfg.CurrentUserName)
			}

			if cfg.DbUrl != tc.wantDBURL {
				t.Fatalf("expected db_url %q, got %q", tc.wantDBURL, cfg.DbUrl)
			}
		})
	}
}

func TestSetUserTable(t *testing.T) {
	tests := []struct {
		name         string
		initialDBURL string
		initialUser  string
		newUser      string
	}{
		{
			name:         "updates user and keeps db url",
			initialDBURL: "postgres://localhost:5432/gator",
			initialUser:  "olduser",
			newUser:      "newuser",
		},
		{
			name:         "allows user names with spaces",
			initialDBURL: "postgres://localhost:5432/gator",
			initialUser:  "someone",
			newUser:      "new user",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tmpHome := t.TempDir()
			t.Setenv("HOME", tmpHome)

			cfg := &Config{
				DbUrl:           tc.initialDBURL,
				CurrentUserName: tc.initialUser,
			}

			if err := cfg.SetUser(tc.newUser); err != nil {
				t.Fatalf("SetUser returned error: %v", err)
			}

			stored, err := ReadConfig()
			if err != nil {
				t.Fatalf("ReadConfig returned error: %v", err)
			}

			if stored.CurrentUserName != tc.newUser {
				t.Fatalf("expected user %q, got %q", tc.newUser, stored.CurrentUserName)
			}

			if stored.DbUrl != tc.initialDBURL {
				t.Fatalf("expected db_url %q, got %q", tc.initialDBURL, stored.DbUrl)
			}

			configPath := filepath.Join(tmpHome, ".gatorconfig.json")
			info, err := os.Stat(configPath)
			if err != nil {
				t.Fatalf("failed to stat config file: %v", err)
			}

			if info.Mode().Perm() != 0600 {
				t.Fatalf("expected file mode 0600, got %o", info.Mode().Perm())
			}

			bytes, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("failed to read config file: %v", err)
			}

			if !strings.Contains(string(bytes), `"current_user_name": "`+tc.newUser+`"`) {
				t.Fatalf("written config did not contain expected user: %s", string(bytes))
			}
		})
	}
}

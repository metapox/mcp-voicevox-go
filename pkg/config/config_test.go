package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Port)
	}

	if cfg.VoicevoxURL != "http://localhost:50021" {
		t.Errorf("Expected default VOICEVOX URL 'http://localhost:50021', got %s", cfg.VoicevoxURL)
	}

	if cfg.DefaultSpeaker != 3 {
		t.Errorf("Expected default speaker 3, got %d", cfg.DefaultSpeaker)
	}

	if cfg.EnablePlayback != false {
		t.Errorf("Expected default playback false, got %t", cfg.EnablePlayback)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set test environment variables
	os.Setenv("MCP_VOICEVOX_PORT", "9090")
	os.Setenv("MCP_VOICEVOX_URL", "http://test:50021")
	os.Setenv("MCP_VOICEVOX_DEFAULT_SPEAKER", "5")
	os.Setenv("MCP_VOICEVOX_ENABLE_PLAYBACK", "true")
	defer func() {
		os.Unsetenv("MCP_VOICEVOX_PORT")
		os.Unsetenv("MCP_VOICEVOX_URL")
		os.Unsetenv("MCP_VOICEVOX_DEFAULT_SPEAKER")
		os.Unsetenv("MCP_VOICEVOX_ENABLE_PLAYBACK")
	}()

	cfg := DefaultConfig()
	err := cfg.LoadFromEnv()

	if err != nil {
		t.Fatalf("LoadFromEnv failed: %v", err)
	}

	if cfg.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", cfg.Port)
	}

	if cfg.VoicevoxURL != "http://test:50021" {
		t.Errorf("Expected VOICEVOX URL 'http://test:50021', got %s", cfg.VoicevoxURL)
	}

	if cfg.DefaultSpeaker != 5 {
		t.Errorf("Expected speaker 5, got %d", cfg.DefaultSpeaker)
	}

	if cfg.EnablePlayback != true {
		t.Errorf("Expected playback true, got %t", cfg.EnablePlayback)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			config: &Config{
				Port:           0,
				VoicevoxURL:    "http://localhost:50021",
				DefaultSpeaker: 3,
				TempDir:        "/tmp",
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: &Config{
				Port:           70000,
				VoicevoxURL:    "http://localhost:50021",
				DefaultSpeaker: 3,
				TempDir:        "/tmp",
			},
			wantErr: true,
		},
		{
			name: "empty voicevox URL",
			config: &Config{
				Port:           8080,
				VoicevoxURL:    "",
				DefaultSpeaker: 3,
				TempDir:        "/tmp",
			},
			wantErr: true,
		},
		{
			name: "negative speaker ID",
			config: &Config{
				Port:           8080,
				VoicevoxURL:    "http://localhost:50021",
				DefaultSpeaker: -1,
				TempDir:        "/tmp",
			},
			wantErr: true,
		},
		{
			name: "empty temp dir",
			config: &Config{
				Port:           8080,
				VoicevoxURL:    "http://localhost:50021",
				DefaultSpeaker: 3,
				TempDir:        "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("New() returned nil config")
	}

	// Should have default values
	if cfg.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Port)
	}
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Config はアプリケーションの設定を管理する構造体です
type Config struct {
	// Server settings
	Port int `json:"port"`

	// VOICEVOX settings
	VoicevoxURL    string `json:"voicevox_url"`
	DefaultSpeaker int    `json:"default_speaker"`

	// Audio synthesis settings
	DefaultSpeedScale      float64 `json:"default_speed_scale"`
	DefaultPitchScale      float64 `json:"default_pitch_scale"`
	DefaultIntonationScale float64 `json:"default_intonation_scale"`
	DefaultVolumeScale     float64 `json:"default_volume_scale"`

	// File settings
	TempDir string `json:"temp_dir"`

	// Audio settings
	EnablePlayback bool `json:"enable_playback"`
}

// DefaultConfig はデフォルト設定を返します
func DefaultConfig() *Config {
	return &Config{
		Port:                   8080,
		VoicevoxURL:            "http://localhost:50021",
		DefaultSpeaker:         3,
		DefaultSpeedScale:      1.0,
		DefaultPitchScale:      0.0,
		DefaultIntonationScale: 1.0,
		DefaultVolumeScale:     1.0,
		TempDir:                os.TempDir(),
		EnablePlayback:         false,
	}
}

// LoadFromEnv は環境変数から設定を読み込みます
func (c *Config) LoadFromEnv() error {
	if envPort := os.Getenv("MCP_VOICEVOX_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			c.Port = p
		} else {
			return fmt.Errorf("invalid port value: %s", envPort)
		}
	}

	if envURL := os.Getenv("MCP_VOICEVOX_URL"); envURL != "" {
		c.VoicevoxURL = envURL
	}

	if envTempDir := os.Getenv("MCP_VOICEVOX_TEMP_DIR"); envTempDir != "" {
		c.TempDir = envTempDir
	}

	if envSpeaker := os.Getenv("MCP_VOICEVOX_DEFAULT_SPEAKER"); envSpeaker != "" {
		if s, err := strconv.Atoi(envSpeaker); err == nil {
			c.DefaultSpeaker = s
		} else {
			return fmt.Errorf("invalid speaker ID value: %s", envSpeaker)
		}
	}

	if envPlayback := os.Getenv("MCP_VOICEVOX_ENABLE_PLAYBACK"); envPlayback != "" {
		c.EnablePlayback = envPlayback == "true"
	}

	// 音声合成パラメータの環境変数読み込み
	if envSpeedScale := os.Getenv("MCP_VOICEVOX_DEFAULT_SPEED_SCALE"); envSpeedScale != "" {
		if s, err := strconv.ParseFloat(envSpeedScale, 64); err == nil {
			if s >= 0.5 && s <= 2.0 {
				c.DefaultSpeedScale = s
			}
		} else {
			return fmt.Errorf("invalid speed scale value: %s", envSpeedScale)
		}
	}

	if envPitchScale := os.Getenv("MCP_VOICEVOX_DEFAULT_PITCH_SCALE"); envPitchScale != "" {
		if p, err := strconv.ParseFloat(envPitchScale, 64); err == nil {
			if p >= -0.15 && p <= 0.15 {
				c.DefaultPitchScale = p
			}
		} else {
			return fmt.Errorf("invalid pitch scale value: %s", envPitchScale)
		}
	}

	if envIntonationScale := os.Getenv("MCP_VOICEVOX_DEFAULT_INTONATION_SCALE"); envIntonationScale != "" {
		if i, err := strconv.ParseFloat(envIntonationScale, 64); err == nil {
			if i >= 0.0 && i <= 2.0 {
				c.DefaultIntonationScale = i
			}
		} else {
			return fmt.Errorf("invalid intonation scale value: %s", envIntonationScale)
		}
	}

	if envVolumeScale := os.Getenv("MCP_VOICEVOX_DEFAULT_VOLUME_SCALE"); envVolumeScale != "" {
		if v, err := strconv.ParseFloat(envVolumeScale, 64); err == nil {
			if v >= 0.0 && v <= 2.0 {
				c.DefaultVolumeScale = v
			}
		} else {
			return fmt.Errorf("invalid volume scale value: %s", envVolumeScale)
		}
	}

	return nil
}

// Validate は設定値の妥当性をチェックします
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	if c.VoicevoxURL == "" {
		return fmt.Errorf("voicevox URL cannot be empty")
	}

	if c.DefaultSpeaker < 0 {
		return fmt.Errorf("default speaker ID must be non-negative, got %d", c.DefaultSpeaker)
	}

	if c.TempDir == "" {
		return fmt.Errorf("temp directory cannot be empty")
	}

	// 音声合成パラメータの妥当性チェック
	if c.DefaultSpeedScale < 0.5 || c.DefaultSpeedScale > 2.0 {
		return fmt.Errorf("default speed scale must be between 0.5 and 2.0, got %f", c.DefaultSpeedScale)
	}

	if c.DefaultPitchScale < -0.15 || c.DefaultPitchScale > 0.15 {
		return fmt.Errorf("default pitch scale must be between -0.15 and 0.15, got %f", c.DefaultPitchScale)
	}

	if c.DefaultIntonationScale < 0.0 || c.DefaultIntonationScale > 2.0 {
		return fmt.Errorf("default intonation scale must be between 0.0 and 2.0, got %f", c.DefaultIntonationScale)
	}

	if c.DefaultVolumeScale < 0.0 || c.DefaultVolumeScale > 2.0 {
		return fmt.Errorf("default volume scale must be between 0.0 and 2.0, got %f", c.DefaultVolumeScale)
	}

	return nil
}

// SetupTempDir は一時ディレクトリを設定・作成します
func (c *Config) SetupTempDir() error {
	if c.TempDir == os.TempDir() {
		// デフォルトの場合は専用ディレクトリを作成
		c.TempDir = filepath.Join(os.TempDir(), "mcp-voicevox")
	}

	if _, err := os.Stat(c.TempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.TempDir, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	return nil
}

// New は新しい設定インスタンスを作成し、環境変数から読み込みます
func New() (*Config, error) {
	config := DefaultConfig()

	if err := config.LoadFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/metapox/mcp-voicevox-go/pkg/mcp"
)

func runHTTPServer() {
	loadEnvVars()
	setupTempDir()

	server := mcp.NewMCPServer(port, voicevoxURL, tempDir, defaultSpeaker)
	log.Printf("MCPサーバーを起動します: ポート %d, VOICEVOX URL: %s", port, voicevoxURL)
	log.Printf("一時ファイルディレクトリ: %s", tempDir)
	log.Printf("デフォルト話者ID: %d", defaultSpeaker)

	if err := server.Start(); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}

func loadEnvVars() {
	if envPort := os.Getenv("MCP_VOICEVOX_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	if envURL := os.Getenv("MCP_VOICEVOX_URL"); envURL != "" {
		voicevoxURL = envURL
	}
	if envTempDir := os.Getenv("MCP_VOICEVOX_TEMP_DIR"); envTempDir != "" {
		tempDir = envTempDir
	}
	if envSpeaker := os.Getenv("MCP_VOICEVOX_DEFAULT_SPEAKER"); envSpeaker != "" {
		if s, err := strconv.Atoi(envSpeaker); err == nil {
			defaultSpeaker = s
		}
	}
}

func setupTempDir() {
	if tempDir == "" {
		userTempDir := os.TempDir()
		tempDir = filepath.Join(userTempDir, "mcp-voicevox")
	}

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			log.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
		}
	}
}

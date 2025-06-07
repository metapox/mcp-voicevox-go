package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/metapox/mcp-voicevox-go/pkg/mcp"
	"github.com/spf13/cobra"
)

var (
	port           int
	voicevoxURL    string
	tempDir        string
	defaultSpeaker int
)

func main() {
	// ルートコマンドの定義
	var rootCmd = &cobra.Command{
		Use:   "mcp-voicevox",
		Short: "VOICEVOXを使用したMCPサーバー",
		Long:  `VOICEVOXのAPIを使用して、テキストを音声に変換するModel Context Protocol (MCP) サーバー`,
	}

	// サーバーコマンドの定義
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "MCPサーバーを起動",
		Long:  `VOICEVOXのAPIを使用したMCPサーバーを起動します`,
		Run: func(cmd *cobra.Command, args []string) {
			// 環境変数から設定を読み込み
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

			// 一時ディレクトリの作成
			if tempDir == "" {
				userTempDir := os.TempDir()
				tempDir = filepath.Join(userTempDir, "mcp-voicevox")
			}

			if _, err := os.Stat(tempDir); os.IsNotExist(err) {
				if err := os.MkdirAll(tempDir, 0755); err != nil {
					log.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
				}
			}

			// サーバーの作成と起動
			server := mcp.NewMCPServer(port, voicevoxURL, tempDir, defaultSpeaker)
			log.Printf("MCPサーバーを起動します: ポート %d, VOICEVOX URL: %s", port, voicevoxURL)
			log.Printf("一時ファイルディレクトリ: %s", tempDir)
			log.Printf("デフォルト話者ID: %d", defaultSpeaker)

			if err := server.Start(); err != nil {
				log.Fatalf("サーバーの起動に失敗しました: %v", err)
			}
		},
	}

	// フラグの設定
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "サーバーのポート番号")
	serverCmd.Flags().StringVarP(&voicevoxURL, "voicevox-url", "u", "http://localhost:50021", "VOICEVOXのAPIエンドポイント")
	serverCmd.Flags().StringVarP(&tempDir, "temp-dir", "t", "", "一時ファイルを保存するディレクトリ（デフォルト: システムの一時ディレクトリ）")
	serverCmd.Flags().IntVarP(&defaultSpeaker, "default-speaker", "s", 3, "デフォルトの話者ID")

	// コマンドの追加
	rootCmd.AddCommand(serverCmd)

	// コマンドの実行
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

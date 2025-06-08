package cmd

import (
	"log"

	"github.com/metapox/mcp-voicevox-go/pkg/config"
	"github.com/metapox/mcp-voicevox-go/pkg/mcp"
	"github.com/spf13/cobra"
)

var (
	port           int
	voicevoxURL    string
	tempDir        string
	defaultSpeaker int
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "HTTPサーバーとして起動",
	Long:  `HTTPサーバーとしてMCPサーバーを起動します`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHTTPServer(cmd)
	},
}

func init() {
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "サーバーのポート番号")
	serverCmd.Flags().StringVarP(&voicevoxURL, "voicevox-url", "u", "http://localhost:50021", "VOICEVOXのAPIエンドポイント")
	serverCmd.Flags().StringVarP(&tempDir, "temp-dir", "t", "", "一時ファイルを保存するディレクトリ")
	serverCmd.Flags().IntVarP(&defaultSpeaker, "default-speaker", "s", 3, "デフォルトの話者ID")
}

func runHTTPServer(cmd *cobra.Command) error {
	// 設定を作成
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// コマンドラインフラグで設定を上書き
	if cmd.Flags().Changed("port") {
		cfg.Port = port
	}
	if cmd.Flags().Changed("voicevox-url") {
		cfg.VoicevoxURL = voicevoxURL
	}
	if cmd.Flags().Changed("temp-dir") {
		cfg.TempDir = tempDir
	}
	if cmd.Flags().Changed("default-speaker") {
		cfg.DefaultSpeaker = defaultSpeaker
	}

	// 一時ディレクトリのセットアップ
	if err := cfg.SetupTempDir(); err != nil {
		return err
	}

	// サーバー起動
	server := mcp.NewMCPServer(cfg.Port, cfg.VoicevoxURL, cfg.TempDir, cfg.DefaultSpeaker)
	log.Printf("MCPサーバーを起動します: ポート %d, VOICEVOX URL: %s", cfg.Port, cfg.VoicevoxURL)
	log.Printf("一時ファイルディレクトリ: %s", cfg.TempDir)
	log.Printf("デフォルト話者ID: %d", cfg.DefaultSpeaker)

	return server.Start()
}

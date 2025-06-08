package cmd

import (
	"log"

	"github.com/metapox/mcp-voicevox-go/pkg/config"
	"github.com/metapox/mcp-voicevox-go/pkg/mcp"
	"github.com/spf13/cobra"
)

var (
	port                   int
	voicevoxURL            string
	tempDir                string
	defaultSpeaker         int
	defaultSpeedScale      float64
	defaultPitchScale      float64
	defaultIntonationScale float64
	defaultVolumeScale     float64
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
	serverCmd.Flags().Float64Var(&defaultSpeedScale, "default-speed-scale", 1.0, "デフォルトの話速（0.5-2.0）")
	serverCmd.Flags().Float64Var(&defaultPitchScale, "default-pitch-scale", 0.0, "デフォルトの音高（-0.15-0.15）")
	serverCmd.Flags().Float64Var(&defaultIntonationScale, "default-intonation-scale", 1.0, "デフォルトの抑揚（0.0-2.0）")
	serverCmd.Flags().Float64Var(&defaultVolumeScale, "default-volume-scale", 1.0, "デフォルトの音量（0.0-2.0）")
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
	if cmd.Flags().Changed("default-speed-scale") {
		cfg.DefaultSpeedScale = defaultSpeedScale
	}
	if cmd.Flags().Changed("default-pitch-scale") {
		cfg.DefaultPitchScale = defaultPitchScale
	}
	if cmd.Flags().Changed("default-intonation-scale") {
		cfg.DefaultIntonationScale = defaultIntonationScale
	}
	if cmd.Flags().Changed("default-volume-scale") {
		cfg.DefaultVolumeScale = defaultVolumeScale
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
	log.Printf("デフォルト音声設定: 話速=%.2f, 音高=%.2f, 抑揚=%.2f, 音量=%.2f", 
		cfg.DefaultSpeedScale, cfg.DefaultPitchScale, cfg.DefaultIntonationScale, cfg.DefaultVolumeScale)

	return server.Start()
}

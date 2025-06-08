package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/metapox/mcp-voicevox-go/pkg/config"
	"github.com/metapox/mcp-voicevox-go/pkg/mcp"
	"github.com/spf13/cobra"
)

var enablePlayback bool

var stdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Stdio経由で起動（MCP用）",
	Long:  `標準入出力経由でMCPサーバーを起動します`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStdioServer(cmd)
	},
}

func init() {
	stdioCmd.Flags().StringVarP(&voicevoxURL, "voicevox-url", "u", "http://localhost:50021", "VOICEVOXのAPIエンドポイント")
	stdioCmd.Flags().StringVarP(&tempDir, "temp-dir", "t", "", "一時ファイルを保存するディレクトリ")
	stdioCmd.Flags().IntVarP(&defaultSpeaker, "default-speaker", "s", 3, "デフォルトの話者ID")
	stdioCmd.Flags().BoolVar(&enablePlayback, "enable-playback", false, "音声の自動再生を有効にする")
}

func runStdioServer(cmd *cobra.Command) error {
	// 設定を作成
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// コマンドラインフラグで設定を上書き
	if cmd.Flags().Changed("voicevox-url") {
		cfg.VoicevoxURL = voicevoxURL
	}
	if cmd.Flags().Changed("temp-dir") {
		cfg.TempDir = tempDir
	}
	if cmd.Flags().Changed("default-speaker") {
		cfg.DefaultSpeaker = defaultSpeaker
	}
	if cmd.Flags().Changed("enable-playback") {
		cfg.EnablePlayback = enablePlayback
	}

	// 一時ディレクトリのセットアップ
	if err := cfg.SetupTempDir(); err != nil {
		return err
	}

	// MCPハンドラーを作成
	handler := mcp.NewHandler(cfg)

	scanner := bufio.NewScanner(os.Stdin)

	log.SetOutput(os.Stderr)
	log.Printf("MCP Stdio Server started - VOICEVOX URL: %s, Default Speaker: %d, Playback: %v",
		cfg.VoicevoxURL, cfg.DefaultSpeaker, cfg.EnablePlayback)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req mcp.MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

		response := handler.HandleRequest(req)
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
	}

	return scanner.Err()
}

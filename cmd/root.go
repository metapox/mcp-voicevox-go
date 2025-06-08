package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	voicevoxURL     string
	tempDir         string
	defaultSpeaker  int
	port            int
	enablePlayback  bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mcp-voicevox",
		Short: "VOICEVOXを使用したMCPサーバー",
		Long:  `VOICEVOXのAPIを使用して、テキストを音声に変換するModel Context Protocol (MCP) サーバー`,
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "HTTPサーバーとしてMCPサーバーを起動",
		Long:  `HTTPサーバーとしてVOICEVOXのAPIを使用したMCPサーバーを起動します`,
		Run: func(cmd *cobra.Command, args []string) {
			runHTTPServer()
		},
	}

	var stdioCmd = &cobra.Command{
		Use:   "stdio",
		Short: "Stdio経由でMCPサーバーを起動",
		Long:  `標準入出力経由でVOICEVOXのAPIを使用したMCPサーバーを起動します`,
		Run: func(cmd *cobra.Command, args []string) {
			runStdioServer()
		},
	}

	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "サーバーのポート番号")
	serverCmd.Flags().StringVarP(&voicevoxURL, "voicevox-url", "u", "http://localhost:50021", "VOICEVOXのAPIエンドポイント")
	serverCmd.Flags().StringVarP(&tempDir, "temp-dir", "t", "", "一時ファイルを保存するディレクトリ（デフォルト: システムの一時ディレクトリ）")
	serverCmd.Flags().IntVarP(&defaultSpeaker, "default-speaker", "s", 3, "デフォルトの話者ID")

	stdioCmd.Flags().StringVarP(&voicevoxURL, "voicevox-url", "u", "http://localhost:50021", "VOICEVOXのAPIエンドポイント")
	stdioCmd.Flags().StringVarP(&tempDir, "temp-dir", "t", "", "一時ファイルを保存するディレクトリ")
	stdioCmd.Flags().IntVarP(&defaultSpeaker, "default-speaker", "s", 3, "デフォルトの話者ID")
	stdioCmd.Flags().BoolVar(&enablePlayback, "enable-playback", false, "音声の自動再生を有効にする")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(stdioCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

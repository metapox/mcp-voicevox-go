package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-voicevox",
	Short: "VOICEVOX MCP Server",
	Long:  `VOICEVOX Model Context Protocol Server - テキストを音声に変換するMCPサーバー`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(stdioCmd)
}

// Execute はrootコマンドを実行します
func Execute() error {
	return rootCmd.Execute()
}

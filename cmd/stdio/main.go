package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/metapox/mcp-voicevox-go/pkg/voicevox"
	"github.com/spf13/cobra"
)

var (
	voicevoxURL    string
	tempDir        string
	defaultSpeaker int
)

// MCP Protocol structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]interface{} `json:"clientInfo"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mcp-voicevox-stdio",
		Short: "VOICEVOX MCP Server (Stdio)",
		Long:  `VOICEVOXのAPIを使用したModel Context Protocol (MCP) サーバー (Stdio版)`,
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "MCPサーバーを起動 (Stdio)",
		Long:  `Stdio経由でMCPサーバーを起動します`,
		Run: func(cmd *cobra.Command, args []string) {
			runStdioServer()
		},
	}

	serverCmd.Flags().StringVarP(&voicevoxURL, "voicevox-url", "u", "http://localhost:50021", "VOICEVOXのAPIエンドポイント")
	serverCmd.Flags().StringVarP(&tempDir, "temp-dir", "t", "", "一時ファイルを保存するディレクトリ")
	serverCmd.Flags().IntVarP(&defaultSpeaker, "default-speaker", "s", 3, "デフォルトの話者ID")

	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runStdioServer() {
	// 環境変数から設定を読み込み
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

	if tempDir == "" {
		tempDir = os.TempDir()
	}

	client := voicevox.NewClient(voicevoxURL)
	scanner := bufio.NewScanner(os.Stdin)

	log.SetOutput(os.Stderr) // ログは標準エラー出力に
	log.Printf("MCP Stdio Server started - VOICEVOX URL: %s, Default Speaker: %d", voicevoxURL, defaultSpeaker)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

		response := handleRequest(req, client)
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
	}
}

func handleRequest(req MCPRequest, client *voicevox.Client) MCPResponse {
	switch req.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "mcp-voicevox-go",
					"version": "1.0.0",
				},
			},
		}

	case "tools/list":
		tools := []Tool{
			{
				Name:        "text_to_speech",
				Description: "テキストを音声に変換します",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "音声に変換するテキスト",
						},
						"speaker_id": map[string]interface{}{
							"type":        "integer",
							"description": "話者ID（省略時はデフォルト話者を使用）",
						},
					},
					"required": []string{"text"},
				},
			},
			{
				Name:        "get_speakers",
				Description: "利用可能な話者一覧を取得します",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": tools,
			},
		}

	case "tools/call":
		var params CallToolParams
		if paramBytes, err := json.Marshal(req.Params); err == nil {
			json.Unmarshal(paramBytes, &params)
		}

		switch params.Name {
		case "text_to_speech":
			return handleTextToSpeech(req.ID, params.Arguments, client)
		case "get_speakers":
			return handleGetSpeakers(req.ID, client)
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Unknown tool: " + params.Name,
			},
		}

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found: " + req.Method,
			},
		}
	}
}

func handleTextToSpeech(id interface{}, args map[string]interface{}, client *voicevox.Client) MCPResponse {
	text, ok := args["text"].(string)
	if !ok {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32602,
				Message: "text parameter is required",
			},
		}
	}

	speakerID := defaultSpeaker
	if sid, ok := args["speaker_id"].(float64); ok {
		speakerID = int(sid)
	}

	// 音声クエリを作成
	query, err := client.CreateAudioQuery(text, speakerID)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32603,
				Message: "Failed to create audio query: " + err.Error(),
			},
		}
	}

	// 音声合成を実行
	audioData, err := client.SynthesizeVoice(query, speakerID)
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32603,
				Message: "Text to speech failed: " + err.Error(),
			},
		}
	}

	// 音声ファイルを一時ディレクトリに保存
	filename := fmt.Sprintf("speech_%d_%d.wav", speakerID, time.Now().Unix())
	filepath := fmt.Sprintf("%s/%s", tempDir, filename)
	
	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32603,
				Message: "Failed to save audio file: " + err.Error(),
			},
		}
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("音声合成が完了しました。\nテキスト: %s\n話者ID: %d\nファイル: %s", text, speakerID, filepath),
				},
			},
		},
	}
}

func handleGetSpeakers(id interface{}, client *voicevox.Client) MCPResponse {
	speakers, err := client.GetSpeakers()
	if err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      id,
			Error: &MCPError{
				Code:    -32603,
				Message: "Failed to get speakers: " + err.Error(),
			},
		}
	}

	speakersJSON, _ := json.MarshalIndent(speakers, "", "  ")

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("利用可能な話者一覧:\n%s", string(speakersJSON)),
				},
			},
		},
	}
}

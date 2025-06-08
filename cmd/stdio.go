package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/metapox/mcp-voicevox-go/pkg/audio"
	"github.com/metapox/mcp-voicevox-go/pkg/voicevox"
)

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

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

func runStdioServer() {
	loadStdioEnvVars()
	if tempDir == "" {
		tempDir = os.TempDir()
	}

	client := voicevox.NewClient(voicevoxURL)
	enablePlayback := os.Getenv("MCP_VOICEVOX_ENABLE_PLAYBACK") == "true"
	player := audio.NewPlayer(enablePlayback)
	
	scanner := bufio.NewScanner(os.Stdin)

	log.SetOutput(os.Stderr)
	log.Printf("MCP Stdio Server started - VOICEVOX URL: %s, Default Speaker: %d, Playback: %v", voicevoxURL, defaultSpeaker, enablePlayback)

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

		response := handleRequest(req, client, player)
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
	}
}

func loadStdioEnvVars() {
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

func handleRequest(req MCPRequest, client *voicevox.Client, player *audio.Player) MCPResponse {
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
			return handleTextToSpeech(req.ID, params.Arguments, client, player)
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

func handleTextToSpeech(id interface{}, args map[string]interface{}, client *voicevox.Client, player *audio.Player) MCPResponse {
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

	playbackStatus := "ファイルに保存されました"
	if err := player.Play(filepath); err != nil {
		log.Printf("Audio playback failed: %v", err)
		playbackStatus = fmt.Sprintf("ファイルに保存されましたが、再生に失敗しました: %v", err)
	} else if player != nil {
		playbackStatus = "ファイルに保存され、音声を再生しました"
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("音声合成が完了しました。\nテキスト: %s\n話者ID: %d\nファイル: %s\n状態: %s", text, speakerID, filepath, playbackStatus),
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

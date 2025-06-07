package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/minamitakumi/mcp-voicevox-go/pkg/voicevox"
	"github.com/rs/cors"
)

// MCPServer はModel Context Protocol サーバーの構造体です
type MCPServer struct {
	Port           int
	VoicevoxURL    string
	VoicevoxClient *voicevox.Client
	TempDir        string
	DefaultSpeaker int
}

// NewMCPServer は新しいMCPサーバーを作成します
func NewMCPServer(port int, voicevoxURL string, tempDir string, defaultSpeaker int) *MCPServer {
	return &MCPServer{
		Port:           port,
		VoicevoxURL:    voicevoxURL,
		VoicevoxClient: voicevox.NewClient(voicevoxURL),
		TempDir:        tempDir,
		DefaultSpeaker: defaultSpeaker,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Start はMCPサーバーを起動します
func (s *MCPServer) Start() error {
	// CORSの設定
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// ルーティングの設定
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/manifest", s.handleManifest)
	mux.HandleFunc("/health", s.handleHealth)

	// サーバーの起動
	handler := c.Handler(mux)
	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("Starting MCP server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

// handleWebSocket はWebSocket接続を処理します
func (s *MCPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var request map[string]interface{}
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

		// リクエストの処理
		response, err := s.handleRequest(request)
		if err != nil {
			log.Printf("Request handling error: %v", err)
			errorResponse := map[string]interface{}{
				"id":    request["id"],
				"error": err.Error(),
			}
			responseJSON, _ := json.Marshal(errorResponse)
			if err := conn.WriteMessage(messageType, responseJSON); err != nil {
				log.Printf("WebSocket write error: %v", err)
				break
			}
			continue
		}

		// レスポンスの送信
		responseJSON, _ := json.Marshal(response)
		if err := conn.WriteMessage(messageType, responseJSON); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

// handleRequest はMCPリクエストを処理します
func (s *MCPServer) handleRequest(request map[string]interface{}) (map[string]interface{}, error) {
	requestID, _ := request["id"].(string)
	method, _ := request["method"].(string)

	switch method {
	case "invoke":
		return s.handleInvoke(requestID, request)
	case "discover":
		return s.handleDiscover(requestID)
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

// handleInvoke はツール呼び出しを処理します
func (s *MCPServer) handleInvoke(requestID string, request map[string]interface{}) (map[string]interface{}, error) {
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}

	toolName, _ := params["name"].(string)
	toolParams, _ := params["parameters"].(map[string]interface{})

	switch toolName {
	case "text_to_speech":
		return s.handleTextToSpeech(requestID, toolParams)
	case "get_speakers":
		return s.handleGetSpeakers(requestID)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

// handleTextToSpeech はテキスト音声合成リクエストを処理します
func (s *MCPServer) handleTextToSpeech(requestID string, params map[string]interface{}) (map[string]interface{}, error) {
	text, _ := params["text"].(string)
	speakerIDFloat, ok := params["speaker_id"].(float64)
	var speakerID int
	if ok {
		speakerID = int(speakerIDFloat)
	} else {
		speakerID = s.DefaultSpeaker // デフォルト話者IDを使用
	}

	if text == "" {
		return nil, fmt.Errorf("text parameter is required")
	}

	// 音声合成クエリの作成
	query, err := s.VoicevoxClient.CreateAudioQuery(text, speakerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio query: %v", err)
	}

	// 音声合成の実行
	audioData, err := s.VoicevoxClient.SynthesizeVoice(query, speakerID)
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize voice: %v", err)
	}

	// 一時ファイルの作成
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("voicevox_%d.wav", timestamp)
	filepath := filepath.Join(s.TempDir, filename)

	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write audio file: %v", err)
	}

	return map[string]interface{}{
		"id": requestID,
		"result": map[string]interface{}{
			"audio_path": filepath,
			"text":       text,
			"speaker_id": speakerID,
		},
	}, nil
}

// handleGetSpeakers は利用可能な話者一覧を取得します
func (s *MCPServer) handleGetSpeakers(requestID string) (map[string]interface{}, error) {
	speakers, err := s.VoicevoxClient.GetSpeakers()
	if err != nil {
		return nil, fmt.Errorf("failed to get speakers: %v", err)
	}

	return map[string]interface{}{
		"id": requestID,
		"result": map[string]interface{}{
			"speakers": speakers,
		},
	}, nil
}

// handleDiscover はツール一覧を返します
func (s *MCPServer) handleDiscover(requestID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"id": requestID,
		"result": map[string]interface{}{
			"name":        "voicevox",
			"display_name": "VOICEVOX",
			"description": "VOICEVOXを使用して日本語テキストを音声に変換するツール",
			"tools": []map[string]interface{}{
				{
					"name":        "text_to_speech",
					"description": "テキストを音声に変換します",
					"parameters": map[string]interface{}{
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
					"name":        "get_speakers",
					"description": "利用可能な話者一覧を取得します",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{},
					},
				},
			},
		},
	}, nil
}

// handleManifest はマニフェスト情報を返します
func (s *MCPServer) handleManifest(w http.ResponseWriter, r *http.Request) {
	manifest := map[string]interface{}{
		"name":        "voicevox",
		"display_name": "VOICEVOX",
		"description": "VOICEVOXを使用して日本語テキストを音声に変換するツール",
		"version":     "1.0.0",
		"api_version": "1",
		"auth": map[string]interface{}{
			"type": "none",
		},
		"endpoints": map[string]interface{}{
			"ws":      fmt.Sprintf("ws://localhost:%d/ws", s.Port),
			"health":  fmt.Sprintf("http://localhost:%d/health", s.Port),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manifest)
}

// handleHealth はヘルスチェックエンドポイントを処理します
func (s *MCPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

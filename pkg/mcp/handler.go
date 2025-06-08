package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/metapox/mcp-voicevox-go/pkg/audio"
	"github.com/metapox/mcp-voicevox-go/pkg/config"
	"github.com/metapox/mcp-voicevox-go/pkg/errors"
	"github.com/metapox/mcp-voicevox-go/pkg/voicevox"
)

// Handler はMCPプロトコルのハンドラーです
type Handler struct {
	config         *config.Config
	voicevoxClient *voicevox.Client
	audioPlayer    *audio.Player
}

// NewHandler は新しいMCPハンドラーを作成します
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config:         cfg,
		voicevoxClient: voicevox.NewClient(cfg.VoicevoxURL),
		audioPlayer:    audio.NewPlayer(cfg.EnablePlayback),
	}
}

// HandleRequest はMCPリクエストを処理します
func (h *Handler) HandleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case MethodInitialize:
		return h.handleInitialize(req.ID)
	case MethodToolsList:
		return h.handleToolsList(req.ID)
	case MethodToolsCall:
		return h.handleToolsCall(req.ID, req.Params)
	default:
		return h.createErrorResponse(req.ID, errors.NewMCPError(errors.MCPMethodNotFound, "Method not found: "+req.Method))
	}
}

// handleInitialize は初期化リクエストを処理します
func (h *Handler) handleInitialize(id interface{}) MCPResponse {
	result := InitializeResult{
		ProtocolVersion: ProtocolVersion,
		Capabilities: map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		ServerInfo: ServerInfo{
			Name:    ServerName,
			Version: ServerVersion,
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// handleToolsList はツール一覧リクエストを処理します
func (h *Handler) handleToolsList(id interface{}) MCPResponse {
	tools := []Tool{
		{
			Name:        ToolTextToSpeech,
			Description: "テキストを音声に変換します",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "音声に変換するテキスト",
						"maxLength":   1000,
					},
					"speaker_id": map[string]interface{}{
						"type":        "integer",
						"description": "話者ID（省略時はデフォルト話者を使用）",
						"minimum":     0,
					},
					"speed_scale": map[string]interface{}{
						"type":        "number",
						"description": "話速（0.5-2.0、デフォルト: 1.0）",
						"minimum":     0.5,
						"maximum":     2.0,
					},
					"pitch_scale": map[string]interface{}{
						"type":        "number",
						"description": "音高（-0.15-0.15、デフォルト: 0.0）",
						"minimum":     -0.15,
						"maximum":     0.15,
					},
					"intonation_scale": map[string]interface{}{
						"type":        "number",
						"description": "抑揚（0.0-2.0、デフォルト: 1.0）",
						"minimum":     0.0,
						"maximum":     2.0,
					},
					"volume_scale": map[string]interface{}{
						"type":        "number",
						"description": "音量（0.0-2.0、デフォルト: 1.0）",
						"minimum":     0.0,
						"maximum":     2.0,
					},
				},
				"required": []string{"text"},
			},
		},
		{
			Name:        ToolGetSpeakers,
			Description: "利用可能な話者一覧を取得します",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	result := ToolsListResult{Tools: tools}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// handleToolsCall はツール呼び出しリクエストを処理します
func (h *Handler) handleToolsCall(id interface{}, params interface{}) MCPResponse {
	var callParams CallToolParams
	if paramBytes, err := json.Marshal(params); err == nil {
		json.Unmarshal(paramBytes, &callParams)
	}

	switch callParams.Name {
	case ToolTextToSpeech:
		return h.handleTextToSpeech(id, callParams.Arguments)
	case ToolGetSpeakers:
		return h.handleGetSpeakers(id)
	default:
		return h.createErrorResponse(id, errors.NewMCPError(errors.MCPInvalidParams, "Unknown tool: "+callParams.Name))
	}
}

// handleTextToSpeech はテキスト音声変換を処理します
func (h *Handler) handleTextToSpeech(id interface{}, args map[string]interface{}) MCPResponse {
	text, ok := args["text"].(string)
	if !ok {
		return h.createErrorResponse(id, errors.NewMCPError(errors.MCPInvalidParams, "text parameter is required"))
	}

	speakerID := h.config.DefaultSpeaker
	if sid, ok := args["speaker_id"].(float64); ok {
		speakerID = int(sid)
	}

	// 音声合成オプションを準備
	var options *voicevox.AudioQueryOptions
	hasOptions := args["speed_scale"] != nil || args["pitch_scale"] != nil || 
		args["intonation_scale"] != nil || args["volume_scale"] != nil
	
	// デフォルト設定またはパラメータ指定の値を使用
	options = &voicevox.AudioQueryOptions{
		SpeedScale:      &h.config.DefaultSpeedScale,
		PitchScale:      &h.config.DefaultPitchScale,
		IntonationScale: &h.config.DefaultIntonationScale,
		VolumeScale:     &h.config.DefaultVolumeScale,
	}
	
	// パラメータで上書き
	if hasOptions {
		if speedScale, ok := args["speed_scale"].(float64); ok {
			if speedScale >= 0.5 && speedScale <= 2.0 {
				options.SpeedScale = &speedScale
			}
		}
		
		if pitchScale, ok := args["pitch_scale"].(float64); ok {
			if pitchScale >= -0.15 && pitchScale <= 0.15 {
				options.PitchScale = &pitchScale
			}
		}
		
		if intonationScale, ok := args["intonation_scale"].(float64); ok {
			if intonationScale >= 0.0 && intonationScale <= 2.0 {
				options.IntonationScale = &intonationScale
			}
		}
		
		if volumeScale, ok := args["volume_scale"].(float64); ok {
			if volumeScale >= 0.0 && volumeScale <= 2.0 {
				options.VolumeScale = &volumeScale
			}
		}
	}

	// 音声クエリ作成
	query, err := h.voicevoxClient.CreateAudioQueryWithOptions(text, speakerID, options)
	if err != nil {
		appErr := errors.NewVoicevoxError("Failed to create audio query", err)
		return h.createErrorResponse(id, appErr)
	}

	// 音声合成
	audioData, err := h.voicevoxClient.SynthesizeVoice(query, speakerID)
	if err != nil {
		appErr := errors.NewAudioSynthesisError("Text to speech failed", err)
		return h.createErrorResponse(id, appErr)
	}

	// ファイル保存
	filename := fmt.Sprintf("speech_%d_%d.wav", speakerID, time.Now().Unix())
	filepath := fmt.Sprintf("%s/%s", h.config.TempDir, filename)

	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		appErr := errors.NewFileOperationError("Failed to save audio file", err)
		return h.createErrorResponse(id, appErr)
	}

	// 音声再生
	playbackStatus := "ファイルに保存されました"
	if err := h.audioPlayer.Play(filepath); err != nil {
		log.Printf("Audio playback failed: %v", err)
		playbackStatus = fmt.Sprintf("ファイルに保存されましたが、再生に失敗しました: %v", err)
	} else if h.audioPlayer != nil && h.config.EnablePlayback {
		playbackStatus = "ファイルに保存され、音声を再生しました"
	}

	// オプション情報を含む結果メッセージ
	optionsInfo := ""
	if options != nil {
		optionsInfo += fmt.Sprintf("\n話速: %.2f", *options.SpeedScale)
		optionsInfo += fmt.Sprintf("\n音高: %.2f", *options.PitchScale)
		optionsInfo += fmt.Sprintf("\n抑揚: %.2f", *options.IntonationScale)
		optionsInfo += fmt.Sprintf("\n音量: %.2f", *options.VolumeScale)
	}

	result := ToolCallResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: fmt.Sprintf("音声合成が完了しました。\nテキスト: %s\n話者ID: %d%s\nファイル: %s\n状態: %s",
					text, speakerID, optionsInfo, filepath, playbackStatus),
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// handleGetSpeakers は話者一覧取得を処理します
func (h *Handler) handleGetSpeakers(id interface{}) MCPResponse {
	speakers, err := h.voicevoxClient.GetSpeakers()
	if err != nil {
		appErr := errors.NewVoicevoxError("Failed to get speakers", err)
		return h.createErrorResponse(id, appErr)
	}

	speakersJSON, _ := json.MarshalIndent(speakers, "", "  ")

	result := ToolCallResult{
		Content: []ContentItem{
			{
				Type: "text",
				Text: fmt.Sprintf("利用可能な話者一覧:\n%s", string(speakersJSON)),
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// createErrorResponse はエラーレスポンスを作成します
func (h *Handler) createErrorResponse(id interface{}, appErr *errors.AppError) MCPResponse {
	mcpErr := appErr.ToMCPError()
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    mcpErr.Code,
			Message: mcpErr.Message,
		},
	}
}

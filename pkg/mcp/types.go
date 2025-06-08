package mcp

// MCPRequest はMCPプロトコルのリクエスト構造体です
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse はMCPプロトコルのレスポンス構造体です
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError はMCPプロトコルのエラー構造体です
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool はMCPツールの定義構造体です
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallToolParams はツール呼び出しのパラメータ構造体です
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// InitializeResult は初期化レスポンスの結果構造体です
type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

// ServerInfo はサーバー情報構造体です
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ToolsListResult はツール一覧レスポンスの結果構造体です
type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

// ToolCallResult はツール呼び出し結果構造体です
type ToolCallResult struct {
	Content []ContentItem `json:"content"`
}

// ContentItem はコンテンツアイテム構造体です
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// 定数定義
const (
	ProtocolVersion = "2024-11-05"
	ServerName      = "mcp-voicevox-go"
	ServerVersion   = "1.0.0"
)

// MCPメソッド名の定数
const (
	MethodInitialize = "initialize"
	MethodToolsList  = "tools/list"
	MethodToolsCall  = "tools/call"
)

// ツール名の定数
const (
	ToolTextToSpeech = "text_to_speech"
	ToolGetSpeakers  = "get_speakers"
)

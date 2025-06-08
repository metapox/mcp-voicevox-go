package errors

import "fmt"

// ErrorCode はアプリケーション固有のエラーコードを定義します
type ErrorCode int

const (
	// MCP関連エラー
	MCPParseError     ErrorCode = -32700 // JSON-RPC Parse error
	MCPInvalidRequest ErrorCode = -32600 // JSON-RPC Invalid Request
	MCPMethodNotFound ErrorCode = -32601 // JSON-RPC Method not found
	MCPInvalidParams  ErrorCode = -32602 // JSON-RPC Invalid params
	MCPInternalError  ErrorCode = -32603 // JSON-RPC Internal error

	// アプリケーション固有エラー
	VoicevoxConnectionError ErrorCode = -40001 // VOICEVOX接続エラー
	AudioSynthesisError     ErrorCode = -40002 // 音声合成エラー
	AudioPlaybackError      ErrorCode = -40003 // 音声再生エラー
	FileOperationError      ErrorCode = -40004 // ファイル操作エラー
	ConfigurationError      ErrorCode = -40005 // 設定エラー
)

// AppError はアプリケーション固有のエラー型です
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error はerrorインターフェースを実装します
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap はerrors.Unwrapに対応します
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError は新しいAppErrorを作成します
func NewAppError(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// MCPError はMCPプロトコル用のエラー構造体です
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToMCPError はAppErrorをMCPErrorに変換します
func (e *AppError) ToMCPError() *MCPError {
	return &MCPError{
		Code:    int(e.Code),
		Message: e.Message,
	}
}

// 便利な関数群

// NewVoicevoxError はVOICEVOX関連のエラーを作成します
func NewVoicevoxError(message string, cause error) *AppError {
	return NewAppError(VoicevoxConnectionError, message, cause)
}

// NewAudioSynthesisError は音声合成エラーを作成します
func NewAudioSynthesisError(message string, cause error) *AppError {
	return NewAppError(AudioSynthesisError, message, cause)
}

// NewAudioPlaybackError は音声再生エラーを作成します
func NewAudioPlaybackError(message string, cause error) *AppError {
	return NewAppError(AudioPlaybackError, message, cause)
}

// NewFileOperationError はファイル操作エラーを作成します
func NewFileOperationError(message string, cause error) *AppError {
	return NewAppError(FileOperationError, message, cause)
}

// NewConfigurationError は設定エラーを作成します
func NewConfigurationError(message string, cause error) *AppError {
	return NewAppError(ConfigurationError, message, cause)
}

// NewMCPError はMCPプロトコルエラーを作成します
func NewMCPError(code ErrorCode, message string) *AppError {
	return NewAppError(code, message, nil)
}

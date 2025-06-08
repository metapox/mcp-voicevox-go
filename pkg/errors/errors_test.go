package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		want     string
	}{
		{
			name: "error without cause",
			appError: &AppError{
				Code:    VoicevoxConnectionError,
				Message: "Connection failed",
				Cause:   nil,
			},
			want: "[-40001] Connection failed",
		},
		{
			name: "error with cause",
			appError: &AppError{
				Code:    AudioSynthesisError,
				Message: "Synthesis failed",
				Cause:   errors.New("network error"),
			},
			want: "[-40002] Synthesis failed: network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appError.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	appErr := &AppError{
		Code:    VoicevoxConnectionError,
		Message: "Connection failed",
		Cause:   cause,
	}

	if got := appErr.Unwrap(); got != cause {
		t.Errorf("AppError.Unwrap() = %v, want %v", got, cause)
	}

	// Test without cause
	appErrNoCause := &AppError{
		Code:    VoicevoxConnectionError,
		Message: "Connection failed",
		Cause:   nil,
	}

	if got := appErrNoCause.Unwrap(); got != nil {
		t.Errorf("AppError.Unwrap() = %v, want nil", got)
	}
}

func TestAppError_ToMCPError(t *testing.T) {
	appErr := &AppError{
		Code:    VoicevoxConnectionError,
		Message: "Connection failed",
		Cause:   nil,
	}

	mcpErr := appErr.ToMCPError()

	if mcpErr.Code != int(VoicevoxConnectionError) {
		t.Errorf("ToMCPError().Code = %v, want %v", mcpErr.Code, int(VoicevoxConnectionError))
	}

	if mcpErr.Message != "Connection failed" {
		t.Errorf("ToMCPError().Message = %v, want %v", mcpErr.Message, "Connection failed")
	}
}

func TestNewAppError(t *testing.T) {
	cause := errors.New("original error")
	appErr := NewAppError(VoicevoxConnectionError, "Connection failed", cause)

	if appErr.Code != VoicevoxConnectionError {
		t.Errorf("NewAppError().Code = %v, want %v", appErr.Code, VoicevoxConnectionError)
	}

	if appErr.Message != "Connection failed" {
		t.Errorf("NewAppError().Message = %v, want %v", appErr.Message, "Connection failed")
	}

	if appErr.Cause != cause {
		t.Errorf("NewAppError().Cause = %v, want %v", appErr.Cause, cause)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	cause := errors.New("test error")

	tests := []struct {
		name     string
		fn       func(string, error) *AppError
		wantCode ErrorCode
	}{
		{"NewVoicevoxError", NewVoicevoxError, VoicevoxConnectionError},
		{"NewAudioSynthesisError", NewAudioSynthesisError, AudioSynthesisError},
		{"NewAudioPlaybackError", NewAudioPlaybackError, AudioPlaybackError},
		{"NewFileOperationError", NewFileOperationError, FileOperationError},
		{"NewConfigurationError", NewConfigurationError, ConfigurationError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := tt.fn("test message", cause)
			if appErr.Code != tt.wantCode {
				t.Errorf("%s().Code = %v, want %v", tt.name, appErr.Code, tt.wantCode)
			}
			if appErr.Message != "test message" {
				t.Errorf("%s().Message = %v, want %v", tt.name, appErr.Message, "test message")
			}
			if appErr.Cause != cause {
				t.Errorf("%s().Cause = %v, want %v", tt.name, appErr.Cause, cause)
			}
		})
	}
}

func TestNewMCPError(t *testing.T) {
	appErr := NewMCPError(MCPInvalidParams, "Invalid parameters")

	if appErr.Code != MCPInvalidParams {
		t.Errorf("NewMCPError().Code = %v, want %v", appErr.Code, MCPInvalidParams)
	}

	if appErr.Message != "Invalid parameters" {
		t.Errorf("NewMCPError().Message = %v, want %v", appErr.Message, "Invalid parameters")
	}

	if appErr.Cause != nil {
		t.Errorf("NewMCPError().Cause = %v, want nil", appErr.Cause)
	}
}

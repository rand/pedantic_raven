package mnemosyne

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- GetErrorMessage Tests ---

func TestGetErrorMessageNil(t *testing.T) {
	msg := GetErrorMessage(nil)
	if msg != "" {
		t.Errorf("GetErrorMessage(nil) should return empty string, got %q", msg)
	}
}

func TestGetErrorMessageConnection(t *testing.T) {
	err := ErrConnection
	msg := GetErrorMessage(err)

	expected := "Cannot connect to mnemosyne server. Working offline with cached data."
	if msg != expected {
		t.Errorf("GetErrorMessage(connection) = %q, want %q", msg, expected)
	}
}

func TestGetErrorMessageServer(t *testing.T) {
	err := ErrInternal
	msg := GetErrorMessage(err)

	expected := "Server error occurred. Changes will be retried automatically."
	if msg != expected {
		t.Errorf("GetErrorMessage(server) = %q, want %q", msg, expected)
	}
}

func TestGetErrorMessageValidation(t *testing.T) {
	err := ErrInvalidArgument
	msg := GetErrorMessage(err)

	if !contains(msg, "Invalid data") {
		t.Errorf("GetErrorMessage(validation) should contain 'Invalid data', got %q", msg)
	}
}

func TestGetErrorMessageTimeout(t *testing.T) {
	err := ErrTimeout
	msg := GetErrorMessage(err)

	expected := "Operation timed out. Check your connection and try again."
	if msg != expected {
		t.Errorf("GetErrorMessage(timeout) = %q, want %q", msg, expected)
	}
}

func TestGetErrorMessageUnknown(t *testing.T) {
	err := errors.New("unknown error")
	msg := GetErrorMessage(err)

	// Unknown errors default to server category, which gives server error message
	if !contains(msg, "Server error") && !contains(msg, "An error occurred") {
		t.Errorf("GetErrorMessage(unknown) should contain server or generic error message, got %q", msg)
	}
}

// --- ErrorNotificationMsg Tests ---

func TestNewErrorNotificationNil(t *testing.T) {
	notif := NewErrorNotification(nil)

	if notif.Category != 0 {
		t.Errorf("Category should be 0 for nil error, got %v", notif.Category)
	}

	if notif.Message != "" {
		t.Errorf("Message should be empty for nil error, got %q", notif.Message)
	}
}

func TestNewErrorNotificationConnection(t *testing.T) {
	err := ErrConnection
	notif := NewErrorNotification(err)

	if notif.Category != ErrCategoryConnection {
		t.Errorf("Category = %v, want %v", notif.Category, ErrCategoryConnection)
	}

	if notif.Action != "Check connection" {
		t.Errorf("Action = %q, want %q", notif.Action, "Check connection")
	}

	if !notif.Retryable {
		t.Error("Connection errors should be retryable")
	}
}

func TestNewErrorNotificationTimeout(t *testing.T) {
	err := status.Error(codes.DeadlineExceeded, "timeout")
	notif := NewErrorNotification(err)

	if notif.Category != ErrCategoryTimeout {
		t.Errorf("Category = %v, want %v", notif.Category, ErrCategoryTimeout)
	}

	if notif.Action != "Retry" {
		t.Errorf("Action = %q, want %q", notif.Action, "Retry")
	}

	if !notif.Retryable {
		t.Error("Timeout errors should be retryable")
	}
}

func TestNewErrorNotificationServerRetryable(t *testing.T) {
	err := status.Error(codes.Unavailable, "service unavailable")
	notif := NewErrorNotification(err)

	if notif.Category != ErrCategoryConnection {
		t.Errorf("Category = %v, want %v", notif.Category, ErrCategoryConnection)
	}

	if notif.Action != "Check connection" {
		t.Errorf("Action = %q, want %q", notif.Action, "Check connection")
	}

	if !notif.Retryable {
		t.Error("Unavailable errors should be retryable")
	}
}

func TestNewErrorNotificationServerNonRetryable(t *testing.T) {
	err := status.Error(codes.Internal, "internal error")
	notif := NewErrorNotification(err)

	if notif.Category != ErrCategoryServer {
		t.Errorf("Category = %v, want %v", notif.Category, ErrCategoryServer)
	}

	// Internal errors are retryable by default
	if !notif.Retryable {
		t.Error("Internal errors should be retryable")
	}

	if notif.Action != "Retry" {
		t.Errorf("Action = %q, want %q", notif.Action, "Retry")
	}
}

func TestNewErrorNotificationValidation(t *testing.T) {
	err := ErrInvalidArgument
	notif := NewErrorNotification(err)

	if notif.Category != ErrCategoryValidation {
		t.Errorf("Category = %v, want %v", notif.Category, ErrCategoryValidation)
	}

	if notif.Action != "Fix input" {
		t.Errorf("Action = %q, want %q", notif.Action, "Fix input")
	}

	if notif.Retryable {
		t.Error("Validation errors should not be retryable")
	}
}

func TestNewErrorNotificationUnknown(t *testing.T) {
	err := errors.New("strange error")
	notif := NewErrorNotification(err)

	if notif.Action != "Retry" {
		t.Errorf("Action for unknown error = %q, want %q", notif.Action, "Retry")
	}
}

// --- Message Content Tests ---

func TestErrorNotificationMessageContent(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantContain string
	}{
		{
			name:        "connection error",
			err:         ErrConnection,
			wantContain: "Cannot connect",
		},
		{
			name:        "timeout error",
			err:         ErrTimeout,
			wantContain: "timed out",
		},
		{
			name:        "validation error",
			err:         ErrInvalidArgument,
			wantContain: "Invalid data",
		},
		{
			name:        "server error",
			err:         ErrInternal,
			wantContain: "Server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notif := NewErrorNotification(tt.err)
			if !contains(notif.Message, tt.wantContain) {
				t.Errorf("Message %q should contain %q", notif.Message, tt.wantContain)
			}
		})
	}
}

// --- Action Suggestions Tests ---

func TestErrorNotificationActionSuggestions(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantAction string
	}{
		{
			name:       "connection",
			err:        errors.New("connection refused"),
			wantAction: "Check connection",
		},
		{
			name:       "timeout",
			err:        errors.New("operation timed out"),
			wantAction: "Retry",
		},
		{
			name:       "validation",
			err:        status.Error(codes.InvalidArgument, "invalid"),
			wantAction: "Fix input",
		},
		{
			name:       "server unavailable",
			err:        status.Error(codes.Unavailable, "unavailable"),
			wantAction: "Check connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notif := NewErrorNotification(tt.err)
			if notif.Action != tt.wantAction {
				t.Errorf("Action = %q, want %q", notif.Action, tt.wantAction)
			}
		})
	}
}

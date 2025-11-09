package mnemosyne

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestCategorizeErrorGRPCCodes verifies gRPC code categorization
func TestCategorizeErrorGRPCCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCategory
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: ErrCategoryUnknown,
		},
		{
			name:     "unavailable",
			err:      status.Error(codes.Unavailable, "service unavailable"),
			expected: ErrCategoryConnection,
		},
		{
			name:     "deadline exceeded",
			err:      status.Error(codes.DeadlineExceeded, "deadline exceeded"),
			expected: ErrCategoryTimeout,
		},
		{
			name:     "invalid argument",
			err:      status.Error(codes.InvalidArgument, "invalid"),
			expected: ErrCategoryValidation,
		},
		{
			name:     "out of range",
			err:      status.Error(codes.OutOfRange, "out of range"),
			expected: ErrCategoryValidation,
		},
		{
			name:     "internal error",
			err:      status.Error(codes.Internal, "internal error"),
			expected: ErrCategoryServer,
		},
		{
			name:     "unknown error",
			err:      status.Error(codes.Unknown, "unknown"),
			expected: ErrCategoryServer,
		},
		{
			name:     "data loss",
			err:      status.Error(codes.DataLoss, "data loss"),
			expected: ErrCategoryServer,
		},
		{
			name:     "failed precondition",
			err:      status.Error(codes.FailedPrecondition, "failed precondition"),
			expected: ErrCategoryConnection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := CategorizeError(tt.err)
			if category != tt.expected {
				t.Errorf("expected category %s, got %s", tt.expected, category)
			}
		})
	}
}

// TestCategorizeErrorKnownErrors verifies known error categorization
func TestCategorizeErrorKnownErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCategory
	}{
		{
			name:     "not connected",
			err:      ErrNotConnected,
			expected: ErrCategoryConnection,
		},
		{
			name:     "connection error",
			err:      ErrConnection,
			expected: ErrCategoryConnection,
		},
		{
			name:     "timeout error",
			err:      ErrTimeout,
			expected: ErrCategoryTimeout,
		},
		{
			name:     "invalid argument",
			err:      ErrInvalidArgument,
			expected: ErrCategoryValidation,
		},
		{
			name:     "memory not found",
			err:      ErrMemoryNotFound,
			expected: ErrCategoryValidation,
		},
		{
			name:     "already exists",
			err:      ErrAlreadyExists,
			expected: ErrCategoryValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := CategorizeError(tt.err)
			if category != tt.expected {
				t.Errorf("expected category %s, got %s", tt.expected, category)
			}
		})
	}
}

// TestCategorizeErrorNetworkErrors verifies network error categorization
func TestCategorizeErrorNetworkErrors(t *testing.T) {
	// Create a network timeout error
	timeoutErr := &net.DNSError{
		Err:       "timeout",
		IsTimeout: true,
	}

	category := CategorizeError(timeoutErr)
	if category != ErrCategoryTimeout {
		t.Errorf("expected timeout category, got %s", category)
	}

	// Create a network error (non-timeout)
	netErr := &net.OpError{
		Op:  "dial",
		Err: errors.New("connection refused"),
	}

	category = CategorizeError(netErr)
	if category != ErrCategoryConnection {
		t.Errorf("expected connection category, got %s", category)
	}
}

// TestCategorizeErrorMessagePatterns verifies error message pattern matching
func TestCategorizeErrorMessagePatterns(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCategory
	}{
		{
			name:     "timeout in message",
			err:      errors.New("operation timeout"),
			expected: ErrCategoryTimeout,
		},
		{
			name:     "timed out in message",
			err:      errors.New("request timed out"),
			expected: ErrCategoryTimeout,
		},
		{
			name:     "deadline in message",
			err:      errors.New("deadline exceeded"),
			expected: ErrCategoryTimeout,
		},
		{
			name:     "connection in message",
			err:      errors.New("connection refused"),
			expected: ErrCategoryConnection,
		},
		{
			name:     "network in message",
			err:      errors.New("network unreachable"),
			expected: ErrCategoryConnection,
		},
		{
			name:     "dial in message",
			err:      errors.New("dial failed"),
			expected: ErrCategoryConnection,
		},
		{
			name:     "unreachable in message",
			err:      errors.New("host unreachable"),
			expected: ErrCategoryConnection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := CategorizeError(tt.err)
			if category != tt.expected {
				t.Errorf("expected category %s, got %s", tt.expected, category)
			}
		})
	}
}

// TestCategorizeErrorMnemosyneError verifies MnemosyneError categorization
func TestCategorizeErrorMnemosyneError(t *testing.T) {
	mnemosyneErr := &MnemosyneError{
		Category:   ErrCategoryConnection,
		Code:       "connection",
		Message:    "test error",
		Retryable:  true,
		Underlying: ErrConnection,
	}

	category := CategorizeError(mnemosyneErr)
	if category != ErrCategoryConnection {
		t.Errorf("expected connection category, got %s", category)
	}
}

// TestIsRetryableGRPCErrors verifies retryability of gRPC errors
func TestIsRetryableGRPCErrors(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "nil error",
			err:       nil,
			retryable: false,
		},
		{
			name:      "unavailable",
			err:       status.Error(codes.Unavailable, "unavailable"),
			retryable: true,
		},
		{
			name:      "resource exhausted",
			err:       status.Error(codes.ResourceExhausted, "exhausted"),
			retryable: true,
		},
		{
			name:      "invalid argument",
			err:       status.Error(codes.InvalidArgument, "invalid"),
			retryable: false,
		},
		{
			name:      "not found",
			err:       status.Error(codes.NotFound, "not found"),
			retryable: false,
		},
		{
			name:      "already exists",
			err:       status.Error(codes.AlreadyExists, "exists"),
			retryable: false,
		},
		{
			name:      "deadline exceeded",
			err:       status.Error(codes.DeadlineExceeded, "timeout"),
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryable := IsRetryable(tt.err)
			if retryable != tt.retryable {
				t.Errorf("expected retryable=%v, got %v", tt.retryable, retryable)
			}
		})
	}
}

// TestIsRetryableKnownErrors verifies retryability of known errors
func TestIsRetryableKnownErrors(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "connection error",
			err:       ErrConnection,
			retryable: true,
		},
		{
			name:      "timeout error",
			err:       ErrTimeout,
			retryable: true,
		},
		{
			name:      "invalid argument",
			err:       ErrInvalidArgument,
			retryable: false,
		},
		{
			name:      "not connected",
			err:       ErrNotConnected,
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryable := IsRetryable(tt.err)
			if retryable != tt.retryable {
				t.Errorf("expected retryable=%v, got %v", tt.retryable, retryable)
			}
		})
	}
}

// TestIsRetryableMnemosyneError verifies MnemosyneError retryability
func TestIsRetryableMnemosyneError(t *testing.T) {
	retryableErr := &MnemosyneError{
		Category:   ErrCategoryConnection,
		Code:       "connection",
		Message:    "test",
		Retryable:  true,
		Underlying: ErrConnection,
	}

	if !IsRetryable(retryableErr) {
		t.Error("expected retryable error to be retryable")
	}

	nonRetryableErr := &MnemosyneError{
		Category:   ErrCategoryValidation,
		Code:       "validation",
		Message:    "test",
		Retryable:  false,
		Underlying: ErrInvalidArgument,
	}

	if IsRetryable(nonRetryableErr) {
		t.Error("expected non-retryable error not to be retryable")
	}
}

// TestWrapErrorWithContext verifies error wrapping with context
func TestWrapErrorWithContext(t *testing.T) {
	originalErr := errors.New("original error")
	wrapped := WrapError(originalErr, ErrCategoryConnection, "test message")

	if wrapped == nil {
		t.Fatal("expected wrapped error")
	}

	if wrapped.Category != ErrCategoryConnection {
		t.Errorf("expected category connection, got %s", wrapped.Category)
	}

	if wrapped.Message != "test message" {
		t.Errorf("expected message 'test message', got %q", wrapped.Message)
	}

	if !errors.Is(wrapped, originalErr) {
		t.Error("expected wrapped error to contain original error")
	}
}

// TestMnemosyneErrorFormatting verifies error message formatting
func TestMnemosyneErrorFormatting(t *testing.T) {
	tests := []struct {
		name     string
		err      *MnemosyneError
		contains string
	}{
		{
			name: "with underlying",
			err: &MnemosyneError{
				Category:   ErrCategoryConnection,
				Code:       "connection",
				Message:    "test message",
				Underlying: errors.New("underlying"),
			},
			contains: "underlying",
		},
		{
			name: "without underlying",
			err: &MnemosyneError{
				Category:   ErrCategoryValidation,
				Code:       "validation",
				Message:    "test message",
				Underlying: nil,
			},
			contains: "test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			if msg == "" {
				t.Error("expected non-empty error message")
			}
		})
	}
}

// TestMnemosyneErrorUnwrapping verifies error unwrapping
func TestMnemosyneErrorUnwrapping(t *testing.T) {
	underlying := errors.New("underlying error")
	mnemosyneErr := &MnemosyneError{
		Category:   ErrCategoryConnection,
		Code:       "connection",
		Message:    "test",
		Underlying: underlying,
	}

	unwrapped := mnemosyneErr.Unwrap()
	if !errors.Is(unwrapped, underlying) {
		t.Error("expected unwrapped error to match underlying")
	}
}

// TestWrapErrorGRPC verifies gRPC error wrapping
func TestWrapErrorGRPC(t *testing.T) {
	tests := []struct {
		name      string
		grpcErr   error
		operation string
		wantErr   error
	}{
		{
			name:      "OK code",
			grpcErr:   status.Error(codes.OK, "ok"),
			operation: "test",
			wantErr:   nil,
		},
		{
			name:      "NotFound",
			grpcErr:   status.Error(codes.NotFound, "not found"),
			operation: "get memory",
			wantErr:   ErrMemoryNotFound,
		},
		{
			name:      "InvalidArgument",
			grpcErr:   status.Error(codes.InvalidArgument, "invalid"),
			operation: "store memory",
			wantErr:   ErrInvalidArgument,
		},
		{
			name:      "AlreadyExists",
			grpcErr:   status.Error(codes.AlreadyExists, "exists"),
			operation: "create",
			wantErr:   ErrAlreadyExists,
		},
		{
			name:      "PermissionDenied",
			grpcErr:   status.Error(codes.PermissionDenied, "denied"),
			operation: "access",
			wantErr:   ErrPermissionDenied,
		},
		{
			name:      "Unavailable",
			grpcErr:   status.Error(codes.Unavailable, "unavailable"),
			operation: "connect",
			wantErr:   ErrUnavailable,
		},
		{
			name:      "Internal",
			grpcErr:   status.Error(codes.Internal, "internal"),
			operation: "process",
			wantErr:   ErrInternal,
		},
		{
			name:      "DeadlineExceeded",
			grpcErr:   status.Error(codes.DeadlineExceeded, "timeout"),
			operation: "request",
			wantErr:   ErrTimeout,
		},
		{
			name:      "Canceled",
			grpcErr:   status.Error(codes.Canceled, "canceled"),
			operation: "request",
			wantErr:   nil, // No specific error, just wrapped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := wrapError(tt.grpcErr, tt.operation)

			if tt.wantErr == nil && wrapped != nil && tt.name == "OK code" {
				t.Errorf("expected nil error for OK code, got %v", wrapped)
			}

			if tt.wantErr != nil && !errors.Is(wrapped, tt.wantErr) {
				t.Errorf("expected wrapped error to contain %v, got %v", tt.wantErr, wrapped)
			}
		})
	}
}

// TestWrapErrorNonGRPCFormat verifies non-gRPC error wrapping format
func TestWrapErrorNonGRPCFormat(t *testing.T) {
	originalErr := errors.New("some error")
	wrapped := wrapError(originalErr, "test operation")

	if wrapped == nil {
		t.Fatal("expected wrapped error")
	}

	if !errors.Is(wrapped, originalErr) {
		t.Error("expected wrapped error to contain original error")
	}
}

// TestErrorHelperFunctions verifies error helper functions
func TestErrorHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		checkFn  func(error) bool
		expected bool
	}{
		{
			name:     "IsNotFound - true",
			err:      ErrMemoryNotFound,
			checkFn:  IsNotFound,
			expected: true,
		},
		{
			name:     "IsNotFound - false",
			err:      ErrConnection,
			checkFn:  IsNotFound,
			expected: false,
		},
		{
			name:     "IsInvalidArgument - true",
			err:      ErrInvalidArgument,
			checkFn:  IsInvalidArgument,
			expected: true,
		},
		{
			name:     "IsUnavailable - true",
			err:      ErrUnavailable,
			checkFn:  IsUnavailable,
			expected: true,
		},
		{
			name:     "IsConnectionError - true (NotConnected)",
			err:      ErrNotConnected,
			checkFn:  IsConnectionError,
			expected: true,
		},
		{
			name:     "IsConnectionError - true (Connection)",
			err:      ErrConnection,
			checkFn:  IsConnectionError,
			expected: true,
		},
		{
			name:     "IsTimeout - true",
			err:      ErrTimeout,
			checkFn:  IsTimeout,
			expected: true,
		},
		{
			name:     "IsTimeout - false",
			err:      ErrConnection,
			checkFn:  IsTimeout,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checkFn(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestErrorCategoryStringRepresentation verifies error category string representation
func TestErrorCategoryStringRepresentation(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		expected string
	}{
		{ErrCategoryConnection, "connection"},
		{ErrCategoryServer, "server"},
		{ErrCategoryValidation, "validation"},
		{ErrCategoryTimeout, "timeout"},
		{ErrCategoryUnknown, "unknown"},
		{ErrorCategory(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			str := tt.category.String()
			if str != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, str)
			}
		})
	}
}

// TestIsConnectionErrorWithGRPC verifies connection error detection with gRPC
func TestIsConnectionErrorWithGRPC(t *testing.T) {
	grpcErr := status.Error(codes.Unavailable, "service unavailable")
	if !IsConnectionError(grpcErr) {
		t.Error("expected Unavailable to be a connection error")
	}
}

// TestIsTimeoutWithNetworkError verifies timeout detection with network errors
func TestIsTimeoutWithNetworkError(t *testing.T) {
	timeoutErr := &net.DNSError{
		Err:       "timeout",
		IsTimeout: true,
	}

	if !IsTimeout(timeoutErr) {
		t.Error("expected network timeout to be detected as timeout")
	}
}

// TestContextCancellationHandling verifies context cancellation error handling
func TestContextCancellationHandling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := ctx.Err()
	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	// Context cancellation should not be retryable
	category := CategorizeError(err)
	// Context errors fall into server category by default
	if category == ErrCategoryUnknown {
		t.Errorf("expected specific category, got unknown")
	}
}

// TestContextDeadlineExceededHandling verifies deadline exceeded error handling
func TestContextDeadlineExceededHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(1 * time.Millisecond)

	err := ctx.Err()
	if err == nil {
		t.Fatal("expected context deadline exceeded error")
	}

	// Deadline exceeded should be timeout category
	category := CategorizeError(err)
	if category != ErrCategoryTimeout {
		t.Errorf("expected timeout category, got %s", category)
	}
}

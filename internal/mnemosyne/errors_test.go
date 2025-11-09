package mnemosyne

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- Error Wrapping Tests ---

func TestWrapErrorNil(t *testing.T) {
	err := wrapError(nil, "test operation")
	if err != nil {
		t.Errorf("wrapError(nil) should return nil, got: %v", err)
	}
}

func TestWrapErrorNonGRPC(t *testing.T) {
	originalErr := errors.New("regular error")
	wrappedErr := wrapError(originalErr, "test operation")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	// Should contain operation name
	errStr := wrappedErr.Error()
	if !contains(errStr, "test operation") {
		t.Errorf("Expected error to contain operation name, got: %s", errStr)
	}
}

func TestWrapErrorNotFound(t *testing.T) {
	grpcErr := status.Error(codes.NotFound, "memory not found")
	wrappedErr := wrapError(grpcErr, "get memory")

	if !IsNotFound(wrappedErr) {
		t.Errorf("Expected NotFound error, got: %v", wrappedErr)
	}
}

func TestWrapErrorInvalidArgument(t *testing.T) {
	grpcErr := status.Error(codes.InvalidArgument, "invalid content")
	wrappedErr := wrapError(grpcErr, "store memory")

	if !IsInvalidArgument(wrappedErr) {
		t.Errorf("Expected InvalidArgument error, got: %v", wrappedErr)
	}
}

func TestWrapErrorAlreadyExists(t *testing.T) {
	grpcErr := status.Error(codes.AlreadyExists, "memory exists")
	wrappedErr := wrapError(grpcErr, "store memory")

	if !errors.Is(wrappedErr, ErrAlreadyExists) {
		t.Errorf("Expected AlreadyExists error, got: %v", wrappedErr)
	}
}

func TestWrapErrorPermissionDenied(t *testing.T) {
	grpcErr := status.Error(codes.PermissionDenied, "access denied")
	wrappedErr := wrapError(grpcErr, "delete memory")

	if !errors.Is(wrappedErr, ErrPermissionDenied) {
		t.Errorf("Expected PermissionDenied error, got: %v", wrappedErr)
	}
}

func TestWrapErrorUnavailable(t *testing.T) {
	grpcErr := status.Error(codes.Unavailable, "service down")
	wrappedErr := wrapError(grpcErr, "health check")

	if !IsUnavailable(wrappedErr) {
		t.Errorf("Expected Unavailable error, got: %v", wrappedErr)
	}
}

func TestWrapErrorInternal(t *testing.T) {
	grpcErr := status.Error(codes.Internal, "internal error")
	wrappedErr := wrapError(grpcErr, "list memories")

	if !errors.Is(wrappedErr, ErrInternal) {
		t.Errorf("Expected Internal error, got: %v", wrappedErr)
	}
}

func TestWrapErrorDeadlineExceeded(t *testing.T) {
	grpcErr := status.Error(codes.DeadlineExceeded, "timeout")
	wrappedErr := wrapError(grpcErr, "recall")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	errStr := wrappedErr.Error()
	if !contains(errStr, "timed out") && !contains(errStr, "deadline exceeded") {
		t.Errorf("Expected error to contain 'timed out' or 'deadline exceeded', got: %s", errStr)
	}
}

func TestWrapErrorCanceled(t *testing.T) {
	grpcErr := status.Error(codes.Canceled, "operation canceled")
	wrappedErr := wrapError(grpcErr, "store memory")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	errStr := wrappedErr.Error()
	if !contains(errStr, "canceled") {
		t.Errorf("Expected error to contain 'canceled', got: %s", errStr)
	}
}

func TestWrapErrorUnknownCode(t *testing.T) {
	grpcErr := status.Error(codes.Unimplemented, "not implemented")
	wrappedErr := wrapError(grpcErr, "test operation")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	// Should include operation name and status code
	errStr := wrappedErr.Error()
	if !contains(errStr, "test operation") {
		t.Errorf("Expected error to contain operation name, got: %s", errStr)
	}
}

// --- Error Check Function Tests ---

func TestIsNotFoundTrue(t *testing.T) {
	err := errors.New("test: memory not found")
	wrappedErr := errors.Join(err, ErrMemoryNotFound)

	if !IsNotFound(wrappedErr) {
		t.Error("Expected IsNotFound to return true")
	}
}

func TestIsNotFoundFalse(t *testing.T) {
	err := errors.New("different error")

	if IsNotFound(err) {
		t.Error("Expected IsNotFound to return false")
	}
}

func TestIsInvalidArgumentTrue(t *testing.T) {
	wrappedErr := errors.Join(errors.New("test"), ErrInvalidArgument)

	if !IsInvalidArgument(wrappedErr) {
		t.Error("Expected IsInvalidArgument to return true")
	}
}

func TestIsInvalidArgumentFalse(t *testing.T) {
	err := errors.New("different error")

	if IsInvalidArgument(err) {
		t.Error("Expected IsInvalidArgument to return false")
	}
}

func TestIsUnavailableTrue(t *testing.T) {
	wrappedErr := errors.Join(errors.New("test"), ErrUnavailable)

	if !IsUnavailable(wrappedErr) {
		t.Error("Expected IsUnavailable to return true")
	}
}

func TestIsUnavailableFalse(t *testing.T) {
	err := errors.New("different error")

	if IsUnavailable(err) {
		t.Error("Expected IsUnavailable to return false")
	}
}

// --- Standard Error Tests ---

func TestStandardErrors(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"ErrNotConnected", ErrNotConnected},
		{"ErrMemoryNotFound", ErrMemoryNotFound},
		{"ErrInvalidArgument", ErrInvalidArgument},
		{"ErrAlreadyExists", ErrAlreadyExists},
		{"ErrPermissionDenied", ErrPermissionDenied},
		{"ErrUnavailable", ErrUnavailable},
		{"ErrInternal", ErrInternal},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Errorf("%s should not be nil", tc.name)
			}

			if tc.err.Error() == "" {
				t.Errorf("%s should have a non-empty error message", tc.name)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}

// --- Error Categorization Tests ---

func TestCategorizeConnectionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorCategory
	}{
		{"ErrNotConnected", ErrNotConnected, ErrCategoryConnection},
		{"ErrConnection", ErrConnection, ErrCategoryConnection},
		{"gRPC Unavailable", status.Error(codes.Unavailable, "unavailable"), ErrCategoryConnection},
		{"gRPC FailedPrecondition", status.Error(codes.FailedPrecondition, "failed"), ErrCategoryConnection},
		{"message with connection", errors.New("connection refused"), ErrCategoryConnection},
		{"message with network", errors.New("network unreachable"), ErrCategoryConnection},
		{"message with dial", errors.New("dial failed"), ErrCategoryConnection},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeError(tt.err)
			if got != tt.want {
				t.Errorf("CategorizeError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestCategorizeServerError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorCategory
	}{
		{"ErrInternal", ErrInternal, ErrCategoryServer},
		{"ErrUnavailable", ErrUnavailable, ErrCategoryServer},
		{"gRPC Internal", status.Error(codes.Internal, "internal"), ErrCategoryServer},
		{"gRPC Unknown", status.Error(codes.Unknown, "unknown"), ErrCategoryServer},
		{"gRPC DataLoss", status.Error(codes.DataLoss, "data loss"), ErrCategoryServer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeError(tt.err)
			if got != tt.want {
				t.Errorf("CategorizeError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestCategorizeValidationError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorCategory
	}{
		{"ErrInvalidArgument", ErrInvalidArgument, ErrCategoryValidation},
		{"gRPC InvalidArgument", status.Error(codes.InvalidArgument, "invalid"), ErrCategoryValidation},
		{"gRPC OutOfRange", status.Error(codes.OutOfRange, "out of range"), ErrCategoryValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeError(tt.err)
			if got != tt.want {
				t.Errorf("CategorizeError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestCategorizeTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorCategory
	}{
		{"ErrTimeout", ErrTimeout, ErrCategoryTimeout},
		{"gRPC DeadlineExceeded", status.Error(codes.DeadlineExceeded, "timeout"), ErrCategoryTimeout},
		{"message with timeout", errors.New("operation timed out"), ErrCategoryTimeout},
		{"message with deadline", errors.New("deadline exceeded"), ErrCategoryTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CategorizeError(tt.err)
			if got != tt.want {
				t.Errorf("CategorizeError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// --- Retry Logic Tests ---

func TestIsRetryableConnectionErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrConnection", ErrConnection, true},
		{"ErrNotConnected", ErrNotConnected, true},
		{"gRPC Unavailable", status.Error(codes.Unavailable, "unavailable"), true},
		{"connection error message", errors.New("connection refused"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsRetryableTimeoutErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrTimeout", ErrTimeout, true},
		{"gRPC DeadlineExceeded", status.Error(codes.DeadlineExceeded, "timeout"), true},
		{"timeout message", errors.New("operation timed out"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsRetryableNonRetryableErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"ErrInvalidArgument", ErrInvalidArgument, false},
		{"ErrMemoryNotFound", ErrMemoryNotFound, false},
		{"gRPC InvalidArgument", status.Error(codes.InvalidArgument, "invalid"), false},
		{"gRPC NotFound", status.Error(codes.NotFound, "not found"), false},
		{"gRPC AlreadyExists", status.Error(codes.AlreadyExists, "exists"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// --- MnemosyneError Tests ---

func TestWrapErrorWithCategory(t *testing.T) {
	baseErr := errors.New("base error")
	wrapped := WrapError(baseErr, ErrCategoryConnection, "connection failed")

	if wrapped == nil {
		t.Fatal("WrapError returned nil")
	}

	if wrapped.Category != ErrCategoryConnection {
		t.Errorf("Category = %v, want %v", wrapped.Category, ErrCategoryConnection)
	}

	if wrapped.Message != "connection failed" {
		t.Errorf("Message = %q, want %q", wrapped.Message, "connection failed")
	}

	if wrapped.Underlying != baseErr {
		t.Errorf("Underlying error mismatch")
	}
}

func TestWrapErrorNilError(t *testing.T) {
	wrapped := WrapError(nil, ErrCategoryConnection, "test")
	if wrapped != nil {
		t.Errorf("WrapError(nil) should return nil, got %v", wrapped)
	}
}

func TestMnemosyneErrorError(t *testing.T) {
	baseErr := errors.New("underlying error")
	mnErr := &MnemosyneError{
		Category:   ErrCategoryConnection,
		Code:       "connection",
		Message:    "failed to connect",
		Underlying: baseErr,
	}

	errStr := mnErr.Error()
	if !contains(errStr, "connection") {
		t.Errorf("Error string should contain code, got: %s", errStr)
	}
	if !contains(errStr, "failed to connect") {
		t.Errorf("Error string should contain message, got: %s", errStr)
	}
}

func TestMnemosyneErrorUnwrap(t *testing.T) {
	baseErr := errors.New("base")
	mnErr := &MnemosyneError{
		Underlying: baseErr,
	}

	unwrapped := mnErr.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

// --- Error Category String Tests ---

func TestErrorCategoryString(t *testing.T) {
	tests := []struct {
		category ErrorCategory
		want     string
	}{
		{ErrCategoryConnection, "connection"},
		{ErrCategoryServer, "server"},
		{ErrCategoryValidation, "validation"},
		{ErrCategoryTimeout, "timeout"},
		{ErrCategoryUnknown, "unknown"},
		{ErrorCategory(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.category.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- New Error Check Function Tests ---

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrConnection", ErrConnection, true},
		{"ErrNotConnected", ErrNotConnected, true},
		{"gRPC Unavailable", status.Error(codes.Unavailable, "unavailable"), true},
		{"connection message", errors.New("connection failed"), true},
		{"validation error", ErrInvalidArgument, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsConnectionError(tt.err)
			if got != tt.want {
				t.Errorf("IsConnectionError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrTimeout", ErrTimeout, true},
		{"gRPC DeadlineExceeded", status.Error(codes.DeadlineExceeded, "timeout"), true},
		{"timeout message", errors.New("timed out"), true},
		{"connection error", ErrConnection, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTimeout(tt.err)
			if got != tt.want {
				t.Errorf("IsTimeout(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

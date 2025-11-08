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
	if !contains(errStr, "deadline exceeded") {
		t.Errorf("Expected error to contain 'deadline exceeded', got: %s", errStr)
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

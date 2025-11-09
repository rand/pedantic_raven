package mnemosyne

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Common errors with clear, actionable messages
var (
	ErrNotConnected     = errors.New("✗ Connection: Not connected to mnemosyne server. Check if running and config.")
	ErrMemoryNotFound   = errors.New("✗ Memory: Not found. May have been deleted. Refresh (Ctrl+R).")
	ErrInvalidArgument  = errors.New("✗ Validation: Invalid argument provided. Check format and required fields.")
	ErrAlreadyExists    = errors.New("✗ Memory: Already exists. Cannot create duplicate.")
	ErrPermissionDenied = errors.New("✗ Permission: Access denied. Check credentials and permissions.")
	ErrUnavailable      = errors.New("✗ Connection: Mnemosyne server unavailable. Check if running.")
	ErrInternal         = errors.New("✗ Server: Internal error. Check server logs for details.")
	ErrTimeout          = errors.New("✗ Timeout: Operation took too long. Check connection or try again.")
	ErrConnection       = errors.New("✗ Connection: Connection error. Check network and server status.")
)

// ErrorCategory classifies errors for appropriate handling
type ErrorCategory int

const (
	ErrCategoryConnection ErrorCategory = iota
	ErrCategoryServer
	ErrCategoryValidation
	ErrCategoryTimeout
	ErrCategoryUnknown
)

// String returns the string representation of the error category.
func (ec ErrorCategory) String() string {
	switch ec {
	case ErrCategoryConnection:
		return "connection"
	case ErrCategoryServer:
		return "server"
	case ErrCategoryValidation:
		return "validation"
	case ErrCategoryTimeout:
		return "timeout"
	case ErrCategoryUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// MnemosyneError provides rich error context
type MnemosyneError struct {
	Category   ErrorCategory
	Code       string
	Message    string
	Retryable  bool
	Underlying error
}

// Error implements the error interface.
func (e *MnemosyneError) Error() string {
	if e.Underlying != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Underlying)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements the error unwrapping interface.
func (e *MnemosyneError) Unwrap() error {
	return e.Underlying
}

// CategorizeError determines the category of an error.
func CategorizeError(err error) ErrorCategory {
	if err == nil {
		return ErrCategoryUnknown
	}

	// Check for MnemosyneError
	var mnemosyneErr *MnemosyneError
	if errors.As(err, &mnemosyneErr) {
		return mnemosyneErr.Category
	}

	// Check for known error types
	if errors.Is(err, ErrNotConnected) || errors.Is(err, ErrConnection) {
		return ErrCategoryConnection
	}

	if errors.Is(err, ErrTimeout) {
		return ErrCategoryTimeout
	}

	if errors.Is(err, ErrInvalidArgument) || errors.Is(err, ErrMemoryNotFound) || errors.Is(err, ErrAlreadyExists) {
		return ErrCategoryValidation
	}

	// Check for network errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return ErrCategoryTimeout
		}
		return ErrCategoryConnection
	}

	// Check gRPC status codes
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.Unavailable, codes.FailedPrecondition:
			return ErrCategoryConnection
		case codes.DeadlineExceeded:
			return ErrCategoryTimeout
		case codes.InvalidArgument, codes.OutOfRange:
			return ErrCategoryValidation
		case codes.Internal, codes.Unknown, codes.DataLoss:
			return ErrCategoryServer
		}
	}

	// Check error message for timeout-related keywords (before connection checks)
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "timed out") ||
		strings.Contains(errMsg, "deadline") {
		return ErrCategoryTimeout
	}

	// Check error message for connection-related keywords
	if strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "network") ||
		strings.Contains(errMsg, "dial") ||
		strings.Contains(errMsg, "unreachable") {
		return ErrCategoryConnection
	}

	// Default to server error
	return ErrCategoryServer
}

// IsRetryable determines if an error should be retried.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Check for MnemosyneError
	var mnemosyneErr *MnemosyneError
	if errors.As(err, &mnemosyneErr) {
		return mnemosyneErr.Retryable
	}

	// Check category
	category := CategorizeError(err)
	switch category {
	case ErrCategoryConnection, ErrCategoryTimeout:
		return true
	case ErrCategoryValidation:
		return false
	case ErrCategoryServer:
		// Check gRPC status for retryability
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unavailable, codes.ResourceExhausted:
				return true
			case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists:
				return false
			}
		}
		return true
	}

	return false
}

// WrapError wraps an error with additional context and categorization.
func WrapError(err error, category ErrorCategory, message string) *MnemosyneError {
	if err == nil {
		return nil
	}

	return &MnemosyneError{
		Category:   category,
		Code:       category.String(),
		Message:    message,
		Retryable:  IsRetryable(err),
		Underlying: err,
	}
}

// wrapError converts gRPC status errors to domain errors.
func wrapError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Extract gRPC status
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error, return as-is with context
		return fmt.Errorf("%s failed: %w", operation, err)
	}

	// Map gRPC status codes to domain errors
	switch st.Code() {
	case codes.OK:
		return nil

	case codes.NotFound:
		return fmt.Errorf("%s: %w: %s", operation, ErrMemoryNotFound, st.Message())

	case codes.InvalidArgument:
		return fmt.Errorf("%s: %w: %s", operation, ErrInvalidArgument, st.Message())

	case codes.AlreadyExists:
		return fmt.Errorf("%s: %w: %s", operation, ErrAlreadyExists, st.Message())

	case codes.PermissionDenied:
		return fmt.Errorf("%s: %w: %s", operation, ErrPermissionDenied, st.Message())

	case codes.Unavailable:
		return fmt.Errorf("%s: %w: %s", operation, ErrUnavailable, st.Message())

	case codes.Internal:
		return fmt.Errorf("%s: %w: %s", operation, ErrInternal, st.Message())

	case codes.DeadlineExceeded:
		return fmt.Errorf("%s: %w: %s", operation, ErrTimeout, st.Message())

	case codes.Canceled:
		return fmt.Errorf("%s: canceled: %s", operation, st.Message())

	default:
		return fmt.Errorf("%s failed: %s (code: %s)", operation, st.Message(), st.Code())
	}
}

// IsNotFound returns true if the error is a "not found" error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrMemoryNotFound)
}

// IsInvalidArgument returns true if the error is an "invalid argument" error.
func IsInvalidArgument(err error) bool {
	return errors.Is(err, ErrInvalidArgument)
}

// IsUnavailable returns true if the error is a "service unavailable" error.
func IsUnavailable(err error) bool {
	return errors.Is(err, ErrUnavailable)
}

// IsConnectionError returns true if the error is a connection error.
func IsConnectionError(err error) bool {
	return errors.Is(err, ErrNotConnected) || errors.Is(err, ErrConnection) ||
		CategorizeError(err) == ErrCategoryConnection
}

// IsTimeout returns true if the error is a timeout error.
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout) || CategorizeError(err) == ErrCategoryTimeout
}

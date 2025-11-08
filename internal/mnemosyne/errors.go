package mnemosyne

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Common errors
var (
	ErrNotConnected    = errors.New("not connected to mnemosyne server")
	ErrMemoryNotFound  = errors.New("memory not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrAlreadyExists   = errors.New("memory already exists")
	ErrPermissionDenied = errors.New("permission denied")
	ErrUnavailable     = errors.New("service unavailable")
	ErrInternal        = errors.New("internal server error")
)

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
		return fmt.Errorf("%s: deadline exceeded: %s", operation, st.Message())

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

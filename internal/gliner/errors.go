package gliner

import (
	"errors"
	"fmt"
)

// Standard errors
var (
	// ErrDisabled is returned when GLiNER is disabled in config
	ErrDisabled = errors.New("gliner: service disabled")

	// ErrServiceUnavailable is returned when the GLiNER service is not reachable
	ErrServiceUnavailable = errors.New("gliner: service unavailable")

	// ErrInvalidRequest is returned for malformed requests
	ErrInvalidRequest = errors.New("gliner: invalid request")

	// ErrExtractionFailed is returned when entity extraction fails
	ErrExtractionFailed = errors.New("gliner: extraction failed")

	// ErrTimeout is returned when a request times out
	ErrTimeout = errors.New("gliner: request timeout")

	// ErrModelNotLoaded is returned when the model is not loaded
	ErrModelNotLoaded = errors.New("gliner: model not loaded")
)

// Error wraps an error with additional context.
type Error struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

func (e *Error) Error() string {
	return fmt.Sprintf("gliner.%s: %v", e.Op, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// wrapError wraps an error with operation context.
func wrapError(op string, err error) error {
	if err == nil {
		return nil
	}
	return &Error{Op: op, Err: err}
}

// IsUnavailable checks if an error is ErrServiceUnavailable.
func IsUnavailable(err error) bool {
	return errors.Is(err, ErrServiceUnavailable)
}

// IsTimeout checks if an error is ErrTimeout.
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// IsDisabled checks if an error is ErrDisabled.
func IsDisabled(err error) bool {
	return errors.Is(err, ErrDisabled)
}

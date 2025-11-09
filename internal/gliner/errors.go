package gliner

import (
	"errors"
	"fmt"
)

// Standard errors with clear, actionable messages
var (
	// ErrDisabled is returned when GLiNER is disabled in config
	ErrDisabled = errors.New("✗ Service: GLiNER is disabled in config. Enable it to use semantic analysis.")

	// ErrServiceUnavailable is returned when the GLiNER service is not reachable
	ErrServiceUnavailable = errors.New("✗ Service: GLiNER unavailable. Start with: ./gliner-server")

	// ErrInvalidRequest is returned for malformed requests
	ErrInvalidRequest = errors.New("✗ Service: Invalid request to GLiNER. Check text encoding (UTF-8) and format.")

	// ErrExtractionFailed is returned when entity extraction fails
	ErrExtractionFailed = errors.New("✗ Service: GLiNER extraction failed. Text may be invalid or too long.")

	// ErrTimeout is returned when a request times out
	ErrTimeout = errors.New("✗ Service: GLiNER request timeout (30s). Text may be too large. Try shorter text.")

	// ErrModelNotLoaded is returned when the model is not loaded
	ErrModelNotLoaded = errors.New("✗ Service: GLiNER model not loaded. Reload with Ctrl+Shift+R.")
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

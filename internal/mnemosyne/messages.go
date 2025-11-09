package mnemosyne

import "fmt"

// GetErrorMessage returns a user-friendly error message.
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	category := CategorizeError(err)
	switch category {
	case ErrCategoryConnection:
		return "Cannot connect to mnemosyne server. Working offline with cached data."
	case ErrCategoryServer:
		return "Server error occurred. Changes will be retried automatically."
	case ErrCategoryValidation:
		return fmt.Sprintf("Invalid data: %v", err)
	case ErrCategoryTimeout:
		return "Operation timed out. Check your connection and try again."
	default:
		return fmt.Sprintf("An error occurred: %v", err)
	}
}

// ErrorNotificationMsg contains information for displaying error notifications in the UI.
type ErrorNotificationMsg struct {
	Category  ErrorCategory
	Message   string
	Retryable bool
	Action    string // User action to resolve (e.g., "Retry", "Check connection")
}

// NewErrorNotification creates an ErrorNotificationMsg from an error.
func NewErrorNotification(err error) ErrorNotificationMsg {
	if err == nil {
		return ErrorNotificationMsg{}
	}

	category := CategorizeError(err)
	retryable := IsRetryable(err)

	msg := ErrorNotificationMsg{
		Category:  category,
		Message:   GetErrorMessage(err),
		Retryable: retryable,
	}

	// Determine user action
	switch category {
	case ErrCategoryConnection:
		msg.Action = "Check connection"
	case ErrCategoryTimeout:
		msg.Action = "Retry"
	case ErrCategoryServer:
		if retryable {
			msg.Action = "Retry"
		} else {
			msg.Action = "Contact support"
		}
	case ErrCategoryValidation:
		msg.Action = "Fix input"
	default:
		msg.Action = "Retry"
	}

	return msg
}

package common

import (
	"strings"

	"github.com/lib/pq"
)

// IsUniqueConstraintError checks if the error is a unique constraint violation
func IsUniqueConstraintError(err error) bool {
	// Check for lib/pq errors
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name() == "unique_violation"
	}

	// Fallback to string matching for other database drivers
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "duplicate key value violates unique constraint") ||
		strings.Contains(errMsg, "unique constraint failed") ||
		strings.Contains(errMsg, "users_username_unique")
}

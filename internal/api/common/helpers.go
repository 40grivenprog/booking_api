package common

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Nullable helpers

// ToNullString converts string pointer to sql.NullString
func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// ToNullStringValue converts string to sql.NullString
func ToNullStringValue(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// ToNullInt64 converts int64 pointer to sql.NullInt64
func ToNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// ToNullInt64Value converts int64 to sql.NullInt64
func ToNullInt64Value(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

// ToNullUUID converts UUID to uuid.NullUUID
func ToNullUUID(id uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: id, Valid: true}
}

// FromNullString converts sql.NullString to string pointer
func FromNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// FromNullInt64 converts sql.NullInt64 to int64 pointer
func FromNullInt64(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

// Time formatting constants
const (
	TimeFormatRFC3339      = time.RFC3339
	TimeFormatDateOnly     = "2006-01-02"
	TimeFormatMonthOnly    = "2006-01"
	TimeFormatWithTimezone = "2006-01-02T15:04:05Z07:00"
)

// FormatTimeRFC3339 formats time to RFC3339 string
func FormatTimeRFC3339(t time.Time) string {
	return t.Format(TimeFormatRFC3339)
}

// FormatTimeWithTimezone formats time with timezone
func FormatTimeWithTimezone(t time.Time) string {
	return t.Format(TimeFormatWithTimezone)
}

// FormatDate formats time to date only (YYYY-MM-DD)
func FormatDate(t time.Time) string {
	return t.Format(TimeFormatDateOnly)
}

// FormatMonth formats time to month only (YYYY-MM)
func FormatMonth(t time.Time) string {
	return t.Format(TimeFormatMonthOnly)
}

// String helpers

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// Int64Ptr returns a pointer to the given int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// StringValue returns the value of string pointer or empty string if nil
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int64Value returns the value of int64 pointer or 0 if nil
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

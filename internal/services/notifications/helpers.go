package notifications

import "time"

// formatAppointmentTime formats appointment time for display
// Input: RFC3339 format (e.g., "2024-01-15T10:00:00Z")
// Output: date (YYYY-MM-DD), startTime (HH:MM), endTime (HH:MM)
func formatAppointmentTime(startTimeStr, endTimeStr string) (string, string, string) {
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return startTimeStr[:10], startTimeStr[11:16], endTimeStr[11:16]
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return startTimeStr[:10], startTimeStr[11:16], endTimeStr[11:16]
	}

	date := startTime.Format("2006-01-02")
	startTimeFormatted := startTime.Format("15:04")
	endTimeFormatted := endTime.Format("15:04")

	return date, startTimeFormatted, endTimeFormatted
}

// formatDate formats a date string for display
// Input: RFC3339 format (e.g., "2024-01-15T10:00:00Z")
// Output: formatted date string (YYYY-MM-DD)
func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// If parsing fails, try to extract date part
		if len(dateStr) >= 10 {
			return dateStr[:10]
		}
		return dateStr
	}
	return t.Format("2006-01-02")
}

package professionals

import (
	"fmt"
	"time"

	db "github.com/vention/booking_api/internal/repository"
)

// TimeSlot represents an availability time slot
type TimeSlot struct {
	StartTime   string
	EndTime     string
	Available   bool
	Type        string
	Description string
}

// AvailabilityConfig contains configuration for availability calculation
type AvailabilityConfig struct {
	WorkingHoursStart int
	WorkingHoursEnd   int
	AppTimezone       *time.Location
}

// GenerateAvailabilitySlots generates time slots for a specific date with availability info
func (s *service) GenerateAvailabilitySlots(date time.Time, appointments []*db.GetAppointmentsByProfessionalAndDateWithClientRow, config AvailabilityConfig) []TimeSlot {
	slots := make([]TimeSlot, 0, 18)

	// Use provided timezone for current time
	localNow := time.Now().In(config.AppTimezone)

	// Create base date in application timezone
	baseDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, config.AppTimezone)

	// Determine the starting hour based on current time
	startHour := config.WorkingHoursStart
	if date.Year() == localNow.Year() && date.Month() == localNow.Month() && date.Day() == localNow.Day() {
		// If it's today, start from current hour (rounded up to next hour)
		currentHour := localNow.Hour()
		if localNow.Minute() > 0 {
			currentHour++ // Round up to next hour if we're past the hour mark
		}
		if currentHour > config.WorkingHoursStart {
			startHour = currentHour
		}
	}

	for hour := startHour; hour < config.WorkingHoursEnd; hour++ {
		startTime := baseDate.Add(time.Duration(hour) * time.Hour)
		endTime := startTime.Add(time.Hour)

		// Skip if the slot is in the past (additional safety check)
		if startTime.Before(localNow) {
			continue
		}

		slot := TimeSlot{
			StartTime: formatTimeRFC3339(startTime),
			EndTime:   formatTimeRFC3339(endTime),
			Available: true,
		}

		// Check if this slot conflicts with any existing appointment
		for _, appointment := range appointments {
			apptStart := appointment.StartTime
			apptEnd := appointment.EndTime

			// Convert appointment times to application timezone for comparison
			apptStartLocal := apptStart.In(config.AppTimezone)
			apptEndLocal := apptEnd.In(config.AppTimezone)

			// Check if the slot overlaps with the appointment (both in application timezone)
			if startTime.Before(apptEndLocal) && endTime.After(apptStartLocal) {
				slot.Available = false
				slot.Type = string(appointment.Type)

				// Generate description based on appointment type and client info
				if appointment.Description.Valid {
					if appointment.ClientID.Valid && appointment.ClientFirstName.Valid && appointment.ClientLastName.Valid {
						// Appointment with client - show client info + description
						slot.Description = fmt.Sprintf("%s %s - %s",
							appointment.ClientFirstName.String,
							appointment.ClientLastName.String,
							appointment.Description.String)
					} else {
						// Unavailable period - show just description
						slot.Description = appointment.Description.String
					}
				}
				break
			}
		}

		slots = append(slots, slot)
	}

	return slots
}

// formatTimeRFC3339 formats time to RFC3339 string
func formatTimeRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

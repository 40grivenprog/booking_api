package api

// import (
// 	common "github.com/vention/booking_api/internal/api/common"
// 	db "github.com/vention/booking_api/internal/repository"
// )

// // mapAppointmentToCreateAppointmentResponse maps database result to API response
// func mapAppointmentToCreateAppointmentResponse(appointment *db.CreateAppointmentWithDetailsRow) CreateAppointmentResponse {
// 	return CreateAppointmentResponse{
// 		Appointment: Appointment{
// 			StartTime:   common.FormatTimeRFC3339(appointment.StartTime),
// 			EndTime:     common.FormatTimeRFC3339(appointment.EndTime),
// 			Description: appointment.Description.String,
// 		},
// 		Client: Client{
// 			FirstName: appointment.ClientFirstName.String,
// 			LastName:  appointment.ClientLastName.String,
// 		},
// 		Professional: Professional{
// 			ChatID: common.Int64Value(common.FromNullInt64(appointment.ProfessionalChatID)),
// 		},
// 	}
// }

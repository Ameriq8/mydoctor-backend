package models

import (
	"time"
)

// Enum for appointment statuses
type AppointmentStatus string

const (
	Scheduled   AppointmentStatus = "Scheduled"
	Completed   AppointmentStatus = "Completed"
	Cancelled   AppointmentStatus = "Cancelled"
	NoShow      AppointmentStatus = "No-Show"
	Rescheduled AppointmentStatus = "Rescheduled"
)

type FacilityAppointment struct {
	BaseModel
	PatientName          string            `json:"patient_name" db:"patient_name"`
	PatientContact       string            `json:"patient_contact" db:"patient_contact"`
	FacilityID           int64             `json:"facility_id" db:"facility_id"`
	DoctorID             int64             `json:"doctor_id" db:"doctor_id"`
	AppointmentTime      time.Time         `json:"appointment_time" db:"appointment_time"`
	Status               AppointmentStatus `json:"status" db:"status"`
	ReasonForAppointment string            `json:"reason_for_appointment" db:"reason_for_appointment"`
}

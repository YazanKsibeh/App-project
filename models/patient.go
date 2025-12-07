package models

// Patient struct represents a patient in the system
type Patient struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Phone             string `json:"phone"`
	Age               int    `json:"age"`
	Gender            string `json:"gender"`
	TotalRequired     int    `json:"total_required"`
	Allergies         string `json:"allergies"`
	CurrentMedications string `json:"current_medications"`
	MedicalConditions string `json:"medical_conditions"`
	SmokingStatus     bool   `json:"smoking_status"`
	PregnancyStatus   bool   `json:"pregnancy_status"`
	DentalHistory     string `json:"dental_history"`
	SpecialNotes      string `json:"special_notes"`
}

// PatientForm represents the data needed to create/update a patient
type PatientForm struct {
	Name              string `json:"name"`
	Phone             string `json:"phone"`
	Age               int    `json:"age"`
	Gender            string `json:"gender"`
	Allergies         string `json:"allergies"`
	CurrentMedications string `json:"current_medications"`
	MedicalConditions string `json:"medical_conditions"`
	SmokingStatus     bool   `json:"smoking_status"`
	PregnancyStatus   bool   `json:"pregnancy_status"`
	DentalHistory     string `json:"dental_history"`
	SpecialNotes      string `json:"special_notes"`
}

// Appointment struct represents an appointment in the system
// DateTime is in RFC3339 format (e.g., "2024-06-01T14:00:00Z")
type Appointment struct {
	ID        int    `json:"id"`
	PatientID int    `json:"patient_id"`
	DateTime  string `json:"datetime"`
	Duration  int    `json:"duration"` // in minutes
	Notes     string `json:"notes"`
}

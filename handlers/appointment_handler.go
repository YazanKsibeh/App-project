package handlers

import (
	"DentistApp/models"
	"database/sql"
)

// AppointmentHandler handles appointment-related database operations
type AppointmentHandler struct {
	db *sql.DB
}

// NewAppointmentHandler creates a new AppointmentHandler
func NewAppointmentHandler(db *sql.DB) *AppointmentHandler {
	return &AppointmentHandler{db: db}
}

// AddAppointment adds a new appointment to the database
func (h *AppointmentHandler) AddAppointment(appt models.Appointment) (int64, error) {
	query := `INSERT INTO appointments (patient_id, datetime, duration, notes) VALUES (?, ?, ?, ?)`
	result, err := h.db.Exec(query, appt.PatientID, appt.DateTime, appt.Duration, appt.Notes)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAppointments returns all appointments
func (h *AppointmentHandler) GetAppointments() ([]models.Appointment, error) {
	query := `SELECT id, patient_id, datetime, duration, notes FROM appointments ORDER BY datetime DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := []models.Appointment{}
	for rows.Next() {
		var appt models.Appointment
		err := rows.Scan(&appt.ID, &appt.PatientID, &appt.DateTime, &appt.Duration, &appt.Notes)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, appt)
	}
	return appointments, nil
}

// GetAppointment returns a specific appointment by ID
func (h *AppointmentHandler) GetAppointment(id int) (models.Appointment, error) {
	var appt models.Appointment
	query := `SELECT id, patient_id, datetime, duration, notes FROM appointments WHERE id = ?`
	err := h.db.QueryRow(query, id).Scan(&appt.ID, &appt.PatientID, &appt.DateTime, &appt.Duration, &appt.Notes)
	return appt, err
}

// UpdateAppointment updates an existing appointment
func (h *AppointmentHandler) UpdateAppointment(appt models.Appointment) error {
	query := `UPDATE appointments SET patient_id = ?, datetime = ?, duration = ?, notes = ? WHERE id = ?`
	_, err := h.db.Exec(query, appt.PatientID, appt.DateTime, appt.Duration, appt.Notes, appt.ID)
	return err
}

// DeleteAppointment deletes an appointment
func (h *AppointmentHandler) DeleteAppointment(id int) error {
	query := `DELETE FROM appointments WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}

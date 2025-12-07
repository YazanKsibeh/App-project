package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"DentistApp/models"
)

// PatientHandler handles patient-related database operations
type PatientHandler struct {
	db *sql.DB
}

// NewPatientHandler creates a new PatientHandler
func NewPatientHandler(db *sql.DB) *PatientHandler {
	return &PatientHandler{db: db}
}

// AddPatient adds a new patient to the database
func (h *PatientHandler) AddPatient(patient models.PatientForm) (int64, error) {
	// Check for phone number uniqueness
	var existingCount int
	err := h.db.QueryRow("SELECT COUNT(*) FROM patients WHERE phone = ?", patient.Phone).Scan(&existingCount)
	if err != nil {
		return 0, fmt.Errorf("failed to check phone uniqueness: %v", err)
	}
	if existingCount > 0 {
		return 0, fmt.Errorf("phone number already exists for another patient")
	}

	// Validate age range
	if patient.Age < 6 || patient.Age > 100 {
		return 0, fmt.Errorf("age must be between 6 and 100 years")
	}

	query := `
	INSERT INTO patients (name, phone, age, gender, allergies, current_medications, medical_conditions, smoking_status, pregnancy_status, dental_history, special_notes)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Convert boolean to integer for SQLite (0 = false, 1 = true)
	smokingStatus := 0
	if patient.SmokingStatus {
		smokingStatus = 1
	}
	pregnancyStatus := 0
	if patient.PregnancyStatus {
		pregnancyStatus = 1
	}

	result, err := h.db.Exec(query, patient.Name, patient.Phone, patient.Age, patient.Gender, 
		patient.Allergies, patient.CurrentMedications, patient.MedicalConditions, 
		smokingStatus, pregnancyStatus, patient.DentalHistory, patient.SpecialNotes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Create patient data directory
	patientDir := filepath.Join("patient_data", fmt.Sprintf("%d", id))
	err = os.MkdirAll(patientDir, 0755)
	if err != nil {
		// Optional: decide if you want to roll back the DB insert on directory creation failure
		return 0, err
	}

	// Create patient name folder inside the ID folder
	// Clean the patient name to be filesystem-safe
	cleanName := cleanPatientName(patient.Name)
	patientNameDir := filepath.Join(patientDir, cleanName)
	err = os.MkdirAll(patientNameDir, 0755)
	if err != nil {
		fmt.Printf("Warning: Failed to create patient name folder '%s': %v\n", patientNameDir, err)
		// Don't fail the entire operation if name folder creation fails
	}

	return id, nil
}

// GetPatients returns all patients from the database
func (h *PatientHandler) GetPatients() ([]models.Patient, error) {
	query := `SELECT id, name, phone, age, gender, total_required, allergies, current_medications, medical_conditions, smoking_status, pregnancy_status, dental_history, special_notes FROM patients ORDER BY name`
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := []models.Patient{}
	for rows.Next() {
		var patient models.Patient
		var smokingStatus, pregnancyStatus int
		err := rows.Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Age, &patient.Gender, &patient.TotalRequired,
			&patient.Allergies, &patient.CurrentMedications, &patient.MedicalConditions, 
			&smokingStatus, &pregnancyStatus, &patient.DentalHistory, &patient.SpecialNotes)
		if err != nil {
			return nil, err
		}
		// Convert integer to boolean
		patient.SmokingStatus = smokingStatus == 1
		patient.PregnancyStatus = pregnancyStatus == 1
		patients = append(patients, patient)
	}

	return patients, nil
}

// GetPatient returns a specific patient by ID
func (h *PatientHandler) GetPatient(id int) (models.Patient, error) {
	var patient models.Patient
	var smokingStatus, pregnancyStatus int
	query := `SELECT id, name, phone, age, gender, total_required, allergies, current_medications, medical_conditions, smoking_status, pregnancy_status, dental_history, special_notes FROM patients WHERE id = ?`
	err := h.db.QueryRow(query, id).Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Age, &patient.Gender, &patient.TotalRequired,
		&patient.Allergies, &patient.CurrentMedications, &patient.MedicalConditions,
		&smokingStatus, &pregnancyStatus, &patient.DentalHistory, &patient.SpecialNotes)
	if err != nil {
		return patient, err
	}
	// Convert integer to boolean
	patient.SmokingStatus = smokingStatus == 1
	patient.PregnancyStatus = pregnancyStatus == 1
	return patient, nil
}

// UpdatePatient updates an existing patient
func (h *PatientHandler) UpdatePatient(patient models.Patient) error {
	// Check for phone number uniqueness (excluding current patient)
	var existingCount int
	err := h.db.QueryRow("SELECT COUNT(*) FROM patients WHERE phone = ? AND id != ?", patient.Phone, patient.ID).Scan(&existingCount)
	if err != nil {
		return fmt.Errorf("failed to check phone uniqueness: %v", err)
	}
	if existingCount > 0 {
		return fmt.Errorf("phone number already exists for another patient")
	}

	// Validate age range
	if patient.Age < 6 || patient.Age > 100 {
		return fmt.Errorf("age must be between 6 and 100 years")
	}

	query := `
	UPDATE patients 
	SET name = ?, phone = ?, age = ?, gender = ?, total_required = ?, allergies = ?, current_medications = ?, medical_conditions = ?, smoking_status = ?, pregnancy_status = ?, dental_history = ?, special_notes = ?
	WHERE id = ?`

	// Convert boolean to integer for SQLite (0 = false, 1 = true)
	smokingStatus := 0
	if patient.SmokingStatus {
		smokingStatus = 1
	}
	pregnancyStatus := 0
	if patient.PregnancyStatus {
		pregnancyStatus = 1
	}

	_, err = h.db.Exec(query, patient.Name, patient.Phone, patient.Age, patient.Gender, patient.TotalRequired,
		patient.Allergies, patient.CurrentMedications, patient.MedicalConditions,
		smokingStatus, pregnancyStatus, patient.DentalHistory, patient.SpecialNotes, patient.ID)
	return err
}

// DeletePatient deletes a patient and their data directory
func (h *PatientHandler) DeletePatient(id int) error {
	// Ensure foreign keys are enabled
	_, err := h.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	// Start a transaction to ensure all operations succeed or fail together
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Enable foreign keys for this transaction
	_, err = tx.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys in transaction: %v", err)
	}

	// Delete patient (should cascade to appointments and payments)
	_, err = tx.Exec("DELETE FROM patients WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %v", err)
	}

	// Explicitly verify that related records were deleted (defensive programming)
	var appointmentCount, paymentCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM appointments WHERE patient_id = ?", id).Scan(&appointmentCount)
	if err != nil {
		return fmt.Errorf("failed to check appointments: %v", err)
	}
	err = tx.QueryRow("SELECT COUNT(*) FROM payments WHERE patient_id = ?", id).Scan(&paymentCount)
	if err != nil {
		return fmt.Errorf("failed to check payments: %v", err)
	}

	// If foreign key constraints didn't work, manually delete related records
	if appointmentCount > 0 {
		_, err = tx.Exec("DELETE FROM appointments WHERE patient_id = ?", id)
		if err != nil {
			return fmt.Errorf("failed to manually delete appointments: %v", err)
		}
	}
	if paymentCount > 0 {
		_, err = tx.Exec("DELETE FROM payments WHERE patient_id = ?", id)
		if err != nil {
			return fmt.Errorf("failed to manually delete payments: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Delete patient data directory
	patientDir := filepath.Join("patient_data", fmt.Sprintf("%d", id))
	return os.RemoveAll(patientDir)
}

// SearchPatients searches patients by name or phone
func (h *PatientHandler) SearchPatients(searchTerm string) ([]models.Patient, error) {
	query := `
	SELECT id, name, phone, age, gender, total_required, allergies, current_medications, medical_conditions, smoking_status, pregnancy_status, dental_history, special_notes
	FROM patients 
	WHERE name LIKE ? OR phone LIKE ?
	ORDER BY name`

	searchPattern := "%" + searchTerm + "%"
	rows, err := h.db.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients := []models.Patient{}
	for rows.Next() {
		var patient models.Patient
		var smokingStatus, pregnancyStatus int
		err := rows.Scan(&patient.ID, &patient.Name, &patient.Phone, &patient.Age, &patient.Gender, &patient.TotalRequired,
			&patient.Allergies, &patient.CurrentMedications, &patient.MedicalConditions,
			&smokingStatus, &pregnancyStatus, &patient.DentalHistory, &patient.SpecialNotes)
		if err != nil {
			return nil, err
		}
		// Convert integer to boolean
		patient.SmokingStatus = smokingStatus == 1
		patient.PregnancyStatus = pregnancyStatus == 1
		patients = append(patients, patient)
	}

	return patients, nil
}

// OpenPatientFolder opens the patient's folder in the system file explorer
func (h *PatientHandler) OpenPatientFolder(id int) error {
	fmt.Printf("[HANDLER] OpenPatientFolder called - Patient ID: %d\n", id)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("[HANDLER] ERROR: Failed to get current working directory: %v\n", err)
		return fmt.Errorf("failed to get current working directory: %v", err)
	}
	fmt.Printf("[HANDLER] Current working directory: %s\n", cwd)

	// Create absolute path
	folderPath := filepath.Join(cwd, "patient_data", fmt.Sprintf("%d", id))
	fmt.Printf("[HANDLER] Target folder path: %s\n", folderPath)

	// Check if directory exists, create if it doesn't
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Printf("[HANDLER] Directory does not exist, creating: %s\n", folderPath)
		err = os.MkdirAll(folderPath, 0755)
		if err != nil {
			fmt.Printf("[HANDLER] ERROR: Failed to create patient folder: %v\n", err)
			return fmt.Errorf("failed to create patient folder: %v", err)
		}
		fmt.Printf("[HANDLER] Directory created successfully\n")
	} else if err != nil {
		fmt.Printf("[HANDLER] ERROR: Failed to stat directory: %v\n", err)
		return fmt.Errorf("failed to check directory: %v", err)
	} else {
		fmt.Printf("[HANDLER] Directory already exists\n")
	}

	// Get patient info to create name folder if it doesn't exist
	// Skip if database is not available (e.g., in tests)
	if h.db != nil {
		patient, err := h.GetPatient(id)
		if err != nil {
			fmt.Printf("[HANDLER] Warning: Could not get patient info for name folder: %v\n", err)
		} else {
			// Create patient name folder inside the ID folder if it doesn't exist
			cleanName := cleanPatientName(patient.Name)
			patientNameDir := filepath.Join(folderPath, cleanName)
			if _, err := os.Stat(patientNameDir); os.IsNotExist(err) {
				fmt.Printf("[HANDLER] Creating patient name folder: %s\n", patientNameDir)
				err = os.MkdirAll(patientNameDir, 0755)
				if err != nil {
					fmt.Printf("[HANDLER] Warning: Failed to create patient name folder: %v\n", err)
				} else {
					fmt.Printf("[HANDLER] Patient name folder created successfully\n")
				}
			} else {
				fmt.Printf("[HANDLER] Patient name folder already exists: %s\n", patientNameDir)
			}
		}
	} else {
		fmt.Printf("[HANDLER] Database not available, skipping patient name folder creation\n")
	}

	// Convert to platform-specific path
	var cmd *exec.Cmd
	fmt.Printf("[HANDLER] Operating system: %s\n", runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		// Use Windows explorer directly with proper path handling
		cmd = exec.Command("explorer", folderPath)
	case "darwin":
		cmd = exec.Command("open", folderPath)
	default: // linux
		cmd = exec.Command("xdg-open", folderPath)
	}

	// Start the command and capture any errors
	fmt.Printf("[HANDLER] === COMMAND EXECUTION DEBUG ===\n")
	fmt.Printf("[HANDLER] Command Path: %s\n", cmd.Path)
	fmt.Printf("[HANDLER] Command Args: %v\n", cmd.Args)
	fmt.Printf("[HANDLER] Working Directory: %s\n", cmd.Dir)
	fmt.Printf("[HANDLER] Environment Variables: %v\n", cmd.Env)

	// Verify the target folder exists before running command
	if stat, err := os.Stat(folderPath); err != nil {
		fmt.Printf("[HANDLER] ERROR: Target folder stat failed: %v\n", err)
		return fmt.Errorf("target folder verification failed: %v", err)
	} else {
		fmt.Printf("[HANDLER] Target folder verified - IsDir: %v, ModTime: %v\n", stat.IsDir(), stat.ModTime())
	}

	fmt.Printf("[HANDLER] Attempting to start command...\n")

	// Run the command in a goroutine to avoid blocking
	go func() {
		fmt.Printf("[HANDLER] [GOROUTINE] Starting command execution...\n")

		// Set working directory for the command
		cmd.Dir = filepath.Dir(folderPath)
		fmt.Printf("[HANDLER] [GOROUTINE] Set command working directory to: %s\n", cmd.Dir)

		err := cmd.Run()
		if err != nil {
			fmt.Printf("[HANDLER] [GOROUTINE] Command FAILED with error: %v\n", err)
			// Try to get more details about the error
			if exitError, ok := err.(*exec.ExitError); ok {
				fmt.Printf("[HANDLER] [GOROUTINE] Exit code: %d\n", exitError.ExitCode())
				fmt.Printf("[HANDLER] [GOROUTINE] Stderr: %s\n", string(exitError.Stderr))
			}
		} else {
			fmt.Printf("[HANDLER] [GOROUTINE] Command completed SUCCESSFULLY\n")
		}
	}()

	fmt.Printf("[HANDLER] Command goroutine started, returning success\n")
	return nil
}

// DeleteAllPatients deletes all patients and their data directories
func (h *PatientHandler) DeleteAllPatients() error {
	// Ensure foreign keys are enabled
	_, err := h.db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %v", err)
	}

	// Start a transaction
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Enable foreign keys for this transaction
	_, err = tx.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys in transaction: %v", err)
	}

	// Delete all appointments first (defensive)
	_, err = tx.Exec("DELETE FROM appointments")
	if err != nil {
		return fmt.Errorf("failed to delete appointments: %v", err)
	}

	// Delete all payments (defensive)
	_, err = tx.Exec("DELETE FROM payments")
	if err != nil {
		return fmt.Errorf("failed to delete payments: %v", err)
	}

	// Delete all patients
	_, err = tx.Exec("DELETE FROM patients")
	if err != nil {
		return fmt.Errorf("failed to delete patients: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Remove all patient data directories
	return os.RemoveAll("patient_data")
}

// cleanPatientName sanitizes patient names for safe filesystem usage
func cleanPatientName(name string) string {
	// Remove/replace characters that are not safe for folder names
	// Replace common problematic characters
	cleanName := strings.ReplaceAll(name, "/", "-")
	cleanName = strings.ReplaceAll(cleanName, "\\", "-")
	cleanName = strings.ReplaceAll(cleanName, ":", "-")
	cleanName = strings.ReplaceAll(cleanName, "*", "-")
	cleanName = strings.ReplaceAll(cleanName, "?", "-")
	cleanName = strings.ReplaceAll(cleanName, "\"", "-")
	cleanName = strings.ReplaceAll(cleanName, "<", "-")
	cleanName = strings.ReplaceAll(cleanName, ">", "-")
	cleanName = strings.ReplaceAll(cleanName, "|", "-")
	cleanName = strings.ReplaceAll(cleanName, ".", "-") // Replace dots with dashes

	// Remove leading/trailing spaces
	cleanName = strings.Trim(cleanName, " ")

	// Use regex to replace multiple consecutive dashes/spaces with single dash
	re := regexp.MustCompile(`[-\s]+`)
	cleanName = re.ReplaceAllString(cleanName, "-")

	// Ensure the name is not empty
	if cleanName == "" {
		cleanName = "Unknown-Patient"
	}

	// Limit length to avoid filesystem issues
	if len(cleanName) > 50 {
		cleanName = cleanName[:50]
	}

	// Remove trailing dash if any
	cleanName = strings.TrimSuffix(cleanName, "-")

	return cleanName
}

package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// DentalLabHandler handles dental lab operations
type DentalLabHandler struct {
	db *sql.DB
}

// NewDentalLabHandler creates new handler
func NewDentalLabHandler(db *sql.DB) *DentalLabHandler {
	return &DentalLabHandler{db: db}
}

// generateLabCode generates the next lab code (LAB-001, LAB-002, etc.)
func (h *DentalLabHandler) generateLabCode() (string, error) {
	// Get the maximum code from database
	var maxCode sql.NullString
	err := h.db.QueryRow(`SELECT MAX(code) FROM dental_labs WHERE code LIKE 'LAB-%'`).Scan(&maxCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get max code: %v", err)
	}

	// If no existing codes, start with LAB-001
	if !maxCode.Valid || maxCode.String == "" {
		return "LAB-001", nil
	}

	// Extract number from code (e.g., "LAB-001" -> 1)
	re := regexp.MustCompile(`LAB-(\d+)`)
	matches := re.FindStringSubmatch(maxCode.String)
	if len(matches) < 2 {
		// If format doesn't match, start fresh
		return "LAB-001", nil
	}

	// Parse and increment
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return "LAB-001", nil
	}

	// Increment and format
	nextNum := num + 1
	return fmt.Sprintf("LAB-%03d", nextNum), nil
}

// CreateDentalLab inserts new dental lab
func (h *DentalLabHandler) CreateDentalLab(lab models.DentalLabForm) (int64, error) {
	// Validate required fields
	if lab.Name == "" {
		return 0, fmt.Errorf("lab name is required")
	}
	if lab.ContactPerson == "" {
		return 0, fmt.Errorf("contact person is required")
	}
	if lab.PhonePrimary == "" {
		return 0, fmt.Errorf("phone primary is required")
	}

	// Validate phone format (10 digits)
	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	if !phoneRegex.MatchString(lab.PhonePrimary) {
		return 0, fmt.Errorf("phone primary must be exactly 10 digits")
	}

	// Validate phone secondary if provided
	if lab.PhoneSecondary != "" && !phoneRegex.MatchString(lab.PhoneSecondary) {
		return 0, fmt.Errorf("phone secondary must be exactly 10 digits")
	}

	// Check name uniqueness
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM dental_labs WHERE name = ?", lab.Name).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check lab name: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("lab name already exists")
	}

	// Generate code
	code, err := h.generateLabCode()
	if err != nil {
		return 0, fmt.Errorf("failed to generate lab code: %v", err)
	}

	// Set defaults
	isActive := lab.IsActive
	if !lab.IsActive {
		isActive = true // Default to active
	}

	query := `INSERT INTO dental_labs (code, name, contact_person, phone_primary, phone_secondary, email, specialties, is_active, notes) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := h.db.Exec(query, code, lab.Name, lab.ContactPerson, lab.PhonePrimary, 
		lab.PhoneSecondary, lab.Email, lab.Specialties, isActive, lab.Notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetDentalLabsPaginated returns paginated dental labs ordered by name
func (h *DentalLabHandler) GetDentalLabsPaginated(page, pageSize int) (*models.DentalLabsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Get total count
	var totalCount int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM dental_labs`).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count dental labs: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, code, name, contact_person, phone_primary, phone_secondary, email, specialties, is_active, notes, created_at, updated_at 
	          FROM dental_labs ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load dental labs: %v", err)
	}
	defer rows.Close()

	labs := make([]models.DentalLab, 0)
	for rows.Next() {
		var lab models.DentalLab
		err := rows.Scan(&lab.ID, &lab.Code, &lab.Name, &lab.ContactPerson, &lab.PhonePrimary, 
			&lab.PhoneSecondary, &lab.Email, &lab.Specialties, &lab.IsActive, &lab.Notes, 
			&lab.CreatedAt, &lab.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan dental lab: %v", err)
		}
		labs = append(labs, lab)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dental lab rows error: %v", err)
	}

	return &models.DentalLabsResponse{
		Labs:        labs,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

// GetDentalLab returns a specific dental lab by id
func (h *DentalLabHandler) GetDentalLab(id int) (*models.DentalLab, error) {
	query := `SELECT id, code, name, contact_person, phone_primary, phone_secondary, email, specialties, is_active, notes, created_at, updated_at 
	          FROM dental_labs WHERE id = ?`
	
	var lab models.DentalLab
	err := h.db.QueryRow(query, id).Scan(&lab.ID, &lab.Code, &lab.Name, &lab.ContactPerson, 
		&lab.PhonePrimary, &lab.PhoneSecondary, &lab.Email, &lab.Specialties, &lab.IsActive, 
		&lab.Notes, &lab.CreatedAt, &lab.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dental lab not found")
		}
		return nil, fmt.Errorf("failed to get dental lab: %v", err)
	}

	return &lab, nil
}

// UpdateDentalLab updates dental lab by id
func (h *DentalLabHandler) UpdateDentalLab(id int, lab models.DentalLabForm) error {
	// Validate required fields
	if lab.Name == "" {
		return fmt.Errorf("lab name is required")
	}
	if lab.ContactPerson == "" {
		return fmt.Errorf("contact person is required")
	}
	if lab.PhonePrimary == "" {
		return fmt.Errorf("phone primary is required")
	}

	// Validate phone format (10 digits)
	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	if !phoneRegex.MatchString(lab.PhonePrimary) {
		return fmt.Errorf("phone primary must be exactly 10 digits")
	}

	// Validate phone secondary if provided
	if lab.PhoneSecondary != "" && !phoneRegex.MatchString(lab.PhoneSecondary) {
		return fmt.Errorf("phone secondary must be exactly 10 digits")
	}

	// Check name uniqueness (excluding current lab)
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM dental_labs WHERE name = ? AND id != ?", lab.Name, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check lab name: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("lab name already exists")
	}

	query := `UPDATE dental_labs SET name = ?, contact_person = ?, phone_primary = ?, phone_secondary = ?, 
	          email = ?, specialties = ?, is_active = ?, notes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = h.db.Exec(query, lab.Name, lab.ContactPerson, lab.PhonePrimary, lab.PhoneSecondary, 
		lab.Email, lab.Specialties, lab.IsActive, lab.Notes, id)
	return err
}

// DeleteDentalLab deletes dental lab by id
func (h *DentalLabHandler) DeleteDentalLab(id int) error {
	query := `DELETE FROM dental_labs WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}


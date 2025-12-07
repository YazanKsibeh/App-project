package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"strings"
)

// SessionHandler handles session-related database operations
type SessionHandler struct {
	db *sql.DB
}

// NewSessionHandler creates a new SessionHandler
func NewSessionHandler(db *sql.DB) *SessionHandler {
	return &SessionHandler{db: db}
}

// CreateSession creates a new session with items
func (h *SessionHandler) CreateSession(session models.SessionForm) (int64, error) {
	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Calculate total amount from items
	totalAmount := 0
	for _, item := range session.Items {
		totalAmount += item.Amount
	}

	// Insert session
	query := `INSERT INTO sessions (patient_id, dentist_id, session_date, total_amount, status, notes)
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, session.PatientID, session.DentistID, session.SessionDate,
		totalAmount, session.Status, session.Notes)
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %v", err)
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get session ID: %v", err)
	}

	// Insert session items
	for _, item := range session.Items {
		itemQuery := `INSERT INTO session_items (session_id, procedure_id, item_name, amount)
		              VALUES (?, ?, ?, ?)`
		var procedureID interface{}
		if item.ProcedureID != nil {
			procedureID = *item.ProcedureID
		}
		_, err := tx.Exec(itemQuery, sessionID, procedureID, item.ItemName, item.Amount)
		if err != nil {
			return 0, fmt.Errorf("failed to create session item: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return sessionID, nil
}

// GetSessions returns paginated sessions (10 per page, newest first)
// If filters is nil, returns all sessions
func (h *SessionHandler) GetSessions(page int, filters *models.SessionFilters) (models.SessionsResponse, error) {
	pageSize := 10
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	// Build WHERE clause and arguments
	whereConditions := []string{}
	args := []interface{}{}

	// Standard filters (AND logic)
	if filters != nil {
		if filters.PatientID != nil {
			whereConditions = append(whereConditions, "s.patient_id = ?")
			args = append(args, *filters.PatientID)
		}

		if filters.Status != nil && *filters.Status != "" {
			whereConditions = append(whereConditions, "s.status = ?")
			args = append(args, *filters.Status)
		}

		if filters.DentistID != nil {
			whereConditions = append(whereConditions, "s.dentist_id = ?")
			args = append(args, *filters.DentistID)
		}

		if filters.DateFrom != nil && *filters.DateFrom != "" {
			whereConditions = append(whereConditions, "DATE(s.session_date) >= ?")
			args = append(args, *filters.DateFrom)
		}

		if filters.DateTo != nil && *filters.DateTo != "" {
			whereConditions = append(whereConditions, "DATE(s.session_date) <= ?")
			args = append(args, *filters.DateTo)
		}

		// Procedure filter (OR logic - session must contain ANY of the selected procedures)
		if len(filters.ProcedureIDs) > 0 {
			// Use EXISTS subquery to check if session has any of the selected procedures
			placeholders := make([]string, len(filters.ProcedureIDs))
			procedureArgs := make([]interface{}, len(filters.ProcedureIDs))
			for i, procID := range filters.ProcedureIDs {
				placeholders[i] = "?"
				procedureArgs[i] = procID
			}
			// Add procedure args to main args
			args = append(args, procedureArgs...)
			procedureFilter := fmt.Sprintf(
				"EXISTS (SELECT 1 FROM session_items si WHERE si.session_id = s.id AND si.procedure_id IN (%s))",
				strings.Join(placeholders, ","),
			)
			whereConditions = append(whereConditions, procedureFilter)
		}
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Build count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT s.id)
		FROM sessions s
		LEFT JOIN patients p ON s.patient_id = p.id
		LEFT JOIN users u ON s.dentist_id = u.id
		LEFT JOIN invoices i ON s.id = i.session_id
		%s`, whereClause)

	var totalCount int
	err := h.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return models.SessionsResponse{}, fmt.Errorf("failed to get session count: %v", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT DISTINCT s.id, s.patient_id, s.dentist_id, s.session_date, s.total_amount, s.status, s.notes,
		       p.name as patient_name, u.username as dentist_name,
		       i.invoice_number
		FROM sessions s
		LEFT JOIN patients p ON s.patient_id = p.id
		LEFT JOIN users u ON s.dentist_id = u.id
		LEFT JOIN invoices i ON s.id = i.session_id
		%s
		ORDER BY s.session_date DESC
		LIMIT ? OFFSET ?`, whereClause)

	// Add pagination parameters
	queryArgs := append(args, pageSize, offset)

	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		return models.SessionsResponse{}, fmt.Errorf("failed to query sessions: %v", err)
	}
	defer rows.Close()

	sessions := make([]models.Session, 0)
	for rows.Next() {
		var session models.Session
		var invoiceNumber sql.NullString
		err := rows.Scan(&session.ID, &session.PatientID, &session.DentistID,
			&session.SessionDate, &session.TotalAmount, &session.Status, &session.Notes,
			&session.PatientName, &session.DentistName, &invoiceNumber)
		if err != nil {
			return models.SessionsResponse{}, fmt.Errorf("failed to scan session: %v", err)
		}
		if invoiceNumber.Valid {
			session.InvoiceNumber = invoiceNumber.String
		}
		sessions = append(sessions, session)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return models.SessionsResponse{
		Sessions:    sessions,
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		TotalCount:  totalCount,
	}, nil
}

// GetSession returns a specific session by ID with all items
func (h *SessionHandler) GetSession(id int) (models.Session, error) {
	var session models.Session

	// Get session with patient and dentist names
	query := `
		SELECT s.id, s.patient_id, s.dentist_id, s.session_date, s.total_amount, s.status, s.notes,
		       p.name as patient_name, u.username as dentist_name
		FROM sessions s
		LEFT JOIN patients p ON s.patient_id = p.id
		LEFT JOIN users u ON s.dentist_id = u.id
		WHERE s.id = ?`

	err := h.db.QueryRow(query, id).Scan(&session.ID, &session.PatientID, &session.DentistID,
		&session.SessionDate, &session.TotalAmount, &session.Status, &session.Notes,
		&session.PatientName, &session.DentistName)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, fmt.Errorf("session not found")
		}
		return session, fmt.Errorf("failed to get session: %v", err)
	}

	// Get session items
	itemsQuery := `SELECT id, session_id, procedure_id, item_name, amount
	               FROM session_items WHERE session_id = ? ORDER BY id`
	itemRows, err := h.db.Query(itemsQuery, id)
	if err != nil {
		return session, fmt.Errorf("failed to get session items: %v", err)
	}
	defer itemRows.Close()

	items := make([]models.SessionItem, 0)
	for itemRows.Next() {
		var item models.SessionItem
		var procedureID sql.NullInt64
		err := itemRows.Scan(&item.ID, &item.SessionID, &procedureID, &item.ItemName, &item.Amount)
		if err != nil {
			return session, fmt.Errorf("failed to scan session item: %v", err)
		}
		if procedureID.Valid {
			procID := int(procedureID.Int64)
			item.ProcedureID = &procID
		}
		items = append(items, item)
	}
	session.Items = items

	return session, nil
}

// UpdateSession updates an existing session and its items
func (h *SessionHandler) UpdateSession(session models.Session, items []models.SessionItemForm) error {
	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Calculate total amount from items
	totalAmount := 0
	for _, item := range items {
		totalAmount += item.Amount
	}

	// Update session
	query := `UPDATE sessions SET patient_id = ?, dentist_id = ?, session_date = ?, 
	          total_amount = ?, status = ?, notes = ? WHERE id = ?`
	_, err = tx.Exec(query, session.PatientID, session.DentistID, session.SessionDate,
		totalAmount, session.Status, session.Notes, session.ID)
	if err != nil {
		return fmt.Errorf("failed to update session: %v", err)
	}

	// Delete existing items
	_, err = tx.Exec("DELETE FROM session_items WHERE session_id = ?", session.ID)
	if err != nil {
		return fmt.Errorf("failed to delete session items: %v", err)
	}

	// Insert new items
	for _, item := range items {
		itemQuery := `INSERT INTO session_items (session_id, procedure_id, item_name, amount)
		              VALUES (?, ?, ?, ?)`
		var procedureID interface{}
		if item.ProcedureID != nil {
			procedureID = *item.ProcedureID
		}
		_, err := tx.Exec(itemQuery, session.ID, procedureID, item.ItemName, item.Amount)
		if err != nil {
			return fmt.Errorf("failed to create session item: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// DeleteSession deletes a session and its items (cascade)
func (h *SessionHandler) DeleteSession(id int) error {
	query := `DELETE FROM sessions WHERE id = ?`
	result, err := h.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

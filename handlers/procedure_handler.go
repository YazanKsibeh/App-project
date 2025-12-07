package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"log"
	"math"
)

// ProcedureHandler handles dental procedure operations
type ProcedureHandler struct {
	db *sql.DB
}

// NewProcedureHandler creates new handler
func NewProcedureHandler(db *sql.DB) *ProcedureHandler {
	return &ProcedureHandler{db: db}
}

// CreateProcedure inserts new procedure
func (h *ProcedureHandler) CreateProcedure(procedure models.ProcedureForm) (int64, error) {
	log.Printf("[ProcedureHandler] CreateProcedure called with name: %s, price: %d", procedure.Name, procedure.Price)
	
	if procedure.Name == "" {
		log.Println("[ProcedureHandler] Validation failed: procedure name is required")
		return 0, fmt.Errorf("procedure name is required")
	}
	if procedure.Price <= 0 {
		log.Printf("[ProcedureHandler] Validation failed: price must be greater than zero, got: %d", procedure.Price)
		return 0, fmt.Errorf("price must be greater than zero")
	}

	query := `INSERT INTO dental_procedures (name, price) VALUES (?, ?)`
	log.Printf("[ProcedureHandler] Executing query: %s with values: name=%s, price=%d", query, procedure.Name, procedure.Price)
	
	result, err := h.db.Exec(query, procedure.Name, procedure.Price)
	if err != nil {
		log.Printf("[ProcedureHandler] Database error: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("[ProcedureHandler] Error getting last insert ID: %v", err)
		return 0, err
	}
	
	log.Printf("[ProcedureHandler] Procedure created successfully with ID: %d", id)
	return id, nil
}

// GetProcedures returns all procedures ordered by name (legacy method for backward compatibility)
func (h *ProcedureHandler) GetProcedures() ([]models.Procedure, error) {
	query := `SELECT id, name, price, created_at FROM dental_procedures ORDER BY name`
	
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	procedures := make([]models.Procedure, 0)
	for rows.Next() {
		var procedure models.Procedure
		err := rows.Scan(&procedure.ID, &procedure.Name, &procedure.Price, &procedure.CreatedAt)
		if err != nil {
			return nil, err
		}
		procedures = append(procedures, procedure)
	}
	
	return procedures, nil
}

// GetProceduresPaginated returns paginated procedures ordered by name
func (h *ProcedureHandler) GetProceduresPaginated(page, pageSize int) (*models.ProceduresResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Get total count
	var totalCount int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM dental_procedures`).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count procedures: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, name, price, created_at FROM dental_procedures ORDER BY name LIMIT ? OFFSET ?`
	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load procedures: %v", err)
	}
	defer rows.Close()

	procedures := make([]models.Procedure, 0)
	for rows.Next() {
		var procedure models.Procedure
		err := rows.Scan(&procedure.ID, &procedure.Name, &procedure.Price, &procedure.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan procedure: %v", err)
		}
		procedures = append(procedures, procedure)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("procedure rows error: %v", err)
	}

	return &models.ProceduresResponse{
		Procedures:  procedures,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

// UpdateProcedure updates procedure by id
func (h *ProcedureHandler) UpdateProcedure(procedure models.Procedure) error {
	if procedure.Name == "" {
		return fmt.Errorf("procedure name is required")
	}
	if procedure.Price <= 0 {
		return fmt.Errorf("price must be greater than zero")
	}

	query := `UPDATE dental_procedures SET name = ?, price = ? WHERE id = ?`
	_, err := h.db.Exec(query, procedure.Name, procedure.Price, procedure.ID)
	return err
}

// DeleteProcedure deletes procedure by id
func (h *ProcedureHandler) DeleteProcedure(id int) error {
	query := `DELETE FROM dental_procedures WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}

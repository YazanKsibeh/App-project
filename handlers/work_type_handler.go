package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"math"
)

// WorkTypeHandler handles work type operations
type WorkTypeHandler struct {
	db *sql.DB
}

// NewWorkTypeHandler creates new handler
func NewWorkTypeHandler(db *sql.DB) *WorkTypeHandler {
	return &WorkTypeHandler{db: db}
}

// CreateWorkType inserts new work type
func (h *WorkTypeHandler) CreateWorkType(workType models.WorkTypeForm, userID int) (int64, error) {
	if workType.Name == "" {
		return 0, fmt.Errorf("work type name is required")
	}

	query := `INSERT INTO work_types (name, description, created_by) VALUES (?, ?, ?)`
	result, err := h.db.Exec(query, workType.Name, workType.Description, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetWorkTypesPaginated returns paginated work types ordered by name
func (h *WorkTypeHandler) GetWorkTypesPaginated(page, pageSize int) (*models.WorkTypesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Get total count
	var totalCount int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM work_types`).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count work types: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, name, description, sort_order, created_by, created_at, updated_at 
	          FROM work_types ORDER BY name LIMIT ? OFFSET ?`
	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load work types: %v", err)
	}
	defer rows.Close()

	workTypes := make([]models.WorkType, 0)
	for rows.Next() {
		var workType models.WorkType
		var createdBy sql.NullInt64
		err := rows.Scan(&workType.ID, &workType.Name, &workType.Description, &workType.SortOrder, 
			&createdBy, &workType.CreatedAt, &workType.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan work type: %v", err)
		}
		if createdBy.Valid {
			createdByInt := int(createdBy.Int64)
			workType.CreatedBy = &createdByInt
		}
		workTypes = append(workTypes, workType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("work type rows error: %v", err)
	}

	return &models.WorkTypesResponse{
		WorkTypes:   workTypes,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

// UpdateWorkType updates work type by id
func (h *WorkTypeHandler) UpdateWorkType(id int, workType models.WorkTypeForm, userID int) error {
	if workType.Name == "" {
		return fmt.Errorf("work type name is required")
	}

	query := `UPDATE work_types SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := h.db.Exec(query, workType.Name, workType.Description, id)
	return err
}

// DeleteWorkType deletes work type by id
func (h *WorkTypeHandler) DeleteWorkType(id int) error {
	query := `DELETE FROM work_types WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}


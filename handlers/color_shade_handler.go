package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"math"
	"strings"
)

// ColorShadeHandler handles color shade operations
type ColorShadeHandler struct {
	db *sql.DB
}

// NewColorShadeHandler creates new handler
func NewColorShadeHandler(db *sql.DB) *ColorShadeHandler {
	return &ColorShadeHandler{db: db}
}

// CreateColorShade inserts new color shade
func (h *ColorShadeHandler) CreateColorShade(shade models.ColorShadeForm, userID int) (int64, error) {
	if shade.Name == "" {
		return 0, fmt.Errorf("color shade name is required")
	}

	// Validate hex color format if provided
	if shade.HexColor != "" {
		hexColor := strings.TrimPrefix(strings.ToUpper(shade.HexColor), "#")
		if len(hexColor) != 6 {
			return 0, fmt.Errorf("hex color must be 6 characters (e.g., F1ECE4)")
		}
		// Add # prefix if not present
		if !strings.HasPrefix(shade.HexColor, "#") {
			shade.HexColor = "#" + shade.HexColor
		}
	}

	query := `INSERT INTO color_shades (name, description, hex_color, is_active, created_by) 
	          VALUES (?, ?, ?, ?, ?)`
	result, err := h.db.Exec(query, shade.Name, shade.Description, shade.HexColor, shade.IsActive, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetColorShadesPaginated returns paginated color shades ordered by name
func (h *ColorShadeHandler) GetColorShadesPaginated(page, pageSize int) (*models.ColorShadesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Get total count
	var totalCount int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM color_shades`).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count color shades: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, name, description, hex_color, is_active, sort_order, created_by, created_at, updated_at 
	          FROM color_shades ORDER BY name LIMIT ? OFFSET ?`
	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load color shades: %v", err)
	}
	defer rows.Close()

	colorShades := make([]models.ColorShade, 0)
	for rows.Next() {
		var shade models.ColorShade
		var createdBy sql.NullInt64
		err := rows.Scan(&shade.ID, &shade.Name, &shade.Description, &shade.HexColor, &shade.IsActive, 
			&shade.SortOrder, &createdBy, &shade.CreatedAt, &shade.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan color shade: %v", err)
		}
		if createdBy.Valid {
			createdByInt := int(createdBy.Int64)
			shade.CreatedBy = &createdByInt
		}
		colorShades = append(colorShades, shade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("color shade rows error: %v", err)
	}

	return &models.ColorShadesResponse{
		ColorShades: colorShades,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

// UpdateColorShade updates color shade by id
func (h *ColorShadeHandler) UpdateColorShade(id int, shade models.ColorShadeForm, userID int) error {
	if shade.Name == "" {
		return fmt.Errorf("color shade name is required")
	}

	// Validate hex color format if provided
	if shade.HexColor != "" {
		hexColor := strings.TrimPrefix(strings.ToUpper(shade.HexColor), "#")
		if len(hexColor) != 6 {
			return fmt.Errorf("hex color must be 6 characters (e.g., F1ECE4)")
		}
		// Add # prefix if not present
		if !strings.HasPrefix(shade.HexColor, "#") {
			shade.HexColor = "#" + shade.HexColor
		}
	}

	query := `UPDATE color_shades SET name = ?, description = ?, hex_color = ?, is_active = ?, 
	          updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := h.db.Exec(query, shade.Name, shade.Description, shade.HexColor, shade.IsActive, id)
	return err
}

// DeleteColorShade deletes color shade by id
func (h *ColorShadeHandler) DeleteColorShade(id int) error {
	query := `DELETE FROM color_shades WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}


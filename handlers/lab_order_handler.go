package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LabOrderHandler handles lab order operations
type LabOrderHandler struct {
	db *sql.DB
}

// NewLabOrderHandler creates new handler
func NewLabOrderHandler(db *sql.DB) *LabOrderHandler {
	return &LabOrderHandler{db: db}
}

// GetLabOrdersPaginated returns paginated lab orders with optimized queries
func (h *LabOrderHandler) GetLabOrdersPaginated(page, pageSize int, searchOrderNumber, searchPatientName, searchLabName, statusFilter string) (*models.LabOrdersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20 // Default 20 per page as specified
	}

	// Build WHERE conditions for search/filter
	whereConditions := []string{}
	args := []interface{}{}

	// Search by order number
	if searchOrderNumber != "" {
		whereConditions = append(whereConditions, "lo.order_number LIKE ?")
		args = append(args, "%"+strings.TrimSpace(searchOrderNumber)+"%")
	}

	// Search by patient name
	if searchPatientName != "" {
		whereConditions = append(whereConditions, "p.name LIKE ?")
		args = append(args, "%"+strings.TrimSpace(searchPatientName)+"%")
	}

	// Search by lab name
	if searchLabName != "" {
		whereConditions = append(whereConditions, "dl.name LIKE ?")
		args = append(args, "%"+strings.TrimSpace(searchLabName)+"%")
	}

	// Filter by status
	if statusFilter != "" && statusFilter != "all" {
		whereConditions = append(whereConditions, "lo.status = ?")
		args = append(args, statusFilter)
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Optimized COUNT query - only counts lab_orders, JOINs only if needed for WHERE
	var totalCount int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM lab_orders lo
		LEFT JOIN patients p ON lo.patient_id = p.id
		LEFT JOIN dental_labs dl ON lo.lab_id = dl.id
		%s`, whereClause)

	if err := h.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count lab orders: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	// Optimized data query with JOINs
	query := fmt.Sprintf(`
		SELECT 
			lo.id,
			lo.order_number,
			COALESCE(p.name, 'Unknown') as patient_name,
			COALESCE(dl.name, 'Unknown') as lab_name,
			COALESCE(wt.name, 'Unknown') as work_type_name,
			lo.status,
			lo.lab_cost,
			lo.order_date
		FROM lab_orders lo
		LEFT JOIN patients p ON lo.patient_id = p.id
		LEFT JOIN dental_labs dl ON lo.lab_id = dl.id
		LEFT JOIN work_types wt ON lo.work_type_id = wt.id
		%s
		ORDER BY lo.created_at DESC
		LIMIT ? OFFSET ?`, whereClause)

	// Add pagination parameters
	queryArgs := append(args, pageSize, offset)

	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to load lab orders: %v", err)
	}
	defer rows.Close()

	orders := make([]models.LabOrderListItem, 0)
	for rows.Next() {
		var order models.LabOrderListItem
		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.PatientName,
			&order.LabName,
			&order.WorkTypeName,
			&order.Status,
			&order.LabCost,
			&order.OrderDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab order: %v", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("lab order rows error: %v", err)
	}

	return &models.LabOrdersResponse{
		Orders:      orders,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

// GetLabOrder returns full details of a specific lab order
func (h *LabOrderHandler) GetLabOrder(id int) (*models.LabOrderDetail, error) {
	query := `
		SELECT 
			lo.id,
			lo.order_number,
			lo.patient_id,
			COALESCE(p.name, 'Unknown') as patient_name,
			lo.lab_id,
			COALESCE(dl.name, 'Unknown') as lab_name,
			lo.created_by,
			COALESCE(u.username, 'Unknown') as dentist_name,
			lo.work_type_id,
			COALESCE(wt.name, 'Unknown') as work_type_name,
			COALESCE(lo.description, '') as description,
			COALESCE(lo.upper_left, '') as upper_left,
			COALESCE(lo.upper_right, '') as upper_right,
			COALESCE(lo.lower_left, '') as lower_left,
			COALESCE(lo.lower_right, '') as lower_right,
			lo.quantity,
			lo.color_shade_id,
			COALESCE(cs.name, '') as color_shade_name,
			lo.lab_cost,
			lo.order_date,
			lo.status,
			COALESCE(lo.notes, '') as notes,
			lo.created_at,
			lo.updated_at
		FROM lab_orders lo
		LEFT JOIN patients p ON lo.patient_id = p.id
		LEFT JOIN dental_labs dl ON lo.lab_id = dl.id
		LEFT JOIN users u ON lo.created_by = u.id
		LEFT JOIN work_types wt ON lo.work_type_id = wt.id
		LEFT JOIN color_shades cs ON lo.color_shade_id = cs.id
		WHERE lo.id = ?`

	var order models.LabOrderDetail
	var colorShadeID sql.NullInt64

	err := h.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.PatientID,
		&order.PatientName,
		&order.LabID,
		&order.LabName,
		&order.CreatedBy,
		&order.DentistName,
		&order.WorkTypeID,
		&order.WorkTypeName,
		&order.Description,
		&order.UpperLeft,
		&order.UpperRight,
		&order.LowerLeft,
		&order.LowerRight,
		&order.Quantity,
		&colorShadeID,
		&order.ColorShadeName,
		&order.LabCost,
		&order.OrderDate,
		&order.Status,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lab order not found")
		}
		return nil, fmt.Errorf("failed to get lab order: %v", err)
	}

	// Handle nullable color_shade_id
	if colorShadeID.Valid {
		val := int(colorShadeID.Int64)
		order.ColorShadeID = &val
	}

	return &order, nil
}

// generateOrderNumber generates the next order number (ORDER-001, ORDER-002, etc.)
func (h *LabOrderHandler) generateOrderNumber() (string, error) {
	// Get the maximum order number from database
	var maxOrderNumber sql.NullString
	err := h.db.QueryRow(`SELECT MAX(order_number) FROM lab_orders WHERE order_number LIKE 'ORDER-%'`).Scan(&maxOrderNumber)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get max order number: %v", err)
	}

	// If no existing order numbers, start with ORDER-001
	if !maxOrderNumber.Valid || maxOrderNumber.String == "" {
		return "ORDER-001", nil
	}

	// Extract number from order number (e.g., "ORDER-001" -> 1)
	re := regexp.MustCompile(`ORDER-(\d+)`)
	matches := re.FindStringSubmatch(maxOrderNumber.String)
	if len(matches) < 2 {
		// If format doesn't match, start fresh
		return "ORDER-001", nil
	}

	// Parse and increment
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return "ORDER-001", nil
	}

	// Increment and format
	nextNum := num + 1
	return fmt.Sprintf("ORDER-%03d", nextNum), nil
}

// CreateLabOrder creates a new lab order
func (h *LabOrderHandler) CreateLabOrder(order models.LabOrderForm, createdBy int) (*models.CreateLabOrderResponse, error) {
	// Validate required fields
	if order.PatientID == 0 {
		return nil, fmt.Errorf("patient is required")
	}
	if order.LabID == 0 {
		return nil, fmt.Errorf("dental lab is required")
	}
	if order.WorkTypeID == 0 {
		return nil, fmt.Errorf("work type is required")
	}
	if order.OrderDate == "" {
		return nil, fmt.Errorf("order date is required")
	}

	// Validate lab cost >= 0
	if order.LabCost < 0 {
		return nil, fmt.Errorf("lab cost must be greater than or equal to 0")
	}

	// Validate quantity >= 1
	if order.Quantity < 1 {
		return nil, fmt.Errorf("quantity must be at least 1")
	}

	// Parse order date - Handle datetime-local format (2006-01-02T15:04) or other formats
	var orderDate time.Time
	var err error
	
	// Try datetime-local format first (from HTML5 input)
	if len(order.OrderDate) >= 16 {
		orderDate, err = time.Parse("2006-01-02T15:04", order.OrderDate[:16])
		if err != nil {
			// Try with seconds
			orderDate, err = time.Parse("2006-01-02T15:04:05", order.OrderDate)
		}
		if err != nil {
			// Try SQLite format
			orderDate, err = time.Parse("2006-01-02 15:04:05", order.OrderDate)
		}
		if err != nil {
			return nil, fmt.Errorf("invalid order date format: %v", err)
		}
	} else {
		return nil, fmt.Errorf("order date is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":     true,
		"sent":      true,
		"delivered": true,
		"cancelled": true,
	}
	if !validStatuses[order.Status] {
		return nil, fmt.Errorf("invalid status: %s", order.Status)
	}

	// Generate order number
	orderNumber, err := h.generateOrderNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate order number: %v", err)
	}

	// Set defaults
	quantity := order.Quantity
	if quantity < 1 {
		quantity = 1
	}

	labCost := order.LabCost
	if labCost < 0 {
		labCost = 0
	}

	// Format order date for database (SQLite uses TEXT for DATETIME)
	orderDateStr := orderDate.Format("2006-01-02 15:04:05")

	// Insert order
	query := `INSERT INTO lab_orders (
		order_number, patient_id, lab_id, created_by, work_type_id,
		description, upper_left, upper_right, lower_left, lower_right,
		quantity, color_shade_id, lab_cost, order_date, status, notes
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var colorShadeID interface{}
	if order.ColorShadeID != nil && *order.ColorShadeID > 0 {
		colorShadeID = *order.ColorShadeID
	} else {
		colorShadeID = nil
	}

	result, err := h.db.Exec(query,
		orderNumber,
		order.PatientID,
		order.LabID,
		createdBy,
		order.WorkTypeID,
		strings.TrimSpace(order.Description),
		strings.TrimSpace(order.UpperLeft),
		strings.TrimSpace(order.UpperRight),
		strings.TrimSpace(order.LowerLeft),
		strings.TrimSpace(order.LowerRight),
		quantity,
		colorShadeID,
		labCost,
		orderDateStr,
		order.Status,
		strings.TrimSpace(order.Notes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lab order: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get order ID: %v", err)
	}

	return &models.CreateLabOrderResponse{
		ID:          id,
		OrderNumber: orderNumber,
	}, nil
}


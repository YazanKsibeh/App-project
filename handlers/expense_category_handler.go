package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"log"
	"math"
)

// ExpenseCategoryHandler handles expense category operations
type ExpenseCategoryHandler struct {
	db *sql.DB
}

// NewExpenseCategoryHandler creates new handler
func NewExpenseCategoryHandler(db *sql.DB) *ExpenseCategoryHandler {
	return &ExpenseCategoryHandler{db: db}
}

// CreateExpenseCategory inserts new expense category
func (h *ExpenseCategoryHandler) CreateExpenseCategory(category models.ExpenseCategoryForm, userID int) (int64, error) {
	log.Printf("[ExpenseCategoryHandler] CreateExpenseCategory called with name: %s, userID: %d", category.Name, userID)

	if category.Name == "" {
		return 0, fmt.Errorf("category name is required")
	}

	// Validate expense type
	validExpenseTypes := map[string]bool{
		"operational":     true,
		"capital":         true,
		"personnel":       true,
		"marketing":       true,
		"administrative":  true,
	}
	if !validExpenseTypes[category.ExpenseType] {
		return 0, fmt.Errorf("invalid expense type: %s", category.ExpenseType)
	}

	// Check if name already exists
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM expense_categories WHERE name = ?", category.Name).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check category name: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("category name already exists")
	}

	// Set defaults for fields not shown in form
	color := category.Color
	if color == "" {
		color = "#3498db"
	}
	costCenter := "main"
	budgetAmount := 0
	budgetPeriod := "monthly"
	isActive := category.IsActive
	isTaxDeductible := true
	requiresApproval := false
	approvalThreshold := 0
	reportingGroup := ""
	sortOrder := 0
	var parentCategoryID *int = nil
	var accountCode interface{} = nil // Set to NULL (not used anymore)
	
	query := `INSERT INTO expense_categories (
	          name, description, color, expense_type, budget_amount, budget_period,
	          is_tax_deductible, cost_center, account_code, parent_category_id,
	          is_active, requires_approval, approval_threshold, reporting_group, sort_order, created_by
	          ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := h.db.Exec(query,
		category.Name, category.Description, color, category.ExpenseType,
		budgetAmount, budgetPeriod,
		isTaxDeductible, costCenter, accountCode, parentCategoryID,
		isActive, requiresApproval, approvalThreshold,
		reportingGroup, sortOrder, userID)
	if err != nil {
		log.Printf("[ExpenseCategoryHandler] Database error: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Printf("[ExpenseCategoryHandler] Expense category created successfully with ID: %d", id)
	return id, nil
}

// GetExpenseCategories returns all expense categories ordered by name
func (h *ExpenseCategoryHandler) GetExpenseCategories() ([]models.ExpenseCategory, error) {
	query := `SELECT id, name, description, color, budget_amount, budget_period, expense_type, 
	          is_tax_deductible, cost_center, account_code, parent_category_id, is_active, 
	          requires_approval, approval_threshold, reporting_group, sort_order, 
	          created_by, updated_by, created_at, updated_at 
	          FROM expense_categories 
	          ORDER BY name`

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.ExpenseCategory, 0)
	for rows.Next() {
		var category models.ExpenseCategory
		var parentCategoryID sql.NullInt64
		var createdBy sql.NullInt64
		var updatedBy sql.NullInt64
		var accountCode sql.NullString

		err := rows.Scan(
			&category.ID, &category.Name, &category.Description, &category.Color,
			&category.BudgetAmount, &category.BudgetPeriod, &category.ExpenseType,
			&category.IsTaxDeductible, &category.CostCenter, &accountCode,
			&parentCategoryID, &category.IsActive, &category.RequiresApproval,
			&category.ApprovalThreshold, &category.ReportingGroup, &category.SortOrder,
			&createdBy, &updatedBy, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if accountCode.Valid {
			category.AccountCode = accountCode.String
		} else {
			category.AccountCode = ""
		}
		if parentCategoryID.Valid {
			val := int(parentCategoryID.Int64)
			category.ParentCategoryID = &val
		}
		if createdBy.Valid {
			val := int(createdBy.Int64)
			category.CreatedBy = &val
		}
		if updatedBy.Valid {
			val := int(updatedBy.Int64)
			category.UpdatedBy = &val
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// GetExpenseCategoriesPaginated returns paginated expense categories ordered by name
func (h *ExpenseCategoryHandler) GetExpenseCategoriesPaginated(page, pageSize int) (*models.ExpenseCategoriesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Get total count
	var totalCount int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM expense_categories`).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count expense categories: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, name, description, color, budget_amount, budget_period, expense_type, 
	          is_tax_deductible, cost_center, account_code, parent_category_id, is_active, 
	          requires_approval, approval_threshold, reporting_group, sort_order, 
	          created_by, updated_by, created_at, updated_at 
	          FROM expense_categories 
	          ORDER BY name 
	          LIMIT ? OFFSET ?`

	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load expense categories: %v", err)
	}
	defer rows.Close()

	categories := make([]models.ExpenseCategory, 0)
	for rows.Next() {
		var category models.ExpenseCategory
		var parentCategoryID sql.NullInt64
		var createdBy sql.NullInt64
		var updatedBy sql.NullInt64
		var accountCode sql.NullString

		err := rows.Scan(
			&category.ID, &category.Name, &category.Description, &category.Color,
			&category.BudgetAmount, &category.BudgetPeriod, &category.ExpenseType,
			&category.IsTaxDeductible, &category.CostCenter, &accountCode,
			&parentCategoryID, &category.IsActive, &category.RequiresApproval,
			&category.ApprovalThreshold, &category.ReportingGroup, &category.SortOrder,
			&createdBy, &updatedBy, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense category: %v", err)
		}

		if accountCode.Valid {
			category.AccountCode = accountCode.String
		} else {
			category.AccountCode = ""
		}
		if parentCategoryID.Valid {
			val := int(parentCategoryID.Int64)
			category.ParentCategoryID = &val
		}
		if createdBy.Valid {
			val := int(createdBy.Int64)
			category.CreatedBy = &val
		}
		if updatedBy.Valid {
			val := int(updatedBy.Int64)
			category.UpdatedBy = &val
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("expense category rows error: %v", err)
	}

	return &models.ExpenseCategoriesResponse{
		Categories:  categories,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

// UpdateExpenseCategory updates expense category by id
func (h *ExpenseCategoryHandler) UpdateExpenseCategory(id int, category models.ExpenseCategoryForm, userID int) error {
	if category.Name == "" {
		return fmt.Errorf("category name is required")
	}

	// Validate expense type
	validExpenseTypes := map[string]bool{
		"operational":    true,
		"capital":        true,
		"personnel":      true,
		"marketing":      true,
		"administrative": true,
	}
	if !validExpenseTypes[category.ExpenseType] {
		return fmt.Errorf("invalid expense type: %s", category.ExpenseType)
	}


	// Check if name already exists for another category
	var existingID int
	err := h.db.QueryRow("SELECT id FROM expense_categories WHERE name = ?", category.Name).Scan(&existingID)
	if err == nil && existingID != id {
		return fmt.Errorf("category name already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check category name: %v", err)
	}

	// Set defaults for fields not shown in form
	color := category.Color
	if color == "" {
		color = "#3498db"
	}
	costCenter := "main"
	budgetAmount := 0
	budgetPeriod := "monthly"
	isActive := category.IsActive
	isTaxDeductible := true
	requiresApproval := false
	approvalThreshold := 0
	reportingGroup := ""
	sortOrder := 0
	var parentCategoryID *int = nil
	var accountCode interface{} = nil // Set to NULL (not used anymore)
	
	query := `UPDATE expense_categories 
	          SET name = ?, description = ?, color = ?, expense_type = ?, budget_amount = ?, budget_period = ?,
	              is_tax_deductible = ?, cost_center = ?, account_code = ?, parent_category_id = ?,
	              is_active = ?, requires_approval = ?, approval_threshold = ?, reporting_group = ?, sort_order = ?,
	              updated_by = ?, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = ?`
	_, err = h.db.Exec(query,
		category.Name, category.Description, color, category.ExpenseType,
		budgetAmount, budgetPeriod,
		isTaxDeductible, costCenter, accountCode, parentCategoryID,
		isActive, requiresApproval, approvalThreshold,
		reportingGroup, sortOrder, userID, id)
	return err
}

// DeleteExpenseCategory deletes expense category by id (soft delete - sets is_active = 0)
func (h *ExpenseCategoryHandler) DeleteExpenseCategory(id int) error {
	query := `UPDATE expense_categories SET is_active = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := h.db.Exec(query, id)
	return err
}

// PermanentlyDeleteExpenseCategory permanently deletes expense category by id (hard delete)
// Returns error if category has child categories or expenses referencing it
func (h *ExpenseCategoryHandler) PermanentlyDeleteExpenseCategory(id int) error {
	// Check for child categories (parent references)
	var childCount int
	err := h.db.QueryRow(
		"SELECT COUNT(*) FROM expense_categories WHERE parent_category_id = ? AND is_active = 1",
		id,
	).Scan(&childCount)
	if err != nil {
		return fmt.Errorf("failed to check child categories: %v", err)
	}

	// Check for expenses referencing this category
	var expenseCount int
	err = h.db.QueryRow(
		"SELECT COUNT(*) FROM expenses WHERE category_id = ?",
		id,
	).Scan(&expenseCount)
	if err != nil {
		return fmt.Errorf("failed to check expenses: %v", err)
	}

	// Build error message if references exist
	if childCount > 0 && expenseCount > 0 {
		return fmt.Errorf("cannot delete category: %d active child categories and %d expenses reference this category", childCount, expenseCount)
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete category: %d active child categories reference this category", childCount)
	}
	if expenseCount > 0 {
		return fmt.Errorf("cannot delete category: %d expenses reference this category", expenseCount)
	}

	// Perform hard delete
	query := `DELETE FROM expense_categories WHERE id = ?`
	result, err := h.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to permanently delete category: %v", err)
	}

	// Verify deletion
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}


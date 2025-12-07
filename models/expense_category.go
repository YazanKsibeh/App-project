package models

// ExpenseCategory represents an expense category in the system
type ExpenseCategory struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Color             string `json:"color"`
	BudgetAmount      int    `json:"budget_amount"`
	BudgetPeriod      string `json:"budget_period"`
	ExpenseType       string `json:"expense_type"`
	IsTaxDeductible   bool   `json:"is_tax_deductible"`
	CostCenter        string `json:"cost_center"`
	AccountCode       string `json:"account_code"`
	ParentCategoryID  *int   `json:"parent_category_id,omitempty"`
	IsActive          bool   `json:"is_active"`
	RequiresApproval  bool   `json:"requires_approval"`
	ApprovalThreshold int    `json:"approval_threshold"`
	ReportingGroup    string `json:"reporting_group"`
	SortOrder         int    `json:"sort_order"`
	CreatedBy         *int   `json:"created_by,omitempty"`
	UpdatedBy         *int   `json:"updated_by,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// ExpenseCategoryForm represents the data needed to create/update an expense category
type ExpenseCategoryForm struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Color             string `json:"color"`
	BudgetAmount      int    `json:"budget_amount"`
	BudgetPeriod      string `json:"budget_period"`
	ExpenseType       string `json:"expense_type"`
	IsTaxDeductible   bool   `json:"is_tax_deductible"`
	CostCenter        string `json:"cost_center"`
	AccountCode       string `json:"account_code"`
	ParentCategoryID  *int   `json:"parent_category_id,omitempty"`
	IsActive          bool   `json:"is_active"`
	RequiresApproval  bool   `json:"requires_approval"`
	ApprovalThreshold int    `json:"approval_threshold"`
	ReportingGroup    string `json:"reporting_group"`
	SortOrder         int    `json:"sort_order"`
}

// ExpenseCategoriesResponse represents the response for GetExpenseCategoriesPaginated
type ExpenseCategoriesResponse struct {
	Categories  []ExpenseCategory `json:"categories"`
	CurrentPage int               `json:"current_page"`
	TotalPages   int               `json:"total_pages"`
	PageSize     int               `json:"page_size"`
	TotalCount   int               `json:"total_count"`
}


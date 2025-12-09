package models

// WorkType represents a lab work type
type WorkType struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Description string `json:"description"`
	SortOrder int    `json:"sort_order"`
	CreatedBy *int   `json:"created_by,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// WorkTypeForm represents data needed to create/update a work type
type WorkTypeForm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// WorkTypesResponse represents paginated work types response
type WorkTypesResponse struct {
	WorkTypes   []WorkType `json:"work_types"`
	CurrentPage int        `json:"current_page"`
	TotalPages  int        `json:"total_pages"`
	TotalCount  int        `json:"total_count"`
	PageSize    int        `json:"page_size"`
}


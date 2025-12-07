package models

// Procedure represents a dental procedure
type Procedure struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	CreatedAt string `json:"created_at"`
}

// ProcedureForm represents data needed to create/update a procedure
type ProcedureForm struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// ProceduresResponse represents paginated procedures response
type ProceduresResponse struct {
	Procedures  []Procedure `json:"procedures"`
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	TotalCount  int         `json:"total_count"`
	PageSize    int         `json:"page_size"`
}

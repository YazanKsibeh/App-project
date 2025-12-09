package models

// ColorShade represents a tooth color shade
type ColorShade struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HexColor    string `json:"hex_color"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
	CreatedBy   *int   `json:"created_by,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ColorShadeForm represents data needed to create/update a color shade
type ColorShadeForm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HexColor    string `json:"hex_color"`
	IsActive    bool   `json:"is_active"`
}

// ColorShadesResponse represents paginated color shades response
type ColorShadesResponse struct {
	ColorShades []ColorShade `json:"color_shades"`
	CurrentPage int           `json:"current_page"`
	TotalPages  int           `json:"total_pages"`
	TotalCount  int           `json:"total_count"`
	PageSize    int           `json:"page_size"`
}


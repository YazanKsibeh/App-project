package models

// Session struct represents a treatment session
type Session struct {
	ID          int    `json:"id"`
	PatientID   int    `json:"patient_id"`
	DentistID   int    `json:"dentist_id"`
	SessionDate string `json:"session_date"`
	TotalAmount int    `json:"total_amount"`
	Status      string `json:"status"` // "completed" or "in-progress"
	Notes       string `json:"notes"`
	// Related data (populated when needed)
	PatientName   string        `json:"patient_name,omitempty"`
	DentistName   string        `json:"dentist_name,omitempty"`
	InvoiceNumber string        `json:"invoice_number,omitempty"`
	Items         []SessionItem `json:"items,omitempty"`
}

// SessionItem struct represents a procedure/item in a session
type SessionItem struct {
	ID          int    `json:"id"`
	SessionID   int    `json:"session_id"`
	ProcedureID *int   `json:"procedure_id,omitempty"` // nullable
	ItemName    string `json:"item_name"`
	Amount      int    `json:"amount"`
}

// SessionForm represents the data needed to create/update a session
type SessionForm struct {
	PatientID   int               `json:"patient_id"`
	DentistID   int               `json:"dentist_id"`
	SessionDate string            `json:"session_date"`
	Status      string            `json:"status"`
	Notes       string            `json:"notes"`
	Items       []SessionItemForm `json:"items"`
}

// SessionItemForm represents the data needed to create/update a session item
type SessionItemForm struct {
	ProcedureID *int   `json:"procedure_id,omitempty"` // nullable
	ItemName    string `json:"item_name"`
	Amount      int    `json:"amount"`
}

// SessionsResponse represents the response for GetSessions (sessions + pagination info)
type SessionsResponse struct {
	Sessions    []Session `json:"sessions"`
	CurrentPage int       `json:"current_page"`
	TotalPages  int       `json:"total_pages"`
	PageSize    int       `json:"page_size"`
	TotalCount  int       `json:"total_count"`
}

// SessionFilters represents filter criteria for sessions
type SessionFilters struct {
	PatientID    *int    `json:"patient_id,omitempty"`    // Optional patient filter
	Status       *string `json:"status,omitempty"`         // Optional status filter ("completed" or "in-progress")
	DentistID    *int    `json:"dentist_id,omitempty"`    // Optional dentist filter
	DateFrom     *string `json:"date_from,omitempty"`     // Optional start date (YYYY-MM-DD)
	DateTo       *string `json:"date_to,omitempty"`       // Optional end date (YYYY-MM-DD)
	ProcedureIDs []int   `json:"procedure_ids,omitempty"` // Optional procedure IDs (OR logic)
}

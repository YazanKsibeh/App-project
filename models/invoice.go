package models

// Invoice struct represents an invoice
type Invoice struct {
	ID            int    `json:"id"`
	SessionID     int    `json:"session_id"`
	PatientID     int    `json:"patient_id"`
	InvoiceNumber string `json:"invoice_number"`
	InvoiceDate   string `json:"invoice_date"`
	TotalAmount   int    `json:"total_amount"`
	Status        string `json:"status"` // "issued", "paid", "cancelled"
	Notes         string `json:"notes"`
}

// InvoicePreview represents preview data for invoice confirmation
type InvoicePreview struct {
	PatientName   string        `json:"patient_name"`
	SessionDate   string        `json:"session_date"`
	InvoiceNumber string        `json:"invoice_number"`
	Procedures    []SessionItem `json:"procedures"`
	TotalAmount   int           `json:"total_amount"`
}

// InvoiceOverview represents aggregated invoice stats over key time periods
type InvoiceOverview struct {
	TodayTotal int `json:"today_total"`
	TodayCount int `json:"today_count"`
	WeekTotal  int `json:"week_total"`
	WeekCount  int `json:"week_count"`
	MonthTotal int `json:"month_total"`
	MonthCount int `json:"month_count"`
}

// InvoiceListItem represents invoice data shown in tables
type InvoiceListItem struct {
	ID            int    `json:"id"`
	InvoiceNumber string `json:"invoice_number"`
	PatientName   string `json:"patient_name"`
	SessionID     int    `json:"session_id"`
	InvoiceDate   string `json:"invoice_date"`
	TotalAmount   int    `json:"total_amount"`
	Status        string `json:"status"`
}

// InvoiceListResponse provides paginated invoices
type InvoiceListResponse struct {
	Invoices    []InvoiceListItem `json:"invoices"`
	CurrentPage int               `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
	TotalCount  int               `json:"total_count"`
	PageSize    int               `json:"page_size"`
}

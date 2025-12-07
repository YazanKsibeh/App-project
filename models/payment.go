package models

type Payment struct {
	ID            int    `json:"id"`
	InvoiceID     int    `json:"invoice_id"`
	PatientID     int    `json:"patient_id"`
	PaymentCode   string `json:"payment_code"`
	Amount        int    `json:"amount"`
	PaymentDate   string `json:"payment_date"`
	Note          string `json:"note"`
	PaymentMethod string `json:"payment_method"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type PaymentSummary struct {
	TotalPaid       int       `json:"total_paid"`
	RemainingAmount int       `json:"remaining_amount"`
	Payments        []Payment `json:"payments"`
}

type InvoicePaymentDetails struct {
	Invoice       Invoice   `json:"invoice"`
	PatientName   string    `json:"patient_name"`
	Status        string    `json:"status"`
	TotalAmount   int       `json:"total_amount"`
	Payments      []Payment `json:"payments"`
	TotalPaid     int       `json:"total_paid"`
	Remaining     int       `json:"remaining"`
	AllowPayments bool      `json:"allow_payments"`
}

type PaymentListItem struct {
	ID            int    `json:"id"`
	PaymentCode   string `json:"payment_code"`
	PaymentAmount int    `json:"payment_amount"`
	PaymentDate   string `json:"payment_date"`
	Note          string `json:"note"`
	InvoiceID     int    `json:"invoice_id"`
	InvoiceNumber string `json:"invoice_number"`
	InvoiceAmount int    `json:"invoice_amount"`
	InvoiceStatus string `json:"invoice_status"`
	InvoiceDate   string `json:"invoice_date"`
	PatientName   string `json:"patient_name"`
}

type PaymentListResponse struct {
	Payments    []PaymentListItem `json:"payments"`
	CurrentPage int               `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
	TotalCount  int               `json:"total_count"`
	PageSize    int               `json:"page_size"`
}

package models

// LabOrderListItem represents lab order data shown in the orders list table
type LabOrderListItem struct {
	ID            int    `json:"id"`
	OrderNumber   string `json:"order_number"`
	PatientName   string `json:"patient_name"`
	LabName       string `json:"lab_name"`
	WorkTypeName  string `json:"work_type_name"`
	Status        string `json:"status"`
	LabCost       int    `json:"lab_cost"`
	OrderDate     string `json:"order_date"`
}

// LabOrderDetail represents full lab order details for the detail modal
type LabOrderDetail struct {
	ID            int    `json:"id"`
	OrderNumber   string `json:"order_number"`
	PatientID     int    `json:"patient_id"`
	PatientName   string `json:"patient_name"`
	LabID         int    `json:"lab_id"`
	LabName       string `json:"lab_name"`
	CreatedBy     int    `json:"created_by"`
	DentistName   string `json:"dentist_name"`
	WorkTypeID    int    `json:"work_type_id"`
	WorkTypeName  string `json:"work_type_name"`
	Description   string `json:"description"`
	UpperLeft     string `json:"upper_left"`
	UpperRight    string `json:"upper_right"`
	LowerLeft     string `json:"lower_left"`
	LowerRight    string `json:"lower_right"`
	Quantity      int    `json:"quantity"`
	ColorShadeID  *int   `json:"color_shade_id"`
	ColorShadeName string `json:"color_shade_name"`
	LabCost       int    `json:"lab_cost"`
	OrderDate     string `json:"order_date"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// LabOrdersResponse provides paginated lab orders
type LabOrdersResponse struct {
	Orders      []LabOrderListItem `json:"orders"`
	CurrentPage int               `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
	TotalCount  int               `json:"total_count"`
	PageSize    int               `json:"page_size"`
}

// LabOrderForm represents data needed to create a lab order
type LabOrderForm struct {
	PatientID    int    `json:"patient_id"`
	LabID        int    `json:"lab_id"`
	WorkTypeID   int    `json:"work_type_id"`
	ColorShadeID *int   `json:"color_shade_id"`
	Description  string `json:"description"`
	UpperLeft    string `json:"upper_left"`
	UpperRight   string `json:"upper_right"`
	LowerLeft    string `json:"lower_left"`
	LowerRight   string `json:"lower_right"`
	Quantity     int    `json:"quantity"`
	LabCost      int    `json:"lab_cost"`
	OrderDate    string `json:"order_date"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
}

// CreateLabOrderResponse represents the response from creating a lab order
type CreateLabOrderResponse struct {
	ID          int64  `json:"id"`
	OrderNumber string `json:"order_number"`
}


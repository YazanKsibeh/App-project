package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"math"
	"time"
)

// InvoiceHandler handles invoice-related database operations
type InvoiceHandler struct {
	db *sql.DB
}

// NewInvoiceHandler creates a new InvoiceHandler
func NewInvoiceHandler(db *sql.DB) *InvoiceHandler {
	return &InvoiceHandler{db: db}
}

// GenerateInvoiceNumber generates the next sequential invoice number (INV-001, INV-002, etc.)
func (h *InvoiceHandler) GenerateInvoiceNumber() (string, error) {
	// Get the highest invoice number
	query := `SELECT invoice_number FROM invoices 
	          WHERE invoice_number LIKE 'INV-%' 
	          ORDER BY CAST(SUBSTR(invoice_number, 5) AS INTEGER) DESC 
	          LIMIT 1`

	var lastNumber string
	err := h.db.QueryRow(query).Scan(&lastNumber)

	if err == sql.ErrNoRows {
		// No invoices exist yet, start with INV-001
		return "INV-001", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to get last invoice number: %v", err)
	}

	// Extract the number part and increment
	var num int
	_, err = fmt.Sscanf(lastNumber, "INV-%d", &num)
	if err != nil {
		// If parsing fails, start fresh
		return "INV-001", nil
	}

	// Increment and format
	num++
	return fmt.Sprintf("INV-%03d", num), nil
}

// GetInvoiceBySession checks if an invoice exists for a session
func (h *InvoiceHandler) GetInvoiceBySession(sessionID int) (*models.Invoice, error) {
	var invoice models.Invoice

	query := `SELECT id, session_id, patient_id, invoice_number, invoice_date, 
	          total_amount, status, notes 
	          FROM invoices WHERE session_id = ?`

	err := h.db.QueryRow(query, sessionID).Scan(
		&invoice.ID, &invoice.SessionID, &invoice.PatientID,
		&invoice.InvoiceNumber, &invoice.InvoiceDate,
		&invoice.TotalAmount, &invoice.Status, &invoice.Notes)

	if err == sql.ErrNoRows {
		return nil, nil // No invoice found, but not an error
	} else if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %v", err)
	}

	return &invoice, nil
}

// CreateInvoice creates an invoice from a session
func (h *InvoiceHandler) CreateInvoice(sessionID int) (*models.Invoice, error) {
	// First, check if invoice already exists
	existing, err := h.GetInvoiceBySession(sessionID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("invoice already exists for this session")
	}

	// Get session data
	var session models.Session
	sessionQuery := `SELECT id, patient_id, dentist_id, session_date, total_amount, status, notes
	                 FROM sessions WHERE id = ?`

	err = h.db.QueryRow(sessionQuery, sessionID).Scan(
		&session.ID, &session.PatientID, &session.DentistID,
		&session.SessionDate, &session.TotalAmount, &session.Status, &session.Notes)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	// Generate invoice number
	invoiceNumber, err := h.GenerateInvoiceNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %v", err)
	}

	// Create invoice
	invoiceDate := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO invoices (session_id, patient_id, invoice_number, invoice_date, 
	          total_amount, status, notes)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := h.db.Exec(query, sessionID, session.PatientID, invoiceNumber,
		invoiceDate, session.TotalAmount, "issued", session.Notes)

	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %v", err)
	}

	invoiceID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice ID: %v", err)
	}

	// Return the created invoice
	invoice := &models.Invoice{
		ID:            int(invoiceID),
		SessionID:     sessionID,
		PatientID:     session.PatientID,
		InvoiceNumber: invoiceNumber,
		InvoiceDate:   invoiceDate,
		TotalAmount:   session.TotalAmount,
		Status:        "issued",
		Notes:         session.Notes,
	}

	return invoice, nil
}

// PreviewInvoice returns preview data for invoice confirmation without creating the invoice
func (h *InvoiceHandler) PreviewInvoice(sessionID int) (*models.InvoicePreview, error) {
	// Check if invoice already exists
	existing, err := h.GetInvoiceBySession(sessionID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("invoice already exists for this session")
	}

	// Get session data with patient name
	var session models.Session
	sessionQuery := `SELECT s.id, s.patient_id, s.dentist_id, s.session_date, s.total_amount, s.status, s.notes,
	                        p.name as patient_name
	                 FROM sessions s
	                 LEFT JOIN patients p ON s.patient_id = p.id
	                 WHERE s.id = ?`

	err = h.db.QueryRow(sessionQuery, sessionID).Scan(
		&session.ID, &session.PatientID, &session.DentistID,
		&session.SessionDate, &session.TotalAmount, &session.Status, &session.Notes,
		&session.PatientName)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	// Get session items
	itemsQuery := `SELECT id, session_id, procedure_id, item_name, amount
	               FROM session_items WHERE session_id = ? ORDER BY id`
	itemRows, err := h.db.Query(itemsQuery, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session items: %v", err)
	}
	defer itemRows.Close()

	items := make([]models.SessionItem, 0)
	for itemRows.Next() {
		var item models.SessionItem
		var procedureID sql.NullInt64
		err := itemRows.Scan(&item.ID, &item.SessionID, &procedureID, &item.ItemName, &item.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session item: %v", err)
		}
		if procedureID.Valid {
			procID := int(procedureID.Int64)
			item.ProcedureID = &procID
		}
		items = append(items, item)
	}

	// Generate invoice number (without creating invoice)
	invoiceNumber, err := h.GenerateInvoiceNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %v", err)
	}

	// Return preview
	preview := &models.InvoicePreview{
		PatientName:   session.PatientName,
		SessionDate:   session.SessionDate,
		InvoiceNumber: invoiceNumber,
		Procedures:    items,
		TotalAmount:   session.TotalAmount,
	}

	return preview, nil
}

// GetInvoiceOverview aggregates invoice totals and counts over key time ranges
func (h *InvoiceHandler) GetInvoiceOverview() (*models.InvoiceOverview, error) {
	now := time.Now()
	loc := now.Location()

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	weekday := int(todayStart.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := todayStart.AddDate(0, 0, -(weekday - 1))

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)

	todayTotal, todayCount, err := h.invoiceSummary(todayStart, now)
	if err != nil {
		return nil, err
	}

	weekTotal, weekCount, err := h.invoiceSummary(weekStart, now)
	if err != nil {
		return nil, err
	}

	monthTotal, monthCount, err := h.invoiceSummary(monthStart, now)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceOverview{
		TodayTotal: todayTotal,
		TodayCount: todayCount,
		WeekTotal:  weekTotal,
		WeekCount:  weekCount,
		MonthTotal: monthTotal,
		MonthCount: monthCount,
	}, nil
}

func (h *InvoiceHandler) invoiceSummary(start, end time.Time) (int, int, error) {
	query := `SELECT COALESCE(SUM(total_amount), 0) AS total_amount, COUNT(*) AS invoice_count
	          FROM invoices
	          WHERE datetime(invoice_date) >= datetime(?) AND datetime(invoice_date) < datetime(?)`

	startStr := start.Format("2006-01-02 15:04:05")
	endStr := end.Format("2006-01-02 15:04:05")

	var total sql.NullInt64
	var count int

	err := h.db.QueryRow(query, startStr, endStr).Scan(&total, &count)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to calculate invoice summary: %v", err)
	}

	sum := 0
	if total.Valid {
		sum = int(total.Int64)
	}

	return sum, count, nil
}

// GetInvoices returns paginated invoices sorted from most recent to oldest
func (h *InvoiceHandler) GetInvoices(page, pageSize int) (*models.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	var totalCount int
	countQuery := `SELECT COUNT(*) FROM invoices`
	if err := h.db.QueryRow(countQuery).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count invoices: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	query := `SELECT i.id, i.invoice_number, i.session_id, i.invoice_date, i.total_amount, i.status, 
	                 COALESCE(p.name, 'Unknown') AS patient_name
	          FROM invoices i
	          LEFT JOIN patients p ON p.id = i.patient_id
	          ORDER BY datetime(i.invoice_date) DESC, i.id DESC
	          LIMIT ? OFFSET ?`

	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load invoices: %v", err)
	}
	defer rows.Close()

	invoices := make([]models.InvoiceListItem, 0)
	for rows.Next() {
		var item models.InvoiceListItem
		if err := rows.Scan(
			&item.ID,
			&item.InvoiceNumber,
			&item.SessionID,
			&item.InvoiceDate,
			&item.TotalAmount,
			&item.Status,
			&item.PatientName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invoices: %v", err)
		}
		invoices = append(invoices, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("invoice rows error: %v", err)
	}

	response := &models.InvoiceListResponse{
		Invoices:    invoices,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}

	return response, nil
}

// GetInvoicePaymentDetails returns invoice info along with payment history and totals
func (h *InvoiceHandler) GetInvoicePaymentDetails(invoiceID int) (*models.InvoicePaymentDetails, error) {
	return h.fetchInvoicePaymentDetails(h.db, invoiceID)
}

// CreatePayment records a payment for an invoice and updates invoice status/totals
func (h *InvoiceHandler) CreatePayment(invoiceID int, amount int, paymentDateStr string, note string) (*models.InvoicePaymentDetails, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("payment amount must be greater than zero")
	}

	paymentDate, err := parsePaymentDate(paymentDateStr)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if paymentDate.After(now) {
		return nil, fmt.Errorf("payment date cannot be in the future")
	}

	tx, err := h.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var invoice models.Invoice
	var patientName string
	invoiceQuery := `SELECT i.id, i.session_id, i.patient_id, i.invoice_number, i.invoice_date, i.total_amount, i.status, i.notes,
							COALESCE(p.name, 'Unknown') AS patient_name
					 FROM invoices i
					 LEFT JOIN patients p ON p.id = i.patient_id
					 WHERE i.id = ?`

	err = tx.QueryRow(invoiceQuery, invoiceID).Scan(
		&invoice.ID,
		&invoice.SessionID,
		&invoice.PatientID,
		&invoice.InvoiceNumber,
		&invoice.InvoiceDate,
		&invoice.TotalAmount,
		&invoice.Status,
		&invoice.Notes,
		&patientName,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to load invoice: %v", err)
	}

	if invoice.Status == "paid" || invoice.Status == "cancelled" {
		return nil, fmt.Errorf("payments are not allowed for invoices with status %s", invoice.Status)
	}

	var totalPaid int
	err = tx.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM payments WHERE invoice_id = ?`, invoiceID).Scan(&totalPaid)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate previous payments: %v", err)
	}

	remaining := invoice.TotalAmount - totalPaid
	if remaining <= 0 {
		return nil, fmt.Errorf("invoice is already fully paid")
	}
	if amount > remaining {
		return nil, fmt.Errorf("payment exceeds remaining balance")
	}

	paymentCode, err := h.generatePaymentCode(tx)
	if err != nil {
		return nil, err
	}

	paymentDateFormatted := paymentDate.Format("2006-01-02 15:04:05")
	noteValue := sql.NullString{String: note, Valid: note != ""}

	insertQuery := `INSERT INTO payments (invoice_id, patient_id, payment_code, amount, payment_date, note, payment_method, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, ?, 'cash', datetime('now'), datetime('now'))`
	_, err = tx.Exec(insertQuery, invoiceID, invoice.PatientID, paymentCode, amount, paymentDateFormatted, noteValue)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %v", err)
	}

	newTotalPaid := totalPaid + amount
	newStatus := "partially_paid"
	if newTotalPaid >= invoice.TotalAmount {
		newStatus = "paid"
	}

	if newStatus != invoice.Status {
		_, err = tx.Exec(`UPDATE invoices SET status = ? WHERE id = ?`, newStatus, invoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to update invoice status: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit payment: %v", err)
	}

	return h.GetInvoicePaymentDetails(invoiceID)
}

func (h *InvoiceHandler) fetchInvoicePaymentDetails(runner queryRunner, invoiceID int) (*models.InvoicePaymentDetails, error) {
	var invoice models.Invoice
	var patientName string

	query := `SELECT i.id, i.session_id, i.patient_id, i.invoice_number, i.invoice_date, i.total_amount, i.status, i.notes,
	                 COALESCE(p.name, 'Unknown') AS patient_name
	          FROM invoices i
	          LEFT JOIN patients p ON p.id = i.patient_id
	          WHERE i.id = ?`

	err := runner.QueryRow(query, invoiceID).Scan(
		&invoice.ID,
		&invoice.SessionID,
		&invoice.PatientID,
		&invoice.InvoiceNumber,
		&invoice.InvoiceDate,
		&invoice.TotalAmount,
		&invoice.Status,
		&invoice.Notes,
		&patientName,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to load invoice: %v", err)
	}

	payments, totalPaid, err := h.fetchPaymentsForInvoice(runner, invoiceID)
	if err != nil {
		return nil, err
	}

	remaining := invoice.TotalAmount - totalPaid
	if remaining < 0 {
		remaining = 0
	}

	return &models.InvoicePaymentDetails{
		Invoice:       invoice,
		PatientName:   patientName,
		Status:        invoice.Status,
		TotalAmount:   invoice.TotalAmount,
		Payments:      payments,
		TotalPaid:     totalPaid,
		Remaining:     remaining,
		AllowPayments: invoice.Status == "issued" || invoice.Status == "partially_paid",
	}, nil
}

func (h *InvoiceHandler) fetchPaymentsForInvoice(runner queryRunner, invoiceID int) ([]models.Payment, int, error) {
	rows, err := runner.Query(`SELECT id, invoice_id, patient_id, COALESCE(payment_code, ''), amount, 
	                                  COALESCE(payment_date, ''), COALESCE(note, ''), COALESCE(payment_method, 'cash'),
	                                  COALESCE(created_at, ''), COALESCE(updated_at, '')
	                           FROM payments
	                           WHERE invoice_id = ?
	                           ORDER BY datetime(payment_date) DESC, id DESC`, invoiceID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load payments: %v", err)
	}
	defer rows.Close()

	payments := make([]models.Payment, 0)
	totalPaid := 0
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.ID,
			&payment.InvoiceID,
			&payment.PatientID,
			&payment.PaymentCode,
			&payment.Amount,
			&payment.PaymentDate,
			&payment.Note,
			&payment.PaymentMethod,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment: %v", err)
		}
		totalPaid += payment.Amount
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("payment rows error: %v", err)
	}

	return payments, totalPaid, nil
}

func (h *InvoiceHandler) generatePaymentCode(tx *sql.Tx) (string, error) {
	var runner queryRunner = h.db
	if tx != nil {
		runner = tx
	}
	return generatePaymentCodeFromRunner(runner)
}

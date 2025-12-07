package handlers

import (
	"DentistApp/models"
	"database/sql"
	"fmt"
	"math"
	"time"
)

type PaymentHandler struct {
	db *sql.DB
}

func NewPaymentHandler(db *sql.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

// AddPayment adds a new payment for a patient
func (h *PaymentHandler) AddPayment(payment models.Payment) (int64, error) {
	if payment.Amount <= 0 {
		return 0, fmt.Errorf("payment amount must be greater than zero")
	}

	paymentDate, err := parsePaymentDate(payment.PaymentDate)
	if err != nil {
		return 0, err
	}
	if paymentDate.After(time.Now()) {
		return 0, fmt.Errorf("payment date cannot be in the future")
	}

	paymentCode, err := generatePaymentCodeFromRunner(h.db)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO payments (invoice_id, patient_id, payment_code, amount, payment_date, note, payment_method, created_at, updated_at)
			  VALUES (NULL, ?, ?, ?, ?, ?, 'cash', datetime('now'), datetime('now'))`
	result, err := h.db.Exec(query, payment.PatientID, paymentCode, payment.Amount, paymentDate.Format("2006-01-02 15:04:05"), payment.Note)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetPaymentsForPatient returns all payments for a patient, most recent first
func (h *PaymentHandler) GetPaymentsForPatient(patientID int) ([]models.Payment, error) {
	query := `SELECT id,
	                COALESCE(invoice_id, 0),
	                patient_id,
	                COALESCE(payment_code, ''),
	                amount,
	                COALESCE(payment_date, ''),
	                COALESCE(note, ''),
	                COALESCE(payment_method, 'cash'),
	                COALESCE(created_at, ''),
	                COALESCE(updated_at, '')
	         FROM payments
	         WHERE patient_id = ? AND invoice_id IS NULL
	         ORDER BY datetime(payment_date) DESC, id DESC`
	rows, err := h.db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		err := rows.Scan(
			&p.ID,
			&p.InvoiceID,
			&p.PatientID,
			&p.PaymentCode,
			&p.Amount,
			&p.PaymentDate,
			&p.Note,
			&p.PaymentMethod,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

// GetLastPaymentForPatient returns the most recent payment for a patient
func (h *PaymentHandler) GetLastPaymentForPatient(patientID int) (*models.Payment, error) {
	query := `SELECT id,
	                 COALESCE(invoice_id, 0),
	                 patient_id,
	                 COALESCE(payment_code, ''),
	                 amount,
	                 COALESCE(payment_date, ''),
	                 COALESCE(note, ''),
	                 COALESCE(payment_method, 'cash'),
	                 COALESCE(created_at, ''),
	                 COALESCE(updated_at, '')
	          FROM payments
	          WHERE patient_id = ? AND invoice_id IS NULL
	          ORDER BY datetime(payment_date) DESC, id DESC
	          LIMIT 1`
	row := h.db.QueryRow(query, patientID)
	var p models.Payment
	err := row.Scan(
		&p.ID,
		&p.InvoiceID,
		&p.PatientID,
		&p.PaymentCode,
		&p.Amount,
		&p.PaymentDate,
		&p.Note,
		&p.PaymentMethod,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateTotalRequired sets the total required amount for a patient
func (h *PaymentHandler) UpdateTotalRequired(patientID int, total int) error {
	query := `UPDATE patients SET total_required = ? WHERE id = ?`
	_, err := h.db.Exec(query, total, patientID)
	return err
}

// PatientBalance struct for returning patient balance info
type PatientBalance struct {
	TotalRequired int `json:"total_required"`
	TotalPaid     int `json:"total_paid"`
	Remaining     int `json:"remaining"`
}

// GetPatientBalance returns the total required, total paid, and remaining for a patient
func (h *PaymentHandler) GetPatientBalance(patientID int) (*PatientBalance, error) {
	var totalRequired, totalPaid int
	err := h.db.QueryRow(`SELECT total_required FROM patients WHERE id = ?`, patientID).Scan(&totalRequired)
	if err != nil {
		return nil, err
	}
	err = h.db.QueryRow(`SELECT COALESCE(SUM(amount),0) FROM payments WHERE patient_id = ? AND invoice_id IS NULL`, patientID).Scan(&totalPaid)
	if err != nil {
		return nil, err
	}
	remaining := totalRequired - totalPaid
	return &PatientBalance{
		TotalRequired: totalRequired,
		TotalPaid:     totalPaid,
		Remaining:     remaining,
	}, nil
}

// DeletePayment deletes a payment by ID
func (h *PaymentHandler) DeletePayment(paymentID int) error {
	query := `DELETE FROM payments WHERE id = ? AND invoice_id IS NULL`
	_, err := h.db.Exec(query, paymentID)
	return err
}

// UpdatePayment updates a payment by ID
func (h *PaymentHandler) UpdatePayment(payment models.Payment) error {
	if payment.Amount <= 0 {
		return fmt.Errorf("payment amount must be greater than zero")
	}

	paymentDate, err := parsePaymentDate(payment.PaymentDate)
	if err != nil {
		return err
	}
	if paymentDate.After(time.Now()) {
		return fmt.Errorf("payment date cannot be in the future")
	}

	query := `UPDATE payments
	          SET amount = ?, payment_date = ?, note = ?, updated_at = datetime('now')
	          WHERE id = ? AND invoice_id IS NULL`
	_, err = h.db.Exec(query, payment.Amount, paymentDate.Format("2006-01-02 15:04:05"), payment.Note, payment.ID)
	return err
}

// GetInvoicePayments returns paginated payments tied to invoices
func (h *PaymentHandler) GetInvoicePayments(page, pageSize int) (*models.PaymentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var totalCount int
	if err := h.db.QueryRow(`SELECT COUNT(*) FROM payments WHERE invoice_id IS NOT NULL`).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count payments: %v", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * pageSize

	query := `SELECT
				p.id,
				COALESCE(p.payment_code, ''),
				p.amount,
				COALESCE(p.payment_date, ''),
				COALESCE(p.note, ''),
				COALESCE(i.id, 0),
				COALESCE(i.invoice_number, ''),
				COALESCE(i.total_amount, 0),
				COALESCE(i.status, ''),
				COALESCE(i.invoice_date, ''),
				COALESCE(pt.name, 'Unknown')
			  FROM payments p
			  LEFT JOIN invoices i ON i.id = p.invoice_id
			  LEFT JOIN patients pt ON pt.id = p.patient_id
			  WHERE p.invoice_id IS NOT NULL
			  ORDER BY datetime(p.payment_date) DESC, p.id DESC
			  LIMIT ? OFFSET ?`

	rows, err := h.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to load payments: %v", err)
	}
	defer rows.Close()

	items := make([]models.PaymentListItem, 0)
	for rows.Next() {
		var item models.PaymentListItem
		if err := rows.Scan(
			&item.ID,
			&item.PaymentCode,
			&item.PaymentAmount,
			&item.PaymentDate,
			&item.Note,
			&item.InvoiceID,
			&item.InvoiceNumber,
			&item.InvoiceAmount,
			&item.InvoiceStatus,
			&item.InvoiceDate,
			&item.PatientName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payments: %v", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("payment rows error: %v", err)
	}

	return &models.PaymentListResponse{
		Payments:    items,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		PageSize:    pageSize,
	}, nil
}

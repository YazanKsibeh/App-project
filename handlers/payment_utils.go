package handlers

import (
	"database/sql"
	"fmt"
	"time"
)

type queryRunner interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

const paymentCodePrefix = "Payment-"

func generatePaymentCodeFromRunner(runner queryRunner) (string, error) {
	query := `SELECT payment_code FROM payments 
	          WHERE payment_code LIKE 'Payment-%'
	          ORDER BY CAST(SUBSTR(payment_code, LENGTH('Payment-') + 1) AS INTEGER) DESC
	          LIMIT 1`

	var lastCode string
	err := runner.QueryRow(query).Scan(&lastCode)
	if err == sql.ErrNoRows {
		return fmt.Sprintf("%s%03d", paymentCodePrefix, 1), nil
	} else if err != nil {
		return "", fmt.Errorf("failed to fetch last payment code: %v", err)
	}

	var num int
	_, err = fmt.Sscanf(lastCode, paymentCodePrefix+"%d", &num)
	if err != nil {
		return fmt.Sprintf("%s%03d", paymentCodePrefix, 1), nil
	}

	num++
	return fmt.Sprintf("%s%03d", paymentCodePrefix, num), nil
}

func parsePaymentDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("payment date is required")
	}

	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid payment date format")
}

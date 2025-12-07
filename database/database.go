package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database and creates tables
func InitDB() (*sql.DB, error) {
	// Open database with foreign key support enabled in connection string
	db, err := sql.Open("sqlite", "./dentist.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	// Ensure foreign key constraints are enabled
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// Verify foreign keys are enabled
	var fkEnabled int
	err = db.QueryRow("PRAGMA foreign_keys;").Scan(&fkEnabled)
	if err != nil {
		return nil, err
	}
	if fkEnabled != 1 {
		return nil, fmt.Errorf("foreign keys could not be enabled")
	}

	// Create patients table
	createPatientsTable := `
	CREATE TABLE IF NOT EXISTS patients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		occupation TEXT,
		total_required INTEGER DEFAULT 0
	);`

	_, err = db.Exec(createPatientsTable)
	if err != nil {
		return nil, err
	}

	// Add total_required column if it doesn't exist (for migrations)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN total_required INTEGER DEFAULT 0;`)

	// Migration: Remove occupation column (SQLite doesn't support DROP COLUMN directly, so we'll handle it in queries)
	// Note: We'll just stop using occupation in queries, existing data will remain but won't be accessed

	// Migration: Add medical history columns if they don't exist
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN allergies TEXT;`)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN current_medications TEXT;`)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN medical_conditions TEXT;`)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN smoking_status INTEGER DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN pregnancy_status INTEGER DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN dental_history TEXT;`)
	_, _ = db.Exec(`ALTER TABLE patients ADD COLUMN special_notes TEXT;`)

	// Create appointments table
	createAppointmentsTable := `
	CREATE TABLE IF NOT EXISTS appointments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		patient_id INTEGER NOT NULL,
		datetime TEXT NOT NULL,
		duration INTEGER,
		notes TEXT,
		FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE
	);`

	_, err = db.Exec(createAppointmentsTable)
	if err != nil {
		return nil, err
	}

	// Create payments table
	createPaymentsTable := `
	CREATE TABLE IF NOT EXISTS payments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		invoice_id INTEGER,
		patient_id INTEGER NOT NULL,
		payment_code TEXT UNIQUE,
		amount INTEGER NOT NULL,
		payment_date TEXT NOT NULL DEFAULT (datetime('now')),
		note TEXT,
		payment_method TEXT NOT NULL DEFAULT 'cash',
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now')),
		FOREIGN KEY(invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
		FOREIGN KEY(patient_id) REFERENCES patients(id)
	);`

	_, err = db.Exec(createPaymentsTable)
	if err != nil {
		return nil, err
	}

	// Migration attempts for legacy payments table
	_, _ = db.Exec(`ALTER TABLE payments ADD COLUMN invoice_id INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE payments ADD COLUMN payment_code TEXT;`)
	_, _ = db.Exec(`ALTER TABLE payments ADD COLUMN payment_date TEXT DEFAULT (datetime('now'));`)
	_, _ = db.Exec(`ALTER TABLE payments ADD COLUMN payment_method TEXT DEFAULT 'cash';`)
	_, _ = db.Exec(`ALTER TABLE payments ADD COLUMN created_at TEXT DEFAULT (datetime('now'));`)
	_, _ = db.Exec(`ALTER TABLE payments ADD COLUMN updated_at TEXT DEFAULT (datetime('now'));`)

	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'Dentist'
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		return nil, err
	}

	// Create dental procedures table
	createProceduresTable := `
	CREATE TABLE IF NOT EXISTS dental_procedures (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime('now'))
	);`

	_, err = db.Exec(createProceduresTable)
	if err != nil {
		return nil, err
	}

	// Create sessions table
	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		patient_id INTEGER NOT NULL,
		dentist_id INTEGER NOT NULL,
		session_date TEXT NOT NULL DEFAULT (datetime('now')),
		total_amount INTEGER DEFAULT 0,
		status TEXT DEFAULT 'completed',
		notes TEXT,
		FOREIGN KEY(patient_id) REFERENCES patients(id) ON DELETE CASCADE,
		FOREIGN KEY(dentist_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createSessionsTable)
	if err != nil {
		return nil, err
	}

	// Create session_items table
	createSessionItemsTable := `
	CREATE TABLE IF NOT EXISTS session_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER NOT NULL,
		procedure_id INTEGER,
		item_name TEXT NOT NULL,
		amount INTEGER NOT NULL,
		FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE,
		FOREIGN KEY(procedure_id) REFERENCES dental_procedures(id)
	);`

	_, err = db.Exec(createSessionItemsTable)
	if err != nil {
		return nil, err
	}

	// Create invoices table
	createInvoicesTable := `
	CREATE TABLE IF NOT EXISTS invoices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id INTEGER UNIQUE NOT NULL,
		patient_id INTEGER NOT NULL,
		invoice_number TEXT UNIQUE,
		invoice_date TEXT NOT NULL DEFAULT (datetime('now')),
		total_amount INTEGER NOT NULL,
		status TEXT DEFAULT 'issued',
		notes TEXT,
		FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE,
		FOREIGN KEY(patient_id) REFERENCES patients(id)
	);`

	_, err = db.Exec(createInvoicesTable)
	if err != nil {
		return nil, err
	}

	// Create expense_categories table
	createExpenseCategoriesTable := `
	CREATE TABLE IF NOT EXISTS expense_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		color TEXT DEFAULT '#3498db',
		budget_amount INTEGER DEFAULT 0,
		budget_period TEXT NOT NULL DEFAULT 'monthly' CHECK(budget_period IN ('monthly', 'quarterly', 'yearly')),
		expense_type TEXT NOT NULL DEFAULT 'operational' CHECK(expense_type IN ('operational', 'capital', 'personnel', 'marketing', 'administrative')),
		is_tax_deductible BOOLEAN DEFAULT 1,
		cost_center TEXT DEFAULT 'main',
		account_code TEXT UNIQUE,
		parent_category_id INTEGER,
		is_active BOOLEAN DEFAULT 1,
		requires_approval BOOLEAN DEFAULT 0,
		approval_threshold INTEGER DEFAULT 0,
		reporting_group TEXT,
		sort_order INTEGER DEFAULT 0,
		created_by INTEGER,
		updated_by INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(parent_category_id) REFERENCES expense_categories(id)
	);`

	_, err = db.Exec(createExpenseCategoriesTable)
	if err != nil {
		return nil, err
	}

	// Migration attempts for expense_categories table (in case table exists but columns are missing)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN description TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN color TEXT DEFAULT '#3498db';`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN budget_amount INTEGER DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN budget_period TEXT NOT NULL DEFAULT 'monthly';`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN expense_type TEXT NOT NULL DEFAULT 'operational';`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN is_tax_deductible BOOLEAN DEFAULT 1;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN cost_center TEXT DEFAULT 'main';`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN account_code TEXT UNIQUE;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN parent_category_id INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN is_active BOOLEAN DEFAULT 1;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN requires_approval BOOLEAN DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN approval_threshold INTEGER DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN reporting_group TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN sort_order INTEGER DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN created_by INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN updated_by INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)
	_, _ = db.Exec(`ALTER TABLE expense_categories ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)

	// Create expenses table
	createExpensesTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expense_code TEXT UNIQUE,
		expense_date TEXT NOT NULL DEFAULT (datetime('now')),
		description TEXT NOT NULL,
		amount INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		payment_status TEXT DEFAULT 'unpaid' CHECK(payment_status IN ('unpaid', 'paid', 'partially_paid')),
		payment_method TEXT DEFAULT 'cash' CHECK(payment_method IN ('cash', 'bank_transfer', 'check', 'card')),
		vendor_name TEXT,
		vendor_contact TEXT,
		receipt_number TEXT,
		notes TEXT,
		receipt_file_path TEXT,
		reporting_period TEXT,
		is_recurring BOOLEAN DEFAULT 0,
		recurring_period TEXT,
		created_by INTEGER,
		updated_by INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(category_id) REFERENCES expense_categories(id),
		FOREIGN KEY(created_by) REFERENCES users(id),
		FOREIGN KEY(updated_by) REFERENCES users(id)
	);`

	_, err = db.Exec(createExpensesTable)
	if err != nil {
		return nil, err
	}

	// Migration attempts for expenses table (in case table exists but columns are missing)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN expense_code TEXT UNIQUE;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN expense_date TEXT NOT NULL DEFAULT (datetime('now'));`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN description TEXT NOT NULL;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN amount INTEGER NOT NULL;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN category_id INTEGER NOT NULL;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN payment_status TEXT DEFAULT 'unpaid';`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN payment_method TEXT DEFAULT 'cash';`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN vendor_name TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN vendor_contact TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN receipt_number TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN notes TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN receipt_file_path TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN reporting_period TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN is_recurring BOOLEAN DEFAULT 0;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN recurring_period TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN created_by INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN updated_by INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)
	_, _ = db.Exec(`ALTER TABLE expenses ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)

	// Create expense_payments table
	createExpensePaymentsTable := `
	CREATE TABLE IF NOT EXISTS expense_payments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expense_id INTEGER NOT NULL,
		payment_code TEXT UNIQUE,
		amount INTEGER NOT NULL,
		payment_date TEXT NOT NULL DEFAULT (datetime('now')),
		payment_method TEXT NOT NULL DEFAULT 'cash' CHECK(payment_method IN ('cash', 'bank_transfer', 'check', 'card')),
		note TEXT,
		created_by INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
		FOREIGN KEY(created_by) REFERENCES users(id)
	);`

	_, err = db.Exec(createExpensePaymentsTable)
	if err != nil {
		return nil, err
	}

	// Migration attempts for expense_payments table (in case table exists but columns are missing)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN expense_id INTEGER NOT NULL;`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN payment_code TEXT UNIQUE;`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN amount INTEGER NOT NULL;`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN payment_date TEXT NOT NULL DEFAULT (datetime('now'));`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN payment_method TEXT NOT NULL DEFAULT 'cash';`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN note TEXT;`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN created_by INTEGER;`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)
	_, _ = db.Exec(`ALTER TABLE expense_payments ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)

	// Create patient_data directory if it doesn't exist
	err = os.MkdirAll("patient_data", 0755)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// EnsureForeignKeys ensures foreign keys are enabled for the given database connection
func EnsureForeignKeys(db *sql.DB) error {
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	return err
}

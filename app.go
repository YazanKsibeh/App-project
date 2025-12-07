package main

import (
	"context"
	"fmt"

	"DentistApp/handlers"
	"DentistApp/models"
)

// App struct
type App struct {
	ctx                   context.Context
	patientHandler        *handlers.PatientHandler
	appointmentHandler    *handlers.AppointmentHandler
	paymentHandler        *handlers.PaymentHandler
	procedureHandler      *handlers.ProcedureHandler
	sessionHandler        *handlers.SessionHandler
	invoiceHandler        *handlers.InvoiceHandler
	expenseCategoryHandler *handlers.ExpenseCategoryHandler
	licenseService        *handlers.LicenseService
	authHandler           *handlers.AuthHandler
}

// NewApp creates a new App application struct
func NewApp(patientHandler *handlers.PatientHandler, appointmentHandler *handlers.AppointmentHandler, paymentHandler *handlers.PaymentHandler, procedureHandler *handlers.ProcedureHandler, sessionHandler *handlers.SessionHandler, invoiceHandler *handlers.InvoiceHandler, expenseCategoryHandler *handlers.ExpenseCategoryHandler, authHandler *handlers.AuthHandler) *App {
	return &App{
		patientHandler:        patientHandler,
		appointmentHandler:    appointmentHandler,
		paymentHandler:        paymentHandler,
		procedureHandler:      procedureHandler,
		sessionHandler:        sessionHandler,
		invoiceHandler:        invoiceHandler,
		expenseCategoryHandler: expenseCategoryHandler,
		licenseService:        handlers.NewLicenseService(),
		authHandler:           authHandler,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is shutting down
func (a *App) shutdown(ctx context.Context) {
	// a.db.Close() // The database is closed in main.go
}

// License Management Methods

// ValidateLicense validates a license key and returns license information
func (a *App) ValidateLicense(licenseKey string) (*handlers.LicenseInfo, error) {
	return a.licenseService.ValidateLicense(licenseKey)
}

// IsLicenseValid checks if a license key is valid
func (a *App) IsLicenseValid(licenseKey string) bool {
	return a.licenseService.IsLicenseValid(licenseKey)
}

// checkLicense is a helper method to validate license before operations
func (a *App) checkLicense(licenseKey string) error {
	if !a.licenseService.IsLicenseValid(licenseKey) {
		return fmt.Errorf("invalid or expired license key")
	}
	return nil
}

// Authentication Methods

// Login authenticates a user and returns a session token
func (a *App) Login(username, password, licenseKey string) (*models.LoginResponse, error) {
	// First validate license
	if err := a.checkLicense(licenseKey); err != nil {
		return &models.LoginResponse{
			Success: false,
			Message: "Invalid or expired license key",
		}, nil
	}

	// Then authenticate user
	return a.authHandler.Login(username, password)
}

// CreateUser creates a new user (admin only)
func (a *App) CreateUser(userForm models.UserForm, createdByID int, licenseKey string) (int64, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return 0, err
	}
	return a.authHandler.CreateUser(userForm, createdByID)
}

// GetCurrentUser returns the current user by ID
func (a *App) GetCurrentUser(userID int, licenseKey string) (*models.User, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.authHandler.GetUserByID(userID)
}

// GetAllUsers returns all users (admin only)
func (a *App) GetAllUsers(licenseKey string) ([]models.User, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.authHandler.GetAllUsers()
}

// Patient Management Methods

// AddPatient adds a new patient
func (a *App) AddPatient(patient models.PatientForm, licenseKey string) (int64, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return 0, err
	}
	return a.patientHandler.AddPatient(patient)
}

// GetPatients returns all patients
func (a *App) GetPatients(licenseKey string) ([]models.Patient, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.patientHandler.GetPatients()
}

// GetPatient returns a specific patient by ID
func (a *App) GetPatient(id int, licenseKey string) (models.Patient, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return models.Patient{}, err
	}
	return a.patientHandler.GetPatient(id)
}

// UpdatePatient updates an existing patient
func (a *App) UpdatePatient(patient models.Patient, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.patientHandler.UpdatePatient(patient)
}

// DeletePatient deletes a patient
func (a *App) DeletePatient(id int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.patientHandler.DeletePatient(id)
}

// SearchPatients searches patients by name or phone
func (a *App) SearchPatients(searchTerm string, licenseKey string) ([]models.Patient, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.patientHandler.SearchPatients(searchTerm)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string, licenseKey string) string {
	if err := a.checkLicense(licenseKey); err != nil {
		return "License validation required"
	}
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// OpenPatientFolder opens the patient's folder in the system file explorer
func (a *App) OpenPatientFolder(id int, licenseKey string) error {
	fmt.Printf("[APP] OpenPatientFolder called - Patient ID: %d\n", id)
	fmt.Printf("[APP] License key length: %d\n", len(licenseKey))
	fmt.Printf("[APP] License key (first 20 chars): %s...\n",
		func() string {
			if len(licenseKey) >= 20 {
				return licenseKey[:20]
			}
			return licenseKey
		}())

	if err := a.checkLicense(licenseKey); err != nil {
		fmt.Printf("[APP] License validation FAILED: %v\n", err)
		return err
	}

	fmt.Printf("[APP] License validation PASSED, calling handler...\n")
	err := a.patientHandler.OpenPatientFolder(id)
	if err != nil {
		fmt.Printf("[APP] Handler returned ERROR: %v\n", err)
	} else {
		fmt.Printf("[APP] Handler returned SUCCESS\n")
	}
	return err
}

// DeleteAllPatients deletes all patients
func (a *App) DeleteAllPatients(licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.patientHandler.DeleteAllPatients()
}

// Appointment Management Methods

func (a *App) AddAppointment(appt models.Appointment, licenseKey string) (int64, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return 0, err
	}
	return a.appointmentHandler.AddAppointment(appt)
}

func (a *App) GetAppointments(licenseKey string) ([]models.Appointment, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.appointmentHandler.GetAppointments()
}

func (a *App) GetAppointment(id int, licenseKey string) (models.Appointment, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return models.Appointment{}, err
	}
	return a.appointmentHandler.GetAppointment(id)
}

func (a *App) UpdateAppointment(appt models.Appointment, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.appointmentHandler.UpdateAppointment(appt)
}

func (a *App) DeleteAppointment(id int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.appointmentHandler.DeleteAppointment(id)
}

// Payment Management Methods

func (a *App) AddPayment(payment models.Payment, licenseKey string) (int64, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return 0, err
	}
	return a.paymentHandler.AddPayment(payment)
}

func (a *App) GetPaymentsForPatient(patientID int, licenseKey string) ([]models.Payment, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.paymentHandler.GetPaymentsForPatient(patientID)
}

func (a *App) GetLastPaymentForPatient(patientID int, licenseKey string) (*models.Payment, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.paymentHandler.GetLastPaymentForPatient(patientID)
}

func (a *App) UpdateTotalRequired(patientID int, total int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.paymentHandler.UpdateTotalRequired(patientID, total)
}

// GetPatientBalance returns the total required, total paid, and remaining for a patient
func (a *App) GetPatientBalance(patientID int, licenseKey string) (*handlers.PatientBalance, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.paymentHandler.GetPatientBalance(patientID)
}

// DeletePayment deletes a payment by ID
func (a *App) DeletePayment(paymentID int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.paymentHandler.DeletePayment(paymentID)
}

// UpdatePayment updates a payment by ID
func (a *App) UpdatePayment(payment models.Payment, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.paymentHandler.UpdatePayment(payment)
}

// Dental Procedure Management Methods

func (a *App) CreateProcedure(procedure models.ProcedureForm, licenseKey string) (int64, error) {
	fmt.Printf("[App] CreateProcedure called with name: %s, price: %d\n", procedure.Name, procedure.Price)
	if err := a.checkLicense(licenseKey); err != nil {
		fmt.Printf("[App] License check failed: %v\n", err)
		return 0, err
	}
	fmt.Println("[App] License check passed, calling procedureHandler.CreateProcedure")
	return a.procedureHandler.CreateProcedure(procedure)
}

func (a *App) GetProcedures(licenseKey string) ([]models.Procedure, error) {
	// fmt.Printf("[App] GetProcedures() called with license key length: %d\n", len(licenseKey))

	if err := a.checkLicense(licenseKey); err != nil {
		// fmt.Printf("[App] GetProcedures() - License check FAILED: %v\n", err)
		return nil, err
	}

	// fmt.Printf("[App] GetProcedures() - License check PASSED, calling handler...\n")
	result, err := a.procedureHandler.GetProcedures()
	// if err != nil {
	// 	fmt.Printf("[App] GetProcedures() - Handler returned ERROR: %v\n", err)
	// } else {
	// 	fmt.Printf("[App] GetProcedures() - Handler returned SUCCESS, count: %d\n", len(result))
	// }
	return result, err
}

// GetProceduresPaginated returns paginated procedures
func (a *App) GetProceduresPaginated(page, pageSize int, licenseKey string) (*models.ProceduresResponse, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.procedureHandler.GetProceduresPaginated(page, pageSize)
}

func (a *App) UpdateProcedure(procedure models.Procedure, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.procedureHandler.UpdateProcedure(procedure)
}

func (a *App) DeleteProcedure(id int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.procedureHandler.DeleteProcedure(id)
}

// Expense Category Management Methods

// CreateExpenseCategory creates a new expense category
func (a *App) CreateExpenseCategory(category models.ExpenseCategoryForm, userID int, licenseKey string) (int64, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return 0, err
	}
	return a.expenseCategoryHandler.CreateExpenseCategory(category, userID)
}

// GetExpenseCategories returns all active expense categories
func (a *App) GetExpenseCategories(licenseKey string) ([]models.ExpenseCategory, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.expenseCategoryHandler.GetExpenseCategories()
}

// GetExpenseCategoriesPaginated returns paginated expense categories
func (a *App) GetExpenseCategoriesPaginated(page, pageSize int, licenseKey string) (*models.ExpenseCategoriesResponse, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.expenseCategoryHandler.GetExpenseCategoriesPaginated(page, pageSize)
}

// UpdateExpenseCategory updates an expense category
func (a *App) UpdateExpenseCategory(id int, category models.ExpenseCategoryForm, userID int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.expenseCategoryHandler.UpdateExpenseCategory(id, category, userID)
}

// DeleteExpenseCategory deletes an expense category (soft delete)
func (a *App) DeleteExpenseCategory(id int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.expenseCategoryHandler.DeleteExpenseCategory(id)
}

// PermanentlyDeleteExpenseCategory permanently deletes an expense category (hard delete)
func (a *App) PermanentlyDeleteExpenseCategory(id int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.expenseCategoryHandler.PermanentlyDeleteExpenseCategory(id)
}

// Session Management Methods

// CreateSession creates a new session
func (a *App) CreateSession(session models.SessionForm, licenseKey string) (int64, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return 0, err
	}
	return a.sessionHandler.CreateSession(session)
}

// GetSessions returns paginated sessions
func (a *App) GetSessions(page int, filters models.SessionFilters, licenseKey string) (models.SessionsResponse, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return models.SessionsResponse{}, err
	}
	// Convert empty filters to nil for backward compatibility
	var filterPtr *models.SessionFilters
	if filters.PatientID != nil || filters.Status != nil || filters.DentistID != nil ||
		filters.DateFrom != nil || filters.DateTo != nil || len(filters.ProcedureIDs) > 0 {
		filterPtr = &filters
	}
	return a.sessionHandler.GetSessions(page, filterPtr)
}

// GetSession returns a specific session by ID
func (a *App) GetSession(id int, licenseKey string) (models.Session, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return models.Session{}, err
	}
	return a.sessionHandler.GetSession(id)
}

// UpdateSession updates an existing session
func (a *App) UpdateSession(session models.Session, items []models.SessionItemForm, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.sessionHandler.UpdateSession(session, items)
}

// DeleteSession deletes a session
func (a *App) DeleteSession(id int, licenseKey string) error {
	if err := a.checkLicense(licenseKey); err != nil {
		return err
	}
	return a.sessionHandler.DeleteSession(id)
}

// Invoice Management Methods

// CreateInvoice creates an invoice from a session
func (a *App) CreateInvoice(sessionID int, licenseKey string) (*models.Invoice, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.CreateInvoice(sessionID)
}

// GetInvoiceBySession gets an invoice by session ID
func (a *App) GetInvoiceBySession(sessionID int, licenseKey string) (*models.Invoice, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.GetInvoiceBySession(sessionID)
}

// PreviewInvoice returns preview data for invoice confirmation
func (a *App) PreviewInvoice(sessionID int, licenseKey string) (*models.InvoicePreview, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.PreviewInvoice(sessionID)
}

// GetInvoiceOverview returns aggregated invoice stats for today, this week, and this month
func (a *App) GetInvoiceOverview(licenseKey string) (*models.InvoiceOverview, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.GetInvoiceOverview()
}

// GetInvoices returns paginated invoices for the financials dashboard
func (a *App) GetInvoices(page int, pageSize int, licenseKey string) (*models.InvoiceListResponse, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.GetInvoices(page, pageSize)
}

// GetInvoicePaymentDetails returns invoice payment summary and history
func (a *App) GetInvoicePaymentDetails(invoiceID int, licenseKey string) (*models.InvoicePaymentDetails, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.GetInvoicePaymentDetails(invoiceID)
}

// CreateInvoicePayment records a payment for an invoice
func (a *App) CreateInvoicePayment(invoiceID int, amount int, paymentDate string, note string, licenseKey string) (*models.InvoicePaymentDetails, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.invoiceHandler.CreatePayment(invoiceID, amount, paymentDate, note)
}

// GetInvoicePayments returns paginated payments linked to invoices
func (a *App) GetInvoicePayments(page int, pageSize int, licenseKey string) (*models.PaymentListResponse, error) {
	if err := a.checkLicense(licenseKey); err != nil {
		return nil, err
	}
	return a.paymentHandler.GetInvoicePayments(page, pageSize)
}

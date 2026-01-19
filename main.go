package main

import (
	"context"
	"embed"
	"log"

	"DentistApp/database"
	"DentistApp/handlers"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	patientHandler := handlers.NewPatientHandler(db)
	appointmentHandler := handlers.NewAppointmentHandler(db)
	paymentHandler := handlers.NewPaymentHandler(db)
	procedureHandler := handlers.NewProcedureHandler(db)
	sessionHandler := handlers.NewSessionHandler(db)
	invoiceHandler := handlers.NewInvoiceHandler(db)
	expenseCategoryHandler := handlers.NewExpenseCategoryHandler(db)
	workTypeHandler := handlers.NewWorkTypeHandler(db)
	colorShadeHandler := handlers.NewColorShadeHandler(db)
	dentalLabHandler := handlers.NewDentalLabHandler(db)
	labOrderHandler := handlers.NewLabOrderHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// Initialize admin user if it doesn't exist
	err = authHandler.InitializeAdmin()
	if err != nil {
		log.Printf("Warning: Failed to initialize admin user: %v", err)
	}

	// Create an instance of the app structure
	app := NewApp(patientHandler, appointmentHandler, paymentHandler, procedureHandler, sessionHandler, invoiceHandler, expenseCategoryHandler, workTypeHandler, colorShadeHandler, dentalLabHandler, labOrderHandler, authHandler)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "DentistApp",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

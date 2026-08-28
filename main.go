package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cargo_management_app/database"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	ctx context.Context
	db  *database.Database
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Initialize database
	dbPath := getDBPath()
	log.Printf("Database path: %s", dbPath)
	
	var err error
	a.db, err = database.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	log.Println("Application started successfully")
}

// shutdown is called when the app is shutting down
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
	log.Println("Application shut down successfully")
}

// getDBPath returns the database file path
func getDBPath() string {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// Create .cargo_management_app directory if it doesn't exist
	appDataDir := filepath.Join(homeDir, ".cargo_management_app")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		log.Fatalf("Failed to create app data directory: %v", err)
	}

	return filepath.Join(appDataDir, "cargo_management.db")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, it's been fun!", name)
}

// ============ CARGO METHODS ============

// CreateCargo adds a new cargo item
func (a *App) CreateCargo(name, status, destination, date string, weight float64) (map[string]interface{}, error) {
	cargo, err := a.db.CreateCargo(name, status, destination, date, weight)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":          cargo.ID,
		"name":        cargo.Name,
		"status":      cargo.Status,
		"weight":      cargo.Weight,
		"destination": cargo.Destination,
		"date":        cargo.Date,
	}, nil
}

// GetAllCargo retrieves all cargo items
func (a *App) GetAllCargo() ([]map[string]interface{}, error) {
	cargoList, err := a.db.GetAllCargo()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, cargo := range cargoList {
		result = append(result, map[string]interface{}{
			"id":          cargo.ID,
			"name":        cargo.Name,
			"status":      cargo.Status,
			"weight":      cargo.Weight,
			"destination": cargo.Destination,
			"date":        cargo.Date,
			"created_at":  cargo.CreatedAt,
			"updated_at":  cargo.UpdatedAt,
		})
	}

	return result, nil
}

// GetCargoByID retrieves a specific cargo by ID
func (a *App) GetCargoByID(id int) (map[string]interface{}, error) {
	cargo, err := a.db.GetCargoByID(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          cargo.ID,
		"name":        cargo.Name,
		"status":      cargo.Status,
		"weight":      cargo.Weight,
		"destination": cargo.Destination,
		"date":        cargo.Date,
		"created_at":  cargo.CreatedAt,
		"updated_at":  cargo.UpdatedAt,
	}, nil
}

// UpdateCargo updates an existing cargo item
func (a *App) UpdateCargo(id int, name, status, destination, date string, weight float64) error {
	return a.db.UpdateCargo(id, name, status, destination, date, weight)
}

// DeleteCargo removes a cargo item
func (a *App) DeleteCargo(id int) error {
	return a.db.DeleteCargo(id)
}

// GetCargoStats returns statistics about cargo
func (a *App) GetCargoStats() (map[string]interface{}, error) {
	return a.db.GetCargoStats()
}

// ============ SHIPMENT METHODS ============

// CreateShipment creates a new shipment
func (a *App) CreateShipment(cargoID int, origin, destination, departureDate, arrivalDate, status string) (map[string]interface{}, error) {
	shipment, err := a.db.CreateShipment(cargoID, origin, destination, departureDate, arrivalDate, status)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":             shipment.ID,
		"cargo_id":       shipment.CargoID,
		"origin":         shipment.Origin,
		"destination":    shipment.Destination,
		"departure_date": shipment.DepartureDate,
		"arrival_date":   shipment.ArrivalDate,
		"status":         shipment.Status,
	}, nil
}

// GetShipmentsByCargoID retrieves all shipments for a cargo
func (a *App) GetShipmentsByCargoID(cargoID int) ([]map[string]interface{}, error) {
	shipments, err := a.db.GetShipmentsByCargoID(cargoID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, shipment := range shipments {
		result = append(result, map[string]interface{}{
			"id":             shipment.ID,
			"cargo_id":       shipment.CargoID,
			"origin":         shipment.Origin,
			"destination":    shipment.Destination,
			"departure_date": shipment.DepartureDate,
			"arrival_date":   shipment.ArrivalDate,
			"status":         shipment.Status,
			"created_at":     shipment.CreatedAt,
			"updated_at":     shipment.UpdatedAt,
		})
	}

	return result, nil
}

// UpdateShipmentStatus updates the status of a shipment
func (a *App) UpdateShipmentStatus(id int, status string) error {
	return a.db.UpdateShipmentStatus(id, status)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Cargo Management App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		panic(err)
	}
}

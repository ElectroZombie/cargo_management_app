package main

import (
	"context"
	"embed"
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

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type App struct {
	ctx        context.Context
	db        *database.Database
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	path := getDBPath()
	var err error
	a.db, err = database.NewDatabase(path)
	if err != nil {
		log.Fatalf("initialize database: %v", err)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

func getDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(home, ".cargo_management_app")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, "cargo_management.db")
}

func (a *App) startup(ctx context.Context) {
	a.Startup(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	a.Shutdown(ctx)
}

// ---------- Driver methods ----------

func (a *App) ListDrivers(filters map[string]interface{}, page int, pageSize int, sortBy string, sortDir string) Response {
	// TODO: wire through to database.ListDrivers(filters, database.QueryOptions{page, pageSize, sortBy, sortDir})
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	return Response{Success: true, Data: []database.Driver{}}
}

func (a *App) GetDriver(id int) Response {
	driver, err := a.db.GetDriver(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Response{Success: false, Error: "Driver not found"}
		}
		return Response{Success: false, Error: "Failed to fetch driver"}
	}
	return Response{Success: true, Data: driver}
}

func (a *App) CreateDriver(dto database.DriverDTO) Response {
	driver, err := a.db.CreateDriver(dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: driver}
}

func (a *App) UpdateDriver(id int, dto database.DriverDTO) Response {
	err := a.db.UpdateDriver(id, dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	// return updated driver
	driver, getErr := a.db.GetDriver(id)
	if getErr != nil {
		return Response{Success: true, Data: driver} // still report success
	}
	return Response{Success: true, Data: driver}
}

func (a *App) DeleteDriver(id int) Response {
	err := a.db.DeleteDriver(id)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true}
}

func (a *App) SearchDrivers(query string) Response {
	items, err := a.db.SearchDrivers(query)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: items}
}

// ---------- Loader methods ----------

func (a *App) ListLoaders(filters map[string]interface{}, page int, pageSize int, sortBy string, sortDir string) Response {
	// TODO: wire through to database.ListLoaders(filters, database.QueryOptions{page, pageSize, sortBy, sortDir})
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	return Response{Success: true, Data: []database.Loader{}}
}

func (a *App) GetLoader(id int) Response {
	loader, err := a.db.GetLoader(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Response{Success: false, Error: "Loader not found"}
		}
		return Response{Success: false, Error: "Failed to fetch loader"}
	}
	return Response{Success: true, Data: loader}
}

func (a *App) CreateLoader(dto database.LoaderDTO) Response {
	loader, err := a.db.CreateLoader(dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: loader}
}

func (a *App) UpdateLoader(id int, dto database.LoaderDTO) Response {
	err := a.db.UpdateLoader(id, dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	loader, getErr := a.db.GetLoader(id)
	if getErr != nil {
		return Response{Success: true, Data: loader}
	}
	return Response{Success: true, Data: loader}
}

func (a *App) DeleteLoader(id int) Response {
	err := a.db.DeleteLoader(id)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true}
}

func (a *App) SearchLoaders(query string) Response {
	items, err := a.db.SearchLoaders(query)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: items}
}

// ---------- Vehicle methods ----------

func (a *App) ListVehicles(filters map[string]interface{}, page int, pageSize int, sortBy string, sortDir string) Response {
	// TODO: wire through to database.ListVehicles(filters, database.QueryOptions{page, pageSize, sortBy, sortDir})
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	return Response{Success: true, Data: []database.Vehicle{}}
}

func (a *App) GetVehicle(id int) Response {
	vehicle, err := a.db.GetVehicle(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Response{Success: false, Error: "Vehicle not found"}
		}
		return Response{Success: false, Error: "Failed to fetch vehicle"}
	}
	return Response{Success: true, Data: vehicle}
}

func (a *App) CreateVehicle(dto database.VehicleDTO) Response {
	vehicle, err := a.db.CreateVehicle(dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: vehicle}
}

func (a *App) UpdateVehicle(id int, dto database.VehicleDTO) Response {
	err := a.db.UpdateVehicle(id, dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	vehicle, getErr := a.db.GetVehicle(id)
	if getErr != nil {
		return Response{Success: true, Data: vehicle}
	}
	return Response{Success: true, Data: vehicle}
}

func (a *App) DeleteVehicle(id int) Response {
	err := a.db.DeleteVehicle(id)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true}
}

func (a *App) SearchVehicles(query string) Response {
	items, err := a.db.SearchVehicles(query)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: items}
}

// ---------- Well methods ----------

func (a *App) ListWells(filters map[string]interface{}, page int, pageSize int, sortBy string, sortDir string) Response {
	// TODO: wire through to database.ListWells(filters, database.QueryOptions{page, pageSize, sortBy, sortDir})
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	return Response{Success: true, Data: []database.Well{}}
}

func (a *App) GetWell(id int) Response {
	well, err := a.db.GetWell(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Response{Success: false, Error: "Well not found"}
		}
		return Response{Success: false, Error: "Failed to fetch well"}
	}
	return Response{Success: true, Data: well}
}

func (a *App) CreateWell(dto database.WellDTO) Response {
	well, err := a.db.CreateWell(dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: well}
}

func (a *App) UpdateWell(id int, dto database.WellDTO) Response {
	err := a.db.UpdateWell(id, dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	well, getErr := a.db.GetWell(id)
	if getErr != nil {
		return Response{Success: true, Data: well}
	}
	return Response{Success: true, Data: well}
}

func (a *App) DeleteWell(id int) Response {
	err := a.db.DeleteWell(id)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true}
}

func (a *App) SearchWells(query string) Response {
	items, err := a.db.SearchWells(query)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: items}
}

// ---------- Load methods ----------

func (a *App) ListLoads(filters map[string]interface{}, page int, pageSize int, sortBy string, sortDir string) Response {
	// TODO: wire through to database.ListLoads(filters, database.QueryOptions{page, pageSize, sortBy, sortDir})
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	_ = filters
	_ = page
	_ = pageSize
	_ = sortBy
	_ = sortDir
	return Response{Success: true, Data: []database.Load{}}
}

func (a *App) GetLoad(id int) Response {
	load, err := a.db.GetLoad(id)
	if err != nil {
		if err.Error() == "record not found" {
			return Response{Success: false, Error: "Load not found"}
		}
		return Response{Success: false, Error: "Failed to fetch load"}
	}
	return Response{Success: true, Data: load}
}

func (a *App) CreateLoad(dto database.LoadDTO) Response {
	load, err := a.db.CreateLoad(dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: load}
}

func (a *App) UpdateLoad(id int, dto database.LoadDTO) Response {
	err := a.db.UpdateLoad(id, dto)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	load, getErr := a.db.GetLoad(id)
	if getErr != nil {
		return Response{Success: true, Data: load}
	}
	return Response{Success: true, Data: load}
}

func (a *App) DeleteLoad(id int) Response {
	err := a.db.DeleteLoad(id)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true}
}

func (a *App) SearchLoads(query string) Response {
	items, err := a.db.SearchLoads(query)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: items}
}

func (a *App) GetLoadsByDriver(driverID int) Response {
	loads, err := a.db.GetLoadsByDriver(driverID)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: loads}
}

func (a *App) GetLoadsByWell(wellID int) Response {
	loads, err := a.db.GetLoadsByWell(wellID)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: loads}
}

func (a *App) GetLoadsByStatus(status bool) Response {
	loads, err := a.db.GetLoadsByStatus(status)
	if err != nil {
		return Response{Success: false, Error: err.Error()}
	}
	return Response{Success: true, Data: loads}
}

func main() {
	app := NewApp()
	if err := wails.Run(&options.App{
		Title:    "Cargo Management App",
		Width:    1024,
		Height:   768,
		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup: app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
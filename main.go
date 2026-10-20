package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"

	"cargo_management_app/go/app"
	"cargo_management_app/go/database"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// dbPath returns the per-user location of the SQLite database.
func dbPath() string {
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

func main() {
	wailsApp := &app.App{}

	err := wails.Run(&options.App{
		Title:       "Cargo Management App",
		Width:       1024,
		Height:      768,
		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup: func(ctx context.Context) {
			db, err := database.NewDatabase(dbPath())
			if err != nil {
				log.Fatalf("initialize database: %v", err)
			}
			wailsApp.SetDatabase(db)
			wailsApp.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			wailsApp.Shutdown(ctx)
		},
		Bind: []interface{}{
			wailsApp,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
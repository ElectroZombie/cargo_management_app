// Package app holds the Wails application struct and its exported bindings.
//
// Bindings are thin wrappers: they translate frontend input into service calls,
// wrap the result in a response envelope, and never contain business logic.
package app

import (
	"context"
	"log/slog"

	"cargo_management_app/go/database"
)

// Response is the envelope every binding returns so the frontend never has to
// interpret raw Go errors.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Total   int64       `json:"total,omitempty"`
	// Code is a stable machine-readable classification of Error, allowing the
	// frontend to react to a failure category without parsing message text.
	Code string `json:"code,omitempty"`
}

// App is the Wails application struct. It owns the database handle and exposes
// the operational services to the frontend.
type App struct {
	ctx context.Context
	db  *database.Database
}

// NewApp creates the application bound to the given database handle.
func NewApp(db *database.Database) *App {
	return &App{db: db}
}

// SetDatabase attaches the database handle used by every binding.
func (a *App) SetDatabase(db *database.Database) {
	a.db = db
}

// Startup is called when the Wails app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	slog.Info("application started")
}

// Shutdown is called when the Wails app closes.
func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
	slog.Info("application shutdown")
}

// failure builds an unsuccessful response from a domain error. The error is
// classified so the frontend receives both a readable message and a code.
func failure(err error) Response {
	return Response{Success: false, Error: formatError(err), Code: errorCode(err)}
}

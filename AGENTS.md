# 🤖 AGENTS.md - Development Best Practices

**Purpose:** Guidelines for coding agents and developers working on the Cargo Management App  
**Last Updated:** 2026-10-02  
**Branch:** main

---

## 📋 Table of Contents

1. [General Principles](#general-principles)
2. [Go Backend Best Practices](#go-backend-best-practices)
3. [Wails Framework Integration](#wails-framework-integration)
4. [Angular Frontend Best Practices](#angular-frontend-best-practices)
5. [Code Organization](#code-organization)
6. [Naming Conventions](#naming-conventions)
7. [Testing Strategy](#testing-strategy)
8. [Error Handling](#error-handling)
9. [Documentation Standards](#documentation-standards)
10. [CI/CD and Git Workflow](#cicd-and-git-workflow)

---

## General Principles

### 1. **Clean Architecture**
- Maintain strict separation of concerns: Database → Service → API → Frontend
- Each layer should be independently testable
- Avoid circular dependencies and tightly coupled code
- Use dependency injection where applicable

### 2. **Consistency Over Cleverness**
- Write code that is easy to read and maintain, not code that impresses
- Follow established patterns in the codebase
- Use domain-driven design terminology consistently
- Prioritize clarity for future developers (including future you)

### 3. **No Authentication or User Roles**
- **Never** add authentication, login, or role-based access control
- Build the application as a single-tenant operational system
- Assume all users have full access to all operations
- Document this assumption in any architectural decisions

### 4. **Data Integrity First**
- Prioritize database constraints and validation over application logic
- Use foreign keys with proper cascade rules (restrict/cascade as appropriate)
- Enforce data uniqueness at the database level
- Validate domain rules in services before database operations

### 5. **Workspace Cleanliness**
- Keep the project structure organized and predictable
- Remove unused files, imports, and dependencies regularly
- Follow the established directory structure strictly
- Use meaningful commit messages that describe *why*, not just *what*

---

## Go Backend Best Practices

### 1. **Package Structure**

```
cmd/
  ├── main.go                    # Application entry point
go/
  ├── database/
  │   ├── database.go            # DB initialization and schema
  │   ├── migrations.go           # Migration logic (future)
  │   └── seed.go                 # Seed data (future)
  ├── models/
  │   ├── driver.go               # Domain entities
  │   ├── loader.go
  │   ├── vehicle.go
  │   ├── well.go
  │   └── load.go
  ├── services/
  │   ├── driver_service.go       # Business logic (CRUD + queries)
  │   ├── loader_service.go
  │   ├── vehicle_service.go
  │   ├── well_service.go
  │   └── load_service.go
  ├── dtos/
  │   ├── driver_dto.go           # Data transfer objects
  │   ├── loader_dto.go
  │   ├── vehicle_dto.go
  │   ├── well_dto.go
  │   └── load_dto.go
  ├── errors/
  │   └── errors.go               # Custom domain error types
  ├── validators/
  │   ├── driver_validator.go     # DTO validation logic
  │   ├── loader_validator.go
  │   └── shared_validator.go     # Common validation rules
  ├── app/
  │   └── app.go                  # Wails application struct
  └── utils/
      └── constants.go             # Shared constants
```

### 2. **Error Handling**

**Never expose raw database errors to the frontend.** Create domain-specific error types:

```go
// go/errors/errors.go
package errors

type ValidationError struct {
	Field   string
	Message string
}

type NotFoundError struct {
	Entity string
	ID     interface{}
}

type ConflictError struct {
	Field   string
	Message string
}

type InternalError struct {
	Message string
}

// Implement error interface for each type
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
```

**Always wrap SQLite errors:**

```go
if err != nil {
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return ConflictError{Field: "email", Message: "already exists"}
	}
	if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		return ConflictError{Field: "vehicle_id", Message: "invalid reference"}
	}
	return InternalError{Message: "database operation failed"}
}
```

### 3. **Service Layer Pattern**

All CRUD operations should follow this pattern:

```go
// go/services/driver_service.go
package services

type DriverService struct {
	db *sql.DB
}

// Create validates and inserts a new driver
func (s *DriverService) Create(ctx context.Context, dto *dtos.DriverDTO) (*models.Driver, error) {
	// 1. Validate DTO
	if err := validators.ValidateDriver(dto); err != nil {
		return nil, err
	}

	// 2. Check business rules (e.g., unique license number)
	exists, err := s.driverExists(ctx, dto.LicenseNumber)
	if exists {
		return nil, ConflictError{Field: "license_number", Message: "already exists"}
	}

	// 3. Execute database operation
	result, err := s.db.ExecContext(ctx, insertDriverSQL, /* args */)
	if err != nil {
		return nil, handleDBError(err)
	}

	// 4. Fetch and return created entity
	id, _ := result.LastInsertId()
	return s.GetByID(ctx, id)
}

// GetByID retrieves a driver by ID
func (s *DriverService) GetByID(ctx context.Context, id int64) (*models.Driver, error) {
	driver := &models.Driver{}
	err := s.db.QueryRowContext(ctx, selectDriverByIDSQL, id).Scan(/* fields */)
	if err == sql.ErrNoRows {
		return nil, NotFoundError{Entity: "driver", ID: id}
	}
	if err != nil {
		return nil, InternalError{Message: "failed to fetch driver"}
	}
	return driver, nil
}

// List returns all drivers with optional filtering
func (s *DriverService) List(ctx context.Context) ([]models.Driver, error) {
	rows, err := s.db.QueryContext(ctx, selectAllDriversSQL)
	if err != nil {
		return nil, InternalError{Message: "failed to fetch drivers"}
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		driver := models.Driver{}
		if err := rows.Scan(/* fields */); err != nil {
			return nil, InternalError{Message: "failed to scan driver"}
		}
		drivers = append(drivers, driver)
	}
	return drivers, rows.Err()
}

// Update modifies an existing driver
func (s *DriverService) Update(ctx context.Context, id int64, dto *dtos.DriverDTO) (*models.Driver, error) {
	// 1. Verify driver exists
	if _, err := s.GetByID(ctx, id); err != nil {
		return nil, err
	}

	// 2. Validate DTO
	if err := validators.ValidateDriver(dto); err != nil {
		return nil, err
	}

	// 3. Execute update
	_, err := s.db.ExecContext(ctx, updateDriverSQL, /* args with id */)
	if err != nil {
		return nil, handleDBError(err)
	}

	// 4. Return updated entity
	return s.GetByID(ctx, id)
}

// Delete removes a driver
func (s *DriverService) Delete(ctx context.Context, id int64) error {
	// Verify driver exists before deletion
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}

	result, err := s.db.ExecContext(ctx, deleteDriverSQL, id)
	if err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return ConflictError{
				Field:   "driver_id",
				Message: "cannot delete driver with assigned loads",
			}
		}
		return InternalError{Message: "failed to delete driver"}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return NotFoundError{Entity: "driver", ID: id}
	}

	return nil
}
```

### 4. **Context Usage**
- Always accept `context.Context` as the first parameter in service methods
- Use `QueryContext`, `ExecContext` for database operations
- Respect context cancellation in loops and long-running operations
- Pass context through the call chain (service → database)

### 5. **SQL Query Management**

Keep SQL queries as named constants in the same file:

```go
const (
	insertDriverSQL = `
		INSERT INTO driver (first_name, last_name, license_number, phone, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`
	selectDriverByIDSQL = `SELECT id, first_name, last_name, license_number, phone, email, created_at, updated_at FROM driver WHERE id = ?`
	selectAllDriversSQL = `SELECT id, first_name, last_name, license_number, phone, email, created_at, updated_at FROM driver ORDER BY created_at DESC`
	updateDriverSQL = `UPDATE driver SET first_name = ?, last_name = ?, phone = ?, email = ?, updated_at = datetime('now') WHERE id = ?`
	deleteDriverSQL = `DELETE FROM driver WHERE id = ?`
)
```

### 6. **Logging**
- Use a structured logging library (e.g., `log/slog` in Go 1.21+)
- Log errors at appropriate levels (DEBUG, INFO, WARN, ERROR)
- Never log sensitive data (passwords, tokens, full credit cards)
- Include context information (operation, entity ID, duration)

```go
import "log/slog"

slog.Error("failed to fetch driver", "driver_id", id, "error", err)
slog.Info("driver created", "driver_id", newDriver.ID)
```

### 7. **Concurrency**
- Use buffered channels with explicit buffer sizes
- Always close channels from the sender side
- Handle goroutine panics with recovery
- Prefer sync.WaitGroup for coordinating goroutines
- Document race conditions and concurrency assumptions in comments

---

## Wails Framework Integration

### 1. **Application Struct Pattern**

The Wails app struct should be a thin orchestration layer:

```go
// go/app/app.go
package app

import (
	"context"
	"database/sql"
	"myapp/go/services"
	"myapp/go/dtos"
	"myapp/go/models"
)

type App struct {
	db               *sql.DB
	driverService    *services.DriverService
	loaderService    *services.LoaderService
	vehicleService   *services.VehicleService
	wellService      *services.WellService
	loadService      *services.LoadService
}

// NewApp creates and initializes the application
func NewApp(db *sql.DB) *App {
	return &App{
		db:            db,
		driverService: services.NewDriverService(db),
		// ... other services
	}
}

// Startup is called when Wails app starts
func (a *App) Startup(ctx context.Context) {
	// Perform any warm-up operations
	slog.Info("application started")
}

// Shutdown is called when Wails app closes
func (a *App) Shutdown(ctx context.Context) {
	a.db.Close()
	slog.Info("application shutdown")
}
```

### 2. **Wails Method Binding**

Wails methods should be **thin wrappers** that:
- Accept TypeScript-compatible types (strings, numbers, structs)
- Return typed responses (DTOs, not raw database models)
- Handle context passing to services
- Convert domain errors to user-friendly response objects

```go
// go/app/driver_bindings.go
package app

import (
	"context"
	"myapp/go/dtos"
	"myapp/go/models"
)

// Response wrapper for consistent frontend consumption
type DriverResponse struct {
	Success bool          `json:"success"`
	Data    *models.Driver `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type DriverListResponse struct {
	Success bool            `json:"success"`
	Data    []models.Driver `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// CreateDriver creates a new driver
func (a *App) CreateDriver(ctx context.Context, dto dtos.DriverDTO) DriverResponse {
	driver, err := a.driverService.Create(ctx, &dto)
	if err != nil {
		return DriverResponse{
			Success: false,
			Error:   formatErrorMessage(err),
		}
	}
	return DriverResponse{
		Success: true,
		Data:    driver,
	}
}

// GetDriver retrieves a driver by ID
func (a *App) GetDriver(ctx context.Context, id int64) DriverResponse {
	driver, err := a.driverService.GetByID(ctx, id)
	if err != nil {
		return DriverResponse{
			Success: false,
			Error:   formatErrorMessage(err),
		}
	}
	return DriverResponse{
		Success: true,
		Data:    driver,
	}
}

// ListDrivers retrieves all drivers
func (a *App) ListDrivers(ctx context.Context) DriverListResponse {
	drivers, err := a.driverService.List(ctx)
	if err != nil {
		return DriverListResponse{
			Success: false,
			Error:   formatErrorMessage(err),
		}
	}
	return DriverListResponse{
		Success: true,
		Data:    drivers,
	}
}

// UpdateDriver updates an existing driver
func (a *App) UpdateDriver(ctx context.Context, id int64, dto dtos.DriverDTO) DriverResponse {
	driver, err := a.driverService.Update(ctx, id, &dto)
	if err != nil {
		return DriverResponse{
			Success: false,
			Error:   formatErrorMessage(err),
		}
	}
	return DriverResponse{
		Success: true,
		Data:    driver,
	}
}

// DeleteDriver removes a driver
func (a *App) DeleteDriver(ctx context.Context, id int64) DriverResponse {
	err := a.driverService.Delete(ctx, id)
	if err != nil {
		return DriverResponse{
			Success: false,
			Error:   formatErrorMessage(err),
		}
	}
	return DriverResponse{
		Success: true,
	}
}

// formatErrorMessage converts domain errors to user-friendly messages
func formatErrorMessage(err error) string {
	switch e := err.(type) {
	case errors.ValidationError:
		return fmt.Sprintf("Validation error on %s: %s", e.Field, e.Message)
	case errors.NotFoundError:
		return fmt.Sprintf("%s not found", e.Entity)
	case errors.ConflictError:
		return fmt.Sprintf("Conflict: %s", e.Message)
	case errors.InternalError:
		return "An internal error occurred. Please try again."
	default:
		return "An unexpected error occurred"
	}
}
```

### 3. **IPC Protocol**
- All Wails methods should be **synchronous** (no goroutines in bindings)
- Return structured response objects with `success` and `error` fields
- Use lowercase field names in JSON responses for JavaScript compatibility
- Always include timestamps in audit operations

### 4. **Context in Wails**
- Extract context from Wails callback (first parameter)
- Create a reasonable timeout if needed (e.g., 30 seconds)
- Pass context through all service calls

---

## Angular Frontend Best Practices

### 1. **Project Structure**

```
src/
  ├── app/
  │   ├── app.component.ts
  │   ├── app.component.html
  │   ├── app.component.scss
  │   ├── app-routing.module.ts
  │   ├── core/
  │   │   ├── services/
  │   │   │   ├── driver.service.ts
  │   │   │   ├── loader.service.ts
  │   │   │   ├── vehicle.service.ts
  │   │   │   ├── well.service.ts
  │   │   │   ├── load.service.ts
  │   │   │   ├── wails.service.ts              # Wails bridge
  │   │   │   └── notification.service.ts
  │   │   ├── models/
  │   │   │   ├── driver.model.ts
  │   │   │   ├── loader.model.ts
  │   │   │   ├── vehicle.model.ts
  │   │   │   ├── well.model.ts
  │   │   │   ├── load.model.ts
  │   │   │   └── api-response.model.ts
  │   │   └── core.module.ts
  │   ├── shared/
  │   │   ├── components/
  │   │   │   ├── button/
  │   │   │   ├── card/
  │   │   │   ├── modal/
  │   │   │   ├── loading-spinner/
  │   │   │   ├── toast/
  │   │   │   └── confirm-dialog/
  │   │   ├── pipes/
  │   │   ├── directives/
  │   │   ├── styles/
  │   │   │   ├── _variables.scss
  │   │   │   ├── _typography.scss
  │   │   │   ├── _mixins.scss
  │   │   │   └── _global.scss
  │   │   └── shared.module.ts
  │   ├── features/
  │   │   ├── dashboard/
  │   │   │   ├── dashboard.component.ts
  │   │   │   ├── dashboard-routing.module.ts
  │   │   │   └── dashboard.module.ts
  │   │   ├── drivers/
  │   │   │   ├── driver-list/
  │   │   │   ├── driver-form/
  │   │   │   ├── driver-detail/
  │   │   │   ├── drivers-routing.module.ts
  │   │   │   └── drivers.module.ts
  │   │   ├── loaders/
  │   │   ├── vehicles/
  │   │   ├── wells/
  │   │   └── loads/
  │   ├── layout/
  │   │   ├── navbar/
  │   │   ├── sidebar/
  │   │   ├── footer/
  │   │   ├── breadcrumb/
  │   │   └── layout.component.ts
  │   └── app.module.ts
  ├── assets/
  ├── styles/
  │   ├── index.scss
  │   └── theme.scss
  ├── environments/
  ├── main.ts
  └── index.html
```

### 2. **Service Layer**

Every Wails binding should have a corresponding Angular service:

```typescript
// src/app/core/services/driver.service.ts
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { Driver } from '../models/driver.model';
import { DriverDTO } from '../models/driver.model';
import { WailsService } from './wails.service';
import { NotificationService } from './notification.service';

@Injectable({ providedIn: 'root' })
export class DriverService {
  private drivers$ = new BehaviorSubject<Driver[]>([]);
  private loading$ = new BehaviorSubject<boolean>(false);
  private selectedDriver$ = new BehaviorSubject<Driver | null>(null);

  constructor(
    private wails: WailsService,
    private notification: NotificationService
  ) {}

  // Observable streams
  getDrivers(): Observable<Driver[]> {
    return this.drivers$.asObservable();
  }

  isLoading(): Observable<boolean> {
    return this.loading$.asObservable();
  }

  getSelectedDriver(): Observable<Driver | null> {
    return this.selectedDriver$.asObservable();
  }

  // CRUD operations
  async createDriver(dto: DriverDTO): Promise<Driver | null> {
    this.loading$.next(true);
    try {
      const response = await this.wails.invoke('CreateDriver', dto);
      if (response.success) {
        this.notification.success('Driver created successfully');
        await this.loadDrivers();
        return response.data;
      } else {
        this.notification.error(response.error);
        return null;
      }
    } catch (error) {
      this.notification.error('Failed to create driver');
      console.error(error);
      return null;
    } finally {
      this.loading$.next(false);
    }
  }

  async loadDrivers(): Promise<void> {
    this.loading$.next(true);
    try {
      const response = await this.wails.invoke('ListDrivers');
      if (response.success) {
        this.drivers$.next(response.data);
      } else {
        this.notification.error(response.error);
      }
    } catch (error) {
      this.notification.error('Failed to load drivers');
      console.error(error);
    } finally {
      this.loading$.next(false);
    }
  }

  async getDriver(id: number): Promise<Driver | null> {
    this.loading$.next(true);
    try {
      const response = await this.wails.invoke('GetDriver', id);
      if (response.success) {
        this.selectedDriver$.next(response.data);
        return response.data;
      } else {
        this.notification.error(response.error);
        return null;
      }
    } catch (error) {
      this.notification.error('Failed to load driver');
      console.error(error);
      return null;
    } finally {
      this.loading$.next(false);
    }
  }

  async updateDriver(id: number, dto: DriverDTO): Promise<Driver | null> {
    this.loading$.next(true);
    try {
      const response = await this.wails.invoke('UpdateDriver', id, dto);
      if (response.success) {
        this.notification.success('Driver updated successfully');
        await this.loadDrivers();
        return response.data;
      } else {
        this.notification.error(response.error);
        return null;
      }
    } catch (error) {
      this.notification.error('Failed to update driver');
      console.error(error);
      return null;
    } finally {
      this.loading$.next(false);
    }
  }

  async deleteDriver(id: number): Promise<boolean> {
    this.loading$.next(true);
    try {
      const response = await this.wails.invoke('DeleteDriver', id);
      if (response.success) {
        this.notification.success('Driver deleted successfully');
        await this.loadDrivers();
        return true;
      } else {
        this.notification.error(response.error);
        return false;
      }
    } catch (error) {
      this.notification.error('Failed to delete driver');
      console.error(error);
      return false;
    } finally {
      this.loading$.next(false);
    }
  }
}
```

### 3. **Wails Bridge Service**

Create a single service to handle all Wails communication:

```typescript
// src/app/core/services/wails.service.ts
import { Injectable } from '@angular/core';

declare global {
  interface Window {
    runtime: any;
  }
}

@Injectable({ providedIn: 'root' })
export class WailsService {
  async invoke<T = any>(method: string, ...args: any[]): Promise<T> {
    return new Promise((resolve, reject) => {
      if (!window.runtime || !window.runtime.Call) {
        reject(new Error('Wails runtime not available'));
        return;
      }

      window.runtime.Call.ByName(method, ...args)
        .then((result: T) => resolve(result))
        .catch((error: any) => reject(error));
    });
  }
}
```

### 4. **Models and DTOs**

Keep models synchronized with Go DTOs:

```typescript
// src/app/core/models/driver.model.ts
export interface Driver {
  id: number;
  firstName: string;
  lastName: string;
  licenseNumber: string;
  phone: string;
  email: string;
  createdAt: string;  // ISO 8601 timestamp
  updatedAt: string;
}

export interface DriverDTO {
  firstName: string;
  lastName: string;
  licenseNumber: string;
  phone: string;
  email: string;
}

// API Response wrapper
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface ListResponse<T> {
  success: boolean;
  data?: T[];
  error?: string;
}
```

### 5. **Component Architecture**

**Smart (Container) Components:**
- Connect to services and manage state
- Pass data to presentational components
- Handle user actions and trigger service calls
- Subscribe to observables

```typescript
// src/app/features/drivers/driver-list/driver-list.component.ts
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { Driver } from '../../../core/models/driver.model';
import { DriverService } from '../../../core/services/driver.service';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-driver-list',
  templateUrl: './driver-list.component.html',
  styleUrls: ['./driver-list.component.scss'],
})
export class DriverListComponent implements OnInit, OnDestroy {
  drivers$ = this.driverService.getDrivers();
  loading$ = this.driverService.isLoading();
  private destroy$ = new Subject<void>();

  constructor(
    private driverService: DriverService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.driverService.loadDrivers();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  onEdit(driver: Driver): void {
    this.router.navigate(['/drivers', driver.id, 'edit']);
  }

  async onDelete(driver: Driver): Promise<void> {
    const confirmed = confirm(`Delete driver ${driver.firstName} ${driver.lastName}?`);
    if (confirmed) {
      await this.driverService.deleteDriver(driver.id);
    }
  }
}
```

**Presentational (Dumb) Components:**
- Accept data via `@Input()`
- Emit events via `@Output()`
- No service dependencies
- Fully reusable

```typescript
// src/app/shared/components/button/button.component.ts
import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-button',
  template: `
    <button
      [class]="'btn btn-' + variant"
      [disabled]="disabled || loading"
      (click)="onClick()"
    >
      <span *ngIf="!loading">{{ text }}</span>
      <app-loading-spinner *ngIf="loading" size="small"></app-loading-spinner>
    </button>
  `,
  styleUrls: ['./button.component.scss'],
})
export class ButtonComponent {
  @Input() text: string = 'Button';
  @Input() variant: 'primary' | 'secondary' | 'danger' = 'primary';
  @Input() disabled: boolean = false;
  @Input() loading: boolean = false;
  @Output() clicked = new EventEmitter<void>();

  onClick(): void {
    this.clicked.emit();
  }
}
```

### 6. **Reactive Forms**

Always use Reactive Forms with strong typing:

```typescript
// src/app/features/drivers/driver-form/driver-form.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { DriverService } from '../../../core/services/driver.service';
import { Driver, DriverDTO } from '../../../core/models/driver.model';

@Component({
  selector: 'app-driver-form',
  templateUrl: './driver-form.component.html',
  styleUrls: ['./driver-form.component.scss'],
})
export class DriverFormComponent implements OnInit {
  form: FormGroup;
  isEditing = false;
  driverId: number | null = null;
  loading$ = this.driverService.isLoading();

  constructor(
    private fb: FormBuilder,
    private driverService: DriverService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.form = this.createForm();
  }

  ngOnInit(): void {
    this.driverId = this.route.snapshot.paramMap.get('id')
      ? Number(this.route.snapshot.paramMap.get('id'))
      : null;

    if (this.driverId) {
      this.isEditing = true;
      this.loadDriver(this.driverId);
    }
  }

  createForm(): FormGroup {
    return this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      licenseNumber: ['', [Validators.required, Validators.pattern(/^[A-Z0-9]+$/)]],
      phone: ['', [Validators.required, Validators.pattern(/^\d{10,}$/)]],
      email: ['', [Validators.required, Validators.email]],
    });
  }

  async loadDriver(id: number): Promise<void> {
    const driver = await this.driverService.getDriver(id);
    if (driver) {
      this.form.patchValue({
        firstName: driver.firstName,
        lastName: driver.lastName,
        licenseNumber: driver.licenseNumber,
        phone: driver.phone,
        email: driver.email,
      });
    }
  }

  async onSubmit(): Promise<void> {
    if (!this.form.valid) {
      return;
    }

    const dto: DriverDTO = this.form.value;

    if (this.isEditing && this.driverId) {
      const success = await this.driverService.updateDriver(this.driverId, dto);
      if (success) {
        this.router.navigate(['/drivers']);
      }
    } else {
      const driver = await this.driverService.createDriver(dto);
      if (driver) {
        this.router.navigate(['/drivers']);
      }
    }
  }
}
```

### 7. **Change Detection**
- Use `OnPush` change detection strategy in all components
- Leverage observables and async pipe
- Minimize manual change detection calls

```typescript
import { ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-example',
  templateUrl: './example.component.html',
  styleUrls: ['./example.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ExampleComponent {}
```

### 8. **Styling Guidelines**

Create a centralized theme:

```scss
// src/app/shared/styles/_variables.scss
// Colors
$primary-color: #5A9FD4;
$secondary-color: #6BBD9F;
$accent-color: #F5A28B;
$background-color: #F7F8FA;
$neutral-light: #E8E9EB;
$text-primary: #2D3436;
$text-secondary: #636E72;
$success-color: #27AE60;
$warning-color: #F39C12;
$danger-color: #E74C3C;

// Spacing
$spacing-xs: 4px;
$spacing-sm: 8px;
$spacing-md: 16px;
$spacing-lg: 24px;
$spacing-xl: 32px;

// Typography
$font-primary: 'Roboto', sans-serif;
$font-size-body: 14px;
$font-size-body-lg: 16px;
$font-size-heading: 24px;
$line-height-normal: 1.5;

// Border radius
$border-radius-sm: 4px;
$border-radius-md: 8px;
$border-radius-lg: 16px;

// Shadows
$shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
$shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
$shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
```

Use variables in components:

```scss
// src/app/features/drivers/driver-list/driver-list.component.scss
@import 'app/shared/styles/variables';

.driver-list {
  padding: $spacing-lg;
  background-color: $background-color;

  &__title {
    font-size: $font-size-heading;
    color: $text-primary;
    margin-bottom: $spacing-md;
    font-weight: bold;
  }

  &__table {
    background: white;
    border-radius: $border-radius-md;
    box-shadow: $shadow-md;
  }
}
```

---

## Code Organization

### 1. **File Naming Conventions**

| Type | Pattern | Example |
|------|---------|---------|
| Components | `kebab-case.component.ts` | `driver-list.component.ts` |
| Services | `kebab-case.service.ts` | `driver.service.ts` |
| Models | `kebab-case.model.ts` | `driver.model.ts` |
| Modules | `kebab-case.module.ts` | `drivers.module.ts` |
| Go packages | `snake_case` | `driver_service.go` |
| Go files | `snake_case.go` | `driver_service.go` |

### 2. **Directory Organization**

**By feature, not by type:**
- ❌ BAD: `src/app/components/drivers`, `src/app/services/drivers`, `src/app/models/drivers`
- ✅ GOOD: `src/app/features/drivers/list`, `src/app/features/drivers/form`, `src/app/core/models/driver.model.ts`

### 3. **Import Organization**

In all files, organize imports in this order:

```typescript
// 1. Angular core
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';

// 2. RxJS
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

// 3. Third-party libraries
import { moment } from 'moment';

// 4. Local app services and models
import { DriverService } from '../../../core/services/driver.service';
import { Driver } from '../../../core/models/driver.model';

// 5. Local component/module files
import { DriverListComponent } from './driver-list.component';
```

---

## Naming Conventions

### Go Naming

| Entity | Convention | Example |
|--------|-----------|---------|
| Package | `snake_case` | `driver_service`, `user_service` |
| Function/Method | `PascalCase` | `CreateDriver()`, `GetByID()` |
| Variable | `camelCase` | `driverID`, `vehicleCount` |
| Constant | `UPPER_SNAKE_CASE` | `MAX_DRIVERS`, `INSERT_DRIVER_SQL` |
| Interface | `PascalCase` | `DriverRepository`, `Validator` |
| Error type | `*Error` suffix | `ValidationError`, `NotFoundError` |
| Private field | `camelCase` | `db`, `driverService` |

### TypeScript Naming

| Entity | Convention | Example |
|--------|-----------|---------|
| Class | `PascalCase` | `DriverService`, `DriverComponent` |
| Interface | `PascalCase` | `Driver`, `IDriver` (no `I` prefix) |
| Method | `camelCase` | `getDriver()`, `createDriver()` |
| Property | `camelCase` | `driverId`, `firstName` |
| Constant | `UPPER_SNAKE_CASE` | `MAX_DRIVERS`, `DEFAULT_TIMEOUT` |
| Enum | `PascalCase` | `LoadStatus`, `VehicleType` |
| Observable | `$` suffix | `drivers$`, `loading$` |
| Private field | `#` or `private` | `#cache`, `private db` |
| DTO | `DTO` suffix | `DriverDTO`, `VehicleDTO` |

---

## Testing Strategy

### Go Testing

**Test file location:** `{source_file}_test.go`

```go
// go/services/driver_service_test.go
package services

import (
	"context"
	"testing"
	"myapp/go/dtos"
	"myapp/go/errors"
)

func TestCreateDriver_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer db.Close()
	service := NewDriverService(db)

	// Execute
	dto := &dtos.DriverDTO{
		FirstName:     "John",
		LastName:      "Doe",
		LicenseNumber: "DL123456",
		Phone:         "5551234567",
		Email:         "john@example.com",
	}
	driver, err := service.Create(context.Background(), dto)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if driver.ID == 0 {
		t.Error("driver should have an ID")
	}
	if driver.FirstName != "John" {
		t.Errorf("expected John, got %s", driver.FirstName)
	}
}

func TestCreateDriver_ValidationError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	service := NewDriverService(db)

	dto := &dtos.DriverDTO{
		FirstName: "", // Invalid: empty name
	}
	_, err := service.Create(context.Background(), dto)

	if err == nil {
		t.Error("expected validation error, got nil")
	}
	if _, ok := err.(errors.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}
```

**Testing guidelines:**
- Test happy path and error cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Clean up resources (close databases, etc.)
- Use descriptive test names: `Test{Function}_{Scenario}_{Expected}`

### Angular Testing

Use `TestBed` for component and service testing:

```typescript
// src/app/core/services/driver.service.spec.ts
import { TestBed } from '@angular/core/testing';
import { DriverService } from './driver.service';
import { WailsService } from './wails.service';
import { NotificationService } from './notification.service';

describe('DriverService', () => {
  let service: DriverService;
  let wailsMock: jasmine.SpyObj<WailsService>;
  let notificationMock: jasmine.SpyObj<NotificationService>;

  beforeEach(() => {
    wailsMock = jasmine.createSpyObj('WailsService', ['invoke']);
    notificationMock = jasmine.createSpyObj('NotificationService', [
      'success',
      'error',
    ]);

    TestBed.configureTestingModule({
      providers: [
        DriverService,
        { provide: WailsService, useValue: wailsMock },
        { provide: NotificationService, useValue: notificationMock },
      ],
    });

    service = TestBed.inject(DriverService);
  });

  it('should load drivers', async () => {
    const mockDrivers = [
      { id: 1, firstName: 'John', lastName: 'Doe' },
    ];
    wailsMock.invoke.and.returnValue(
      Promise.resolve({ success: true, data: mockDrivers })
    );

    await service.loadDrivers();

    service.getDrivers().subscribe((drivers) => {
      expect(drivers).toEqual(mockDrivers);
    });
  });
});
```

---

## Error Handling

### Go Error Handling

**Always define domain-specific errors:**

```go
// go/errors/errors.go
package errors

import "fmt"

type DomainError interface {
	error
	IsValidation() bool
	IsNotFound() bool
	IsConflict() bool
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func (e ValidationError) IsValidation() bool { return true }
func (e ValidationError) IsNotFound() bool   { return false }
func (e ValidationError) IsConflict() bool   { return false }
```

**Handle errors at the appropriate layer:**

```go
// Service layer: wrap and transform errors
func (s *DriverService) Create(ctx context.Context, dto *dtos.DriverDTO) (*models.Driver, error) {
	if err := validators.ValidateDriver(dto); err != nil {
		return nil, err  // Already a domain error
	}

	// Database operation
	result, err := s.db.ExecContext(ctx, insertDriverSQL, /* args */)
	if err != nil {
		// Transform SQLite error to domain error
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, ValidationError{Field: "license_number", Message: "already exists"}
		}
		return nil, InternalError{Message: "database error"}
	}

	return s.GetByID(ctx, id)
}
```

### Angular Error Handling

**Create a centralized error handler:**

```typescript
// src/app/core/services/error-handler.service.ts
import { Injectable, ErrorHandler, Injector } from '@angular/core';
import { NotificationService } from './notification.service';

@Injectable()
export class AppErrorHandler implements ErrorHandler {
  constructor(private injector: Injector) {}

  handleError(error: Error | any): void {
    const notificationService = this.injector.get(NotificationService);

    let message = 'An unexpected error occurred';

    if (error instanceof Error) {
      message = error.message;
    } else if (typeof error === 'string') {
      message = error;
    }

    console.error('Error:', error);
    notificationService.error(message);
  }
}

// In app.module.ts
import { ErrorHandler } from '@angular/core';
import { AppErrorHandler } from './core/services/error-handler.service';

@NgModule({
  // ...
  providers: [
    { provide: ErrorHandler, useClass: AppErrorHandler },
  ],
})
export class AppModule {}
```

---

## Documentation Standards

### Go Documentation

**Document public functions with comments:**

```go
// CreateDriver inserts a new driver into the database.
// It validates the DTO, checks for duplicate license numbers,
// and returns the created driver or a domain error.
//
// Errors:
// - ValidationError: if the DTO is invalid
// - ConflictError: if the license number already exists
// - InternalError: if the database operation fails
func (s *DriverService) Create(ctx context.Context, dto *dtos.DriverDTO) (*models.Driver, error) {
	// ...
}
```

**Use clear, domain-specific language:**

```go
// Bad comment
// do the thing with the value

// Good comment
// CalculateLoadWeight computes the total weight of all items in a load,
// excluding any items marked as damaged or recalled.
func (s *LoadService) CalculateLoadWeight(ctx context.Context, loadID int64) (float64, error) {
	// ...
}
```

### TypeScript Documentation

**Use JSDoc for public APIs:**

```typescript
/**
 * Loads all drivers from the backend and updates the drivers$ observable.
 *
 * @returns A promise that resolves when drivers are loaded
 *
 * @example
 * await this.driverService.loadDrivers();
 * this.driverService.getDrivers().subscribe(drivers => {
 *   console.log(drivers);
 * });
 */
async loadDrivers(): Promise<void> {
  // ...
}
```

**Document component inputs and outputs:**

```typescript
@Component({
  selector: 'app-driver-list',
  templateUrl: './driver-list.component.html',
})
export class DriverListComponent {
  /** Array of drivers to display in the table */
  @Input() drivers: Driver[] = [];

  /** Emitted when the user clicks the edit button for a driver */
  @Output() edit = new EventEmitter<Driver>();

  /** Emitted when the user clicks the delete button for a driver */
  @Output() delete = new EventEmitter<number>();
}
```

### README and Markdown

**Keep documentation current with code changes:**

- Update documentation in the same PR as code changes
- Use clear, short sentences
- Provide examples for complex features
- Link to relevant documentation

---

## CI/CD and Git Workflow

### Git Commit Messages

**Format:** `<type>(<scope>): <subject>`

```
feat(driver-service): add license number validation
^    ^               ^
|    |               └─ Imperative, present tense
|    └─ Scope of the commit
└─ Type of commit
```

**Types:**

- `feat:` A new feature
- `fix:` A bug fix
- `refactor:` Code change that neither fixes a bug nor adds a feature
- `test:` Adding or updating tests
- `docs:` Documentation-only changes
- `chore:` Changes to build process, dependencies, etc.
- `style:` Code style changes (formatting, semicolons, etc.)

**Examples:**

```
feat(api): expose driver CRUD operations via Wails
fix(driver-service): handle duplicate license numbers correctly
docs(AGENTS.md): add Go testing guidelines
test(driver-service): add unit tests for validation
refactor(error-handling): create domain-specific error types
```

### Branch Naming

`<type>/<description>`

```
feature/driver-crud
bugfix/license-validation
docs/update-agents
hotfix/database-error
```

### Pull Request Workflow

1. **Create a feature branch** from `main`
2. **Implement changes** following these guidelines
3. **Test thoroughly** with unit and integration tests
4. **Create a PR** with a clear description
5. **Request review** from team members
6. **Address feedback** in subsequent commits
7. **Merge** when approved (use "Squash and merge" for clean history)

### Code Review Checklist

Before approving a PR, verify:

- [ ] Code follows these guidelines
- [ ] Tests are included and pass
- [ ] Documentation is updated
- [ ] No hardcoded values or secrets
- [ ] Error handling is appropriate
- [ ] No performance regressions
- [ ] Commit messages are clear

---

## Quick Reference Checklist

When starting new work, use this checklist:

### Backend (Go)
- [ ] Create service in `go/services/{entity}_service.go`
- [ ] Define DTOs in `go/dtos/{entity}_dto.go`
- [ ] Create validators in `go/validators/{entity}_validator.go`
- [ ] Implement CRUD methods following the pattern
- [ ] Handle errors with domain-specific types
- [ ] Add unit tests
- [ ] Create Wails bindings in `go/app/{entity}_bindings.go`

### Frontend (Angular)
- [ ] Create models in `src/app/core/models/{entity}.model.ts`
- [ ] Create service in `src/app/core/services/{entity}.service.ts`
- [ ] Create container component in `src/app/features/{entity}/{entity}-list`
- [ ] Create form component in `src/app/features/{entity}/{entity}-form`
- [ ] Create presentational components as needed
- [ ] Add unit tests
- [ ] Update module imports
- [ ] Add routing

### General
- [ ] Write clear commit messages
- [ ] Update TODO.md if scope changes
- [ ] Add/update documentation
- [ ] Request code review

---

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Angular Documentation](https://angular.io/docs)
- [Wails Documentation](https://wails.io/docs)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Clean Code: A Handbook of Agile Software Craftsmanship](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)

---

**Last Updated:** 2026-10-02  
**Author:** ElectroZombie  
**Branch:** main

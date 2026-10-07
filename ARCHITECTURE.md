# 🏗️ ARCHITECTURE.md - Cargo Management App

**Purpose:** Describe the system architecture, layers, data model, and conventions of the Cargo Management App  
**Last Updated:** 2026-10-06  
**Branch:** main

---

## 📖 Table of Contents

1. [Overview](#overview)
2. [Technology Stack](#technology-stack)
3. [High-Level Architecture](#high-level-architecture)
4. [Layered Design](#layered-design)
5. [Data Model](#data-model)
6. [Frontend Architecture](#frontend-architecture)
7. [IPC and Error Handling](#ipc-and-error-handling)
8. [Directory Structure](#directory-structure)
9. [Request Lifecycle](#request-lifecycle)
10. [Theming and Design System](#theming-and-design-system)
11. [Build and Distribution](#build-and-distribution)
12. [Design Decisions and Constraints](#design-decisions-and-constraints)

---

## Overview

The Cargo Management App is a **single-tenant desktop application** for managing the day-to-day
operations of a hauling/cargo business: drivers, loaders, vehicles, wells, and the loads that tie
them together.

The application is built with **Wails**, which packages a **Go** backend and an **Angular**
frontend into a single native desktop binary. The browser-based Angular UI communicates with the
Go backend through Wails' generated JavaScript bindings rather than a network API.

### Goals

- Fast, offline-first operation backed by a local **SQLite** database.
- A clean separation between persistence, business logic, and presentation.
- A consistent, responsive UI across desktop, tablet, and mobile window sizes.
- Domain validation and data integrity enforced close to the database.

### Non-Goals

The following are **intentionally excluded** from the product:

- Authentication, login, registration, sessions, or tokens.
- User accounts, roles, permissions, AuthGuards, or RoleGuards.
- Multi-tenancy or remote/server deployments.
- A public HTTP/REST API surface.

---

## Technology Stack

| Layer | Technology | Notes |
|-------|-----------|-------|
| Desktop shell | Wails v2 | Native window + asset server + Go bindings |
| Backend language | Go 1.21 | Standard library + `database/sql` |
| Database | SQLite | Embedded, single-file, foreign keys enabled |
| Frontend framework | Angular | NgModule-based feature architecture |
| Frontend language | TypeScript | Strict mode enabled |
| State | RxJS `BehaviorSubject` | Service-owned reactive state |
| Styling | SCSS | Shared variables/mixins, soft color palette |
| Charts | Chart.js + ng2-charts | Optional, dashboard visualizations |
| Testing | Go `testing` + Angular TestBed | Unit and integration tests |

---

## High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                        Wails Desktop Shell                    │
│                                                              │
│   ┌────────────────────────┐       ┌──────────────────────┐  │
│   │    Angular Frontend    │       │    Go Backend        │  │
│   │  (WebView / TypeScript)│       │  (native process)    │  │
│   │                        │  IPC  │                      │  │
│   │  Feature Components    │◄─────►│  Wails App Bindings  │  │
│   │  Core Services         │       │  Service Layer       │  │
│   │  Wails Bridge          │       │  Validators          │  │
│   │  Models / DTOs         │       │  Database (SQLite)   │  │
│   └────────────────────────┘       └──────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

- The Angular app is compiled to static assets and embedded in the Go binary via
  `//go:embed all:frontend/dist`.
- Wails generates typed bindings that expose exported `App` methods to the frontend.
- All calls are **synchronous request/response** over the Wails runtime bridge.

---

## Layered Design

The backend follows a strict, unidirectional layering model:

```
HTTP-like boundary          Wails Bindings            (go/app or main.go)
                             │  thin wrappers only
Domain services             Service Layer             (database/*_service.go)
                             │  validation + business rules
Domain entities             Models + DTOs             (database/models.go)
                             │  scanner helpers
Persistence                 SQLite via database/sql   (database/*.go)
```

### Rules

1. **Bindings are thin.** They translate frontend input into service calls, wrap the result in a
   response envelope, and never contain business logic.
2. **Services own business rules.** Validation, reference checks, and conflict handling live here.
3. **Models are plain data.** Read models include generated IDs and timestamps; write DTOs exclude
   them.
4. **The database is the source of truth.** Constraints, uniqueness, and foreign keys are enforced
   at the schema level as a final safety net.

### Persistence Layer

- `database/database.go` opens the SQLite connection, enables foreign-key enforcement, and
  creates the schema.
- `database/schema.go` is the single source of truth for the operational schema and indexes.
- `database/migrations.go` provides a versioned migration manager with rollback support, tracked
  in the `schema_migrations` table.
- `database/helpers.go` provides validation primitives, domain errors, row scanners, and
  SQLite-error-to-domain-error translation.
- `database/transactions.go` provides transaction helpers for multi-entity operations.
- `database/listing.go` provides pagination, sorting, and combined filtering for list queries.
- `database/lcs.go` provides fuzzy string similarity used by search.

---

## Data Model

The operational schema is centered on five core entities plus financial records.

### Core entities

| Entity | Key fields | Relationships |
|--------|-----------|---------------|
| `driver` | `name`, `license_number` (unique), `phone`, `email`, `address` | referenced by vehicles |
| `loader` | `name`, `phone`, `email`, `address`, `company` | referenced by loads |
| `vehicle` | `vec_id`, `license_id`, `trailer_number`, `vin` (unique), `license_expiration`, `overweight_permit_id` | belongs to a driver |
| `well` | `name`, `client`, `location` | referenced by loads |
| `load` | `date`, `load_number` (unique), `ticket_number`, `miles`, `net_weight`, `tons`, `ton_rate`, `total_value`, `status`, `ticket_id` | belongs to loader, well, driver, vehicle |

### Additional entities

| Entity | Purpose |
|--------|---------|
| `category` | Expense/income classification |
| `expense` | Operating expenses linked to a category |
| `income` | Revenue linked to a category |
| `fuel_expense` | Fuel-specific expense linked to expense + driver |
| `stats` | Aggregated weekly financial summaries |

### Integrity rules

- Primary keys are auto-incrementing integers.
- `created_at` / `updated_at` are managed by the database.
- Unique constraints protect `driver.license_number`, `vehicle.vin`, `load.load_number`, and
  `category.type`.
- Foreign keys restrict deletion of referenced records.
- Indexes cover relationship columns (`driver_id`, `well_id`, `vec_id`, `status`) and date columns.

> **Synchronization note:** The Go read model and the SQL schema must be kept in sync. The
> frontend TypeScript models mirror the Go read model field-for-field using camelCase.

---

## Frontend Architecture

The Angular frontend is organized **by feature, not by type**:

```
Core        →  singletons: models, Wails bridge, entity services, error handling
Shared      →  reusable presentational components, pipes, styles
Layout      →  shell: navbar, sidebar, breadcrumb, footer
Features    →  self-contained screens (dashboard, drivers, loaders, vehicles, wells, loads)
```

### State management

There is no global store (no NgRx). Each entity service owns its state through RxJS
`BehaviorSubject`s and exposes it as observables:

- `items$` — the current collection
- `selected$` — the currently selected entity
- `loading$` — in-flight request indicator

Components subscribe via the `async` pipe and dispatch actions back to the service.

### Component model

- **Container/smart components** live in `features/*`. They inject services, manage state, and
  handle routing.
- **Presentational components** live in `shared/components`. They accept data through `@Input()`,
  emit through `@Output()`, and are free of service dependencies.
- All components use `ChangeDetectionStrategy.OnPush`.

### Forms

- Reactive Forms with strongly typed `FormGroup`s.
- Validators mirror backend validation rules (required, min length, patterns).
- Backend domain errors are surfaced as field-level or toast messages.

### Tables

A shared `data-table` component provides sorting, filtering, pagination, and row selection. It is
generic and driven by column definitions supplied by each feature.

---

## IPC and Error Handling

### Response envelope

Every Wails binding returns a consistent JSON envelope so the frontend never has to interpret raw
Go errors:

```typescript
interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
```

### Domain errors

The backend defines domain error types (validation, not-found, conflict, internal). Raw SQLite
errors are never exposed. `databaseError` maps constraint failures to `ErrConflict` /
`ErrValidation`, and bindings convert domain errors into human-readable `error` strings.

### Frontend handling

- A centralized Angular `ErrorHandler` catches uncaught errors.
- A `NotificationService` renders toast notifications for success/failure.
- Services translate a `success: false` response into a user-facing message and log the detail.

---

## Directory Structure

```
Cargo Management App/
├── main.go                     # Entry point + Wails bootstrapping + bindings
├── wails.json                  # Wails project configuration
├── go.mod                      # Go module definition
├── database/                   # Go backend (schema, services, models, validation)
│   ├── database.go
│   ├── schema.go
│   ├── migrations.go
│   ├── models.go
│   ├── helpers.go
│   ├── crud.go
│   ├── listing.go
│   ├── transactions.go
│   ├── lcs.go
│   └── *_service.go
└── frontend/                   # Angular application
    ├── package.json
    ├── angular.json
    ├── tsconfig*.json
    └── src/
        ├── index.html
        ├── main.ts
        ├── styles/
        ├── environments/
        └── app/
            ├── app.module.ts
            ├── app-routing.module.ts
            ├── core/
            │   ├── models/
            │   └── services/
            ├── shared/
            │   ├── components/
            │   ├── styles/
            │   └── shared.module.ts
            ├── layout/
            └── features/
```

---

## Request Lifecycle

A typical read operation (for example, listing drivers):

```
Angular Component
  └─ DriverService.loadDrivers()
       └─ WailsService.invoke('ListDrivers', params)
            └─ window.runtime.Call.ByName('ListDrivers', params)   [IPC]
                 └─ App.ListDrivers(ctx, params)                   [Go binding]
                      └─ Database.ListDrivers(filters, opts)        [service]
                           └─ SQLite SELECT + COUNT                [persistence]
                 ◄─ DriverListResponse{success, data}              [envelope]
            ◄─ resolved envelope
       └─ items$.next(response.data)
  ◄─ async pipe re-renders the view
```

A write operation additionally runs validation and reference checks, then re-fetches the created
or updated entity before returning it to the frontend.

---

## Theming and Design System

A soft, professional palette is defined once in `shared/styles/_variables.scss` and reused through
SCSS `@import`. It covers:

- **Colors:** primary `#5A9FD4`, secondary `#6BBD9F`, accent `#F5A28B`, plus neutrals, text,
  and semantic colors.
- **Spacing:** 4 / 8 / 16 / 24 / 32 px scale.
- **Typography:** Roboto; headings 24–32px, body 14–16px.
- **Radius:** 4 / 8 / 16 px.
- **Shadows:** small / medium / large elevation tokens.

Dark mode is implemented by overriding the same CSS custom properties under a `.dark-theme`
class/attribute, so components need no changes.

---

## Build and Distribution

| Command | Purpose |
|---------|---------|
| `wails dev` | Run with hot reload (Angular dev server + Go) |
| `wails build` | Produce a production desktop binary in `/build` |
| `npm run build` | Build Angular assets only (inside `frontend/`) |
| `go test ./...` | Run Go unit tests |
| `npm test` | Run Angular unit tests |

The compiled Angular assets in `frontend/dist` are embedded into the Go binary at build time. The
SQLite database is stored per-user at `~/.cargo_management_app/cargo_management.db`.

---

## Design Decisions and Constraints

| Decision | Rationale |
|----------|-----------|
| Local SQLite, no server | Offline-first, single-user operational tool |
| Embedded asset server | No separate frontend deployment or CORS |
| Wails bindings instead of REST | Type-safe, low-latency, no network surface |
| No authentication | Product is single-tenant and runs on the user's machine |
| BehaviorSubject over NgRx | Sufficient for entity-scoped state; less ceremony |
| Database-level constraints | Guarantees integrity even if a service path is bypassed |
| Feature-based folders | Scales with new operational modules; easy navigation |

---

## Related Documents

- [README.md](./README.md) — Project overview and setup
- [AGENTS.md](./AGENTS.md) — Development best practices
- [CONTRIBUTING.md](./CONTRIBUTING.md) — Contribution workflow
- [TODO.md](./TODO.md) — Development backlog and status

---

**Author:** ElectroZombie  
**Branch:** main

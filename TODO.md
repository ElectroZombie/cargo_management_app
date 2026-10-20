# 📋 Cargo Management App - Development To-Do List

**Project Status:** In Development — backend schema, validation, migration, seed, and transaction foundation completed; Angular Material and Chart.js installed but **not yet integrated**; Go package layout does not yet match AGENTS.md
**Last Updated:** 2026-10-20  
**Branch:** main

---

## ✅ Completed

### Backend database foundation
- [x] Removed user-login, authentication, roles, tokens, and user-account functionality.
- [x] Removed the `users` table from the operational database design.
- [x] Defined the operational schemas for:
  - [x] `driver`
  - [x] `loader`
  - [x] `vehicle`
  - [x] `well`
  - [x] `load`
- [x] Added primary keys, required fields, unique constraints, timestamps, and foreign keys.
- [x] Enabled SQLite foreign-key enforcement.
- [x] Added relationship indexes for vehicles and loads.
- [x] Added delete restrictions for referenced drivers, loaders, vehicles, and wells.

### DTOs and domain models
- [x] Created read models for drivers, loaders, vehicles, wells, and loads.
- [x] Created write DTOs:
  - [x] `DriverDTO`
  - [x] `LoaderDTO`
  - [x] `VehicleDTO`
  - [x] `WellDTO`
  - [x] `LoadDTO`
- [x] Excluded generated IDs and timestamps from write DTOs.
- [x] Created read models and write DTOs for the finance entities (`category`, `expense`, `income`, `fuel_expense`, `stats`) — **models and DTOs only; see Finance domain below.**

### Database services
- [x] Implemented driver CRUD services.
- [x] Implemented loader CRUD services.
- [x] Implemented vehicle CRUD services.
- [x] Implemented well CRUD services.
- [x] Implemented load CRUD services.
- [x] Added load queries by driver, well, and active/inactive status.
- [x] Added required-field validation for core create operations.
- [x] Added not-found handling for reads, updates, and deletes.
- [x] Added search support for drivers, loaders, vehicles, wells, and loads.

### Database quality and operations
- [x] Added a versioned migration system with rollback support.
- [x] Added database seed data and repeatable test fixtures.
- [x] Added complete validation for DTOs, including:
  - [x] Date and date-range validation.
  - [x] Positive numeric values for mileage, weights, rates, and totals.
  - [x] Valid foreign-key references before writes.
  - [x] VIN, license, ticket, and identifier format validation.
- [x] Returned domain-specific validation and conflict errors instead of raw SQLite errors.
- [x] Added transaction helpers for multi-entity operations.
- [x] Added pagination, sorting, and combined filtering to list services.
- [x] Added relationship, status, and date indexes; parameterized filter construction via `listWhere`; `LIMIT`/`OFFSET` pagination.

### Wails/API integration
- [x] Added Wails application methods that expose the new DTO-based services.
- [x] Added consistent response DTOs for the frontend (`Response` envelope with `success`/`data`/`error`/`total`).
- [x] Added application-level error handling and logging (`log/slog`).
- [x] Added frontend service wrappers for drivers, loaders, vehicles, wells, and loads.
- [x] Added dashboard statistics aggregation (`GetDashboardStats`).

### Frontend foundation
- [x] Created shared Button, Card, Modal, LoadingSpinner, Toast, and ConfirmDialog components.
- [x] Created a generic `data-table` component with sorting, filtering, pagination, and row selection.
- [x] Created driver, loader, vehicle, well, and load management list/form screens.
- [x] Added reusable form validation and error messages (`shared/validators/positive-number.validator.ts`).
- [x] Added loading states, success notifications, and error handling.
- [x] Added responsive mobile, tablet, and desktop layouts.
- [x] Completed Navbar, Sidebar, Layout, Breadcrumb, and Footer components.
- [x] Created the dashboard layout with header, sidebar, content area, and footer.
- [x] Added operational statistics for loads, weight, active loads, drivers, vehicles, and wells.
- [x] Added recent activity and quick-access actions.
- [x] Created routes for dashboard, drivers, loaders, vehicles, wells, and loads.
- [x] Added collapsible sidebar navigation and active route indicators.
- [x] Added top navigation with search, notifications, and settings.
- [x] Added dynamic breadcrumb navigation.
- [x] Added a `ThemeService` with light/dark switching persisted to `localStorage`, wired to the navbar toggle.

### Application cleanup
- [x] Removed the old cargo/shipment/user-oriented Wails bindings from `main.go`.
- [x] Kept application startup and shutdown connected to the new database package.
- [x] Moved the backend under `go/` to match the AGENTS.md package layout.

### Error contract and binding layer
- [x] Consolidated all domain errors into `go/errors` with sentinel + typed pairs and `Unwrap` support.
- [x] Added a domain-aware error formatter that returns friendly messages and stable error codes.
- [x] Split the Wails bindings into per-entity files under `go/app/`.
- [x] Added service tests for the loader, vehicle, and well entities.
- [x] Added binding-layer integration tests covering the response envelope and error codes.

---

## 🚧 Remaining backend work

### Go package layout alignment (AGENTS.md conformance)
> The implementation currently lives in `go/database/`, but AGENTS.md designates a separate
> `go/services/` layer. `go/errors/`, `go/utils/`, and `go/validators/` are currently dead or
> stub code, so the documented boundaries are not enforced in code.

- [ ] Create the `go/services/` package and move the five entity services (`driver`, `loader`, `vehicle`, `well`, `load`) out of `go/database/`.
- [x] Make `go/errors/errors.go` the single source of domain errors and delete the duplicate error types in `go/database/helpers.go`. — `go/errors` now owns the sentinels (`ErrValidation`, `ErrNotFound`, `ErrConflict`, `ErrInternal`) and the typed errors (`ValidationError`, `NotFoundError`, `ConflictError`, `InternalError`), each implementing `Unwrap` so `errors.Is`/`errors.As` work through wrapping. `go/database/helpers.go` re-exports the sentinels and the `ValidationError` alias.
- [ ] Replace the `go/validators/driver_validator.go` stub (currently returns `nil`) with real implementations, and add:
  - [ ] `go/validators/driver_validator.go`
  - [ ] `go/validators/loader_validator.go`
  - [ ] `go/validators/vehicle_validator.go`
  - [ ] `go/validators/well_validator.go`
  - [ ] `go/validators/load_validator.go`
  - [ ] `go/validators/shared_validator.go`
- [ ] Move the inlined validation helpers from `go/database/helpers.go` into `go/validators/`.
- [ ] Either use `go/utils/constants.go` for pagination/sort defaults or delete it.
- [x] Split the monolithic `go/app/app.go` into per-entity binding files (`driver_bindings.go`, `loader_bindings.go`, `vehicle_bindings.go`, `well_bindings.go`, `load_bindings.go`, `dashboard_bindings.go`). — `app.go` now holds only the `App` struct, lifecycle hooks, and the `Response` envelope; shared error handling lives in `go/app/errors.go`.
- [ ] Accept `context.Context` as the first parameter of every Wails binding, per AGENTS.md and the ARCHITECTURE.md request lifecycle. — **Blocked on the framework:** Wails v2.8.1's `BoundMethod.ParseArgs` requires an exact argument-count match and JSON-decodes every argument, so a `context.Context` parameter cannot be injected and would break the frontend call. The app context is already captured in `App.Startup` and held on the `App` struct; revisit if the project moves to Wails v3.
- [x] Replace `formatError` with a domain-aware formatter that maps validation/not-found/conflict/internal errors to user-facing messages instead of returning `err.Error()`. — `go/app/errors.go` classifies errors into `validation`, `not_found`, `conflict`, `internal`, and `unknown`, returns a friendly message, and never leaks raw SQLite internals (internal failures are logged via `slog` and reduced to a generic message). The `Response` envelope now carries a machine-readable `code` alongside `error`.
- [ ] Split `go/database/` into a persistence layer and a service layer so the layered design in ARCHITECTURE.md is real.

### Finance domain
> Tables, models, and DTOs exist. Services, bindings, and UI do not. TODO previously had no finance items at all.

- [ ] Implement `category` CRUD service and validators.
- [ ] Implement `expense` CRUD service and validators.
- [ ] Implement `income` CRUD service and validators.
- [ ] Implement `fuel_expense` CRUD service and validators.
- [ ] Implement `stats` aggregation service (weekly summaries).
- [ ] Expose the finance entities through Wails bindings.
- [ ] Add finance screens: categories, expenses, income, fuel expenses, and a financial summary report.
- [ ] Add the Reports route as a real feature instead of a `ComingSoonComponent` placeholder.

### Wails/API integration (remaining)
- [ ] Add Wails methods for the finance entities.
- [ ] Add an application-level logger with request-level context (operation, entity ID, duration).

### Database quality and operations (remaining)
- [ ] Add service and database tests for the missing entities:
  - [x] `loader_service_test.go`
  - [x] `vehicle_service_test.go`
  - [x] `well_service_test.go`
  - [ ] Finance service tests
- [ ] Add tests for the listing layer (pagination, sorting, combined filtering).
- [ ] Add tests for migrations (apply + rollback) and transaction helpers.
- [ ] Add tests for the dashboard aggregation queries.
- [ ] Add tests for the LCS search implementation.
- [ ] Add validator unit tests in `go/validators/`.
- [x] Add Wails binding-layer integration tests in `go/app/`. — `go/app/app_test.go` covers the response envelope, success and failure paths, error codes (`validation`, `not_found`, `conflict`, `internal`), pagination totals, filtering, dashboard stats, and the guarantee that internal error detail never reaches the frontend.

---

## 🚧 Remaining frontend work

### Angular Material integration
> `@angular/material` and `@angular/cdk` are installed and a theme is emitted, but there is
> **zero `mat-` usage in any template**. Material is currently unused bundle weight.

- [ ] Migrate shared components to Material widgets (buttons, form fields, cards, dialogs, menus, tooltips, snackbars).
- [ ] Migrate `data-table` to `MatTable` with `MatSort` and `MatPaginator`.
- [ ] Migrate all feature forms to `MatFormField`, `MatInput`, `MatSelect`, and `MatDatepicker`.
- [ ] Replace the custom `ModalComponent` with `MatDialog`, or keep the custom modal and adopt Material styling consistently.
- [ ] Replace the custom `ToastComponent` with `MatSnackBar`, or keep the custom toast and adopt Material styling consistently.
- [ ] Add Material icons and typography to the global stylesheet.
- [ ] Remove the duplicate theme emission (see Build and tooling below).

### Charts and dashboard visualization
> `chart.js` and `ng2-charts` are installed and `provideCharts(withDefaultRegisterables())` is
> registered, but **no chart components exist**. The backend aggregation queries the charts need
> are also missing.

- [ ] Add backend aggregation endpoints for chart data:
  - [ ] Loads over time (per day/week/month)
  - [ ] Weight and tons trends
  - [ ] Load status distribution
  - [ ] Delivery performance (on-time vs. late)
  - [ ] Revenue vs. expense summary
- [ ] Create shared, responsive chart components:
  - [ ] Line chart component
  - [ ] Pie/doughnut chart component
  - [ ] Bar chart component
- [ ] Add the load, weight, status, and delivery-performance charts to the dashboard.
- [ ] Make charts respect the active theme (light and dark) and re-render on theme change.

### Feature completeness
- [ ] Add detail screens for each entity (`driver-detail`, `loader-detail`, `vehicle-detail`, `well-detail`, `load-detail`) — AGENTS.md specifies these; only list and form exist today.
- [ ] Add filters to the load list (driver, well, loader, vehicle, status, date range).
- [ ] Add status controls to the load list (activate/deactivate) and to the load form.
- [ ] Add filters to the driver, loader, vehicle, and well lists.
- [ ] Wire the `ConfirmDialogComponent` into the delete flows instead of the native `confirm()`.
- [ ] Implement the Settings feature (theme preference, page size, defaults) instead of the `ComingSoonComponent` placeholder.
- [ ] Implement the Reports feature (see Finance domain above).

### Animations and transitions
- [ ] Add route/page transitions using `@angular/animations` (currently installed but unused).
- [ ] Add hover effects and list/stagger animations.
- [ ] Add loading-state animations for the data table and dashboard cards.

### Frontend testing
- [ ] Add Angular unit tests — there are currently **zero** `.spec.ts` files, so `npm test` has nothing to run.
  - [ ] Service specs for driver, loader, vehicle, well, load, dashboard, theme, and notification services.
  - [ ] Component specs for presentational shared components (input/output behavior).
  - [ ] Container component specs for the list and form screens.
  - [ ] A shared `WailsService` test double / mock helper.

### Shared and operational components
- [ ] Populate the empty `frontend/src/app/shared/pipes/` directory or remove it.
- [x] Surface backend error codes in the frontend. — `ApiResponse`/`ListResponse` carry an optional `code`, and `CrudService` records the last error in a `lastError()` observable so forms can display a backend validation message next to the offending field.

---

## 🛠️ Build and tooling

- [ ] Add a CI pipeline (`.github/workflows/`) — CONTRIBUTING.md references "CI is green" but no CI exists. Steps: `go vet`, `go test ./...`, `npm run lint`, `npm test`, `npm run build`, and `wails build`.
- [ ] Add ESLint: install `@angular-eslint/*` and add an `.eslintrc.json`. `npm run lint` is currently declared but will fail.
- [ ] Add Prettier: install it and add a `.prettierrc` matching the documented style (2-space indent, single quotes).
- [ ] Add `.editorconfig`.
- [ ] Add a `LICENSE` file.
- [ ] Add a `CODEOWNERS` file.
- [ ] Fix the duplicate Material theme emission — two files generate a theme (`frontend/src/styles/theme.scss` and `frontend/src/app/shared/styles/theme.scss`), producing "same color/density/typography styles are generated multiple times" warnings.
- [ ] Raise or restructure the Angular bundle budget — the initial bundle is 825.39 kB against a 500 kB warning / 1 MB error budget.
- [ ] Remove the leftover empty `database/` directory at the repository root.
- [ ] Add a release/binary distribution step to the build pipeline.

---

## 📚 Documentation maintenance

- [ ] Update ARCHITECTURE.md §Directory Structure to show the real `go/` package layout (`go/database`, `go/models`, `go/dtos`, `go/errors`, `go/validators`, `go/utils`, `go/app`) and remove the stale root-level `database/` tree.
- [ ] Update ARCHITECTURE.md §Request Lifecycle to match the actual binding signatures.
- [ ] Reconcile the dark-mode selector across all docs — ARCHITECTURE.md says `.dark-theme`, TODO.md says `[data-theme='dark']`, the code uses the `data-theme` attribute on `<html>`.
- [ ] Update ARCHITECTURE.md §Additional entities to state clearly that finance entities are schema-only until the finance work lands.
- [ ] Update README.md — it states Go 1.18+ / Node 16+ (CONTRIBUTING.md says 1.21+ / 18+), omits the `go/` directory, and lists a non-existent `frontend/styles.scss`.
- [ ] Replace MISSING_ELEMENTS.md and IMPLEMENTED_SUMMARY.md with the current analysis (both contain obsolete claims, e.g. that the root `database/` package is still functional and that `MatChartsModule` is used instead of `provideCharts`).
- [ ] Align the "Last Updated" dates across AGENTS.md, ARCHITECTURE.md, and TODO.md.
- [ ] Keep TODO.md status current when completing or adding work.

---

## 📊 Current implementation status

| Area | Status | Notes |
|------|--------|-------|
| Authentication and login | ✅ Removed | Intentionally excluded from the product design |
| Operational database schema | ✅ Complete | Driver, loader, vehicle, well, and load tables defined |
| DTOs and read models | ✅ Complete | DTOs created for all operational entities |
| Database CRUD services | ✅ Complete | Core CRUD and load-specific queries implemented |
| Validation | ✅ Complete | Full DTO validation lives in `go/database/helpers.go` — **not yet extracted to `go/validators/`** |
| Migrations and seed data | ✅ Complete | Versioned migration system and repeatable seed fixtures implemented |
| Transaction helpers | ✅ Complete | Transaction wrappers for multi-entity operations implemented |
| Wails service bindings | ✅ Complete | Go bindings implemented in `go/app` — **split per entity, domain-aware error formatting, machine-readable `code`** |
| Frontend screens and services | ⚠️ Partial | List and form screens for all five entities; **detail screens, filters, and status controls missing** |
| Dashboard and navigation | ✅ Complete | Layout, navigation, stats, routes, recent activity, and quick actions implemented |
| Angular Material | ⚠️ Partial | Installed and themed, but **zero `mat-` component usage** |
| Chart.js / ng2-charts | ⚠️ Partial | Installed and registered, but **zero chart components and no chart data endpoints** |
| Dark mode | ✅ Complete | `ThemeService` with persisted preference, navbar toggle, CSS custom properties |
| Animations and transitions | ❌ Not started | `@angular/animations` installed but unused |
| Domain error handling | ✅ Complete | Single source in `go/errors`; typed errors with `Unwrap`; friendly messages and stable error codes at the binding layer |
| Finance domain | ❌ Not started | Tables, models, and DTOs only — no services, bindings, or UI |
| Go package layout | ⚠️ Partial | `go/errors` is now authoritative; **`go/services/` still absent, `go/utils` still dead, `go/validators` still a stub** |
| Go tests | ⚠️ Partial | All five core entities plus the binding layer covered; **validator, listing, migration, transaction, and finance tests missing** |
| Angular tests | ❌ Not started | Zero `.spec.ts` files |
| CI/CD | ❌ Not started | No pipeline despite CONTRIBUTING.md referencing CI |
| Linting and formatting | ❌ Not started | ESLint and Prettier referenced in docs but not installed |
| Documentation accuracy | ⚠️ Partial | ARCHITECTURE.md and README.md contradict the current code layout |

---

## 🎨 Design guidelines

### Color palette
```scss
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
```

### Typography and spacing
- **Font:** Roboto, sans-serif
- **Headings:** 24–32px, bold
- **Body:** 14–16px
- **Spacing:** 4px, 8px, 16px, 24px, 32px
- **Border radius:** 4px, 8px, 16px

---

## 🔗 Related documents

- [x] ARCHITECTURE.md — System architecture
- [x] CONTRIBUTING.md — Contribution guidelines
- [x] MISSING_ELEMENTS.md — Current gap analysis
- [README.md](./README.md) — Project overview

---

**Updated By:** ElectroZombie  
**Branch:** main

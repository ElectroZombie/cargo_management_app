# 📋 Cargo Management App - Development To-Do List

**Project Status:** In Development — backend schema and service foundation completed  
**Last Updated:** 2026-09-18  
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

### Database services
- [x] Implemented driver CRUD services.
- [x] Implemented loader CRUD services.
- [x] Implemented vehicle CRUD services.
- [x] Implemented well CRUD services.
- [x] Implemented load CRUD services.
- [x] Added load queries by driver, well, and active/inactive status.
- [x] Added required-field validation for core create operations.
- [x] Added not-found handling for reads, updates, and deletes.

### Application cleanup
- [x] Removed the old cargo/shipment/user-oriented Wails bindings from `main.go`.
- [x] Kept application startup and shutdown connected to the new database package.

---

## 🚧 Remaining backend work

### Database quality and operations
- [ ] Add a versioned migration system with rollback support.
- [ ] Add database seed data and repeatable test fixtures.
- [ ] Add complete validation for every DTO, including:
  - [ ] Date and date-range validation.
  - [ ] Positive numeric values for mileage, weights, rates, and totals.
  - [ ] Valid foreign-key references before writes.
  - [ ] VIN, license, ticket, and identifier format validation.
- [ ] Return domain-specific validation and conflict errors instead of raw SQLite errors.
- [ ] Add transaction helpers for multi-entity operations.
- [ ] Add pagination, sorting, and combined filtering to list services.
- [ ] Add search support where required.
- [ ] Add service and database tests, including foreign-key and uniqueness tests.

### Wails/API integration
- [ ] Add Wails application methods that expose the new DTO-based services.
- [ ] Add consistent response DTOs for the frontend.
- [ ] Add application-level error handling and logging.
- [ ] Add frontend service wrappers for drivers, loaders, vehicles, wells, and loads.

---

## 🚧 Remaining frontend work

### Shared and operational components
- [ ] Create shared Button, Card, Modal, LoadingSpinner, Toast, and ConfirmDialog components.
- [ ] Create driver management screens.
- [ ] Create loader management screens.
- [ ] Create vehicle management screens.
- [ ] Create well management screens.
- [ ] Create load management screens, including forms, filters, details, and status controls.
- [ ] Add reusable form validation and error messages.
- [ ] Add data tables with sorting, filtering, pagination, and row selection.
- [ ] Add loading states, success notifications, retry actions, and error handling.
- [ ] Add responsive mobile, tablet, and desktop layouts.
- [ ] Complete Navbar, Sidebar, Layout, and Breadcrumb components.

### Dashboard and navigation
- [ ] Create the dashboard layout with header, sidebar, content area, and footer.
- [ ] Add operational statistics for loads, weight, active loads, drivers, vehicles, and wells.
- [ ] Add recent activity and quick-access actions.
- [ ] Create routes for dashboard, drivers, loaders, vehicles, wells, loads, reports, and settings.
- [ ] Add collapsible sidebar navigation and active route indicators.
- [ ] Add top navigation with search, notifications, and settings.
- [ ] Add dynamic breadcrumb navigation.
- [ ] Do not add login, registration, user-management, AuthGuard, or RoleGuard features.

### Styling and visualization
- [ ] Install and configure Angular Material.
- [ ] Create the Material Design theme and soft color palette.
- [ ] Add global SCSS variables for colors, spacing, typography, shadows, and radius.
- [ ] Build reusable styled tables, forms, cards, and modals.
- [ ] Add dark-mode support.
- [ ] Add page transitions, hover effects, loading animations, and toast notifications.
- [ ] Install Chart.js and `ng2-charts` if charting is still required.
- [ ] Create responsive line, pie, and bar chart components.
- [ ] Add load, weight, status, and delivery-performance charts.

---

## 📊 Current implementation status

| Area | Status | Notes |
|------|--------|-------|
| Authentication and login | ✅ Removed | Intentionally excluded from the product design |
| Operational database schema | ✅ Complete | Driver, loader, vehicle, well, and load tables defined |
| DTOs and read models | ✅ Complete | DTOs created for all operational entities |
| Database CRUD services | ✅ Complete | Core CRUD and load-specific queries implemented |
| Validation | ⚠️ Partial | Required-field validation exists; full domain validation remains |
| Migrations and seed data | ❌ Not started | Still required |
| Wails service bindings | ❌ Not started | New services are not yet exposed to the frontend |
| Frontend screens and services | ❌ Not started | Operational UI remains to be built |
| Dashboard and navigation | ❌ Not started | Must be designed without authentication screens |
| Charts and Material styling | ❌ Not started | Optional chart dependencies remain |

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

- [README.md](./README.md) — Project overview
- [ARCHITECTURE.md](./ARCHITECTURE.md) — System architecture (to be created)
- [CONTRIBUTING.md](./CONTRIBUTING.md) — Contribution guidelines (to be created)

---

**Updated By:** ElectroZombie  
**Branch:** main

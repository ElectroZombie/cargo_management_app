# 📋 Cargo Management App - Development To-Do List

**Project Status:** In Development  
**Last Updated:** 2026-09-15  
**Branch:** main

---

## 🎯 Main Development Phases

### Phase 1: Database Schema & Backend CRUD Operations
#### 1.1 Create and Fill Database Schema ⭐ **CRITICAL**
- [ ] **1.1.1** Finalize and validate cargo table schema
  - Priority: 🔴 **CRITICAL**
  - Status: Schema exists, needs review
  - Details: Ensure all fields are properly defined (id, name, status, weight, destination, date, created_at, updated_at)
  - Extra: Add indexes on frequently queried fields (status, destination)
- [ ] **1.1.2** Finalize and validate shipments table schema
  - Priority: 🔴 **CRITICAL**
  - Status: Schema exists, needs review
  - Details: Validate foreign key relationships and constraints
  - Extra: Add indexes on cargo_id and status
- [ ] **1.1.3** Complete users table implementation
  - Priority: 🔴 **CRITICAL**
  - Status: Schema exists, not implemented
  - Details: username, email, password, role, timestamps
  - Extra: Add unique constraints, password hashing support
- [ ] **1.1.4** Create database migration system
  - Priority: 🟡 **MEDIUM**
  - Status: Not started
  - Details: Set up version control for database changes
  - Extra: Add rollback capability
- [ ] **1.1.5** Add database seed data (test fixtures)
  - Priority: 🟡 **MEDIUM**
  - Status: Not started
  - Details: Create sample cargo and shipment data for testing
  - Extra: Generate realistic test data with multiple statuses

#### 1.2 Write CRUD Code for Database Tables ⭐ **CRITICAL**
- [ ] **1.2.1** Implement User CRUD operations
  - Priority: 🔴 **CRITICAL**
  - Status: Schema only
  - Tasks: registration, login, profile read/update, account deletion, and admin user listing
  - Extra: Add password validation, email verification
- [ ] **1.2.2** Enhance Cargo CRUD operations
  - Priority: 🔴 **CRITICAL**
  - Status: Partially implemented
  - Tasks: validate CreateCargo, add status/destination/date-range queries, bulk update/delete
  - Extra: Add cargo archival, soft delete support
- [ ] **1.2.3** Enhance Shipment CRUD operations
  - Priority: 🔴 **CRITICAL**
  - Status: Partially implemented
  - Tasks: implement list/get/update/delete and status/date-range queries
  - Extra: Add shipment history tracking
- [ ] **1.2.4** Add advanced query methods
  - Priority: 🟡 **MEDIUM**
  - Status: Not started
  - Tasks: search, pagination, sorting, and combined filtering
  - Extra: Full-text search support
- [ ] **1.2.5** Implement data validation
  - Priority: 🟡 **MEDIUM**
  - Status: Not started
  - Details: Add validation for all input parameters
  - Extra: Custom error messages, input sanitization

---

### Phase 2: Frontend Component Generation & Enhancement
#### 2.1 Generate Components for Frontend ⭐ **CRITICAL**
- [ ] **2.1.1** Create Shared Components Module
  - Priority: 🔴 **CRITICAL**
  - Status: Not started
  - Components: Button, Card, Modal, LoadingSpinner, ToastNotification, ConfirmDialog
  - Styling: Minimalistic, soft colors, Material Design
- [ ] **2.1.2** Generate Cargo Management Components
  - Priority: 🔴 **CRITICAL**
  - Status: Not started
  - Components: CargoList, CargoDetail, CargoForm, CargoFilter, CargoBulkActions
  - Styling: Soft colors, high contrast text
  - Extra: Export to CSV, print view
- [ ] **2.1.3** Generate Shipment Management Components
  - Priority: 🔴 **CRITICAL**
  - Status: Not started
  - Components: ShipmentList, ShipmentDetail, ShipmentForm, ShipmentTracking, ShipmentFilter
  - Styling: Soft colors with Material Design
  - Extra: Real-time status updates, tracking map
- [ ] **2.1.4** Generate User Management Components
  - Priority: 🟡 **MEDIUM**
  - Status: Not started
  - Components: Login, Register, UserProfile, UserList, UserForm
  - Extra: Password reset, two-factor authentication UI
- [ ] **2.1.5** Generate Layout Components
  - Priority: 🔴 **CRITICAL**
  - Status: Partially started
  - Components: Navbar, Sidebar, Layout, Breadcrumb
  - Styling: Minimalistic, soft colors, Material Design

#### 2.2 Connect Components to Database Tables 🔗
- [ ] **2.2.1** Create CargoService: getAllCargo, getCargoById, createCargo, updateCargo, deleteCargo, getCargoStats, filterCargo
- [ ] **2.2.2** Create ShipmentService: list/get shipments, get by cargo ID, create/update/delete
- [ ] **2.2.3** Create UserService: login, register, current user, update profile, logout
- [ ] **2.2.4** Create AuthService: authentication, roles, token access, token refresh
- [ ] **2.2.5** Integrate with Wails runtime: wrapper, error handling, and logging

#### 2.3 Upgrade Components Functionality 🚀
- [ ] **2.3.1** Add Angular Reactive Forms validation and real-time feedback
- [ ] **2.3.2** Add data tables with sorting, filtering, pagination, and row selection
- [ ] **2.3.3** Add loading states, errors, success notifications, and retry mechanisms
- [ ] **2.3.4** Add responsive design for mobile, tablet, and desktop
- [ ] **2.3.5** Install and apply Angular Material components, theme, and icons

---

### Phase 3: Main View & Navigation Design
#### 3.1 Design Main Dashboard View ⭐ **CRITICAL**
- [ ] **3.1.1** Create dashboard layout with header, sidebar, content area, and footer
- [ ] **3.1.2** Add quick stats for cargo count, total weight, in-transit shipments, and completed deliveries
- [ ] **3.1.3** Add Chart.js/ng2-charts visualizations for status, shipments, and weight
- [ ] **3.1.4** Add a recent activity feed with the last ten activities and timestamps
- [ ] **3.1.5** Create quick access actions for cargo, shipments, and reports

#### 3.2 Design Navigation System 🗺️
- [ ] **3.2.1** Create routes: /dashboard, /cargo, /cargo/:id, /cargo/create, /shipment, /shipment/:id, /reports, /settings, /users
- [ ] **3.2.2** Create a collapsible sidebar with icons, labels, active state, and animations
- [ ] **3.2.3** Create top navigation with title, global search, user menu, notifications, and settings
- [ ] **3.2.4** Implement AuthGuard, RoleGuard, and UnsavedChangesGuard
- [ ] **3.2.5** Add dynamic, clickable breadcrumb navigation

#### 3.3 Design UI/UX with Material & Soft Colors 🎨
- [ ] **3.3.1** Create the soft color scheme and dark-mode support
- [ ] **3.3.2** Implement the Material Design theme, typography, and icons
- [ ] **3.3.3** Create `frontend/src/styles/_variables.scss` with colors, spacing, typography, shadows, and radius
- [ ] **3.3.4** Build reusable styled buttons, cards, forms, tables, and modals
- [ ] **3.3.5** Add page transitions, hover effects, loading animations, and toast notifications

---

### Phase 4: Chart Integration & Dependencies
#### 4.1 Install & Configure Charts Library ⭐ **CRITICAL**
- [ ] **4.1.1** Install Chart.js and ng2-charts: `npm install chart.js ng2-charts`
- [ ] **4.1.2** Create reusable line, pie, and bar chart components
- [ ] **4.1.3** Implement cargo status, shipment timeline, weight distribution, and delivery performance charts
- [ ] **4.1.4** Add responsive canvas sizing and touch interactions

#### 4.2 Install Additional Dependencies 📦
- [ ] **4.2.1** Angular Material: `ng add @angular/material`
- [ ] **4.2.2** Angular Forms & Validators (built-in)
- [ ] **4.2.3** RxJS operators (built-in)
- [ ] **4.2.4** Optional date utilities: `date-fns` or `ng-date-pipe`
- [ ] **4.2.5** HTTP Client (built-in)

---

## 📊 Implementation Summary

| Phase | Component | Priority | Estimated Effort |
|-------|-----------|----------|------------------|
| 1 | Database Schema & CRUD | 🔴 CRITICAL | 8-10 days |
| 2 | Frontend Components | 🔴 CRITICAL | 10-12 days |
| 3 | Main View & Navigation | 🔴 CRITICAL | 5-7 days |
| 4 | Charts & Styling | 🔴 CRITICAL | 3-4 days |
| **Total** | **All Features** | - | **26-33 days** |

---

## 🎨 Design Guidelines

### Color Palette (Soft Colors with High Contrast)
```scss
$primary-color: #5A9FD4;
$secondary-color: #6BBD9F;
$accent-color: #F5A28B;
$background-color: #F7F8FA;
$neutral-light: #E8E9EB;
$neutral-medium: #B8BCC4;
$neutral-dark: #636E72;
$text-primary: #2D3436;
$text-secondary: #636E72;
$text-light: #95A5A6;
$success-color: #27AE60;
$warning-color: #F39C12;
$danger-color: #E74C3C;
$info-color: #3498DB;
$status-pending: #F39C12;
$status-in-transit: #3498DB;
$status-delivered: #27AE60;
$status-cancelled: #E74C3C;
```

### Typography
- **Font Family:** 'Roboto', sans-serif (Material Design standard)
- **Headings:** Bold, 24-32px
- **Body:** Regular, 14-16px
- **Small:** Regular, 12-13px

### Spacing Scale
- **xs:** 4px
- **sm:** 8px
- **md:** 16px
- **lg:** 24px
- **xl:** 32px

### Border Radius
- **sm:** 4px
- **md:** 8px
- **lg:** 16px

---

## 📝 Status Indicators

- 🔴 **CRITICAL** - Must be completed for MVP
- 🟡 **MEDIUM** - Important but not blocking
- 🟢 **LOW** - Nice-to-have features
- ✅ **DONE** - Completed
- ⚠️ **IN PROGRESS** - Currently being worked on
- ❌ **BLOCKED** - Waiting for dependencies

---

## 🔗 Related Documents

- [README.md](./README.md) - Project overview
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture (to be created)
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guidelines (to be created)

---

**Last Updated:** 2026-09-15  
**Updated By:** ElectroZombie  
**Branch:** main  

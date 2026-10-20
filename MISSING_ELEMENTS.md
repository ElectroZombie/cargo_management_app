# Missing Elements for Cargo Management App

**Analysis date:** 2026-10-20  
**Sources:** `TODO.md`, `ARCHITECTURE.md`, `CONTRIBUTING.md`, `AGENTS.md`, plus a full read of the `go/` backend and `frontend/src/` tree.  
**Branch:** main

> This document is a point-in-time gap analysis. It records what the documentation claims versus
> what the code actually delivers. `TODO.md` is the working backlog — keep that in sync when items
> here are closed.

---

## Summary

| Severity | Count | Description |
|----------|-------|-------------|
| 🔴 Critical | 15 | Claimed as done, but not implemented — or a documented architecture that the code does not follow |
| 🟠 High | 10 | Real product features and tooling that do not exist at all |
| 🟡 Medium | 12 | Documentation drift and build/structure defects |
| 🟢 Low | 6 | Polish, ergonomics, and distribution concerns |

**Headline:** the documentation has drifted ahead of the implementation. `TODO.md` marks Material
styling and chart support as complete in places, `MISSING_ELEMENTS.md` (previous revision) and
`IMPLEMENTED_SUMMARY.md` assert facts that are false, and `ARCHITECTURE.md` still describes a
package layout the repository no longer has.

**Progress:** the domain-error contract and the Wails binding layer have been brought in line
with `AGENTS.md` (see *Resolved since the previous revision*). The largest remaining structural
gap is the absent `go/services/` package.

---

## 🔴 Critical — documented as done, but not implemented

| # | Item | Claimed in | Reality |
|---|------|-----------|---------|
| 1 | Responsive line, pie, and bar chart components | `TODO.md` (chart section) | `chart.js` + `ng2-charts` installed and `provideCharts(withDefaultRegisterables())` registered in `AppModule`, but **zero chart components exist**. `dashboard.component.html` contains no chart reference. |
| 2 | Load, weight, status, and delivery-performance charts | `TODO.md` (chart section) | Not implemented, and **not implementable yet** — see #3. |
| 3 | Chart data endpoints / aggregation queries | implied by the chart goals | `go/database/dashboard.go` exposes only flat totals. There are no time-series queries for loads over time, weight/tons trends, status distribution, or on-time vs. late delivery. The backend must be built first. |
| 4 | Angular Material integration | `TODO.md` "Build reusable styled tables, forms, cards, and modals" ✅ | **Zero `mat-` usages in any template** across the whole app. No `MatTable`, `MatFormField`, `MatDialog`, `MatPaginator`, `MatSort`, `MatDatepicker`, `MatSelect`, `MatSnackBar`. Material is installed but entirely unused — pure bundle cost. |
| 5 | Page transitions and animations | `TODO.md` (unchecked item) | `BrowserAnimationsModule` is imported, but **no `@angular/animations` usage** and no route-transition code anywhere. |
| 6 | Hover effects and loading animations | `TODO.md` (unchecked item) | Not implemented. |
| 7 | Toast notifications | `TODO.md` (unchecked item) | ✅ Actually complete — `ToastComponent` exists, is declared in `SharedModule`, and is mounted in `app.component.html`. The checkbox was simply never ticked. |
| 8 | Angular unit tests | `AGENTS.md` §Testing, `CONTRIBUTING.md` §Testing | **Zero `.spec.ts` files** in the repository. `npm test` is declared in `package.json` but has nothing to execute. `frontend/src/app/core/services/driver.service.spec.ts` and the other specs shown in AGENTS.md do not exist. |
| 9 | Detail screens | `TODO.md` ("forms, filters, details, and status controls"), `AGENTS.md` directory listing | Only `*-list` and `*-form` exist per feature. No `driver-detail`, `loader-detail`, `vehicle-detail`, `well-detail`, or `load-detail`. |
| 10 | Load list filters | `TODO.md` | `load-list.component.ts` contains a single `{ key: 'status', header: 'Status', type: 'boolean' }` column definition and no filter logic at all. |
| 11 | Load status controls | `TODO.md` | Status is a read-only boolean column. No activate/deactivate action in the list or the form. |
| 12 | Reports feature | `TODO.md` ("reports and settings are placeholders") | `app-routing.module.ts` routes both `reports` and `settings` to `ComingSoonComponent`. |
| 13 | Settings feature | `TODO.md` | Same as #12 — placeholder only. |
| 14 | `go/services/` package | `AGENTS.md` §Package Structure | **The package does not exist.** All five entity services live in `go/database/`, so the layer AGENTS.md designates for business logic is absent and the documented boundary is not enforced in code. |
| 15 | `go/errors/` as the domain-error source | `AGENTS.md` §Error Handling | ✅ **Resolved.** `go/errors` now owns the sentinels and typed errors; `go/database/helpers.go` re-exports them instead of defining a second set. |
| 16 | `go/validators/` package | `AGENTS.md` §Package Structure | Only `driver_validator.go` exists and it is a **stub that returns `nil`**. All real validation is inlined in `go/database/helpers.go`. `loader`, `vehicle`, `well`, `load`, and `shared` validators are missing. |

---

## 🟠 High — features and tooling that do not exist

### Finance domain — schema only

`ARCHITECTURE.md` §Data Model lists five additional entities: `category`, `expense`, `income`,
`fuel_expense`, `stats`.

| Layer | State |
|-------|-------|
| Tables in `go/database/schema.go` | ✅ present |
| Read models in `go/models/finance.go` | ✅ present |
| Write DTOs in `go/dtos/finance_dto.go` | ✅ present |
| Validators in `go/database/helpers.go` | ✅ present (`validateCategory`, `validateExpense`, `validateIncome`, `validateFuelExpense`, `validateStats`) |
| **CRUD services** | ❌ **missing** |
| **Wails bindings** | ❌ **missing** |
| **Frontend screens** | ❌ **missing** |
| **Tests** | ❌ **missing** |

The previous `TODO.md` contained **zero** finance items, so the backlog did not even acknowledge the gap.

### Go test coverage

`TODO.md` claims "Added service and database tests" ✅, but coverage is still incomplete.

- ✅ `go/database/driver_service_test.go`
- ✅ `go/database/load_service_test.go`
- ✅ `go/database/loader_service_test.go`
- ✅ `go/database/vehicle_service_test.go`
- ✅ `go/database/well_service_test.go`
- ✅ `go/app/app_test.go` — binding-layer integration tests (envelope, error codes, pagination totals, filtering, dashboard stats)
- ❌ finance service tests
- ❌ validator unit tests (`go/validators/`)
- ❌ listing-layer tests (pagination, sorting, combined filtering)
- ❌ migration tests (apply and rollback)
- ❌ transaction-helper tests
- ❌ dashboard aggregation tests
- ❌ LCS search tests (`go/database/lcs.go`)

### Tooling referenced in the docs but absent

| Tool | Referenced in | Reality |
|------|---------------|---------|
| **CI pipeline** | `CONTRIBUTING.md` — "merge once approved and **CI is green**"; `AGENTS.md` §CI/CD | No `.github/` directory and no pipeline config of any kind. |
| **ESLint** | `frontend/package.json` → `"lint": "ng lint"` | No `@angular-eslint/*` dependency and no `.eslintrc*` file — `npm run lint` fails. |
| **Prettier** | `CONTRIBUTING.md` — "Format TypeScript with Prettier (2-space indent, single quotes)" | Not a dependency; no `.prettierrc`. |
| **EditorConfig** | — | Missing. |
| **LICENSE** | — | Missing. |
| **CODEOWNERS** | — | Missing. |

---

## 🟡 Medium — documentation drift and build defects

### Documentation contradictions

| Document | Stale or wrong claim |
|----------|---------------------|
| `ARCHITECTURE.md` §Directory Structure | Shows a root-level `database/` package; the `go/` tree is absent entirely. Reality: `go/{database,models,dtos,errors,validators,utils,app}`. Also states `main.go` holds the bindings — they live in `go/app/app.go`. |
| `ARCHITECTURE.md` §Data Model (Additional entities) | Presents the five finance entities as part of the live data model; no services, bindings, or UI exist for them. |
| `ARCHITECTURE.md` §Request Lifecycle | Shows `App.ListDrivers(ctx, params)` and `Database.ListDrivers(filters, opts)`; no binding accepts a `context.Context`. |
| `ARCHITECTURE.md` §IPC and Error Handling | Describes domain-error-to-message mapping; the code's `formatError` returns raw `err.Error()`. |
| `ARCHITECTURE.md` §Theming | States dark mode uses a `.dark-theme` class/attribute; the code sets a `data-theme` attribute on `<html>`. `TODO.md` says `[data-theme='dark']` — a **three-way inconsistency**. |
| `README.md` | States Go 1.18+ / Node 16+ (`CONTRIBUTING.md` says 1.21+ / 18+); the project tree omits `go/` and `database/`; lists a non-existent `frontend/styles.scss` (actual: `frontend/src/styles/index.scss`); written in Spanish while every other document is in English. |
| `MISSING_ELEMENTS.md` (previous revision) | Asserted that "the existing `database/` directory with `package database` remains functional" and that the `cargo_management_app/database` import path still works — both false after the move to `go/database/`. Claimed "MatChartsModule added to AppModule" — the code uses `provideCharts()`. Claimed `go/errors` was created *and* used — it is imported by nothing. |
| `IMPLEMENTED_SUMMARY.md` | Same `MatChartsModule` error and the same over-claiming about `go/errors`. |
| Document metadata | "Last Updated" dates differ across files: `AGENTS.md` 2026-10-02, `ARCHITECTURE.md` 2026-10-06, `TODO.md` 2026-10-20. |

### Build and structure defects

1. **Duplicate Material theme emission.** Two files each generate a Material theme — `frontend/src/styles/theme.scss` and `frontend/src/app/shared/styles/theme.scss`. The Angular build emits *"the same color / density / typography styles are generated multiple times"* warnings.
2. **Bundle budget exceeded.** The initial bundle is 825.39 kB against a 500 kB warning / 1 MB error budget in `frontend/angular.json`. Material and Chart.js are both in the bundle while being largely unused.
3. **Leftover empty `database/` directory** at the repository root (untracked, but present on disk and visible in the tree).
4. **Empty `frontend/src/app/shared/pipes/`** — AGENTS.md lists pipes as a shared concern; the directory has no files.
5. **Dead Go packages.** `go/errors` and `go/utils` are imported by nothing; `go/validators` contains a single no-op stub.
6. **No shared Go test fixture.** The test-DB setup helper is not factored out for reuse by the missing test files.
7. **Monolithic binding file.** All 388 lines of Wails bindings live in `go/app/app.go` instead of the per-entity files AGENTS.md prescribes.

---

## 🟢 Low — polish and ergonomics

1. `confirm-dialog` component exists but delete flows use the native `confirm()` instead.
2. `environment.ts` / `environment.prod.ts` exist, but nothing branches on environment (expected for a Wails build, but worth confirming intent).
3. No release or binary-distribution step; `wails build` output is not published anywhere.
4. `wails.json` has no project-level dev/build script polish (e.g. no `npm run` wrapper scripts).
5. README has no feature tour, screenshots, or data-model overview.
6. `frontend/src/app/core/core.module.ts` guards against double import, but no service is actually provided through it — all services rely on `providedIn: 'root'`.

---

## Non-Goals (intentionally excluded — do not implement)

- ❌ Authentication, login, registration, sessions, or tokens
- ❌ User accounts, roles, permissions, AuthGuards, or RoleGuards
- ❌ Multi-tenancy or remote/server deployments
- ❌ A public HTTP/REST API surface

These are confirmed exclusions in `ARCHITECTURE.md` §Non-Goals, `AGENTS.md` §General Principles, and
`CONTRIBUTING.md`. They are **not** missing elements.

---

## Recommended order of work

1. **Backend chart data endpoints** — unblocks the entire chart backlog (#1–#3).
2. **Material integration** — either adopt it across the shared components and forms, or remove the
   dependency and its theme to stop the bundle-budget and duplicate-theme warnings (#4, build defects 1–2).
3. **Go package consolidation** — create `go/services/`, replace the `go/validators` stub (#14, #16).
4. **Finance domain** — services, bindings, and screens, then replace the Reports placeholder.
5. **Feature completeness** — detail screens, list filters, load status controls (#9–#11).
6. **Test coverage** — the Angular specs (#8) and the remaining Go gaps (validators, listing,
   migrations, transactions, finance).
7. **Tooling** — CI, ESLint, Prettier, EditorConfig, LICENSE.
8. **Documentation reconciliation** — bring `ARCHITECTURE.md` and `README.md` back in line with the code.

---

## ✅ Resolved since the previous revision

The following were implemented in this pass and are no longer gaps:

| Item | Change |
|------|--------|
| **Single source of domain errors** (#15) | `go/errors` now owns the sentinels (`ErrValidation`, `ErrNotFound`, `ErrConflict`, `ErrInternal`) and the typed errors (`ValidationError`, `NotFoundError`, `ConflictError`, `InternalError`), each with an `Unwrap` method so `errors.Is` and `errors.As` work through wrapping. `go/database/helpers.go` re-exports the sentinels and aliases `ValidationError` instead of defining a competing set. |
| **Domain-aware error formatting** | `go/app/errors.go` classifies failures into `validation`, `not_found`, `conflict`, `internal`, and `unknown`; internal failures are logged via `slog` and reduced to a generic message so raw SQLite internals never reach the UI. The `Response` envelope now carries a machine-readable `code` next to `error`. |
| **Per-entity binding files** | `go/app/app.go` was split into `app.go` (struct, lifecycle, envelope), `errors.go` (error handling), and `driver/loader/vehicle/well/load/dashboard_bindings.go`. |
| **Missing entity service tests** | Added `loader_service_test.go`, `vehicle_service_test.go`, and `well_service_test.go`, asserting the domain error types rather than only "an error occurred". |
| **Binding-layer integration tests** | Added `go/app/app_test.go` covering the envelope, success/failure paths, error codes, pagination totals, filtering, dashboard stats, and the guarantee that internal error detail is hidden. |
| **Frontend error surfacing** | `ApiResponse`/`ListResponse` gained an optional `code`, and `CrudService` exposes a `lastError()` observable so forms can show the backend's field-level validation message. |

### Known framework constraint

The `context.Context`-first binding signature that `AGENTS.md` and `ARCHITECTURE.md` prescribe is
**not achievable on Wails v2.8.1**: `BoundMethod.ParseArgs` requires an exact argument-count match
and JSON-decodes every argument, so a `context.Context` parameter cannot be injected and would
break the frontend call. The application context is captured in `App.Startup` and held on the
`App` struct. This should be revisited if the project upgrades to Wails v3.

---

## Related documents

- [TODO.md](./TODO.md) — Working backlog
- [ARCHITECTURE.md](./ARCHITECTURE.md) — System architecture
- [CONTRIBUTING.md](./CONTRIBUTING.md) — Contribution workflow
- [AGENTS.md](./AGENTS.md) — Development best practices
- [README.md](./README.md) — Project overview

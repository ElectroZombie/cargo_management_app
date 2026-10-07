# 🤝 CONTRIBUTING.md - Cargo Management App

**Purpose:** Guidelines for contributing code, documentation, and reviews to the Cargo Management App  
**Last Updated:** 2026-10-06  
**Branch:** main

---

## 📖 Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Ways to Contribute](#ways-to-contribute)
3. [Prerequisites](#prerequisites)
4. [Getting Started](#getting-started)
5. [Project Structure](#project-structure)
6. [Development Workflow](#development-workflow)
7. [Branching Strategy](#branching-strategy)
8. [Commit Conventions](#commit-conventions)
9. [Coding Standards](#coding-standards)
10. [Testing](#testing)
11. [Documentation](#documentation)
12. [Pull Request Process](#pull-request-process)
13. [Definition of Done](#definition-of-done)

---

## Code of Conduct

- Be respectful and constructive in all communication.
- Assume good intent and focus feedback on the code, not the person.
- Keep discussions on-topic and evidence-based.
- Report unacceptable behavior to the maintainer (`ElectroZombie`).

---

## Ways to Contribute

- **Bug fixes** — reproduce, isolate, fix, and add a regression test.
- **Features** — pick an unchecked item from [TODO.md](./TODO.md) and open a PR.
- **Tests** — increase coverage for services, validators, and components.
- **Documentation** — improve README, ARCHITECTURE, or inline docs.
- **Refactoring** — reduce duplication and improve clarity without changing behavior.

Before starting significant work, open an issue or comment on an existing one so the direction can
be agreed on early.

---

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.21+ | https://go.dev/dl/ |
| Node.js | 18+ | https://nodejs.org/ |
| npm | 9+ | bundled with Node |
| Wails CLI | v2.8+ | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |

Platform build dependencies for Wails (WebKit2GTK on Linux, WebView2 on Windows) must also be
installed. See the [Wails docs](https://wails.io/docs/gettingstarted/installation) for details.

---

## Getting Started

```bash
# 1. Clone the repository
git clone <repository-url>
cd "Cargo Management App"

# 2. Install Go dependencies
go mod download
go mod tidy

# 3. Install frontend dependencies
cd frontend && npm install && cd ..

# 4. Verify the toolchain
wails doctor

# 5. Run in development mode (hot reload)
wails dev
```

Build a production binary with:

```bash
wails build
```

The database is created automatically at `~/.cargo_management_app/cargo_management.db` on first
launch.

---

## Project Structure

| Path | Responsibility |
|------|----------------|
| `main.go` | Wails bootstrap, startup/shutdown, exported bindings |
| `database/` | Go persistence, services, validation, migrations |
| `frontend/src/app/core/` | Models, Wails bridge, entity services |
| `frontend/src/app/shared/` | Reusable components, pipes, styles |
| `frontend/src/app/layout/` | Application shell (navbar, sidebar, breadcrumb, footer) |
| `frontend/src/app/features/` | Feature screens |

See [ARCHITECTURE.md](./ARCHITECTURE.md) for a detailed breakdown.

---

## Development Workflow

1. **Sync** with `main`: `git checkout main && git pull`.
2. **Create a branch** using the naming convention below.
3. **Implement** the change following the coding standards.
4. **Test** locally (`go test ./...`, `npm test`).
5. **Commit** using conventional commit messages.
6. **Push** and open a **pull request** against `main`.
7. **Address review** feedback with follow-up commits.
8. **Merge** once approved and CI is green (squash and merge).

---

## Branching Strategy

```
<type>/<short-description>
```

| Type | Use for | Example |
|------|---------|---------|
| `feature/` | New functionality | `feature/load-management` |
| `bugfix/` | Non-critical fixes | `bugfix/license-validation` |
| `hotfix/` | Urgent production fixes | `hotfix/database-lock` |
| `docs/` | Documentation only | `docs/update-architecture` |
| `refactor/` | Behavior-preserving changes | `refactor/service-helpers` |
| `test/` | Test-only changes | `test/driver-service` |

Keep branches focused and short-lived. Rebase on `main` before opening a PR.

---

## Commit Conventions

Use the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>
```

- **Type:** `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `style`, `perf`, `build`.
- **Scope:** the affected area (`driver-service`, `frontend`, `api`, `docs`).
- **Subject:** imperative, present tense, no trailing period, ≤ 72 characters.

Examples:

```
feat(load-service): add pagination and status filtering
fix(vehicle-service): reject invalid VIN before insert
docs(architecture): document request lifecycle
test(driver-service): cover duplicate license numbers
refactor(helpers): centralize SQLite error translation
```

**Rules**

- One logical change per commit.
- Never commit secrets, credentials, or the SQLite database file.
- Do not commit generated `frontend/dist` or `node_modules`.
- Reference issue numbers in the body when applicable (`Closes #12`).

---

## Coding Standards

Follow [AGENTS.md](./AGENTS.md) for the full standards. Key points:

### Go

- Accept `context.Context` as the first argument in service methods.
- Never expose raw SQLite errors; use domain error types (`ValidationError`, `NotFoundError`,
  `ConflictError`, `InternalError`).
- Keep SQL as named constants near their usage.
- Validate DTOs before writes and check foreign-key references explicitly.
- Keep functions small and single-purpose.

### TypeScript / Angular

- Use `ChangeDetectionStrategy.OnPush` in every component.
- Prefer the `async` pipe over manual subscriptions; unsubscribe when you must subscribe.
- Use Reactive Forms with typed `FormGroup`s.
- Model fields mirror Go models in camelCase.
- Container components inject services; presentational components use `@Input()` / `@Output()`.
- Organize imports: Angular core → RxJS → third-party → local.

### Naming

| Entity | Convention |
|--------|-----------|
| Go files/packages | `snake_case` |
| Go exported funcs | `PascalCase` |
| Go constants | `UPPER_SNAKE_CASE` |
| TS classes/interfaces | `PascalCase` |
| TS methods/properties | `camelCase` |
| TS observables | `$` suffix (`items$`) |
| TS DTOs | `DTO` suffix |
| Components/services files | `kebab-case.component.ts` / `.service.ts` |

### Style

- Format Go with `gofmt` / `go vet`.
- Format TypeScript with Prettier (2-space indent, single quotes).
- Do not add code comments unless they explain *why*, not *what*.
- No emojis in source files.

---

## Testing

### Go

- Place tests next to the source as `*_test.go`.
- Use table-driven tests for multiple scenarios.
- Cover happy paths, validation failures, not-found, and conflict cases.
- Clean up databases and use temporary files per test.

```bash
go test ./...
go vet ./...
```

### Angular

- Use `TestBed` with mocked collaborators (`jasmine.SpyObj`).
- Test services for response handling and state updates.
- Test presentational components for input/output behavior.

```bash
cd frontend
npm test
```

A PR should not reduce overall test coverage. Bug fixes should include a regression test.

---

## Documentation

- Update documentation in the **same PR** as the code change.
- Keep [README.md](./README.md) for setup and usage.
- Keep [ARCHITECTURE.md](./ARCHITECTURE.md) for system design.
- Keep [TODO.md](./TODO.md) status current when completing or adding work.
- Document public Go functions and TypeScript APIs with comments/JSDoc.

---

## Pull Request Process

1. Fill out the PR description: what changed, why, and how it was tested.
2. Link the related issue(s).
3. Ensure the branch is up to date with `main`.
4. Confirm the build and tests pass locally.
5. Request review from the maintainer.
6. Address feedback in new commits (do not force-push over review history unless asked).

### Review checklist

- [ ] Change is focused and scoped to one concern.
- [ ] Code follows AGENTS.md and this guide.
- [ ] Validation and error handling are correct.
- [ ] Tests cover the change and pass.
- [ ] Documentation and TODO are updated.
- [ ] No secrets, generated files, or stray debug code committed.
- [ ] No regressions in existing behavior.

---

## Definition of Done

A contribution is complete when:

- [ ] The feature/fix works as described.
- [ ] Go tests and Angular tests pass.
- [ ] Validation, error handling, and data integrity are handled.
- [ ] Responsive behavior is verified for UI changes.
- [ ] Documentation is updated.
- [ ] TODO.md reflects the new status.
- [ ] The change is reviewed and merged.

---

## Related Documents

- [README.md](./README.md) — Project overview and setup
- [ARCHITECTURE.md](./ARCHITECTURE.md) — System architecture
- [AGENTS.md](./AGENTS.md) — Development best practices
- [TODO.md](./TODO.md) — Development backlog

---

**Author:** ElectroZombie  
**Branch:** main

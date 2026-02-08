- To regenerate the JavaScript SDK, run `./packages/sdk/js/script/build.ts`.
- ALWAYS USE PARALLEL TOOLS WHEN APPLICABLE.
- The default branch in this repo is `dev`.
- Local `main` ref may not exist; use `dev` or `origin/dev` for diffs.
- Prefer automation: execute requested actions without confirmation unless blocked by missing info or safety/irreversibility.

## Style Guide

### General Principles

- Keep things in one function unless composable or reusable
- Avoid `try`/`catch` where possible
- Avoid using the `any` type
- Prefer single word variable names where possible
- Use Bun APIs when possible, like `Bun.file()`
- Rely on type inference when possible; avoid explicit type annotations or interfaces unless necessary for exports or clarity
- Prefer functional array methods (flatMap, filter, map) over for loops; use type guards on filter to maintain type inference downstream

### Naming

Prefer single word names for variables and functions. Only use multiple words if necessary.

```ts
// Good
const foo = 1
function journal(dir: string) {}

// Bad
const fooBar = 1
function prepareJournal(dir: string) {}
```

Reduce total variable count by inlining when a value is only used once.

```ts
// Good
const journal = await Bun.file(path.join(dir, "journal.json")).json()

// Bad
const journalPath = path.join(dir, "journal.json")
const journal = await Bun.file(journalPath).json()
```

### Destructuring

Avoid unnecessary destructuring. Use dot notation to preserve context.

```ts
// Good
obj.a
obj.b

// Bad
const { a, b } = obj
```

### Variables

Prefer `const` over `let`. Use ternaries or early returns instead of reassignment.

```ts
// Good
const foo = condition ? 1 : 2

// Bad
let foo
if (condition) foo = 1
else foo = 2
```

### Control Flow

Avoid `else` statements. Prefer early returns.

```ts
// Good
function foo() {
  if (condition) return 1
  return 2
}

// Bad
function foo() {
  if (condition) return 1
  else return 2
}
```

### Schema Definitions (Drizzle)

Use snake_case for field names so column names don't need to be redefined as strings.

```ts
// Good
const table = sqliteTable("session", {
  id: text().primaryKey(),
  project_id: text().notNull(),
  created_at: integer().notNull(),
})

// Bad
const table = sqliteTable("session", {
  id: text("id").primaryKey(),
  projectID: text("project_id").notNull(),
  createdAt: integer("created_at").notNull(),
})
```

## Testing

- Avoid mocks as much as possible
- Test actual implementation, do not duplicate logic into tests

## Golang Refactor Project (Multi-Agent Collaboration)

### Project Overview
Complete rewrite of OpenCode from TypeScript/JavaScript to Pure Golang with TUI and GUI interfaces.

### Active Branches
- **`golang-refactor`** (Agent 1): TUI development, backend services, core infrastructure
- **`golang-refactor-v2`** (Agent 2): Architecture design, performance optimization, GUI planning, documentation

### Collaboration Rules

#### Synchronization Protocol
1. **Shared Code** lives in `/go/pkg/{types,api,core}` - controlled by Agent 1
2. **Agent 1 Code** in `/go/cmd/tui`, `/go/pkg/tui`, `/go/test`
3. **Agent 2 Code** in `/go/doc`, `/go/test/benchmarks`, `/go/config`
4. **Merge Conflict** → Agent 1 has final say on shared packages

#### Git Workflow
```
Developer in golang-refactor:
  1. Make changes on golang-refactor branch
  2. Run: go test -v ./... && go build ./cmd/{tui,backend}
  3. Commit with: git commit -m "feat: xxx\n\nCo-authored-by: Agent 2"
  4. Periodically: git fetch origin golang-refactor-v2 && git merge (if needed)

Developer in golang-refactor-v2:
  1. Work on architecture, perf, docs
  2. Fetch Agent 1's changes: git fetch origin golang-refactor
  3. Create PR to golang-refactor when ready (don't merge directly)
  4. Wait for Agent 1 review before merging
```

#### Task Allocation (Feb 2026)
- **Agent 1**: Phase 1 (✅ done) → Phase 2 TUI completion → Phase 3 backend features
- **Agent 2**: Phase 1 architecture validation → Phase 2 performance → Phase 3 GUI groundwork

#### Go Style Rules (Specific to this project)
- Use `casing` following Go conventions (PascalCase for exported, camelCase for package-level private)
- Minimize struct tags - use only when marshaling (json, db tags)
- Avoid package `init()` functions; prefer explicit initialization
- Prefer `context.Context` as first parameter for long-running operations
- Use `error` return as last value; panic only for programming errors
- No global state except in `func init()` for package-level singletons

#### Testing Requirements
- **Minimum coverage**: 80% overall, 95% for core packages
- **Core packages** (`pkg/core`, `pkg/api`): zero-mock tests on real file ops
- **Unit tests** for all public functions in `pkg/*` packages
- **Integration tests** in `test/` folder for TUI+backend interaction
- Run tests before every commit: `go test -v -cover ./...`

### Current Status (Feb 8, 2026)
- ✅ Project scaffold complete
- ✅ Backend services (FileService, ProjectService, API)
- ✅ TUI skeleton (basic layout, file tree, editor, status bar)
- ✅ API client and HTTP server
- ✅ Unit tests for core (4 tests passing, all green)
- 🔄 Phase 2: Expanding TUI features (syntax highlighting, key bindings)
- See full plan in `/go/COLLABORATION.md`

### Key Files
- `/go/COLLABORATION.md` - Full agent cooperation guide
- `/go/go.mod` - Golang module dependencies
- `/go/README.md` - Project-level documentation
- `/go/Makefile` - Build and run targets

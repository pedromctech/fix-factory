# AGENTS.md

fix-factory is a Go app for fixing security vulnerabilities. The user adds GitHub repos, the app
lists their vulnerabilities, and when the user picks one it runs the factory (plan → build → check
→ review) to open a PR that fixes it.

## Status

This is only a skeleton. Right now the code parses repo references (`internal/repo`) and has a small
CLI (`fix-factory add <owner/name>`). Planned packages, each added by its own issue:
- `repo`: parse and validate `owner/name` and GitHub URLs (exists)
- `store`: save the repos the user added
- `vulns`: get and normalize vulnerabilities for a repo (e.g. Dependabot/code-scanning alerts via the GitHub API)
- `run`: start a factory run for one vulnerability and follow it until a PR exists

## Commands
- Test: `go test ./...` (the TEST_CMD in factory.conf; it must pass before any PR)
- Lint: `go vet ./...` (the LINT_CMD); code must also be `gofmt`-clean
- Setup on a fresh clone: `go build ./...`
- Run: `go run ./cmd/fix-factory add octo/hello-world`

## Layout
- `go.mod`: module `github.com/pedromctech/fix-factory`.
- `internal/<package>/<package>.go`: one package per planned module (see Status). Logic lives here.
- `internal/<package>/<package>_test.go`: tests in the external `<package>_test` package, so they
  only use the public API.
- `cmd/fix-factory/main.go`: a thin CLI. It parses args, calls the packages and maps errors to exit
  codes (0 ok, 1 error, 2 usage). No logic.

## Conventions
- Go 1.27, standard toolchain. Code must be `gofmt`-clean and pass `go vet` with no warnings.
- Errors: return `error`, never panic on user input. Each package exports sentinel errors
  (`ErrBadFormat`) and, when the input matters, an error type that wraps them (`*ParseError` with
  `Unwrap`), so callers use `errors.Is`/`errors.As`. Wrap with `fmt.Errorf("...: %w", err)`.
  Don't ignore returned errors.
- Parse at the edge: turn strings into typed values (`repo.Ref`, later `Vulnerability`, `Severity`)
  through `Parse`/constructor functions that validate, and pass typed values around after that.
- Security: this app handles untrusted input (repo names, API responses, advisory text). Validate
  input before using it in URLs, paths or commands. Never build a shell command by joining strings:
  use `exec.Command(name, args...)` with separate args, never `sh -c`. Never log or print tokens.
  Read the GitHub token from the environment (`GITHUB_TOKEN`), never from code or files in the repo.
- Dependencies: `go.mod`/`go.sum` are outside SCOPE. If a change needs a new module, say so in the
  plan; a human adds it. Prefer the standard library (`net/http`, `encoding/json`).
- Tests: a good test checks behavior through the public API. Name it for what it checks
  (`TestRejectsPathTraversal`); prefer table-driven cases. Cover the happy path and at least one
  rejected or edge input. Tests run offline: no network and no real GitHub calls (use fixtures,
  `httptest.Server` or an interface you can fake).

## Factory rules
- Phase state lives in `.factory/` as markdown files. Never put plan or review text in code comments.
- Only the paths in SCOPE (factory.conf: `cmd/ internal/`) may change. Do not touch factory.sh,
  doctor.sh, lib/, gates/, .claude/, dashboard/, factory.conf, go.mod, go.sum, .gitignore or this file.

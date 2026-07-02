# Apartment Building Maintenance CLI

A command-line tool that gives a building manager a fast, structured way to
register, track, and resolve apartment maintenance issues directly from the
terminal.

## Business Value

Building managers often track maintenance through paper logs or informal
messaging, which leads to lost records, missed follow-ups, and no audit trail.
This tool replaces that with a single source of truth:

- **Structured records** — every issue has a unit, priority, status, and timestamps.
- **A clear workflow** — issues move through `open → in-progress → resolved → closed` with validated transitions.
- **Automatic stakeholder notifications** — an email is sent when an issue is created or resolved (via SendGrid).
- **Local, durable state** — data is persisted in SQLite so nothing is lost between sessions.

## Architecture

Layered design with dependency injection, so the business logic is isolated and
testable without a real database or network:

```
cmd/                     → Cobra commands (thin CLI layer)
internal/
  model/                 → Issue struct, status workflow, domain errors
  repository/            → GORM/SQLite data layer (implements service.Repository)
  service/               → business logic (IssueService) + SendGrid notifier
  database/              → GORM connection + AutoMigrate
```

The service depends on `Repository` and `Notifier` **interfaces**, not concrete
types — the real GORM repository and SendGrid notifier are injected in `cmd/`,
while tests inject in-memory mocks.

## Prerequisites

- **Docker** and **Docker Compose** (the only requirement to run the system).
- For local (non-Docker) development: **Go 1.26+**.

## Quick Start (Docker)

Bringing up the entire system requires a single command:

```bash
docker compose up --build
```

This builds the image and starts a container with the CLI installed and a
persistent SQLite volume mounted at `/data`. The container stays alive so you
can run commands against it:

```bash
docker compose exec app maintenance issue create \
  --title "Broken boiler" --unit 4B --description "No heat on 4th floor" --priority high

docker compose exec app maintenance issue list
docker compose exec app maintenance resolve 1
```

The database schema is created and migrated automatically (GORM AutoMigrate) on
the first command.

## Configuration

Copy the template and fill in real values. `.env` is git-ignored, so secrets
stay local.

```bash
cp .env.example .env
```

### Environment variables (`.env.example`)

| Variable            | Required | Description                                                        |
|---------------------|----------|--------------------------------------------------------------------|
| `DB_PATH`           | Yes      | Path to the SQLite file (e.g. `/data/maintenance.db` in Docker).   |
| `SENDGRID_API_KEY`  | No       | SendGrid API key. If blank, email notifications are skipped safely.|
| `EMAIL_FROM`        | No       | Verified SendGrid sender address.                                  |
| `MANAGER_EMAIL`     | No       | Recipient of create/resolve notifications.                         |

Email is best-effort: if the key is missing the app logs a notice and continues;
if SendGrid returns an error it is logged to stderr and never crashes the CLI.

## Usage

Canonical form is `maintenance <resource> <action> [flags]`:

```bash
# Create a new issue
maintenance issue create --title "Broken boiler" --unit 4B --description "No heat" --priority high

# List all issues
maintenance issue list

# View one issue
maintenance issue view --id 1

# Update status (validated transitions)
maintenance issue update --id 1 --status in-progress

# Resolve an issue (sends email notification)
maintenance issue resolve --id 1

# Delete an issue
maintenance issue delete --id 1
```

### Shortcuts

Top-level shortcuts take the ID as a positional argument and share the same
logic (including notifications):

```bash
maintenance resolve 1   # == maintenance issue resolve --id 1
maintenance delete 1    # == maintenance issue delete --id 1
```

### Issue workflow

`open → in-progress → resolved → closed`

| From          | Allowed targets                     |
|---------------|-------------------------------------|
| `open`        | `in-progress`, `resolved`, `closed` |
| `in-progress` | `open`, `resolved`, `closed`        |
| `resolved`    | `in-progress` (reopen), `closed`    |
| `closed`      | _(terminal — no transitions)_       |

Invalid transitions and unknown issue IDs return a clear error, and the CLI
exits with a non-zero POSIX status code.

## Running Locally (without Docker)

```bash
cp .env.example .env         # set DB_PATH=maintenance.db for a local file
go run ./cmd issue create --title "Broken boiler" --unit 4B --description "No heat" --priority high
go run ./cmd issue list
```

## Testing

The business-logic layer is covered by unit tests that use mocked dependencies —
no real database or SendGrid calls (the SendGrid path is exercised with an
`httptest` server):

```bash
go test ./...
go test ./... -race    # as run in CI
```

## Continuous Integration

GitHub Actions runs on every pull request (`.github/workflows/ci.yml`) and must
be green before review:

1. **Lint** — `gofmt`, `go vet`, and `golangci-lint`.
2. **Test** — the full suite with the race detector.
3. **Docker build** — builds the image to prove the infrastructure is valid.

### Secrets management

`SENDGRID_API_KEY` is injected into the test job from **GitHub Secrets**
(`${{ secrets.SENDGRID_API_KEY }}`), never committed to the repo. GitHub
automatically masks secret values in logs. To configure it:

**Repository → Settings → Secrets and variables → Actions → New repository secret**, name it `SENDGRID_API_KEY`.

## Tech Stack

| Concern        | Choice                        |
|----------------|-------------------------------|
| Language       | Go                            |
| CLI Framework  | Cobra                         |
| Database       | SQLite (`glebarez/sqlite`)    |
| ORM            | GORM (+ AutoMigrate)          |
| Config         | godotenv (`.env`)             |
| Email          | SendGrid API                  |
| Container      | Docker + docker-compose       |
| CI             | GitHub Actions                |

See [RFC.md](RFC.md) for the full design rationale.

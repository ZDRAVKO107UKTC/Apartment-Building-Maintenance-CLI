# RFC: Apartment Building Maintenance CLI

**Status:** Proposed  
**Author:** ZDRAVKO107UKTC  
**Date:** 2026-05-26

---

## 1. Problem Statement

Building managers currently track maintenance issues through paper logs or
informal messaging, leading to lost records, missed follow-ups, and no
audit trail. This CLI application gives a building manager a fast,
structured way to register, track, and resolve maintenance issues directly
from the terminal. All state is persisted locally and stakeholders are
notified automatically via email when an issue is created or resolved.

---

## 2. Interface

All commands follow the pattern: `maintenance <resource> <action> [flags]`

```bash
# Create a new maintenance issue
maintenance issue create --unit 4B --description "Broken boiler" --priority high

# List all active issues
maintenance issue list

# View a single issue
maintenance issue view --id 1

# Update issue status
maintenance issue update --id 1 --status in-progress

# Resolve an issue
maintenance issue resolve --id 1

# Delete an issue
maintenance issue delete --id 1
```

**Issue States:** `open` → `in-progress` → `resolved`

---

## 3. Tech Stack

| Concern        | Choice                        |
|----------------|-------------------------------|
| Language       | Go 1.22                       |
| CLI Framework  | Cobra                         |
| Database       | SQLite                        |
| ORM / DB Layer | GORM                          |
| Migrations     | Raw SQL migration scripts     |
| Config         | Godotenv (.env files)         |
| Email          | SendGrid API                  |
| Testing        | Go testing package + testify  |
| Container      | Docker + docker-compose       |
| CI             | GitHub Actions                |

---

## 4. External Dependency

**SendGrid Email API**  
When a maintenance issue is **created** or **resolved**, the system will
send an automated email notification to a configured building manager
address. The HTTP client will enforce a timeout and log failures to stderr
without crashing the application.

Required environment variable: `SENDGRID_API_KEY`

---

## 5. Architecture
/cmd
main.go                  → entry point
/internal
/repository
issue_repository.go    → all SQLite queries (Data Layer)
/service
issue_service.go       → business logic (Business Logic Layer)
email_service.go       → SendGrid integration (Business Logic Layer)
/model
issue.go               → Issue struct and types
/db
migrations/
001_create_issues.sql  → schema migration
.env.example               → committed config template
Dockerfile
docker-compose.yml

---

## 6. MVP Scope

In scope:
- Full CRUD for maintenance issues
- Issue status workflow (open → in-progress → resolved)
- Email notification on create and resolve via SendGrid
- Dockerised local environment
- GitHub Actions CI (lint, test, build)

Out of scope for MVP:
- User authentication
- Multiple buildings
- Web UI or REST API
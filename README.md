# Dividr

**HMRC MTD Income Tax Bridge for Sole Traders & Landlords.**

Dividr simplifies the transition to Making Tax Digital (MTD) for Income Tax Self Assessment (ITSA). It acts as a smart bridge between your spreadsheets and HMRC, handling the complexity of quarterly updates without forcing you to abandon your existing workflow.

> 📖 **Want to know what Dividr can do?** See the comprehensive [Capabilities Guide](docs/CAPABILITIES.md) for detailed feature documentation.


## 🏗️ Architecture: The Session-Based Bridge
Dividr utilizes a **Session-Based Persistence** model. This hybrid approach balances user convenience with strict data minimization:

* **Resumable Workflows:** We temporarily persist your session data (parsed CSVs and mapping choices) to our secure database. This ensures your work is saved even if you close your browser or need to finish the return later.
* **Transient Lifecycle:** Your raw transaction data exists in our system **only** for the duration of the submission process.
* **Submission & Purge:** Once a return is successfully filed with HMRC and the receipt is generated, the raw line-item data is purged. We retain only the **Submission Receipt** (Correlation IDs) and the high-level **Dashboard Status** to prove compliance.


## ✨ Key Features

* **The 8-Slot Dashboard:** A persistent visual grid tracking the 4 quarterly updates for both Sole Trade and UK Property streams.
* **Stream Separator:** Intelligently maps and divides "Hybrid" income records into the correct legal buckets.
* **Digital Link Compliance:** Automated mapping from spreadsheet to HMRC categories, preserving the "Digital Link" required by MTD legislation.
* **Receipt Vault:** Generates immutable PDF/JSON receipts with HMRC Correlation IDs for record preservation.


## 🛠️ Tech Stack
Dividr is built for speed, security, and type-safety:

- **Backend:** [Go](https://go.dev/) (Golang)
- **UI Architecture:** [Templ](https://templ.guide/) + [HTMX](https://htmx.org/) + [Tailwind CSS](https://tailwindcss.com/)
- **Database (Metadata only):** [Postgres](https://www.postgresql.org/) + [sqlc](https://sqlc.dev/) for type-safe queries.
- **API:** REST /api/v1 (Versioned)
- **Communication:** HMRC OAuth2 & Fraud Header Engine.

---

## 🏗️ Project Structure
```text
.
├── cmd/dividr/                 # Main application entry point
├── internal/
│   ├── app/                    # Application lifecycle & graceful shutdown
│   ├── config/                 # Configuration management
│   ├── database/               # Database layer (sqlc generated code)
│   │   ├── migrations/         # SQL migration files
│   │   └── queries/            # SQL query definitions for sqlc
│   ├── domain/                 # Core business logic & models
│   ├── hmrc/                   # HMRC API client & Fraud Header engine
│   ├── http/                   # HTTP handlers, middleware, routing
│   └── ui/                     # Templ components
├── docker/                     # Docker Compose & secrets
├── web/                        # Static assets (Tailwind, CSS, JS)
└── Makefile                    # Dev automation (build, test, migrate)
```

---

## 🗄️ Database Configuration

**PostgreSQL 17** via Docker Compose with Adminer UI for local development.

### Schema Overview
Dividr uses **session-based persistence** with automatic purging post-submission:

#### Core Tables
- **users** – Authentication (email + bcrypt hash)
- **sessions** – Session tokens with expiry
- **oauth_tokens** – HMRC OAuth credentials (application-layer encrypted)
- **audit_events** – Security and compliance logging

#### Workflow Tables (Transient Data)
- **transaction_import_batches** – Upload metadata and batch tracking
- **transactions** – Staging workspace for parsed CSV rows (purged after submission)
- **mapping_profiles** – Reusable CSV column-to-HMRC category mappings

#### Compliance Vault
- **receipts** – Immutable submission proofs with HMRC correlation IDs and payload hashes

**Data Lifecycle:** Raw transaction data in the `transactions` table exists only during active sessions. Once a return is filed with HMRC, the line-item data is purged, retaining only aggregated totals and receipt metadata for compliance.

### Local Setup
```bash
# Start database + Adminer UI
docker compose -f docker/docker-compose.yaml up -d

# Run migrations
make migrate-up

# Generate type-safe queries (sqlc)
make sqlc
```

**Adminer UI:** http://localhost:8081

### Migration Management
Powered by [golang-migrate](https://github.com/golang-migrate/migrate):
```bash
make migrate-create name=add_new_table  # Create new migration
make migrate-up                          # Apply all migrations
make migrate-down                        # Rollback last migration
```

All migrations live in `internal/database/migrations/` as sequential `.up.sql` and `.down.sql` files.

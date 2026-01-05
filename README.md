# Dividr | The 2026 MTD ITSA Bridge

![Dividr Banner](https://via.placeholder.com/800x200/0F172A/0EA5E9?text=Dividr+|+Separate+your+income.+Simplify+your+reporting.)

**Dividr** is a high-performance, stateless bridge for UK MTD ITSA compliance. Built with Go and HTMX, it acts as a secure conduit between raw financial records and HMRC, enabling sole traders and property investors to fulfill their 2026 mandates without ever storing private transaction data.

---

## 🛡️ The Stateless Guarantee
Unlike traditional accounting software, **Dividr is not a ledger.** We believe your financial privacy is paramount. Dividr is built on a "Stateless Architecture":
- **Zero-Persistence:** We never store uploaded CSVs, XLSX files, or raw transaction rows.
- **In-Memory Processing:** Data is parsed, transformed into HMRC-compliant JSON, and transmitted in-session.
- **Immediate Purge:** Once the submission receipt is generated, the session data is destroyed.

---

## ✨ Key Features
- **The 8-Slot Dashboard:** A clear, visual grid tracking the 4 quarterly updates for both Sole Trade and UK Property streams.
- **Stream Separator:** Intelligently maps and divides "Hybrid" income records into the correct legal buckets.
- **Digital Link Compliance:** Automated mapping from spreadsheet to HMRC categories—no manual re-typing, keeping you fully compliant with HMRC’s digital link requirements.
- **Receipt Vault:** Generates immutable PDF/JSON receipts with HMRC Correlation IDs for record preservation.

---

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
├── cmd/dividr            # Main application entry point
├── internal/
│   ├── api            # JSON API handlers
│   ├── domain         # Core business logic & models
│   ├── hmrc           # HMRC API client & Fraud Header engine
│   ├── parser         # Stateless in-memory file ingestion (CSV/XLSX)
│   ├── storage        # Database wrappers (sqlc generated)
│   └── ui             # Templ components & HTMX fragments
├── db/                # Migrations and sqlc queries
├── web/               # Static assets (Tailwind, JS, Images)
└── Makefile           # Dev automation (build, test, migrate)
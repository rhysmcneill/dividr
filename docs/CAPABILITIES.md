# What Can Dividr Do?

## Overview
**Dividr** is a specialized web application that helps UK sole traders and landlords comply with HMRC's Making Tax Digital (MTD) for Income Tax Self Assessment (ITSA) requirements. It bridges the gap between traditional spreadsheet-based bookkeeping and modern digital tax reporting.

---

## 🎯 Core Capabilities

### 1. **Making Tax Digital (MTD) Compliance Bridge**
Dividr acts as an intelligent intermediary between your existing spreadsheet workflow and HMRC's digital tax requirements:

- **Digital Link Compliance**: Maintains the required electronic data flow from spreadsheets to HMRC without manual intervention
- **Automated Mapping**: Intelligently categorizes transactions into HMRC-approved income and expense categories
- **Quarterly Updates**: Manages and tracks the 4 mandatory quarterly submissions per tax year
- **Dual-Stream Support**: Handles both Sole Trade income and UK Property income separately as required by HMRC

### 2. **Session-Based Transaction Processing**
A unique hybrid approach that prioritizes data minimization while maintaining user convenience:

- **CSV Import**: Upload and parse transaction data from spreadsheets
- **Temporary Storage**: Transaction data is temporarily stored only during the active submission workflow
- **Resume Capability**: Save your progress and continue later without losing work
- **Automatic Purge**: Raw transaction data is automatically deleted after successful submission to HMRC
- **Receipt Retention**: Only submission receipts and correlation IDs are permanently stored for compliance

### 3. **8-Slot Dashboard System**
Visual tracking interface for managing compliance across tax periods:

- **Quarterly Tracking**: Monitor all 4 quarters of the tax year at a glance
- **Dual Streams**: Separate tracking for Sole Trade and UK Property income
- **Status Indicators**: See which quarters have been submitted, are in progress, or are pending
- **Compliance Overview**: Understand your MTD obligations and submission deadlines

### 4. **Stream Separator & Transaction Mapping**
Intelligent categorization system for hybrid business structures:

- **Income Classification**: Automatically identify and separate sole trade vs. property income
- **Expense Allocation**: Correctly attribute expenses to the appropriate income stream
- **Category Mapping**: Map your spreadsheet columns to HMRC's standardized categories
- **Reusable Profiles**: Save mapping configurations for future imports

### 5. **HMRC Integration**
Secure, direct communication with HMRC systems:

- **OAuth2 Authentication**: Secure authorization using HMRC's official OAuth2 flow
- **Test & Production**: Support for both HMRC Sandbox (testing) and Production environments
- **Fraud Prevention Headers**: Automatic generation of required anti-fraud metadata
- **Token Management**: Secure storage and automatic refresh of authentication tokens
- **API Versioning**: Uses HMRC's versioned API endpoints for stability

### 6. **User Authentication & Session Management**
Secure user account system:

- **User Registration**: Create accounts with email and password
- **Secure Login**: BCrypt password hashing for security
- **Session Management**: PostgreSQL-backed session storage with automatic expiry
- **Protected Routes**: Middleware-based authentication for private pages
- **Graceful Logout**: Clean session termination

### 7. **Receipt Vault & Audit Trail**
Compliance-focused record keeping:

- **Immutable Receipts**: Generate tamper-proof submission records
- **HMRC Correlation IDs**: Store official submission identifiers for proof of filing
- **PDF Generation**: Create printable receipt documents
- **JSON Archives**: Machine-readable submission records
- **Audit Events**: Comprehensive logging of security and compliance events
- **Data Retention**: Automatic purging of transient data per privacy policies

### 8. **Security & Data Protection**
Enterprise-grade security features:

- **Encryption at Rest**: OAuth tokens and sensitive data encrypted using NaCl/AES-GCM
- **Content Security Policy**: XSS prevention via CSP headers with nonce validation
- **HTTPS Only**: Secure transport layer
- **Secret Management**: Docker secrets integration for credentials
- **Rate Limiting**: Protection against brute force attacks
- **Audit Logging**: Complete trail of security-relevant events

---

## 🔧 Technical Capabilities

### Application Features
- **Web Server**: HTTP server with graceful shutdown
- **RESTful API**: Versioned API at `/api/v1` for programmatic access
- **Real-time Updates**: HTMX-powered dynamic UI without page reloads
- **Responsive Design**: Tailwind CSS for mobile-friendly layouts
- **Server-Side Rendering**: Templ template engine for type-safe HTML generation

### Database Operations
- **PostgreSQL 17**: Modern relational database with full ACID compliance
- **Automatic Migrations**: Database schema versioning and updates on startup
- **Type-Safe Queries**: sqlc-generated Go code for compile-time SQL validation
- **Connection Pooling**: Efficient database resource management
- **Transaction Support**: ACID guarantees for critical operations

### Development & Operations
- **Docker Support**: Containerized deployment with Docker Compose
- **Multi-Architecture**: Supports AMD64 and ARM64 (Raspberry Pi compatible)
- **Hot Reload**: Live development server with automatic rebuilds
- **Semantic Versioning**: Automated version management and changelog generation
- **CI/CD Integration**: GitHub Actions for testing, building, and deployment
- **Health Checks**: Application health monitoring endpoints

---

## 📊 Data Flow

### Typical User Journey

1. **Registration & Authentication**
   - Create account or log in
   - Connect to HMRC via OAuth2 (one-time setup)

2. **Import Transactions**
   - Upload CSV file with income/expense data
   - System parses and validates the data
   - Data temporarily stored in session

3. **Map Categories**
   - Map CSV columns to HMRC categories
   - Optionally save mapping profile for reuse
   - Review automatic categorizations

4. **Separate Streams** (if applicable)
   - System identifies hybrid records
   - Split transactions between Sole Trade and UK Property
   - Validate totals and allocations

5. **Review & Submit**
   - Review aggregated quarterly figures
   - Submit to HMRC via API
   - Receive confirmation and correlation ID

6. **Receipt & Cleanup**
   - Generate immutable receipt (PDF/JSON)
   - Store correlation ID for compliance
   - **Automatically purge** raw transaction data
   - Update dashboard status

---

## 🔐 Privacy & Data Minimization

Dividr is designed with **privacy by design** principles:

- **No Permanent Transaction Storage**: Your detailed transaction records are never permanently stored
- **Transient Data Model**: Line-item data exists only during active sessions
- **Automatic Purging**: Raw data deleted immediately after successful HMRC submission
- **Minimal Retention**: Only submission receipts and correlation IDs kept for legal compliance
- **User Control**: Clear data lifecycle with no hidden long-term storage

---

## 🚀 Deployment Options

### Local Development
```bash
# Start PostgreSQL + Adminer
docker compose -f docker/docker-compose.yaml up -d

# Run migrations
make migrate-up

# Start development server (with hot reload)
make dev
```

### Production Deployment
```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/rhysmcneill/dividr:latest

# Run with docker-compose
docker compose up -d
```

### Supported Platforms
- Linux (AMD64, ARM64)
- macOS (via Docker)
- Windows (via Docker or WSL2)
- Raspberry Pi (ARM64)

---

## 📋 Current Status

**Version**: v0.8.0 (January 2026)

### ✅ Implemented Features
- ✅ User authentication and session management
- ✅ OAuth2 integration with HMRC
- ✅ Fraud prevention header generation
- ✅ Database migrations and schema management
- ✅ Responsive UI with Tailwind CSS
- ✅ Session-based data persistence
- ✅ Audit event logging
- ✅ Docker containerization
- ✅ Multi-architecture builds
- ✅ Automated versioning and releases

### 🚧 In Development
- 🚧 CSV import and parsing
- 🚧 Transaction mapping profiles
- 🚧 Stream separator logic
- 🚧 8-slot dashboard UI
- 🚧 HMRC API submission endpoints
- 🚧 Receipt generation (PDF/JSON)
- 🚧 Automated data purging

### 🔮 Planned Features
- 📅 Calendar-based deadline tracking
- 📱 Mobile app support
- 📊 Basic reporting and analytics
- 🔔 Submission reminders
- 📧 Email notifications
- 🌍 Multi-language support (Welsh, Scottish Gaelic)

---

## 🎓 Who Is This For?

### Primary Users
- **UK Sole Traders**: Self-employed individuals needing to file MTD returns
- **UK Landlords**: Property owners with UK rental income
- **Hybrid Workers**: Those with both sole trade and property income

### User Scenarios
1. **Spreadsheet Users**: Already track finances in Excel/Sheets and need MTD compliance
2. **Small Business Owners**: Don't need full accounting software, just MTD submission
3. **Accountant-Prepared Returns**: Using Dividr to digitally submit accountant-prepared figures
4. **DIY Tax Filers**: Comfortable managing own books but need HMRC connectivity

---

## 🔗 Integration Points

### External Services
- **HMRC OAuth2**: Authentication and authorization
- **HMRC MTD API**: Quarterly update submissions
- **PostgreSQL**: Data persistence
- **Docker Registry**: Container image distribution

### API Endpoints (Planned)
- `POST /api/v1/upload` - Upload CSV transactions
- `POST /api/v1/mapping` - Save category mappings
- `POST /api/v1/submit` - Submit to HMRC
- `GET /api/v1/receipts/:id` - Retrieve submission receipt
- `GET /api/v1/dashboard` - Get dashboard status

---

## 🛡️ Compliance & Standards

### Legal Compliance
- **MTD for ITSA**: Complies with Making Tax Digital regulations
- **Digital Link**: Maintains electronic data flow as required by law
- **Fraud Prevention**: Implements HMRC's mandatory fraud headers
- **Data Protection**: GDPR-compliant data handling

### Technical Standards
- **OAuth 2.0**: Industry-standard authorization
- **RESTful API**: Standard HTTP methods and status codes
- **Semantic Versioning**: Predictable version numbering
- **12-Factor App**: Cloud-native architecture principles

---

## 📖 Further Reading

- **README.md**: Quick start and architecture overview
- **HMRC_LIMITATIONS.md**: HMRC Sandbox testing constraints
- **AUTOMATED_VERSIONING.md**: Release management details
- **CHANGELOG.md**: Complete version history

---

## 💡 Key Differentiators

What makes Dividr unique:

1. **Session-Based Architecture**: Temporary storage with automatic purging
2. **Spreadsheet Continuity**: Keep using your existing workflow
3. **Privacy First**: Minimal data retention by design
4. **Open Source**: Transparent, auditable codebase
5. **Self-Hostable**: Run on your own infrastructure
6. **No Vendor Lock-in**: Control your own data and deployment

---

**In summary**: Dividr enables UK sole traders and landlords to digitally submit quarterly tax updates to HMRC while maintaining their existing spreadsheet workflows, with strong privacy guarantees and minimal permanent data storage.

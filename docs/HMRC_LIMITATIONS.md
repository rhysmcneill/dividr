## HMRC Sandbox Limitations:

**Statefulness**: The Sandbox is "stateful" for individual test users. Data you submit (like a Quarterly Update) persists for that user. However, you must rely on the Test Support API to reset data if a test gets messy.

**Date Logic**: The Sandbox strictly enforces tax years. Submitting data for a "future" tax year often returns 400 Bad Request unless you explicitly configure the test user for a future simulation.

**latency**: Sandbox responses can be slower than Production. Do not tune timeouts based solely on Sandbox performance.

**Authentication**: Grant_type=authorization_code flows work exactly like production, but the login screen URL starts with test-www.tax.service.gov.uk.

## Production Readiness for HMRC Portal Evidence (Draft)

1. OAuth Flow:

    [ ] Implemented Authorization Code Grant.

    [ ] State parameter used to prevent CSRF.

    [ ] Tokens encrypted at rest (NaCl SecretBox / AES-GCM).

2. Fraud Prevention Headers:

    [ ] All 5 mandatory headers present (Gov-Client-Connection-Method, Gov-Client-Public-IP, etc.).

    [ ] Validated against /test/fraud-prevention-headers/validate.

3. Data Handling (Staging & Purge):

    [ ] No file storage: Original Excel files are processed in-memory and discarded.

    [ ] Staging: Rows stored in transactions table only for draft duration.

    [ ] Purge: Successful POST /submission triggers immediate deletion of staged rows.

    [ ] TTL: Abandoned drafts purged automatically after 90 days.

4. Digital Link:

    [ ] Data flows electronically from Excel import -> JSON payload -> HMRC API without manual intervention/copy-pasting.

# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.0] - 2026-01-09

**Type:** MINOR

### Changes
- minor: defining API blueprint including /api/v1 scaffolding, hx and app routes, Set up HTTP logging, updated docker builds to inject version into app
- fix: Add sqlc generate to CI, fail if generates and the generated files arent in branch.


## [v0.1.0] - 2026-01-08

**Type:** MINOR

### Changes
- fix: fixing migrate path for database CI
- minor: Setting up docker-compose with adminer and postgres17, configuring migrations and sqlc for database schema, creating basic CI tests for testing database connectivity, creating an Run() function to keep main.go clean and implement a graceful shutdown for Dividr
- chore: clean up CHANGELOG duplicates


## [v0.0.10] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Updating CHANGELOG management in release CI


## [v0.0.9] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Adding logging to release asset CI


## [v0.0.8] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Resolving issues with CHANGELOG duplication


## [v0.0.7] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Resolving issues with CI


## [v0.0.6] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Resolving issues with CI
- patch: Resolving issues with release artifact build and attachment


## [v0.0.5] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Resolving issues with docker image tags in release management


## [v0.0.4] - 2026-01-07

**Type:** PATCH

### Changes
- patch: Resolving issues with docker image tags in release management


## [v0.0.3] - 2026-01-06

**Type:** PATCH

### Changes
- patch: Resolving CHANGELOG issues - excluding merge commits


## [v0.0.2] - 2026-01-06

**Type:** PATCH

### Changes
- patch: fixing semver script


## [v0.0.1] - 2026-01-06

**Type:** INITIAL

### Changes
- patch: setting up golang logging, Dockerfile, docker compose, semantic versioning, Go Error handling, application config
- chore: update ci
- chore: pre-commit config
- chore: initial commit

---

This changelog is automatically generated based on commit messages and semantic versioning.

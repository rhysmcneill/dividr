# Changelog

All notable changes to this project will be documented in this file.


## [v0.6.9] - 2026-01-15

**Type:** PATCH

### Changes
- patch: typo in landing page


## [v0.6.8] - 2026-01-15

**Type:** PATCH

### Changes
- patch: improvements to meta html
- patch: improvements to hero section in landing


## [v0.6.7] - 2026-01-15

**Type:** PATCH

### Changes
- patch: improvements to landing page to explain the purpose more clearly


## [v0.6.6] - 2026-01-15

**Type:** PATCH

### Changes
- patch: ignoring templ and css output in .dockerignore


## [v0.6.5] - 2026-01-15

**Type:** PATCH

### Changes
- Triggering commit for base and landing tmepl
- patch: Implementing nonce CSP header for security and preventing XSS attacks


## [v0.6.4] - 2026-01-15

**Type:** PATCH

### Changes
- patch: renaming 002 migration


## [v0.6.3] - 2026-01-15

**Type:** PATCH

### Changes
- patch: Updating build artifact web/static/css directory
- chore(release): update CHANGELOG for v0.6.3 [skip ci]
- patch: Updating build artifact web/static/css directory


## [v0.6.2] - 2026-01-15

**Type:** PATCH

### Changes
- reverting errornous 0.6.2 release
- chore(release): update CHANGELOG for v0.6.2 [skip ci]
- patch: refactoring base template


## [v0.6.1] - 2026-01-15

**Type:** PATCH

### Changes
- patch: fixing docker build to include static files


## [v0.6.0] - 2026-01-15

**Type:** MINOR

### Changes
- generating sqlc
- minor: Introducing db migrations in the app, updating meta data in html and fixing html on smaller screens
- Removing mogrations from CI
- Running CI test


## [v0.5.1] - 2026-01-14

**Type:** PATCH

### Changes
- patch: optimising builds for multi-arch docker images
- patch: Solving deployment issues with raspberry pi arm64 architecture - integrating mult-arch image builds via QEMU
- bugfix: SOlving deployment issues with raspberry pi arm64 architecture


## [v0.5.0] - 2026-01-12

**Type:** MINOR

### Changes
- minor: Create product branding, landing page, set up privacy, security and terms pages. Integrating SEO (robot.txt and sitemap.xml), integrating a waitlist and refined my frontend folder structure


## [v0.4.1] - 2026-01-11

**Type:** PATCH

### Changes
- patch: Fixing build issue compatibilities with tailwind and alpine images, changing base image to debian


## [v0.4.0] - 2026-01-11

**Type:** MINOR

### Changes
- minor: implementing auth system scs (postgres), and password hashing using bcrypt. Integrating sign up and login pages and RequireAuth handler middleware for protecting private endpoints using sessions
- chore: Adding output.css to gitignore and container names to d-c.yml
- chore: Adding output.css to gitignore and container names to d-c.yml
- chore: Adding css output.css and container names to d-c.yml


## [v0.3.0] - 2026-01-10

**Type:** MINOR

### Changes
- minor: Integrating tailwindcss, basic templ pages, base and components. Added all required routes, handlers - dashboard and landing are configured. Release build process updated to include tailwindcss.
- chore: updating v0.2.1 entry in changelog


## [v0.2.2] - 2026-01-09

**Type:** PATCH

### Changes
- patch: solving duplication in CHANGELOG via semver.sh

## [v0.2.1] - 2026-01-09

**Type:** PATCH

### Changes
- fix: Fixing semver script to check if there are duplicates in the CHANGELOG.md
- patch: Resolving app version var injection at build time

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

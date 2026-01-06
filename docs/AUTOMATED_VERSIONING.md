# Automated Semantic Versioning & Docker Build

This project uses **automated semantic versioning** based on commit messages. When you merge to `main`, the system automatically creates version tags and publishes Docker images.

---

## 🔄 How It Works

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Developer commits with conventional message              │
│    git commit -m "feat: add user authentication"            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Merge to main branch                                      │
│    - Via PR or direct push                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Auto-Semver Workflow Runs                                │
│    - Analyzes commits since last tag                         │
│    - Determines bump type (major/minor/patch)                │
│    - Creates new version tag (e.g., v1.2.0)                 │
│    - Pushes tag to repository                                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Docker Build Workflow Triggers                           │
│    - Triggered by new version tag                            │
│    - Builds Docker image                                     │
│    - Tags with multiple versions (1.2.0, 1.2, 1, latest)    │
│    - Pushes to ghcr.io/rhysmcneill/dividr                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 📝 Commit Message Convention

Use the simple **semver prefix** format:

### Format
```
<semver>: <description>
```

Where `<semver>` is one of: `major`, `minor`, or `patch`

### Version Bump Rules

| Commit Prefix | Version Bump | When to Use |
|---------------|--------------|-------------|
| `major:` | **Major** (1.0.0 → 2.0.0) | Breaking changes, incompatible API changes |
| `minor:` | **Minor** (1.0.0 → 1.1.0) | New features, backwards-compatible |
| `patch:` | **Patch** (1.0.0 → 1.0.1) | Bug fixes, performance improvements |

### Examples

#### Major Version Bump (Breaking Change)
```bash
git commit -m "major: redesign API endpoints"
git commit -m "major: remove deprecated features"
```

#### Minor Version Bump (New Feature)
```bash
git commit -m "minor: add user profile page"
git commit -m "minor: implement OAuth2 login"
```

#### Patch Version Bump (Bug Fix)
```bash
git commit -m "patch: correct validation error message"
git commit -m "patch: optimize database queries"
```

---

## 🚀 Workflow in Practice

### Scenario 1: Adding a New Feature

```bash
# 1. Create feature branch
git checkout -b feature/user-dashboard

# 2. Make changes and commit
git add .
git commit -m "minor: add user dashboard with analytics"

# 3. Push and create PR
git push origin feature/user-dashboard

# 4. Merge PR to main
# ✅ Auto-Semver runs → creates tag v1.1.0
# ✅ Docker Build runs → publishes ghcr.io/rhysmcneill/dividr:1.1.0
```

### Scenario 2: Fixing a Bug

```bash
# 1. Create fix branch
git checkout -b fix/login-redirect

# 2. Commit fix
git commit -m "patch: correct redirect after login"

# 3. Merge to main
# ✅ Auto-Semver runs → creates tag v1.1.1
# ✅ Docker Build runs → publishes ghcr.io/rhysmcneill/dividr:1.1.1
```

### Scenario 3: Breaking Change

```bash
git commit -m "major: migrate to new database schema"

# Merge to main
# ✅ Auto-Semver runs → creates tag v2.0.0
# ✅ Docker Build runs → publishes ghcr.io/rhysmcneill/dividr:2.0.0
```

---

## 🏷️ Version Tags Created

When Auto-Semver creates a tag (e.g., `v1.2.3`), Docker Build creates multiple image tags:

| Tag | Description | Example |
|-----|-------------|---------|
| `1.2.3` | Full semantic version | `ghcr.io/rhysmcneill/dividr:1.2.3` |
| `1.2` | Major + Minor | `ghcr.io/rhysmcneill/dividr:1.2` |
| `1` | Major only | `ghcr.io/rhysmcneill/dividr:1` |
| `latest` | Latest from main | `ghcr.io/rhysmcneill/dividr:latest` |

---

## 🐳 Using Docker Images

### In docker-compose.yml

```yaml
services:
  dividr:
    # Pin to a specific version (recommended for production)
    image: ghcr.io/rhysmcneill/dividr:1.2.3

    # Or use major.minor (auto-updates patches)
    image: ghcr.io/rhysmcneill/dividr:1.2

    # Or use latest (for development)
    image: ghcr.io/rhysmcneill/dividr:latest
```

### Pulling Private Images Locally

```bash
# 1. Create a GitHub Personal Access Token (PAT)
#    Settings → Developer settings → Personal access tokens → Generate new token
#    Permissions: read:packages

# 2. Login to ghcr.io
echo YOUR_GITHUB_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin

# 3. Pull the image
docker pull ghcr.io/rhysmcneill/dividr:latest

# 4. Run with docker-compose
cd docker
docker compose pull
docker compose up -d
```

---

## 🔍 Viewing Releases

### In GitHub
- **Releases**: https://github.com/rhysmcneill/dividr/releases
- **Tags**: https://github.com/rhysmcneill/dividr/tags
- **Container Registry**: https://github.com/rhysmcneill/dividr/pkgs/container/dividr

### In CLI
```bash
# List all tags
git fetch --tags
git tag -l

# View latest tag
git describe --tags --abbrev=0

# View Docker images
docker images | grep dividr
```

---

## 🛠️ Manual Override (Emergency)

If you need to create a version manually:

```bash
# Create and push tag manually
git tag v1.3.0 -m "Manual release: hotfix deployment"
git push origin v1.3.0

# Docker Build will still trigger automatically
```

---

## ⚙️ Configuration Files

- **Auto-Semver Workflow**: `.github/workflows/auto-semver.yml`
- **Docker Build Workflow**: `.github/workflows/docker-build.yml`
- **Semver Bump Script**: `.github/scripts/semver-bump.sh`

---

## 🚨 Important Notes

1. **First Release**: If no tags exist, the script starts from `v0.0.0`
2. **Skip CI**: Tag commits include `[skip ci]` to prevent infinite loops
3. **Merge Commits**: The script analyzes all commits between the last tag and HEAD
4. **No Default Behavior**: If no `major:`, `minor:`, or `patch:` prefix is found, **no tag is created**
5. **Changelog**: `CHANGELOG.md` is automatically updated with each release

---

## 📋 Changelog

The repository maintains an automated `CHANGELOG.md` file at the root. Each release includes:
- Version number
- Release date
- Bump type (MAJOR/MINOR/PATCH)
- All commit messages since the last version

Example entry:
```markdown
## [v1.2.0] - 2026-01-06

**Type:** MINOR

### Changes
- minor: add user dashboard with analytics
- patch: fix navigation bug
```

View the full changelog: [CHANGELOG.md](../CHANGELOG.md)

---

## 📊 Version History Example

```
v0.1.0 → Initial release
v0.2.0 → minor: add authentication
v0.2.1 → patch: login redirect issue
v0.3.0 → minor: add user dashboard
v1.0.0 → major: production-ready API
v1.0.1 → patch: database connection pool
v1.1.0 → minor: add export functionality
```

---

## 🔗 References

- [Semantic Versioning](https://semver.org/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)

#!/usr/bin/env bash
set -euo pipefail

# Local test version of semver-bump.sh (no GitHub API calls or pushes)
# This simulates what would happen in CI

echo "=========================================="
echo "LOCAL SEMVER SCRIPT TEST (DRY RUN)"
echo "=========================================="
echo ""

# Ensure full history
git fetch --tags --prune 2>/dev/null || true

last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -z "$last_tag" ]; then
  echo "[INFO]: No tags found. Would create initial tag v0.0.1"
  exit 0
fi

echo "[INFO]: Last tag: $last_tag"

# Validate that the tag actually exists in git history
if ! git rev-parse "$last_tag" >/dev/null 2>&1; then
  echo "[ERROR]: Tag $last_tag does not exist in git history"
  exit 1
fi

commits=$(git log --no-merges --pretty=%B "$last_tag"..HEAD)
if [ -z "$commits" ]; then
  echo "[INFO]: No new commits since $last_tag; nothing to do."
  exit 0
fi

echo "[INFO]: Found commits since $last_tag:"
echo "$commits" | head -5
if [ "$(echo "$commits" | wc -l)" -gt 5 ]; then
  echo "... and $(($(echo "$commits" | wc -l) - 5)) more"
fi
echo ""

bump="none"
# Major
if echo "$commits" | grep -qi "^major:"; then
  bump=major
  echo "[INFO]: Found 'major:' commit"
# Minor
elif echo "$commits" | grep -qi "^minor:"; then
  bump=minor
  echo "[INFO]: Found 'minor:' commit"
# Patch
elif echo "$commits" | grep -qi "^patch:"; then
  bump=patch
  echo "[INFO]: Found 'patch:' commit"
fi

if [ "$bump" = "none" ]; then
  echo "[INFO]: No semver-relevant commits found; nothing to do."
  exit 0
fi

# strip leading v
ver=${last_tag#v}
IFS='.' read -r major minor patch <<< "$ver"

major=${major:-0}
minor=${minor:-0}
patch=${patch:-0}

case "$bump" in
  major)
    major=$((major+1))
    minor=0
    patch=0
    ;;
  minor)
    minor=$((minor+1))
    patch=0
    ;;
  patch)
    patch=$((patch+1))
    ;;
esac

new_tag="v${major}.${minor}.${patch}"

echo "[INFO]: Bump type: $bump -> New tag: $new_tag"
echo ""

# CRITICAL: Check if this tag already exists (prevents duplicate releases)
echo "[TEST]: Checking if tag $new_tag already exists..."
if git rev-parse "$new_tag" >/dev/null 2>&1; then
  echo "❌ [WOULD EXIT]: Tag $new_tag already exists in repository"
  echo "[INFO]: Script would skip release to prevent duplicates"
  exit 0
else
  echo "✅ [PASS]: Tag $new_tag does not exist yet"
fi
echo ""

changelog_file="CHANGELOG.md"

# Check if this version already exists in the changelog (additional safety check)
echo "[TEST]: Checking if $new_tag already exists in CHANGELOG.md..."
if grep -q "^## \[$new_tag\]" "$changelog_file"; then
  echo "❌ [WOULD EXIT]: $new_tag already exists in CHANGELOG.md but tag doesn't exist"
  echo "[INFO]: This is inconsistent state - script would exit to prevent duplicates"
  exit 0
else
  echo "✅ [PASS]: $new_tag not found in CHANGELOG.md"
fi
echo ""

# Get current date
release_date=$(date +%Y-%m-%d)

# Extract commit messages for this release (exclude merge commits)
commit_messages=$(git log --no-merges --pretty="- %s" "$last_tag"..HEAD)

# Create changelog entry
changelog_entry="## [$new_tag] - $release_date

**Type:** ${bump^^}

### Changes
$commit_messages

"

echo "[TEST]: Testing CHANGELOG insertion logic..."
# Test the awk command without modifying the file
temp_test=$(mktemp)
awk -v entry="$changelog_entry" '/^## \[/{if(!found){print entry; found=1}} {print}' "$changelog_file" > "$temp_test"

# Verify insertion
entry_count=$(grep -c "^## \[$new_tag\]" "$temp_test" || echo "0")
first_version=$(grep -m 1 "^## \[" "$temp_test" | grep -oP '\[v[^\]]+\]' | tr -d '[]')

if [ "$entry_count" = "1" ] && [ "$first_version" = "$new_tag" ]; then
  echo "✅ [PASS]: Entry would be inserted exactly once at the top"
  echo ""
  echo "Preview of CHANGELOG.md changes:"
  echo "-----------------------------------"
  head -20 "$temp_test"
  echo "-----------------------------------"
else
  echo "❌ [FAIL]: Insertion test failed (count=$entry_count, first=$first_version)"
  rm "$temp_test"
  exit 1
fi

rm "$temp_test"

echo ""
echo "=========================================="
echo "✅ ALL CHECKS PASSED"
echo "=========================================="
echo ""
echo "In CI, the script would:"
echo "  1. Update CHANGELOG.md with the entry above"
echo "  2. Commit: 'chore(release): update CHANGELOG for $new_tag [skip ci]'"
echo "  3. Create annotated tag: $new_tag"
echo "  4. Push commit and tag to origin/main"
echo "  5. Create GitHub Release"
echo ""
echo "No duplicates would be created! ✨"

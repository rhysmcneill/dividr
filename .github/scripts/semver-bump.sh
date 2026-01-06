#!/usr/bin/env bash
set -euo pipefail

# Simple semantic version bump based on commit messages since last tag.
# Rules:
# - If any commit starts with "major:" -> major
# - Else if any commit starts with "minor:" -> minor
# - Else if any commit starts with "patch:" -> patch
# - Else -> no tag (exit 0)

# Require these env vars in CI: GITHUB_TOKEN, GITHUB_REPOSITORY

GITHUB_TOKEN=${GITHUB_TOKEN:-}
GITHUB_REPOSITORY=${GITHUB_REPOSITORY:-$(git config --get remote.origin.url | sed -E 's|.*/(.*)\.git$|\1|')}

if [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_TOKEN not set; exiting (script intended for CI)."
  exit 1
fi

# Ensure full history
git fetch --tags --prune

last_tag=$(git describe --tags --abbrev=0 2>/dev/null || true)
if [ -z "$last_tag" ]; then
  last_tag="v0.0.0"
fi

echo "Last tag: $last_tag"

commits=$(git log --pretty=%B "$last_tag"..HEAD || true)
if [ -z "$commits" ]; then
  echo "No new commits since $last_tag; nothing to do."
  exit 0
fi

bump="none"
# Major
if echo "$commits" | grep -qi "^major:"; then
  bump=major
# Minor
elif echo "$commits" | grep -qi "^minor:"; then
  bump=minor
# Patch
elif echo "$commits" | grep -qi "^patch:"; then
  bump=patch
fi

if [ "$bump" = "none" ]; then
  echo "No semver-relevant commits found; nothing to do."
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

echo "Bump type: $bump -> New tag: $new_tag"

# Update CHANGELOG.md
changelog_file="CHANGELOG.md"
if [ ! -f "$changelog_file" ]; then
  echo "# Changelog" > "$changelog_file"
  echo "" >> "$changelog_file"
  echo "All notable changes to this project will be documented in this file." >> "$changelog_file"
  echo "" >> "$changelog_file"
fi

# Get current date
release_date=$(date +%Y-%m-%d)

# Extract commit messages for this release
commit_messages=$(git log --pretty="- %s" "$last_tag"..HEAD)

# Create changelog entry
changelog_entry="## [$new_tag] - $release_date

**Type:** ${bump^^}

### Changes
$commit_messages

"

# Insert at the top of the changelog (after the header)
temp_file=$(mktemp)
if grep -q "^## \[" "$changelog_file"; then
  # Insert before first version entry
  awk -v entry="$changelog_entry" '/^## \[/{print entry; found=1} {print}' "$changelog_file" > "$temp_file"
else
  # No version entries yet, append after header
  awk -v entry="$changelog_entry" '1; /^All notable changes/{print ""; print entry}' "$changelog_file" > "$temp_file"
fi
mv "$temp_file" "$changelog_file"

echo "Updated CHANGELOG.md"

# Create annotated tag and push it
git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

# Commit changelog
git add "$changelog_file"
git commit -m "chore(release): update CHANGELOG for $new_tag [skip ci]"

git tag -a "$new_tag" -m "chore(release): $new_tag [skip ci]"

remote_repo="https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"

echo "Pushing changelog commit and tag to $remote_repo"

git push "$remote_repo" HEAD:main
git push "$remote_repo" "$new_tag"

echo "Changelog updated and tag pushed: $new_tag"

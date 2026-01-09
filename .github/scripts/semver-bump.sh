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
  echo "[ERROR]: GITHUB_TOKEN not set; exiting (script intended for CI)."
  exit 1
fi

# Ensure full history
git fetch --tags --prune

last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -z "$last_tag" ]; then
  echo "[INFO]: No tags found in repository. Creating initial tag v0.0.1"
  # For the very first tag, check if there are any commits at all
  if ! git rev-parse HEAD >/dev/null 2>&1; then
    echo "[ERROR]: No commits in repository"
    exit 1
  fi
  # Create initial tag and exit
  new_tag="v0.0.1"

  # Update CHANGELOG.md
  changelog_file="CHANGELOG.md"
  if [ ! -f "$changelog_file" ]; then
    echo "# Changelog" > "$changelog_file"
    echo "" >> "$changelog_file"
    echo "All notable changes to this project will be documented in this file." >> "$changelog_file"
    echo "" >> "$changelog_file"
  fi

  release_date=$(date +%Y-%m-%d)
  commit_messages=$(git log --no-merges --pretty="- %s" HEAD)

  changelog_entry="## [$new_tag] - $release_date

**Type:** INITIAL

### Changes
$commit_messages

"

  temp_file=$(mktemp)
  if grep -q "^## \[" "$changelog_file"; then
    awk -v entry="$changelog_entry" '/^## \[/{print entry; found=1} {print}' "$changelog_file" > "$temp_file"
  else
    awk -v entry="$changelog_entry" '1; /^All notable changes/{print ""; print entry}' "$changelog_file" > "$temp_file"
  fi
  mv "$temp_file" "$changelog_file"

  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

  git add "$changelog_file"
  git commit -m "chore(release): update CHANGELOG for $new_tag [skip ci]"
  git tag -a "$new_tag" -m "chore(release): $new_tag [skip ci]"

  remote_repo="https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
  git push "$remote_repo" HEAD:main
  git push "$remote_repo" "$new_tag"

  echo "[INFO]: Initial tag created and pushed: $new_tag"

  # Create initial GitHub Release
  echo "[INFO]: Creating initial GitHub Release for $new_tag"

  release_response=$(jq -n \
    --arg tag "$new_tag" \
    --arg name "Release $new_tag" \
    --arg body "## Initial Release\n\n$commit_messages" \
    '{
      tag_name: $tag,
      name: $name,
      body: $body,
      draft: false,
      prerelease: false
    }' | curl -X POST \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/${GITHUB_REPOSITORY}/releases" \
    -d @-)

  if echo "$release_response" | grep -q '"id"'; then
    echo "[INFO]: GitHub Release created successfully: $new_tag"
  else
    echo "[ERROR]: Failed to create GitHub Release"
    echo "Response: $release_response"
  fi

  exit 0
fi

echo "Last tag: $last_tag"

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

echo "Bump type: $bump -> New tag: $new_tag"

# CRITICAL: Check if this tag already exists (prevents duplicate releases)
if git rev-parse "$new_tag" >/dev/null 2>&1; then
  echo "[WARN]: Tag $new_tag already exists in repository, skipping release"
  exit 0
fi

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

# Extract commit messages for this release (exclude merge commits)
commit_messages=$(git log --no-merges --pretty="- %s" "$last_tag"..HEAD)

# Create changelog entry
changelog_entry="## [$new_tag] - $release_date

**Type:** ${bump^^}

### Changes
$commit_messages

"

# Check if this version already exists in the changelog (additional safety check)
if grep -q "^## \[$new_tag\]" "$changelog_file"; then
  echo "[WARN]: $new_tag already exists in CHANGELOG.md but tag doesn't exist - this is inconsistent!"
  echo "[INFO]: Cleaning up and exiting to prevent duplicates"
  exit 0
fi

# Insert at the top of the changelog (after the header)
temp_file=$(mktemp)
if grep -q "^## \[" "$changelog_file"; then
  # Insert before first version entry
  awk -v entry="$changelog_entry" '/^## \[/{if(!found){print entry; found=1}} {print}' "$changelog_file" > "$temp_file"
else
  # No version entries yet, append after header
  awk -v entry="$changelog_entry" '1; /^All notable changes/{print ""; print entry}' "$changelog_file" > "$temp_file"
fi
mv "$temp_file" "$changelog_file"
echo "[INFO]: Updated CHANGELOG.md"

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

echo "[INFO]: Changelog updated and tag pushed: $new_tag"

# Create GitHub Release
echo "[INFO]: Creating GitHub Release for $new_tag"

# Extract the changelog entry for this version
changelog_body=$(awk -v tag="$new_tag" '
  /^## \[/ {
    if (found) exit;
    if ($0 ~ tag) { found=1; next }
  }
  found && /^## \[/ { exit }
  found { print }
' "$changelog_file")

# Create release using GitHub API with proper JSON encoding
release_response=$(jq -n \
  --arg tag "$new_tag" \
  --arg name "Release $new_tag" \
  --arg body "$changelog_body" \
  '{
    tag_name: $tag,
    name: $name,
    body: $body,
    draft: false,
    prerelease: false
  }' | curl -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/${GITHUB_REPOSITORY}/releases" \
  -d @-)

if echo "$release_response" | grep -q '"id"'; then
  echo "[INFO]: ✅ GitHub Release created successfully: $new_tag"
else
  echo "[ERROR]: ⚠️ Failed to create GitHub Release"
  echo "Response: $release_response"
  # Don't fail the script if release creation fails
fi

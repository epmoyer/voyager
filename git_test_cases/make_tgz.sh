#!/usr/bin/env bash
#
# Re-create git_test_cases.tgz from the test-case directories in this folder.
#
# Each sub-directory here is a small git repo used by demo.py to exercise the
# prompt under different conditions. We tar them up so they can be committed to
# git as a single blob (see README). This script regenerates that archive.
#
# Works on both macOS (BSD tar) and Linux (GNU tar). On macOS it avoids
# embedding AppleDouble (._*) resource-fork files in the archive.
#
set -euo pipefail

# Operate on this script's own directory regardless of where it's invoked from.
cd "$(dirname "${BASH_SOURCE[0]}")"

ARCHIVE="git_test_cases.tgz"

# Prefer GNU tar (gtar) when available. macOS ships BSD tar, which embeds
# extended attributes that surface as ._* files when the archive is extracted
# on Linux; gtar avoids that. Either tar works with the settings below.
if command -v gtar >/dev/null 2>&1; then
    TAR=gtar
else
    TAR=tar
fi

# Stop macOS's BSD tar from writing AppleDouble (._*) entries for xattrs/resource
# forks. GNU tar ignores this and doesn't store xattrs unless asked, so setting
# it unconditionally is safe on both platforms.
export COPYFILE_DISABLE=1

# Gather the test-case directories: every sub-directory here, which naturally
# excludes this script, the archive itself, and .gitkeep. nullglob makes the
# loop produce nothing (rather than a literal "*/") when there are no matches.
shopt -s nullglob
items=()
for entry in */ ; do
    items+=("${entry%/}")
done

if [ ${#items[@]} -eq 0 ]; then
    echo "No test-case directories found in $(pwd); nothing to archive." >&2
    exit 1
fi

# Overwrite any existing archive.
rm -f "$ARCHIVE"

echo "Creating $ARCHIVE with $TAR from: ${items[*]}"
"$TAR" \
    --exclude='._*' \
    --exclude='.DS_Store' \
    -czf "$ARCHIVE" \
    "${items[@]}"

echo "Done: $(pwd)/$ARCHIVE"

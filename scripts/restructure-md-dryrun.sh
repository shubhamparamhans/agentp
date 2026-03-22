#!/usr/bin/env bash
set -euo pipefail

ROOT="$(pwd)"
OUT_DIR="docs/features"

dry_run=false
force=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run) dry_run=true; shift;;
    --force) force=true; shift;;
    *) echo "Unknown arg: $1"; exit 1;;
  esac
done

echo "Scanning for .md files (excluding docs/features)..."
files=$(find . -maxdepth 3 -type f -name '*.md' ! -path './docs/features/*' | sed 's|^./||' | sort)

if [ -z "$files" ]; then
  echo "No markdown files found."; exit 0
fi

printf "%s\n" "Planned mappings (source -> target):"
while IFS= read -r src; do
  name=$(echo "$src" | tr '/' '-' | tr '[:upper:]' '[:lower:]' | sed 's/\.md$//; s/[^a-z0-9.-]/-/g; s/--*/-/g')
  target_dir="$OUT_DIR/$name"
  target_detailed="$target_dir/detailed.md"
  printf "%s -> %s\n" "$src" "$target_detailed"
done <<EOF
$files
EOF

if [ "$dry_run" = true ]; then
  echo "\nDry-run complete. No files were written."; exit 0
fi

echo "Running actual migration..."
for src in "${files[@]}"; do
  name=$(echo "$src" | tr '/' '-' | tr '[:upper:]' '[:lower:]' | sed 's/\.md$//; s/[^a-z0-9.-]/-/g; s/--*/-/g')
  target_dir="$OUT_DIR/$name"
  target_detailed="$target_dir/detailed.md"
  target_summary="$target_dir/summary.md"
  target_changelog="$target_dir/changelog.md"

  mkdir -p "$target_dir"

  if [ -f "$target_detailed" ] && [ "$force" = false ]; then
    echo "Skipping $target_detailed (exists). Use --force to overwrite."; continue
  fi

  printf "# Source: %s\n\n" "$src" > "$target_detailed"
  cat "$src" >> "$target_detailed"

  if [ ! -f "$target_summary" ]; then
    cat > "$target_summary" <<EOF
# TODO: summary for $name
Branch: TODO
Status: TODO
EOF
  fi

  if [ ! -f "$target_changelog" ]; then
    printf "# Changelog\n\n%s — migrated from %s\n" "$(date +%F)" "$src" > "$target_changelog"
  fi
done

echo "Migration complete. Review created folders under $OUT_DIR"

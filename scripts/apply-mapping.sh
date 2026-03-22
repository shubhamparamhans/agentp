#!/usr/bin/env bash
set -euo pipefail

OUT_CSV="docs/features/feat-restructure-md-files/mapping.csv"
force=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --force) force=true; shift;;
    *) echo "Unknown arg: $1"; exit 1;;
  esac
done

if [ ! -f "$OUT_CSV" ]; then
  echo "Mapping CSV not found: $OUT_CSV"; exit 1
fi

echo "Applying mapping from $OUT_CSV"

tail -n +2 "$OUT_CSV" | while IFS=, read -r src target; do
  # trim possible whitespace
  src="$(echo "$src" | sed 's/^ *//; s/ *$//')"
  target="$(echo "$target" | sed 's/^ *//; s/ *$//')"
  if [ -z "$src" ] || [ -z "$target" ]; then
    continue
  fi

  if [ ! -f "$src" ]; then
    echo "Warning: source file not found: $src"; continue
  fi

  target_dir=$(dirname "$target")
  mkdir -p "$target_dir"

  if [ -f "$target" ] && [ "$force" = false ]; then
    echo "Skipping existing: $target"
    continue
  fi

  printf "# Source: %s\n\n" "$src" > "$target"
  cat "$src" >> "$target"

  summary_file="$target_dir/summary.md"
  changelog_file="$target_dir/changelog.md"

  if [ ! -f "$summary_file" ]; then
    cat > "$summary_file" <<EOF
# TODO: summary for $(basename "$target_dir")
Branch: TODO
Status: TODO
EOF
  fi

  if [ ! -f "$changelog_file" ]; then
    printf "# Changelog\n\n%s — migrated from %s\n" "$(date +%F)" "$src" > "$changelog_file"
  fi

  echo "Created: $target"
done

echo "Apply mapping complete. Review new files under docs/features/" 

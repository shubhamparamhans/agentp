#!/usr/bin/env bash
set -euo pipefail

# Consolidate redundant feature folders into canonical feature assets.
# Usage: bash scripts/consolidate-features.sh

ROOT="$(pwd)"
FEATURE_DIR="docs/features"
ARCHIVE_DIR="$FEATURE_DIR/_archive_$(date +%s)"

GROUP_KEYS=(mongodb frontend postgres query)
CANONICALS=(mongodb-support frontend-support postgres-support query-support)

mkdir -p "$ARCHIVE_DIR"

echo "Canonical groups:"
for i in "${!GROUP_KEYS[@]}"; do
  k=${GROUP_KEYS[$i]}
  echo " - $k -> ${CANONICALS[$i]}"
done

for i in "${!GROUP_KEYS[@]}"; do
  key=${GROUP_KEYS[$i]}
  canon=${CANONICALS[$i]}
  canon_assets="$FEATURE_DIR/$canon/assets"
  mkdir -p "$canon_assets"

  # find folders matching patterns (case-insensitive)
  matches=$(ls -1 "$FEATURE_DIR" | tr '[:upper:]' '[:lower:]' | grep "$key" || true)
  if [ -z "$matches" ]; then
    echo "No folders found for group '$key'"
    continue
  fi

  echo "Processing group '$key' -> canonical '$canon'"
  # iterate actual folder names
  for f in $(ls -1 "$FEATURE_DIR"); do
    low=$(echo "$f" | tr '[:upper:]' '[:lower:]')
    if [[ "$low" == *"$key"* ]]; then
      # skip the canonical folder itself
      if [ "$f" = "$canon" ]; then
        continue
      fi
      src_dir="$FEATURE_DIR/$f"
      if [ -d "$src_dir" ]; then
        echo " - Archiving $src_dir to $ARCHIVE_DIR/$f"
        mv "$src_dir" "$ARCHIVE_DIR/"
        # copy files from archived folder into canonical assets with prefix
        mkdir -p "$canon_assets"
        if [ -d "$ARCHIVE_DIR/$f" ]; then
          for file in $(find "$ARCHIVE_DIR/$f" -maxdepth 1 -type f -name '*.md' -print); do
            base=$(basename "$file")
            dest="$canon_assets/${f}__${base}"
            echo "   copying $file -> $dest"
            cp "$file" "$dest"
          done
        fi
      fi
    fi
  done
done

echo "Consolidation complete. Archived folders moved to $ARCHIVE_DIR"
echo "Review copies in each canonical assets/ folder, then remove $ARCHIVE_DIR if satisfied."

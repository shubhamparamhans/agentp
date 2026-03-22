# Feature: Restructure existing Markdown files into feature folders

Feature Name: feat-restructure-md-files
Branch: feat/restructure-md-files
Author: shubhamparamhans
Date: 2026-03-15
Status: draft

Short Summary
-------------
Move (copy) every existing Markdown documentation file into a dedicated feature folder under `docs/features/` using the naming and layout rules in `docs/features/README.md`. This is a two-phase change: documentation-only migration (this feature) followed by optional cleanup.

Goals
-----
- Produce a deterministic mapping from current `.md` files to `docs/features/<feature>/detailed.md`.
- Keep originals unchanged during review; each copy must contain a header indicating original path.
- Provide a dry-run script and clear review/rollback steps.

Design & Rationale
-------------------
Approach:
- Implement the migration as a shell script that performs a dry-run and an actual copy step.
- The script will:
  1. Traverse project for `.md` files (excluding `docs/features/*`).
  2. Generate a mapping file (CSV) of source -> target folder (cleaned name).
  3. Create `docs/features/<feature>/detailed.md` with a header referencing the original path and then the original content.
  4. Create `summary.md` with a `TODO: add summary` placeholder and `changelog.md` noting the migration.

Safety measures:
- Dry-run mode prints the planned actions without writing files.
- The script refuses to overwrite an existing `detailed.md` unless `--force` is passed.
- No deletions are performed by default.

Script (example)
-----------------
Save this as `scripts/restructure-md-dryrun.sh` and run with `bash scripts/restructure-md-dryrun.sh --dry-run`.

```sh
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

echo "Scanning for .md files..."
mapfile -t files < <(find . -maxdepth 3 -type f -name '*.md' ! -path './docs/features/*' | sed 's|^./||')

if [ ${#files[@]} -eq 0 ]; then
  echo "No markdown files found."; exit 0
fi

for src in "${files[@]}"; do
  # derive clean name: replace slashes with -, lowercase
  name=$(echo "$src" | tr '/' '-' | tr '[:upper:]' '[:lower:]' | sed 's/\.md$//; s/[^a-z0-9.-]/-/g; s/--*/-/g')
  target_dir="$OUT_DIR/$name"
  target_detailed="$target_dir/detailed.md"
  target_summary="$target_dir/summary.md"
  target_changelog="$target_dir/changelog.md"

  echo "Will create: $target_detailed (from $src)"

  if [ "$dry_run" = true ]; then
    continue
  fi

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

echo "Done. Review created folders under $OUT_DIR"
```

Testing & Verification
----------------------
- Run in `--dry-run` mode first, review the planned mapping.
- Run without `--dry-run` to create the files; run `git status` to inspect added files.
- Review diffs and open a draft PR for team review.

Rollback Plan
-------------
- Because this script only copies files, rollback is simply removing the created `docs/features/<...>` folders (git revert/commit can be used to remove them after review).

Edge cases
----------
- Files with the same cleaned name (e.g., `docs/readme.md` and `README.md`) will collide; the script warns and uses the cleaned name. Manual review required for collisions.

Next steps
----------
1. Review this detailed plan and the example script.
2. Approve a dry-run execution.
3. After dry-run review, run the script to create the feature folders and open a draft PR.

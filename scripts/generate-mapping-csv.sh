#!/usr/bin/env bash
set -euo pipefail

OUT_CSV="docs/features/feat-restructure-md-files/mapping.csv"

bash scripts/restructure-md-dryrun.sh --dry-run | awk '/->/ {gsub(/ -> /,","); print}' > "$OUT_CSV"
echo "CSV written to $OUT_CSV"

echo "\nCollisions (targets with more than one source):"
awk -F, '{print $2}' "$OUT_CSV" | sort | uniq -c | awk '$1>1{print $0}' || true

echo "\nDetailed collision groups:"
awk -F, '{a[$2]=(a[$2]?a[$2]";":"")$1} END{for(i in a) if(split(a[i],b,";")>1) print i ":" a[i]}' "$OUT_CSV" || true

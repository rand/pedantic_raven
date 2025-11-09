#!/bin/bash
# Regenerate all pedantic_raven diagrams with light/dark variants

set -e

cd "$(dirname "$0")"
mkdir -p images

echo "Regenerating pedantic_raven diagrams..."

# Generate both light and dark variants for all diagrams
for diagram in diagrams/*.d2; do
  if [ -f "$diagram" ]; then
    basename=$(basename "$diagram" .d2)
    echo "  Rendering: $basename (light + dark)"
    d2 --theme=0 "$diagram" "images/${basename}-light.svg"
    d2 --theme=200 "$diagram" "images/${basename}-dark.svg"
  fi
done

echo "Done! Generated diagrams:"
ls -lh images/

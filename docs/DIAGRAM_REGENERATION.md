# Diagram Regeneration Guide for Pedantic Raven

This guide explains how to regenerate diagrams for the Pedantic Raven GitHub Pages site.

## Overview

Pedantic Raven uses [D2](https://d2lang.com/) for diagram generation with automated GitHub Actions workflow. Diagrams are stored in `docs/diagrams/` as `.d2` source files and automatically rendered to SVG with both light and dark theme variants on every push.

## Prerequisites

### Local Development

Install D2:
```bash
curl -fsSL https://d2lang.com/install.sh | sh -s --
```

Add to your PATH (if not already):
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Automated (GitHub Actions)

The GitHub Actions workflow automatically installs D2 and renders diagrams on every push to `main` that modifies `docs/` files.

## Current Diagram Inventory

The following 10 diagrams exist in `docs/diagrams/`:
1. `01-system-context.d2` - System context and boundaries
2. `02-component-architecture.d2` - Component structure
3. `03-semantic-pipeline.d2` - Semantic processing pipeline
4. `04-entity-comparison.d2` - Entity extraction comparison
5. `05-triple-graph.d2` - RDF triple graph structure
6. `06-data-flow.d2` - Data flow through system
7. `07-offline-state.d2` - Offline-first state management
8. `08-mode-switching.d2` - Mode switching architecture
9. `09-connection-lifecycle.d2` - Connection lifecycle
10. `10-event-system.d2` - Event system pub/sub

## Automatic Regeneration (Recommended)

Diagrams are **automatically regenerated** on every push via GitHub Actions:

```bash
# 1. Edit diagram source
vim docs/diagrams/01-system-context.d2

# 2. Commit and push
git add docs/diagrams/01-system-context.d2
git commit -m "docs: Update system context diagram"
git push

# 3. GitHub Actions automatically renders light and dark SVGs
# Output: docs/images/01-system-context-light.svg
#         docs/images/01-system-context-dark.svg
```

The workflow is defined in `.github/workflows/deploy-pages.yml`.

## Manual Local Regeneration

### Regenerate All Diagrams

```bash
cd /Users/rand/src/pedantic_raven/docs

# Create output directory
mkdir -p images

# Regenerate all with light and dark themes
for diagram in diagrams/*.d2; do
  basename=$(basename "$diagram" .d2)
  echo "Rendering: $basename (light and dark)"

  # Light theme (theme 0)
  d2 --theme=0 "$diagram" "images/${basename}-light.svg"

  # Dark theme (theme 200)
  d2 --theme=200 "$diagram" "images/${basename}-dark.svg"
done

echo "Done! Generated $(ls images/*.svg | wc -l) SVG files"
```

### Regenerate Single Diagram

```bash
cd /Users/rand/src/pedantic_raven/docs

# For a specific diagram (e.g., semantic-pipeline)
d2 --theme=0 diagrams/03-semantic-pipeline.d2 images/03-semantic-pipeline-light.svg
d2 --theme=200 diagrams/03-semantic-pipeline.d2 images/03-semantic-pipeline-dark.svg
```

## GitHub Actions Workflow

The automated workflow (`.github/workflows/deploy-pages.yml`) performs:

1. **Checkout**: Clones repository
2. **Setup D2**: Installs D2 diagram tool
3. **Render Diagrams**: Generates light/dark SVGs for all `.d2` files
4. **Deploy**: Publishes to GitHub Pages

Key workflow snippet:
```yaml
- name: Render D2 diagrams
  run: |
    mkdir -p docs/images
    for diagram in docs/diagrams/*.d2; do
      basename=$(basename "$diagram" .d2)
      echo "Rendering: $basename (light and dark)"
      d2 --theme=0 "$diagram" "docs/images/${basename}-light.svg"
      d2 --theme=200 "$diagram" "docs/images/${basename}-dark.svg"
    done
```

## Theme Variants

D2 theme options:
- `--theme=0` - Light theme (white background, dark text)
- `--theme=200` - Dark theme (dark background, light text)

## Diagram Styling Consistency

To match RUNE's diagram styling:
1. Use teal accent color (#16A085 light, #1ABC9C dark) for Pedantic Raven theme
2. Set proper padding and margins in D2 files
3. Use clear, readable fonts (Geist, JetBrains Mono)
4. Test both light and dark themes locally before pushing

## HTML Integration

Diagrams are integrated with theme-aware CSS in `index.html`:

```html
<!-- Picture element with media query for automatic theme switching -->
<div class="diagram-container">
    <picture>
        <source srcset="images/01-system-context-dark.svg"
                media="(prefers-color-scheme: dark)">
        <img src="images/01-system-context-light.svg"
             alt="System Context Diagram"
             class="diagram-svg">
    </picture>
</div>
```

CSS handles theme toggle:
```css
/* Diagram styling */
.diagram-container {
    margin: 2rem 0;
    text-align: center;
}

.diagram-svg {
    max-width: 100%;
    height: auto;
    border: 1px solid var(--border);
    border-radius: 8px;
}

/* Theme-aware switching via picture element handles this automatically */
```

## Verification Checklist

After regeneration:
- [ ] All 10 diagrams have both `-light.svg` and `-dark.svg` variants
- [ ] File sizes reasonable (< 200KB per diagram)
- [ ] View in browser with light theme - text readable
- [ ] View in browser with dark theme - text readable
- [ ] Toggle theme - diagrams switch correctly
- [ ] Colors match Pedantic Raven teal theme
- [ ] No rendering artifacts or clipping

## Troubleshooting

**Missing dark variants**:
```bash
# Check which diagrams are missing dark variants
ls docs/images/*-light.svg | while read f; do
  dark="${f/-light/-dark}"
  [ ! -f "$dark" ] && echo "Missing: $dark"
done
```

**D2 not found locally**:
```bash
# Verify installation
which d2
d2 --version

# Reinstall if needed
curl -fsSL https://d2lang.com/install.sh | sh -s --
```

**GitHub Actions failing**:
1. Check workflow logs in GitHub Actions tab
2. Verify `.d2` files have valid syntax
3. Test locally with same D2 version as workflow
4. Check file permissions on `.d2` files

**Diagrams not updating on site**:
1. Hard refresh browser (Cmd+Shift+R / Ctrl+Shift+F5)
2. Check GitHub Pages deployment status
3. Verify workflow completed successfully
4. Check `images/` directory has updated timestamps

## Adding New Diagrams

1. Create new `.d2` file in `docs/diagrams/`:
```bash
vim docs/diagrams/11-new-feature.d2
```

2. Write D2 diagram code:
```d2
# Example diagram
title: New Feature Architecture {
  near: top-center
  shape: text
}

component A: Frontend {
  shape: rectangle
}

component B: Backend {
  shape: rectangle
}

A -> B: API Calls
```

3. Test locally:
```bash
d2 --theme=0 docs/diagrams/11-new-feature.d2 docs/images/11-new-feature-light.svg
d2 --theme=200 docs/diagrams/11-new-feature.d2 docs/images/11-new-feature-dark.svg
```

4. View in browser to verify both themes

5. Commit and push - GitHub Actions will regenerate automatically:
```bash
git add docs/diagrams/11-new-feature.d2
git commit -m "docs: Add new feature architecture diagram"
git push
```

6. Add to `index.html`:
```html
<div class="diagram-container">
    <picture>
        <source srcset="images/11-new-feature-dark.svg"
                media="(prefers-color-scheme: dark)">
        <img src="images/11-new-feature-light.svg"
             alt="New Feature Architecture"
             class="diagram-svg">
    </picture>
</div>
```

## Best Practices

1. **Version Control**: Always commit `.d2` source files, not just SVG outputs
2. **Naming**: Use numbered prefixes (01-, 02-, etc.) for logical ordering
3. **Testing**: Test both themes locally before pushing
4. **Size**: Keep diagrams focused - split complex diagrams into multiple files
5. **Consistency**: Follow existing diagram style and color scheme
6. **Documentation**: Update this guide when adding new diagram types

## References

- [D2 Documentation](https://d2lang.com/)
- [D2 Language Tour](https://d2lang.com/tour/intro/)
- [D2 Themes](https://d2lang.com/tour/themes/)
- [GitHub Actions Workflow](.github/workflows/deploy-pages.yml)
- [RUNE Diagram Style Guide](../../RUNE/docs/diagrams/)

# Pedantic Raven GitHub Pages Site

This directory contains the GitHub Pages website for Pedantic Raven.

## Structure

```
docs/site/
├── index.html          # Main landing page
├── css/
│   └── style.css       # Raven-themed stylesheet
├── js/
│   └── main.js         # Minimal JavaScript enhancements
├── images/             # SVG diagrams (generated from ../diagrams/*.d2)
│   ├── 01-system-context.svg
│   ├── 02-component-architecture.svg
│   ├── 03-semantic-pipeline.svg
│   ├── 04-entity-comparison.svg
│   ├── 05-triple-graph.svg
│   ├── 06-data-flow.svg
│   ├── 07-offline-state.svg
│   ├── 08-mode-switching.svg
│   ├── 09-connection-lifecycle.svg
│   └── 10-event-system.svg
└── README.md           # This file
```

## Deployment

The site is automatically deployed to GitHub Pages via GitHub Actions whenever changes are pushed to the `main` branch that affect:
- `docs/site/**`
- `docs/diagrams/**`
- `docs/WHITEPAPER.md`

See `.github/workflows/deploy-pages.yml` for the deployment workflow.

## Local Development

To test the site locally:

```bash
# Option 1: Python HTTP server
cd docs/site
python -m http.server 8000
# Visit: http://localhost:8000

# Option 2: Node.js http-server
npm install -g http-server
cd docs/site
http-server -p 8000
# Visit: http://localhost:8000
```

## Diagram Generation

Diagrams are generated from D2 source files in `docs/diagrams/`:

```bash
# Render all diagrams
for diagram in docs/diagrams/*.d2; do
  basename=$(basename "$diagram" .d2)
  d2 "$diagram" "docs/site/images/${basename}.svg"
done

# Or use the provided script (future)
make diagrams
```

## Design

The site uses a Raven-themed color palette inspired by the Pedantic Raven name:

- **Raven Black** (#0A0E27): Primary dark color
- **Raven Teal** (#16A085): Primary accent color
- **Raven Amber** (#F39C12): Current/warning states
- **Raven Purple** (#9B59B6): Special/analysis features

Design is based on the mnemosyne GitHub Pages site with adaptations for Pedantic Raven's identity.

## Fonts

- **Inter**: Sans-serif for body text
- **JetBrains Mono**: Monospace for code and technical elements

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- Responsive design for mobile/tablet
- Graceful degradation for older browsers

## License

MIT License - See LICENSE file in repository root.

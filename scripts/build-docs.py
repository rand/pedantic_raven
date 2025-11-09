#!/usr/bin/env python3
"""Custom documentation build system for Pedantic Raven."""
import os, shutil
from pathlib import Path
import markdown
from jinja2 import Environment, FileSystemLoader

PROJECT_NAME, PROJECT_VERSION, PROJECT_TAGLINE = "Pedantic Raven", "v0.5.0", "// Context Engineering"
PROJECT_GLYPH, PROJECT_ACCENT_COLOR = "∴", "#DC2626"  # Therefore symbol, Red/Orange
GITHUB_URL, SITE_URL = "https://github.com/rand/pedantic_raven", "https://rand.github.io/pedantic_raven/"
BASE_DIR, DOCS_DIR = Path(__file__).parent.parent, Path(__file__).parent.parent / "docs"
TEMPLATES_DIR, SITE_DIR = BASE_DIR / "templates", BASE_DIR / "docs"
NAV_LINKS = [{"title": "Abstract", "href": "#abstract"}, {"title": "Architecture", "href": "#architecture"},
             {"title": "Features", "href": "#features"}, {"title": "Validation", "href": "#validation"},
             {"title": "Source", "href": f"{GITHUB_URL}/tree/400da43", "external": True}]

def setup_markdown():
    return markdown.Markdown(extensions=["extra", "codehilite", "toc", "sane_lists"],
        extension_configs={"codehilite": {"css_class": "highlight", "linenums": False}, "toc": {"permalink": False, "toc_depth": 3}})

def copy_static_files():
    # Skip copying if source and destination are the same
    if DOCS_DIR == SITE_DIR:
        print("  Static files already in place (source == destination)")
        return
    for dir_name in ["css", "js", "assets", "images"]:
        if (src_dir := DOCS_DIR / dir_name).exists():
            if (dest_dir := SITE_DIR / dir_name).exists(): shutil.rmtree(dest_dir)
            shutil.copytree(src_dir, dest_dir); print(f"  Copied {dir_name}/")
    for f in DOCS_DIR.glob("favicon*"): shutil.copy(f, SITE_DIR / f.name); print(f"  Copied {f.name}")

def strip_yaml_frontmatter(content):
    return content.split("---", 2)[2].strip() if content.startswith("---") and len(content.split("---", 2)) >= 3 else content

def render_page(env, md, tpl, md_file, out, ctx=None):
    if not (md_path := DOCS_DIR / md_file).exists(): return print(f"  ⚠️  Skip {md_file}")
    content = strip_yaml_frontmatter(md_path.read_text())
    html = md.convert(content); md.reset()
    context = {"project_name": PROJECT_NAME, "project_version": PROJECT_VERSION, "project_tagline": PROJECT_TAGLINE,
               "project_glyph": PROJECT_GLYPH, "github_url": GITHUB_URL, "site_url": SITE_URL,
               "nav_links": NAV_LINKS, "content": html, **(ctx or {})}
    (SITE_DIR / out).parent.mkdir(parents=True, exist_ok=True)
    (SITE_DIR / out).write_text(env.get_template(tpl).render(**context))
    print(f"  ✓ {md_file} → {out}")

def build():
    print(f"\n🔨 Building {PROJECT_NAME} documentation...\n")
    SITE_DIR.mkdir(parents=True, exist_ok=True)
    env, md = Environment(loader=FileSystemLoader(str(TEMPLATES_DIR))), setup_markdown()
    print("Rendering pages:")
    render_page(env, md, "index.html", "index.md", "index.html")
    print("\nCopying static assets:"); copy_static_files()
    (SITE_DIR / ".nojekyll").touch(); print("  ✓ Created .nojekyll")
    print(f"\n✅ Build complete! Site in: {SITE_DIR}\n")

if __name__ == "__main__": build()

package syntax

import (
	"path/filepath"
	"strings"
)

// DetectLanguage detects the programming language from a file path.
func DetectLanguage(path string) Language {
	if path == "" {
		return LangNone
	}

	ext := strings.ToLower(filepath.Ext(path))
	base := strings.ToLower(filepath.Base(path))

	// Check by extension
	switch ext {
	case ".go":
		return LangGo
	case ".md", ".markdown":
		return LangMarkdown
	case ".py":
		return LangPython
	case ".js", ".jsx":
		return LangJavaScript
	case ".ts", ".tsx":
		return LangTypeScript
	case ".rs":
		return LangRust
	case ".json":
		return LangJSON
	}

	// Check by filename
	switch base {
	case "go.mod", "go.sum":
		return LangGo
	case "readme", "readme.txt":
		return LangMarkdown
	case "package.json", "tsconfig.json":
		return LangJSON
	}

	return LangNone
}

// DetectLanguageFromContent attempts to detect language from file content.
// This is a fallback when file path detection fails.
func DetectLanguageFromContent(content string) Language {
	if content == "" {
		return LangNone
	}

	// Check for shebangs
	if strings.HasPrefix(content, "#!") {
		if strings.Contains(content, "python") {
			return LangPython
		}
		if strings.Contains(content, "node") || strings.Contains(content, "javascript") {
			return LangJavaScript
		}
	}

	// Check for package declarations
	if strings.Contains(content, "package ") && strings.Contains(content, "func ") {
		return LangGo
	}

	// Check for Markdown indicators
	if strings.Contains(content, "# ") || strings.Contains(content, "## ") || strings.Contains(content, "```") {
		return LangMarkdown
	}

	// Check for JSON
	trimmed := strings.TrimSpace(content)
	if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
		return LangJSON
	}

	return LangNone
}

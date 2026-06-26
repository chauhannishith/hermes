package hermes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// ParseStylesDefinition parses a raw CSS string into a StylesDefinition.
// It supports a very small subset of CSS adequate for simple selector { prop: value; } rules.
// - Removes standalone comments (/* ... */) but preserves inline comments with selectors
// - Does not support nested rules, media queries, or at-rules (they should be injected separately)
// - Multiple selectors separated by commas are split and each receives the full property set
// Consumers can use this to transform their custom CSS overrides into a StylesDefinition for merging.
func ParseStylesDefinition(css string) StylesDefinition {
	styles := StylesDefinition{}

	lines := strings.Split(css, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "/*") && strings.HasSuffix(trimmed, "*/") && !strings.Contains(trimmed, "{") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	cleanedCSS := strings.Join(cleanedLines, "\n")

	blockRE := regexp.MustCompile(`(?s)([^{}]+)\{([^{}]+)\}`)
	matches := blockRE.FindAllStringSubmatch(cleanedCSS, -1)
	for _, m := range matches {
		selectorPart := strings.TrimSpace(m[1])
		propsPart := strings.TrimSpace(m[2])
		if selectorPart == "" || propsPart == "" {
			continue
		}
		selectors := strings.Split(selectorPart, ",")
		commentRE := regexp.MustCompile(`(?s)/\*.*?\*/`)
		propsPartNoComments := commentRE.ReplaceAllString(propsPart, "")
		decls := strings.Split(propsPartNoComments, ";")
		props := map[string]any{}
		for _, d := range decls {
			d = strings.TrimSpace(d)
			if d == "" {
				continue
			}
			if colon := strings.Index(d, ":"); colon != -1 {
				key := strings.TrimSpace(d[:colon])
				val := strings.TrimSpace(d[colon+1:])
				if key != "" && val != "" {
					props[key] = val
				}
			}
		}
		if len(props) == 0 {
			continue
		}
		for _, sel := range selectors {
			s := strings.TrimSpace(sel)
			if s == "" {
				continue
			}
			if existing, ok := styles[s]; ok {
				for k, v := range props {
					existing[k] = v
				}
			} else {
				cp := map[string]any{}
				for k, v := range props {
					cp[k] = v
				}
				styles[s] = cp
			}
		}
	}
	return styles
}

// loadCustomCSS reads CSS from a URL, file path, or treats the input as inline CSS.
func loadCustomCSS(source string) (string, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return "", fmt.Errorf("hermes: fetch custom CSS: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("hermes: fetch custom CSS: status %d", resp.StatusCode)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("hermes: read custom CSS: %w", err)
		}
		return string(b), nil
	}

	if b, err := os.ReadFile(source); err == nil {
		return string(b), nil
	}

	return source, nil
}

func mergeStylesWithTheme(overrides StylesDefinition, theme Theme) StylesDefinition {
	themeStyles := theme.Styles()
	for sel, props := range overrides {
		if defProps, exists := themeStyles[sel]; exists {
			for k, v := range props {
				defProps[k] = v
			}
			themeStyles[sel] = defProps
		} else {
			themeStyles[sel] = props
		}
	}
	return themeStyles
}

func resolveEmailStyles(theme Theme, customCSS string) (StylesDefinition, error) {
	if customCSS == "" {
		return theme.Styles(), nil
	}

	raw, err := loadCustomCSS(customCSS)
	if err != nil {
		return nil, err
	}

	return mergeStylesWithTheme(ParseStylesDefinition(raw), theme), nil
}

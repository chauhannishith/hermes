package hermes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	errCustomCSSFetchStatus = errors.New("hermes: fetch custom CSS unexpected status")
	errCustomCSSRead        = errors.New("hermes: read custom CSS")
)

// ParseStylesDefinition parses a raw CSS string into a StylesDefinition.
// It supports a very small subset of CSS adequate for simple selector { prop: value; } rules.
func ParseStylesDefinition(css string) StylesDefinition {
	styles := StylesDefinition{}
	cleanedCSS := stripStandaloneCSSComments(css)

	blockRE := regexp.MustCompile(`(?s)([^{}]+)\{([^{}]+)\}`)
	matches := blockRE.FindAllStringSubmatch(cleanedCSS, -1)

	for _, m := range matches {
		mergeCSSBlock(styles, m[1], m[2])
	}

	return styles
}

func stripStandaloneCSSComments(css string) string {
	lines := strings.Split(css, "\n")
	cleanedLines := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "/*") && strings.HasSuffix(trimmed, "*/") && !strings.Contains(trimmed, "{") {
			continue
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}

func mergeCSSBlock(styles StylesDefinition, selectorPart, propsPart string) {
	selectorPart = strings.TrimSpace(selectorPart)
	propsPart = strings.TrimSpace(propsPart)
	if selectorPart == "" || propsPart == "" {
		return
	}

	props := parseCSSProperties(propsPart)
	if len(props) == 0 {
		return
	}

	for _, sel := range strings.Split(selectorPart, ",") {
		s := strings.TrimSpace(sel)
		if s == "" {
			continue
		}

		if existing, ok := styles[s]; ok {
			maps.Copy(existing, props)

			continue
		}

		cp := map[string]any{}
		maps.Copy(cp, props)
		styles[s] = cp
	}
}

func parseCSSProperties(propsPart string) map[string]any {
	commentRE := regexp.MustCompile(`(?s)/\*.*?\*/`)
	propsPartNoComments := commentRE.ReplaceAllString(propsPart, "")
	props := map[string]any{}

	for _, d := range strings.Split(propsPartNoComments, ";") {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}

		key, val, ok := strings.Cut(d, ":")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if key != "" && val != "" {
			props[key] = val
		}
	}

	return props
}

// loadCustomCSS reads CSS from a URL, file path, or treats the input as inline CSS.
func loadCustomCSS(source string) (string, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return loadCustomCSSFromURL(source)
	}

	if b, err := os.ReadFile(source); err == nil { //nolint:gosec // user-provided theme override path
		return string(b), nil
	}

	return source, nil
}

func loadCustomCSSFromURL(source string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, source, nil)
	if err != nil {
		return "", fmt.Errorf("hermes: fetch custom CSS: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("hermes: fetch custom CSS: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d", errCustomCSSFetchStatus, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errCustomCSSRead, err)
	}

	return string(b), nil
}

func mergeStylesWithTheme(overrides StylesDefinition, theme Theme) StylesDefinition {
	themeStyles := theme.Styles()

	for sel, props := range overrides {
		if defProps, exists := themeStyles[sel]; exists {
			maps.Copy(defProps, props)
			themeStyles[sel] = defProps

			continue
		}

		themeStyles[sel] = props
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

const defaultMediaQueries = `
@media only screen and (max-width: 600px) {
  .email-body_inner,
  .email-footer {
    width: 100% !important;
  }
}

@media only screen and (max-width: 500px) {
  .button {
    width: 100% !important;
  }
}
`

func renderStylesCSS(styles StylesDefinition) string {
	var b strings.Builder

	for selector, props := range styles {
		b.WriteString(selector)
		b.WriteString(" {\n")

		for key, val := range props {
			fmt.Fprintf(&b, "  %s: %v;\n", key, val)
		}

		b.WriteString("}\n\n")
	}

	b.WriteString(defaultMediaQueries)

	return b.String()
}

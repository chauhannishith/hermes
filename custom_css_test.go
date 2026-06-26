package hermes

import (
	"path/filepath"
	"strings"
	"testing"
)

const exampleCustomCSS = "examples/custom.css"

func testHermes() *Hermes {
	return &Hermes{
		Product:            Product{Name: "Test", Link: "https://example.com"},
		DisableCSSInlining: true,
	}
}

func TestDefaultThemeWithoutCustomCSS(t *testing.T) {
	html, err := testHermes().GenerateHTML(Email{Body: Body{Intros: []string{"Hi"}}})
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	lower := strings.ToLower(html)
	for _, want := range []string{
		"background-color: #f2f4f6",
		"background-color: #3869d4",
		"font-family:",
	} {
		if !strings.Contains(lower, want) {
			t.Errorf("expected default theme to contain %q", want)
		}
	}
}

func TestCustomCSSFromExampleFile(t *testing.T) {
	cssPath, err := filepath.Abs(exampleCustomCSS)
	if err != nil {
		t.Fatalf("resolve css path: %v", err)
	}

	html, err := testHermes().GenerateHTML(Email{Body: Body{
		Intros:    []string{"Hi"},
		CustomCSS: cssPath,
	}})
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	lower := strings.ToLower(html)

	// Overrides from examples/custom.css
	for _, want := range []string{
		"background-color: #ede7f6",
		"color: #4a148c",
		"color: #6a1b9a",
		"background-color: #7b1fa2",
	} {
		if !strings.Contains(lower, want) {
			t.Errorf("expected override from %s to contain %q", exampleCustomCSS, want)
		}
	}

	// Default styles for selectors not in custom.css
	for _, want := range []string{
		"background-color: #fff",
		"font-family:",
		".email-body_inner",
	} {
		if !strings.Contains(lower, want) {
			t.Errorf("expected default theme to still contain %q", want)
		}
	}

	// Original default values replaced where overridden
	if strings.Contains(lower, "background-color: #3869d4") {
		t.Error("expected .button default color to be overridden by custom.css")
	}
}

func TestCustomCSSPartialPropertyMerge(t *testing.T) {
	html, err := testHermes().GenerateHTML(Email{Body: Body{
		Intros:    []string{"Hi"},
		CustomCSS: `.button { background-color: #FF5722; }`,
	}})
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	lower := strings.ToLower(html)
	if !strings.Contains(lower, "background-color: #ff5722") {
		t.Error("expected button background override")
	}
	if !strings.Contains(lower, "border-radius: 3px") {
		t.Error("expected other .button properties from default theme to remain")
	}
}

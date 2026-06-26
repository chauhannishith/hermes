package hermes

import (
	"strings"
	"testing"
)

func TestParseStylesDefinitionBasic(t *testing.T) {
	css := `body { color: #111; background-color: #fff; }
.a, .b { font-size: 14px; }`

	styles := ParseStylesDefinition(css)

	if styles["body"]["color"] != "#111" {
		t.Fatalf("expected body color, got %#v", styles["body"])
	}
	if styles[".a"]["font-size"] != "14px" || styles[".b"]["font-size"] != "14px" {
		t.Fatalf("expected shared font-size for .a and .b, got %#v", styles)
	}
}

func TestRenderStylesCSSIncludesMediaQueries(t *testing.T) {
	css := renderStylesCSS(StylesDefinition{
		"body": {"color": "#111"},
	})
	if !strings.Contains(css, "body {\n  color: #111;\n}") {
		t.Fatal("expected selector block in rendered CSS")
	}
	if !strings.Contains(css, "@media only screen and (max-width: 600px)") {
		t.Fatal("expected responsive media queries in rendered CSS")
	}
}

func TestParseStylesDefinitionIgnoresInvalid(t *testing.T) {
	css := `div { color: red; }
span { missing colon }`

	styles := ParseStylesDefinition(css)
	if styles["div"]["color"] != "red" {
		t.Fatal("expected div color red")
	}
	if _, ok := styles["span"]; ok {
		t.Fatal("expected malformed span rule to be ignored")
	}
}

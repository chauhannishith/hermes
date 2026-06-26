package hermes

import "testing"

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

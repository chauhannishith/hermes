package hermes

import (
	"embed"
	"fmt"
)

//go:embed templates
var staticFS embed.FS

const (
	htmlEmail      = "templates/%s.tpl.html"
	plainTextEmail = "templates/%s.tpl.txt"
)

func getTemplate(name string) string {
	htmlBytes, err := staticFS.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("hermes: read template %s: %v", name, err))
	}

	return string(htmlBytes)
}

// GetDefaultStyles returns the parsed CSS styles for the default theme.
func GetDefaultStyles() StylesDefinition {
	cssBytes, err := staticFS.ReadFile("templates/default.css")
	if err != nil {
		panic(fmt.Sprintf("hermes: read default.css: %v", err))
	}

	return ParseStylesDefinition(string(cssBytes))
}

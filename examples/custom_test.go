package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matcornic/hermes"
)

func TestGenerateCustomCSSExamples(t *testing.T) {
	h := hermes.Hermes{
		DefaultGreeting: "Hi",
		Theme:           new(hermes.Default),
		Product: hermes.Product{
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			Logo: "https://github.com/matcornic/hermes/blob/master/examples/gopher.png?raw=true",
		},
	}

	cssPath, err := filepath.Abs("custom.css")
	if err != nil {
		t.Fatalf("resolve custom.css: %v", err)
	}

	examples := []example{
		new(welcome),
		new(reset),
		new(receipt),
		new(maintenance),
		new(inviteCode),
	}

	const outDir = "custom"

	for _, e := range examples {
		email := e.Email()
		email.Body.CustomCSS = cssPath

		html, err := h.GenerateHTML(email)
		if err != nil {
			t.Fatalf("GenerateHTML(%s): %v", e.Name(), err)
		}

		txt, err := h.GeneratePlainText(email)
		if err != nil {
			t.Fatalf("GeneratePlainText(%s): %v", e.Name(), err)
		}

		if err := os.MkdirAll(outDir, 0750); err != nil {
			t.Fatalf("mkdir %s: %v", outDir, err)
		}

		htmlFile := fmt.Sprintf("%s/%s.%s.html", outDir, outDir, e.Name())
		if err := os.WriteFile(htmlFile, []byte(html), 0600); err != nil {
			t.Fatalf("write %s: %v", htmlFile, err)
		}

		txtFile := fmt.Sprintf("%s/%s.%s.txt", outDir, outDir, e.Name())
		if err := os.WriteFile(txtFile, []byte(txt), 0600); err != nil {
			t.Fatalf("write %s: %v", txtFile, err)
		}

		if !strings.Contains(strings.ToLower(html), "#f3e5f5") {
			t.Errorf("%s: expected custom.css email-body background in generated HTML", htmlFile)
		}

		if e.Name() == "welcome" || e.Name() == "reset" {
			lower := strings.ToLower(html)
			if !strings.Contains(lower, "border-radius:24px") && !strings.Contains(lower, "border-radius: 24px") {
				t.Errorf("%s: expected custom.css pill button style in generated HTML", htmlFile)
			}
		}

		if e.Name() == "invite_code" && !strings.Contains(strings.ToLower(html), "#e1bee7") {
			t.Errorf("%s: expected custom.css invite-code background in generated HTML", htmlFile)
		}
	}
}

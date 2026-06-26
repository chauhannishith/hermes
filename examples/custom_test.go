package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chauhannishith/hermes"
)

func TestGenerateCustomCSSExamples(t *testing.T) {
	h := hermes.Hermes{
		DefaultGreeting: "Hi",
		Theme:           new(hermes.Default),
		Product: hermes.Product{
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			Logo: "https://github.com/chauhannishith/hermes/blob/main/examples/gopher.png?raw=true",
		},
	}

	cssPath, err := filepath.Abs("custom.css")
	if err != nil {
		t.Fatalf("resolve custom.css: %v", err)
	}

	for _, e := range []example{
		new(welcome),
		new(reset),
		new(receipt),
		new(maintenance),
		new(inviteCode),
	} {
		writeCustomExample(t, h, e, cssPath)
	}
}

func writeCustomExample(t *testing.T, h hermes.Hermes, e example, cssPath string) {
	t.Helper()

	const outDir = "custom"

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

	err = os.MkdirAll(outDir, 0750)
	if err != nil {
		t.Fatalf("mkdir %s: %v", outDir, err)
	}

	htmlFile := fmt.Sprintf("%s/%s.%s.html", outDir, outDir, e.Name())
	err = os.WriteFile(htmlFile, []byte(html), 0600)
	if err != nil {
		t.Fatalf("write %s: %v", htmlFile, err)
	}

	txtFile := fmt.Sprintf("%s/%s.%s.txt", outDir, outDir, e.Name())
	err = os.WriteFile(txtFile, []byte(txt), 0600)
	if err != nil {
		t.Fatalf("write %s: %v", txtFile, err)
	}

	assertCustomCSSApplied(t, htmlFile, e.Name(), html)
}

func assertCustomCSSApplied(t *testing.T, htmlFile, name, html string) {
	t.Helper()

	lower := strings.ToLower(html)
	if !strings.Contains(lower, "#f3e5f5") {
		t.Errorf("%s: expected custom.css email-body background in generated HTML", htmlFile)
	}

	if name == "welcome" || name == "reset" {
		if !strings.Contains(lower, "border-radius:24px") && !strings.Contains(lower, "border-radius: 24px") {
			t.Errorf("%s: expected custom.css pill button style in generated HTML", htmlFile)
		}
	}

	if name == "invite_code" && !strings.Contains(lower, "#e1bee7") {
		t.Errorf("%s: expected custom.css invite-code background in generated HTML", htmlFile)
	}
}

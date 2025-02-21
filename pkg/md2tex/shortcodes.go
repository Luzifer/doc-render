package md2tex

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

var (
	shortCodes = template.FuncMap{
		"graphic": shortCodeGraphic,
		"part":    shortCodePart,
		"raw":     shortCodeRaw,
		"vspace":  shortCodeVSpace,
	}
	shortCodeDef = regexp.MustCompile(`{% (.*?) %}`)
)

func renderShortCode(content string) (string, error) {
	tpl, err := template.New("shortCode").
		Funcs(shortCodes).
		Parse(fmt.Sprintf(`{{- %s -}}`, content))
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, nil); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func shortCodeGraphic(path string) string {
	return fmt.Sprintf(`\includegraphics{%s}`, path)
}

func shortCodePart(title string) string {
	return fmt.Sprintf(`\part{%s}`, title)
}

func shortCodeRaw(latex string) string { return latex }

func shortCodeVSpace(dist string) string {
	return fmt.Sprintf(`\vspace{%s}`, dist)
}

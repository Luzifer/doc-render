package md2tex

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	shortCodes = map[string]shortCode{
		"graphic": shortCodeGraphic,
		"part":    shortCodePart,
		"vspace":  shortCodeVSpace,
	}
	shortCodeDef = regexp.MustCompile(`{% (.*?) %}`)
)

func shortCodeGraphic(args []string) (string, error) {
	return fmt.Sprintf(`\includegraphics{%s}`, args[1]), nil
}

func shortCodePart(args []string) (string, error) {
	return fmt.Sprintf(`\part{%s}`, strings.Join(args[1:], " ")), nil
}

func shortCodeVSpace(args []string) (string, error) {
	return fmt.Sprintf(`\vspace{%s}`, args[1]), nil
}

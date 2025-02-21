// Package md2tex wraps Goldmark to provide an easy-to-use API to
// convert Markdown documents to LaTeX.
package md2tex

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const prio = 1000

// Convert takes a Markdown document and returns LaTex source from it
func Convert(md []byte) (tex []byte, err error) {
	rd := renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(newGenerator(), prio)))
	gm := goldmark.New(
		goldmark.WithRenderer(rd),
	)

	output := new(bytes.Buffer)
	if err = gm.Convert(md, output); err != nil {
		return nil, fmt.Errorf("rendering TeX: %w", err)
	}

	return bytes.TrimSpace(output.Bytes()), nil
}

package md2tex

//revive:disable:flag-parameter // Third-party API, no way to change

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type (
	generator struct{}
)

var _ renderer.NodeRenderer = &generator{}

func newGenerator() *generator {
	return &generator{}
}

func (g generator) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, g.renderDocument)
	reg.Register(ast.KindHeading, g.renderHeading)
	reg.Register(ast.KindBlockquote, g.renderBlockquote)
	reg.Register(ast.KindCodeBlock, g.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, g.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, g.renderHTMLBlock)
	reg.Register(ast.KindList, g.renderList)
	reg.Register(ast.KindListItem, g.renderListItem)
	reg.Register(ast.KindParagraph, g.renderParagraph)
	reg.Register(ast.KindTextBlock, g.renderTextBlock)
	reg.Register(ast.KindThematicBreak, g.renderThematicBreak)

	// inlines
	reg.Register(ast.KindAutoLink, g.renderAutoLink)
	reg.Register(ast.KindCodeSpan, g.renderCodeSpan)
	reg.Register(ast.KindEmphasis, g.renderEmphasis)
	reg.Register(ast.KindImage, g.renderImage)
	reg.Register(ast.KindLink, g.renderLink)
	reg.Register(ast.KindRawHTML, g.renderRawHTML)
	reg.Register(ast.KindText, g.renderText)
	reg.Register(ast.KindString, g.renderString)
}

func (generator) escapeLaTeX(data []byte) []byte {
	buf := new(bytes.Buffer)
	for _, b := range data {
		switch b {
		case '\\':
			_, _ = buf.WriteString("\\textbackslash~")

		case '~':
			_, _ = buf.WriteString("\\textasciitilde~")

		case '^':
			_, _ = buf.WriteString("\\textasciicircum~")

		case '&', '%', '$', '#', '_', '{', '}':
			_, _ = buf.Write([]byte{'\\', b})

		default:
			_ = buf.WriteByte(b)
		}
	}

	return buf.Bytes()
}

func (generator) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	var (
		n   = node.(*ast.AutoLink)
		url = n.URL(source)
	)

	if n.AutoLinkType == ast.AutoLinkEmail {
		url = bytes.Join([][]byte{[]byte("mailto"), url}, []byte{':'})
	}

	if _, err := w.WriteString(fmt.Sprintf("\\href{%s}{%s}", url, n.Label(source))); err != nil {
		return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
	}

	return ast.WalkContinue, nil
}

func (generator) renderBlockquote(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("\\begin{framed}\n\\begin{quote}\n"); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if _, err := w.WriteString("\\end{quote}\n\\end{framed}\n\n"); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("\\begin{lstlisting}\n"); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}

		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			if _, err := w.Write(line.Value(source)); err != nil {
				return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
			}
		}
	} else {
		if _, err := w.WriteString("\\end{lstlisting}\n\n"); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderCodeSpan(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("\\texttt{"); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if err := w.WriteByte('}'); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderDocument(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	// nothing to do
	return ast.WalkContinue, nil
}

func (generator) renderEmphasis(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var (
		n       = node.(*ast.Emphasis)
		command = "textit"
	)

	if n.Level == 2 { //nolint:mnd
		command = "textbf"
	}

	if entering {
		if _, err := w.WriteString(fmt.Sprintf("\\%s{", command)); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if err := w.WriteByte('}'); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (g generator) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return g.renderCodeBlock(w, source, node, entering)
}

func (generator) renderHeading(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var (
		n        = node.(*ast.Heading)
		headings = []string{
			"",               // SKIP
			"section*",       // H1
			"subsection*",    // H2
			"subsubsection*", // H3
			"paragraph*",     // H4
			"subparagraph*",  // H5
			"textbf",         // H6
		}
	)

	if entering {
		if _, err := w.WriteString(fmt.Sprintf("\\%s{", headings[n.Level])); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if _, err := w.WriteString("}\n\n"); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderHTMLBlock(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, fmt.Errorf("unsupported: HTMLBlock")
}

func (generator) renderImage(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, fmt.Errorf("unsupported: Image")
}

func (generator) renderLink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)

	if entering {
		if _, err := w.WriteString(fmt.Sprintf("\\href{%s}{", n.Destination)); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if err := w.WriteByte('}'); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderList(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var (
		n        = node.(*ast.List)
		listType = "itemize"
	)

	if n.IsOrdered() {
		listType = "enumerate"
	}

	if entering {
		if _, err := w.WriteString(fmt.Sprintf("\\begin{%s}\n", listType)); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if _, err := w.WriteString(fmt.Sprintf("\\end{%s}", listType)); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}

		if _, ok := n.Parent().(*ast.Document); ok {
			// Prevent double line breaks on sub-lists
			if _, err := w.Write([]byte{'\n', '\n'}); err != nil {
				return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
			}
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderListItem(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if _, err := w.WriteString("\\item "); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else {
		if err := w.WriteByte('\n'); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderParagraph(w util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, err := w.Write([]byte("\n\n")); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderRawHTML(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, fmt.Errorf("unsupported: RawHTML")
}

func (generator) renderString(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, fmt.Errorf("unsupported: String")
}

func (g generator) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)

	if _, err := w.Write(g.replaceShortCode(n.Segment.Value(source))); err != nil {
		return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
	}

	if n.HardLineBreak() {
		if _, err := w.Write([]byte("\\\\\n")); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	} else if n.SoftLineBreak() {
		if err := w.WriteByte(' '); err != nil {
			return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderTextBlock(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if node.NextSibling() != nil && node.FirstChild() != nil {
			if err := w.WriteByte('\n'); err != nil {
				return ast.WalkStop, fmt.Errorf("writing TeX: %w", err)
			}
		}
	}

	return ast.WalkContinue, nil
}

func (generator) renderThematicBreak(util.BufWriter, []byte, ast.Node, bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, fmt.Errorf("unsupported: ThematicBreak")
}

func (g generator) replaceShortCode(value []byte) []byte {
	if !shortCodeDef.Match(value) {
		return g.escapeLaTeX(value)
	}

	match := shortCodeDef.FindSubmatch(value)

	args := strings.Fields(string(match[1]))
	sc, ok := shortCodes[args[0]]
	if !ok {
		return value
	}

	repl, err := sc(args)
	if err != nil {
		return []byte(fmt.Sprintf("%% Shortcode %q error: %s", args[0], err))
	}

	return []byte(repl)
}

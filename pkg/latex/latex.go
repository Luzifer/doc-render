// Package latex contains a rendering helper for the letter source
// using a hosted TeX-API instance
package latex

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/Luzifer/doc-render/pkg/recipientcsv"
	"github.com/sirupsen/logrus"
)

type (
	// RenderOpts define what to render into the template
	RenderOpts struct {
		// Instance of the TeX-API to use for rendering
		TexAPIURL string

		// Folder containing the source-sets
		SourceBaseFolder string
		// Source-set to include in the zip
		SourceSet string

		// Recipients contains the recipients for the letter (might not be
		// supported by the chosen template)
		Recipients []recipientcsv.Person
		// Values to be used in the template (usage depends on the template)
		Values any
	}
)

// Render takes the options and the included template / source files,
// generate the TeX document and renders it through the provided API.
//
// The returned io.ReadCloser MUST be closed after usage to free up resources.
func Render(ctx context.Context, opts RenderOpts) (pdf io.ReadCloser, err error) {
	sourceFiles := os.DirFS(path.Join(opts.SourceBaseFolder, opts.SourceSet))

	tpl, err := readTemplate(sourceFiles, "main.tex.tpl")
	if err != nil {
		return nil, fmt.Errorf("reading template: %w", err)
	}

	// Prepare a ZIP to upload to the API
	zipFile := new(bytes.Buffer)
	if err = packSource(zipFile, sourceFiles, tpl, opts); err != nil {
		return nil, fmt.Errorf("building ZIP: %w", err)
	}

	if pdf, err = renderDocument(ctx, opts, zipFile); err != nil {
		return nil, fmt.Errorf("rendering PDF: %w", err)
	}

	return pdf, nil
}

func packSource(dst io.Writer, sourceFiles fs.FS, tpl *template.Template, opts RenderOpts) (err error) {
	zw := zip.NewWriter(dst)

	// Add all files from the source (including the template which will
	// not be used by the TeX-API as of the .tpl suffix)
	if err = zw.AddFS(sourceFiles); err != nil {
		return fmt.Errorf("adding source-files: %w", err)
	}

	// Add the TeX document
	texFile, err := zw.Create("main.tex")
	if err != nil {
		return fmt.Errorf("creating main.tex: %w", err)
	}

	if err = tpl.Execute(texFile, opts); err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	// Close and finalize archive
	if err = zw.Close(); err != nil {
		return fmt.Errorf("closing archive: %w", err)
	}

	return nil
}

func readTemplate(src fs.FS, name string) (*template.Template, error) {
	f, err := src.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening template file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.WithError(err).Error("closing template file")
		}
	}()

	tplSource, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading template: %w", err)
	}

	tpl, err := template.New("letter").Funcs(templateFuncs()).Parse(string(tplSource))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	return tpl, nil
}

func renderDocument(ctx context.Context, opts RenderOpts, zipFile io.Reader) (pdf io.ReadCloser, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, opts.TexAPIURL, zipFile)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/pdf")
	req.Header.Set("Content-Type", "application/zip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("reading body: %s", err))
		}
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %d (%s)", resp.StatusCode, body)
	}

	return resp.Body, nil
}

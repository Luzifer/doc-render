package latex

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/invopop/jsonschema"
	"github.com/sirupsen/logrus"
)

// GetSourceSets returns all available source-sets and their definition
func GetSourceSets(base string) (schemas map[string]jsonschema.Schema, err error) {
	schemas = make(map[string]jsonschema.Schema)

	if err = filepath.WalkDir(base, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "schema.json" {
			return nil
		}

		f, err := os.Open(filePath) //#nosec:G304 // Intended to traverse custom path
		if err != nil {
			return fmt.Errorf("opening schema file: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				logrus.WithError(err).Error("closing schema file")
			}
		}()

		var s jsonschema.Schema
		if err = json.NewDecoder(f).Decode(&s); err != nil {
			return fmt.Errorf("parsing schema: %w", err)
		}

		schemas[path.Base(path.Dir(filePath))] = s

		return nil
	}); err != nil {
		return nil, fmt.Errorf("reading schemas: %w", err)
	}

	return schemas, nil
}

// HasSourceSet checks whether the given directory exists and contains
// at least the main template
func HasSourceSet(base, name string) bool {
	sourceFiles := os.DirFS(path.Join(base, name))

	f, err := sourceFiles.Open("main.tex.tpl")
	if err != nil {
		return false
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.WithError(err).Error("closing template file")
		}
	}()

	info, err := f.Stat()
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

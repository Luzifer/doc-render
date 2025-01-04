package latex

import (
	"text/template"

	"github.com/Luzifer/doc-render/pkg/md2tex"
	"github.com/Masterminds/sprig/v3"
)

func templateFuncs() template.FuncMap {
	fm := make(template.FuncMap)

	for fn, f := range sprig.FuncMap() {
		fm[fn] = f
	}

	fm["md2tex"] = func(s string) (string, error) {
		l, err := md2tex.Convert([]byte(s))
		return string(l), err
	}

	return fm
}

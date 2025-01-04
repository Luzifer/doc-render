package md2tex

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed test/input.md
	testMD []byte
	//go:embed test/output.tex
	testTeX []byte
)

func TestConvert(t *testing.T) {
	tex, err := Convert(testMD)
	require.NoError(t, err)

	assert.Equal(t, bytes.TrimSpace(testTeX), tex)
}

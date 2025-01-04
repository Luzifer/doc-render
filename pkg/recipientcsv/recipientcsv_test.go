package recipientcsv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	data := strings.TrimSpace(`
NACHNAME;VORNAME;STRASSE;HAUSNR;PLZ;ORT
Muster;Karl;Musterstraße;123;12345;Musterhausen
Muster;Birgit;Musterstraße;123;12345;Musterhausen
`)

	p, err := Parse(strings.NewReader(data))
	require.NoError(t, err)
	require.Len(t, p, 2)

	assert.Equal(t, Person{
		Lastname:     "Muster",
		Firstname:    "Karl",
		Street:       "Musterstraße",
		StreetNumber: "123",
		PostalCode:   "12345",
		City:         "Musterhausen",
	}, p[0])
}

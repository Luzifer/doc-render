// Package recipientcsv contains a method to parse an address-data CSV
package recipientcsv

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type (
	// Person represents a person the CSV
	Person struct {
		Lastname     string `json:"NACHNAME"`
		Firstname    string `json:"VORNAME"`
		Street       string `json:"STRASSE"`
		StreetNumber string `json:"HAUSNR"`
		PostalCode   string `json:"PLZ"`
		City         string `json:"ORT"`
	}
)

// Parse reads the given CSV and returns the Person data included
func Parse(data io.Reader) (out []Person, err error) {
	r := csv.NewReader(data)
	r.Comma = ';'
	r.FieldsPerRecord = -1

	// We need the headers to assemble the JSON object
	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("reading headers: %w", err)
	}

	// Now we walk through the lines and build the output
	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading record: %w", err)
		}

		obj := make(map[string]any)
		for i := range headers {
			obj[headers[i]] = record[i]
		}

		raw, err := json.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("marshalling JSON: %w", err)
		}

		var p Person
		if err = json.Unmarshal(raw, &p); err != nil {
			return nil, fmt.Errorf("unmarshalling JSON: %w", err)
		}

		out = append(out, p)
	}

	return out, nil
}

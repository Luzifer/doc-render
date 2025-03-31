// Package persist defines an interface to persist entered document
// contents server-side
package persist

type (
	// Backend defines the interface to implement when implementing a persist
	// backend
	Backend interface {
		// Get retrieves the JSON encoded template by its content-hash
		Get(uid string) (templateJSON []byte, err error)
		// Store takes the JSON encoded template, stores it and returns the
		// content-hash as uid and optionally an error
		Store(templateJSON []byte) (uid string, err error)
	}
)

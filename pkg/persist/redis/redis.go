// Package redis implements a storage backend to hold the templates
// inside a Redis cache
package redis

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/Luzifer/doc-render/pkg/persist"
	"github.com/redis/go-redis/v9"
)

type (
	// Backend implements the persist.Backend interface for Redis storage
	Backend struct {
		c *redis.Client
	}
)

var _ persist.Backend = (*Backend)(nil)

// New creates a new redis persistence backend
func New() (*Backend, error) {
	dsn := os.Getenv("PERSIST_REDIS")
	if dsn == "" {
		return nil, fmt.Errorf("no DSN set")
	}

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing redis DSN: %w", err)
	}

	return &Backend{redis.NewClient(opts)}, nil
}

// Get retrieves the JSON encoded template by its content-hash
func (b Backend) Get(uid string) (templateJSON []byte, err error) {
	if templateJSON, err = b.c.Get(context.Background(), b.storageKey(uid)).Bytes(); err != nil {
		return nil, fmt.Errorf("fetching template: %w", err)
	}

	return templateJSON, nil
}

// Store takes the JSON encoded template, stores it and returns the
// content-hash as uid and optionally an error
func (b Backend) Store(templateJSON []byte) (uid string, err error) {
	uid = b.contentHash(templateJSON)
	if err = b.c.Set(context.Background(), b.storageKey(uid), templateJSON, 0).Err(); err != nil {
		return "", fmt.Errorf("storing template: %w", err)
	}

	return uid, nil
}

func (Backend) contentHash(templateJSON []byte) string {
	return fmt.Sprintf("sha256:%x", sha256.Sum256(templateJSON))
}

func (Backend) storageKey(contentHash string) string {
	var parts []string

	if prefix := os.Getenv("PERSIST_REDIS_PREFIX"); prefix != "" {
		parts = append(parts, prefix)
	}

	return strings.Join(append(parts, contentHash), ":")
}

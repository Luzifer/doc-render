// Package k8s implements a storage backend to hold the templates
// inside a Kubernets cluster
package k8s

import (
	"context"
	"crypto/sha1" //#nosec G505: Used for content hash, not crypto
	"fmt"
	"os"

	"github.com/Luzifer/doc-render/pkg/persist"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyCoreV1 "k8s.io/client-go/applyconfigurations/core/v1"
	applyMetaV1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type (
	// Backend implements the persist.Backend interface for Memory storage
	Backend struct {
		c *kubernetes.Clientset
	}
)

var _ persist.Backend = (*Backend)(nil)

// New creates a new k8s persistence backend
func New() (*Backend, error) {
	if os.Getenv("PERSIST_NAMESPACE") == "" {
		return nil, fmt.Errorf("no PERSIST_NAMESPACE set")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("getting in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}

	return &Backend{
		c: clientset,
	}, nil
}

// Get retrieves the JSON encoded template by its content-hash
func (b *Backend) Get(uid string) (templateJSON []byte, err error) {
	cm, err := b.c.CoreV1().
		ConfigMaps(os.Getenv("PERSIST_NAMESPACE")).
		Get(context.Background(), b.cmName(uid), metaV1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting config-map: %w", err)
	}

	return []byte(cm.Data["template"]), nil
}

// Store takes the JSON encoded template, stores it and returns the
// content-hash as uid and optionally an error
func (b *Backend) Store(templateJSON []byte) (uid string, err error) {
	uid = b.contentHash(templateJSON)

	if _, err = b.c.CoreV1().
		ConfigMaps(os.Getenv("PERSIST_NAMESPACE")).
		Apply(context.Background(), &applyCoreV1.ConfigMapApplyConfiguration{
			ObjectMetaApplyConfiguration: &applyMetaV1.ObjectMetaApplyConfiguration{
				Name:      ptr(b.cmName(uid)),
				Namespace: ptr(os.Getenv("PERSIST_NAMESPACE")),
			},
			Data: map[string]string{
				"template": string(templateJSON),
			},
		}, metaV1.ApplyOptions{
			FieldManager: "application/apply-patch",
		}); err != nil {
		return "", fmt.Errorf("creating config-map: %w", err)
	}

	return uid, nil
}

func (Backend) cmName(uid string) string {
	return fmt.Sprintf("tpl-%s", uid)
}

func (Backend) contentHash(templateJSON []byte) string {
	// K8s resource names are limited to 63 chars, sha256 is 65 chars so
	// we use sha1-hashes in this backend
	return fmt.Sprintf("%x", sha1.Sum(templateJSON)) //#nosec G401: Used for content hash, not crypto
}

func ptr[T comparable](v T) *T { return &v }

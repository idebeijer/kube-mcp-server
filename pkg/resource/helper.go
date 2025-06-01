package resource

import (
	"errors"
	"strings"
)

// ExtractNamespaceFromURI extracts the namespace from a URI in the format k8s://{namespace}.
func ExtractNamespaceFromURI(uri string) (string, error) {
	if !strings.HasPrefix(uri, "k8s://") {
		return "", errors.New("invalid URI format")
	}

	// Remove the k8s:// prefix
	trimmedURI := strings.TrimPrefix(uri, "k8s://")
	parts := strings.Split(trimmedURI, "/")
	if len(parts) < 1 || parts[0] == "" {
		return "", errors.New("namespace not found in URI")
	}

	return parts[0], nil
}

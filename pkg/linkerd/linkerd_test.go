package linkerd

import (
	"testing"
)

// TestPackageExists tests that the linkerd package can be imported
func TestPackageExists(t *testing.T) {
	// This test simply verifies the package compiles and imports correctly
	// More comprehensive tests would require a real Kubernetes cluster
	t.Log("Linkerd package imported successfully")
}

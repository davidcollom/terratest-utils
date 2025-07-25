// Package velero provides Terratest-style helpers for testing Velero backups,
// restores, and storage configurations. Helpers include wait functions for
// BackupStorageLocations, Backups, and Restores using status conditions and phases.
package velero

import (
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewVeleroClient creates and returns a new controller-runtime client for interacting with Velero resources.
// It initializes a new runtime scheme, adds the Velero v1 API to the scheme, and constructs the client
// using the provided Kubernetes REST configuration.
//
// Parameters:
//   - cfg: A pointer to a rest.Config object containing the Kubernetes API server configuration.
//
// Returns:
//   - client.Client: A controller-runtime client configured for Velero resources.
//   - error: An error if the client could not be created.
func NewVeleroClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = velerov1.AddToScheme(scheme)
	return client.New(cfg, client.Options{Scheme: scheme})
}

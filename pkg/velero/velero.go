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

// NewVeleroClient creates a new controller-runtime client for Velero resources using the provided Kubernetes REST config.
func NewVeleroClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = velerov1.AddToScheme(scheme)
	return client.New(cfg, client.Options{Scheme: scheme})
}

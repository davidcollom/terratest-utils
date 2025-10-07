// Package cd provides Terratest-style helpers for testing Argo CD Applications and
// AppProjects. It includes functions to wait for Applications to become Synced and Healthy,
// as well as utilities to verify AppProject presence.
package cd

import (
	"testing"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	argocd "github.com/argoproj/argo-cd/v3/pkg/client/clientset/versioned"
	apphealth "github.com/argoproj/gitops-engine/pkg/health"
	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/client-go/rest"
)

// IsApplicationHealthyAndSynced returns true if the given Argo CD application is both healthy and synced.
// It checks that the application's health status is 'Healthy' and its sync status is 'Synced'.
func IsApplicationHealthyAndSynced(app *argocdv1alpha1.Application) bool {
	return app.Status.Health.Status == apphealth.HealthStatusHealthy &&
		app.Status.Sync.Status == argocdv1alpha1.SyncStatusCodeSynced
}

// NewArgoCDClient creates a new ArgoCD client interface for use in tests.
// This function attempts to use the REST configuration from options if available,
// otherwise it will create a new REST configuration using the kubectl options.
//
// Parameters:
//   - t: Testing context that implements the testing.TB interface
//   - options: Kubectl configuration options that may contain a REST client configuration
//
// Returns:
//   - An ArgoCD client interface that can be used to interact with ArgoCD
//   - An error if the client creation fails
var NewArgoCDClient = newArgoCDClient

func newArgoCDClient(t *testing.T, options *k8s.KubectlOptions) (argocd.Interface, error) {
	t.Helper()
	var cfg *rest.Config
	var err error
	if options.RestConfig == nil {
		cfg, err = utils.GetRestConfigE(t, options)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = options.RestConfig
	}

	return argocd.NewForConfig(cfg)
}

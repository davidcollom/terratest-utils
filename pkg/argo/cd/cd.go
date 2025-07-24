// Package cd provides Terratest-style helpers for testing Argo CD Applications and
// AppProjects. It includes functions to wait for Applications to become Synced and Healthy,
// as well as utilities to verify AppProject presence.
package cd

import (
	argocdv1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	apphealth "github.com/argoproj/gitops-engine/pkg/health"
)

// IsApplicationHealthyAndSynced returns true if the given Argo CD application is both healthy and synced.
// It checks that the application's health status is 'Healthy' and its sync status is 'Synced'.
func IsApplicationHealthyAndSynced(app *argocdv1alpha1.Application) bool {
	return app.Status.Health.Status == apphealth.HealthStatusHealthy &&
		app.Status.Sync.Status == argocdv1alpha1.SyncStatusCodeSynced
}

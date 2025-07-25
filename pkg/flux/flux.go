// Package flux provides Terratest-style helpers for testing Flux resources such as
// Kustomizations, HelmReleases, GitRepositories, and HelmRepositories. These functions
// wait for Flux CRDs to become Ready using status conditions and standard polling logic.
package flux

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

// hasReadyCondition checks if the provided slice of metav1.Condition contains a condition
// of type "Ready" with a status of metav1.ConditionTrue. It returns true if such a condition
// is found, otherwise it returns false.
//
// This function is commonly used to determine if a Kubernetes custom resource managed by Flux
// (such as a GitRepository, Kustomization, or HelmRelease) has reached the "Ready" state.
// It iterates through all conditions attached to the resource and returns true as soon as it
// finds a "Ready" condition with a status of "True". If no such condition is found, it returns false.
//
// Example usage:
//
//	ready := hasReadyCondition(resource.Status.Conditions)
//	if ready {
//	    // Resource is ready
//	}
func hasReadyCondition(conds []metav1.Condition) bool {
	for _, cond := range conds {
		if cond.Type == "Ready" && cond.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}

// NewFluxClient creates and returns a new controller-runtime client for interacting with Flux resources.
// It initializes a new runtime scheme, adds the Flux Kustomize, Helm, and Source controller APIs to the scheme,
// and constructs the client using the provided Kubernetes REST configuration.
//
// Parameters:
//   - t: The testing context.
//   - options: The KubectlOptions containing the Kubernetes REST config.
//
// Returns:
//   - client.Client: A controller-runtime client configured for Flux resources.
//   - error: An error if the client could not be created.
func NewFluxClient(t *testing.T, options *k8s.KubectlOptions) (client.Client, error) {
	cfg, err := utils.GetRestConfigE(t, options)

	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	_ = kustomizev1.AddToScheme(scheme)
	_ = helmv2.AddToScheme(scheme)
	_ = sourcev1.AddToScheme(scheme)

	return client.New(cfg, client.Options{Scheme: scheme})
}

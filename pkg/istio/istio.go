// Package istio provides Terratest-style helpers for testing Istio resources
// in Kubernetes clusters. It offers utility functions to create Istio clients,
// check resource readiness conditions, and interact with Istio custom resources
// like ServiceEntries and other networking components.
package istio

import (
	"testing"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	istiometa "istio.io/api/meta/v1alpha1"
	istionetworking "istio.io/api/networking/v1alpha3"
	istioClientset "istio.io/client-go/pkg/clientset/versioned"
)

// NewClient creates and returns a new Istio Client for use in tests.
// It initializes the Kubernetes REST configuration and the Istio clientset,
// failing the test if any errors occur during setup.
//
// Parameters:
//   - t: The testing context used for logging and error handling.
//
// Returns:
//   - *Client: A pointer to the initialized Istio client.
var NewClient = newClient

func newClient(t *testing.T, options *k8s.KubectlOptions) *istioClientset.Clientset {
	cfg, err := utils.GetRestConfigE(t, options)
	require.NoError(t, err)

	client, err := istioClientset.NewForConfig(cfg)
	require.NoError(t, err, "Failed to create Istio client")

	return client
}

func istioConditionReady(t *testing.T, status *istiometa.IstioStatus) bool {
	require.NotNil(t, status)
	var found bool
	for _, condition := range status.Conditions {
		if condition.Type == "Ready" && condition.Status == "true" {
			found = true
			break
		}
	}
	return found
}

func serviceEntryConditionReady(t *testing.T, status *istionetworking.ServiceEntryStatus) bool {
	require.NotNil(t, status)
	var found bool
	for _, condition := range status.Conditions {
		if condition.Type == "Ready" && condition.Status == "true" {
			found = true
			break
		}
	}
	return found
}

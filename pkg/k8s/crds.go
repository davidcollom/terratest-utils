package k8s

import (
	"context"
	"testing"
	"time"

	terrak8s "github.com/gruntwork-io/terratest/modules/k8s"

	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/stretchr/testify/require"
)

// GetCustomResourceDefinition retrieves a Kubernetes CustomResourceDefinition (CRD) by name using the provided KubectlOptions.
// It fails the test immediately if the CRD cannot be retrieved, reporting the encountered error.
// Returns the retrieved CustomResourceDefinition object.
//
// Parameters:
//
//	t        - The testing context.
//	options  - The kubectl options to use for the request.
//	crdName  - The name of the CRD to retrieve.
//
// Returns:
//
//	*apixv1.CustomResourceDefinition - The requested CRD object.
func GetCustomResourceDefinition(t *testing.T, options *terrak8s.KubectlOptions, crdName string, opts metav1.GetOptions) *apixv1.CustomResourceDefinition {
	t.Helper()

	crd, err := GetCustomResourceDefinitionE(t, options, crdName, opts)
	require.NoError(t, err)
	return crd
}

// GetCustomResourceDefinitionE retrieves a Kubernetes CustomResourceDefinition (CRD) by name using the provided KubectlOptions.
// It returns the CRD object if found, or an error if the retrieval fails.
// This function is intended for use in tests and will mark the test as a helper.
//
// Parameters:
//   - t: The testing context.
//   - options: The KubectlOptions to use for connecting to the cluster.
//   - crdName: The name of the CustomResourceDefinition to retrieve.
//
// Returns:
//   - *apixv1.CustomResourceDefinition: The retrieved CRD object.
//   - error: An error if the CRD could not be retrieved.
func GetCustomResourceDefinitionE(t *testing.T, options *terrak8s.KubectlOptions, crdName string, opts metav1.GetOptions) (*apixv1.CustomResourceDefinition, error) {
	t.Helper()

	client, err := NewAPIXClient(t, options)
	if err != nil {
		return nil, err
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().Get(t.Context(), crdName, opts)
}

// ListCustomResourceDefinitionsE retrieves a list of CustomResourceDefinitions (CRDs) from the Kubernetes cluster
// using the provided KubectlOptions and ListOptions. It fails the test immediately if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the Kubernetes context and namespace.
//   - opts: The list options to filter the CRDs.
//
// Returns:
//   - A pointer to a CustomResourceDefinitionList containing the CRDs found in the cluster.
//   - error: An error if the list could not be retrieved.
func ListCustomResourceDefinitionsE(t *testing.T, options *terrak8s.KubectlOptions, opts metav1.ListOptions) (*apixv1.CustomResourceDefinitionList, error) {
	t.Helper()

	client, err := NewAPIXClient(t, options)
	if err != nil {
		return nil, err
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().List(t.Context(), opts)
}

// WaitForCustomResourceDefinitionIsReady waits until the specified CustomResourceDefinition (CRD) is ready in the Kubernetes cluster.
// It polls the CRD status at regular intervals until it is ready or the timeout is reached.
// If the CRD does not become ready within the given timeout, the test fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the cluster.
//   - crdName: The name of the CRD to check for readiness.
//   - timeout: The maximum duration to wait for the CRD to become ready.
func WaitForCustomResourceDefinitionIsReady(t *testing.T, options *terrak8s.KubectlOptions, crdName string, timeout time.Duration) {
	t.Helper()

	client, err := NewAPIXClient(t, options)
	if err != nil {
		t.Fatalf("Failed to create APIX client: %v", err)
	}

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		crd, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crdName, metav1.GetOptions{})
		if err != nil {
			return false, nil // retry
		}
		if IsCustomResourceDefinitionReady(crd) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("CustomResourceDefinition %s was not Ready in time: %v", crdName, err)
	}
}

// IsCustomResourceDefinitionReady checks whether the given CustomResourceDefinition (CRD)
// is ready by verifying that both the 'Established' and 'NamesAccepted' conditions are true.
// It returns true if both conditions are met, indicating the CRD is fully established and
// its name has been accepted by the Kubernetes API server.
//
// Parameters:
//   - crd: A pointer to the CustomResourceDefinition object to check.
//
// Returns:
//   - bool: True if the CRD is ready (both 'Established' and 'NamesAccepted' conditions are true), false otherwise.
func IsCustomResourceDefinitionReady(crd *apixv1.CustomResourceDefinition) bool {
	conds := crd.Status.Conditions
	var (
		established bool
		accepted    bool
	)

	for _, cond := range conds {
		if cond.Type == apixv1.Established && cond.Status == apixv1.ConditionTrue {
			established = true
		}
		if cond.Type == apixv1.NamesAccepted && cond.Status == apixv1.ConditionTrue {
			accepted = true
		}
	}
	return established && accepted
}

package k8s

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
//
// GetCustomResourceDefinition gets a resource by name.
func GetCustomResourceDefinition(t testing.TestingT, options *KubectlOptions, crdName string, opts metav1.GetOptions) *apixv1.CustomResourceDefinition {
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
//
// GetCustomResourceDefinitionE gets a resource by name.
func GetCustomResourceDefinitionE(t testing.TestingT, options *KubectlOptions, crdName string, opts metav1.GetOptions) (*apixv1.CustomResourceDefinition, error) {
	client, err := NewAPIXClient(t, options)
	if err != nil {
		return nil, err
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), crdName, opts)
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
//
// ListCustomResourceDefinitionsE lists matching resources.
func ListCustomResourceDefinitionsE(t testing.TestingT, options *KubectlOptions, opts metav1.ListOptions) (*apixv1.CustomResourceDefinitionList, error) {
	client, err := NewAPIXClient(t, options)
	if err != nil {
		return nil, err
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().List(context.Background(), opts)
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
//
// WaitForCustomResourceDefinitionIsReady waits for the resource condition to be satisfied.
func WaitForCustomResourceDefinitionIsReady(t testing.TestingT, options *KubectlOptions, crdName string, timeout time.Duration) {
	err := WaitForCustomResourceDefinitionIsReadyE(t, options, crdName, timeout)
	require.NoError(t, err, "CustomResourceDefinition %s was not Ready in time", crdName)
}

// WaitForCustomResourceDefinitionIsReadyE waits for the resource condition to be satisfied.
func WaitForCustomResourceDefinitionIsReadyE(t testing.TestingT, options *KubectlOptions, crdName string, timeout time.Duration) error {
	client, err := NewAPIXClient(t, options)
	if err != nil {
		return err
	}

	return wait.PollUntilContextTimeout(context.Background(), 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		crd, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crdName, metav1.GetOptions{})
		if err != nil {
			return false, nil // retry
		}
		if IsCustomResourceDefinitionReady(crd) {
			return true, nil
		}
		return false, nil
	})
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
//
// IsCustomResourceDefinitionReady returns whether the resource matches the expected state.
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

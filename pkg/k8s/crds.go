package k8s

import (
	"testing"

	terrak8s "github.com/gruntwork-io/terratest/modules/k8s"

	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/davidcollom/terratest-utils/pkg/utils"
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
func GetCustomResourceDefinition(t *testing.T, options *terrak8s.KubectlOptions, crdName string) *apixv1.CustomResourceDefinition {
	t.Helper()

	crd, err := GetCustomResourceDefinitionE(t, options, crdName)
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
func GetCustomResourceDefinitionE(t *testing.T, options *terrak8s.KubectlOptions, crdName string) (*apixv1.CustomResourceDefinition, error) {
	t.Helper()

	restConfig, err := utils.GetRestConfigE(t, options)
	if err != nil {
		return nil, err
	}

	client, err := apixclientset.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().Get(t.Context(), crdName, metav1.GetOptions{})
}

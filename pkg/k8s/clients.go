package k8s

import (
	"testing"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	terrak8s "github.com/gruntwork-io/terratest/modules/k8s"

	apixclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// NewClient gets a new standard kubernetes clientset
var NewClient = terrak8s.GetKubernetesClientFromOptionsE

// NewAPIXClient creates a new API Extensions (apix) clientset using the provided
// terrak8s.KubectlOptions. It returns an apixclientset.Interface for interacting
// with Kubernetes API extensions resources, or an error if the client could not
// be created. The testing.T object is used for test context and error reporting.
//
// Parameters:
//   - t: The testing context, used for helper annotation and error reporting.
//   - options: The kubectl options containing cluster access configuration.
//
// Returns:
//   - apixclientset.Interface: The API Extensions clientset for interacting with CRDs and other API extensions.
//   - error: An error if the configuration or clientset could not be created.
//
// Example usage:
//
//	client, err := newAPIXClient(t, options)
//	require.NoError(t, err)
//	crds, err := client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
//	require.NoError(t, err)
var NewAPIXClient = newAPIXClient

func newAPIXClient(t *testing.T, options *terrak8s.KubectlOptions) (apixclientset.Interface, error) {
	t.Helper()

	restConfig, err := utils.GetRestConfigE(t, options)
	if err != nil {
		return nil, err
	}

	client, err := apixclientset.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

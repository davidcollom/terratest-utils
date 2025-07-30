package linkerd

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	linkerdpolicyv1alpha1 "github.com/linkerd/linkerd2/controller/gen/apis/policy/v1alpha1"
	"github.com/stretchr/testify/require"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListNetworkAuthentications retrieves all Linkerd NetworkAuthentication resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list NetworkAuthentications from.
//
// Returns:
//   - A slice of pointers to NetworkAuthentication objects found in the namespace.
func ListNetworkAuthentications(t *testing.T, options *k8s.KubectlOptions, namespace string) []*linkerdpolicyv1alpha1.NetworkAuthentication {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	networkAuthentications, err := linkerdClient.PolicyV1alpha1().NetworkAuthentications(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list NetworkAuthentications in namespace %s", namespace)

	// Convert slice of values to slice of pointers
	var result []*linkerdpolicyv1alpha1.NetworkAuthentication
	for i := range networkAuthentications.Items {
		result = append(result, &networkAuthentications.Items[i])
	}

	return result
}

// GetNetworkAuthentication retrieves a specific Linkerd NetworkAuthentication resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the NetworkAuthentication to retrieve.
//   - namespace: The namespace of the NetworkAuthentication.
//
// Returns:
//   - A pointer to the NetworkAuthentication object.
func GetNetworkAuthentication(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *linkerdpolicyv1alpha1.NetworkAuthentication {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	networkAuthentication, err := linkerdClient.PolicyV1alpha1().NetworkAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get NetworkAuthentication %s in namespace %s", name, namespace)

	return networkAuthentication
}

// WaitForNetworkAuthenticationExists waits until the specified NetworkAuthentication exists in the given namespace or the timeout is reached.
// It polls the NetworkAuthentication every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the NetworkAuthentication to check.
//   - namespace: The namespace of the NetworkAuthentication.
//   - timeout: The maximum duration to wait for the resource to exist.
func WaitForNetworkAuthenticationExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.PolicyV1alpha1().NetworkAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("NetworkAuthentication %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}

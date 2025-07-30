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

// ListMeshTLSAuthentications retrieves all Linkerd MeshTLSAuthentication resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list MeshTLSAuthentications from.
//
// Returns:
//   - A slice of pointers to MeshTLSAuthentication objects found in the namespace.
func ListMeshTLSAuthentications(t *testing.T, options *k8s.KubectlOptions, namespace string) []*linkerdpolicyv1alpha1.MeshTLSAuthentication {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	meshTLSAuthentications, err := linkerdClient.PolicyV1alpha1().MeshTLSAuthentications(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list MeshTLSAuthentications in namespace %s", namespace)

	// Convert slice of values to slice of pointers
	var result []*linkerdpolicyv1alpha1.MeshTLSAuthentication
	for i := range meshTLSAuthentications.Items {
		result = append(result, &meshTLSAuthentications.Items[i])
	}

	return result
}

// GetMeshTLSAuthentication retrieves a specific Linkerd MeshTLSAuthentication resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the MeshTLSAuthentication to retrieve.
//   - namespace: The namespace of the MeshTLSAuthentication.
//
// Returns:
//   - A pointer to the MeshTLSAuthentication object.
func GetMeshTLSAuthentication(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *linkerdpolicyv1alpha1.MeshTLSAuthentication {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	meshTLSAuthentication, err := linkerdClient.PolicyV1alpha1().MeshTLSAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get MeshTLSAuthentication %s in namespace %s", name, namespace)

	return meshTLSAuthentication
}

// WaitForMeshTLSAuthenticationExists waits until the specified MeshTLSAuthentication exists in the given namespace or the timeout is reached.
// It polls the MeshTLSAuthentication every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the MeshTLSAuthentication to check.
//   - namespace: The namespace of the MeshTLSAuthentication.
//   - timeout: The maximum duration to wait for the resource to exist.
func WaitForMeshTLSAuthenticationExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.PolicyV1alpha1().MeshTLSAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("MeshTLSAuthentication %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}

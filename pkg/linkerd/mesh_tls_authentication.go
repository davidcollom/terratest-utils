package linkerd

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
//
// ListMeshTLSAuthentications lists matching resources.
func ListMeshTLSAuthentications(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*linkerdpolicyv1alpha1.MeshTLSAuthentication {
	meshTLSAuthentications, err := ListMeshTLSAuthenticationsE(t, options, namespace)
	require.NoError(t, err, "Failed to list MeshTLSAuthentications in namespace %s", namespace)
	return meshTLSAuthentications
}

// ListMeshTLSAuthenticationsE lists matching resources.
func ListMeshTLSAuthenticationsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*linkerdpolicyv1alpha1.MeshTLSAuthentication, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	meshTLSAuthentications, err := linkerdClient.PolicyV1alpha1().MeshTLSAuthentications(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Convert slice of values to slice of pointers
	var result []*linkerdpolicyv1alpha1.MeshTLSAuthentication
	for i := range meshTLSAuthentications.Items {
		result = append(result, &meshTLSAuthentications.Items[i])
	}

	return result, nil
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
//
// GetMeshTLSAuthentication gets a resource by name.
func GetMeshTLSAuthentication(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) *linkerdpolicyv1alpha1.MeshTLSAuthentication {
	meshTLSAuthentication, err := GetMeshTLSAuthenticationE(t, options, name, namespace)
	require.NoError(t, err, "Failed to get MeshTLSAuthentication %s in namespace %s", name, namespace)
	return meshTLSAuthentication
}

// GetMeshTLSAuthenticationE gets a resource by name.
func GetMeshTLSAuthenticationE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) (*linkerdpolicyv1alpha1.MeshTLSAuthentication, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	meshTLSAuthentication, err := linkerdClient.PolicyV1alpha1().MeshTLSAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	return meshTLSAuthentication, nil
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
//
// WaitForMeshTLSAuthenticationExists waits for the resource condition to be satisfied.
func WaitForMeshTLSAuthenticationExists(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForMeshTLSAuthenticationExistsE(t, options, name, namespace, timeout)
	require.NoError(t, err, "MeshTLSAuthentication %s/%s did not exist within timeout", namespace, name)
}

// WaitForMeshTLSAuthenticationExistsE waits for the resource condition to be satisfied.
func WaitForMeshTLSAuthenticationExistsE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.PolicyV1alpha1().MeshTLSAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}

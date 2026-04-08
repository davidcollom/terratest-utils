package linkerd

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	linkerdserverauthorizationv1beta1 "github.com/linkerd/linkerd2/controller/gen/apis/serverauthorization/v1beta1"
	"github.com/stretchr/testify/require"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListServerAuthorizations retrieves all Linkerd ServerAuthorization resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list ServerAuthorizations from.
//
// Returns:
//   - A slice of pointers to ServerAuthorization objects found in the namespace.
// ListServerAuthorizations lists matching resources.
func ListServerAuthorizations(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*linkerdserverauthorizationv1beta1.ServerAuthorization {
	serverAuthorizations, err := ListServerAuthorizationsE(t, options, namespace)
	require.NoError(t, err, "Failed to list ServerAuthorizations in namespace %s", namespace)
	return serverAuthorizations
}

// ListServerAuthorizationsE lists matching resources.
func ListServerAuthorizationsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*linkerdserverauthorizationv1beta1.ServerAuthorization, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	serverAuthorizations, err := linkerdClient.ServerauthorizationV1beta1().ServerAuthorizations(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Convert slice of values to slice of pointers
	var result []*linkerdserverauthorizationv1beta1.ServerAuthorization
	for i := range serverAuthorizations.Items {
		result = append(result, &serverAuthorizations.Items[i])
	}

	return result, nil
}

// GetServerAuthorization retrieves a specific Linkerd ServerAuthorization resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the ServerAuthorization to retrieve.
//   - namespace: The namespace of the ServerAuthorization.
//
// Returns:
//   - A pointer to the ServerAuthorization object.
// GetServerAuthorization gets a resource by name.
func GetServerAuthorization(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) *linkerdserverauthorizationv1beta1.ServerAuthorization {
	serverAuthorization, err := GetServerAuthorizationE(t, options, name, namespace)
	require.NoError(t, err, "Failed to get ServerAuthorization %s in namespace %s", name, namespace)
	return serverAuthorization
}

// GetServerAuthorizationE gets a resource by name.
func GetServerAuthorizationE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) (*linkerdserverauthorizationv1beta1.ServerAuthorization, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	serverAuthorization, err := linkerdClient.ServerauthorizationV1beta1().ServerAuthorizations(namespace).Get(ctx, name, v1meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	return serverAuthorization, nil
}

// WaitForServerAuthorizationExists waits until the specified ServerAuthorization exists in the given namespace or the timeout is reached.
// It polls the ServerAuthorization every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the ServerAuthorization to check.
//   - namespace: The namespace of the ServerAuthorization.
//   - timeout: The maximum duration to wait for the resource to exist.
// WaitForServerAuthorizationExists waits for the resource condition to be satisfied.
func WaitForServerAuthorizationExists(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForServerAuthorizationExistsE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ServerAuthorization %s/%s did not exist within timeout", namespace, name)
}

// WaitForServerAuthorizationExistsE waits for the resource condition to be satisfied.
func WaitForServerAuthorizationExistsE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.ServerauthorizationV1beta1().ServerAuthorizations(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}

package linkerd

import (
	"context"
	"testing"
	"time"

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
func ListServerAuthorizations(t *testing.T, options *k8s.KubectlOptions, namespace string) []*linkerdserverauthorizationv1beta1.ServerAuthorization {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	serverAuthorizations, err := linkerdClient.ServerauthorizationV1beta1().ServerAuthorizations(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list ServerAuthorizations in namespace %s", namespace)

	// Convert slice of values to slice of pointers
	var result []*linkerdserverauthorizationv1beta1.ServerAuthorization
	for i := range serverAuthorizations.Items {
		result = append(result, &serverAuthorizations.Items[i])
	}

	return result
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
func GetServerAuthorization(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *linkerdserverauthorizationv1beta1.ServerAuthorization {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	serverAuthorization, err := linkerdClient.ServerauthorizationV1beta1().ServerAuthorizations(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get ServerAuthorization %s in namespace %s", name, namespace)

	return serverAuthorization
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
func WaitForServerAuthorizationExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.ServerauthorizationV1beta1().ServerAuthorizations(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("ServerAuthorization %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}

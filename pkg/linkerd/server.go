package linkerd

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	linkerdserverv1beta1 "github.com/linkerd/linkerd2/controller/gen/apis/server/v1beta1"
	"github.com/stretchr/testify/require"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListServers retrieves all Linkerd Server resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list Servers from.
//
// Returns:
//   - A slice of pointers to Server objects found in the namespace.
func ListServers(t *testing.T, options *k8s.KubectlOptions, namespace string) []*linkerdserverv1beta1.Server {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	servers, err := linkerdClient.ServerV1beta1().Servers(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Servers in namespace %s", namespace)

	// Convert slice of values to slice of pointers
	var result []*linkerdserverv1beta1.Server
	for i := range servers.Items {
		result = append(result, &servers.Items[i])
	}

	return result
}

// GetServer retrieves a specific Linkerd Server resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the Server to retrieve.
//   - namespace: The namespace of the Server.
//
// Returns:
//   - A pointer to the Server object.
func GetServer(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *linkerdserverv1beta1.Server {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	server, err := linkerdClient.ServerV1beta1().Servers(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get Server %s in namespace %s", name, namespace)

	return server
}

// WaitForServerExists waits until the specified Server exists in the given namespace or the timeout is reached.
// It polls the Server every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the Server to check.
//   - namespace: The namespace of the Server.
//   - timeout: The maximum duration to wait for the resource to exist.
func WaitForServerExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.ServerV1beta1().Servers(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("Server %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}

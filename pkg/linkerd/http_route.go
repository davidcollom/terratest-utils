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

// ListHTTPRoutes retrieves all Linkerd HTTPRoute resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list HTTPRoutes from.
//
// Returns:
//   - A slice of pointers to HTTPRoute objects found in the namespace.
//
// ListHTTPRoutes lists matching resources.
func ListHTTPRoutes(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*linkerdpolicyv1alpha1.HTTPRoute {
	httpRoutes, err := ListHTTPRoutesE(t, options, namespace)
	require.NoError(t, err, "Failed to list HTTPRoutes in namespace %s", namespace)
	return httpRoutes
}

// ListHTTPRoutesE lists matching resources.
func ListHTTPRoutesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*linkerdpolicyv1alpha1.HTTPRoute, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	httpRoutes, err := linkerdClient.PolicyV1alpha1().HTTPRoutes(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Convert slice of values to slice of pointers
	var result []*linkerdpolicyv1alpha1.HTTPRoute
	for i := range httpRoutes.Items {
		result = append(result, &httpRoutes.Items[i])
	}

	return result, nil
}

// GetHTTPRoute retrieves a specific Linkerd HTTPRoute resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the HTTPRoute to retrieve.
//   - namespace: The namespace of the HTTPRoute.
//
// Returns:
//   - A pointer to the HTTPRoute object.
//
// GetHTTPRoute gets a resource by name.
func GetHTTPRoute(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) *linkerdpolicyv1alpha1.HTTPRoute {
	httpRoute, err := GetHTTPRouteE(t, options, name, namespace)
	require.NoError(t, err, "Failed to get HTTPRoute %s in namespace %s", name, namespace)
	return httpRoute
}

// GetHTTPRouteE gets a resource by name.
func GetHTTPRouteE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) (*linkerdpolicyv1alpha1.HTTPRoute, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	httpRoute, err := linkerdClient.PolicyV1alpha1().HTTPRoutes(namespace).Get(ctx, name, v1meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	return httpRoute, nil
}

// WaitForHTTPRouteExists waits until the specified HTTPRoute exists in the given namespace or the timeout is reached.
// It polls the HTTPRoute every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the HTTPRoute to check.
//   - namespace: The namespace of the HTTPRoute.
//   - timeout: The maximum duration to wait for the resource to exist.
//
// WaitForHTTPRouteExists waits for the resource condition to be satisfied.
func WaitForHTTPRouteExists(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForHTTPRouteExistsE(t, options, name, namespace, timeout)
	require.NoError(t, err, "HTTPRoute %s/%s did not exist within timeout", namespace, name)
}

// WaitForHTTPRouteExistsE waits for the resource condition to be satisfied.
func WaitForHTTPRouteExistsE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.PolicyV1alpha1().HTTPRoutes(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}

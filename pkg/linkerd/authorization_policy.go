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

// ListAuthorizationPolicies retrieves all Linkerd AuthorizationPolicy resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list AuthorizationPolicies from.
//
// Returns:
//   - A slice of pointers to AuthorizationPolicy objects found in the namespace.
func ListAuthorizationPolicies(t *testing.T, options *k8s.KubectlOptions, namespace string) []*linkerdpolicyv1alpha1.AuthorizationPolicy {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	authorizationPolicies, err := linkerdClient.PolicyV1alpha1().AuthorizationPolicies(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list AuthorizationPolicies in namespace %s", namespace)

	// Convert slice of values to slice of pointers
	var result []*linkerdpolicyv1alpha1.AuthorizationPolicy
	for i := range authorizationPolicies.Items {
		result = append(result, &authorizationPolicies.Items[i])
	}

	return result
}

// GetAuthorizationPolicy retrieves a specific Linkerd AuthorizationPolicy resource by name in the specified namespace.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the AuthorizationPolicy to retrieve.
//   - namespace: The namespace of the AuthorizationPolicy.
//
// Returns:
//   - A pointer to the AuthorizationPolicy object.
func GetAuthorizationPolicy(t *testing.T, options *k8s.KubectlOptions, name, namespace string) *linkerdpolicyv1alpha1.AuthorizationPolicy {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	authorizationPolicy, err := linkerdClient.PolicyV1alpha1().AuthorizationPolicies(namespace).Get(ctx, name, v1meta.GetOptions{})
	require.NoError(t, err, "Failed to get AuthorizationPolicy %s in namespace %s", name, namespace)

	return authorizationPolicy
}

// WaitForAuthorizationPolicyExists waits until the specified AuthorizationPolicy exists in the given namespace or the timeout is reached.
// It polls the AuthorizationPolicy every 2 seconds until it exists.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the AuthorizationPolicy to check.
//   - namespace: The namespace of the AuthorizationPolicy.
//   - timeout: The maximum duration to wait for the resource to exist.
func WaitForAuthorizationPolicyExists(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	linkerdClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.PolicyV1alpha1().AuthorizationPolicies(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		t.Fatalf("AuthorizationPolicy %s/%s did not exist within timeout: %v", namespace, name, err)
	}
}

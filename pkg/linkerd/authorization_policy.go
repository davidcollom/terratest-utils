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
// ListAuthorizationPolicies lists matching resources.
func ListAuthorizationPolicies(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*linkerdpolicyv1alpha1.AuthorizationPolicy {
	authorizationPolicies, err := ListAuthorizationPoliciesE(t, options, namespace)
	require.NoError(t, err, "Failed to list AuthorizationPolicies in namespace %s", namespace)
	return authorizationPolicies
}

// ListAuthorizationPoliciesE lists matching resources.
func ListAuthorizationPoliciesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*linkerdpolicyv1alpha1.AuthorizationPolicy, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	authorizationPolicies, err := linkerdClient.PolicyV1alpha1().AuthorizationPolicies(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Convert slice of values to slice of pointers
	var result []*linkerdpolicyv1alpha1.AuthorizationPolicy
	for i := range authorizationPolicies.Items {
		result = append(result, &authorizationPolicies.Items[i])
	}

	return result, nil
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
// GetAuthorizationPolicy gets a resource by name.
func GetAuthorizationPolicy(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) *linkerdpolicyv1alpha1.AuthorizationPolicy {
	authorizationPolicy, err := GetAuthorizationPolicyE(t, options, name, namespace)
	require.NoError(t, err, "Failed to get AuthorizationPolicy %s in namespace %s", name, namespace)
	return authorizationPolicy
}

// GetAuthorizationPolicyE gets a resource by name.
func GetAuthorizationPolicyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string) (*linkerdpolicyv1alpha1.AuthorizationPolicy, error) {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	authorizationPolicy, err := linkerdClient.PolicyV1alpha1().AuthorizationPolicies(namespace).Get(ctx, name, v1meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	return authorizationPolicy, nil
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
// WaitForAuthorizationPolicyExists waits for the resource condition to be satisfied.
func WaitForAuthorizationPolicyExists(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForAuthorizationPolicyExistsE(t, options, name, namespace, timeout)
	require.NoError(t, err, "AuthorizationPolicy %s/%s did not exist within timeout", namespace, name)
}

// WaitForAuthorizationPolicyExistsE waits for the resource condition to be satisfied.
func WaitForAuthorizationPolicyExistsE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	linkerdClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		_, err := linkerdClient.PolicyV1alpha1().AuthorizationPolicies(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
}

package istio

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	istiosecurityv1 "istio.io/client-go/pkg/apis/security/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListAuthorizationPolicies retrieves all Istio AuthorizationPolicy resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list AuthorizationPolicies from.
//
// Returns:
//   - A slice of pointers to AuthorizationPolicy objects found in the namespace.
//
// ListAuthorizationPolicies lists matching resources.
func ListAuthorizationPolicies(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istiosecurityv1.AuthorizationPolicy {
	authorizationPolicies, err := ListAuthorizationPoliciesE(t, options, namespace)
	require.NoError(t, err, "Failed to list Authorization Policies in namespace %s", namespace)
	return authorizationPolicies
}

// ListAuthorizationPoliciesE lists matching resources.
func ListAuthorizationPoliciesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istiosecurityv1.AuthorizationPolicy, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	authorizationPolicies, err := istioClient.SecurityV1().AuthorizationPolicies(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return authorizationPolicies.Items, nil
}

// WaitForAuthorizationPolicyReady waits until the specified AuthorizationPolicy in the given namespace is Ready or the timeout is reached.
// It polls the AuthorizationPolicy status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the AuthorizationPolicy to check.
//   - namespace: The namespace of the AuthorizationPolicy.
//   - timeout: The maximum duration to wait for the resource to become Ready.
//
// WaitForAuthorizationPolicyReady waits for the resource condition to be satisfied.
func WaitForAuthorizationPolicyReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForAuthorizationPolicyReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "AuthorizationPolicy %s/%s did not become Ready", namespace, name)
}

// WaitForAuthorizationPolicyReadyE waits for the resource condition to be satisfied.
func WaitForAuthorizationPolicyReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var authorizationPolicy *istiosecurityv1.AuthorizationPolicy
		authorizationPolicy, err := istioClient.SecurityV1().AuthorizationPolicies(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if authorizationPolicy.Status.Conditions != nil {
			return istioConditionReady(t, &authorizationPolicy.Status), nil
		}
		return false, nil
	})
}

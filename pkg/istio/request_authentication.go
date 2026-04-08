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

// ListRequestAuthentications retrieves all Istio RequestAuthentication resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list RequestAuthentications from.
//
// Returns:
//   - A slice of pointers to RequestAuthentication objects found in the namespace.
// ListRequestAuthentications lists matching resources.
func ListRequestAuthentications(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istiosecurityv1.RequestAuthentication {
	requestAuthentications, err := ListRequestAuthenticationsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Request Authentications in namespace %s", namespace)
	return requestAuthentications
}

// ListRequestAuthenticationsE lists matching resources.
func ListRequestAuthenticationsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istiosecurityv1.RequestAuthentication, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	requestAuthentications, err := istioClient.SecurityV1().RequestAuthentications(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return requestAuthentications.Items, nil
}

// WaitForRequestAuthenticationReady waits until the specified RequestAuthentication in the given namespace is Ready or the timeout is reached.
// It polls the RequestAuthentication status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the RequestAuthentication to check.
//   - namespace: The namespace of the RequestAuthentication.
//   - timeout: The maximum duration to wait for the resource to become Ready.
// WaitForRequestAuthenticationReady waits for the resource condition to be satisfied.
func WaitForRequestAuthenticationReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForRequestAuthenticationReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "RequestAuthentication %s/%s did not become Ready", namespace, name)
}

// WaitForRequestAuthenticationReadyE waits for the resource condition to be satisfied.
func WaitForRequestAuthenticationReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var requestAuthentication *istiosecurityv1.RequestAuthentication
		requestAuthentication, err := istioClient.SecurityV1().RequestAuthentications(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if requestAuthentication.Status.Conditions != nil {
			return istioConditionReady(t, &requestAuthentication.Status), nil
		}
		return false, nil
	})
}

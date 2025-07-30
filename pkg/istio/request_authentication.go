package istio

import (
	"context"
	"testing"
	"time"

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
func ListRequestAuthentications(t *testing.T, options *k8s.KubectlOptions, namespace string) []*istiosecurityv1.RequestAuthentication {
	t.Helper()

	istioClient := NewClient(t, options)

	ctx := t.Context()
	requestAuthentications, err := istioClient.SecurityV1().RequestAuthentications(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Request Authentications in namespace %s", namespace)

	return requestAuthentications.Items
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
func WaitForRequestAuthenticationReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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

	if err != nil {
		t.Fatalf("RequestAuthentication %s/%s did not become Ready: %v", namespace, name, err)
	}
}

package istio

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListGateways retrieves all Istio Gateway resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list Gateways from.
//
// Returns:
//   - A slice of pointers to Gateway objects found in the namespace.
func ListGateways(t *testing.T, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.Gateway {
	t.Helper()

	istioClient := NewClient(t, options)

	ctx := t.Context()
	gateways, err := istioClient.NetworkingV1alpha3().Gateways(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Gateways in namespace %s", namespace)

	return gateways.Items
}

// WaitForGatewayReady waits until the specified Gateway in the given namespace is Ready or the timeout is reached.
// It polls the Gateway status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the Gateway to check.
//   - namespace: The namespace of the Gateway.
//   - timeout: The maximum duration to wait for the resource to become Ready.
func WaitForGatewayReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var gateway *istionetworkingv1alpha3.Gateway
		gateway, err := istioClient.NetworkingV1alpha3().Gateways(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if gateway.Status.Conditions != nil {
			return istioConditionReady(t, &gateway.Status), nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Gateway %s/%s did not become Ready: %v", namespace, name, err)
	}
}

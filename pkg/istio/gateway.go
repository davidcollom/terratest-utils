package istio

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
//
// ListGateways lists matching resources.
func ListGateways(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.Gateway {
	gateways, err := ListGatewaysE(t, options, namespace)
	require.NoError(t, err, "Failed to list Gateways in namespace %s", namespace)
	return gateways
}

// ListGatewaysE lists matching resources.
func ListGatewaysE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.Gateway, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	gateways, err := istioClient.NetworkingV1alpha3().Gateways(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return gateways.Items, nil
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
//
// WaitForGatewayReady waits for the resource condition to be satisfied.
func WaitForGatewayReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForGatewayReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Gateway %s/%s did not become Ready", namespace, name)
}

// WaitForGatewayReadyE waits for the resource condition to be satisfied.
func WaitForGatewayReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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
}

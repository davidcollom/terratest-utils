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

// ListVirtualServices retrieves all Istio VirtualService resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list VirtualServices from.
//
// Returns:
//   - A slice of pointers to VirtualService objects found in the namespace.
// ListVirtualServices lists matching resources.
func ListVirtualServices(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.VirtualService {
	virtualServices, err := ListVirtualServicesE(t, options, namespace)
	require.NoError(t, err, "Failed to list Virtual Services in namespace %s", namespace)
	return virtualServices
}

// ListVirtualServicesE lists matching resources.
func ListVirtualServicesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.VirtualService, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	virtualServices, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return virtualServices.Items, nil
}

// WaitForVirtualServiceReady waits until the specified VirtualService in the given namespace is Ready or the timeout is reached.
// It polls the VirtualService status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the VirtualService to check.
//   - namespace: The namespace of the VirtualService.
//   - timeout: The maximum duration to wait for the resource to become Ready.
// WaitForVirtualServiceReady waits for the resource condition to be satisfied.
func WaitForVirtualServiceReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForVirtualServiceReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "VirtualService %s/%s did not become Ready", namespace, name)
}

// WaitForVirtualServiceReadyE waits for the resource condition to be satisfied.
func WaitForVirtualServiceReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var virtualService *istionetworkingv1alpha3.VirtualService
		virtualService, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if virtualService.Status.Conditions != nil {
			return istioConditionReady(t, &virtualService.Status), nil
		}
		return false, nil
	})
}

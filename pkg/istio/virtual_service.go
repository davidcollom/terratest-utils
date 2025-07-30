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
func ListVirtualServices(t *testing.T, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.VirtualService {
	t.Helper()

	istioClient := NewClient(t, options)

	ctx := t.Context()
	virtualServices, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(ctx, v1meta.ListOptions{})
	require.NoError(t, err, "Failed to list Virtual Services in namespace %s", namespace)

	return virtualServices.Items
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
func WaitForVirtualServiceReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := t.Context()
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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

	if err != nil {
		t.Fatalf("VirtualService %s/%s did not become Ready: %v", namespace, name, err)
	}
}

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

// ListSidecars retrieves all Istio Sidecar resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list Sidecars from.
//
// Returns:
//   - A slice of pointers to Sidecar objects found in the namespace.
//
// ListSidecars lists matching resources.
func ListSidecars(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.Sidecar {
	sidecars, err := ListSidecarsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Sidecars in namespace %s", namespace)
	return sidecars
}

// ListSidecarsE lists matching resources.
func ListSidecarsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.Sidecar, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	sidecars, err := istioClient.NetworkingV1alpha3().Sidecars(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return sidecars.Items, nil
}

// WaitForSidecarReady waits until the specified Sidecar in the given namespace is Ready or the timeout is reached.
// It polls the Sidecar status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the Sidecar to check.
//   - namespace: The namespace of the Sidecar.
//   - timeout: The maximum duration to wait for the resource to become Ready.
//
// WaitForSidecarReady waits for the resource condition to be satisfied.
func WaitForSidecarReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForSidecarReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Sidecar %s/%s did not become Ready", namespace, name)
}

// WaitForSidecarReadyE waits for the resource condition to be satisfied.
func WaitForSidecarReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var sidecar *istionetworkingv1alpha3.Sidecar
		sidecar, err := istioClient.NetworkingV1alpha3().Sidecars(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if sidecar.Status.Conditions != nil {
			return istioConditionReady(t, &sidecar.Status), nil
		}
		return false, nil
	})
}

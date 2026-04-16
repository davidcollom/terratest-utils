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

// ListEnvoyFilters retrieves all Istio EnvoyFilter resources in the specified namespace using the provided KubectlOptions.
// It fails the test if an error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - namespace: The namespace to list EnvoyFilters from.
//
// Returns:
//   - A slice of pointers to EnvoyFilter objects found in the namespace.
// ListEnvoyFilters lists matching resources.
func ListEnvoyFilters(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []*istionetworkingv1alpha3.EnvoyFilter {
	envoyFilters, err := ListEnvoyFiltersE(t, options, namespace)
	require.NoError(t, err, "Failed to list Envoy Filters in namespace %s", namespace)
	return envoyFilters
}

// ListEnvoyFiltersE lists matching resources.
func ListEnvoyFiltersE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]*istionetworkingv1alpha3.EnvoyFilter, error) {
	istioClient := NewClient(t, options)

	ctx := context.Background()
	envoyFilters, err := istioClient.NetworkingV1alpha3().EnvoyFilters(namespace).List(ctx, v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	return envoyFilters.Items, nil
}

// WaitForEnvoyFilterReady waits until the specified EnvoyFilter in the given namespace is Ready or the timeout is reached.
// It polls the EnvoyFilter status every 2 seconds and checks for the Ready condition.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options specifying the context and namespace.
//   - name: The name of the EnvoyFilter to check.
//   - namespace: The namespace of the EnvoyFilter.
//   - timeout: The maximum duration to wait for the resource to become Ready.
// WaitForEnvoyFilterReady waits for the resource condition to be satisfied.
func WaitForEnvoyFilterReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForEnvoyFilterReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "EnvoyFilter %s/%s did not become Ready", namespace, name)
}

// WaitForEnvoyFilterReadyE waits for the resource condition to be satisfied.
func WaitForEnvoyFilterReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	options = k8s.NewKubectlOptions("", "", namespace)
	istioClient := NewClient(t, options)

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var envoyFilter *istionetworkingv1alpha3.EnvoyFilter
		envoyFilter, err := istioClient.NetworkingV1alpha3().EnvoyFilters(namespace).Get(ctx, name, v1meta.GetOptions{})
		if err != nil {
			return false, nil
		}
		if envoyFilter.Status.Conditions != nil {
			return istioConditionReady(t, &envoyFilter.Status), nil
		}
		return false, nil
	})
}

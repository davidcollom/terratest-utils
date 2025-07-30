package flux

import (
	"context"
	"testing"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ListHelmCharts retrieves a list of HelmChart resources from the specified namespace using the provided
// Kubernetes options. It requires a testing context and will fail the test if unable to create the Flux client
// or if listing the HelmCharts fails. Returns a slice of HelmChart objects found in the namespace.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list HelmChart resources.
//
// Returns:
//   - A slice of sourcev1.HelmChart objects present in the specified namespace.
func ListHelmCharts(t *testing.T, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []sourcev1.HelmChart {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := t.Context()
	var charts sourcev1.HelmChartList
	err = fluxclient.List(ctx, &charts, opts...)
	require.NoError(t, err, "Failed to list HelmCharts in namespace %s", namespace)

	return charts.Items
}

// WaitForHelmChartReady waits until the specified HelmChart resource in the given namespace becomes Ready within the provided timeout.
// It uses the Flux client to poll the HelmChart status and checks for the Ready condition.
// If the HelmChart does not become Ready within the timeout, the test fails with a fatal error.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the Kubernetes REST config.
//	name     - The name of the HelmChart resource.
//	namespace- The namespace where the HelmChart is deployed.
//	timeout  - The maximum duration to wait for the HelmChart to become Ready.
//
// Fails the test if the HelmChart does not reach the Ready condition within the timeout.
func WaitForHelmChartReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var chart sourcev1.HelmChart
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &chart)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(chart.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("HelmChart %s/%s did not become Ready: %v", namespace, name, err)
	}
}

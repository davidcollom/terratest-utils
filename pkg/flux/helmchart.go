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

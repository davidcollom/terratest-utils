package flux

import (
	"context"
	"testing"
	"time"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ListHelmReleases retrieves all HelmRelease resources in the specified namespace using the provided kubectl options.
// It requires a testing context and will fail the test if the Flux client cannot be created or if listing the HelmReleases fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list HelmRelease resources.
//
// Returns:
//   - A slice of helmv2.HelmRelease objects found in the specified namespace.
func ListHelmReleases(t *testing.T, options *k8s.KubectlOptions, namespace string) []helmv2.HelmRelease {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	var releases helmv2.HelmReleaseList
	err = fluxclient.List(ctx, &releases, client.InNamespace(namespace))
	require.NoError(t, err, "Failed to list HelmReleases in namespace %s", namespace)

	return releases.Items
}

// WaitForHelmReleaseReady waits until the specified HelmRelease resource in the given namespace
// reaches the Ready condition or the timeout is exceeded. It polls the resource status at regular
// intervals and fails the test if the resource does not become Ready within the timeout period.
//
// Parameters:
//
//	t        - The testing context.
//	options  - The kubectl options containing the Kubernetes REST config.
//	name     - The name of the HelmRelease resource.
//	namespace- The namespace where the HelmRelease is located.
//	timeout  - The maximum duration to wait for the HelmRelease to become Ready.
//
// The function will call t.Fatalf if the HelmRelease does not become Ready within the timeout.
func WaitForHelmReleaseReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var release helmv2.HelmRelease
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &release)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(release.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("HelmRelease %s/%s did not become Ready: %v", namespace, name, err)
	}
}

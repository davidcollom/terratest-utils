package flux

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
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
//
// ListHelmReleases lists matching resources.
func ListHelmReleases(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []helmv2.HelmRelease {
	releases, err := ListHelmReleasesE(t, options, namespace, opts...)
	require.NoError(t, err, "Failed to list HelmReleases in namespace %s", namespace)
	return releases
}

// ListHelmReleasesE lists matching resources.
func ListHelmReleasesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) ([]helmv2.HelmRelease, error) {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return nil, err
	}

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := context.Background()
	var releases helmv2.HelmReleaseList
	err = fluxclient.List(ctx, &releases, opts...)
	if err != nil {
		return nil, err
	}

	return releases.Items, nil
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
// WaitForHelmReleaseReady waits for the resource condition to be satisfied.
func WaitForHelmReleaseReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForHelmReleaseReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "HelmRelease %s/%s did not become Ready", namespace, name)
}

// WaitForHelmReleaseReadyE waits for the resource condition to be satisfied.
func WaitForHelmReleaseReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var release helmv2.HelmRelease
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &release)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(release.Status.Conditions), nil
	})
}

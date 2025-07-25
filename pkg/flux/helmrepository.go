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

// ListHelmRepositories retrieves all HelmRepository resources in the specified namespace using the provided
// Kubernetes options. It returns a slice of HelmRepository objects. The function fails the test if it is unable
// to create a Flux client or if it encounters an error while listing the HelmRepositories.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list HelmRepository resources.
//
// Returns:
//   - A slice of sourcev1.HelmRepository objects found in the specified namespace.
func ListHelmRepositories(t *testing.T, options *k8s.KubectlOptions, namespace string) []sourcev1.HelmRepository {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	var repos sourcev1.HelmRepositoryList
	err = fluxclient.List(ctx, &repos, client.InNamespace(namespace))
	require.NoError(t, err, "Failed to list HelmRepositories in namespace %s", namespace)

	return repos.Items
}

// WaitForHelmRepositoryReady waits until the specified Flux HelmRepository resource becomes Ready within the given timeout.
// It polls the resource status every 2 seconds and fails the test if the resource does not become Ready in time.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the HelmRepository resource.
//   - namespace: The namespace of the HelmRepository resource.
//   - timeout: The maximum duration to wait for the resource to become Ready.
//
// Fails the test if the HelmRepository does not reach the Ready condition within the timeout.
func WaitForHelmRepositoryReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var helmrepo sourcev1.HelmRepository
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &helmrepo)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(helmrepo.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("HelmRepository %s/%s did not become Ready: %v", namespace, name, err)
	}
}

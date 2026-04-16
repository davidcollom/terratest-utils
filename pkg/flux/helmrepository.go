package flux

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
// ListHelmRepositories lists matching resources.
func ListHelmRepositories(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []sourcev1.HelmRepository {
	repos, err := ListHelmRepositoriesE(t, options, namespace, opts...)
	require.NoError(t, err, "Failed to list HelmRepositories in namespace %s", namespace)
	return repos
}

// ListHelmRepositoriesE lists matching resources.
func ListHelmRepositoriesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) ([]sourcev1.HelmRepository, error) {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return nil, err
	}

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := context.Background()
	var repos sourcev1.HelmRepositoryList
	err = fluxclient.List(ctx, &repos, opts...)
	if err != nil {
		return nil, err
	}

	return repos.Items, nil
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
// WaitForHelmRepositoryReady waits for the resource condition to be satisfied.
func WaitForHelmRepositoryReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForHelmRepositoryReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "HelmRepository %s/%s did not become Ready", namespace, name)
}

// WaitForHelmRepositoryReadyE waits for the resource condition to be satisfied.
func WaitForHelmRepositoryReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var helmrepo sourcev1.HelmRepository
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &helmrepo)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(helmrepo.Status.Conditions), nil
	})
}

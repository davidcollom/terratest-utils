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

// ListOCIRepositories retrieves a list of OCIRepository resources from the specified namespace
// using the provided Kubernetes options. It returns a slice of sourcev1.OCIRepository objects.
// The function requires a testing.T instance for error handling and test context propagation.
// It fails the test if the Flux client cannot be created or if listing the OCIRepositories fails.
//
// Parameters:
//   - t:        The testing.T instance used for test context and assertions.
//   - options:  The KubectlOptions specifying the Kubernetes context and configuration.
//   - namespace: The namespace from which to list OCIRepository resources.
//
// Returns:
//   - []sourcev1.OCIRepository: A slice containing the retrieved OCIRepository resources.
// ListOCIRepositories lists matching resources.
func ListOCIRepositories(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []sourcev1.OCIRepository {
	repos, err := ListOCIRepositoriesE(t, options, namespace, opts...)
	require.NoError(t, err, "Failed to list OCIRepositories in namespace %s", namespace)
	return repos
}

// ListOCIRepositoriesE lists matching resources.
func ListOCIRepositoriesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) ([]sourcev1.OCIRepository, error) {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return nil, err
	}

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := context.Background()
	var repos sourcev1.OCIRepositoryList
	err = fluxclient.List(ctx, &repos, opts...)
	if err != nil {
		return nil, err
	}

	return repos.Items, nil
}

// WaitForOCIRepositoryReady waits until the specified Flux OCIRepository resource becomes Ready within the given timeout.
// It polls the resource status at regular intervals and fails the test if the resource does not become Ready in time.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for the Kubernetes cluster.
//	name     - The name of the OCIRepository resource.
//	namespace- The namespace where the OCIRepository resource is located.
//	timeout  - The maximum duration to wait for the resource to become Ready.
//
// Fails the test if the OCIRepository does not reach the Ready condition within the timeout.
// WaitForOCIRepositoryReady waits for the resource condition to be satisfied.
func WaitForOCIRepositoryReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForOCIRepositoryReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "OCIRepository %s/%s did not become Ready", namespace, name)
}

// WaitForOCIRepositoryReadyE waits for the resource condition to be satisfied.
func WaitForOCIRepositoryReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var ocirepo sourcev1.OCIRepository
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &ocirepo)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(ocirepo.Status.Conditions), nil
	})
}

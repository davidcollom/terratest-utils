package flux

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ListKustomization retrieves all Flux Kustomization resources in the specified namespace.
// It uses the provided testing context and kubectl options to create a Flux client,
// then lists the Kustomizations within the given namespace. The function fails the test
// if the client cannot be created or if listing the Kustomizations fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list Kustomizations.
//
// Returns:
//   - A slice of kustomizev1.Kustomization objects found in the specified namespace.
//
// ListKustomization lists matching resources.
func ListKustomization(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []kustomizev1.Kustomization {
	kustomizations, err := ListKustomizationE(t, options, namespace, opts...)
	require.NoError(t, err, "Failed to list Kustomizations in namespace %s", namespace)
	return kustomizations
}

// ListKustomizationE lists matching resources.
func ListKustomizationE(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) ([]kustomizev1.Kustomization, error) {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return nil, err
	}

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := context.Background()
	var kustomizations kustomizev1.KustomizationList
	err = fluxclient.List(ctx, &kustomizations, opts...)
	if err != nil {
		return nil, err
	}

	return kustomizations.Items, nil
}

// WaitForKustomizationReady waits until the specified Flux Kustomization resource reaches the Ready condition within the given timeout.
// It polls the resource status at regular intervals and fails the test if the resource does not become Ready in time.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the Kustomization resource.
//   - namespace: The namespace of the Kustomization resource.
//   - timeout: The maximum duration to wait for the resource to become Ready.
//
// The function will fail the test if the Kustomization does not become Ready within the timeout.
// WaitForKustomizationReady waits for the resource condition to be satisfied.
func WaitForKustomizationReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForKustomizationReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Kustomization %s/%s did not become Ready", namespace, name)
}

// WaitForKustomizationReadyE waits for the resource condition to be satisfied.
func WaitForKustomizationReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	fluxclient, err := NewFluxClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var kust kustomizev1.Kustomization
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &kust)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(kust.Status.Conditions), nil
	})
}

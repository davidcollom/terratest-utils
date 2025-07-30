package flux

import (
	"context"
	"testing"
	"time"

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
func ListKustomization(t *testing.T, options *k8s.KubectlOptions, namespace string, opts ...client.ListOption) []kustomizev1.Kustomization {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	// Append the namespace to the list options
	opts = append(opts, client.InNamespace(namespace))

	ctx := t.Context()
	var kustomizations kustomizev1.KustomizationList
	err = fluxclient.List(ctx, &kustomizations, opts...)
	require.NoError(t, err, "Failed to list Kustomizations in namespace %s", namespace)

	return kustomizations.Items
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
func WaitForKustomizationReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(t, options)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var kust kustomizev1.Kustomization
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &kust)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(kust.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("Kustomization %s/%s did not become Ready: %v", namespace, name, err)
	}
}

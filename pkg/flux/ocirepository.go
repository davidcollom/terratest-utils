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
func WaitForOCIRepositoryReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	fluxclient, err := NewFluxClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Flux client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {

		var ocirepo sourcev1.OCIRepository
		err = fluxclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &ocirepo)
		if err != nil {
			return false, nil
		}
		return hasReadyCondition(ocirepo.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("OCIRepository %s/%s did not become Ready: %v", namespace, name, err)
	}
}

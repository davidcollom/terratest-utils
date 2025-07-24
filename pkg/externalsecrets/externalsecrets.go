package externalsecrets

import (
	"context"
	"testing"
	"time"

	esov1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForExternalSecretReady waits until the specified ExternalSecret resource in the given namespace
// becomes ready within the provided timeout duration. It polls the resource status at regular intervals
// and fails the test if the resource does not become ready in time.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for Kubernetes client.
//	name     - The name of the ExternalSecret resource.
//	namespace- The namespace where the ExternalSecret is located.
//	timeout  - The maximum duration to wait for the resource to become ready.
//
// The function uses the External Secrets Operator client to fetch the resource and checks its readiness
// using IsExternalSecretReady. If the resource does not become ready within the timeout, the test fails.
func WaitForExternalSecretReady(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	esoclient, err := NewESOClient(options.RestConfig)
	require.NoError(t, err, "Unable to create External Secrets client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var eso *esov1.ExternalSecret
		err := esoclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, eso)
		if err != nil {
			return false, nil
		}

		if IsExternalSecretReady(eso) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Application %s/%s did not become Healthy & Synced: %v", namespace, name, err)
	}
}

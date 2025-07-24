package externalsecrets

import (
	"context"
	"testing"
	"time"

	esov1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForClusterSecretStoreReady waits until the specified ClusterSecretStore resource is in a "Ready" state.
// It polls the Kubernetes API at regular intervals until the ClusterSecretStore's status condition
// `ReasonStoreValid` is `ConditionTrue`, or until the provided timeout is reached.
// If the ClusterSecretStore does not become ready within the timeout, the test fails.
//
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for Kubernetes API access.
//   - name: The name of the ClusterSecretStore resource.
//   - namespace: The namespace of the ClusterSecretStore resource.
//   - timeout: The maximum duration to wait for the ClusterSecretStore to become ready.
//
// This function is intended for use in integration tests to ensure that ClusterSecretStore resources
// are fully initialized before proceeding.
func WaitForClusterSecretStoreReady(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	esoclient, err := NewESOClient(options.RestConfig)
	require.NoError(t, err, "Unable to create External Secrets client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var store esov1.ClusterSecretStore
		err := esoclient.Get(context.TODO(), ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &store)
		if err != nil {
			t.Logf("SecretStore %s/%s not yet available: %v", namespace, name, err)
			return false, nil // keep retrying
		}
		for _, cond := range store.Status.Conditions {
			if cond.Type == esov1.ReasonStoreValid && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		t.Fatalf("SecretStore %s/%s did not become Ready: %v", namespace, name, err)
	}
}

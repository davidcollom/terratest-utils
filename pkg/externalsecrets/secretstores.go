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

// WaitForSecretStoreReady waits until the specified SecretStore resource in the given namespace becomes Ready.
// It polls the SecretStore status at regular intervals until the Ready condition is met or the timeout is reached.
// If the SecretStore does not become Ready within the timeout, the test fails.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for Kubernetes API access.
//	name     - The name of the SecretStore resource.
//	namespace- The namespace where the SecretStore is located.
//	timeout  - The maximum duration to wait for the SecretStore to become Ready.
//
// This function requires the External Secrets Operator client to be available and the SecretStore resource to be present.
func WaitForSecretStoreReady(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	esoclient, err := NewESOClient(options.RestConfig)
	require.NoError(t, err, "Unable to create External Secrets client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var store esov1.SecretStore
		err := esoclient.Get(context.TODO(), ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &store)
		if err != nil {
			t.Logf("SecretStore %s/%s not yet available: %v", namespace, name, err)
			return false, nil // keep retrying
		}
		for _, cond := range store.Status.Conditions {
			if cond.Type == esov1.SecretStoreReady && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		t.Fatalf("SecretStore %s/%s did not become Ready: %v", namespace, name, err)
	}
}

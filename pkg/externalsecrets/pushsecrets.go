package externalsecrets

import (
	"context"
	"testing"
	"time"

	esov1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitForPushSecretReady waits until the specified PushSecret resource in the given namespace becomes Ready within the provided timeout.
// It polls the Kubernetes API at regular intervals to check the status of the PushSecret's conditions.
// If the PushSecret does not become Ready within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for Kubernetes API access.
//   - name: The name of the PushSecret resource.
//   - namespace: The namespace where the PushSecret is located.
//   - timeout: The maximum duration to wait for the PushSecret to become Ready.
func WaitForPushSecretReady(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()
	esoclient, err := NewESOClient(options.RestConfig)
	require.NoError(t, err, "Unable to create External Secrets client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var ps esov1alpha1.PushSecret
		err := esoclient.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &ps)
		if err != nil {
			t.Logf("PushSecret %s/%s not found yet: %v", namespace, name, err)
			return false, nil
		}
		return hasReadyCondition(ps.Status.Conditions), nil
	})

	if err != nil {
		t.Fatalf("PushSecret %s/%s did not become Ready: %v", namespace, name, err)
	}
}

func hasReadyCondition(conds []esov1alpha1.PushSecretStatusCondition) bool {
	for _, cond := range conds {
		if cond.Type == esov1alpha1.PushSecretReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

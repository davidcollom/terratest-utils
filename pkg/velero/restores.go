package velero

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitForRestoreCompleted waits until a Velero Restore resource reaches the "Completed" phase or the specified timeout is reached.
// It polls the status of the Restore resource every 2 seconds. If the Restore does not reach the "Completed" phase within the timeout,
// the test fails with a fatal error. If the Restore resource is not found during polling, it logs a retry message and continues polling.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options containing the Kubernetes REST config.
//   - name: The name of the Velero Restore resource.
//   - namespace: The namespace where the Restore resource is located.
//   - timeout: The maximum duration to wait for the Restore to complete.
func WaitForRestoreCompleted(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()
	client, err := NewVeleroClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Velero client")
	ctx := t.Context()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var restore velerov1.Restore
		err := client.Get(ctx, key, &restore)
		if err != nil {
			t.Logf("Retrying: Restore %s/%s not found: %v", namespace, name, err)
			return false, nil
		}
		return restore.Status.Phase == velerov1.RestorePhaseCompleted, nil
	})

	if err != nil {
		t.Fatalf("Restore %s/%s did not complete: %v", namespace, name, err)
	}
}

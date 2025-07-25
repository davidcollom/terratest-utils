package velero

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func WaitForBackupStorageLocationReady(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewVeleroClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Velero client")
	ctx := t.Context()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var bsl velerov1.BackupStorageLocation
		err := client.Get(ctx, key, &bsl)
		if err != nil {
			t.Logf("Retrying: BSL %s/%s not found: %v", namespace, name, err)
			return false, nil
		}
		return bsl.Status.Phase == velerov1.BackupStorageLocationPhaseAvailable, nil
	})

	if err != nil {
		t.Fatalf("BackupStorageLocation %s/%s did not become Available: %v", namespace, name, err)
	}
}

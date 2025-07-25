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

// ListBackups retrieves a list of Velero Backup resources in the specified namespace.
// It uses the provided testing context and Kubernetes options to create a Velero client,
// then lists all Backup objects within the given namespace. The function fails the test
// if the client cannot be created or if the list operation fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The Kubernetes KubectlOptions to use for client configuration.
//   - namespace: The namespace from which to list Velero Backups.
//
// Returns:
//   - A slice of velerov1.Backup objects found in the specified namespace.
func ListBackups(t *testing.T, options *k8s.KubectlOptions, namespace string) []velerov1.Backup {
	t.Helper()

	client, err := NewVeleroClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Velero client")

	ctx := t.Context()
	var backups velerov1.BackupList
	err = client.List(ctx, &backups, ctrlclient.InNamespace(namespace))
	require.NoError(t, err, "Failed to list Backups in namespace %s", namespace)

	return backups.Items
}

// WaitForBackupSucceeded waits until the specified Velero backup reaches the "Completed" phase or the timeout is reached.
// It polls the backup status every 2 seconds and fails the test if the backup does not complete successfully within the given timeout.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options containing the Kubernetes REST config.
//   - name: The name of the Velero backup to check.
//   - namespace: The namespace where the backup resides.
//   - timeout: The maximum duration to wait for the backup to complete.
//
// This function will call t.Fatalf if the backup does not complete successfully within the timeout.
func WaitForBackupSucceeded(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewVeleroClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Velero client")
	ctx := t.Context()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var backup velerov1.Backup
		err := client.Get(ctx, key, &backup)
		if err != nil {
			t.Logf("Retrying: Backup %s/%s not found: %v", namespace, name, err)
			return false, nil
		}
		return backup.Status.Phase == velerov1.BackupPhaseCompleted, nil
	})

	if err != nil {
		t.Fatalf("Backup %s/%s did not complete successfully: %v", namespace, name, err)
	}
}

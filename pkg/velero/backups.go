package velero

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

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
//
// ListBackups lists matching resources.
func ListBackups(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []velerov1.Backup {
	backups, err := ListBackupsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Backups in namespace %s", namespace)
	return backups
}

// ListBackupsE lists matching resources.
func ListBackupsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]velerov1.Backup, error) {
	client, err := NewVeleroClient(options.RestConfig)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var backups velerov1.BackupList
	err = client.List(ctx, &backups, ctrlclient.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	return backups.Items, nil
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
// WaitForBackupSucceeded waits for the resource condition to be satisfied.
func WaitForBackupSucceeded(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForBackupSucceededE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Backup %s/%s did not complete successfully", namespace, name)
}

// WaitForBackupSucceededE waits for the resource condition to be satisfied.
func WaitForBackupSucceededE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewVeleroClient(options.RestConfig)
	if err != nil {
		return err
	}
	ctx := context.Background()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var backup velerov1.Backup
		err := client.Get(ctx, key, &backup)
		if err != nil {
			fmt.Printf("Retrying: Backup %s/%s not found: %v\n", namespace, name, err)
			return false, nil
		}
		return backup.Status.Phase == velerov1.BackupPhaseCompleted, nil
	})
}

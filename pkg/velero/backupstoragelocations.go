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

// ListBackupStorageLocation retrieves a list of Velero BackupStorageLocation resources in the specified namespace.
// It uses the provided testing context and Kubernetes options to create a Velero client and perform the list operation.
// The function fails the test if the client cannot be created or if the list operation fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options containing the Kubernetes REST config.
//   - namespace: The namespace from which to list BackupStorageLocations.
//
// Returns:
//   - A slice of velerov1.BackupStorageLocation objects found in the specified namespace.
func ListBackupStorageLocation(t *testing.T, options *k8s.KubectlOptions, namespace string) []velerov1.BackupStorageLocation {
	t.Helper()

	client, err := NewVeleroClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Velero client")

	ctx := t.Context()
	var bsl velerov1.BackupStorageLocationList
	err = client.List(ctx, &bsl, ctrlclient.InNamespace(namespace))
	require.NoError(t, err, "Failed to list BackupStorageLocations in namespace %s", namespace)

	return bsl.Items
}

// WaitForBackupStorageLocationReady waits until the specified Velero BackupStorageLocation resource
// reaches the "Available" phase or the provided timeout is reached. It polls the resource status
// every 2 seconds. If the resource does not become available within the timeout, the test fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options containing the Kubernetes REST config.
//   - name: The name of the BackupStorageLocation resource.
//   - namespace: The namespace of the BackupStorageLocation resource.
//   - timeout: The maximum duration to wait for the resource to become available.
//
// This function is intended for use in integration or end-to-end tests to ensure that
// a Velero BackupStorageLocation is ready before proceeding.
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

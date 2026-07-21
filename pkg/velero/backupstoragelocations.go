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
//
// ListBackupStorageLocation lists matching resources.
func ListBackupStorageLocation(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []velerov1.BackupStorageLocation {
	locations, err := ListBackupStorageLocationE(t, options, namespace)
	require.NoError(t, err, "Failed to list BackupStorageLocations in namespace %s", namespace)
	return locations
}

// ListBackupStorageLocationE lists matching resources.
func ListBackupStorageLocationE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]velerov1.BackupStorageLocation, error) {
	client, err := NewVeleroClient(options.RestConfig)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var bsl velerov1.BackupStorageLocationList
	err = client.List(ctx, &bsl, ctrlclient.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	return bsl.Items, nil
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
// WaitForBackupStorageLocationReady waits for the resource condition to be satisfied.
func WaitForBackupStorageLocationReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForBackupStorageLocationReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "BackupStorageLocation %s/%s did not become Available", namespace, name)
}

// WaitForBackupStorageLocationReadyE waits for the resource condition to be satisfied.
func WaitForBackupStorageLocationReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewVeleroClient(options.RestConfig)
	if err != nil {
		return err
	}
	ctx := context.Background()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var bsl velerov1.BackupStorageLocation
		err := client.Get(ctx, key, &bsl)
		if err != nil {
			fmt.Printf("Retrying: BSL %s/%s not found: %v\n", namespace, name, err)
			return false, nil
		}
		return bsl.Status.Phase == velerov1.BackupStorageLocationPhaseAvailable, nil
	})
}

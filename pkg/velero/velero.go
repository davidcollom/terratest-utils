// Package velero provides Terratest-style helpers for testing Velero backups,
// restores, and storage configurations. Helpers include wait functions for
// BackupStorageLocations, Backups, and Restores using status conditions and phases.
package velero

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// NewVeleroClient creates a new controller-runtime client for Velero resources using the provided Kubernetes REST config.
func NewVeleroClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = velerov1.AddToScheme(scheme)
	return client.New(cfg, client.Options{Scheme: scheme})
}

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

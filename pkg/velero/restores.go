package velero

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ListRestores retrieves a list of Velero Restore resources in the specified namespace.
// It uses the provided testing context and Kubernetes options to create a Velero client,
// then lists all Restore objects within the given namespace. The function fails the test
// if the client cannot be created or if the list operation fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The Kubernetes options containing the REST config.
//   - namespace: The namespace from which to list Restore resources.
//
// Returns:
//   - A slice of velerov1.Restore objects found in the specified namespace.
//
// ListRestores lists matching resources.
func ListRestores(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []velerov1.Restore {
	restores, err := ListRestoresE(t, options, namespace)
	require.NoError(t, err, "Failed to list Restores in namespace %s", namespace)
	return restores
}

// ListRestoresE lists matching resources.
func ListRestoresE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]velerov1.Restore, error) {
	client, err := NewVeleroClient(options.RestConfig)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var restores velerov1.RestoreList
	err = client.List(ctx, &restores, ctrlclient.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	return restores.Items, nil
}

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
//
// WaitForRestoreCompleted waits for the resource condition to be satisfied.
func WaitForRestoreCompleted(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForRestoreCompletedE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Restore %s/%s did not complete", namespace, name)
}

// WaitForRestoreCompletedE waits for the resource condition to be satisfied.
func WaitForRestoreCompletedE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewVeleroClient(options.RestConfig)
	if err != nil {
		return err
	}
	ctx := context.Background()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var restore velerov1.Restore
		err := client.Get(ctx, key, &restore)
		if err != nil {
			fmt.Printf("Retrying: Restore %s/%s not found: %v\n", namespace, name, err)
			return false, nil
		}
		return restore.Status.Phase == velerov1.RestorePhaseCompleted, nil
	})
}

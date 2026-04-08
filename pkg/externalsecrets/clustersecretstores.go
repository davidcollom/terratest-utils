package externalsecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	esov1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListClusterSecretStores retrieves a list of ClusterSecretStore resources from the specified namespace
// using the provided Kubernetes options. It returns a slice of ClusterSecretStore objects.
// The function fails the test if the External Secrets client cannot be created or if the list operation fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list ClusterSecretStores.
//
// Returns:
//   - A slice of esov1.ClusterSecretStore representing the ClusterSecretStores found in the namespace.
// ListClusterSecretStores lists matching resources.
func ListClusterSecretStores(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []esov1.ClusterSecretStore {
	stores, err := ListClusterSecretStoresE(t, options, namespace)
	require.NoError(t, err, "Failed to list ClusterSecretStores in namespace %s", namespace)
	return stores
}

// ListClusterSecretStoresE lists matching resources.
func ListClusterSecretStoresE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]esov1.ClusterSecretStore, error) {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var stores esov1.ClusterSecretStoreList
	err = esoclient.List(ctx, &stores, ctrlclient.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	return stores.Items, nil
}

// WaitForClusterSecretStoreReady waits until the specified ClusterSecretStore resource is in a "Ready" state.
// It polls the Kubernetes API at regular intervals until the ClusterSecretStore's status condition
// `ReasonStoreValid` is `ConditionTrue`, or until the provided timeout is reached.
// If the ClusterSecretStore does not become ready within the timeout, the test fails.
//
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for Kubernetes API access.
//   - name: The name of the ClusterSecretStore resource.
//   - namespace: The namespace of the ClusterSecretStore resource.
//   - timeout: The maximum duration to wait for the ClusterSecretStore to become ready.
//
// This function is intended for use in integration tests to ensure that ClusterSecretStore resources
// are fully initialized before proceeding.
// WaitForClusterSecretStoreReady waits for the resource condition to be satisfied.
func WaitForClusterSecretStoreReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForClusterSecretStoreReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "SecretStore %s/%s did not become Ready", namespace, name)
}

// WaitForClusterSecretStoreReadyE waits for the resource condition to be satisfied.
func WaitForClusterSecretStoreReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var store esov1.ClusterSecretStore
		err := esoclient.Get(context.TODO(), ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &store)
		if err != nil {
			fmt.Printf("SecretStore %s/%s not yet available: %v\n", namespace, name, err)
			return false, nil // keep retrying
		}
		for _, cond := range store.Status.Conditions {
			if cond.Type == esov1.ReasonStoreValid && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
}

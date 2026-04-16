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

// ListSecretStores retrieves a list of External Secrets SecretStore resources from the specified Kubernetes namespace.
// It uses the provided testing context and kubectl options to create an External Secrets client and perform the list operation.
// The function fails the test if the client cannot be created or if the list operation encounters an error.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace from which to list SecretStores.
//
// Returns:
//   - A slice of esov1.SecretStore objects found in the specified namespace.
//
// ListSecretStores lists matching resources.
func ListSecretStores(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []esov1.SecretStore {
	stores, err := ListSecretStoresE(t, options, namespace)
	require.NoError(t, err, "Failed to list SecretStores in namespace %s", namespace)
	return stores
}

// ListSecretStoresE lists matching resources.
func ListSecretStoresE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]esov1.SecretStore, error) {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var stores esov1.SecretStoreList
	err = esoclient.List(ctx, &stores, ctrlclient.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	return stores.Items, nil
}

// WaitForSecretStoreReady waits until the specified SecretStore resource in the given namespace becomes Ready.
// It polls the SecretStore status at regular intervals until the Ready condition is met or the timeout is reached.
// If the SecretStore does not become Ready within the timeout, the test fails.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for Kubernetes API access.
//	name     - The name of the SecretStore resource.
//	namespace- The namespace where the SecretStore is located.
//	timeout  - The maximum duration to wait for the SecretStore to become Ready.
//
// This function requires the External Secrets Operator client to be available and the SecretStore resource to be present.
// WaitForSecretStoreReady waits for the resource condition to be satisfied.
func WaitForSecretStoreReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForSecretStoreReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "SecretStore %s/%s did not become Ready", namespace, name)
}

// WaitForSecretStoreReadyE waits for the resource condition to be satisfied.
func WaitForSecretStoreReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var store esov1.SecretStore
		err := esoclient.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &store)
		if err != nil {
			fmt.Printf("SecretStore %s/%s not yet available: %v\n", namespace, name, err)
			return false, nil // keep retrying
		}
		for _, cond := range store.Status.Conditions {
			if cond.Type == esov1.SecretStoreReady && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
}

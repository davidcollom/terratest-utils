// Package externalsecrets provides Terratest-style helpers for testing External Secrets Operator (ESO)
// resources including SecretStores, ClusterSecretStores, ExternalSecrets, and PushSecrets.
// Helpers wait for readiness conditions and validate that secrets have been reconciled properly.
package externalsecrets

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	esov1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListExternalSecrets retrieves all ExternalSecret resources in the specified namespace using the provided
// Kubernetes options. It returns a slice of ExternalSecret objects. The function fails the test if the client
// cannot be created or if listing the ExternalSecrets fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list ExternalSecrets.
//
// Returns:
//   - A slice of ExternalSecret objects found in the specified namespace.
//
// ListExternalSecrets lists matching resources.
func ListExternalSecrets(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []esov1.ExternalSecret {
	secrets, err := ListExternalSecretsE(t, options, namespace)
	require.NoError(t, err, "Failed to list ExternalSecrets in namespace %s", namespace)
	return secrets
}

// ListExternalSecretsE lists matching resources.
func ListExternalSecretsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]esov1.ExternalSecret, error) {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var secrets esov1.ExternalSecretList
	err = esoclient.List(ctx, &secrets, client.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	return secrets.Items, nil
}

// WaitForExternalSecretReady waits until the specified ExternalSecret resource in the given namespace
// becomes ready within the provided timeout duration. It polls the resource status at regular intervals
// and fails the test if the resource does not become ready in time.
//
// Parameters:
//
//	t        - The testing context.
//	options  - Kubectl options containing the REST config for Kubernetes client.
//	name     - The name of the ExternalSecret resource.
//	namespace- The namespace where the ExternalSecret is located.
//	timeout  - The maximum duration to wait for the resource to become ready.
//
// The function uses the External Secrets Operator client to fetch the resource and checks its readiness
// using IsExternalSecretReady. If the resource does not become ready within the timeout, the test fails.
// WaitForExternalSecretReady waits for the resource condition to be satisfied.
func WaitForExternalSecretReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForExternalSecretReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "ExternalSecret %s/%s did not become Ready", namespace, name)
}

// WaitForExternalSecretReadyE waits for the resource condition to be satisfied.
func WaitForExternalSecretReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var eso esov1.ExternalSecret
		err := esoclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &eso)
		if err != nil {
			return false, nil
		}

		if IsExternalSecretReady(eso.Status) {
			return true, nil
		}
		return false, nil
	})
}

// IsExternalSecretReady checks if the provided ExternalSecret resource has a condition
// of type ExternalSecretReady with a status of ConditionTrue, indicating that the
// external secret is ready. It returns true if such a condition is found, otherwise false.
//
// Parameters:
//
//	sec - Pointer to an esov1.ExternalSecret resource.
//
// Returns:
//
//	bool - true if the ExternalSecret is ready, false otherwise.
//
// IsExternalSecretReady returns whether the resource matches the expected state.
func IsExternalSecretReady(secStatus esov1.ExternalSecretStatus) bool {
	for _, condition := range secStatus.Conditions {
		if condition.Type == esov1.ExternalSecretReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

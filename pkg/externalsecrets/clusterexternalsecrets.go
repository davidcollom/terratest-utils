package externalsecrets

import (
	"context"
	"testing"
	"time"

	esov1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListClusterExternalSecrets retrieves a list of ClusterExternalSecret resources from the specified namespace
// using the provided Kubernetes KubectlOptions. It returns a slice of ClusterExternalSecret objects.
// The function requires a testing.T instance for error handling and test context propagation.
// It fails the test if the External Secrets client cannot be created or if listing the secrets fails.
//
// Parameters:
//   - t:        The testing.T instance used for test context and assertions.
//   - options:  The KubectlOptions specifying the Kubernetes context and configuration.
//   - namespace: The namespace from which to list ClusterExternalSecrets.
//
// Returns:
//   - []esov1.ClusterExternalSecret: A slice containing the ClusterExternalSecret resources found in the namespace.
func ListClusterExternalSecrets(t *testing.T, options *k8s.KubectlOptions, namespace string) []esov1.ClusterExternalSecret {
	t.Helper()

	esoclient, err := NewESOClient(t, options)
	require.NoError(t, err, "Unable to create External Secrets client")

	ctx := t.Context()
	var secrets esov1.ClusterExternalSecretList
	err = esoclient.List(ctx, &secrets, client.InNamespace(namespace))
	require.NoError(t, err, "Failed to list ClusterExternalSecrets in namespace %s", namespace)

	return secrets.Items
}

// WaitForClusterExternalSecretReady waits until the specified ClusterExternalSecret resource in the given namespace
// becomes ready within the provided timeout duration. It polls the resource status at regular intervals and fails the test
// if the resource does not become ready in time. This function requires a valid External Secrets client and uses the
// provided k8s.KubectlOptions for cluster access.
//
// Parameters:
//   - t: The testing context.
//   - options: Kubernetes KubectlOptions containing cluster access configuration.
//   - name: The name of the ClusterExternalSecret resource.
//   - namespace: The namespace where the resource is located.
//   - timeout: The maximum duration to wait for the resource to become ready.
//
// Fails the test if the ClusterExternalSecret does not become ready within the timeout.
func WaitForClusterExternalSecretReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	esoclient, err := NewESOClient(t, options)
	require.NoError(t, err, "Unable to create External Secrets client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var eso esov1.ClusterExternalSecret
		err := esoclient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &eso)
		if err != nil {
			return false, nil
		}

		if IsClusterExternalSecretReady(eso.Status) {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Application %s/%s did not become Healthy & Synced: %v", namespace, name, err)
	}
}

// IsClusterExternalSecretReady checks if the ClusterExternalSecret resource is in a ready state.
// It returns true if any of the conditions in the provided ClusterExternalSecretStatus
// has a type of ClusterExternalSecretReady and a status of ConditionTrue, indicating readiness.
// Otherwise, it returns false.
func IsClusterExternalSecretReady(secStatus esov1.ClusterExternalSecretStatus) bool {
	for _, condition := range secStatus.Conditions {
		if condition.Type == esov1.ClusterExternalSecretReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

package externalsecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	esov1alpha1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ListPushSecrets retrieves all PushSecret resources in the specified Kubernetes namespace.
// It uses the provided testing context and KubectlOptions to create an External Secrets Operator (ESO) client,
// then lists all PushSecrets within the given namespace. The function fails the test if any error occurs during
// client creation or resource listing.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace to search for PushSecrets.
//
// Returns:
//   - A slice of PushSecret resources found in the specified namespace.
// ListPushSecrets lists matching resources.
func ListPushSecrets(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...ctrlclient.ListOption) []esov1alpha1.PushSecret {
	pushSecrets, err := ListPushSecretsE(t, options, namespace, opts...)
	require.NoError(t, err, "Failed to list PushSecrets in namespace %s", namespace)
	return pushSecrets
}

// ListPushSecretsE lists matching resources.
func ListPushSecretsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string, opts ...ctrlclient.ListOption) ([]esov1alpha1.PushSecret, error) {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return nil, err
	}

	// Append the namespace to the list options.
	opts = append(opts, ctrlclient.InNamespace(namespace))

	ctx := context.Background()
	var pushSecrets esov1alpha1.PushSecretList
	err = esoclient.List(ctx, &pushSecrets, opts...)
	if err != nil {
		return nil, err
	}

	return pushSecrets.Items, nil
}

// WaitForPushSecretReady waits until the specified PushSecret resource in the given namespace becomes Ready within the provided timeout.
// It polls the Kubernetes API at regular intervals to check the status of the PushSecret's conditions.
// If the PushSecret does not become Ready within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for Kubernetes API access.
//   - name: The name of the PushSecret resource.
//   - namespace: The namespace where the PushSecret is located.
//   - timeout: The maximum duration to wait for the PushSecret to become Ready.
// WaitForPushSecretReady waits for the resource condition to be satisfied.
func WaitForPushSecretReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForPushSecretReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "PushSecret %s/%s did not become Ready", namespace, name)
}

// WaitForPushSecretReadyE waits for the resource condition to be satisfied.
func WaitForPushSecretReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	esoclient, err := NewESOClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var ps esov1alpha1.PushSecret
		err := esoclient.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: namespace}, &ps)
		if err != nil {
			fmt.Printf("PushSecret %s/%s not found yet: %v\n", namespace, name, err)
			return false, nil
		}
		return hasReadyCondition(ps.Status.Conditions), nil
	})
}

// hasReadyCondition checks if the provided slice of PushSecretStatusCondition contains
// a condition of type PushSecretReady with a status of ConditionTrue. It returns true
// if such a condition is found, otherwise it returns false.
func hasReadyCondition(conds []esov1alpha1.PushSecretStatusCondition) bool {
	for _, cond := range conds {
		if cond.Type == esov1alpha1.PushSecretReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

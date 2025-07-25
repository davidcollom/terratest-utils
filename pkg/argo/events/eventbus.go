package events

import (
	"context"
	"testing"
	"time"

	argoeventsv1alpha1 "github.com/argoproj/argo-events/pkg/apis/events/v1alpha1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func ListEventBuses(t *testing.T, options *k8s.KubectlOptions, namespace string) []argoeventsv1alpha1.EventBus {
	t.Helper()

	client, err := NewArgoEventsClient(t, options)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	eventBusList, err := client.ArgoprojV1alpha1().EventBus(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list EventBuses in namespace %s", namespace)

	return eventBusList.Items
}

// WaitForEventBusReady waits until the specified Argo Events EventBus resource is Ready, or times out.
// Useful for integration tests to ensure event infrastructure is available before proceeding.
func WaitForEventBusReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewArgoEventsClient(t, options)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		eventBus, err := client.ArgoprojV1alpha1().EventBus(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		var (
			configured = false
			deployed   = false
		)

		for _, cond := range eventBus.Status.Conditions {
			if cond.Type == argoeventsv1alpha1.EventBusConditionDeployed && cond.IsTrue() {
				deployed = true
			}
			if cond.Type == argoeventsv1alpha1.EventBusConditionConfigured && cond.IsTrue() {
				configured = true
			}
		}
		return configured && deployed, nil
	})

	if err != nil {
		t.Fatalf("EventBus %s/%s did not become Ready: %v", namespace, name, err)
	}
}

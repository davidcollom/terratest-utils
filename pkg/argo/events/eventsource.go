package events

import (
	"context"
	"testing"
	"time"

	argoeventsv1alpha1 "github.com/argoproj/argo-events/pkg/apis/events/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
)

func ListEventSources(t *testing.T, options *k8s.KubectlOptions, namespace string) []argoeventsv1alpha1.EventSource {
	t.Helper()

	client, err := NewArgoEventsClient(t, options)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	eventSourceList, err := client.ArgoprojV1alpha1().EventSources(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list EventSources in namespace %s", namespace)

	return eventSourceList.Items
}

// WaitForEventSourceReady waits until the specified Argo Events EventSource resource is Ready, or times out.
// Useful for integration tests to ensure event sources are available before proceeding.
func WaitForEventSourceReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewArgoEventsClient(t, options)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		es, err := client.ArgoprojV1alpha1().EventSources(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil // keep retrying
		}

		var (
			deployed   = false
			hasSources = false
		)

		for _, cond := range es.Status.Conditions {
			if cond.Type == argoeventsv1alpha1.EventSourceConditionDeployed && cond.IsTrue() {
				deployed = true
			}
			if cond.Type == argoeventsv1alpha1.EventSourceConditionSourcesProvided && cond.IsTrue() {
				hasSources = true
			}
		}
		return deployed && hasSources, nil
	})

	if err != nil {
		t.Fatalf("EventSource %s/%s did not become Ready: %v", namespace, name, err)
	}
}

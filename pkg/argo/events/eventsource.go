package events

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	argoeventsv1alpha1 "github.com/argoproj/argo-events/pkg/apis/events/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
)

// ListEventSources retrieves a list of Argo EventSource resources from the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Events client,
// then lists all EventSources in the given namespace. The function fails the test if any error occurs.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list EventSources.
//
// Returns:
//   - A slice of argoeventsv1alpha1.EventSource objects representing the EventSources in the namespace.
// ListEventSources lists matching resources.
func ListEventSources(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []argoeventsv1alpha1.EventSource {
	eventSources, err := ListEventSourcesE(t, options, namespace)
	require.NoError(t, err, "Failed to list EventSources in namespace %s", namespace)
	return eventSources
}

// ListEventSourcesE lists matching resources.
func ListEventSourcesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]argoeventsv1alpha1.EventSource, error) {
	client, err := NewArgoEventsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	eventSourceList, err := client.ArgoprojV1alpha1().EventSources(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return eventSourceList.Items, nil
}

// WaitForEventSourceReady waits until the specified Argo Events EventSource resource is Ready, or times out.
// Useful for integration tests to ensure event sources are available before proceeding.
// WaitForEventSourceReady waits for the resource condition to be satisfied.
func WaitForEventSourceReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForEventSourceReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "EventSource %s/%s did not become Ready", namespace, name)
}

// WaitForEventSourceReadyE waits for the resource condition to be satisfied.
func WaitForEventSourceReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoEventsClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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
}

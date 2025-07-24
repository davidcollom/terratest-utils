package events

import (
	"context"
	"testing"
	"time"

	argoeventsv1alpha1 "github.com/argoproj/argo-events/pkg/apis/events/v1alpha1"
	argoclientset "github.com/argoproj/argo-events/pkg/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
)

// WaitForEventSourceReady waits until the specified Argo Events EventSource resource is Ready, or times out.
// Useful for integration tests to ensure event sources are available before proceeding.
func WaitForEventSourceReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := argoclientset.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		es, err := client.ArgoprojV1alpha1().EventSources(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil // keep retrying
		}

		for _, cond := range es.Status.Conditions {
			if cond.Type == argoeventsv1alpha1.ConditionReady && cond.Status == "True" {
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("EventSource %s/%s did not become Ready: %v", namespace, name, err)
	}
}

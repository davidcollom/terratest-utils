// Package rollouts provides Terratest-style helpers for testing Argo Rollouts.
// It includes polling-based utilities for checking rollout phases, pause states,
// and progressive deployment status using the Argo Rollouts clientset.
package rollouts

import (
	"context"
	"testing"
	"time"

	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	rolloutClientSet "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForRolloutHealthy waits until the specified Argo Rollout resource reaches a Healthy phase within the given timeout.
// It polls the rollout status every 2 seconds and checks for the "Progressing" condition with status "True" and phase "Healthy".
// If the rollout does not become healthy within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for the Kubernetes client.
//   - name: The name of the rollout resource.
//   - namespace: The namespace of the rollout resource.
//   - timeout: The maximum duration to wait for the rollout to become healthy.
func WaitForRolloutHealthy(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := rolloutClientSet.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create Argo Rollouts clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		ro, err := client.ArgoprojV1alpha1().Rollouts(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		for _, cond := range ro.Status.Conditions {
			if cond.Type == rolloutsv1alpha1.RolloutProgressing && cond.Status == "True" {
				if ro.Status.Phase == rolloutsv1alpha1.RolloutPhaseHealthy {
					return true, nil
				}
			}
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Rollout %s/%s did not become Healthy in time: %v", namespace, name, err)
	}
}

// WaitForRolloutPaused waits until the specified Argo Rollout resource enters the "Paused" phase within the given timeout.
// It polls the rollout status every 2 seconds using the provided Kubernetes options, rollout name, and namespace.
// If the rollout does not reach the paused phase within the timeout, the test fails with a fatal error.
// Requires a valid Argo Rollouts clientset and test context.
//
// Parameters:
//
//	t        - The testing context.
//	options  - The kubectl options containing REST config for the client.
//	name     - The name of the rollout resource.
//	namespace- The namespace of the rollout resource.
//	timeout  - The maximum duration to wait for the rollout to pause.
func WaitForRolloutPaused(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := rolloutClientSet.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create Argo Rollouts clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		ro, err := client.ArgoprojV1alpha1().Rollouts(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return ro.Status.Phase == rolloutsv1alpha1.RolloutPhasePaused, nil
	})

	if err != nil {
		t.Fatalf("Rollout %s/%s did not pause in time: %v", namespace, name, err)
	}
}

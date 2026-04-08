// Package rollouts provides Terratest-style helpers for testing Argo Rollouts.
// It includes polling-based utilities for checking rollout phases, pause states,
// and progressive deployment status using the Argo Rollouts clientset.
package rollouts

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	rolloutClientSet "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	"github.com/davidcollom/terratest-utils/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"k8s.io/apimachinery/pkg/util/wait"
)

// NewArgoRolloutsClient creates a new Argo Rollouts client using the provided testing context and kubectl options.
// It retrieves the Kubernetes REST configuration from the given options or generates one if not present.
// Returns an Argo Rollouts client interface and an error if the client could not be created.
//
// Parameters:
//   - t: The testing context, used for logging and error handling.
//   - options: The kubectl options containing cluster configuration and optional REST config.
//
// Returns:
//   - rolloutClientSet.Interface: The Argo Rollouts client interface for interacting with Rollouts resources.
//   - error: An error if the client could not be created.
// NewArgoRolloutsClient creates a new client or helper instance.
func NewArgoRolloutsClient(t testing.TestingT, options *k8s.KubectlOptions) (rolloutClientSet.Interface, error) {
	var cfg *rest.Config
	var err error
	if options.RestConfig == nil {
		cfg, err = utils.GetRestConfigE(t, options)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = options.RestConfig
	}

	return rolloutClientSet.NewForConfig(cfg)
}

// ListRollouts retrieves all Argo Rollouts in the specified namespace using the provided kubectl options.
// It requires a testing context and will fail the test if the client cannot be created or the Rollouts cannot be listed.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the cluster.
//   - namespace: The namespace from which to list the Rollouts.
//
// Returns:
//   - A slice of rolloutsv1alpha1.Rollout objects representing the Rollouts in the given namespace.
// ListRollouts lists matching resources.
func ListRollouts(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []rolloutsv1alpha1.Rollout {
	rollouts, err := ListRolloutsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Rollouts in namespace %s", namespace)
	return rollouts
}

// ListRolloutsE lists matching resources.
func ListRolloutsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]rolloutsv1alpha1.Rollout, error) {
	client, err := NewArgoRolloutsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rolloutList, err := client.ArgoprojV1alpha1().Rollouts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return rolloutList.Items, nil
}

// WaitForRolloutHealthy waits until the specified Argo Rollout resource reaches a Healthy phase within the given timeout.
// It polls the rollout status every 2 seconds and checks for the "Progressing" condition with status "True" and phase "Healthy".
// If the rollout does not become healthy within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the REST config for the Kubernetes client.
//   - name: The name of the rollout resource.
//   - namespace: The namespace of the rollout resource.
//   - timeout: The maximum duration to wait for the rollout to become healthy.
// WaitForRolloutHealthy waits for the resource condition to be satisfied.
func WaitForRolloutHealthy(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForRolloutHealthyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Rollout %s/%s did not become Healthy in time", namespace, name)
}

// WaitForRolloutHealthyE waits for the resource condition to be satisfied.
func WaitForRolloutHealthyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoRolloutsClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
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
// WaitForRolloutPaused waits for the resource condition to be satisfied.
func WaitForRolloutPaused(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForRolloutPausedE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Rollout %s/%s did not pause in time", namespace, name)
}

// WaitForRolloutPausedE waits for the resource condition to be satisfied.
func WaitForRolloutPausedE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoRolloutsClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		ro, err := client.ArgoprojV1alpha1().Rollouts(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		return ro.Status.Phase == rolloutsv1alpha1.RolloutPhasePaused, nil
	})
}

package cd

import (
	"context"
	"testing"
	"time"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForApplicationSetHealthyAndSynced waits until the specified Argo CD ApplicationSet in the given namespace
// is healthy and its resources are up to date, or until the provided timeout is reached.
// It polls the ApplicationSet status every 2 seconds, checking for the "ResourcesUpToDate" condition with a "True" status.
// If the ApplicationSet does not become healthy and synced within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the ApplicationSet.
//   - namespace: The namespace where the ApplicationSet resides.
//   - timeout: The maximum duration to wait for the ApplicationSet to become healthy and synced.
func WaitForApplicationSetHealthyAndSynced(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := NewArgoCDClient(t, options)
	require.NoError(t, err, "Unable to create Argo CD client")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		app, err := client.ArgoprojV1alpha1().ApplicationSets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		for _, cond := range app.Status.Conditions {
			if cond.Type == argocdv1alpha1.ApplicationSetConditionResourcesUpToDate && cond.Status == argocdv1alpha1.ApplicationSetConditionStatusTrue {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		t.Fatalf("Application %s/%s did not become Healthy & Synced: %v", namespace, name, err)
	}
}

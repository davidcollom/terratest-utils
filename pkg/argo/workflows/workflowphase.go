package workflows

import (
	"context"
	"testing"
	"time"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	workflowsClientSet "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForWorkflowPhase waits until the specified Argo Workflow reaches the desired phase within the given timeout.
// It polls the workflow status every 2 seconds using the provided Kubernetes options and namespace.
// If the workflow does not reach the desired phase in time, the test fails with a fatal error.
//
// Parameters:
//
//	t            - The testing context.
//	options      - The Kubernetes KubectlOptions containing REST config.
//	name         - The name of the workflow to monitor.
//	namespace    - The namespace where the workflow resides.
//	desiredPhase - The target WorkflowPhase to wait for.
//	timeout      - The maximum duration to wait for the workflow to reach the desired phase.
//
// Fails the test if the workflow does not reach the desired phase within the timeout.
func WaitForWorkflowPhase(t *testing.T, options *k8s.KubectlOptions, name, namespace string, desiredPhase workflowv1alpha1.WorkflowPhase, timeout time.Duration) {
	t.Helper()

	client, err := workflowsClientSet.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		wf, err := client.ArgoprojV1alpha1().Workflows(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		return wf.Status.Phase == desiredPhase, nil
	})

	if err != nil {
		t.Fatalf("Workflow %s/%s did not reach phase %q in time: %v", namespace, name, desiredPhase, err)
	}
}

package workflows

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListWorkflowPhases retrieves the phases of all Argo Workflows in the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Workflows client,
// lists all workflows in the given namespace, and returns a slice containing the phase of each workflow.
//
// Parameters:
//   - t: The testing context used for logging and error handling.
//   - options: The kubectl options used to configure the Kubernetes client.
//   - namespace: The namespace from which to list the workflows.
//
// Returns:
//   - A slice of workflowv1alpha1.WorkflowPhase representing the phase of each workflow in the namespace.
//
// Panics if there is an error creating the client or listing the workflows.
// ListWorkflowPhases lists matching resources.
func ListWorkflowPhases(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowPhase {
	phases, err := ListWorkflowPhasesE(t, options, namespace)
	require.NoError(t, err, "Failed to list Workflow phases in namespace %s", namespace)
	return phases
}

// ListWorkflowPhasesE retrieves the phases of all Argo Workflows in the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListWorkflowPhasesE lists matching resources.
func ListWorkflowPhasesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]workflowv1alpha1.WorkflowPhase, error) {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	workflowList, err := client.ArgoprojV1alpha1().Workflows(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	phases := make([]workflowv1alpha1.WorkflowPhase, 0, len(workflowList.Items))
	for _, wf := range workflowList.Items {
		phases = append(phases, wf.Status.Phase)
	}

	return phases, nil
}

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
// WaitForWorkflowPhase waits for the resource condition to be satisfied.
func WaitForWorkflowPhase(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, desiredPhase workflowv1alpha1.WorkflowPhase, timeout time.Duration) {
	err := WaitForWorkflowPhaseE(t, options, name, namespace, desiredPhase, timeout)
	require.NoError(t, err, "Workflow %s/%s did not reach phase %q in time", namespace, name, desiredPhase)
}

// WaitForWorkflowPhaseE waits until the specified Argo Workflow reaches the desired phase within the given timeout.
func WaitForWorkflowPhaseE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, desiredPhase workflowv1alpha1.WorkflowPhase, timeout time.Duration) error {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		wf, err := client.ArgoprojV1alpha1().Workflows(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		return wf.Status.Phase == desiredPhase, nil
	})
}

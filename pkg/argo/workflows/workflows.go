package workflows

import (
	"testing"
	"time"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

// WaitForWorkflowFailed waits until the specified Argo workflow reaches the "Failed" phase within the given timeout.
// It uses the provided KubectlOptions to interact with the Kubernetes cluster.
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for accessing the cluster.
//   - name: The name of the workflow to monitor.
//   - namespace: The namespace where the workflow resides.
//   - timeout: The maximum duration to wait for the workflow to enter the "Running" phase.
func WaitForWorkflowFailed(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowError, timeout)
}

// WaitForWorkflowRunning waits until the specified Argo workflow reaches the "Running" phase within the given timeout.
// It uses the provided KubectlOptions to interact with the Kubernetes cluster.
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for accessing the cluster.
//   - name: The name of the workflow to monitor.
//   - namespace: The namespace where the workflow resides.
//   - timeout: The maximum duration to wait for the workflow to enter the "Running" phase.
func WaitForWorkflowRunning(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowRunning, timeout)
}

// WaitForWorkflowPending waits until the specified Argo workflow reaches the "Pending" phase within the given timeout.
// It uses the provided testing context, kubectl options, workflow name, and namespace.
// If the workflow does not reach the "Pending" phase within the timeout, the test will fail.
func WaitForWorkflowPending(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowPending, timeout)
}

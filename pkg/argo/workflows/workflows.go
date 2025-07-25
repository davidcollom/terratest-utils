// Package workflows provides Terratest-style helpers for testing Argo Workflows.
// These include functions to wait for specific workflow phases (Succeeded, Failed, etc.),
// and utilities to assert workflow conditions using the Argo Workflows clientset.
package workflows

import (
	"testing"
	"time"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	workflowsClientSet "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/davidcollom/terratest-utils/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

// ListArgoWorkflows retrieves all Argo Workflows in the specified namespace using the provided KubectlOptions.
// It returns a slice of Workflow objects. The function will fail the test if it cannot create the client
// or list the workflows.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list the workflows.
//
// Returns:
//   - A slice of Workflow objects present in the specified namespace.
func ListArgoWorkflows(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.Workflow {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowList, err := client.ArgoprojV1alpha1().Workflows(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Workflows in namespace %s", namespace)

	return workflowList.Items
}

// NewArgoWorkflowsClient creates a new Argo Workflows client using the provided testing context and Kubernetes options.
// It returns an implementation of the workflowv1alpha1.Interface for interacting with Argo Workflows resources.
// If the provided KubectlOptions does not include a RestConfig, it attempts to generate one.
// Returns an error if the client cannot be created.
func NewArgoWorkflowsClient(t *testing.T, options *k8s.KubectlOptions) (workflowsClientSet.Interface, error) {
	t.Helper()
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

	return workflowsClientSet.NewForConfig(cfg)
}

// ListWorkflows retrieves all Argo Workflows in the specified namespace using the provided KubectlOptions.
// It requires a testing.T instance for error handling and context propagation.
// The function returns a slice of Workflow objects present in the given namespace.
// If the client creation or workflow listing fails, the test will be failed with an appropriate error message.
//
// Parameters:
//   - t: The testing.T instance used for test context and assertions.
//   - options: The KubectlOptions used to configure access to the Kubernetes cluster.
//   - namespace: The namespace from which to list the workflows.
//
// Returns:
//   - []workflowv1alpha1.Workflow: A slice containing the workflows found in the specified namespace.
func ListWorkflows(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.Workflow {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	workflowList, err := client.ArgoprojV1alpha1().Workflows(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list Workflows in namespace %s", namespace)

	return workflowList.Items
}

// WaitForWorkflowRunning waits until the specified Argo workflow reaches the "Running" phase or the timeout is reached.
// It uses the provided testing context, kubectl options, workflow name, namespace, and timeout duration.
// Fails the test if the workflow does not reach the "Running" phase within the timeout.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - name: The name of the workflow to check.
//   - namespace: The namespace where the workflow resides.
//   - timeout: The maximum duration to wait for the workflow to reach the "Running" phase.
func WaitForWorkflowRunning(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowRunning, timeout)
}

// WaitForWorkflowError waits until the specified Argo workflow reaches the "Error" phase or the timeout is reached.
// It fails the test if the workflow does not enter the "Error" phase within the given duration.
//
// Parameters:
//
//	t         - The testing context.
//	options   - The kubectl options to use for interacting with the Kubernetes cluster.
//	name      - The name of the workflow to monitor.
//	namespace - The namespace where the workflow is running.
//	timeout   - The maximum duration to wait for the workflow to reach the "Error" phase.
func WaitForWorkflowError(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowError, timeout)
}

// WaitForWorkflowPending waits until the specified Argo workflow reaches the "Pending" phase within the given timeout.
// It uses the provided testing context, kubectl options, workflow name, and namespace.
// If the workflow does not reach the "Pending" phase within the timeout, the test will fail.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - name: The name of the workflow to check.
//   - namespace: The namespace where the workflow resides.
//   - timeout: The maximum duration to wait for the workflow to reach the "Pending" phase.
//
// This function delegates to WaitForWorkflowPhase with the "Pending" phase.
func WaitForWorkflowPending(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForWorkflowPhase(t, options, name, namespace, workflowv1alpha1.WorkflowPending, timeout)
}

package workflows

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListArgoWorkflowTaskResults retrieves a list of Argo WorkflowTaskResult resources from the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Workflows client.
// The function fails the test if the client cannot be created or if listing the WorkflowTaskResults fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace from which to list WorkflowTaskResults.
//
// Returns:
//   - A slice of WorkflowTaskResult resources found in the specified namespace.
// ListArgoWorkflowTaskResults lists matching resources.
func ListArgoWorkflowTaskResults(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTaskResult {
	results, err := ListArgoWorkflowTaskResultsE(t, options, namespace)
	require.NoError(t, err, "Failed to list WorkflowTaskResults in namespace %s", namespace)
	return results
}

// ListArgoWorkflowTaskResultsE retrieves a list of Argo WorkflowTaskResult resources from the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListArgoWorkflowTaskResultsE lists matching resources.
func ListArgoWorkflowTaskResultsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]workflowv1alpha1.WorkflowTaskResult, error) {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	workflowTaskResultList, err := client.ArgoprojV1alpha1().WorkflowTaskResults(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return workflowTaskResultList.Items, nil
}

// ListArgoWorkflowTaskSet retrieves all Argo WorkflowTaskSet resources in the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Workflows client,
// then lists the WorkflowTaskSets in the given namespace. If any error occurs during client creation
// or listing, the test will fail. Returns a slice of WorkflowTaskSet objects found in the namespace.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list WorkflowTaskSets.
//
// Returns:
//   - A slice of WorkflowTaskSet objects present in the specified namespace.
// ListArgoWorkflowTaskSet lists matching resources.
func ListArgoWorkflowTaskSet(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTaskSet {
	taskSets, err := ListArgoWorkflowTaskSetE(t, options, namespace)
	require.NoError(t, err, "Failed to list WorkflowTaskSets in namespace %s", namespace)
	return taskSets
}

// ListArgoWorkflowTaskSetE retrieves all Argo WorkflowTaskSet resources in the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListArgoWorkflowTaskSetE lists matching resources.
func ListArgoWorkflowTaskSetE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]workflowv1alpha1.WorkflowTaskSet, error) {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	workflowTaskSetList, err := client.ArgoprojV1alpha1().WorkflowTaskSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return workflowTaskSetList.Items, nil
}

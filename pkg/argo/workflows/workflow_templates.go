package workflows

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
)

// ListArgoWorkflowTemplates retrieves all Argo WorkflowTemplates in the specified namespace.
//
// It uses the provided testing context and kubectl options to create an Argo Workflows client,
// then lists all WorkflowTemplates in the given namespace. The function fails the test if any
// errors occur during client creation or listing.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options for accessing the Kubernetes cluster.
//   - namespace: The namespace from which to list WorkflowTemplates.
//
// Returns:
//   - A slice of WorkflowTemplate objects found in the specified namespace.
//
// ListArgoWorkflowTemplates lists matching resources.
func ListArgoWorkflowTemplates(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.WorkflowTemplate {
	templates, err := ListArgoWorkflowTemplatesE(t, options, namespace)
	require.NoError(t, err, "Failed to list WorkflowTemplates in namespace %s", namespace)
	return templates
}

// ListArgoWorkflowTemplatesE retrieves all Argo WorkflowTemplates in the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListArgoWorkflowTemplatesE lists matching resources.
func ListArgoWorkflowTemplatesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]workflowv1alpha1.WorkflowTemplate, error) {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	workflowTemplateList, err := client.ArgoprojV1alpha1().WorkflowTemplates(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return workflowTemplateList.Items, nil
}

// ListArgoClusterWorkflowTemplates retrieves all Argo ClusterWorkflowTemplates in the specified namespace using the provided KubectlOptions.
// It returns a slice of ClusterWorkflowTemplate objects. The function will fail the test if there is an error creating the client
// or listing the ClusterWorkflowTemplates.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the cluster.
//   - namespace: The namespace to list ClusterWorkflowTemplates from.
//
// Returns:
//   - A slice of ClusterWorkflowTemplate objects.
//
// ListArgoClusterWorkflowTemplates lists matching resources.
func ListArgoClusterWorkflowTemplates(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.ClusterWorkflowTemplate {
	templates, err := ListArgoClusterWorkflowTemplatesE(t, options, namespace)
	require.NoError(t, err, "Failed to list ClusterWorkflowTemplates")
	return templates
}

// ListArgoClusterWorkflowTemplatesE retrieves all Argo ClusterWorkflowTemplates.
// It returns an error to the caller instead of failing the test directly.
// ListArgoClusterWorkflowTemplatesE lists matching resources.
func ListArgoClusterWorkflowTemplatesE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]workflowv1alpha1.ClusterWorkflowTemplate, error) {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	workflowTemplateList, err := client.ArgoprojV1alpha1().ClusterWorkflowTemplates().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return workflowTemplateList.Items, nil
}

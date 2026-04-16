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

// ListCronWorkflows retrieves all Argo CronWorkflows in the specified namespace using the provided kubectl options.
// It returns a slice of CronWorkflow objects. The function fails the test if the client cannot be created or if
// listing the CronWorkflows fails.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options to use for connecting to the Kubernetes cluster.
//   - namespace: The namespace from which to list the CronWorkflows.
//
// Returns:
//   - A slice of workflowv1alpha1.CronWorkflow representing the CronWorkflows found in the namespace.
// ListCronWorkflows lists matching resources.
func ListCronWorkflows(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.CronWorkflow {
	cronWorkflows, err := ListCronWorkflowsE(t, options, namespace)
	require.NoError(t, err, "Failed to list CronWorkflows in namespace %s", namespace)
	return cronWorkflows
}

// ListCronWorkflowsE retrieves all Argo CronWorkflows in the specified namespace.
// It returns an error to the caller instead of failing the test directly.
// ListCronWorkflowsE lists matching resources.
func ListCronWorkflowsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]workflowv1alpha1.CronWorkflow, error) {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	cronWorkflowList, err := client.ArgoprojV1alpha1().CronWorkflows(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return cronWorkflowList.Items, nil
}

// WaitForCronWorkflowActive waits until the specified Argo CronWorkflow reaches the 'Active' phase within the given timeout.
// It uses the provided KubectlOptions, workflow name, and namespace for the check.
// Fails the test if the CronWorkflow does not become active within the timeout.
// WaitForCronWorkflowActive waits for the resource condition to be satisfied.
func WaitForCronWorkflowActive(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForCronWorkflowActiveE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Workflow %s/%s did not reach phase %q in time", namespace, name, workflowv1alpha1.ActivePhase)
}

// WaitForCronWorkflowActiveE waits until the specified Argo CronWorkflow reaches the Active phase.
func WaitForCronWorkflowActiveE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	return WaitForCronWorkflowPhaseE(t, options, name, namespace, workflowv1alpha1.ActivePhase, timeout)
}

// WaitForCronWorkflowStopped waits until the specified Argo CronWorkflow reaches the "Stopped" phase within the given timeout.
// It uses the provided testing context, kubectl options, workflow name, and namespace.
// If the workflow does not reach the "Stopped" phase within the timeout, the test will fail.
// WaitForCronWorkflowStopped waits for the resource condition to be satisfied.
func WaitForCronWorkflowStopped(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForCronWorkflowStoppedE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Workflow %s/%s did not reach phase %q in time", namespace, name, workflowv1alpha1.StoppedPhase)
}

// WaitForCronWorkflowStoppedE waits until the specified Argo CronWorkflow reaches the Stopped phase.
func WaitForCronWorkflowStoppedE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	return WaitForCronWorkflowPhaseE(t, options, name, namespace, workflowv1alpha1.StoppedPhase, timeout)
}

// WaitForCronWorkflowPhase waits until the specified Argo CronWorkflow reaches the desired phase within the given timeout.
// It polls the CronWorkflow status at regular intervals and fails the test if the desired phase is not reached in time.
//
// Parameters:
//
//	t            - The testing context.
//	options      - Kubectl options containing Kubernetes REST config.
//	name         - The name of the CronWorkflow.
//	namespace    - The namespace of the CronWorkflow.
//	desiredPhase - The target phase to wait for.
//	timeout      - The maximum duration to wait for the desired phase.
//
// Fails the test if the CronWorkflow does not reach the desired phase within the timeout.
// WaitForCronWorkflowPhase waits for the resource condition to be satisfied.
func WaitForCronWorkflowPhase(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, desiredPhase workflowv1alpha1.CronWorkflowPhase, timeout time.Duration) {
	err := WaitForCronWorkflowPhaseE(t, options, name, namespace, desiredPhase, timeout)
	require.NoError(t, err, "Workflow %s/%s did not reach phase %q in time", namespace, name, desiredPhase)
}

// WaitForCronWorkflowPhaseE waits until the specified Argo CronWorkflow reaches the desired phase within the given timeout.
func WaitForCronWorkflowPhaseE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, desiredPhase workflowv1alpha1.CronWorkflowPhase, timeout time.Duration) error {
	client, err := NewArgoWorkflowsClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		wf, err := client.ArgoprojV1alpha1().CronWorkflows(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		return wf.Status.Phase == desiredPhase, nil
	})
}

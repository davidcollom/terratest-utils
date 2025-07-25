package workflows

import (
	"context"
	"testing"
	"time"

	workflowv1alpha1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
)

func ListCronWorkflows(t *testing.T, options *k8s.KubectlOptions, namespace string) []workflowv1alpha1.CronWorkflow {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	cronWorkflowList, err := client.ArgoprojV1alpha1().CronWorkflows(namespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "Failed to list CronWorkflows in namespace %s", namespace)

	return cronWorkflowList.Items
}

// WaitForCronWorkflowActive waits until the specified Argo CronWorkflow reaches the 'Active' phase within the given timeout.
// It uses the provided KubectlOptions, workflow name, and namespace for the check.
// Fails the test if the CronWorkflow does not become active within the timeout.
func WaitForCronWorkflowActive(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForCronWorkflowPhase(t, options, name, namespace, workflowv1alpha1.ActivePhase, timeout)
}

// WaitForCronWorkflowStopped waits until the specified Argo CronWorkflow reaches the "Stopped" phase within the given timeout.
// It uses the provided testing context, kubectl options, workflow name, and namespace.
// If the workflow does not reach the "Stopped" phase within the timeout, the test will fail.
func WaitForCronWorkflowStopped(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	WaitForCronWorkflowPhase(t, options, name, namespace, workflowv1alpha1.StoppedPhase, timeout)
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
func WaitForCronWorkflowPhase(t *testing.T, options *k8s.KubectlOptions, name, namespace string, desiredPhase workflowv1alpha1.CronWorkflowPhase, timeout time.Duration) {
	t.Helper()

	client, err := NewArgoWorkflowsClient(t, options)
	require.NoError(t, err, "Failed to create Argo Workflows clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		wf, err := client.ArgoprojV1alpha1().CronWorkflows(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		return wf.Status.Phase == desiredPhase, nil
	})

	if err != nil {
		t.Fatalf("Workflow %s/%s did not reach phase %q in time: %v", namespace, name, desiredPhase, err)
	}
}

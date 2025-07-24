package events

import (
	"context"
	"testing"
	"time"

	argoeventsv1alpha1 "github.com/argoproj/argo-events/pkg/apis/events/v1alpha1"
	argoclientset "github.com/argoproj/argo-events/pkg/client/clientset/versioned"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForSensorReady waits until the specified Argo Sensor resource in the given namespace becomes Ready.
// It polls the sensor's status conditions at regular intervals until the ConditionReady is True or the timeout is reached.
// If the sensor does not become Ready within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the Sensor resource.
//   - namespace: The namespace where the Sensor is located.
//   - timeout: The maximum duration to wait for the sensor to become Ready.
func WaitForSensorReady(t *testing.T, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()

	client, err := argoclientset.NewForConfig(options.RestConfig)
	require.NoError(t, err, "Failed to create Argo clientset")

	ctx := t.Context()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		sensor, err := client.ArgoprojV1alpha1().Sensors(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		for _, cond := range sensor.Status.Conditions {
			if cond.Type == argoeventsv1alpha1.ConditionReady && cond.Status == "True" {
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		t.Fatalf("Sensor %s/%s did not become Ready: %v", namespace, name, err)
	}
}

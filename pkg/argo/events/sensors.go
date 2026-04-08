package events

import (
	"context"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"

	argoeventsv1alpha1 "github.com/argoproj/argo-events/pkg/apis/events/v1alpha1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// ListSensors retrieves a list of Argo Events Sensor resources from the specified namespace.
// It uses the provided testing context and kubectl options to create an Argo Events client,
// then lists all Sensor resources in the given namespace. The function fails the test if
// any error occurs during client creation or sensor listing.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options used to configure the client.
//   - namespace: The Kubernetes namespace from which to list the sensors.
//
// Returns:
//   - A slice of argoeventsv1alpha1.Sensor objects representing the sensors found in the namespace.
// ListSensors lists matching resources.
func ListSensors(t testing.TestingT, options *k8s.KubectlOptions, namespace string) []argoeventsv1alpha1.Sensor {
	sensors, err := ListSensorsE(t, options, namespace)
	require.NoError(t, err, "Failed to list Sensors in namespace %s", namespace)
	return sensors
}

// ListSensorsE lists matching resources.
func ListSensorsE(t testing.TestingT, options *k8s.KubectlOptions, namespace string) ([]argoeventsv1alpha1.Sensor, error) {
	client, err := NewArgoEventsClient(t, options)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	sensorList, err := client.ArgoprojV1alpha1().Sensors(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return sensorList.Items, nil
}

// WaitForSensorReady waits until the specified Argo Sensor resource in the given namespace becomes Ready.
// It polls the sensor's status conditions at regular intervals until the ConditionReady is True or the timeout is reached.
// If the sensor does not become Ready within the timeout, the test fails.
// Parameters:
//   - t: The testing context.
//   - options: Kubectl options containing the Kubernetes REST config.
//   - name: The name of the Sensor resource.
//   - namespace: The namespace where the Sensor is located.
//   - timeout: The maximum duration to wait for the sensor to become Ready.
// WaitForSensorReady waits for the resource condition to be satisfied.
func WaitForSensorReady(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	err := WaitForSensorReadyE(t, options, name, namespace, timeout)
	require.NoError(t, err, "Sensor %s/%s did not become Ready", namespace, name)
}

// WaitForSensorReadyE waits for the resource condition to be satisfied.
func WaitForSensorReadyE(t testing.TestingT, options *k8s.KubectlOptions, name, namespace string, timeout time.Duration) error {
	client, err := NewArgoEventsClient(t, options)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		sensor, err := client.ArgoprojV1alpha1().Sensors(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		var (
			hasTriggers = false
			hasDeployed = false
			hasDeps     = false
		)

		for _, cond := range sensor.Status.Conditions {
			if cond.Type == argoeventsv1alpha1.SensorConditionTriggersProvided && cond.IsTrue() {
				hasTriggers = true
			}
			if cond.Type == argoeventsv1alpha1.SensorConditionDeployed && cond.IsTrue() {
				hasDeployed = true
			}
			if cond.Type == argoeventsv1alpha1.SensorConditionDepencencyProvided && cond.IsTrue() {
				hasDeps = true
			}
		}

		return hasTriggers && hasDeployed && hasDeps, nil
	})
}

package velero

import (
	"context"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// WaitForScheduleToExist waits until a Velero Schedule resource with the specified name and namespace exists
// and is in the "Enabled" phase, or until the given timeout is reached. It polls the Kubernetes API at regular
// intervals and fails the test if the schedule does not become enabled within the timeout period.
//
// Parameters:
//   - t: The testing context.
//   - options: The kubectl options containing the Kubernetes REST config.
//   - name: The name of the Velero Schedule resource to wait for.
//   - namespace: The namespace where the Velero Schedule resource is expected to exist.
//   - timeout: The maximum duration to wait for the schedule to become enabled.
//
// This function logs retries and fails the test with a fatal error if the schedule does not become enabled in time.
func WaitForScheduleToExist(t *testing.T, options k8s.KubectlOptions, name, namespace string, timeout time.Duration) {
	t.Helper()
	client, err := NewVeleroClient(options.RestConfig)
	require.NoError(t, err, "Unable to create Velero client")
	ctx := t.Context()

	key := ctrlclient.ObjectKey{Name: name, Namespace: namespace}

	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		var schedule velerov1.Schedule
		err := client.Get(ctx, key, &schedule)
		if err != nil {
			t.Logf("Retrying: Schedule %s/%s not found: %v", namespace, name, err)
			return false, nil
		}
		return schedule.Status.Phase == velerov1.SchedulePhaseEnabled, nil
	})

	if err != nil {
		t.Fatalf("Schedule %s/%s did not become enabled: %v", namespace, name, err)
	}
}

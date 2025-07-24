package certmanager

import (
	"context"
	"testing"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	fakecm "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestWaitForCertificateRequestReady(t *testing.T) {
	tests := []struct {
		name        string
		conditions  []cmv1.CertificateRequestCondition
		expectError bool
	}{
		{
			name: "ready true",
			conditions: []cmv1.CertificateRequestCondition{
				{Type: cmv1.CertificateRequestConditionReady, Status: cmmetav1.ConditionTrue},
			},
			expectError: false,
		},
		{
			name: "ready false",
			conditions: []cmv1.CertificateRequestCondition{
				{Type: cmv1.CertificateRequestConditionReady, Status: cmmetav1.ConditionFalse},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = cmv1.AddToScheme(scheme)

			client := fakecm.NewSimpleClientset(&cmv1.CertificateRequest{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cr",
					Namespace: "default",
				},
				Status: cmv1.CertificateRequestStatus{
					Conditions: tc.conditions,
				},
			})

			err := wait.PollImmediate(100*time.Millisecond, 1*time.Second, func() (bool, error) {
				cr, _ := client.CertmanagerV1().CertificateRequests("default").Get(context.TODO(), "test-cr", metav1.GetOptions{})
				for _, cond := range cr.Status.Conditions {
					if cond.Type == cmv1.CertificateRequestConditionReady && cond.Status == cmmetav1.ConditionTrue {
						return true, nil
					}
				}
				return false, nil
			})

			if tc.expectError && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

package certmanager

import (
	"testing"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/tj/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestWaitForCertificateRequestReady(t *testing.T) {
	tests := []struct {
		name        string
		conditions  []cmv1.CertificateRequestCondition
		expectError bool
	}{
		{
			name:        "No Conditions",
			conditions:  []cmv1.CertificateRequestCondition{},
			expectError: true,
		},
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

			NewTestClient(t, &cmv1.CertificateRequest{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cr",
					Namespace: "default",
				},
				Status: cmv1.CertificateRequestStatus{
					Conditions: tc.conditions,
				},
			})

			err := WaitForCertificateRequestReadyE(t, k8soptions, "test-cr", "default", 10*time.Second)

			if tc.expectError && err == nil {
				assert.Error(t, err)
				t.Fatalf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				assert.NoError(t, err)
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

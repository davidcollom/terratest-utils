package istio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	istiometa "istio.io/api/meta/v1alpha1"
)

// TestIstioConditionReady tests the istioConditionReady helper function
func TestIstioConditionReady(t *testing.T) {
	tests := []struct {
		name       string
		status     *istiometa.IstioStatus
		expectReady bool
	}{
		{
			name: "Ready condition true",
			status: &istiometa.IstioStatus{
				Conditions: []*istiometa.IstioCondition{
					{
						Type:   "Ready",
						Status: "true",
					},
				},
			},
			expectReady: true,
		},
		{
			name: "Ready condition false",
			status: &istiometa.IstioStatus{
				Conditions: []*istiometa.IstioCondition{
					{
						Type:   "Ready",
						Status: "false",
					},
				},
			},
			expectReady: false,
		},
		{
			name: "No ready condition",
			status: &istiometa.IstioStatus{
				Conditions: []*istiometa.IstioCondition{
					{
						Type:   "Other",
						Status: "true",
					},
				},
			},
			expectReady: false,
		},
		{
			name: "Empty conditions",
			status: &istiometa.IstioStatus{
				Conditions: []*istiometa.IstioCondition{},
			},
			expectReady: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := istioConditionReady(t, tt.status)
			assert.Equal(t, tt.expectReady, result)
		})
	}
}

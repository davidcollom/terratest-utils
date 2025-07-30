package k8s

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetCustomResourceDefinitionE(t *testing.T) {
	type fields struct {
		crdName string
		crds    []*apixv1.CustomResourceDefinition
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		expectedCRD *apixv1.CustomResourceDefinition
	}{
		{
			name: "success - CRD exists",
			fields: fields{
				crdName: "mycrd.example.com",
				crds: []*apixv1.CustomResourceDefinition{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "mycrd.example.com",
						},
					},
				},
			},
			wantErr: false,
			expectedCRD: &apixv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: "mycrd.example.com",
				},
			},
		},
		{
			name: "failure - CRD does not exist",
			fields: fields{
				crdName: "missingcrd.example.com",
				crds: []*apixv1.CustomResourceDefinition{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "othercrd.example.com",
						},
					},
				},
			},
			wantErr:     true,
			expectedCRD: nil,
		},
		{
			name: "failure - invalid rest config",
			fields: fields{
				crdName: "mycrd.example.com",
				crds:    nil,
			},
			wantErr:     true,
			expectedCRD: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cl := NewAPIXTestClient(t, []runtime.Object{})
			if tt.fields.crds != nil {
				// Create the CRDs in the fake client
				for _, crd := range tt.fields.crds {
					_, err := cl.ApiextensionsV1().CustomResourceDefinitions().Create(t.Context(), crd, metav1.CreateOptions{})
					require.NoError(t, err, "Failed to create test CRD %s", crd.Name)
				}
			}

			got, err := GetCustomResourceDefinitionE(t, k8soptions, tt.fields.crdName, metav1.GetOptions{})
			if tt.wantErr {
				require.NotNil(t, err, "Expected error but got none")
				assert.Error(t, err, "Expected error but got none")
				return
			}
			if tt.expectedCRD == nil && got != nil {
				assert.Equal(t, got.Name, tt.expectedCRD.Name)
			}
			assert.Equal(t, tt.expectedCRD, got)
		})
	}
}

// k8soptions a global k8s.KubectlOptions instance to be used within many tests..
var k8soptions = &k8s.KubectlOptions{}

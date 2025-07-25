package k8s

import (
	"context"
	"testing"

	terrak8s "github.com/gruntwork-io/terratest/modules/k8s"
	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// fakeRestConfig is a dummy rest.Config for testing.
var fakeRestConfig = &rest.Config{Host: "https://fake"}

func TestGetCustomResourceDefinitionE(t *testing.T) {
	type fields struct {
		restConfig *rest.Config
		crdName    string
		crds       []*apixv1.CustomResourceDefinition
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
				restConfig: fakeRestConfig,
				crdName:    "mycrd.example.com",
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
				restConfig: fakeRestConfig,
				crdName:    "missingcrd.example.com",
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
				restConfig: nil,
				crdName:    "mycrd.example.com",
				crds:       nil,
			},
			wantErr:     true,
			expectedCRD: nil,
		},
	}
	t.Skipf("Skipping test %s as it requires a real Kubernetes cluster", t.Name())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup fake clientset
			client := apixfake.NewSimpleClientset()
			for _, crd := range tt.fields.crds {
				_, err := client.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to create fake CRD: %v", err)
				}
			}

			options := &terrak8s.KubectlOptions{
				RestConfig: tt.fields.restConfig,
			}
			got, err := GetCustomResourceDefinitionE(t, options, tt.fields.crdName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCustomResourceDefinitionE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expectedCRD != nil && got != nil && got.Name != tt.expectedCRD.Name {
				t.Errorf("expected CRD name %v, got %v", tt.expectedCRD.Name, got.Name)
			}
			if tt.expectedCRD == nil && got != nil {
				t.Errorf("expected nil CRD, got %v", got)
			}
		})
	}
}

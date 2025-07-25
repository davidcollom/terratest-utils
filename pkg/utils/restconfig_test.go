package utils

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestGetRestConfigE_ReturnsRestConfigIfPresent(t *testing.T) {
	expectedConfig := &rest.Config{Host: "https://example.com"}
	options := &k8s.KubectlOptions{
		RestConfig: expectedConfig,
	}
	cfg, err := GetRestConfigE(t, options)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)
}

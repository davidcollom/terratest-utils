module github.com/davidcollom/terratest-utils/pkg/argo/cd

go 1.27.1

require (
	github.com/argoproj/argo-cd/v3 v3.6.0-rc2
	github.com/davidcollom/terratest-utils/pkg/utils v0.1.0
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/k8s/v2 v2.0.0
	github.com/stretchr/testify v1.12.1
)

require github.com/argoproj/argo-cd/gitops-engine/v3 v3.6.0-rc2 // indirect

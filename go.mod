module github.com/davidcollom/terratest-utils

go 1.25.1

// Needed this to avoid issues with sigs.k8s.io/gateway-api v1.1.0
replace sigs.k8s.io/gateway-api => sigs.k8s.io/gateway-api v1.0.0

# Terratest Utils

This package provides a set of helpers for [Terratest](https://terratest.gruntwork.io/) and general integration testing. It is designed to help validate that platforms and deployments have been successfully provisioned and configured.

## Features

- Utilities for testing Kubernetes resources, including cert-manager, external-secrets, ArgoCD, and Flux
- Functions to check readiness, status, and correctness of deployed resources
- Simplifies writing robust integration tests for cloud-native platforms
- Can be used with Terratest or standalone in Go test suites

## Usage

Import the relevant package(s) in your Terratest or Go integration tests:

```go
import (
    "github.com/davidcollom/terratest-utils/pkg/certmanager"
    "github.com/davidcollom/terratest-utils/pkg/externalsecrets"
    // ...other helpers
)
```

Use the provided functions to validate resources, e.g.:

```go
ready := certmanager.IsCertificateReady(cert)
if !ready {
    t.Fatalf("Certificate is not ready!")
}
```

## Structure

- `pkg/certmanager/` - Helpers for cert-manager resources
- `pkg/externalsecrets/` - Helpers for external-secrets resources
- `pkg/argo/` - Helpers for ArgoCD resources
- `pkg/flux/` - Helpers for Flux resources

## Purpose

These helpers are intended to:

- Accelerate writing integration tests for Kubernetes platforms
- Provide reusable checks for resource readiness and correctness
- Help ensure deployments are successful and meet expected criteria

## License

MIT

# Contributing to terratest-utils

Thank you for your interest in contributing! We welcome improvements, bug fixes, and new helpers for integration testing.

## How to Contribute

1. **Fork the repository** and create your branch from `main`.
2. **Make your changes** (code, documentation, tests, etc.).
3. **Run tests and linting** to ensure your changes do not break existing functionality.

    The repo is a monorepo of independent Go modules wired together by `go.work`. The repo root is an empty marker module, so `go build ./...` at the root only matches packages inside the root — it does not recurse into the 12 nested modules. Build and test per module, or loop over them:

    ```sh
    # Build a single module
    cd pkg/certmanager
    go build ./...
    go test ./...

    # Build every module
    for mod in pkg/utils pkg/k8s pkg/certmanager pkg/externalsecrets pkg/flux pkg/istio pkg/linkerd pkg/velero pkg/argo/cd pkg/argo/events pkg/argo/rollouts pkg/argo/workflows; do
      (cd "$mod" && go build ./... && go vet ./...)
    done
    ```

    When you add a new function, work inside the specific module the change belongs to and run `go build ./...` and `go vet ./...` from that module's directory before committing.

4. **Commit your changes** with clear and descriptive messages.
5. **Push to your fork** and open a Pull Request (PR) against the `main` branch.
6. **Describe your changes** in the PR, including motivation and any relevant context.

## Release process

Releases are automated via `.github/workflows/release.yaml`. To cut a release, run the `Release` workflow from the Actions tab (`workflow_dispatch`) and pick a release type (`auto`, `major`, `minor`, or `patch`). The workflow computes the next version with `svu`, runs `./scripts/bump-version.sh` to update cross-module `require` directives and drop the local `replace` directives, commits the bump, creates the version and alias tags, and publishes a GitHub Release. See the README for full details.

## Guidelines

- Follow existing code style and conventions.
- Update documentation as needed for your changes.
- Ensure all tests pass before submitting your PR.
- If adding new features, include corresponding tests.

## Code of Conduct

Be respectful and constructive in all interactions. For questions or help, open an issue or start a discussion.

Thank you for helping make this project better!

# AGENTS.md

## Project Purpose

`secret-contract-operator` is a Kubernetes operator written in Go using Kubebuilder and controller-runtime.

The operator manages `SecretContract` resources. A `SecretContract` validates that a Kubernetes `Secret` satisfies an application-level contract, including:

- Required keys.
- Non-empty values.
- Length and pattern rules.
- Optional `ExternalSecret` readiness checks.
- Optional workload reference validation.
- Optional workload restart on `Secret` updates.

The operator is complementary to secret management tools such as External Secrets Operator and Secrets Store CSI Driver. It must not become a new secret-manager syncer.

## Security Rules

- Never log `Secret` values.
- Never write `Secret` values to status, events, metrics, annotations, errors, tests, or docs.
- Status can contain key names and violation reasons, but not values.
- Prefer `Secret` `resourceVersion` for restart annotations, not a hash of `Secret` data.
- Default behavior must be validation-only.
- Workload mutation must be opt-in.

## Development Commands

Use these commands during development:

```sh
go test ./...
make test
make generate
make manifests
make install
make run
make docker-build IMG=...
make deploy IMG=...
```

Run `make generate` and `make manifests` whenever API types, CRD schema, RBAC markers, or generated controller metadata change.

## Codex Git Workflow

- After completing each user prompt, run `git add`, `git commit`, and `git push` for the completed changes.
- Keep commits focused on the work from the current prompt.
- Do not include unrelated user changes in commits.

## Code Style

- Use `gofmt` for all Go code.
- Keep controller logic small and testable.
- Put pure validation logic in `internal/contract` or a similar package so it can be unit-tested without Kubernetes.
- Use `metav1.Condition` and `meta.SetStatusCondition` for status.
- Use controller-runtime client and events.
- Handle `NotFound` gracefully.
- Use `context.Context` everywhere.

## Testing Expectations

- Add unit tests for pure validation functions.
- Add envtest tests for reconciler behavior.
- Test missing `Secret`, missing key, empty key, `minLength` failure, pattern failure, success, and status condition updates.
- Test that `Secret` values do not appear in status or events.

## RBAC Expectations

- Use least privilege.
- Generate RBAC via Kubebuilder markers.
- Read `Secrets` but do not create, update, or delete `Secrets`.
- Update only `SecretContract` status.
- Patch workloads only when opt-in mutation or restart policy is enabled.

## Documentation Expectations

- Update `README.md` and examples for every user-facing behavior.
- Include examples under `examples/`.
- Include limitations and the security model.

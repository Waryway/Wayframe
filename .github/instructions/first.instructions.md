---
applyTo: '**'
description: 'description'
---
Provide project context and coding guidelines that AI should follow when generating code, answering questions, or reviewing changes.

# Wayframe GoLand AI Agent Instructions

## Purpose
Provide concise, enforceable rules for AI-assisted edits inside the `Wayframe` repository so changes remain buildable, idiomatic, and CI-safe.

## Global constraints
- Use Go SDK `1.25`. Configure the IDE to use Go 1.25 for builds and tests.
- Bazel must use Bzlmod (Bazel 8). **Do not** add or modify `WORKSPACE` files. Use `MODULE.bazel` and `bazel_dep()` for dependencies.
- Use `rules_go` and Go extensions provided via Bzlmod. Do not use `io_bazel_rules_go`.
- Prefer `pkg/logger` (slog-style) over `fmt.Println` for logging.

## Coding & package rules
- Respect `internal/` boundaries. Do not make internal packages importable from external modules.
- Follow idiomatic Go conventions and modern features (context usage, error wrapping).
- When changing code that affects interfaces in `internal/web`, ensure implementations for `stdlib`, `fiber`, and `gorilla` remain compatible with `internal/web.Server`.

## Configuration conventions
- Support JSON, YAML, and key-value files.
- Use struct tags:
    - ``config:"key"``
    - ``env:"ENV_VAR"``
    - ``default:"value"``
    - ``file:"path"``
- Configuration precedence: Environment variables → File values → Default values.

## Frontend / JS rules
- Use `pnpm` for package management.
- Build React apps via Bazel rules; frontend outputs are embedded into Go servers.
- Do not modify generated Bazel targets for React without coordinating with frontend owners.

## Build & test commands
- Bazel:
    - `bazel build //...`
    - `bazel test //...`
- Go (local dev):
    - `go build ./...`
    - `go test ./...`
- Always run `verify-ci-setup.sh` before pushing CI-impacting changes.

## GoLand setup recommendations
- Set Project SDK to Go `1.25`.
- Install/configure Bazel plugin and create Bazel Run/Debug configurations for common targets.
- Add a File Watcher or External Tool to run `pnpm install` and `pnpm build` for frontend edits if needed.
- Configure GOPATH/module support as default (repo uses `go.mod`).

## Agent behavior rules
- Never introduce or modify `WORKSPACE` files. If dependency changes are needed, update `MODULE.bazel` and request review.
- Avoid editing files under `bazel-*`, `bazel-bin`, `vendor/` or other generated output directories.
- Prefer small, test-covered changes. Add or update tests when behavior changes.
- If uncertain about a cross-cutting change (build rules, MODULE.bazel, major dependency upgrades), add a review request and a short PR description explaining the rationale.
- When making changes that affect frontend embedding or server build, ensure both Bazel and `go build` succeed locally.

## Quick checklist for AI agents before producing a patch
1. Ensure Go version is 1.25 in local run configuration.
2. Confirm no `WORKSPACE` edits; use `MODULE.bazel`.
3. Run `go test ./...` and `bazel test //...` locally (or report failures).
4. Do not change generated files or vendor without explicit reason and review.
5. Add/modify tests to cover behavioral changes.

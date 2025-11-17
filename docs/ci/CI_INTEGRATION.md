# CI/CD Pipeline Integration Guide

This document explains how to use the newly added GitHub Actions CI/CD pipeline with Bazel caching.

## Quick Start

### Triggering CI
Simply push to `main` or `develop`, or open a PR:

```bash
git push origin feature-branch
# or create a PR on GitHub
```

The pipeline will automatically:
1. ✅ Build all targets
2. ✅ Run all tests
3. ✅ Check code quality
4. ✅ Scan for security issues

### Viewing Results
1. Go to GitHub **Actions** tab
2. Click **"CI Pipeline"** workflow
3. View job details or logs

## Pipeline Overview

### Jobs (run in parallel)

| Job | Purpose | Duration |
|-----|---------|----------|
| **Build & Test** | Compiles code, runs unit tests | 30-60s (with cache) |
| **Coverage** | Generates coverage reports | 45-90s |
| **Lint** | Checks code quality (golangci-lint) | 15-30s |
| **Security** | Scans for vulnerabilities (Trivy) | 20-40s |

### Caching Strategy

**Cache Key:** `bazel-cache-{OS}-{hash(MODULE.bazel, MODULE.bazel.lock, go.mod)}`

**Saves Time By:**
- Storing compiled binaries between runs
- Skipping rebuilds for unchanged dependencies
- Reusing Go module downloads

**Typical Times:**
- First build: 2-5 minutes
- With cache: 30-60 seconds
- No code changes: ~20 seconds

## Local Development

### Use the Same Cache Locally

1. Copy the template:
```bash
cp .bazelrc.local.template .bazelrc.local
```

2. Build with cache:
```bash
bazel build //...
bazel test //...
```

Your local builds will now use the same caching as CI!

### Clearing Cache

If needed, clear the local cache:
```bash
rm -rf ~/.bazel-cache
rm -rf ~/.bazel-repo-cache
```

## Documentation

Detailed documentation available in `.github/`:

- **`CI_QUICK_START.md`** - Developer quick reference
- **`CICD_DOCUMENTATION.md`** - Technical deep dive
- **`CI_SETUP_SUMMARY.md`** - Setup details

And in root:
- **`CI_PIPELINE_SETUP.md`** - Complete setup guide

## Status Badge

Show CI status in your README:

```markdown
[![CI Pipeline](https://github.com/Waryway/Wayframe/actions/workflows/ci.yml/badge.svg)](https://github.com/Waryway/Wayframe/actions/workflows/ci.yml)
```

## Troubleshooting

### My build is slow even with cache?

Check the logs:
1. Look at "Restore Bazel cache" step
2. Verify it says "Cache hit" (not "Cache miss")
3. If cache miss, check if dependencies changed

### Cache seems stale?

Cache is keyed on:
- `MODULE.bazel` (Bazel deps)
- `MODULE.bazel.lock` (dependency resolution)
- `go.mod` (Go deps)

If any of these change, cache is invalidated (which is correct!).

### Tests pass locally but fail in CI?

Try running with CI configuration:
```bash
bazel test //... --disk_cache=~/.bazel-cache
```

If that passes, it's likely a local setup issue.

## Configuration Files

### `.bazelrc`
Global Bazel settings including:
- Disk cache paths
- Build resource limits
- Test output format

### `.bazelrc.local` (optional)
Local overrides for development. Copy from `.bazelrc.local.template`.

### `.github/workflows/ci.yml`
GitHub Actions workflow definition.

## Performance Tips

1. **Keep dependencies minimal** - Only add what you need to `go.mod`
2. **Commit lock files** - Always commit `MODULE.bazel.lock` and `go.sum`
3. **Test locally first** - Use `bazel test //...` before pushing
4. **Use specific targets** - `bazel test //pkg/...` is faster than `//...` if possible

## What's Cached?

✅ Compiled Go binaries
✅ Test artifacts
✅ Downloaded Go modules
✅ Gazelle outputs
✅ Bazel analysis results

❌ Source code (always fresh)
❌ Secrets (never cached)
❌ External network requests (not usually cached)

## Next Steps

1. Review documentation in `.github/` folder
2. Push a change to trigger a build
3. Monitor the pipeline in Actions tab
4. Set up `.bazelrc.local` locally
5. Enjoy faster builds! 🚀

## Questions?

See the detailed documentation:
- Technical details: `.github/CICD_DOCUMENTATION.md`
- Developer guide: `.github/CI_QUICK_START.md`
- Setup guide: `CI_PIPELINE_SETUP.md`


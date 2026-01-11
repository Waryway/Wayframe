I'll update the plan to add Node.js setup to your CI workflow and show exactly where and what to add next; I won't make any repo edits until you confirm.

Checklist
- [ ] Add `actions/setup-node` step to `/.github/workflows/ci.yml` (recommended Node LTS)
- [ ] Add `actions/cache` step to cache `~/.npm` keyed on `examples/react/frontend/package-lock.json`
- [ ] (Optional) Add a fast-fail step: `bazel build //examples/react/frontend:dist` before `bazel build //...`
- [ ] Run CI and verify `bazel build //...` no longer fails with `could not embed dist/*`

## Plan: Install Node in CI (Add steps to `/.github/workflows/ci.yml`)

TL;DR — Insert Node setup and npm cache steps into the `build-and-test` job of `/.github/workflows/ci.yml`, immediately after the `Checkout code` step, and optionally add a dedicated frontend build step to surface frontend failures early. This makes `npm` available for the `react_app` genrule so the frontend `dist` files are produced and `//go:embed dist/*` succeeds.

### Steps (concrete)
1. Edit `/.github/workflows/ci.yml` and locate the `build-and-test` job. Find the `steps:` list and the `- name: Checkout code` step.
2. Immediately after the `Checkout code` step, insert a `Set up Node.js` step that uses `actions/setup-node@v4` with `node-version: '18.x'` (or your preferred LTS).
   - Purpose: ensures `node` and `npm` exist in the runner so the frontend genrule can run `npm install` and `npm run build`.
3. Add an `npm cache` step using `actions/cache@v4` right after Node setup:
   - Cache `path: ~/.npm`
   - Use `key: npm-cache-${{ runner.os }}-${{ hashFiles('examples/react/frontend/package-lock.json') }}`
   - Use `restore-keys: npm-cache-${{ runner.os }}-`
   - Purpose: speed up CI and reduce flakiness by caching npm packages.
4. (Optional but recommended) Add an explicit `Build frontend` step before the global `bazel build //...`:
   - Run `bazel build //examples/react/frontend:dist`
   - Purpose: fail fast with clear logs if the frontend build (npm) fails or if network/npm issues exist.
5. Keep the existing `Build all targets` step (the `bazel build //...`) but consider temporarily appending `--verbose_failures --sandbox_debug` while we validate the fix, then remove verbosity after success.
6. Commit these CI changes and trigger the pipeline. If the pipeline still fails:
   - Inspect the `Build frontend` step logs. If npm errors (e.g., network or permission issues), consider adding `npm ci` or pinning Node version.
   - If the genrule runs but `go:embed` still complains, we will revisit `examples/react/BUILD.bazel` wiring and the `react_app` macro to make outputs explicit.

### Exactly where and what to add (description)
- Insert the following step sequence in `build-and-test.steps` after `- name: Checkout code`:
  1. Step named `Set up Node.js` that uses `actions/setup-node@v4` with input `node-version: '18.x'`.
  2. Step named `Cache npm` that uses `actions/cache@v4` with:
     - `path` = `~/.npm`
     - `key` = `npm-cache-${{ runner.os }}-${{ hashFiles('examples/react/frontend/package-lock.json') }}`
     - `restore-keys` = `npm-cache-${{ runner.os }}-`
  3. (Optional) Step named `Build frontend` that runs `bazel build //examples/react/frontend:dist`.
- No other workflow changes are required to begin; `bazel build //...` should then trigger the genrule successfully because `npm`/`node` are available.

### Further considerations
1. Node version: pick an LTS (18.x or 20.x). Using `18.x` is widely compatible. If your frontend needs a particular Node version, pin that instead.
2. Network/npm flakiness: if CI still fails due to registry/network issues, add a step to run `npm ci` in `examples/react/frontend` and/or use a private npm cache or registry mirror.
3. Bazel-native alternative: longer-term, consider integrating `rules_nodejs` to make frontend builds hermetic under Bazel and reduce reliance on runner-installed Node. For now, installing Node in the runner is the fastest path for front-end dev ergonomics.

If you approve, I will prepare the exact YAML edits to `/.github/workflows/ci.yml` and show the patch (or apply it if you want me to edit directly). Which do you want next: (A) show the exact YAML snippet to paste, or (B) I should apply the edit to the file now?


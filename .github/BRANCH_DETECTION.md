# CI Pipeline - Branch Detection Update

## Changes Made

### 1. **Extended Push Trigger**
The pipeline now runs on **ALL branches** when code is pushed, not just `main` and `develop`.

```yaml
on:
  push:
    branches:
      - '**'  # Now runs on all branches
  pull_request:
    branches:
      - main
      - develop
```

### 2. **Branch Detection Job**
A new `detect-branch` job identifies branch type and provides outputs for other jobs:

```yaml
detect-branch:
  name: Detect Branch Type
  outputs:
    branch-type: main | develop | feature
    branch-name: actual branch name
    is-main: true | false
```

**Detection Logic:**
- `main` branch → `branch-type=main`, `is-main=true`
- `develop` branch → `branch-type=develop`, `is-main=false`
- Any other branch → `branch-type=feature`, `is-main=false`

### 3. **Branch Information Display**
All jobs now depend on `detect-branch` and display branch information:

```yaml
- name: Display branch information
  run: |
    echo "🔍 Branch Type: ${{ needs.detect-branch.outputs.branch-type }}"
    echo "📝 Branch Name: ${{ needs.detect-branch.outputs.branch-name }}"
    echo "✓ Is Main: ${{ needs.detect-branch.outputs.is-main }}"
```

## Job Dependencies

All jobs now depend on the `detect-branch` job:
- **Build and Test** - `needs: detect-branch`
- **Coverage** - `needs: detect-branch`
- **Lint** - `needs: detect-branch`
- **Security Scan** - `needs: detect-branch`

This ensures branch information is available to all jobs.

## Usage Examples

### Access Branch Information in Workflows

In any job that needs branch information:

```yaml
steps:
  - name: Do something based on branch
    run: |
      if [ "${{ needs.detect-branch.outputs.is-main }}" == "true" ]; then
        echo "Running on main branch - deploying to production"
      else
        echo "Running on ${{ needs.detect-branch.outputs.branch-type }} branch"
      fi
```

### Conditional Steps Based on Branch

```yaml
- name: Deploy (only on main)
  if: needs.detect-branch.outputs.is-main == 'true'
  run: ./deploy.sh

- name: Run extra tests (only on feature branches)
  if: needs.detect-branch.outputs.branch-type == 'feature'
  run: bazel test //... --test_filter="*integration*"
```

## Log Output

When you view the GitHub Actions logs, you'll see branch detection output:

```
🔍 Branch Detection Results
Branch Type: feature
Branch Name: feature/awesome-feature
Is Main: false
```

This appears in every job that depends on branch detection.

## Pipeline Flow

```
1. Push/PR Triggered
        ↓
2. detect-branch job runs (determines branch type)
        ↓
3. All 4 jobs run in parallel with branch info:
   - build-and-test
   - coverage
   - lint
   - security-scan
```

## Examples

### Example 1: Push to feature branch
```
Event: git push origin feature/new-feature
Branch Detection:
  Type: feature
  Name: feature/new-feature
  Is Main: false

Result: Pipeline runs, shows all logs tagged as "feature branch"
```

### Example 2: Push to main
```
Event: git push origin main
Branch Detection:
  Type: main
  Name: main
  Is Main: true

Result: Pipeline runs, shows all logs tagged as "main branch"
```

### Example 3: PR to main
```
Event: Pull request to main
Branch Detection:
  Type: feature (or whatever your branch is)
  Name: feature/awesome-feature
  Is Main: false (PR target doesn't change source branch type)

Result: Pipeline runs with branch info from your feature branch
```

## Next Steps

You can now use the branch information to:

1. **Add deployment jobs** - Only deploy to production from main
2. **Run different test suites** - More comprehensive tests on main
3. **Generate different artifacts** - Different release channels per branch
4. **Add branch-specific notifications** - Alert team leads on main branch failures

Example deployment job:

```yaml
deploy:
  name: Deploy to Production
  needs: [build-and-test, detect-branch]
  if: needs.detect-branch.outputs.is-main == 'true'
  runs-on: ubuntu-latest
  steps:
    # Deployment steps here
```

## Troubleshooting

### Branch detection shows "unknown"?
- Check that the branch name is being passed correctly
- The script uses `${{ github.ref }}` which is set by GitHub Actions
- Should be `refs/heads/branch-name`

### Jobs not waiting for branch detection?
- Make sure all jobs have `needs: detect-branch`
- If a job doesn't depend on it, it won't wait
- Jobs can still run in parallel with `needs`

### Using branch info in conditionals?
- Use: `if: needs.detect-branch.outputs.is-main == 'true'`
- Not: `if: ${{ needs.detect-branch.outputs.is-main }}`
- The single braces are for the if condition directly

## Summary

✅ Pipeline now runs on **all branches** with `git push`
✅ **Branch type detection** identifies main/develop/feature
✅ **All jobs** receive branch information
✅ **Logs show** branch context for debugging
✅ **Future-proof** for adding branch-specific logic


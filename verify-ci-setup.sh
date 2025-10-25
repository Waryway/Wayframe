#!/bin/bash
# CI Pipeline Verification Checklist
# Run this script to verify the CI setup is working correctly

echo "🔍 CI Pipeline Verification Checklist"
echo "===================================="
echo ""

# Check 1: Workflow file exists
echo "✓ Checking workflow file..."
if [ -f ".github/workflows/ci.yml" ]; then
    echo "  ✅ .github/workflows/ci.yml exists"
else
    echo "  ❌ .github/workflows/ci.yml NOT FOUND"
fi

# Check 2: .bazelrc has cache settings
echo ""
echo "✓ Checking .bazelrc configuration..."
if grep -q "disk_cache" .bazelrc; then
    echo "  ✅ Disk cache configured"
else
    echo "  ❌ Disk cache NOT configured"
fi

if grep -q "repository_cache" .bazelrc; then
    echo "  ✅ Repository cache configured"
else
    echo "  ❌ Repository cache NOT configured"
fi

# Check 3: Tools directory
echo ""
echo "✓ Checking tools directory..."
if [ -f "tools/workspace_status.sh" ]; then
    echo "  ✅ tools/workspace_status.sh exists"
else
    echo "  ❌ tools/workspace_status.sh NOT FOUND"
fi

if [ -f "tools/BUILD.bazel" ]; then
    echo "  ✅ tools/BUILD.bazel exists"
else
    echo "  ❌ tools/BUILD.bazel NOT FOUND"
fi

# Check 4: Documentation files
echo ""
echo "✓ Checking documentation..."
docs=(
    ".github/CI_QUICK_START.md"
    ".github/CICD_DOCUMENTATION.md"
    "CI_SETUP_SUMMARY.md"
    "CI_PIPELINE_SETUP.md"
    "CI_INTEGRATION.md"
)

for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        echo "  ✅ $doc exists"
    else
        echo "  ❌ $doc NOT FOUND"
    fi
done

# Check 5: Template files
echo ""
echo "✓ Checking template files..."
if [ -f ".bazelrc.local.template" ]; then
    echo "  ✅ .bazelrc.local.template exists"
else
    echo "  ❌ .bazelrc.local.template NOT FOUND"
fi

# Check 6: Git status
echo ""
echo "✓ Checking git status..."
if [ -d ".git" ]; then
    echo "  ✅ Git repository exists"

    # Show uncommitted changes
    echo ""
    echo "  Uncommitted CI changes:"
    git status --short | grep -E "(\.github|CI_|tools/|\.bazelrc)" || echo "    (none or all committed)"
else
    echo "  ❌ Git repository NOT FOUND"
fi

# Check 7: Go module
echo ""
echo "✓ Checking Go configuration..."
if [ -f "go.mod" ]; then
    echo "  ✅ go.mod exists"
    go_version=$(grep "^go " go.mod | awk '{print $2}')
    echo "  📝 Go version: $go_version"
else
    echo "  ❌ go.mod NOT FOUND"
fi

# Check 8: MODULE.bazel
echo ""
echo "✓ Checking Bazel configuration..."
if [ -f "MODULE.bazel" ]; then
    echo "  ✅ MODULE.bazel exists"
    if grep -q "enable_bzlmod" .bazelrc; then
        echo "  ✅ Bzlmod enabled"
    fi
else
    echo "  ❌ MODULE.bazel NOT FOUND"
fi

# Check 9: Build test
echo ""
echo "✓ Testing Bazel build..."
if command -v bazel &> /dev/null; then
    echo "  ✅ Bazel installed"

    # Try a quick query
    if bazel query "//..." > /dev/null 2>&1; then
        echo "  ✅ Bazel query successful"
    else
        echo "  ⚠️  Bazel query failed (may need setup)"
    fi
else
    echo "  ⚠️  Bazel not in PATH (install Bazelisk)"
fi

# Check 10: Summary
echo ""
echo "===================================="
echo "✅ CI Pipeline Setup Verification Complete!"
echo ""
echo "📝 Next steps:"
echo "  1. Commit all CI-related files"
echo "  2. Push to main or open a PR"
echo "  3. Monitor GitHub Actions tab"
echo "  4. Verify cache hit on 2nd run"
echo ""
echo "📚 Documentation:"
echo "  - .github/CI_QUICK_START.md (quick start)"
echo "  - CI_INTEGRATION.md (integration guide)"
echo "  - CI_SETUP_SUMMARY.md (detailed setup)"
echo ""


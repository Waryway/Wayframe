"""Rules for building React/npm frontend applications with Bazel."""

def react_app(
        name,
        srcs,
        package_json = "package.json",
        build_script = "build",
        output_dir = "../dist",
        outs = None,
        visibility = None):
    """
    Builds a React application using npm.

    This rule runs `npm install` and `npm run <build_script>` to build a frontend
    application, then makes the output files available to Bazel.

    The rule automatically tracks source file changes and rebuilds when needed.

    Args:
        name: Name of the build target
        srcs: Source files that trigger rebuilds (TS, TSX, CSS, config files, etc.)
        package_json: Path to package.json file (default: "package.json")
        build_script: npm script to run from package.json (default: "build")
        output_dir: Directory where build outputs files, relative to package.json (default: "../dist")
        outs: List of output files to track (default: ["dist/index.html", "dist/bundle.js"])
        visibility: Visibility of the generated targets

    Example:
        react_app(
            name = "frontend",
            srcs = glob([
                "src/**/*.ts",
                "src/**/*.tsx",
                "src/**/*.css",
                "public/**/*",
            ]) + [
                "package-lock.json",
                "tsconfig.json",
                "webpack.config.js",
            ],
            package_json = "package.json",
            build_script = "build",
            output_dir = "../dist",
        )

    This creates:
        - :<name> - The build target
        - :<name>_dist - Filegroup with all outputs
    """

    if outs == None:
        outs = [
            "dist/index.html",
            "dist/bundle.js",
        ]

    # Get the package directory
    package_dir = native.package_name()
    if package_dir:
        frontend_dir = package_dir
    else:
        frontend_dir = "."

    # Create the genrule
    native.genrule(
        name = name,
        srcs = srcs + [package_json],
        outs = outs,
        cmd_bash = """
            set -euo pipefail

            # Create output directories for all outputs
            for out in $(OUTS); do
                mkdir -p $$(dirname $$out)
            done

            # Change to the frontend directory
            cd {frontend_dir}

            # Install dependencies
            echo "Installing npm packages..."
            npm install --silent 2>&1 | grep -v "npm WARN" || true

            # Build the application
            echo "Building frontend with npm run {build_script}..."
            npm run {build_script} --silent

            # Copy outputs to Bazel's output tree
            echo "Copying outputs..."
            cd {output_dir}

            # Copy index.html
            cp index.html $(location dist/index.html)

            # Copy the first .js file found (webpack hash changes)
            for js in *.js; do
                if [ -f "$$js" ]; then
                    cp "$$js" $(location dist/bundle.js)
                    break
                fi
            done
        """.format(
            frontend_dir = frontend_dir,
            build_script = build_script,
            output_dir = output_dir,
        ),
        local = 1,
        message = "Building React frontend",
        visibility = visibility,
    )

    # Create a filegroup for easy dependency reference
    native.filegroup(
        name = name + "_dist",
        srcs = [":%s" % name],
        visibility = visibility,
    )

def npm_build(
        name,
        srcs,
        package_json,
        script = "build",
        output_dir = "../dist",
        **kwargs):
    """
    Generic npm build rule.

    Lower-level rule for running any npm script. Use react_app for React applications.

    Args:
        name: Name of the build target
        srcs: Source files to track
        package_json: Path to package.json
        script: npm script to run (default: "build")
        output_dir: Output directory relative to package.json
        **kwargs: Additional arguments passed to genrule
    """
    react_app(
        name = name,
        srcs = srcs,
        package_json = package_json,
        build_script = script,
        output_dir = output_dir,
        **kwargs
    )


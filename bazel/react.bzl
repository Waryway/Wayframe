# bazel/react.bzl
def react_app(
        name,
        srcs,
        package_json = "package.json",
        build_script = "build",
        output_dir = "dist",
        outs = None,
        visibility = None):
    if outs == None:
        outs = [
            "dist/index.html",
            "dist/bundle.js",
        ]

    package_dir = native.package_name()
    if package_dir:
        frontend_dir = package_dir
    else:
        frontend_dir = "."

    native.genrule(
        name = name,
        srcs = srcs + [package_json],
        outs = outs,
        cmd_bash = """
            set -euo pipefail
            set -x

            # Grab the Bazel-expanded list of output paths into a shell variable.
            OUTS_STR="$(OUTS)"
            IFS=' ' read -r OUT_HTML OUT_JS <<< "$$OUTS_STR" || true

            echo "OUTS=$$OUTS_STR"
            echo "OUT_HTML=$$OUT_HTML"
            echo "OUT_JS=$$OUT_JS"

            # Use workspace-relative paths for the frontend build dir so we don't
            # depend on the genrule's current working directory.
            BUILD_DIR="{frontend_dir}/{output_dir}"

            # Run npm from the frontend package directory.
            (cd "{frontend_dir}" && \
              npm install --silent 2>&1 | grep -v "npm WARN" || true && \
              npm run {build_script} --silent)

            if [ ! -d "$$BUILD_DIR" ]; then
                echo "Build output directory '$$BUILD_DIR' not found" >&2
                exit 1
            fi

            # Ensure destination directories exist (Bazel expects genrule to create declared files)
            mkdir -p "$$(dirname "$$OUT_HTML")"

            if [ -n "$$OUT_JS" ]; then
                mkdir -p "$$(dirname "$$OUT_JS")"
            fi

            # Copy index.html
            if [ -f "$$BUILD_DIR/index.html" ]; then
                cp "$$BUILD_DIR/index.html" "$$OUT_HTML"
            else
                echo "index.html not found in $$BUILD_DIR" >&2
                exit 1
            fi

            # Find a JS bundle in the build dir and copy it if an OUT_JS was declared
            if [ -n "$$OUT_JS" ]; then
                JS_FILE=$$(ls "$$BUILD_DIR"/*.js 2>/dev/null | head -n 1 || true)
                if [ -z "$$JS_FILE" ]; then
                    echo "No JS bundle found in $$BUILD_DIR" >&2
                    exit 1
                fi
                cp "$$JS_FILE" "$$OUT_JS"
            fi
        """.format(
            frontend_dir = frontend_dir,
            build_script = build_script,
            output_dir = output_dir,
        ),
        local = 1,
        message = "Building React frontend",
        visibility = visibility,
    )

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
        output_dir = "dist",
        **kwargs):
    react_app(
        name = name,
        srcs = srcs,
        package_json = package_json,
        build_script = script,
        output_dir = output_dir,
        **kwargs
    )

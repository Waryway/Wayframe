package react

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BazelBuildConfig holds configuration for building React apps with Bazel.
type BazelBuildConfig struct {
	// WorkspaceDir is the root of the Bazel workspace
	WorkspaceDir string

	// Target is the Bazel target to build (e.g., "//apps/my-app:build")
	Target string

	// OutputDir is where to copy the build output (optional)
	OutputDir string

	// Config is the Bazel config to use (e.g., "opt" for optimized builds)
	Config string

	// AdditionalFlags are extra flags to pass to Bazel
	AdditionalFlags []string
}

// BuildWithBazel builds a React application using Bazel.
func BuildWithBazel(cfg BazelBuildConfig) error {
	if cfg.Target == "" {
		return fmt.Errorf("Target is required")
	}

	if cfg.WorkspaceDir == "" {
		cfg.WorkspaceDir = "."
	}

	fmt.Printf("Building React app with Bazel: %s\n", cfg.Target)

	// Build the Bazel command
	args := []string{"build"}

	if cfg.Config != "" {
		args = append(args, fmt.Sprintf("--config=%s", cfg.Config))
	}

	args = append(args, cfg.AdditionalFlags...)
	args = append(args, cfg.Target)

	// Execute Bazel build
	cmd := exec.Command("bazel", args...)
	cmd.Dir = cfg.WorkspaceDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build with Bazel: %w", err)
	}

	fmt.Println("✓ Bazel build completed successfully!")

	// Copy output if requested
	if cfg.OutputDir != "" {
		if err := copyBazelOutput(cfg); err != nil {
			return fmt.Errorf("failed to copy build output: %w", err)
		}
		fmt.Printf("✓ Build output copied to: %s\n", cfg.OutputDir)
	}

	return nil
}

// copyBazelOutput copies the Bazel build output to the specified directory.
func copyBazelOutput(cfg BazelBuildConfig) error {
	// Determine the Bazel output path
	// For a target like //apps/my-app:build, the output is typically in bazel-bin/apps/my-app/build
	targetParts := strings.Split(cfg.Target, ":")
	if len(targetParts) != 2 {
		return fmt.Errorf("invalid target format: %s", cfg.Target)
	}

	packagePath := strings.TrimPrefix(targetParts[0], "//")
	targetName := targetParts[1]

	bazelBinPath := filepath.Join(cfg.WorkspaceDir, "bazel-bin", packagePath, targetName)

	// Check if the output exists
	if _, err := os.Stat(bazelBinPath); os.IsNotExist(err) {
		return fmt.Errorf("build output not found at: %s", bazelBinPath)
	}

	// Create output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy the output
	// Use cp -r on Unix-like systems, xcopy on Windows
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.Command("xcopy", bazelBinPath, cfg.OutputDir, "/E", "/I", "/Y")
	} else {
		cmd = exec.Command("cp", "-r", bazelBinPath+"/.", cfg.OutputDir)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// isWindows checks if the current OS is Windows.
func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}

// GenerateReactBuildFile generates a BUILD.bazel file for a React application.
func GenerateReactBuildFile(appDir, appName string, typescript bool) error {
	buildContent := generateReactBuildContent(appName, typescript)

	buildPath := filepath.Join(appDir, "BUILD.bazel")

	if err := os.WriteFile(buildPath, []byte(buildContent), 0644); err != nil {
		return fmt.Errorf("failed to write BUILD.bazel: %w", err)
	}

	fmt.Printf("✓ Generated BUILD.bazel for %s\n", appName)
	return nil
}

// generateReactBuildContent generates the content for a React BUILD.bazel file.
func generateReactBuildContent(appName string, typescript bool) string {
	template := `load("@aspect_rules_js//js:defs.bzl", "js_library")
load("@aspect_rules_webpack//webpack:defs.bzl", "webpack_bundle")
load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

package(default_visibility = ["//visibility:public"])

# Source files
filegroup(
    name = "srcs",
    srcs = glob([
        "src/**/*.js",
        "src/**/*.jsx",
%s
        "src/**/*.css",
        "src/**/*.json",
    ]),
)

# Public files (favicon, manifest, etc.)
filegroup(
    name = "public",
    srcs = glob([
        "public/**/*",
    ]),
)
%s
# React application bundle
webpack_bundle(
    name = "build",
    srcs = [
        ":srcs",
        ":public",
        "package.json",
%s
    ],
    output_dir = True,
    webpack_config = "webpack.config.js",
    deps = [
        "//:node_modules/react",
        "//:node_modules/react-dom",
        "//:node_modules/webpack",
        "//:node_modules/webpack-cli",
        "//:node_modules/html-webpack-plugin",
        "//:node_modules/css-loader",
        "//:node_modules/style-loader",
%s
    ],
)

# Export for use by Go server
filegroup(
    name = "%s_build",
    srcs = [":build"],
    visibility = ["//visibility:public"],
)
`

	tsFiles := ""
	tsProject := ""
	tsDeps := ""
	tsProjectDep := ""

	if typescript {
		tsFiles = `        "src/**/*.ts",
        "src/**/*.tsx",`

		tsProject = `
# TypeScript compilation
ts_project(
    name = "ts",
    srcs = [":srcs"],
    declaration = True,
    tsconfig = "tsconfig.json",
    deps = [
        "//:node_modules/@types/react",
        "//:node_modules/@types/react-dom",
    ],
)
`

		tsDeps = `        "//:node_modules/@types/react",
        "//:node_modules/@types/react-dom",
        "//:node_modules/typescript",
        "//:node_modules/ts-loader",`

		tsProjectDep = `        ":ts",`
	}

	return fmt.Sprintf(template, tsFiles, tsProject, tsProjectDep, tsDeps, appName)
}

// GenerateWebpackConfig generates a webpack.config.js file for React.
func GenerateWebpackConfig(appDir string, typescript bool) error {
	config := generateWebpackConfigContent(typescript)

	configPath := filepath.Join(appDir, "webpack.config.js")

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write webpack.config.js: %w", err)
	}

	fmt.Println("✓ Generated webpack.config.js")
	return nil
}

// generateWebpackConfigContent generates the content for webpack.config.js.
func generateWebpackConfigContent(typescript bool) string {
	tsConfig := ""
	tsExtensions := ""

	if typescript {
		tsConfig = `,
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      }`
		tsExtensions = ", '.ts', '.tsx'"
	}

	return fmt.Sprintf(`const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  entry: './src/index.%s',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'static/js/[name].[contenthash:8].js',
    chunkFilename: 'static/js/[name].[contenthash:8].chunk.js',
    publicPath: '/',
    clean: true,
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env', '@babel/preset-react'],
          },
        },
      }%s,
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
      {
        test: /\.(png|svg|jpg|jpeg|gif|ico)$/,
        type: 'asset/resource',
        generator: {
          filename: 'static/media/[name].[contenthash:8][ext]',
        },
      },
    ],
  },
  resolve: {
    extensions: ['.js', '.jsx'%s],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './public/index.html',
      favicon: './public/favicon.ico',
    }),
  ],
  optimization: {
    moduleIds: 'deterministic',
    runtimeChunk: 'single',
    splitChunks: {
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          chunks: 'all',
        },
      },
    },
  },
  devtool: process.env.NODE_ENV === 'production' ? 'source-map' : 'eval-source-map',
  mode: process.env.NODE_ENV || 'production',
};
`, func() string {
		if typescript {
			return "tsx"
		}
		return "jsx"
	}(), tsConfig, tsExtensions)
}

// SetupReactBazelWorkspace sets up a React app for Bazel builds.
func SetupReactBazelWorkspace(appDir, appName string, typescript bool) error {
	fmt.Printf("Setting up Bazel workspace for %s...\n", appName)

	// Generate BUILD.bazel
	if err := GenerateReactBuildFile(appDir, appName, typescript); err != nil {
		return err
	}

	// Generate webpack.config.js
	if err := GenerateWebpackConfig(appDir, typescript); err != nil {
		return err
	}

	// Generate .bazelrc (optional, for better defaults)
	bazelrcPath := filepath.Join(appDir, ".bazelrc")
	bazelrcContent := `# React build configuration
build --strategy=Webpack=worker
build --worker_max_instances=4
build --experimental_remote_merkle_tree_cache

# Optimization for production
build:opt --compilation_mode=opt
build:opt --define=NODE_ENV=production

# Development mode
build:dev --compilation_mode=fastbuild
build:dev --define=NODE_ENV=development
`

	if err := os.WriteFile(bazelrcPath, []byte(bazelrcContent), 0644); err != nil {
		fmt.Printf("Warning: failed to write .bazelrc: %v\n", err)
	} else {
		fmt.Println("✓ Generated .bazelrc")
	}

	fmt.Println("✓ Bazel workspace setup complete!")
	fmt.Println("\nTo build with Bazel:")
	fmt.Printf("  bazel build //:%s_build\n", appName)
	fmt.Println("\nTo build optimized:")
	fmt.Printf("  bazel build --config=opt //:%s_build\n", appName)

	return nil
}

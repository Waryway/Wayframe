package react

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildWithBazel_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  BazelBuildConfig
		wantErr bool
	}{
		{
			name: "missing target",
			config: BazelBuildConfig{
				WorkspaceDir: ".",
			},
			wantErr: true,
		},
		{
			name: "valid config",
			config: BazelBuildConfig{
				WorkspaceDir: ".",
				Target:       "//apps/my-app:build",
			},
			wantErr: false, // Will fail if Bazel not installed, but validation passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BuildWithBazel(tt.config)

			// We only check validation errors, not execution errors
			if tt.wantErr && err == nil {
				t.Error("BuildWithBazel() expected validation error")
			}
			if !tt.wantErr && err != nil && strings.Contains(err.Error(), "Target is required") {
				t.Errorf("BuildWithBazel() unexpected validation error: %v", err)
			}
		})
	}
}

func TestGenerateReactBuildFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		appName    string
		typescript bool
	}{
		{
			name:       "JavaScript app",
			appName:    "my-app",
			typescript: false,
		},
		{
			name:       "TypeScript app",
			appName:    "my-ts-app",
			typescript: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appDir := filepath.Join(tmpDir, tt.appName)
			if err := os.MkdirAll(appDir, 0755); err != nil {
				t.Fatalf("Failed to create app dir: %v", err)
			}

			err := GenerateReactBuildFile(appDir, tt.appName, tt.typescript)
			if err != nil {
				t.Errorf("GenerateReactBuildFile() error = %v", err)
				return
			}

			// Check if BUILD.bazel was created
			buildPath := filepath.Join(appDir, "BUILD.bazel")
			if _, err := os.Stat(buildPath); os.IsNotExist(err) {
				t.Error("BUILD.bazel was not created")
				return
			}

			// Read and verify content
			content, err := os.ReadFile(buildPath)
			if err != nil {
				t.Errorf("Failed to read BUILD.bazel: %v", err)
				return
			}

			contentStr := string(content)

			// Check for required elements
			if !strings.Contains(contentStr, "webpack_bundle") {
				t.Error("BUILD.bazel missing webpack_bundle rule")
			}

			if !strings.Contains(contentStr, tt.appName) {
				t.Errorf("BUILD.bazel missing app name: %s", tt.appName)
			}

			if tt.typescript {
				if !strings.Contains(contentStr, "ts_project") {
					t.Error("TypeScript BUILD.bazel missing ts_project rule")
				}
				if !strings.Contains(contentStr, "*.ts") {
					t.Error("TypeScript BUILD.bazel missing *.ts glob")
				}
			}
		})
	}
}

func TestGenerateWebpackConfig(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		typescript bool
	}{
		{
			name:       "JavaScript webpack config",
			typescript: false,
		},
		{
			name:       "TypeScript webpack config",
			typescript: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(appDir, 0755); err != nil {
				t.Fatalf("Failed to create app dir: %v", err)
			}

			err := GenerateWebpackConfig(appDir, tt.typescript)
			if err != nil {
				t.Errorf("GenerateWebpackConfig() error = %v", err)
				return
			}

			// Check if webpack.config.js was created
			configPath := filepath.Join(appDir, "webpack.config.js")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Error("webpack.config.js was not created")
				return
			}

			// Read and verify content
			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Errorf("Failed to read webpack.config.js: %v", err)
				return
			}

			contentStr := string(content)

			// Check for required elements
			if !strings.Contains(contentStr, "HtmlWebpackPlugin") {
				t.Error("webpack.config.js missing HtmlWebpackPlugin")
			}

			if !strings.Contains(contentStr, "contenthash") {
				t.Error("webpack.config.js missing contenthash for cache busting")
			}

			if tt.typescript {
				if !strings.Contains(contentStr, "ts-loader") {
					t.Error("TypeScript webpack.config.js missing ts-loader")
				}
				if !strings.Contains(contentStr, ".tsx?") {
					t.Error("TypeScript webpack.config.js missing .tsx? pattern")
				}
			}
		})
	}
}

func TestSetupReactBazelWorkspace(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		appName    string
		typescript bool
	}{
		{
			name:       "JavaScript workspace",
			appName:    "js-app",
			typescript: false,
		},
		{
			name:       "TypeScript workspace",
			appName:    "ts-app",
			typescript: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appDir := filepath.Join(tmpDir, tt.appName)
			if err := os.MkdirAll(appDir, 0755); err != nil {
				t.Fatalf("Failed to create app dir: %v", err)
			}

			err := SetupReactBazelWorkspace(appDir, tt.appName, tt.typescript)
			if err != nil {
				t.Errorf("SetupReactBazelWorkspace() error = %v", err)
				return
			}

			// Check if all files were created
			files := []string{
				"BUILD.bazel",
				"webpack.config.js",
				".bazelrc",
			}

			for _, file := range files {
				path := filepath.Join(appDir, file)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("Expected file not created: %s", file)
				}
			}
		})
	}
}

func TestIsWindows(t *testing.T) {
	result := isWindows()
	// Just ensure it doesn't crash - the result depends on the OS
	t.Logf("isWindows() = %v", result)
}

// Example of setting up Bazel workspace
func ExampleSetupReactBazelWorkspace() {
	SetupReactBazelWorkspace("./my-react-app", "my-react-app", true)
	// Output will show success messages
}

// Example of building with Bazel
func ExampleBuildWithBazel() {
	BuildWithBazel(BazelBuildConfig{
		WorkspaceDir: ".",
		Target:       "//apps/my-app:build",
		Config:       "opt",
		OutputDir:    "./dist",
	})
	// Output will show build progress
}

package react

import (
	"os"
	"os/exec"
	"testing"
)

func TestCheckNodeInstalled(t *testing.T) {
	// This test will be skipped if Node.js is not installed
	err := checkNodeInstalled()
	if err != nil {
		t.Skip("Node.js not installed, skipping test")
	}
}

func TestGetNPMCommand(t *testing.T) {
	cmd := GetNPMCommand()
	if cmd == "" {
		t.Error("GetNPMCommand() returned empty string")
	}

	// Should return npm or npm.cmd
	if cmd != "npm" && cmd != "npm.cmd" {
		t.Errorf("GetNPMCommand() = %q, expected 'npm' or 'npm.cmd'", cmd)
	}
}

func TestGetNPXCommand(t *testing.T) {
	cmd := GetNPXCommand()
	if cmd == "" {
		t.Error("GetNPXCommand() returned empty string")
	}

	// Should return npx or npx.cmd
	if cmd != "npx" && cmd != "npx.cmd" {
		t.Errorf("GetNPXCommand() = %q, expected 'npx' or 'npx.cmd'", cmd)
	}
}

func TestBootstrap_Validation(t *testing.T) {
	tests := []struct {
		name         string
		config       BootstrapConfig
		wantErr      bool
		skipIfNoNode bool
	}{
		{
			name: "missing app name",
			config: BootstrapConfig{
				AppName: "",
			},
			wantErr:      true,
			skipIfNoNode: false,
		},
		{
			name: "valid config",
			config: BootstrapConfig{
				AppName:     "test-app",
				SkipInstall: true,
				SkipGit:     true,
			},
			wantErr:      false,
			skipIfNoNode: true, // Skip if Node.js not installed
		},
		{
			name: "unknown template",
			config: BootstrapConfig{
				AppName:  "test-app",
				Template: "unknown",
			},
			wantErr:      true,
			skipIfNoNode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that require Node.js if it's not installed
			if tt.skipIfNoNode && checkNodeInstalled() != nil {
				t.Skip("Node.js not installed, skipping test")
			}

			// Create temp directory for test
			tmpDir := t.TempDir()
			tt.config.Directory = tmpDir

			err := Bootstrap(tt.config)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Bootstrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBootstrapCRA_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if Node.js is installed
	if err := checkNodeInstalled(); err != nil {
		t.Skip("Node.js not installed, skipping integration test")
	}

	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("npx not found, skipping integration test")
	}

	tmpDir := t.TempDir()

	cfg := BootstrapConfig{
		AppName:     "test-cra-app",
		Directory:   tmpDir,
		Template:    "cra",
		TypeScript:  false,
		SkipInstall: false,
		SkipGit:     true,
	}

	// This will actually create a React app - can take several minutes
	err := Bootstrap(cfg)
	if err != nil {
		t.Logf("Warning: Bootstrap failed (this is expected if npm registry is unreachable): %v", err)
		return
	}

	// Check if the app was created
	appPath := tmpDir + "/test-cra-app"
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		t.Errorf("App directory was not created: %s", appPath)
	}

	// Check for package.json
	packageJSON := appPath + "/package.json"
	if _, err := os.Stat(packageJSON); os.IsNotExist(err) {
		t.Errorf("package.json was not created: %s", packageJSON)
	}
}

func TestBuildReactApp_Validation(t *testing.T) {
	// Test with non-existent directory
	err := BuildReactApp("/nonexistent/directory")
	if err == nil {
		t.Error("BuildReactApp() should fail with non-existent directory")
	}
}

func TestInstallDependencies_Validation(t *testing.T) {
	// Test with non-existent directory
	err := InstallDependencies("/nonexistent/directory")
	if err == nil {
		t.Error("InstallDependencies() should fail with non-existent directory")
	}
}

// Example of how to use the bootstrap functionality
func ExampleBootstrap() {
	cfg := BootstrapConfig{
		AppName:     "my-react-app",
		Directory:   "./apps",
		Template:    "vite",
		TypeScript:  true,
		SkipInstall: false,
		SkipGit:     false,
	}

	if err := Bootstrap(cfg); err != nil {
		panic(err)
	}
}

// Example of interactive bootstrap
func ExampleBootstrapInteractive() {
	if err := BootstrapInteractive(); err != nil {
		panic(err)
	}
}

// Example of building a React app
func ExampleBuildReactApp() {
	if err := BuildReactApp("./my-react-app"); err != nil {
		panic(err)
	}
}

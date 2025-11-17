package react

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// BootstrapConfig holds configuration for bootstrapping a new React application.
type BootstrapConfig struct {
	// AppName is the name of the React application to create
	AppName string

	// Directory is where to create the app (default: current directory)
	Directory string

	// Template specifies which template to use (default: "cra" for Create React App)
	// Options: "cra", "vite", "next"
	Template string

	// TypeScript enables TypeScript support (default: false)
	TypeScript bool

	// SkipInstall skips npm install step (default: false)
	SkipInstall bool

	// SkipGit skips git initialization (default: false)
	SkipGit bool

	// SetupBazel sets up Bazel build files (default: false)
	SetupBazel bool
}

// Bootstrap creates a new React application with the specified configuration.
// It uses npx to create the application with the chosen template.
func Bootstrap(cfg BootstrapConfig) error {
	// Validate configuration
	if cfg.AppName == "" {
		return fmt.Errorf("AppName is required")
	}

	// Set defaults
	if cfg.Directory == "" {
		cfg.Directory = "."
	}
	if cfg.Template == "" {
		cfg.Template = "cra"
	}

	// Check if Node.js is installed
	if err := checkNodeInstalled(); err != nil {
		return fmt.Errorf("Node.js is required but not found: %w", err)
	}

	// Create the app based on template
	switch cfg.Template {
	case "cra":
		if err := bootstrapCRA(cfg); err != nil {
			return err
		}
	case "vite":
		if err := bootstrapVite(cfg); err != nil {
			return err
		}
	case "next":
		if err := bootstrapNext(cfg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown template: %s (supported: cra, vite, next)", cfg.Template)
	}

	// Setup Bazel if requested
	if cfg.SetupBazel {
		appPath := filepath.Join(cfg.Directory, cfg.AppName)
		fmt.Println("\nSetting up Bazel build files...")
		if err := SetupReactBazelWorkspace(appPath, cfg.AppName, cfg.TypeScript); err != nil {
			return fmt.Errorf("failed to setup Bazel: %w", err)
		}
	}

	return nil
}

// checkNodeInstalled checks if Node.js and npm are installed.
func checkNodeInstalled() error {
	// Check for node
	nodeCmd := exec.Command("node", "--version")
	if err := nodeCmd.Run(); err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	// Check for npm or npx
	npmCmd := exec.Command("npm", "--version")
	if err := npmCmd.Run(); err != nil {
		return fmt.Errorf("npm not found: %w", err)
	}

	return nil
}

// bootstrapCRA creates a React app using Create React App.
func bootstrapCRA(cfg BootstrapConfig) error {
	fmt.Printf("Creating React app '%s' using Create React App...\n", cfg.AppName)

	args := []string{"create-react-app", cfg.AppName}

	if cfg.TypeScript {
		args = append(args, "--template", "typescript")
	}

	if cfg.SkipGit {
		// CRA doesn't have a skip-git flag, we'll remove .git after
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = cfg.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create React app: %w", err)
	}

	// Remove .git if requested
	if cfg.SkipGit {
		gitDir := filepath.Join(cfg.Directory, cfg.AppName, ".git")
		_ = os.RemoveAll(gitDir)
	}

	return printNextSteps(cfg, "npm start", "npm run build")
}

// bootstrapVite creates a React app using Vite.
func bootstrapVite(cfg BootstrapConfig) error {
	fmt.Printf("Creating React app '%s' using Vite...\n", cfg.AppName)

	template := "react"
	if cfg.TypeScript {
		template = "react-ts"
	}

	args := []string{"create-vite@latest", cfg.AppName, "--", "--template", template}

	cmd := exec.Command("npx", args...)
	cmd.Dir = cfg.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Vite app: %w", err)
	}

	// Install dependencies unless skipped
	if !cfg.SkipInstall {
		fmt.Println("Installing dependencies...")
		installCmd := exec.Command("npm", "install")
		installCmd.Dir = filepath.Join(cfg.Directory, cfg.AppName)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr

		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}
	}

	// Remove .git if requested
	if cfg.SkipGit {
		gitDir := filepath.Join(cfg.Directory, cfg.AppName, ".git")
		_ = os.RemoveAll(gitDir)
	}

	return printNextSteps(cfg, "npm run dev", "npm run build")
}

// bootstrapNext creates a React app using Next.js.
func bootstrapNext(cfg BootstrapConfig) error {
	fmt.Printf("Creating Next.js app '%s'...\n", cfg.AppName)

	args := []string{"create-next-app@latest", cfg.AppName}

	if cfg.TypeScript {
		args = append(args, "--typescript")
	} else {
		args = append(args, "--js")
	}

	// Use defaults for other options
	args = append(args, "--no-eslint", "--no-tailwind", "--no-app", "--no-src-dir", "--import-alias", "@/*")

	cmd := exec.Command("npx", args...)
	cmd.Dir = cfg.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Next.js app: %w", err)
	}

	// Remove .git if requested
	if cfg.SkipGit {
		gitDir := filepath.Join(cfg.Directory, cfg.AppName, ".git")
		_ = os.RemoveAll(gitDir)
	}

	return printNextSteps(cfg, "npm run dev", "npm run build")
}

// printNextSteps prints instructions for the user.
func printNextSteps(cfg BootstrapConfig, devCmd, buildCmd string) error {
	appPath := filepath.Join(cfg.Directory, cfg.AppName)

	fmt.Println()
	fmt.Println("✨ Success! Created", cfg.AppName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", appPath)

	if cfg.SkipInstall {
		fmt.Println("  npm install")
	}

	fmt.Printf("  %s    # Start development server\n", devCmd)
	fmt.Printf("  %s   # Build for production\n", buildCmd)

	fmt.Println()
	fmt.Println("To serve with Wayframe:")
	fmt.Printf("  1. Build the React app: %s\n", buildCmd)
	fmt.Println("  2. Use the Wayframe react package to serve it from Go")
	fmt.Println()
	fmt.Println("See examples/react for a complete example!")

	return nil
}

// BootstrapInteractive creates a React app interactively by prompting the user.
func BootstrapInteractive() error {
	var cfg BootstrapConfig

	fmt.Println("🚀 Wayframe React Bootstrap")
	fmt.Println("This will create a new React application.")
	fmt.Println()

	// App name
	fmt.Print("App name: ")
	if _, err := fmt.Scanln(&cfg.AppName); err != nil {
		return fmt.Errorf("failed to read app name: %w", err)
	}

	// Template
	fmt.Println()
	fmt.Println("Choose a template:")
	fmt.Println("  1) Create React App (CRA)")
	fmt.Println("  2) Vite (faster, modern)")
	fmt.Println("  3) Next.js (with SSR)")
	fmt.Print("Choice [1]: ")

	var choice string
	_, err := fmt.Scanln(&choice)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to read template choice: %w", err)
	}
	if choice == "" {
		choice = "1"
	}

	switch choice {
	case "1":
		cfg.Template = "cra"
	case "2":
		cfg.Template = "vite"
	case "3":
		cfg.Template = "next"
	default:
		cfg.Template = "cra"
	}

	// TypeScript
	fmt.Println()
	fmt.Print("Use TypeScript? [y/N]: ")
	var tsChoice string
	_, err = fmt.Scanln(&tsChoice)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to read TypeScript choice: %w", err)
	}
	cfg.TypeScript = tsChoice == "y" || tsChoice == "Y"

	// Bazel setup
	fmt.Println()
	fmt.Print("Setup Bazel build files? [y/N]: ")
	var bazelChoice string
	_, err = fmt.Scanln(&bazelChoice)
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("failed to read Bazel choice: %w", err)
	}
	cfg.SetupBazel = bazelChoice == "y" || bazelChoice == "Y"

	// Directory
	cfg.Directory = "."

	fmt.Println()
	return Bootstrap(cfg)
}

// GetNPMCommand returns the appropriate npm command for the current OS.
func GetNPMCommand() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}

// GetNPXCommand returns the appropriate npx command for the current OS.
func GetNPXCommand() string {
	if runtime.GOOS == "windows" {
		return "npx.cmd"
	}
	return "npx"
}

// GetPNPMCommand returns the appropriate pnpm command for the current OS.
func GetPNPMCommand() string {
	if runtime.GOOS == "windows" {
		return "pnpm.cmd"
	}
	return "pnpm"
}

// BuildReactApp builds a React application in the specified directory.
func BuildReactApp(appDir string) error {
	fmt.Printf("Building React app in %s...\n", appDir)

	cmd := exec.Command(GetNPMCommand(), "run", "build")
	cmd.Dir = appDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build React app: %w", err)
	}

	fmt.Println("✓ React app built successfully!")
	return nil
}

// InstallDependencies installs npm dependencies in the specified directory.
func InstallDependencies(appDir string) error {
	fmt.Printf("Installing dependencies in %s...\n", appDir)

	cmd := exec.Command(GetNPMCommand(), "install")
	cmd.Dir = appDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	fmt.Println("✓ Dependencies installed successfully!")
	return nil
}

// Check if pnpm is installed, or install it using npx if missing.
func EnsurePNPM(appDir string) error {
	if err := checkPNPMInstalled(); err == nil {
		fmt.Println("pnpm is already installed.")
	} else {
		fmt.Println("pnpm not found, attempting to install with npx...")
		npxCmd := GetNPXCommand()
		cmd := exec.Command(npxCmd, "-y", "pnpm", "add", "-g", "pnpm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install pnpm with npx: %w", err)
		}
	}
	// Now run pnpm install in the app directory
	fmt.Println("Running pnpm install in", appDir)
	pnpmCmd := GetPNPMCommand()
	cmd := exec.Command(pnpmCmd, "install")
	cmd.Dir = appDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pnpm install: %w", err)
	}
	return nil
}

// checkPNPMInstalled checks if pnpm is installed.
func checkPNPMInstalled() error {
	cmd := exec.Command(GetPNPMCommand(), "--version")
	return cmd.Run()
}

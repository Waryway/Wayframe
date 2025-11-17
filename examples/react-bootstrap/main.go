// Package main provides a CLI tool to bootstrap React applications for Wayframe.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Waryway/Wayframe/pkg/react"
)

var (
	appName     = flag.String("name", "", "Name of the React application to create (required)")
	directory   = flag.String("dir", ".", "Directory where to create the app")
	template    = flag.String("template", "vite", "Template to use: cra, vite, or next")
	typescript  = flag.Bool("typescript", false, "Use TypeScript")
	skipInstall = flag.Bool("skip-install", false, "Skip npm install step")
	skipGit     = flag.Bool("skip-git", false, "Skip git initialization")
	setupBazel  = flag.Bool("bazel", false, "Setup Bazel build files")
	interactive = flag.Bool("interactive", false, "Run in interactive mode")
	build       = flag.String("build", "", "Build an existing React app (provide app directory)")
	install     = flag.String("install", "", "Install dependencies for an existing React app (provide app directory)")
	buildBazel  = flag.String("build-bazel", "", "Build a React app with Bazel (provide Bazel target)")
	ensurePNPM  = flag.Bool("ensure-pnpm", false, "Ensure pnpm is installed and run pnpm install in the app directory after bootstrapping")
)

func main() {
	flag.Parse()

	// Interactive mode
	if *interactive {
		if err := react.BootstrapInteractive(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Build mode
	if *build != "" {
		if err := react.BuildReactApp(*build); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error building app: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Bazel build mode
	if *buildBazel != "" {
		if err := react.BuildWithBazel(react.BazelBuildConfig{
			Target: *buildBazel,
			Config: "opt",
		}); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error building with Bazel: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Install mode
	if *install != "" {
		if err := react.InstallDependencies(*install); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing dependencies: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Bootstrap mode (default)
	if *appName == "" {
		fmt.Fprintln(os.Stderr, "Error: -name is required")
		fmt.Fprintln(os.Stderr, "\nUsage:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "\nExamples:")
		_, _ = fmt.Fprintln(os.Stderr, "  # Interactive mode")
		_, _ = fmt.Fprintln(os.Stderr, "  wayframe-react -interactive")
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "  # Create a Vite app with TypeScript and Bazel")
		_, _ = fmt.Fprintln(os.Stderr, "  wayframe-react -name my-app -template vite -typescript -bazel")
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "  # Create a Create React App")
		_, _ = fmt.Fprintln(os.Stderr, "  wayframe-react -name my-app -template cra")
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "  # Build an existing app")
		_, _ = fmt.Fprintln(os.Stderr, "  wayframe-react -build ./my-app")
		_, _ = fmt.Fprintln(os.Stderr, "")
		_, _ = fmt.Fprintln(os.Stderr, "  # Build with Bazel")
		_, _ = fmt.Fprintln(os.Stderr, "  wayframe-react -build-bazel //apps/my-app:build")
		os.Exit(1)
	}

	cfg := react.BootstrapConfig{
		AppName:     *appName,
		Directory:   *directory,
		Template:    *template,
		TypeScript:  *typescript,
		SkipInstall: *skipInstall,
		SkipGit:     *skipGit,
		SetupBazel:  *setupBazel,
	}

	if err := react.Bootstrap(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *ensurePNPM {
		appPath := fmt.Sprintf("%s/%s", *directory, *appName)
		if err := react.EnsurePNPM(appPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error ensuring pnpm: %v\n", err)
			os.Exit(1)
		}
	}
}

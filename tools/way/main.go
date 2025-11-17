package main

import (
	"fmt"
	"github.com/Waryway/Wayframe/pkg/cli"
	"github.com/Waryway/Wayframe/pkg/logger"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var log = logger.New(logger.ErrorLevel)

func main() {
	// For this example, use Kong as the CLI backend
	var cliApp struct {
		Setup SetupCmd `cmd:"setup" help:"Setup the Wayframe repository (pnpm, lockfile, etc.)"`
		Node  NodeCmd  `cmd:"node" help:"Check for Node.js and prompt if not found"`
		Pnpm  PnpmCmd  `cmd:"pnpm" help:"Check for pnpm and install if missing"`
	}
	app := cli.NewKongCLI(&cliApp)

	args := os.Args[1:]
	if len(args) > 2 {
		_, err := fmt.Fprintln(os.Stderr, "Error: Too many arguments. Maximum 2 allowed.")
		if err != nil {
			return
		}
		os.Exit(1)
	}

	if len(args) > 0 {
		// Single command mode (argument launcher)
		appArgs := append([]string{os.Args[0]}, args...)
		os.Args = appArgs
		if err := app.Run(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Interactive command loop (REPL)
	fmt.Println("Wayframe Interactive CLI. Type 'help' for commands, 'exit' to quit.")
	for {
		fmt.Print("> ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			if err.Error() == "unexpected newline" || err.Error() == "EOF" {
				continue
			}
			_, err := fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
			if err != nil {
				return
			}
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting.")
			break
		}
		if input == "help" {
			fmt.Println("Available commands: setup, help, exit")
			continue
		}
		// Split input into args (max 2)
		cmdArgs := []string{}
		for _, arg := range splitArgs(input) {
			if len(cmdArgs) < 2 {
				cmdArgs = append(cmdArgs, arg)
			}
		}
		if len(cmdArgs) == 0 {
			continue
		}
		// Use RunArgs to pass the command directly
		//app := cli.NewKongCLI(&cliApp)
		if err := app.RunArgs(cmdArgs); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}

// splitArgs splits a string into arguments (max 2, space separated)
func splitArgs(input string) []string {
	args := []string{}
	curr := ""
	for i := 0; i < len(input); i++ {
		if input[i] == ' ' {
			if curr != "" {
				args = append(args, curr)
				curr = ""
				if len(args) == 2 {
					break
				}
			}
		} else {
			curr += string(input[i])
		}
	}
	if curr != "" && len(args) < 2 {
		args = append(args, curr)
	}
	return args
}

type SetupCmd struct{}

func (s *SetupCmd) Run() error {
	fmt.Println("Setting up the Wayframe repository...")
	if err := ensurePNPM(); err != nil {
		return err
	}
	if err := ensureNode(); err != nil {
		return err
	}
	if err := runPNPMInstall(); err != nil {
		return err
	}
	fmt.Println("Wayframe repository setup complete.")
	return nil
}

// yesNoPrompt prompts the user for a yes/no answer and returns true for yes, false for no.
func yesNoPrompt(prompt string) bool {
	for {
		fmt.Printf("%s [y/N]: ", prompt)
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil && err.Error() != "unexpected newline" && err.Error() != "EOF" {
			fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
			return false
		}
		if input == "y" || input == "Y" {
			return true
		}
		if input == "n" || input == "N" || input == "" {
			return false
		}
		fmt.Println("Please enter 'y' or 'n'.")
	}
}

// NodeCmd allows running ensureNode directly from CLI
// Usage: way node
// Description: Checks for Node.js and prompts if not found

type NodeCmd struct{}

func (n *NodeCmd) Run() error {
	return ensureNode()
}

// PnpmCmd allows running ensurePNPM directly from CLI
// Usage: way pnpm
// Description: Checks for pnpm and installs if missing

type PnpmCmd struct{}

func (p *PnpmCmd) Run() error {
	return ensurePNPM()
}

func ensureNode() error {
	if _, err := os.Stat("node_modules"); err == nil {
		log.Error("node_modules already exists")
		return nil // Already set up
	}
	if _, err := execLookPath("node"); err == nil {
		return nil // node already available
	}
	log.Errorf("Node.js is required but not found in PATH.")
	if !yesNoPrompt("Node.js is missing. Would you like to install Node.js automatically using pnpm?") {
		log.Errorf("Node.js was not installed. Please visit https://nodejs.org/en/download for manual installation.")
		return fmt.Errorf("Node.js is required but not found in PATH. Please install node from https://nodejs.org/en/download")
	}
	// Try to install Node.js using pnpm
	if _, err := execLookPath("pnpm"); err != nil {
		log.Errorf("pnpm is required to install Node.js automatically. Please run 'way pnpm' or install pnpm first.")
		return fmt.Errorf("pnpm is required to install Node.js automatically. Please run 'way pnpm' or install pnpm first.")
	}
	cmd := execCommand("pnpm", "env", "use", "--global", "lts")
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		output := outBuf.String() + errBuf.String()
		if strings.Contains(output, "ERR_PNPM_CANNOT_MANAGE_NODE") {
			log.Errorf("pnpm cannot manage Node.js because it was not installed using the standalone script. Please install Node.js manually from https://nodejs.org/en/download or reinstall pnpm using the standalone script from https://pnpm.io/installation.")
			fmt.Println("\nTo fix this:")
			fmt.Println("- Install Node.js manually from https://nodejs.org/en/download")
			fmt.Println("- Or uninstall all Node.js versions and install pnpm using the standalone script from https://pnpm.io/installation")
			return fmt.Errorf("pnpm cannot manage Node.js. Please install Node.js manually or reinstall pnpm as described above.")
		}
		log.Errorf("failed to install Node.js with pnpm: %v. Please visit https://nodejs.org/en/download for manual installation.", err)
		return fmt.Errorf("failed to install Node.js with pnpm: %w. Please install node from https://nodejs.org/en/download", err)
	}
	log.Info("Node.js installed successfully using pnpm.")
	return nil
}

func ensurePNPM() error {
	fullpath, err := execLookPath("pnpm")
	if err == nil {
		log.WithField("pnpm", fullpath).Info("pnpm is available in PATH and is standalone.")
		return nil
	}
	log.Errorf("pnpm is required but not found in PATH.")
	if !yesNoPrompt("pnpm is missing. Would you like to install pnpm automatically?") {
		log.Errorf("pnpm was not installed. Please visit https://pnpm.io/installation for manual installation.")
		return fmt.Errorf("pnpm is required but not found in PATH. Please install pnpm from https://pnpm.io/installation")
	}

	// Detect if running in a bash-like environment (Git Bash, MSYS2, Cygwin, WSL)
	shell := os.Getenv("SHELL")
	if strings.Contains(strings.ToLower(shell), "bash") || strings.Contains(strings.ToLower(shell), "msys") || strings.Contains(strings.ToLower(shell), "cygwin") || strings.Contains(strings.ToLower(shell), "wsl") {
		// Use POSIX install script for bash-like environments
		var shCmd string
		if _, err := execLookPath("curl"); err == nil {
			shCmd = "curl -fsSL https://get.pnpm.io/install.sh | sh -"
		} else if _, err := execLookPath("wget"); err == nil {
			shCmd = "wget -qO- https://get.pnpm.io/install.sh | sh -"
		} else {
			log.Errorf("Neither curl nor wget is available to install pnpm. Please install pnpm from https://pnpm.io/installation")
			return fmt.Errorf("Neither curl nor wget is available to install pnpm. Please install pnpm from https://pnpm.io/installation")
		}
		cmd := execCommand("sh", "-c", shCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Errorf("failed to install pnpm with script: %v. Please visit https://pnpm.io/installation for manual installation.", err)
			return fmt.Errorf("failed to install pnpm with script: %w. Please install pnpm from https://pnpm.io/installation", err)
		}
		fullpath, err = execLookPath("pnpm")
		if err != nil {
			log.Errorf("pnpm was installed but is still not found in PATH. Please open a new terminal or add the pnpm binary location to your PATH.")
			return fmt.Errorf("pnpm was installed but is still not found in PATH. Please open a new terminal or add the pnpm binary location to your PATH")
		}
		fmt.Println("Detected pnpm binary:", fullpath)
		cmd = execCommand(fullpath, "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		log.WithField("pnpm", fullpath).Info("pnpm installed and available in PATH and is standalone.")
		// Install Node.js using pnpm
		if err := installNodeWithPNPM(fullpath); err != nil {
			return err
		}
		return nil
	}

	// Fallback: Use PowerShell script for native Windows shells
	if runtime.GOOS == "windows" {
		pwsh, err := execLookPath("powershell")
		if err != nil {
			log.Errorf("PowerShell is required to install pnpm automatically. Please install pnpm from https://pnpm.io/installation")
			return fmt.Errorf("PowerShell is required to install pnpm automatically. Please install pnpm from https://pnpm.io/installation")
		}
		cmd := execCommand(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "Invoke-WebRequest https://get.pnpm.io/install.ps1 -UseBasicParsing | Invoke-Expression")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Errorf("failed to install pnpm with PowerShell: %v. Please visit https://pnpm.io/installation for manual installation.", err)
			return fmt.Errorf("failed to install pnpm with PowerShell: %w. Please install pnpm from https://pnpm.io/installation", err)
		}
		// Try to find pnpm.exe in common install locations and update PATH
		pnpmHome := os.Getenv("PNPM_HOME")
		if pnpmHome == "" {
			localAppData := os.Getenv("LOCALAPPDATA")
			if localAppData != "" {
				pnpmHome = localAppData + "\\pnpm"
			}
			if pnpmHome == "" || !dirExists(pnpmHome) {
				userProfile := os.Getenv("USERPROFILE")
				if userProfile != "" {
					pnpmHome = userProfile + "\\AppData\\Local\\pnpm"
				}
			}
		}
		if pnpmHome != "" && dirExists(pnpmHome) {
			path := os.Getenv("PATH")
			if !strings.Contains(strings.ToLower(path), strings.ToLower(pnpmHome)) {
				_ = os.Setenv("PATH", pnpmHome+string(os.PathListSeparator)+path)
			}
		}
		fullpath, err = execLookPath("pnpm")
		if err != nil {
			log.Errorf("pnpm was installed but is still not found in PATH. Please open a new terminal or add %s to your PATH.", pnpmHome)
			return fmt.Errorf("pnpm was installed but is still not found in PATH. Please open a new terminal or add %s to your PATH", pnpmHome)
		}
		fmt.Println("Detected pnpm binary:", fullpath)
		cmd = execCommand(fullpath, "--version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		log.WithField("pnpm", fullpath).Info("pnpm installed and available in PATH and is standalone.")
		// Install Node.js using pnpm
		if err := installNodeWithPNPM(fullpath); err != nil {
			return err
		}
		log.Info("If you experience slow installs, consider running: Add-MpPreference -ExclusionPath $(pnpm store path) in an administrator PowerShell window.")
		return nil
	}

	log.Errorf("Could not determine shell environment for pnpm install. Please install pnpm manually from https://pnpm.io/installation.")
	return fmt.Errorf("Could not determine shell environment for pnpm install. Please install pnpm manually from https://pnpm.io/installation.")
}

// installNodeWithPNPM installs Node.js using pnpm env use --global lts
func installNodeWithPNPM(pnpmPath string) error {
	fmt.Println("Installing Node.js using pnpm...")
	cmd := execCommand(pnpmPath, "env", "use", "--global", "lts")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Failed to install Node.js with pnpm: %v. Please install Node.js manually from https://nodejs.org/en/download or use pnpm to manage Node.js.", err)
		return fmt.Errorf("Failed to install Node.js with pnpm: %w. Please install Node.js manually from https://nodejs.org/en/download or use pnpm to manage Node.js.", err)
	}
	fmt.Println("Node.js installed successfully using pnpm.")
	return nil
}

// execLookPath and execCommand are wrappers for testability
var execLookPath = func(file string) (string, error) { return exec.LookPath(file) }
var execCommand = func(name string, arg ...string) *exec.Cmd { return exec.Command(name, arg...) }

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func runPNPMInstall() error {
	fmt.Println("Running pnpm install at repo root...")
	cmd := execCommand("pnpm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("failed to run pnpm install: %v", err)
		return fmt.Errorf("failed to run pnpm install: %w", err)
	}
	return nil
}

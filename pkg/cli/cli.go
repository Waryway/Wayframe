package cli

// Command is a generic interface for CLI commands.
type Command interface {
	Run() error
}

// CLI is a generic interface for CLI applications.
type CLI interface {
	AddCommand(name string, cmd Command)
	Run() error
	RunArgs(args []string) error
}

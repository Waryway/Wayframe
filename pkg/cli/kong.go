package cli

import (
	"github.com/alecthomas/kong"
)

type kongCLI struct {
	app *kong.Kong
}

type kongCommand struct {
	runFunc func() error
}

func (k *kongCommand) Run() error {
	return k.runFunc()
}

func NewKongCLI(ctx interface{}) CLI {
	app := kong.Must(ctx)
	return &kongCLI{app: app}
}

func (k *kongCLI) AddCommand(name string, cmd Command) {
	// Kong uses struct tags for commands, so this is a no-op for compatibility
}

func (k *kongCLI) Run() error {
	ctx, err := k.app.Parse(nil)
	if err != nil {
		return err
	}
	return ctx.Run()
}

func (k *kongCLI) RunArgs(args []string) error {
	ctx, err := k.app.Parse(args)
	if err != nil {
		return err
	}
	return ctx.Run()
}

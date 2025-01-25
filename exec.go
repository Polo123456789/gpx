package main

import "context"

type ExecCommand struct {
	pkg     Package
	binPath string
}

var _ Command = &ExecCommand{}

func (c *ExecCommand) Run(ctx context.Context, args []string) error {
	if !PackageIsInstalled(c.pkg, c.binPath) {
		if err := InstallPackage(ctx, c.pkg, c.binPath); err != nil {
			return err
		}
	}

	return ExecPackage(ctx, c.pkg, c.binPath, args...)
}

func (c *ExecCommand) Name() string {
	return c.pkg.CommandName()
}

func (c *ExecCommand) Synopsis() string {
	return "Execute " + c.pkg.BinName()
}

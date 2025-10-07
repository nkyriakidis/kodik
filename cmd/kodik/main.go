package main

import (
	"log"
	"os"

	"kodik/internal/kodik"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kodik",
		Usage: "Manage kodik repository configurations (preserves user files; provides backups)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Skip backups and confirmations",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Preview actions without execution",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "github",
				Usage: "Install/update .github directory (selective merge; preserves workflows, CODEOWNERS, custom files)",
				Action: func(c *cli.Context) error {
					force := c.Bool("force")
					dryRun := c.Bool("dry-run")
					if err := kodik.InstallOrUpdateGithub(force, dryRun); err != nil {
						if ke, ok := err.(*kodik.KodikError); ok {
							return cli.Exit(ke.Error(), ke.Code)
						}
						return err
					}
					return nil
				},
			},
			{
				Name:  "roo",
				Usage: "Install/update .roomodes file (graceful when missing; fresh install safe)",
				Action: func(c *cli.Context) error {
					force := c.Bool("force")
					dryRun := c.Bool("dry-run")
					if err := kodik.InstallOrUpdateRoomodes(force, dryRun); err != nil {
						if ke, ok := err.(*kodik.KodikError); ok {
							return cli.Exit(ke.Error(), ke.Code)
						}
						return err
					}
					return nil
				},
			},
			{
				Name:  "opencode",
				Usage: "Install/update .opencode directory (graceful when missing; creates placeholder)",
				Action: func(c *cli.Context) error {
					force := c.Bool("force")
					dryRun := c.Bool("dry-run")
					if err := kodik.InstallOrUpdateOpencode(force, dryRun); err != nil {
						if ke, ok := err.(*kodik.KodikError); ok {
							return cli.Exit(ke.Error(), ke.Code)
						}
						return err
					}
					return nil
				},
			},
			{
				Name:  "all",
				Usage: "Install/update all components (preservation + graceful handling across github, roomodes, opencode)",
				Action: func(c *cli.Context) error {
					force := c.Bool("force")
					dryRun := c.Bool("dry-run")
					if err := kodik.InstallOrUpdateGithub(force, dryRun); err != nil {
						if ke, ok := err.(*kodik.KodikError); ok {
							return cli.Exit(ke.Error(), ke.Code)
						}
						return err
					}
					if err := kodik.InstallOrUpdateRoomodes(force, dryRun); err != nil {
						if ke, ok := err.(*kodik.KodikError); ok {
							return cli.Exit(ke.Error(), ke.Code)
						}
						return err
					}
					if err := kodik.InstallOrUpdateOpencode(force, dryRun); err != nil {
						if ke, ok := err.(*kodik.KodikError); ok {
							return cli.Exit(ke.Error(), ke.Code)
						}
						return err
					}
					println("All components installed/updated successfully.")
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

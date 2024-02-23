package main

import (
	"fmt"
	"gophercises/secrets/secret"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filepath := homeDir + "/.secret"

	keyFlag := &cli.StringFlag{
		Name:  "k",
		Value: "",
		Usage: "Key for encrypting/decrypting the secrets",
	}

	app := &cli.App{
		Name:  "secret",
		Usage: "Manages your secrets from the CLI",
		Commands: []*cli.Command{
			{
				Name:  "set",
				Usage: "Sets a secret in the vault",
				Flags: []cli.Flag{
					keyFlag,
				},
				Action: func(ctx *cli.Context) error {
					key := ctx.String("k")
					if key == "" {
						return fmt.Errorf("please provide an encryption key")
					}

					if ctx.Args().Len() != 2 {
						return fmt.Errorf("invalid number of arguments")
					}

					name := ctx.Args().Get(0)
					if name == "" {
						return fmt.Errorf("please provide a name")
					}
					value := ctx.Args().Get(1)
					if value == "" {
						return fmt.Errorf("please provide a value")
					}

					v := secret.NewVault(key, filepath)

					err := v.Set(name, value)
					return err
				},
			},
			{
				Name:  "get",
				Usage: "Gets a secret from the vault",
				Flags: []cli.Flag{
					keyFlag,
				},
				Action: func(ctx *cli.Context) error {
					key := ctx.String("k")
					if key == "" {
						return fmt.Errorf("please provide an encryption key")
					}

					if ctx.Args().Len() != 1 {
						return fmt.Errorf("invalid number of arguments")
					}

					name := ctx.Args().Get(0)
					if name == "" {
						return fmt.Errorf("please provide a name")
					}

					v := secret.NewVault(key, filepath)
					if val, ok, err := v.Get(name); err != nil {
						return err
					} else if ok {
						fmt.Println(val)
					} else {
						fmt.Println("no value found")
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

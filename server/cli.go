package server

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/spf13/cobra"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func MakeCli() {
	var api huma.API

	// Create a CLI app which takes a port option.
	cli := humacli.New(RunKaraberus(&api))

	// Add a command to print the OpenAPI spec.
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			// Use downgrade to return OpenAPI 3.0.3 YAML since oapi-codegen doesn't
			// support OpenAPI 3.1 fully yet. Use `.YAML()` instead for 3.1.
			b, _ := api.OpenAPI().DowngradeYAML()
			fmt.Println(string(b))
		},
	})

	cli.Root().AddCommand(&cobra.Command{
		Use:   "create-token",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := CreateSystemToken()
			if err != nil {
				getLogger().Fatalln(err)
			}
			println(token)
		},
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}

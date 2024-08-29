package server

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func MakeCli() {
	app, api := setupKaraberus()

	rootCmd := &cobra.Command{
		Use:   "karaberus",
		Short: "Start the karaberus server",
		Run: func(cmd *cobra.Command, args []string) {
			RunKaraberus(app, api)
		},
	}

	// Add a command to print the OpenAPI spec.
	rootCmd.AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			// Use downgrade to return OpenAPI 3.0.3 YAML since oapi-codegen doesn't
			// support OpenAPI 3.1 fully yet. Use `.YAML()` instead for 3.1.
			b, _ := api.OpenAPI().DowngradeYAML()
			fmt.Println(string(b))
		},
	})

	create_user_admin_flag := false
	create_user_cmd := &cobra.Command{
		Use:   "create-user <user_id>",
		Short: "Create an user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			db := GetDB(cmd.Context())
			user := User{ID: args[0], Admin: create_user_admin_flag}
			err := db.Create(&user).Error
			if err != nil {
				panic(err)
			}
			if user.Admin {
				fmt.Printf("created administrator %s.", user.ID)
			} else {
				fmt.Printf("created user %s.", user.ID)
			}
		},
	}
	create_user_cmd.Flags().BoolVar(&create_user_admin_flag, "admin", false, "create an administrator.")

	rootCmd.AddCommand(create_user_cmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "create-token <user_id> <name>",
		Short: "Create a token for the given user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			db := GetDB(cmd.Context())
			user := User{ID: args[0]}
			if err := db.First(&user).Error; err != nil {
				panic(err)
			}
			token, err := createTokenForUser(context.TODO(), user, args[1], AllScopes)
			if err != nil {
				panic(err)
			}
			fmt.Println(token.Token)
		},
	})

	rootCmd.PersistentFlags().IntVarP(
		&CONFIG.Listen.Port,
		"port", "p",
		CONFIG.Listen.Port,
		"Port to listen on",
	)

	// Run the CLI. When passed no commands, it starts the server.
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

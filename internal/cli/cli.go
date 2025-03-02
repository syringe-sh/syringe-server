package cli

import (
	"context"
	"fmt"

	"github.com/nixpig/syringe.sh/internal/cmd"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "syringe",
		Short:   "",
		Long:    "",
		Example: "",
		Version: "",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(
		setCmd(),
		getCmd(),
		listCmd(),
		deleteCmd(),
	)

	return cmd
}

func setCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set [flags] KEY VALUE",
		Short:   "Set a key-value",
		Args:    cobra.ExactArgs(2),
		Example: "  syringe set username nixpig",
		RunE: func(c *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			if err := cmd.Set(context.Background(), key, []byte(value)); err != nil {
				return fmt.Errorf("set: %w", err)
			}

			return nil
		},
	}

	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get [flags] KEY",
		Short:   "Get a value from the store",
		Args:    cobra.ExactArgs(1),
		Example: "  syringe get username",
		RunE: func(c *cobra.Command, args []string) error {
			key := args[0]

			value, err := cmd.Get(context.Background(), key)
			if err != nil {
				return fmt.Errorf("set: %w", err)
			}

			fmt.Printf("%s", value)

			return nil
		},
	}

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [flags] KEY",
		Short:   "Delete a record from the store",
		Args:    cobra.ExactArgs(1),
		Example: "  syringe delete username",
		RunE: func(c *cobra.Command, args []string) error {
			key := args[0]

			if err := cmd.Delete(context.Background(), key); err != nil {
				return fmt.Errorf("delete: %w", err)
			}

			return nil
		},
	}

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [flags]",
		Short:   "List all records in store",
		Args:    cobra.ExactArgs(0),
		Example: "  syringe list",
		RunE: func(c *cobra.Command, args []string) error {
			list, err := cmd.List(context.Background())
			if err != nil {
				return fmt.Errorf("list: %w", err)
			}

			fmt.Printf("%v", list)

			return nil
		},
	}

	return cmd
}

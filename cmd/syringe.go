package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/nixpig/syringe.sh/internal/syringe"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func New(v *viper.Viper) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "syringe",
		Short:   "Encrypted key-value store",
		Version: "",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			applyFlags(c, v)

			identity, _ := c.Flags().GetString("identity")
			log.Debug("", "identity", identity)

			store, _ := c.Flags().GetString("store")
			log.Debug("", "store", store)

			// TODO: get identity and set on context??

			// TODO: make database connection and set on ctx

			return nil
		},
	}

	rootCmd.PersistentFlags().StringP(
		"identity",
		"i",
		"",
		fmt.Sprintf("Path to SSH key (default: %s)", v.GetString("identity")),
	)

	rootCmd.PersistentFlags().StringP(
		"store",
		"s",
		"",
		fmt.Sprintf("Store as name, path or URL (default: %s)", v.GetString("store")),
	)

	bindFlags(rootCmd, v)

	rootCmd.AddCommand(
		setCmd,
		getCmd,
		listCmd,
		deleteCmd,
	)

	return rootCmd
}

var setCmd = &cobra.Command{
	Use:     "set [flags] KEY VALUE",
	Short:   "Set a key-value",
	Args:    cobra.ExactArgs(2),
	Example: "  syringe set username nixpig",
	RunE: func(c *cobra.Command, args []string) error {
		return syringe.Set(c.Context(), c.OutOrStdout(), args[0], args[1])
	},
}

var getCmd = &cobra.Command{
	Use:     "get [flags] KEY",
	Short:   "Get a value from the store",
	Args:    cobra.ExactArgs(1),
	Example: "  syringe get username",
	RunE: func(c *cobra.Command, args []string) error {
		return syringe.Get(c.Context(), c.OutOrStdout(), args[0])
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete [flags] KEY",
	Short:   "Delete a record from the store",
	Args:    cobra.ExactArgs(1),
	Example: "  syringe delete username",
	RunE: func(c *cobra.Command, args []string) error {
		return syringe.Delete(c.Context(), c.OutOrStdout(), args[0])
	},
}

var listCmd = &cobra.Command{
	Use:     "list [flags]",
	Short:   "List all records in store",
	Args:    cobra.ExactArgs(0),
	Example: "  syringe list",
	RunE: func(c *cobra.Command, args []string) error {
		return syringe.List(c.Context(), c.OutOrStdout())
	},
}

func bindFlags(c *cobra.Command, v *viper.Viper) {
	c.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		v.BindPFlag(f.Name, f)
	})
}

func applyFlags(c *cobra.Command, v *viper.Viper) {
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if v.IsSet(f.Name) {
			c.Flags().Set(f.Name, v.GetString(f.Name))
		}
	})
}

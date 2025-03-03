package cli

import (
	"fmt"

	"github.com/nixpig/syringe.sh/internal/cmd"
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
			return nil
		},

		RunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	rootCmd.PersistentFlags().StringP("identity", "i", "", "Path to SSH key (optional)")
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
		return cmd.Set(args[0], args[1], c.OutOrStdout())
	},
}

var getCmd = &cobra.Command{
	Use:     "get [flags] KEY",
	Short:   "Get a value from the store",
	Args:    cobra.ExactArgs(1),
	Example: "  syringe get username",
	RunE: func(c *cobra.Command, args []string) error {
		// fmt.Println("here: ", c.Context().Value("qux"))
		v, e := c.Flags().GetString("identity")
		if e != nil {
			fmt.Println("e: ", e)
		}
		fmt.Println("here: ", v)
		return cmd.Get(args[0], c.OutOrStdout())
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete [flags] KEY",
	Short:   "Delete a record from the store",
	Args:    cobra.ExactArgs(1),
	Example: "  syringe delete username",
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.Delete(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:     "list [flags]",
	Short:   "List all records in store",
	Args:    cobra.ExactArgs(0),
	Example: "  syringe list",
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.List(c.OutOrStdout())
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

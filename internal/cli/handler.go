package cli

import (
	"errors"
	"fmt"
	"io"
	"os/user"
	"strings"

	"github.com/nixpig/syringe.sh/pkg"
	"github.com/nixpig/syringe.sh/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

func NewHandlerCLI(host string, port int, out io.Writer) pkg.CobraHandler {
	return func(cmd *cobra.Command, args []string) error {
		currentUser, err := user.Current()
		if err != nil || currentUser.Username == "" {
			return fmt.Errorf("failed to determine username: %w", err)
		}

		identity, err := cmd.Flags().GetString("identity")
		if err != nil {
			return err
		}

		if identity == "" {
			return errors.New("no identity provided")
		}

		authMethod, err := ssh.AuthMethod(identity, cmd.OutOrStdout())
		if err != nil {
			return err
		}

		configFile, err := ssh.ConfigFile()
		if err != nil {
			return err
		}

		defer configFile.Close()

		if err := ssh.AddIdentityToSSHConfig(identity, configFile); err != nil {
			return fmt.Errorf("failed to add or update identity in ssh config file: %w", err)
		}

		client, err := ssh.NewSSHClient(
			host,
			port,
			currentUser.Username,
			authMethod,
		)
		if err != nil {
			return err
		}

		defer client.Close()

		if cmd.CalledAs() == "inject" {
			sshcmd := buildCommand(cmd, args)

			privateKey, err := ssh.GetPrivateKey(identity, cmd.OutOrStderr(), term.ReadPassword)
			if err != nil {
				return fmt.Errorf("failed to read private key: %w", err)
			}

			if err := client.Run(
				sshcmd,
				InjectResponseParser{
					w:          out,
					privateKey: privateKey,
				},
			); err != nil {
				return err
			}

			return nil
		}

		if cmd.Parent().Use == "secret" {
			switch cmd.CalledAs() {

			case "list":
				sshcmd := buildCommand(cmd, args)

				privateKey, err := ssh.GetPrivateKey(identity, cmd.OutOrStderr(), term.ReadPassword)
				if err != nil {
					return fmt.Errorf("failed to read private key: %w", err)
				}

				if err := client.Run(
					sshcmd,
					ListResponseParser{
						w:          out,
						privateKey: privateKey,
					},
				); err != nil {
					return err
				}

				return nil

			case "get":
				sshcmd := buildCommand(cmd, args)

				privateKey, err := ssh.GetPrivateKey(identity, cmd.OutOrStderr(), term.ReadPassword)
				if err != nil {
					return fmt.Errorf("failed to read private key: %w", err)
				}

				if err := client.Run(
					sshcmd,
					GetResponseParser{
						w:          out,
						privateKey: privateKey,
					},
				); err != nil {
					return err
				}

				return nil
			}
		}

		sshcmd := buildCommand(cmd, args)
		if err := client.Run(sshcmd, out); err != nil {
			return err
		}

		return nil
	}
}

func buildCommand(cmd *cobra.Command, args []string) string {
	var flags string

	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Name == "identity" {
			return
		}

		flags = fmt.Sprintf("%s --%s %s", flags, flag.Name, flag.Value)
	})

	scmd := []string{
		strings.Join(strings.Split(cmd.CommandPath(), " ")[1:], " "),
		strings.Join(args, " "),
		flags,
	}

	return strings.Join(scmd, " ")
}

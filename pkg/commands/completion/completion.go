package completion

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	bashCmd := &cobra.Command{
		Use:   "bash",
		Short: "Generate the autocompletion script for bash",
		Long:  "Generate the autocompletion script for bash",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		},
	}

	zshCmd := &cobra.Command{
		Use:   "zsh",
		Short: "Generate the autocompletion script for zsh",
		Long:  "Generate the autocompletion script for zsh",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		},
	}

	fishCmd := &cobra.Command{
		Use:   "fish",
		Short: "Generate the autocompletion script for fish",
		Long:  "Generate the autocompletion script for fish",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		},
	}

	powershellCmd := &cobra.Command{
		Use:   "powershell",
		Short: "Generate the autocompletion script for powershell",
		Long:  "Generate the autocompletion script for powershell",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
		},
	}

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate the autocompletion script for norden for the specified shell",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(bashCmd)
	cmd.AddCommand(zshCmd)
	cmd.AddCommand(fishCmd)
	cmd.AddCommand(powershellCmd)

	return cmd
}

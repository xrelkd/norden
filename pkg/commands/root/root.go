package root

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	"github.com/xrelkd/norden/pkg/commands/attach"
	"github.com/xrelkd/norden/pkg/commands/completion"
	"github.com/xrelkd/norden/pkg/commands/create"
	"github.com/xrelkd/norden/pkg/commands/delete"
	"github.com/xrelkd/norden/pkg/commands/list"
	"github.com/xrelkd/norden/pkg/commands/version"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "norden [command]",
		Short:        "Norden",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			HiddenDefaultCmd:    true,
			DisableDescriptions: true,
		},
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				logger.Debug("ignoring cobra error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(create.Command())
	cmd.AddCommand(attach.Command())
	cmd.AddCommand(list.Command())
	cmd.AddCommand(delete.Command())

	cmd.AddCommand(version.Command())
	cmd.AddCommand(completion.Command())

	return cmd
}

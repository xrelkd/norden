package image

import (
	"github.com/spf13/cobra"

	"github.com/xrelkd/norden/pkg/commands/image/list"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Short:   "Manage images",
		Long:    "Manage images",
		Aliases: []string{"i"},
		RunE:    list.RunE,
	}

	cmd.AddCommand(list.Command())

	return cmd
}

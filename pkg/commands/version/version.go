package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xrelkd/norden/pkg/version"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current version",
		Long:  "Show current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s %s\n", version.AppName, version.Version)
		},
	}
}

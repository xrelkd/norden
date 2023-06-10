package list

import (
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/xrelkd/norden/pkg/config"
)

func runList() error {
	conf, err := config.Load()
	if err != nil {
		logger.Warning("%v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "IMAGE", "PULL POLICY", "INTERACTIVE SHELL", "COMMAND", "ARGS"})
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	for i := range conf.Images {
		image := &conf.Images[i]

		table.Append([]string{
			image.Name,
			image.Image,
			string(image.ImagePullPolicy),
			strings.Join(image.InteractiveShell, " "),
			strings.Join(image.Command, " "),
			strings.Join(image.Args, " "),
		})
	}

	table.Render()

	return nil
}

func RunE(_ *cobra.Command, _ []string) error {
	return runList()
}

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List images",
		Long:    "List images",
		Aliases: []string{"l", "ls"},
		RunE:    RunE,
	}

	return cmd
}

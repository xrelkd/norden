package list

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xrelkd/norden/internal/cmdutils"
)

type ListOptions struct {
	Namespace     string
	AllNamespaces bool
}

func runList(opts *ListOptions) error {
	clientset, _, err := cmdutils.CreateClientset()
	if err != nil {
		return err
	}

	if opts.Namespace == "" && !opts.AllNamespaces {
		var ns string
		if ns, err = cmdutils.GetCurrentNamespace(); err != nil {
			return err
		} else {
			opts.Namespace = ns
		}
	}

	if opts.AllNamespaces {
		opts.Namespace = ""
	}

	pods, err := clientset.
		CoreV1().
		Pods(opts.Namespace).
		List(context.Background(), metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/managed-by=norden",
		})
	if err != nil {
		return err
	}

	displayPods(pods)

	return nil
}

func displayPods(pods *v1.PodList) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "IMAGE", "STATUS", "NAMESPACE", "NODE"})
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

	for i := range pods.Items {
		pod := &pods.Items[i]
		table.Append([]string{
			pod.Name,
			pod.Spec.Containers[0].Image,
			string(pod.Status.Phase),
			pod.Namespace,
			pod.Spec.NodeName,
		})
	}

	table.Render()
}

func Command() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all pods created by norden",
		Long:    "List all pods created by norden",
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "", "Namespace, use current namespace if not provided")
	cmd.Flags().BoolVarP(&opts.AllNamespaces, "all-namespaces", "a", false, "List all pods created by norden in all namespaces")

	return cmd
}

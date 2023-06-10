package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xrelkd/norden/internal/cmdutils"
)

type DeleteOptions struct {
	Namespace string
	PodName   string
}

func runDelete(opts *DeleteOptions) error {
	if opts.PodName == "" {
		return fmt.Errorf("no pod name is provided")
	}

	clientset, _, err := cmdutils.CreateClientset()
	if err != nil {
		return err
	}

	if opts.Namespace == "" {
		var ns string
		if ns, err = cmdutils.GetCurrentNamespace(); err != nil {
			return err
		} else {
			opts.Namespace = ns
		}
	}

	err = clientset.
		CoreV1().
		Pods(opts.Namespace).
		Delete(context.Background(), opts.PodName, metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	fmt.Printf("pod/%v deleted from namespace %v\n", opts.PodName, opts.Namespace)

	return nil
}

func Command() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a pod in a specified namespace",
		Long:    "Delete a pod in a specified namespace",
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "", "Namespace, use current namespace if not provided")
	cmd.Flags().StringVarP(&opts.PodName, "pod-name", "p", "", "Pod name")

	return cmd
}

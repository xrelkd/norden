package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xrelkd/norden/internal/cmdutils"
	"github.com/xrelkd/norden/internal/consts"
	"github.com/xrelkd/norden/pkg/version"
)

type CreateOptions struct {
	Namespace string
	PodName   string
	Image     string
}

func runCreate(opts *CreateOptions) error {
	clientset, _, err := cmdutils.CreateClientset()
	if err != nil {
		return err
	}

	if len(opts.Namespace) == 0 {
		if ns, err := cmdutils.GetCurrentNamespace(); err != nil {
			return err
		} else {
			opts.Namespace = ns
		}
	}

	if len(opts.Image) == 0 {
		opts.Image = consts.DefaultImage
	}

	if len(opts.PodName) == 0 {
		opts.PodName = consts.DefaultPodName
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.PodName,
			Labels:    map[string]string{"app.kubernetes.io/managed-by": version.AppName},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  consts.ContainerName,
					Image: opts.Image,
				},
			},
		},
	}
	pod, err = clientset.
		CoreV1().
		Pods(opts.Namespace).
		Create(context.Background(), pod, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	fmt.Printf("pod/%v created in namespace %v\n", pod.Name, opts.Namespace)

	return nil
}

func Command() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new pod in a specified namespace",
		Long:    "Create new pod in a specified namespace",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "", "Namespace used to create a pod, use current namespace if not provided")
	cmd.Flags().StringVarP(&opts.PodName, "pod-name", "p", consts.DefaultPodName, "Pod name")
	cmd.Flags().StringVarP(&opts.Image, "image", "i", consts.DefaultImage, "Container image")

	return cmd
}

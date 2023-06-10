package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xrelkd/norden/internal/cmdutils"
	"github.com/xrelkd/norden/internal/consts"
	"github.com/xrelkd/norden/pkg/config"
	"github.com/xrelkd/norden/pkg/version"
)

type CreateOptions struct {
	Namespace        string
	PodName          string
	Image            string
	ImagePullPolicy  v1.PullPolicy
	Command          []string
	Args             []string
	InteractiveShell []string
}

func runCreate(opts *CreateOptions) error {
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

	if opts.Image == "" {
		opts.Image = consts.DefaultImage
	}

	if opts.PodName == "" {
		opts.PodName = consts.DefaultPodName
	}

	nordenShellJSON, _ := json.Marshal(opts.InteractiveShell)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      opts.PodName,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": version.AppName,
			},
			Annotations: map[string]string{
				consts.InteractiveShellAnnotationKey: string(nordenShellJSON),
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            consts.ContainerName,
					Image:           opts.Image,
					ImagePullPolicy: opts.ImagePullPolicy,
					Command:         opts.Command,
					Args:            opts.Args,
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
	conf, err := config.Load()

	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create new pod in a specified namespace",
		Long:    "Create new pod in a specified namespace",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err != nil {
				logger.Warning("%v", err)
			}

			return runCreate(opts)
		},
	}

	confImage := conf.GetImage()

	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "", "Namespace used to create a pod, use current namespace if not provided")
	cmd.Flags().StringVarP(&opts.PodName, "pod-name", "p", conf.DefaultPodName, "Pod name")
	cmd.Flags().StringVarP(&opts.Image, "image", "i", confImage.Image, "Container image")

	imagePullPolicy := ""
	cmd.Flags().StringVar(&imagePullPolicy, "image-pull-policy", string(confImage.ImagePullPolicy), "Image pull policy")
	opts.ImagePullPolicy = v1.PullPolicy(imagePullPolicy)

	cmd.Flags().StringArrayVar(&opts.Command, "command", confImage.Command, "Command")
	cmd.Flags().StringArrayVar(&opts.Args, "args", confImage.Args, "Arguments")
	cmd.Flags().StringArrayVar(&opts.InteractiveShell, "shell", confImage.InteractiveShell, "Interactive shell")

	return cmd
}

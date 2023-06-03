package attach

import (
	"bytes"
	"context"
	"fmt"

	dockerterm "github.com/moby/term"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/util/term"

	"github.com/xrelkd/norden/internal/cmdutils"
	"github.com/xrelkd/norden/internal/consts"
)

type AttachOptions struct {
	Namespace string
	PodName   string
	Command   []string
}

func runAttach(opts *AttachOptions) error {
	_, restConfig, err := cmdutils.CreateClientset()
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

	if len(opts.PodName) == 0 {
		opts.PodName = consts.DefaultPodName
	}

	if len(opts.Command) == 0 {
		opts.Command = consts.DefaultCommand
	}

	fmt.Printf("Attach to pod/%v in namespace %v\n", opts.PodName, opts.Namespace)

	term := setupTTY()
	errOut := &bytes.Buffer{}
	tty := true
	sizeQueue := term.MonitorSize(term.GetSize())

	fn := func() error {
		restClient, err := restclient.RESTClientFor(restConfig)
		if err != nil {
			return err
		}

		req := restClient.
			Post().
			Resource("pods").
			Namespace(opts.Namespace).
			Name(opts.PodName).
			SubResource("exec")

		req.VersionedParams(&corev1.PodExecOptions{
			Container: consts.ContainerName,
			Command:   opts.Command,
			Stdin:     term.In != nil,
			Stdout:    term.Out != nil,
			Stderr:    errOut != nil,
			TTY:       tty,
		}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
		if err != nil {
			return err
		}

		return exec.StreamWithContext(context.Background(),
			remotecommand.StreamOptions{
				Stdin:             term.In,
				Stdout:            term.Out,
				Stderr:            errOut,
				Tty:               tty,
				TerminalSizeQueue: sizeQueue,
			})
	}

	return term.Safe(fn)
}

func Command() *cobra.Command {
	opts := &AttachOptions{}

	cmd := &cobra.Command{
		Use:     "attach",
		Short:   "Attach to a pod in a specific namespace",
		Long:    "Attach to a pod in a specific namespace",
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAttach(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "", "Namespace")
	cmd.Flags().StringVarP(&opts.PodName, "pod-name", "p", consts.DefaultPodName, "Pod name")
	cmd.Flags().StringArrayVarP(&opts.Command, "command", "c", consts.DefaultCommand, "Command")

	return cmd
}

func setupTTY() term.TTY {
	stdin, stdout, _ := dockerterm.StdStreams()
	term := term.TTY{
		Parent: nil,
		In:     stdin,
		Out:    stdout,
		Raw:    true,
	}

	return term
}

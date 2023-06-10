package attach

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	dockerterm "github.com/moby/term"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/util/term"

	"github.com/xrelkd/norden/internal/cmdutils"
	"github.com/xrelkd/norden/internal/consts"
	"github.com/xrelkd/norden/pkg/config"
)

type AttachOptions struct {
	Namespace       string
	PodName         string
	InterativeShell []string
}

func getShell(opts *AttachOptions) []string {
	clientset, _, err := cmdutils.CreateClientset()
	if err != nil {
		return consts.DefaultInteractiveShell
	}

	if opts.Namespace == "" {
		var ns string
		if ns, err = cmdutils.GetCurrentNamespace(); err != nil {
			return consts.DefaultInteractiveShell
		} else {
			opts.Namespace = ns
		}
	}

	pods, err := clientset.
		CoreV1().
		Pods(opts.Namespace).
		List(context.Background(), metav1.ListOptions{
			LabelSelector: consts.PodLabelSelector,
		})
	if err != nil {
		return consts.DefaultInteractiveShell
	}

	for i := range pods.Items {
		pod := &pods.Items[i]

		if pod.Name == opts.PodName {
			var shell []string
			err := json.Unmarshal([]byte(pod.Annotations[consts.InteractiveShellAnnotationKey]), &shell)
			if err != nil || len(shell) == 0 {
				return consts.DefaultInteractiveShell
			}

			return shell
		}
	}

	return consts.DefaultInteractiveShell
}

func runAttach(opts *AttachOptions) error {
	_, restConfig, err := cmdutils.CreateClientset()
	if err != nil {
		return err
	}

	if opts.Namespace == "" {
		if ns, err := cmdutils.GetCurrentNamespace(); err != nil {
			return err
		} else {
			opts.Namespace = ns
		}
	}

	if opts.PodName == "" {
		opts.PodName = consts.DefaultPodName
	}

	if len(opts.InterativeShell) == 0 {
		opts.InterativeShell = getShell(opts)
	}

	fmt.Printf(
		"Attach to pod/%s in namespace %s with shell: `%s`\n",
		opts.PodName, opts.Namespace, strings.Join(opts.InterativeShell, " "))

	terminal := setupTTY()
	errOut := &bytes.Buffer{}
	tty := true
	sizeQueue := terminal.MonitorSize(terminal.GetSize())

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
			Command:   opts.InterativeShell,
			Stdin:     terminal.In != nil,
			Stdout:    terminal.Out != nil,
			Stderr:    errOut != nil,
			TTY:       tty,
		}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
		if err != nil {
			return err
		}

		return exec.StreamWithContext(context.Background(),
			remotecommand.StreamOptions{
				Stdin:             terminal.In,
				Stdout:            terminal.Out,
				Stderr:            errOut,
				Tty:               tty,
				TerminalSizeQueue: sizeQueue,
			})
	}

	return terminal.Safe(fn)
}

func Command() *cobra.Command {
	conf, err := config.Load()

	opts := &AttachOptions{}

	cmd := &cobra.Command{
		Use:     "attach",
		Short:   "Attach to a pod in a specific namespace",
		Long:    "Attach to a pod in a specific namespace",
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err != nil {
				logger.Warning("%v", err)
			}

			return runAttach(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Namespace, "namespace", "n", "", "Namespace")
	cmd.Flags().StringVarP(&opts.PodName, "pod-name", "p", conf.DefaultPodName, "Pod name")
	cmd.Flags().StringArrayVarP(&opts.InterativeShell, "shell", "s",
		[]string{}, "Interactive shell used to attach container")

	return cmd
}

func setupTTY() term.TTY {
	stdin, stdout, _ := dockerterm.StdStreams()

	return term.TTY{
		Parent: nil,
		In:     stdin,
		Out:    stdout,
		Raw:    true,
	}
}

package consts

const (
	ContainerName                 string = "default"
	DefaultPodName                string = "norden"
	DefaultImage                  string = "ghcr.io/xrelkd/norden:latest"
	InteractiveShellAnnotationKey string = "norden/interactive-shell"
	PodLabelSelector              string = "app.kubernetes.io/managed-by=norden"
)

var DefaultInteractiveShell []string = []string{"/bin/sh"}

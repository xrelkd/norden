package consts

const (
	ContainerName  string = "default"
	DefaultPodName string = "norden"
	DefaultImage   string = "ghcr.io/xrelkd/norden:latest"

	NordenVersionAnnotationKey    string = "norden/version"
	ImageNameAnnotationKey        string = "norden/imageName"
	InteractiveShellAnnotationKey string = "norden/interactiveShell"

	PodLabelSelector string = "app.kubernetes.io/managed-by=norden"
)

var DefaultInteractiveShell []string = []string{"/bin/sh"}

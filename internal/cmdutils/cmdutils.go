package cmdutils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	GroupName          = ""
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
)

func GetCurrentNamespace() (string, error) {
	pathOptions := clientcmd.NewDefaultPathOptions()

	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return "", err
	}

	if config.CurrentContext == "" {
		err = fmt.Errorf("current-context is not set")
		return "", err
	}

	currentNamespace := config.Contexts[config.CurrentContext].Namespace

	if currentNamespace == "" {
		return "default", nil
	} else {
		return currentNamespace, nil
	}
}

func CreateClientset() (*kubernetes.Clientset, *rest.Config, error) {
	pathOptions := clientcmd.NewDefaultPathOptions()

	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return nil, nil, err
	}

	kubeConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) { return config, nil })
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		return nil, nil, err
	}

	kubeConfig.APIPath = "/api"
	kubeConfig.GroupVersion = &SchemeGroupVersion
	kubeConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, nil, err
	}

	return clientset, kubeConfig, nil
}

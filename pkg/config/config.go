package config

import (
	"os"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"

	"github.com/xrelkd/norden/internal/consts"
)

type Config struct {
	DefaultPodName string  `yaml:"defaultPodName" default:"norden"`
	DefaultImage   string  `yaml:"defaultImage" default:"norden"`
	Images         []Image `yaml:"images"`
}

type Image struct {
	Name             string        `yaml:"name"`
	Image            string        `yaml:"image"`
	ImagePullPolicy  v1.PullPolicy `yaml:"imagePullPolicy" default:"IfNotPresent"`
	Command          []string      `yaml:"command"`
	Args             []string      `yaml:"args"`
	InteractiveShell []string      `yaml:"interactiveShell"`
}

func Load() (*Config, error) {
	confFilePath, err := SearchConfigFilePath()
	if err != nil {
		return NewDefaultConfig(), err
	}

	conf := NewDefaultConfig()
	content, err := os.ReadFile(confFilePath)
	if err != nil {
		return conf, err
	}

	if err := yaml.Unmarshal(content, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func SearchConfigFilePath() (string, error) {
	return xdg.SearchConfigFile("norden/config.yaml")
}

func NewDefaultConfig() *Config {
	return &Config{
		DefaultImage: "norden",
		Images: []Image{
			*NewDefaultImage(),
		},
	}
}

func NewDefaultImage() *Image {
	return &Image{
		Name:             "norden",
		Image:            consts.DefaultImage,
		ImagePullPolicy:  v1.PullIfNotPresent,
		Command:          []string{"pause"},
		Args:             []string{},
		InteractiveShell: []string{"/bin/zsh"},
	}
}

func (c *Config) GetImage() *Image {
	for i := range c.Images {
		if c.Images[i].Name == c.DefaultImage {
			return &c.Images[i]
		}
	}

	return NewDefaultImage()
}

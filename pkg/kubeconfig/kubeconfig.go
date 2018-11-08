package config

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Read parses the file and returns the Config.
func Read(path string) (*api.Config, error) {
	return clientcmd.LoadFromFile(path)
}

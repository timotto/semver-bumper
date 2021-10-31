package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

func FromFile(filename string) (*Options, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %v: %w", filename, err)
	}

	result := Options{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	return &result, nil
}

func (o *Options) WriteToFile(filename string) error {
	var data []byte
	var err error
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".json"):
		data, err = json.Marshal(o)

	default:
		data, err = yaml.Marshal(o)
	}
	if err != nil {
		return fmt.Errorf("failed to encode configuration: %w", err)
	}

	if err := os.WriteFile(o.WriteConfig, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %v: %w", filename, err)
	}

	return nil
}

func SearchConfigFile(gitRepoPath, configFileArgument string) (string, error) {
	if configFileArgument != "" {
		return configFileArgument, nil
	}

	var projectConfigs []string
	for _, suffix := range []string{"", ".yaml", ".yml", ".json"} {
		filename := path.Join(gitRepoPath, fmt.Sprintf(".semver-bumper.conf%s", suffix))
		if !fileExists(filename) {
			continue
		}
		projectConfigs = append(projectConfigs, filename)
	}

	switch len(projectConfigs) {
	case 0:
		return "", nil
	case 1:
		return projectConfigs[0], nil
	default:
		return "", fmt.Errorf("multiple project configuration files found: %v", projectConfigs)
	}
}

func fileExists(filename string) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

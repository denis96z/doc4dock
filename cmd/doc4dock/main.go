package main

import (
	"log"
	"os"

	"io/ioutil"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	srcFileName = ".doc4dock.yaml"
)

func main() {
	curDir := getCurrentDir()
	yamlConfig := readConfigFile(curDir)

	cfg := parseYamlConfig(yamlConfig)

	tag := makeImageTag(cfg)
	callDocker(tag)
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Failed to get current directory: %s\n", err.Error())
	}
	return dir
}

func readConfigFile(dir string) string {
	path := dir + string(os.PathSeparator) + srcFileName

	yamlData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file in directory [path = %q]: %s\n", path, err.Error())
	}

	s := string(yamlData)
	log.Printf("Config:\n%s\n", s)

	return s
}

type Config struct {
	Common       EnvInfo            `yaml:"common"`
	Environments map[string]EnvInfo `yaml:"env"`
}

type EnvInfo struct {
	Registry string `yaml:"registry"`
	Group    string `yaml:"group"`
	Project  string `yaml:"project"`
	Version  string `yaml:"version"`
}

func parseYamlConfig(s string) Config {
	var c Config

	err := yaml.Unmarshal([]byte(s), &c)
	if err != nil {
		log.Fatalf("Failed to parse config: %s\n", err.Error())
	}

	return c
}

func callDocker(tag string) {
	out, err := exec.Command(
		"docker", "build", "-t", tag, ".",
	).Output()
	if err != nil {
		log.Fatalf("Failed to execute `docker` command: %s\n", err.Error())
	}
	log.Printf("`docker` command output:\n%s\n", string(out))
}

func makeImageTag(c Config) string {
	t := ""

	if c.Common.Registry != "" {
		t += c.Common.Registry + "/"
	}
	if c.Common.Group != "" {
		t += c.Common.Group + "/"
	}

	if c.Common.Project == "" {
		log.Fatal("Project name cannot be empty")
	}
	t += c.Common.Project + ":"

	if c.Common.Version == "" {
		log.Fatal("Project version cannot be empty")
	}
	t += c.Common.Version

	return t
}

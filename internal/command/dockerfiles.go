package command

import (
	"embed"
	"fmt"
)

type Dockerfile string

const (
	DockerfileGo = "go"
)

//go:embed dockerfiles/**
var dockerfilesFS embed.FS

func GetDockerfile(name Dockerfile) ([]byte, error) {
	switch name {
	case DockerfileGo:
		return dockerfilesFS.ReadFile("dockerfiles/go/Dockerfile")
	default:
		return nil, fmt.Errorf("unknown dockerfile: %s", name)
	}
}

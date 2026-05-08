package command

import (
	"context"
	"os"
	"path/filepath"

	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/testcontainers/testcontainers-go"

	_ "embed"
)

type DockerBuilder struct {
	workdir string
	env     map[string]string
	cmd     []string
	mounts  map[string]string
}

func NewDockerBuilder() *DockerBuilder {
	return &DockerBuilder{
		env:    map[string]string{},
		mounts: make(map[string]string),
	}
}

func (b *DockerBuilder) WithWorkdir(dir string) *DockerBuilder {
	b.workdir = dir
	return b
}

func (b *DockerBuilder) WithEnv(key, value string) *DockerBuilder {
	b.env[key] = value
	return b
}

func (b *DockerBuilder) WithMount(hostPath, containerPath string) *DockerBuilder {
	b.mounts[hostPath] = containerPath
	return b
}

func (b *DockerBuilder) WithCmd(cmd []string) *DockerBuilder {
	b.cmd = cmd
	return b
}

func (b *DockerBuilder) BuildImage(ctx context.Context, image string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:      image,
		WorkingDir: b.workdir,
		Env:        b.env,
		Cmd:        b.cmd,
	}

	b.applyMounts(&req)

	tc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start image container: %w", err)
	}

	return tc, nil
}

func (b *DockerBuilder) BuildContext(ctx context.Context, contextDir, dockerfile string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    contextDir,
			Dockerfile: dockerfile,
		},
		WorkingDir: b.workdir,
		Env:        b.env,
		Cmd:        b.cmd,
	}

	b.applyMounts(&req)

	tc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start build container: %w", err)
	}

	return tc, nil
}

func (b *DockerBuilder) BuildFromEmbed(
	ctx context.Context,
	name Dockerfile,
) (testcontainers.Container, error) {
	content, err := GetDockerfile(name)
	if err != nil {
		return nil, fmt.Errorf("read dockerfile: %w", err)
	}

	dir := filepath.Join(os.TempDir(), "docker-"+string(name))

	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "Dockerfile")

	if err := os.WriteFile(path, content, 0600); err != nil {
		return nil, err
	}

	return b.BuildContext(ctx, dir, "Dockerfile")
}

func (b *DockerBuilder) applyMounts(req *testcontainers.ContainerRequest) {
	if len(b.mounts) == 0 {
		return
	}

	oldModifier := req.HostConfigModifier

	req.HostConfigModifier = func(hc *container.HostConfig) {
		if oldModifier != nil {
			oldModifier(hc)
		}

		for hostPath, containerPath := range b.mounts {
			hc.Mounts = append(hc.Mounts, mount.Mount{
				Type:     mount.TypeBind,
				Source:   hostPath,
				Target:   containerPath,
				ReadOnly: false,
			})
		}
	}
}

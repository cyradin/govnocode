package executor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cyradin/govnocode/internal/docker"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

func (s *DockerSuite) TestDocker_Execute() {
	out, err := s.exec.Execute(s.ctx, []string{
		"sh", "-c", "echo hello",
	}, nil)

	s.Require().NoError(err)
	s.Require().Equal("hello\n", out.Stdout)
}

func (s *DockerSuite) TestDocker_ExecuteFail() {
	_, err := s.exec.Execute(s.ctx, []string{
		"sh", "-c", "exit 42",
	}, nil)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "exit code 42")
}

type DockerSuite struct {
	suite.Suite

	ctx  context.Context //nolint:containedctx
	c    testcontainers.Container
	exec *Docker

	hostFile string
}

func (s *DockerSuite) SetupSuite() {
	s.ctx = context.Background()

	tmpDir := s.T().TempDir()

	s.hostFile = filepath.Join(tmpDir, "hello.txt")

	err := os.WriteFile(s.hostFile, []byte("initial"), 0644)
	s.Require().NoError(err)

	c, err := docker.NewDockerBuilder().
		WithWorkdir("/app").
		WithMount(s.hostFile, "/app/hello.txt").
		WithCmd([]string{"sleep", "infinity"}).
		BuildImage(s.ctx, "alpine:3.23")

	s.Require().NoError(err)

	s.c = c
	s.exec = NewDocker(c)
}

func (s *DockerSuite) TearDownSuite() {
	if s.c != nil {
		_ = s.c.Terminate(s.ctx)
	}
}

func TestDockerSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DockerSuite))
}

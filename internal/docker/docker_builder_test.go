package docker

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type DockerSuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx
	c   testcontainers.Container

	hostFile string
}

func (s *DockerSuite) SetupSuite() {
	s.ctx = context.Background()

	tmpDir := s.T().TempDir()

	s.hostFile = filepath.Join(tmpDir, "hello.txt")

	err := os.WriteFile(s.hostFile, []byte("initial"), 0644)
	s.Require().NoError(err)

	c, err := NewDockerBuilder().
		WithWorkdir("/app").
		WithMount(s.hostFile, "/app/hello.txt").
		WithCmd([]string{"sleep", "infinity"}).
		BuildImage(s.ctx, "alpine:3.23")

	s.Require().NoError(err)

	s.c = c
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

func (s *DockerSuite) TestDockerBuilder_BuildImage() {
	exitCode, _, err := s.c.Exec(s.ctx, []string{
		"sh", "-c", "echo changed > /app/hello.txt",
	})

	s.Require().NoError(err)
	s.Require().Equal(0, exitCode)

	content, err := os.ReadFile(s.hostFile)
	s.Require().NoError(err)

	s.Require().Equal("changed\n", string(content))
}

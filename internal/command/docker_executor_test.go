package command

func (s *DockerSuite) TestDockerExecutor_Run() {
	out, err := s.exec.Run(s.ctx, []string{
		"sh", "-c", "echo hello",
	})

	s.Require().NoError(err)
	s.Require().Equal("hello\n", out.Stdout)
}

func (s *DockerSuite) TestDockerExecutor_RunFail() {
	_, err := s.exec.Run(s.ctx, []string{
		"sh", "-c", "exit 42",
	})

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "exit code 42")
}

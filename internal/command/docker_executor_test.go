package command

func (s *DockerSuite) TestDockerExecutor_Execute() {
	out, err := s.exec.Execute(s.ctx, []string{
		"sh", "-c", "echo hello",
	}, nil)

	s.Require().NoError(err)
	s.Require().Equal("hello\n", out.Stdout)
}

func (s *DockerSuite) TestDockerExecutor_ExecuteFail() {
	_, err := s.exec.Execute(s.ctx, []string{
		"sh", "-c", "exit 42",
	}, nil)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "exit code 42")
}

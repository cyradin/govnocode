package container

import "github.com/cyradin/govnocode/tools"

func (c *Container) ToolRegistry() *tools.Registry {
	if c.toolRegistry == nil {
		registry := tools.NewRegistry()

		err := registry.Register(
			tools.NewGitInit(),
			tools.NewGitCreateBranch(),
			tools.NewGitCheckout(),
			tools.NewGitStatus(),
			tools.NewGitAdd(),
			tools.NewGitCommit(),
			tools.NewGitPush(),
			tools.NewGitPull(),
			tools.NewGitFetch(),
			tools.NewGitFetch(),
			tools.NewGitBranchList(),

			tools.NewFileRead(),
			tools.NewFileWrite(),
			tools.NewFileDelete(),
			tools.NewDirList(),
			tools.NewMkdir(),
			tools.NewFileMove(),
		)
		if err != nil {
			panic(err)
		}

		c.toolRegistry = registry
	}

	return c.toolRegistry
}

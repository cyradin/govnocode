package container

import "github.com/cyradin/govnocode/tools"

func (c *Container) NewToolRegistry() *tools.Registry {
	return tools.NewRegistry()
}

package infrastructure

type Container struct {
}

func New() *Container {
	c := &Container{
	}
	return c
}

func (c *Container) Init() error {
	return nil
}

func (c *Container) Close() {
}

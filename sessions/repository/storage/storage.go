package storage

type taskStorageImpl struct {
	c *container
}

func newStorage(c *container) *taskStorageImpl {
	return &taskStorageImpl{c: c}
}


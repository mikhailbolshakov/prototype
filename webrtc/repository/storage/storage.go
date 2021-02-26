package storage

type storageImpl struct {
	c *container
}

func newStorage(c *container) *storageImpl {
	s := &storageImpl{c}
	return s
}

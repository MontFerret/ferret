package drivers

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	DriverEntry struct {
		Driver  Driver
		Options []GlobalOption
	}

	Container struct {
		drivers map[string]DriverEntry
	}
)

func NewContainer() *Container {
	return &Container{
		drivers: map[string]DriverEntry{},
	}
}

func (c *Container) Has(name string) bool {
	_, exists := c.drivers[name]

	return exists
}

func (c *Container) Register(drv Driver, opts ...GlobalOption) error {
	if drv == nil {
		return core.Error(core.ErrMissedArgument, "driver")
	}

	name := drv.Name()
	_, exists := c.drivers[name]

	if exists {
		return core.Errorf(core.ErrNotUnique, "driver: %s", name)
	}

	c.drivers[name] = DriverEntry{
		Driver:  drv,
		Options: opts,
	}

	return nil
}

func (c *Container) Remove(name string) {
	delete(c.drivers, name)
}

func (c *Container) Get(name string) (Driver, bool) {
	found, exists := c.drivers[name]

	return found.Driver, exists
}

func (c *Container) GetAll() []Driver {
	res := make([]Driver, 0, len(c.drivers))

	for _, entry := range c.drivers {
		res = append(res, entry.Driver)
	}

	return res
}

func (c *Container) WithContext(ctx context.Context) context.Context {
	next := ctx

	for _, entry := range c.drivers {
		next = withContext(next, entry.Driver, entry.Options)
	}

	return next
}

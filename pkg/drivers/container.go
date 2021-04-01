package drivers

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Container struct {
	drivers map[string]Driver
}

func NewContainer() *Container {
	return &Container{
		drivers: map[string]Driver{},
	}
}

func (c *Container) HasDriver(name string) bool {
	_, exists := c.drivers[name]

	return exists
}

func (c *Container) RegisterDriver(drv Driver) error {
	if drv == nil {
		return core.Error(core.ErrMissedArgument, "driver")
	}

	name := drv.Name()
	_, exists := c.drivers[name]

	if exists {
		return core.Errorf(core.ErrNotUnique, "driver: %s", name)
	}

	c.drivers[name] = drv

	return nil
}

func (c *Container) RegisteredDriver(name string) (Driver, bool) {
	found, exists := c.drivers[name]

	return found, exists
}

func (c *Container) RegisteredDrivers() []Driver {
	res := make([]Driver, 0, len(c.drivers))

	for _, drv := range c.drivers {
		res = append(res, drv)
	}

	return res
}

func (c *Container) RemoveDriver(name string) {
	delete(c.drivers, name)
}

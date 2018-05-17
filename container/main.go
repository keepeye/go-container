package container

import (
	"sync"
)

//sliceCap cap of release-function slice
const sliceCap = 10

// default container
var c = NewContainer()

//Definition definition of a service
type Definition struct {
	Name    interface{}
	Service func(c *Container) interface{}
	Shared  bool
}

//Container define container struct
type Container struct {
	definitions *sync.Map
	resolved    *sync.Map
	r           []func()
}

//NewContainer create a new container
func NewContainer() *Container {
	return &Container{
		definitions: new(sync.Map), // service definitions
		resolved:    new(sync.Map), // cached instances
		r:           make([]func(), 0, sliceCap), // release functions
	}
}

//Get resolve a service
func Get(name interface{}) interface{} { return c.Get(name) }
func (c *Container) Get(name interface{}) interface{} {
	if resolved, ok := c.resolved.Load(name); ok {
		return resolved
	}
	v, ok := c.definitions.Load(name)
	if !ok {
		return nil
	}
	def := v.(*Definition)
	service := def.Service(c)
	if def.Shared {
		c.resolved.Store(name, service)
	}
	return service
}

//Bind register a service, if shared is true, the service will resolve only once
func Bind(name interface{}, service func(c *Container) interface{}, shared bool) { c.Bind(name, service, shared) }
func (c *Container) Bind(name interface{}, service func(c *Container) interface{}, shared bool) {
	def := &Definition{
		Name:    name,
		Shared:  shared,
		Service: service,
	}
	c.definitions.Store(name, def)
	c.resolved.Delete(name)
}

//Singleton register a shared service
func Singleton(name interface{}, service func(c *Container) interface{}) { c.Singleton(name, service) }
func (c *Container) Singleton(name interface{}, service func(c *Container) interface{}) {
	c.Bind(name, service, true)
}

//Instance register a resolved instance
//it's actually going to be converted to a shared service
func Instance(name interface{}, instance interface{}) { c.Instance(name, instance) }
func (c *Container) Instance(name interface{}, instance interface{}) {
	service := func(c *Container) interface{} {
		return instance
	}
	c.Bind(name, service, true)
}

//Has detect if has a service
func Has(name interface{}) bool { return c.Has(name) }
func (c *Container) Has(name interface{}) bool {
	_, ok := c.definitions.Load(name)
	return ok
}

//Remove remove a service from container
func Remove(name interface{}) { c.Remove(name) }
func (c *Container) Remove(name interface{}) {
	c.definitions.Delete(name)
	c.resolved.Delete(name)
}

//BeforeRelease add a function which will be called at releasing
func BeforeRelease(f func()) { c.BeforeRelease(f) }
func (c *Container) BeforeRelease(f func()) {
	c.r = append(c.r, f)
}

//Release call all release functions
func Release() { c.Release() }
func (c *Container) Release() {
	for _, f := range c.r {
		f()
	}
}

//Refresh clear resolved service and release functions
func Refresh() { c.Refresh() }
func (c *Container) Refresh() {
	c.resolved = new(sync.Map)
	c.r = make([]func(), 0, sliceCap)
}

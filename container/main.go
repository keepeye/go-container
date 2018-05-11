package container

import (
"sync"
)

//sliceCap cap of release-function slice
const sliceCap  = 10

// default container
var c = NewContainer()

//Definition definition of a service
type Definition struct {
	Name interface{}
	Service func(c *Container) interface{}
	Shared bool
}

//NewContainer create a new container
func NewContainer() *Container {
	return &Container{
		definitions: make(map[interface{}]*Definition),
		resolved: make(map[interface{}]interface{}),
		r: make([]func(), 0, sliceCap),
	}
}

//Container define container struct
type Container struct {
	definitions map[interface{}]*Definition
	resolved map[interface{}]interface{}
	r []func()
	mutex sync.Mutex
}

//Get resolve a service
func Get(name interface{}) interface{} { return c.Get(name) }
func (c *Container) Get(name interface{}) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if resolved, ok := c.resolved[name];ok {
		return resolved
	}
	def, ok := c.definitions[name]
	if !ok {
		return nil
	}
	service := def.Service(c)
	if def.Shared {
		c.resolved[name] = service
	}
	return service
}

//Bind register a service, if shared is true, the service will resolve only once
func Bind(name interface{}, service func(c *Container) interface{}, shared bool) { c.Bind(name, service, shared) }
func (c *Container) Bind(name interface{}, service func(c *Container) interface{}, shared bool) {
	def := &Definition{
		Name: name,
		Shared: shared,
		Service: service,
	}
	c.mutex.Lock()
	c.definitions[name] = def
	delete(c.resolved, name)
	c.mutex.Unlock()
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.definitions[name]
	return ok
}

//Remove remove a service from container
func Remove(name interface{}) { c.Remove(name) }
func (c *Container) Remove(name interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.definitions, name)
	delete(c.resolved, name)
}

//BeforeRelease add a function which will be called at releasing
func BeforeRelease(f func()) { c.BeforeRelease(f) }
func (c *Container) BeforeRelease(f func()) {
	c.r = append(c.r, f)
}

//Release call all release functions
func Release() { c.Release() }
func (c *Container) Release() {
	for _,f := range c.r {
		f()
	}
}

//Refresh clear resolved service and release functions
func Refresh() { c.Refresh() }
func (c *Container) Refresh() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.resolved = make(map[interface{}]interface{})
	c.r = make([]func(), 0, sliceCap)
}
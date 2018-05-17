package container

import (
	"testing"
	"time"
)

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	if c == nil {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	c := NewContainer()
	c.definitions.Store("foo", &Definition{
		Name: "foo",
		Service: func(c *Container) interface{} {
			return "foo"
		},
		Shared: false,
	})
	v := c.Get("foo")
	if v.(string) != "foo" {
		t.Fail()
	}
}

func TestBind(t *testing.T) {
	c := NewContainer()
	c.Bind("foo", func(c *Container) interface{} {
		return time.Now().Nanosecond()
	}, false)
	c.Bind("bar", func(c *Container) interface{} {
		return time.Now().Nanosecond()
	}, true)
	v1 := c.Get("foo")
	v2 := c.Get("foo")
	if v1 == v2 {
		t.Error("v1 should not equal to v2")
	}
	v3 := c.Get("bar")
	v4 := c.Get("bar")
	if v3 != v4 {
		t.Error("v3 should equal to v4")
	}
}

func TestInstance(t *testing.T) {
	c := NewContainer()
	c.Instance("foo", "foo")
	c.Instance("bar", time.Now().UnixNano())
	v := c.Get("foo")
	if v.(string) != "foo" {
		t.Fail()
	}
	v1 := c.Get("bar")
	v2 := c.Get("bar")
	if v1 != v2 {
		t.Fail()
	}
}

func TestHas(t *testing.T) {
	c := NewContainer()
	c.Instance("foo", "foo")
	if ! c.Has("foo") {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	c := NewContainer()
	c.Instance("foo", "foo")
	c.Remove("foo")
	if c.Has("foo") {
		t.Fail()
	}
}

func TestRelease(t *testing.T) {
	c := NewContainer()
	a := 0
	c.BeforeRelease(func() {
		a++
	})
	c.BeforeRelease(func() {
		a++
	})
	c.Release()
	if a != 2 {
		t.Fail()
	}
}

func TestRefresh(t *testing.T) {
	c := NewContainer()
	c.Bind("foo", func(c *Container) interface{} {
		return time.Now().UnixNano()
	}, true)
	v1 := c.Get("foo")
	v2 := c.Get("foo")
	c.Refresh()
	v3 := c.Get("foo")
	if v1 != v2 || v2 == v3 {
		t.Fail()
	}
}

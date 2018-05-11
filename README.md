介绍(Introduce)
=====

这是一个非常简易的容器实现，主要方便嵌入到项目中，方便初始化一些服务，如：数据库、redis、缓存等。


特点(Feature)
=======

- 支持共享/非共享服务类型
- Support shared / non shared service type
- 注册服务支持工厂函数、实例对象
- You can register factory function or an instance to container
- 支持刷新容器，清空实例缓存，下次获取服务重新调用工厂函数生成
- Support refresh instance cache, then you get service, will resolve again
- 支持添加释放回调函数，当释放容器的时候触发，用于释放数据库等连接
- Support register release function so when your program end, you can release some resource such as database connection.

用法(Usage)
=======

### 安装 (Install)

```
go get github.com/keepeye/go-container
```


### 代码 (Code)

```
import (
    "github.com/keepeye/go-container/container"
    "time"
)

func main() {
    // 我们有一个默认的容器实例，直接通过 container调用所有方法
    // We have a default container instance, call all methods through container.XXX
    // 或者你可以创建一个新的container对象  通过container.NewContainer()
    // or you can make a new container, through container.NewContainer() 
    container.Bind("foo", func(c *Container) interface{}{
        return time.Now().UnixNano()
    }, true)
    
    v1 := container.Get("foo")
    v2 := container.Get("foo")
    fmt.Println(v1,v2) // v1 == v2
    
    container.Instance("bar", "12345")
    fmt.Println(container.Get("bar").(string)) // output:12345
    
    // 注册一个数据库服务
    // let's register a db service
    container.Bind("db", func(c *Container) interface{} {
        db, err := gorm.Open("mysql", "YOUR DATABASE CONNECTION DSN")
        if err != nil {
            panic(err)
        }
        // 注册释放函数
        // add release function
        c.BeforeRelease(func() {
            db.Close()
        })
        return db
    }, true)
    // 别忘了程序结束的时候调用'Release()'方法
    // Don't forget to call the'Release () method at the end of the program.
    defer container.Release()
}
```


方法列表(Method List)
=======

### Get(name interface{}) interface{}

获取一个服务

Get Resolve a service

### Bind(name interface{}, service func(c *Container) interface{}, shared bool)

绑定一个服务，如果share=true，那么这个服务是单例的，工厂函数只会执行一次

Bind Bind a service, if shared is true, the service will resolve only once


### Instance(name interface{}, instance interface{})

注册一个实例对象，实际上会被封装成一个工厂方法，注册成为一个共享服务

Instance Register a resolved instance, it's actually going to be converted to a shared service

### Has(name interface{}) bool

判断容器中是否含有某个服务

Has Detect if has a service

### func Remove(name interface{})

移除某个服务

Remove Remove a service from container

### BeforeRelease(f func())

添加一个释放时回调函数

BeforeRelease Add a function which will be called at releasing

### Release()

释放资源，触发所有释放回调函数

Release Call all release functions

### Refresh()

刷新容器实例缓存和释放函数

Refresh Clear resolved service and release functions

Dependency Injection (DI) for the GO language
=============================================

Introduction
------------

This DI-Packages enables the creation of a Service-Container by adding Struct-Tags to your Structs (Services).
In Context of this Package, Structs registered to the Container a called `Services`.
The general Workflow is to

1. Tag your Struct-Fields
2. optionaly register some Parameters (Serviceparameters) 
3. Register your Structs to the Container and make it a Service


Installation and usage
----------------------

To install it, run 

```
go get github.com/HenryVolkmer/di
```

### Struct Tags

The Tag `service:"service_alias"` Tag a Structs Field as Dependency which will be injected by the Container.
In order to inject Serviceparameters, tag your Field with `serviceparam:"param_id"`.

```golang
type Controller struct {
    // this injects a Service "my.logger"
    Logger *Logger `service:"my.logger"`
    // this injects a Parameter "dbuser"
    Logfile string `serviceparam:"logfile"`
}
```

After that, you have to add the Serviceparameter `logfile` and the Service `Logger` to the Container:

```golang
container := di.NewContainer()
// add a Serviceparameter
container.AddParameter("logfile","var/app.log")
// add Services
container.Add("my.controller",&Controller{})
container.Add("my.logger",&Logger{})
```

Now you can fetch the compiled Service with:

```golang
conn := container.Get("my.controller").(*Controller)
fmt.Sprintf(conn.Logger)
```

Example
-------

```golang
package main

import (
    "fmt"
    "net/http"
    "github.com/HenryVolkmer/di"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type MyController struct {
    Db *Connection `service:"gorm.dbconn"`
    Logger *Logger `service:"my.logger"`
}
func (this *MyController) ServeHTTP(http.ResponseWriter, *http.Request) {
    // do something with db
    db := this.Db.GetConnection()
    fmt.Sprintf("%t",db)

    // Log something
    fmt.Sprintf("%t",this.Logger)
}

type Connection struct {
    // another Dep
    Logger *Logger `service:"my.logger"`
    // some params
    User string `serviceparam:"dbuser"`
    Password string `serviceparam:"dbpass"`
    Host string `serviceparam:"dbhost"`
    DbName string `serviceparam:"dbname"`
    Port string `serviceparam:"dbport"`
    SslMode string `serviceparam:"sslmode"`
    // misc
    conn *gorm.DB
}
func (this *Connection) GetConnection() *gorm.DB {
    if this.conn == nil {
        dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",this.Host,this.User,this.Password,this.DbName,this.Port,this.SslMode)
        var err error
        this.conn,err = gorm.Open(postgres.Open(dsn), &gorm.Config{})   
        if err != nil {
            panic("Could not connect to db")
        }
    }
    return this.conn
}

type Logger struct {
    Logfile string `serviceparam:"logfile"`
    // ... some fancy logger implementation
}

func main() {
    // create di-container
    container := di.NewContainer()

    // add Serviceparameter
    container.AddParameter("sslmode","disable")
    container.AddParameter("dbhost","localhost")
    container.AddParameter("dbname","foobar")
    container.AddParameter("dbport","5432")
    container.AddParameter("logfile","var/log/foo.log")

    // add Serviceparameter from .env-File
    // you have to ensure, that env is setted properly
    container.AddParameter("dbuser","env(DB_USER)")
    container.AddParameter("dbpass","env(DB_PASS)")

    container.Add("my.controller",&MyController{})
    container.Add("my.logger",&Logger{})
    container.Add("gorm.dbconn",&Connection{})
    
    // fetch a Service
    var myController *MyController = container.Get("my.controller").(*MyController)

    // use the Service and all his Deps
    mux := http.NewServeMux()
    mux.Handle("/", myController)
    http.ListenAndServe(":8080", mux)
}
```
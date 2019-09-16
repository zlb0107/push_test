# Daenerys Service Framework
[![pipeline status](https://code.inke.cn/BackendPlatform/daenerys/badges/master/pipeline.svg)](https://code.inke.cn/BackendPlatform/daenerys/commits/master)
[![coverage report](https://code.inke.cn/BackendPlatform/daenerys/badges/master/coverage.svg)](https://code.inke.cn/BackendPlatform/daenerys/commits/master)

## 目录

- [安装](#安装)
- [代码生成工具](#代码生成工具)
- [初始化框架](#初始化框架)
- [配置](#配置)
    - [RPC服务端](#RPC服务端)
    - [RPC客户端](#RPC客户端)
    - [HTTP服务端](#HTTP服务端)
    - [HTTP客户端](#HTTP客户端)
- [编写一个RPC实例](#编写一个RPC实例)
    - [创建proto文件](#创建proto文件)
    - [安装protoc编译器和插件](#安装protoc编译器和插件)
    - [编译proto文件](#编译proto文件)
    - [RPC服务端](#RPC服务端)
    - [RPC客户端](#RPC客户端)

- [Daenerys HTTP](##Deanerys HTTP)
    - [使用前注意事项](###使用前注意事项)
    - [HTTP安装](###HTTP安装)
    - [快速入门](###快速入门)
        - [HTTP客户端](####HTTP客户端)
        - [HTTP服务端](####HTTP服务端)
    - [代码示例](###代码示例) 
    
- [Design Wiki](#design-wiki)
- [Contributing](#contributing)
- [Related projects](#related-projects)
- [Additional reading](#additional-reading)
- [Testing](#testing)

## 安装

为了安装Daenerys包，你需要安装Go，然后配置好Golang的开发环境。

1. 首先需要安装[Go](https://golang.org/dl/) (**version 1.10+ is required**)

2. 配置好GOPATH路径，GOPATH的介绍可以参考文档[GOPATH](https://github.com/golang/go/wiki/GOPATH)

3. 然后你可以使用如下命令来安装Daenerys.

```sh
$ go get -u git.inke.cn/inkelogic/daenerys
```

4. 使用的时候，把它导入到你的代码中:

```go
import "git.inke.cn/inkelogic/daenerys"
```
## 代码生成工具

代码生成工具用于创建一个新的项目，新生成的项目包含了一些通用代码。这样开发人员只需要关注业务逻辑从而提升开发效率。

具体安装与使用方法请参考[daenerys](cmd/daenerys)

建议大家在开发新项目的时候使用代码生成工具来生成框架的代码，生成的代码中会处理框架初始化、框架退出等逻辑。

## 初始化框架

在使用Daenerys框架之前，你首先需要初始化它，就像下面这样：

```go
func main() {
	daenerys.Init(daenerys.RunMode(daenerys.Development))
}
```

大多数时候你并不需要了解Init函数里面发生了什么，它会为你做很多初始化的工作。

但是为了使你更好地了解框架初始化流程，这里提供了一些细节：

Init函数主要会做以下几个事情：

1. 读取配置文件。

2. 读取与服务相关的信息（服务发现名、APP名、基础库版本等信息）, 这些信息由部署系统来提供。

3. 初始化框架内部所需要的组件（服务发现、日志、Trace等）。

4. 根据配置文件初始化rpc、http、redis、kafka等组件。

## 配置

Daenerys框架有以下几种配置项：
	1. 日志
	2. 熔断
	3. 基础组件(rpc、http、redis、kafka等)

### RPC服务端

```toml
[server]
    service_name = "a.b.c"
    port = 10000
    [server.tcp]
        idle_timeout = 3000 #ms
        keeplive_interval = 10000 #ms
```

### RPC客户端

```toml
	
[[server_client]]
        service_name="a.b.c"#remote service name
        proto="http"
        endpoints="1.1.1.1:12345,1.1.1.2:12345"
        balancetype="roundrobin"
        read_timeout=300#ms
        write_timeout=300#ms
        retry_times=0
        endpoints_from="consul"	

```
使用时可以参考[RPC](#服务端)

### HTTP服务端

```toml
[server]
    service_name="a.b.c"
    port = 10001
    
    #以下部分可不填写
    [server.http]
        location="/a/b,/aa/bb/cc"
        logResponse="true,true"

```

### HTTP客户端
```toml
[[server_client]]
        service_name="a.b.c"#remote service name
        proto="http"
        endpoints="1.1.1.1:12345,1.1.1.2:12345"
        balancetype="roundrobin"
        read_timeout=300#ms
        write_timeout=300#ms
        retry_times=0
        endpoints_from="consul"	
```

### Redis客户端

```toml
[[redis]]
        server_name="redis"
        addr="ip:port"
        password="xxxx"
        max_idle=500
        max_active=500
        idle_timeout=1000
        connect_timeout=1000
        read_timeout=3000
        write_timeout=1000
        database=0
```

## 写一个RPC服务 

下面是一个简单的echo RPC服务的例子

### 创建proto文件

RPC服务最关键的一个使用条件就是有一个强定义的接口，Daenerys使用protobuf来实现这个条件。

这里我们定义了一个有Echo方法的EchoService。

```
syntax = "proto2";

package echo;

message EchoRequest {
	required string message = 1;
};

message EchoResponse {
	required string response = 1;
};

service EchoService {
    rpc Echo (EchoRequest) returns (EchoResponse);
};
```
这个例子在[echopb](examples/echopb)


### 安装protoc编译器和插件

在我们写完proto定义之后，我们必须使用protoc和Daenerys RPC插件来编译它。

1. 你可以访问这篇文档[protocol](https://github.com/protocolbuffers/protobuf)来安装protoc编译器。

2. 你可以运行下面的Go命令来安装Daenerys RPC插件。

```sh
GOBIN=/usr/local/bin go install git.inke.cn/inkelogic/daenerys/cmd/protoc-gen-daenerys
```
当执行完上面的命令时，Daenerys RPC插件protoc-gen-daenerys将会安装在/usr/local/bin目录下面

### 编译proto文件

运行下面的命令来编译proto文件。
```sh
# assume you are in examples/rpcserver/echo directory.
protoc  --daenerys_out=plugins=rpc:./ *.proto
```

### RPC服务端

下面的代码是echo服务的服务端，强烈建议大家使用代码生成工具[daenerys](cmd/daenerys)来生成服务端的代码。


它做了下面几个事情：

1. 实现了为Echo handler定义的接口。
2. 初始化Daenerys框架。
3. 注册Echo handler。
4. 运行服务。

```sh
# assume in the following directory
$GOPATH/src/git.inke.cn/inkelogic/daenerys/examples/rpcserver
```

```go
package main

import (
	"golang.org/x/net/context"
	"git.inke.cn/inkelogic/daenerys"
	echo "git.inke.cn/inkelogic/daenerys/examples/echopb"
	proto "github.com/golang/protobuf/proto"
)

type EchoHandler struct{}

// Implements Echo handler.
func (e *EchoHandler) Echo(ctx context.Context, r *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Response: proto.String(r.GetMessage()),
	}, nil
}

func main() {
	// Initialises Daenerys. Optionally include some options here.
	daenerys.Init(
		daenerys.RunMode(daenerys.Development),
	)
	// RPCServer will initializing a rpc server from config file which default is config.toml.
	server := daenerys.RPCServer()

	// Register handler
	echo.RegisterEchoServiceHandler(server, new(EchoHandler))

	// Run the server
	if err := server.Start(); err != nil {
		panic(err)
	}
}
```
这个例子在[rpcserver](examples/rpcserver)

### RPC客户端

下面是请求echo服务端的客户端代码。
生成的proto代码中包含了echo客户端来减少重复代码。

```go
package main

import (
	"golang.org/x/net/context"
	"git.inke.cn/inkelogic/daenerys"
	echo "git.inke.cn/inkelogic/daenerys/examples/echopb"
	proto "github.com/golang/protobuf/proto"
)

func main() {
	// Initialises Daenerys. Optionally include some options here.
	daenerys.Init(daenerys.RunMode(daenerys.Development))

	// Create new rpc client instance from config.
	client := daenerys.RPCFactory("echoclient")

	// Create new echo client.
	service := echo.NewEchoService(client)

	// Call the echo server.
	response, err := service.Echo(context.TODO(), &echo.EchoRequest{
		Message: proto.String("hello"),
	})
	if err != nil {
		panic(err)
	}
	if response.GetResponse() != "hello" {
		panic(response.GetResponse())
	}
}
```
这个例子在[rpcclient](examples/rpcclient)



-------------------------------

## Deanerys HTTP

Daenerys HTTP是公司内部的http框架，为业务提供便捷的搭建HTTP服务的方式。
    
### HTTP安装

使用inkedep工具来下载daenerys框架包，inkedep工具的安装方法请见[inkedep](https://wiki.inkept.cn/display/INKE/inkedep-v2)。

```shell
inkedep get git.inke.cn/inkelogic/daenerys
```

在项目中导入http相关的依赖包：
```go
package demo

//导入http client包
import ikclient "git.inke.cn/inkelogic/daenerys/http/client"

//导入http server包
import ikserver "git.inke.cn/inkelogic/daenerys/http/server"
```

### 快速入门

#### HTTP客户端
```go
package main

import (
	"bytes"
	"fmt"
	"time"
	"git.inke.cn/inkelogic/daenerys"
	ikclient "git.inke.cn/inkelogic/daenerys/http/client"
	"golang.org/x/net/context"
)

func main() {
	//使用方式1:直接使用完整url方式访问
	quickStart()
	
	//使用方式2:使用默认http client或自定义client,并设置http request的方式
	defineClient()

	//使用方式3:通过配置文件方式获取可用的http client,并设置http request的方式
	withConfig()
}

func quickStart() {
	rsp, err := ikclient.HTTPGet("http://127.0.0.1:12345/root/user?a=b&c=d")
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Code())
	fmt.Println(rsp.String())
	fmt.Println(rsp.Bytes())
}

func defineClient() {
	//使用默认http client
	//client := ikclient.NewClient()
	
	//配置自定义client
	client := ikclient.NewClient(
		ikclient.DialTimeout(10*time.Second),
		ikclient.RequestTimeout(5*time.Second),
		ikclient.KeepAliveTimeout(30*time.Second),
	)
	//设置http request
	request := ikclient.NewRequest().WithURL("http://www.baidu.com/")
	//发起http请求,等待http响应
	resp, err := client.Call(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.String())
}

func withConfig() {
	//通过配置文件指定http client
	//参数“a.b.c”即是配置文件中[[server_client]]下的sevice_name
	httpclient := daenerys.HTTPClient("a.b.c")

	//设置http request
	body := bytes.NewBuffer([]byte(`hello`))
	req := ikclient.NewRequest()
	req.WithMethod("POST").
		WithCtxInfo(context.Background()).
		WithURL("http://localhost:12345/root/user").
		WithBody(body)

	//http call
	rsp, err := httpclient.Call(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.Code())
	fmt.Println(rsp.String())
	fmt.Println(rsp.Bytes())

	//save rsp body
	err = rsp.Save("/tmp/a.txt")
	if err != nil {
		panic(err)
	}
}

```
这个例子在:[httpclient](examples/httpclient)

#### HTTP服务端

```go
package main

import (
	"git.inke.cn/inkelogic/daenerys"
	ikserver "git.inke.cn/inkelogic/daenerys/http/server"
	"time"
)

func main() {
  //通过配置文件获取http server
	go quickStart()

  //自定义http server
	go definedServer()

	//blocked
	select {}

}

func quickStart() {
	//使用默认http server
	s := daenerys.HTTPServer()

	defer s.Stop()

	//http路由注册
	s.GETPOST("/hello", func(c *ikserver.Context) {
		data := map[string]string{"a": "b"}
    //以json方式响应结果
		c.JSON(data, nil)
	})

	//start server
	if err := s.Run(":12345"); err != nil {
		panic(err)
	}
}

func definedServer() {
	//use defined config init server
	// name should be set
	s := ikserver.NewServer(
		//指定name选项, 用于接入consul服务发现
		ikserver.Name("a.b.c"),
		ikserver.Port(12345),
		ikserver.WriteTimeout(6*time.Second),
		ikserver.ReadTimeout(3*time.Second),
		ikserver.IdleTimeout(15*time.Second),
	)

	s.GET("/login", func(c *ikserver.Context) {
		data := map[string]string{"login": "hello world"}
		c.JSON(data, nil)
	})

	defer s.Stop()

	//start server
	if err := s.Run(); err != nil {
		panic(err)
	}
}

```
这个例子在:[httpserver](examples/httpserver)

### 代码示例

### HTTP Client

#### 1.自定义client
```go
package main

import (
	"bytes"
	"fmt"
	"time"
	"git.inke.cn/inkelogic/daenerys"
	"git.inke.cn/inkelogic/daenerys/log"
	ikclient "git.inke.cn/inkelogic/daenerys/http/client"
)

func main(){
	client := ikclient.NewClient(
		ikclient.RetryTimes(1),
		ikclient.DialTimeout(30*time.Second),
		ikclient.IdleConnTimeout(10*time.Second),
		ikclient.KeepAliveTimeout(10*time.Second),
		ikclient.MaxIdleConns(100),
		ikclient.MaxIdleConnsPerHost(10),
		ikclient.RequestTimeout(10*time.Second),
	)
	req := ikclient.NewRequest()
	rsp, err := client.Call(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}
```

##### 1.1client option
通过option配置一个http client，option相关接口如下：
```go
//配置建连超时时间
func DialTimeout(d time.Duration) Option
//配置通用的request超时时间
func RequestTimeout(d time.Duration) Option
//配置每个连接的空闲超时时间
func IdleConnTimeout(d time.Duration) Option
//配置keepalive超时时间
func KeepAliveTimeout(d time.Duration) Option
//配置是否开启keepalive
func KeepAlivesDisable(t bool) Option
//配置重试次数
func RetryTimes(d int) Option
//配置与单台server最大的连接数
func MaxIdleConnsPerHost(d int) Option
//配置整体最大连接数
func MaxIdleConns(d int) Option
//使用自定义http client
func WithClient(client *http.Client) Option
```


#### 2.自定义request

```go
package main

import (
	"bytes"
	"fmt"
	"time"
	"git.inke.cn/inkelogic/daenerys"
	"git.inke.cn/inkelogic/daenerys/log"
	ikclient "git.inke.cn/inkelogic/daenerys/http/client"
)

func main(){
	client := ikclient.NewClient()
	val := "xxx"
	req := ikclient.NewRequest().
		//path: a=xxx
		WithURL(fmt.Sprintf("/v1?a=%s", val)).
		WithScheme("").
		//will override WithURL
		WithPath("/:name").
		WithPathParams(map[string]string{"name":"jake"}).
		WithBody(nil).
		WithMethod("POST").
		//effect on replace path: a=bbb
		WithQueryParam("a", "bbb")
	
	rsp, err := client.Call(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rsp)
}
```

##### 2.1 Request提供的方法
```go
//使用指定的http request
func (r *Request) WithRequest(req *http.Request) *Request
//绑定context
func (r *Request) WithCtxInfo(ctx context.Context) *Request
//设置method
func (r *Request) WithMethod(name string) *Request
//设置完整url
func (r *Request) WithURL(uri string) *Request
//设置scheme
func (r *Request) WithScheme(scheme string) *Request
//设置请求的path
func (r *Request) WithPath(path string) *Request
//提供请求path中带参数解析
//例如：/v1/users/:userId/details
// r.WithPathParams(map[string]string{"userId": "12345"})
func (r *Request) WithPathParams(params map[string]string) *Request
//设置请求的body
func (r *Request) WithBody(body io.Reader) *Request
//设置header
func (r *Request) AddHeader(key, value string) *Request
func (r *Request) DelHeader(key string) *Request
func (r *Request) WithMultiHeader(headers map[string]string) *Request
//设置url中的参数：
//例如：/v1/user?a=b&c=d；r.WithQueryParam(a, b).WithQueryParam(c,d)
func (r *Request) WithQueryParam(param, value string) *Request
func (r *Request) WithMultiQueryParam(params map[string]string) *Request
//设置cookie
func (r *Request) WithCookie(ck *http.Cookie) *Request
func (r *Request) WithMultiCookie(cks []*http.Cookie) *Request
//获取原始http request
func (r *Request) RawRequest() *http.Request
```


#### 3.response使用
ikclient包对http response进行了封装，提供了更为简便的使用方式。response提供的接口如下。
```go
//http status code
func (r *Response) Code() int
//http请求响应的错误信息
func (r *Response) Error() error
//以byte方式读取http响应的body
func (r *Response) Bytes() []byte
//以string方式读取http响应的body
func (r *Response) String() string
//以json方式读取http响应的body
func (r *Response) JSON(obj interface{}) error
//将http响应的body保存到文件中
func (r *Response) Save(fileName string) error
//原始的http请求
func (r *Response) RawRequest() *http.Request
//原始的http响应
func (r *Response) RawResponse() *http.Response
```


### HTTP Server

#### 1.使用 GET, POST, PUT, PATCH, DELETE, OPTIONS
```go
package main

import (
	"git.inke.cn/inkelogic/daenerys"
	ikserver "git.inke.cn/inkelogic/daenerys/http/server"
)

func main() {
	foo()
}

func foo() {
		router := ikserver.NewServer()
		router.GET("/someGet", getting)
		router.POST("/somePost", posting)
		router.PUT("/somePut", putting)
		router.DELETE("/someDelete", deleting)
		router.PATCH("/somePatch", patching)
		router.HEAD("/someHead", head)
		router.OPTIONS("/someOptions", options)
		err  := router.Run()
		if err != nil{
			panic(err)
		}
}

func getting(c *ikserver.Context) {
}

func posting(c *ikserver.Context) {
}

func putting(c *ikserver.Context) {
}

func deleting(c *ikserver.Context) {
}

func patching(c *ikserver.Context) {
}

func head(c *ikserver.Context) {
}

func options(c *ikserver.Context) {
}

```

#### 2.获取路径中的参数

```go
package main

import (
	"git.inke.cn/inkelogic/daenerys"
	ikserver "git.inke.cn/inkelogic/daenerys/http/server"
)

func main() {
	boo()
}

func boo() {
	s := ikserver.NewServer(ikserver.Name("a.b.c"))
	// 此规则能够匹配/user/john这种格式，但不能匹配/user/ 或 /user这种格式
	s.GET("/user/:name", func(c *ikserver.Context) {
		name, _ := c.Params.Get("name")
		c.Response.WriteString(name)
	})

	// 但是，这个规则既能匹配/user/john/格式也能匹配/user/john/send这种格式
	s.GET("/user/:name/*action", func(c *ikserver.Context) {
		name, _ := c.Params.Get("name")
		action, _ := c.Params.Get("action")
		message := name + " is " + action
		c.Response.WriteString(message)
	})

	err := s.Run(":8080")
	if err != nil {
		panic(err)
	}
}
```

#### 3.路由分组

```go
package main

import (
	"git.inke.cn/inkelogic/daenerys"
	ikserver "git.inke.cn/inkelogic/daenerys/http/server"
)

func main() {
	coo()
}

func coo() {
	s := ikserver.NewServer()

	// Simple group: v1
	v1 := s.GROUP("/v1")
	{
		//match path: /v1/login
		v1.POST("/login", loginHandler1)
	}

	// Simple group: v2
	v2 := s.GROUP("/v2")
	{
		//match path: /v2/user
		v2.POST("/user", loginHandler2)
	}

	s.Run(":8080")
}


func loginHandler1(c *ikserver.Context) {
}

func loginHandler2(c *ikserver.Context) {
}

```
#### 4.使用自定义插件
```go
package main

import (
	"git.inke.cn/inkelogic/daenerys"
	ikserver "git.inke.cn/inkelogic/daenerys/http/server"
)

func main() {
	eoo()
}


func eoo() {
	s := ikserver.NewServer(ikserver.Name("a.b.c"))

	// 全局中间件
	s.Use(func(c *ikserver.Context) {
		//do some logic
	})
	
	// 路由添加中间件，可以添加任意多个
	s.GET("/benchmark", func(c *ikserver.Context) {
		//do sth
	}, func(c *ikserver.Context) {
		//do sth
	})

	// 路由组中添加中间件
	authorized := s.GROUP("/", groupHandler)
	{
		authorized.POST("/login", func(c *ikserver.Context) {})
		authorized.POST("/submit", func(c *ikserver.Context) {})
	}

	// Listen and serve on 0.0.0.0:8080
	s.Run(":8080")
}

func groupHandler(c *ikserver.Context) {
}

```





## Design Wiki

[Daenerys Design Docs](http://wiki.inkept.cn/pages/viewpage.action?pageId=50836818)

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Added some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Related projects

Projects with a ★ have had particular influence on Go kit's design

- [go-kit](https://github.com/go-kit/kit) ★
- [Gin](https://gin-gonic.github.io/gin/) ★
- [go-micro](https://github.com/myodc/go-micro), a microservices client/server library ★
- [Martini](https://github.com/go-martini/martini)

## Additional reading

- [Your Server as a Function](http://monkey.org/~marius/funsrv.pdf) (PDF) — Twitter
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) (HTML)

---

Development supported by [INF](http://wiki.inkept.cn/pages/viewpage.action?pageId=1639983).

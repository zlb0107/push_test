# 使用说明
## 安装
- 编译安装标准C++实现的protobuf， 参考  https://developers.google.com/protocol-buffers/

- 安装golang环境，参考https://golang.org/doc/install

- 安装golang protobuf，可以使用```go get -u github.com/golang/protobuf/protoc-gen-go```

- 克隆rpc-go工程
  ```shell
  mkdir -p ${GOPATH}/src/git.inke.cn/inkelogic;
  cd ${GOPATH}/src/git.inke.cn/inkelogic
  git clone git@code.inke.cn:BackendPlatform/rpc-go.git
  ```
  或者使用go get命令直接安装
  ```
  go get git.inke.cn/inkelogic/rpc-go
  ```

- 安装rpc代码生成工具

  ```sheel
  go install git.inke.cn/inkelogic/rpc-go/tools/protoc-gen-inkerpc
  ```

## 使用

- 通过proto文件，生成rpc-go的interface代码， 命令如下

  ```shell
   protoc  --inkerpc_out=plugins=rpc:./ *.proto
  ```

## Example

 demo代码参见 https://code.inke.cn/BackendPlatform/rpc-go/demo

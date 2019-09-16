package main

import (
	"flag"
	"git.inke.cn/inkelogic/rpc-go"
	logs "github.com/cihub/seelog"
	"net/http"
	"os"
	"pu/controllers"
)

type Config struct {
}

var configPath string

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	flag.StringVar(&configPath, "config", "./conf/config.toml", "rpc config file")
	flag.Parse()

	var AppConfig Config
	rpcConfig, err := rpc.NewConfigToml(configPath, &AppConfig)
	if err != nil {
		logs.Error("rpc.NewConfigToml error: ", err)
		logs.Flush()
		os.Exit(2)
	}

	server, err := InitServer(rpcConfig)
	if err != nil {
		logs.Error("server.InitServer error: ", err)
		logs.Flush()
		os.Exit(2)
	}

	err = server.Serve(rpcConfig.Port())
	if err != nil {
		logs.Error("server.Serve error: ", err.Error())
		logs.Flush()
		os.Exit(2)
	}
}

func InitServer(conf rpc.Config) (*rpc.HTTPServer, error) {
	server := rpc.NewHTTPServerWithConfig(conf)
	server.Register(&controllers.PingHandler{}, &controllers.LogicController{})

	logs.Info("init server:", conf.GetServiceName())

	return server, nil
}

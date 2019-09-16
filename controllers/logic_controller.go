package controllers

import (
	"context"
	"fmt"
	"git.inke.cn/inkelogic/rpc-go"
	"go_common_lib/controllers"
	"net/http"
)

type LogicController struct {
	controllers.LogicController
}

func (this *LogicController) Serve(ctx context.Context, req *http.Request) (interface{}, int) {
	serviceName := "rec.server.push"
	uri := "/push?uid=727248455&typeid=E0"
	result, _ := rpc.HttpGet(context.TODO(), serviceName, uri, nil)
	fmt.Println(string(result))
	return nil, 0
}

func (this *LogicController) Text() bool {
	return true
}

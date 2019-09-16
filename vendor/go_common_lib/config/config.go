package config

import (
	"os"
	"path/filepath"

	"github.com/astaxie/beego/config"
)

var AppConfig config.Configer

func init() {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var filename = "app.conf"
	appConfigPath := filepath.Join(workPath, "conf", filename)

	ac, err := config.NewConfig("ini", appConfigPath)
	if err != nil {
		panic(err)
	}

	AppConfig = ac
}

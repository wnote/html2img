package conf

import (
	"os"
)

var GConf = map[string]string{
	"font_path": "",
}

func init() {
	goPath, exist := os.LookupEnv("GOPATH")
	if !exist {
		panic("GOPATH error")
	}
	GConf["font_path"] = goPath + "/src/github.com/wnote/html2img/conf/fonts/"
}

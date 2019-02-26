package conf

import (
	"os"
)

var GConf = map[string]string{
	"font_path": "",
}

var DPI = float64(72)

func init() {
	goPath, exist := os.LookupEnv("GOPATH")
	if !exist {
		panic("GOPATH error")
	}
	GConf["font_path"] = goPath + "/src/github.com/wnote/html2img/conf/fonts/"
}

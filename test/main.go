package main

import (
	"fmt"
	"os"

	"github.com/wnote/html2img"
	"github.com/wnote/html2img/conf"
)

func main() {
	conf.GConf["font_path"] = "../conf/fonts/"

	imgByte, err := html2img.ParseHtml("./example.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	fp, err := os.Create("./generated.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fp.Chmod(0777)
	fp.Write(imgByte)
	fp.Close()
}

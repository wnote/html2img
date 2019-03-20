package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/wnote/html2img"
)

func main() {

	htmlUrl := "https://www.baidu.com/"
	rsp, err := http.Get(htmlUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	htmlBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp.Body.Close()
	fmt.Println(htmlBytes)
	imgByte, err := html2img.Html2Img(htmlBytes, 750)
	if err != nil {
		fmt.Println(err)
		return
	}
	fh, err := os.Create("./generated.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fh.Chmod(0755)
	fh.Write(imgByte)
	fh.Close()
}

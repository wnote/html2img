package main

import (
	"fmt"
	"os"

	"github.com/wnote/html2img"
)

func main() {
	imgByte, err := html2img.Html2Img("./example.html")
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

	/*jsonStr, err := json.MarshalIndent(parsedBodyDom, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fp, err := os.Create("./example.json")
	fp.Chmod(0755)
	fp.Write(jsonStr)
	fp.Close()*/
}

package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/wnote/html2img"
	"github.com/wnote/html2img/dom"
	"golang.org/x/net/html"
)

func OutputImg() {
	htmlPath := "./example.html"
	imgByte, err := html2img.Html2Img(htmlPath)
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

// For test
func ExportJson(htmlPath string) {
	htmlBytes, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		log.Fatal(err)
	}
	htmlIoReader := bytes.NewReader(htmlBytes)
	htmlNode, err := html.Parse(htmlIoReader)
	if err != nil {
		log.Fatal(err)
	}

	body, styleList := dom.GetBodyStyle(htmlNode)

	var styleString []string
	for _, value := range styleList {
		styleString = append(styleString, value.FirstChild.Data)
	}
	tagStyleList := dom.ParseStyle(styleString)

	parsedBodyDom := dom.GetHtmlDom(body, tagStyleList)

	jsonStr, err := json.MarshalIndent(parsedBodyDom, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fp, err := os.Create("./example.json")
	fp.Chmod(0755)
	fp.Write(jsonStr)
	fp.Close()
}

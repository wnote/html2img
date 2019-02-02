package html2img

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/wnote/html2img/dom"
	"github.com/wnote/html2img/img"
	"golang.org/x/net/html"
)

func ParseHtml(htmlPath string) ([]byte, error) {
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

	parsedBodyDom := dom.GetHtmlDom(body, tagStyleList, nil)

	/* jsonStr, err := json.Marshal(parsedBodyDom)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonStr)) */

	return img.BodyDom2Img(parsedBodyDom)
}

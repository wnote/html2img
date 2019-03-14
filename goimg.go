package html2img

import (
	"bytes"
	"io/ioutil"
	"log"

	"golang.org/x/net/html"
)

func Html2Img(htmlPath string) ([]byte, error) {
	htmlBytes, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		log.Fatal(err)
	}
	htmlIoReader := bytes.NewReader(htmlBytes)
	htmlNode, err := html.Parse(htmlIoReader)
	if err != nil {
		log.Fatal(err)
	}

	body, styleList := GetBodyStyle(htmlNode)

	var styleString []string
	for _, value := range styleList {
		styleString = append(styleString, value.FirstChild.Data)
	}
	tagStyleList := ParseStyle(styleString)

	parsedBodyDom := GetHtmlDom(body, tagStyleList)

	return bodyDom2Img(parsedBodyDom)
}

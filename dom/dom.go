package dom

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	DOM_TYPE_TEXT        = 1
	DOM_TYPE_ELEMENT     = 3
	DOM_TYPE_COMMENTNODE = 4

	DOM_TYPE_IMAGE = 99
)

type ImageData struct {
	Fm  string
	Img image.Image
}

type Dom struct {
	CalcOffsetX int
	CalcOffsetY int
	CalcWidth   int
	CalcHeight  int
	DomType     int8
	TagName     string
	TagClass    string
	TagData     interface{}

	TagStyle *TagStyle

	Parents  []*Dom
	Children []*Dom
}

func (d *Dom) IsPositionAbsolute() bool {
	return d.TagStyle.Position == "absolute"
}

func (d *Dom) IsAutoHeight() bool {
	return d.TagStyle.Height == "auto" || d.TagStyle.Height == ""
}

func GetHtmlDom(htmlNode *html.Node, tagStyleList []*TagStyle, parents []*Dom) *Dom {
	bodyDom := &Dom{}
	SetDomAttr(bodyDom, htmlNode)
	domStyle := GetDomStyle(bodyDom, tagStyleList)
	bodyDom.TagStyle = domStyle
	bodyDom.Children = GetChildren(htmlNode, tagStyleList, []*Dom{bodyDom})
	return bodyDom
}

func GetChildren(htmlNode *html.Node, tagStyleList []*TagStyle, parents []*Dom) []*Dom {
	var children []*Dom
	for ch := htmlNode.FirstChild; ch != nil; {
		if ch.Type != html.ElementNode && ch.Type != html.TextNode {
			ch = ch.NextSibling
			continue
		}
		// ignore empty text node
		if ch.Type == html.TextNode {
			textData := strings.Trim(ch.Data, CUT_SET_LIST)
			if textData == "" {
				ch = ch.NextSibling
				continue
			}
		}
		dom := &Dom{}
		SetDomAttr(dom, ch)
		domStyle := GetDomStyle(dom, tagStyleList)
		dom.TagStyle = domStyle
		switch ch.Data {
		case "img":
			src := GetAttr(ch, "src")
			resp, err := http.Get(src)
			if err != nil {
				panic(fmt.Sprintf("http.GetImage err :%v", err))
			}
			img, fm, err := image.Decode(resp.Body)
			resp.Body.Close()
			if err != nil {
				panic(fmt.Sprintf("image.Decode err :%v", err))
			}
			imgData := ImageData{
				Fm:  fm,
				Img: img,
			}
			dom.TagData = imgData
		default:

		}
		par := append(parents, dom)
		child := GetChildren(ch, tagStyleList, par)
		dom.Children = child
		children = append(children, dom)

		ch = ch.NextSibling
	}
	return children
}

func GetDomStyle(dom *Dom, tagStyleList []*TagStyle) *TagStyle {
	for _, style := range tagStyleList {
		if style.Selected(dom) {
			return style
		}
	}
	return nil
}

func SetDomAttr(dom *Dom, htmlNode *html.Node) {
	dom.DomType = int8(htmlNode.Type)
	if htmlNode.Type == html.ElementNode {
		dom.DomType = DOM_TYPE_ELEMENT
		dom.TagName = htmlNode.Data
		dom.TagClass = GetAttr(htmlNode, "class")
	} else if htmlNode.Type == html.TextNode {
		dom.DomType = DOM_TYPE_TEXT
		dom.TagData = htmlNode.Data
	}
}

func GetAttr(htmlNode *html.Node, attrKey string) string {
	for _, attr := range htmlNode.Attr {
		if attr.Key == attrKey {
			return attr.Val
		}
	}
	return ""
}

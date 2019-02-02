package dom

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/wnote/html2img/conf"
	"golang.org/x/net/html"
)

const CUT_SET_LIST = "\n\t\b "

var FontMapping = make(map[string]*truetype.Font)

type TagStyle struct {
	// Style selector
	Selector string

	// Inheritable
	Color      string
	FontSize   string
	LineHeight string
	FontFamily string

	// Not Inheritable
	BackgroundColor string
	BackgroundImage string
	Width           string
	Height          string
	Left            string
	Top             string
	Bottom          string
	Right           string
	MarginLeft      string
	MarginTop       string
	MarginRight     string
	MarginBottom    string
	PaddingLeft     string
	PaddingRight    string
	PaddingTop      string
	PaddingBottom   string
	Display         string
	BorderRadius    string
	Position        string
}

// 暂时不考虑优先级问题
func ParseStyle(styleList []string) []*TagStyle {
	tagStyleMap := make(map[string]*TagStyle)
	for _, style := range styleList {
		subList := strings.Split(style, "}")
		for _, subTag := range subList {
			tag := strings.Split(subTag, "{")
			if len(tag) > 1 {
				selector := strings.Trim(tag[0], CUT_SET_LIST)
				re := regexp.MustCompile("/\\s+/")
				selector = re.ReplaceAllString(selector, " ")
				classStyle := strings.Trim(tag[1], CUT_SET_LIST)
				classStyleList := strings.Split(classStyle, ";")
				tagStyle := &TagStyle{}
				if oldStyle, exist := tagStyleMap[selector]; exist {
					tagStyle = oldStyle
				}
				for _, cStyle := range classStyleList {
					SetTagStyle(tagStyle, cStyle)
				}
				tagStyleMap[selector] = tagStyle
			}
		}
	}

	var tagStyleList []*TagStyle
	for selector, tagStyle := range tagStyleMap {
		tagStyle.Selector = strings.Trim(selector, " ")
		tagStyleList = append(tagStyleList, tagStyle)
	}
	return tagStyleList
}

func SetTagStyle(tagStyle *TagStyle, cStyle string) {
	cStyle = strings.Trim(cStyle, CUT_SET_LIST)
	if cStyle == "" {
		return
	}
	css := strings.Split(cStyle, ":")
	if len(css) != 2 {
		panic(fmt.Sprintf("unsupported style %v", cStyle))
	}
	cssKey := strings.Trim(css[0], CUT_SET_LIST)
	cssValue := strings.Trim(css[1], CUT_SET_LIST)
	if cssValue == "" {
		return
	}

	switch cssKey {
	case "background-color":
		tagStyle.BackgroundColor = cssValue
	case "background-image":
		tagStyle.BackgroundImage = cssValue
	case "width":
		tagStyle.Width = cssValue
	case "height":
		tagStyle.Height = cssValue
	case "color":
		tagStyle.Color = cssValue
	case "font-size":
		tagStyle.FontSize = cssValue
	case "left":
		tagStyle.Left = cssValue
	case "top":
		tagStyle.Top = cssValue
	case "bottom":
		tagStyle.Bottom = cssValue
	case "right":
		tagStyle.Right = cssValue
	case "margin-left":
		tagStyle.MarginLeft = cssValue
	case "margin-top":
		tagStyle.MarginTop = cssValue
	case "margin-right":
		tagStyle.MarginRight = cssValue
	case "margin-bottom":
		tagStyle.MarginBottom = cssValue
	case "padding-left":
		tagStyle.PaddingLeft = cssValue
	case "padding-right":
		tagStyle.PaddingRight = cssValue
	case "padding-top":
		tagStyle.PaddingTop = cssValue
	case "padding-bottom":
		tagStyle.PaddingBottom = cssValue
	case "display":
		tagStyle.Display = cssValue
	case "border-radius":
		tagStyle.BorderRadius = cssValue
	case "line-height":
		tagStyle.LineHeight = cssValue
	case "font-family":
		cssValue = strings.Trim(cssValue, "'\"")
		if cssValue == "" {
			return
		}
		initFontMap(cssValue)
		tagStyle.FontFamily = cssValue
	case "position":
		tagStyle.Position = cssValue
	case "padding":
		attrList := strings.Split(cssValue, " ")
		switch len(attrList) {
		case 1:
			tagStyle.PaddingTop = attrList[0]
			tagStyle.PaddingBottom = attrList[0]
			tagStyle.PaddingLeft = attrList[0]
			tagStyle.PaddingRight = attrList[0]
		case 2:
			tagStyle.PaddingTop = attrList[0]
			tagStyle.PaddingBottom = attrList[0]
			tagStyle.PaddingLeft = attrList[1]
			tagStyle.PaddingRight = attrList[1]
		case 3:
			tagStyle.PaddingTop = attrList[0]
			tagStyle.PaddingLeft = attrList[1]
			tagStyle.PaddingRight = attrList[1]
			tagStyle.PaddingBottom = attrList[2]
		case 4:
			tagStyle.PaddingTop = attrList[0]
			tagStyle.PaddingRight = attrList[1]
			tagStyle.PaddingBottom = attrList[2]
			tagStyle.PaddingLeft = attrList[3]
		default:
			panic(fmt.Sprintf("unsupported padding value %v", cStyle))
		}
	case "margin":
		attrList := strings.Split(cssValue, " ")
		switch len(attrList) {
		case 1:
			tagStyle.MarginTop = attrList[0]
			tagStyle.MarginBottom = attrList[0]
			tagStyle.MarginLeft = attrList[0]
			tagStyle.MarginRight = attrList[0]
		case 2:
			tagStyle.MarginTop = attrList[0]
			tagStyle.MarginBottom = attrList[0]
			tagStyle.MarginLeft = attrList[1]
			tagStyle.MarginRight = attrList[1]
		case 3:
			tagStyle.MarginTop = attrList[0]
			tagStyle.MarginLeft = attrList[1]
			tagStyle.MarginRight = attrList[1]
			tagStyle.MarginBottom = attrList[2]
		case 4:
			tagStyle.MarginTop = attrList[0]
			tagStyle.MarginRight = attrList[1]
			tagStyle.MarginBottom = attrList[2]
			tagStyle.MarginLeft = attrList[3]
		default:
			panic(fmt.Sprintf("unsupported margin value %v", cStyle))
		}
	default:
		panic(fmt.Sprintf("unsupported %v", cStyle))

	}
}

func GetBodyStyle(htmlNode *html.Node) (body *html.Node, styleList []*html.Node) {
	for ch := htmlNode.FirstChild; ch != nil; {
		switch ch.Data {
		case "body":
			body = ch
			_, tmpStyle := GetBodyStyle(ch)
			if len(tmpStyle) > 0 {
				styleList = append(styleList, tmpStyle...)
			}
		case "style":
			styleList = append(styleList, ch)
		default:
			tmpBody, tmpStyle := GetBodyStyle(ch)
			if tmpBody != nil {
				body = tmpBody
			}
			if len(tmpStyle) > 0 {
				styleList = append(styleList, tmpStyle...)
			}
		}
		ch = ch.NextSibling
	}
	return
}

// 暂时不考虑多个选择器问题
func (p *TagStyle) Selected(currentDom *Dom) bool {
	selectors := strings.Split(p.Selector, " ")
	// ignore parents
	if len(selectors) > 1 {
		panic(fmt.Sprintf("multiple selector will be supported later  %v", p.Selector))
	}

	subSels := strings.Split(selectors[0], ".")
	if subSels[0] != "" {
		if subSels[0] != currentDom.TagName {
			return false
		}
	}
	classList := subSels[1:]
	if len(classList) == 0 {
		return true
	}
	re := regexp.MustCompile("/\\s+/")
	domClasses := re.Split(currentDom.TagClass, -1)
	domClassMap := make(map[string]bool)
	for _, class := range domClasses {
		if class != "" {
			domClassMap[class] = true
		}
	}
	matched := true
	for _, class := range classList {
		if _, exist := domClassMap[class]; !exist {
			matched = false
		}
	}
	return matched
}

func initFontMap(fontFamily string) {
	if _, exist := FontMapping[fontFamily]; exist {
		return
	}
	fontPath := conf.GConf["font_path"]
	f, err := getFontFromFile(fontPath + "/" + fontFamily)
	if err != nil {
		panic(err)
	}
	FontMapping[fontFamily] = f
}

func getFontFromFile(fontfile string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	return f, nil
}

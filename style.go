package html2img

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

type Pos struct {
	Left   string
	Top    string
	Bottom string
	Right  string
}

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
	Display         string
	Position        string

	BorderRadius Pos
	Offset       Pos
	Margin       Pos
	Padding      Pos
	BorderWidth  Pos
	BorderColor  Pos
	BorderStyle  Pos
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
		tagStyle.Offset.Left = cssValue
	case "top":
		tagStyle.Offset.Top = cssValue
	case "bottom":
		tagStyle.Offset.Bottom = cssValue
	case "right":
		tagStyle.Offset.Right = cssValue
	case "margin-left":
		tagStyle.Margin.Left = cssValue
	case "margin-top":
		tagStyle.Margin.Top = cssValue
	case "margin-right":
		tagStyle.Margin.Right = cssValue
	case "margin-bottom":
		tagStyle.Margin.Bottom = cssValue
	case "padding-left":
		tagStyle.Padding.Left = cssValue
	case "padding-right":
		tagStyle.Padding.Right = cssValue
	case "padding-top":
		tagStyle.Padding.Top = cssValue
	case "padding-bottom":
		tagStyle.Padding.Bottom = cssValue
	case "display":
		tagStyle.Display = cssValue
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
			tagStyle.Padding.Top = attrList[0]
			tagStyle.Padding.Bottom = attrList[0]
			tagStyle.Padding.Left = attrList[0]
			tagStyle.Padding.Right = attrList[0]
		case 2:
			tagStyle.Padding.Top = attrList[0]
			tagStyle.Padding.Bottom = attrList[0]
			tagStyle.Padding.Left = attrList[1]
			tagStyle.Padding.Right = attrList[1]
		case 3:
			tagStyle.Padding.Top = attrList[0]
			tagStyle.Padding.Left = attrList[1]
			tagStyle.Padding.Right = attrList[1]
			tagStyle.Padding.Bottom = attrList[2]
		case 4:
			tagStyle.Padding.Top = attrList[0]
			tagStyle.Padding.Right = attrList[1]
			tagStyle.Padding.Bottom = attrList[2]
			tagStyle.Padding.Left = attrList[3]
		default:
			panic(fmt.Sprintf("unsupported padding value %v", cStyle))
		}
	case "margin":
		attrList := strings.Split(cssValue, " ")
		switch len(attrList) {
		case 1:
			tagStyle.Margin.Top = attrList[0]
			tagStyle.Margin.Bottom = attrList[0]
			tagStyle.Margin.Left = attrList[0]
			tagStyle.Margin.Right = attrList[0]
		case 2:
			tagStyle.Margin.Top = attrList[0]
			tagStyle.Margin.Bottom = attrList[0]
			tagStyle.Margin.Left = attrList[1]
			tagStyle.Margin.Right = attrList[1]
		case 3:
			tagStyle.Margin.Top = attrList[0]
			tagStyle.Margin.Left = attrList[1]
			tagStyle.Margin.Right = attrList[1]
			tagStyle.Margin.Bottom = attrList[2]
		case 4:
			tagStyle.Margin.Top = attrList[0]
			tagStyle.Margin.Right = attrList[1]
			tagStyle.Margin.Bottom = attrList[2]
			tagStyle.Margin.Left = attrList[3]
		default:
			panic(fmt.Sprintf("unsupported margin value %v", cStyle))
		}
	case "border":
		attrList := strings.Split(cssValue, " ")
		if len(attrList) == 3 {
			tagStyle.BorderWidth.Left = attrList[0]
			tagStyle.BorderStyle.Left = attrList[1]
			tagStyle.BorderColor.Left = attrList[2]

			tagStyle.BorderWidth.Right = attrList[0]
			tagStyle.BorderStyle.Right = attrList[1]
			tagStyle.BorderColor.Right = attrList[2]

			tagStyle.BorderWidth.Top = attrList[0]
			tagStyle.BorderStyle.Top = attrList[1]
			tagStyle.BorderColor.Top = attrList[2]

			tagStyle.BorderWidth.Bottom = attrList[0]
			tagStyle.BorderStyle.Bottom = attrList[1]
			tagStyle.BorderColor.Bottom = attrList[2]
		} else {
			panic(fmt.Sprintf("unsupported border value %v", cStyle))
		}
	case "border-left":
		attrList := strings.Split(cssValue, " ")
		if len(attrList) == 3 {
			tagStyle.BorderWidth.Left = attrList[0]
			tagStyle.BorderStyle.Left = attrList[1]
			tagStyle.BorderColor.Left = attrList[2]
		} else {
			panic(fmt.Sprintf("unsupported border-left value %v", cStyle))
		}
	case "border-right":
		attrList := strings.Split(cssValue, " ")
		if len(attrList) == 3 {
			tagStyle.BorderWidth.Right = attrList[0]
			tagStyle.BorderStyle.Right = attrList[1]
			tagStyle.BorderColor.Right = attrList[2]
		} else {
			panic(fmt.Sprintf("unsupported border-right value %v", cStyle))
		}
	case "border-top":
		attrList := strings.Split(cssValue, " ")
		if len(attrList) == 3 {
			tagStyle.BorderWidth.Top = attrList[0]
			tagStyle.BorderStyle.Top = attrList[1]
			tagStyle.BorderColor.Top = attrList[2]
		} else {
			panic(fmt.Sprintf("unsupported border-top value %v", cStyle))
		}
	case "border-bottom":
		attrList := strings.Split(cssValue, " ")
		if len(attrList) == 3 {
			tagStyle.BorderWidth.Bottom = attrList[0]
			tagStyle.BorderStyle.Bottom = attrList[1]
			tagStyle.BorderColor.Bottom = attrList[2]
		} else {
			panic(fmt.Sprintf("unsupported border-bottom value %v", cStyle))
		}
	case "border-radius":
		attrList := strings.Split(cssValue, " ")
		switch len(attrList) {
		case 1:
			tagStyle.BorderRadius.Top = attrList[0]
			tagStyle.BorderRadius.Bottom = attrList[0]
			tagStyle.BorderRadius.Left = attrList[0]
			tagStyle.BorderRadius.Right = attrList[0]
		case 2:
			tagStyle.BorderRadius.Top = attrList[0]
			tagStyle.BorderRadius.Bottom = attrList[0]
			tagStyle.BorderRadius.Left = attrList[1]
			tagStyle.BorderRadius.Right = attrList[1]
		case 3:
			tagStyle.BorderRadius.Top = attrList[0]
			tagStyle.BorderRadius.Left = attrList[1]
			tagStyle.BorderRadius.Right = attrList[1]
			tagStyle.BorderRadius.Bottom = attrList[2]
		case 4:
			tagStyle.BorderRadius.Top = attrList[0]
			tagStyle.BorderRadius.Right = attrList[1]
			tagStyle.BorderRadius.Bottom = attrList[2]
			tagStyle.BorderRadius.Left = attrList[3]
		default:
			panic(fmt.Sprintf("unsupported border-radius value %v", cStyle))
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

func GetInheritStyle(pStyle *TagStyle, curStyle *TagStyle) *TagStyle {
	if pStyle == nil {
		pStyle = &TagStyle{}
	}
	if curStyle == nil {
		curStyle = &TagStyle{}
	}
	if curStyle.Color == "" && pStyle.Color != "" {
		curStyle.Color = pStyle.Color
	}
	if curStyle.FontSize == "" && pStyle.FontSize != "" {
		curStyle.FontSize = pStyle.FontSize
	}
	if curStyle.LineHeight == "" && pStyle.LineHeight != "" {
		curStyle.LineHeight = pStyle.LineHeight
	}
	if curStyle.FontFamily == "" && pStyle.FontFamily != "" {
		curStyle.FontFamily = pStyle.FontFamily
	}
	return curStyle
}

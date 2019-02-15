package img

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"

	"github.com/golang/freetype/truetype"
	"github.com/wnote/html2img/dom"
	"github.com/wnote/html2img/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type DrawCursor struct {
	OffsetX int
	OffsetY int

	FromX int
	EndX  int

	FromY int
	EndY  int

	NeedNewLine bool
}

func BodyDom2Img(bodyDom *dom.Dom) ([]byte, error) {
	bodyWidth := utils.GetIntSize(bodyDom.TagStyle.Width)
	bodyHeight := utils.GetIntSize(bodyDom.TagStyle.Height)
	dst := image.NewRGBA(image.Rect(0, 0, bodyWidth, bodyHeight))
	if bodyDom.TagStyle.BackgroundColor != "" {
		col := utils.GetColor(bodyDom.TagStyle.BackgroundColor)
		draw.Draw(dst, dst.Bounds(), &image.Uniform{C: col}, image.ZP, draw.Src)
	}
	DrawChildren(dst, bodyDom, bodyDom.TagStyle, bodyDom.Children)

	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, dst, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DrawChildren(dst *image.RGBA, parent *dom.Dom, pStyle *dom.TagStyle, children []*dom.Dom) {
	for _, d := range children {
		calcStyle := dom.GetInheritStyle(pStyle, d.TagStyle)

		if d.DomType == dom.DOM_TYPE_ELEMENT {
			switch d.TagName {
			case "img":
				imgData := d.TagData.(dom.ImageData)
				draw.Draw(dst, dst.Bounds().Add(image.Pt(d.Inner.X1, d.Inner.Y1)), imgData.Img, image.ZP, draw.Over)
			case "hr":
				if calcStyle.BackgroundColor == "" {
					calcStyle.BackgroundColor = "#000000"
				}
				fallthrough
			default:
				if calcStyle.BackgroundColor != "" {
					borderColor := utils.GetColor(calcStyle.BackgroundColor)
					for y := d.Inner.Y1; y <= d.Inner.Y2; y++ {
						for x := d.Inner.X1; x <= d.Inner.X2; x++ {
							dst.Set(x, y, borderColor)
						}
					}
				}
			}
			DrawChildren(dst, d, calcStyle, d.Children)
		} else if d.DomType == dom.DOM_TYPE_TEXT {
			f, exist := dom.FontMapping[calcStyle.FontFamily]
			if !exist {
				panic("Font-Family " + calcStyle.FontFamily + " not exist")
			}
			fontSize := utils.GetIntSize(calcStyle.FontSize)
			col := calcStyle.Color
			if col == "" {
				col = "#000000"
			}
			fontColor := utils.GetColor(col)
			AddText(f, float64(fontSize), 72, dst, image.NewUniform(fontColor), d.TagData.(string), d.Inner.X1, d.Inner.Y1+fontSize)
		} else {
			// Comments or other document type
		}
	}
}

func AddText(f *truetype.Font, size float64, dpi float64, dst *image.RGBA, src *image.Uniform, text string, x int, y int) {
	h := font.HintingNone
	fd := &font.Drawer{
		Dst: dst,
		Src: src,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     dpi,
			Hinting: h,
		}),
	}

	fd.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	fd.DrawString(text)
}

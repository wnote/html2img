package img

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"

	"github.com/golang/freetype/truetype"
	"github.com/wnote/html2img/conf"
	"github.com/wnote/html2img/dom"
	"github.com/wnote/html2img/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

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
				if calcStyle.Height == "" {
					calcStyle.Height = "1px"
				}
				fallthrough
			default:
				if calcStyle.BackgroundColor != "" {
					box := d.Container
					borderColor := utils.GetColor(calcStyle.BackgroundColor)
					for y := box.Y1; y <= box.Y2; y++ {
						for x := box.X1; x <= box.X2; x++ {
							dst.Set(x, y, borderColor)
						}
					}
				}
				if calcStyle.BorderStyle.Top != "" && calcStyle.BorderWidth.Top != "" && calcStyle.BorderColor.Top != "" {
					borderWidth := utils.GetIntSize(calcStyle.BorderWidth.Top)
					borderColor := utils.GetColor(calcStyle.BorderColor.Top)
					switch calcStyle.BorderStyle.Top {
					case "solid":
						box := d.Container
						for width := borderWidth - 1; width >= 0; width-- {
							for x := box.X1; x <= box.X2; x++ {
								dst.Set(x, d.Container.Y1+width, borderColor)
							}
						}
					default:
						panic("border-style " + calcStyle.BorderStyle.Top + " not support")

					}
				}

				if calcStyle.BorderStyle.Bottom != "" && calcStyle.BorderWidth.Bottom != "" && calcStyle.BorderColor.Bottom != "" {
					borderWidth := utils.GetIntSize(calcStyle.BorderWidth.Bottom)
					borderColor := utils.GetColor(calcStyle.BorderColor.Bottom)
					switch calcStyle.BorderStyle.Bottom {
					case "solid":
						box := d.Container
						for width := borderWidth - 1; width >= 0; width-- {
							for x := box.X1; x <= box.X2; x++ {
								dst.Set(x, d.Container.Y2-width, borderColor)
							}
						}
					default:
						panic("border-style " + calcStyle.BorderStyle.Top + " not support")

					}
				}

				if calcStyle.BorderStyle.Left != "" && calcStyle.BorderWidth.Left != "" && calcStyle.BorderColor.Left != "" {
					borderWidth := utils.GetIntSize(calcStyle.BorderWidth.Left)
					borderColor := utils.GetColor(calcStyle.BorderColor.Left)
					switch calcStyle.BorderStyle.Left {
					case "solid":
						box := d.Container
						for width := borderWidth - 1; width >= 0; width-- {
							for y := box.Y1; y <= box.Y2; y++ {
								dst.Set(d.Container.X1+width, y, borderColor)
							}
						}
					default:
						panic("border-style " + calcStyle.BorderStyle.Top + " not support")

					}
				}

				if calcStyle.BorderStyle.Right != "" && calcStyle.BorderWidth.Right != "" && calcStyle.BorderColor.Right != "" {
					borderWidth := utils.GetIntSize(calcStyle.BorderWidth.Right)
					borderColor := utils.GetColor(calcStyle.BorderColor.Right)
					switch calcStyle.BorderStyle.Right {
					case "solid":
						box := d.Container
						for width := borderWidth - 1; width >= 0; width-- {
							for y := box.Y1; y <= box.Y2; y++ {
								dst.Set(d.Container.X2-width, y, borderColor)
							}
						}
					default:
						panic("border-style " + calcStyle.BorderStyle.Top + " not support")

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
			AddText(f, float64(fontSize), dst, image.NewUniform(fontColor), d.TagData.(string), d.Inner.X1, d.Inner.Y1+11*fontSize/12)
		} else {
			// Comments or other document type
		}
	}
}

func AddText(f *truetype.Font, size float64, dst *image.RGBA, src *image.Uniform, text string, x int, y int) {
	fd := &font.Drawer{
		Dst: dst,
		Src: src,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    size,
			DPI:     conf.DPI,
			Hinting: font.HintingNone,
		}),
	}

	fd.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	fd.DrawString(text)
}

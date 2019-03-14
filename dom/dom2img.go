package dom

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/wnote/html2img/conf"
	"github.com/wnote/html2img/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func BodyDom2Img(bodyDom *Dom) ([]byte, error) {
	bodyWidth := utils.GetIntSize(bodyDom.TagStyle.Width)
	bodyHeight := utils.GetIntSize(bodyDom.TagStyle.Height)
	dst := image.NewRGBA(image.Rect(0, 0, bodyWidth, bodyHeight))
	if bodyDom.TagStyle.BackgroundColor != "" {
		col := utils.GetColor(bodyDom.TagStyle.BackgroundColor)
		draw.Draw(dst, dst.Bounds(), &image.Uniform{C: col}, image.ZP, draw.Src)
	}
	DrawChildren(dst, bodyDom.TagStyle, bodyDom.Children)

	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, dst, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DrawChildren(dst *image.RGBA, pStyle *TagStyle, children []*Dom) {
	for _, d := range children {
		calcStyle := GetInheritStyle(pStyle, d.TagStyle)

		if d.DomType == DOM_TYPE_ELEMENT {
			switch d.TagName {
			case "img":
				imgData := d.TagData.(ImageData)
				draw.Draw(dst, dst.Bounds().Add(image.Pt(d.Inner.X1, d.Inner.Y1)), imgData.Img, image.ZP, draw.Over)
				DrawBoxRadius(dst, d.Container, calcStyle, pStyle)
			default:
				box := d.Container
				if calcStyle.BackgroundColor != "" {
					borderColor := utils.GetColor(calcStyle.BackgroundColor)
					for y := box.Y1; y <= box.Y2; y++ {
						for x := box.X1; x <= box.X2; x++ {
							dst.Set(x, y, borderColor)
						}
					}
				}
				DrawBoxRadius(dst, box, calcStyle, pStyle)

				borderTopRadius := utils.GetIntSize(calcStyle.BorderRadius.Top)
				borderRightRadius := utils.GetIntSize(calcStyle.BorderRadius.Right)
				borderBottomRadius := utils.GetIntSize(calcStyle.BorderRadius.Bottom)
				borderLeftRadius := utils.GetIntSize(calcStyle.BorderRadius.Left)

				width := d.Container.X2 - d.Container.X1 + 1
				height := d.Container.Y2 - d.Container.Y1 + 1
				var halfSize int
				if width > height {
					halfSize = height / 2
				} else {
					halfSize = width / 2
				}
				if borderTopRadius > halfSize {
					borderTopRadius = halfSize
				}
				if borderRightRadius > halfSize {
					borderRightRadius = halfSize
				}
				if borderBottomRadius > halfSize {
					borderBottomRadius = halfSize
				}
				if borderLeftRadius > halfSize {
					borderLeftRadius = halfSize
				}
				if calcStyle.BorderStyle.Top != "" && calcStyle.BorderWidth.Top != "" && calcStyle.BorderColor.Top != "" {
					borderWidth := utils.GetIntSize(calcStyle.BorderWidth.Top)
					borderColor := utils.GetColor(calcStyle.BorderColor.Top)
					switch calcStyle.BorderStyle.Top {
					case "solid":
						for width := borderWidth - 1; width >= 0; width-- {
							r := borderTopRadius - width
							for xxOffset := r; xxOffset >= 0; xxOffset-- {
								yyOffset := int(math.Sqrt(float64(r*r - xxOffset*xxOffset)))
								dst.Set(d.Container.X1+r-int(xxOffset), box.Y1+r-yyOffset, borderColor)
								dst.Set(d.Container.X1+r-int(yyOffset), box.Y1+r-xxOffset, borderColor)
							}
							for x := box.X1 + borderTopRadius; x <= box.X2-borderRightRadius; x++ {
								dst.Set(x, d.Container.Y1+width, borderColor)
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
						for width := borderWidth - 1; width >= 0; width-- {
							r := borderRightRadius - width
							for xxOffset := r; xxOffset >= 0; xxOffset-- {
								yyOffset := int(math.Sqrt(float64(r*r - xxOffset*xxOffset)))
								dst.Set(d.Container.X2-r+int(xxOffset), box.Y1+r-yyOffset, borderColor)
								dst.Set(d.Container.X2-r+int(yyOffset), box.Y1+r-xxOffset, borderColor)
							}
							for y := box.Y1 + borderRightRadius; y <= box.Y2-borderBottomRadius; y++ {
								dst.Set(d.Container.X2-width, y, borderColor)
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
						for width := borderWidth - 1; width >= 0; width-- {
							r := borderBottomRadius - width
							for xxOffset := r; xxOffset >= 0; xxOffset-- {
								yyOffset := int(math.Sqrt(float64(r*r - xxOffset*xxOffset)))
								dst.Set(d.Container.X2-r+int(xxOffset), box.Y2-r+yyOffset, borderColor)
								dst.Set(d.Container.X2-r+int(yyOffset), box.Y2-r+xxOffset, borderColor)
							}
							for x := box.X1 + borderLeftRadius; x <= box.X2-borderBottomRadius; x++ {
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
						for width := borderWidth - 1; width >= 0; width-- {
							r := borderLeftRadius - width
							for xxOffset := r; xxOffset >= 0; xxOffset-- {
								yyOffset := int(math.Sqrt(float64(r*r - xxOffset*xxOffset)))
								dst.Set(d.Container.X1+r-int(xxOffset), box.Y2-r+yyOffset, borderColor)
								dst.Set(d.Container.X1+r-int(yyOffset), box.Y2-r+xxOffset, borderColor)
							}
							for y := box.Y1 + borderTopRadius; y <= box.Y2-borderLeftRadius; y++ {
								dst.Set(d.Container.X1+width, y, borderColor)
							}
						}
					default:
						panic("border-style " + calcStyle.BorderStyle.Top + " not support")

					}
				}

			}
			DrawChildren(dst, calcStyle, d.Children)
		} else if d.DomType == DOM_TYPE_TEXT {
			f, exist := FontMapping[calcStyle.FontFamily]
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

func OutCircle(x, y, radius int) bool {
	xf := float64(x) + 0.5
	yf := float64(y) + 0.5
	rf := float64(radius)
	return yf*yf+xf*xf > rf*rf
}

func DrawBoxRadius(dst *image.RGBA, box Rectangle, cStyle *TagStyle, pStyle *TagStyle) {
	borderTopRadius := utils.GetIntSize(cStyle.BorderRadius.Top)
	borderRightRadius := utils.GetIntSize(cStyle.BorderRadius.Right)
	borderBottomRadius := utils.GetIntSize(cStyle.BorderRadius.Bottom)
	borderLeftRadius := utils.GetIntSize(cStyle.BorderRadius.Left)

	width := box.X2 - box.X1 + 1
	height := box.Y2 - box.Y1 + 1
	var halfSize int
	if width > height {
		halfSize = height / 2
	} else {
		halfSize = width / 2
	}
	if borderTopRadius > halfSize {
		borderTopRadius = halfSize
	}
	if borderRightRadius > halfSize {
		borderRightRadius = halfSize
	}
	if borderBottomRadius > halfSize {
		borderBottomRadius = halfSize
	}
	if borderLeftRadius > halfSize {
		borderLeftRadius = halfSize
	}

	col := color.RGBA{
		R: uint8(255),
		G: uint8(255),
		B: uint8(255),
		A: uint8(255),
	}

	if pStyle.BackgroundColor != "" {
		pColor := utils.GetColor(pStyle.BackgroundColor)
		r, g, b, a := pColor.RGBA()
		col = color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		}
	}

	for x := 0; x <= borderTopRadius; x++ {
		for y := 0; y <= borderTopRadius; y++ {
			if OutCircle(x, y, borderTopRadius) {
				offsetX := borderTopRadius - x
				offsetY := borderTopRadius - y
				dst.Set(box.X1+offsetX, box.Y1+offsetY, col)
			}
		}
	}
	for x := 0; x <= borderRightRadius; x++ {
		for y := 0; y <= borderRightRadius; y++ {
			if OutCircle(x, y, borderRightRadius) {
				offsetX := borderRightRadius - x
				offsetY := borderRightRadius - y
				dst.Set(box.X2-offsetX, box.Y1+offsetY, col)
			}
		}
	}
	for x := 0; x <= borderBottomRadius; x++ {
		for y := 0; y <= borderBottomRadius; y++ {
			if OutCircle(x, y, borderBottomRadius) {
				offsetX := borderBottomRadius - x
				offsetY := borderBottomRadius - y

				dst.Set(box.X2-offsetX, box.Y2-offsetY, col)
			}
		}
	}
	for x := 0; x <= borderLeftRadius; x++ {
		for y := 0; y <= borderLeftRadius; y++ {
			if OutCircle(x, y, borderLeftRadius) {
				offsetX := borderLeftRadius - x
				offsetY := borderLeftRadius - y
				dst.Set(box.X1+offsetX, box.Y2-offsetY, col)
			}
		}
	}
}

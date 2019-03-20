package html2img

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/wnote/html2img/conf"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func bodyDom2Img(bodyDom *Dom) ([]byte, error) {
	bodyWidth := getIntSize(bodyDom.TagStyle.Width)
	bodyHeight := getIntSize(bodyDom.TagStyle.Height)
	dst := image.NewRGBA(image.Rect(0, 0, bodyWidth, bodyHeight))
	if bodyDom.TagStyle.BackgroundColor != "" {
		col := getColor(bodyDom.TagStyle.BackgroundColor)
		draw.Draw(dst, dst.Bounds(), &image.Uniform{C: col}, image.ZP, draw.Src)
	}
	drawChildren(dst, bodyDom.TagStyle, bodyDom.Children)

	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, dst, &jpeg.Options{
		Quality: 100,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func drawChildren(dst *image.RGBA, pStyle *TagStyle, children []*Dom) {
	for _, d := range children {
		calcStyle := getInheritStyle(pStyle, d.TagStyle)

		if d.DomType == DOM_TYPE_ELEMENT {
			switch d.TagName {
			case "img":
				imgData := d.TagData.(ImageData)
				draw.Draw(dst, dst.Bounds().Add(image.Pt(d.Inner.X1, d.Inner.Y1)), imgData.Img, image.ZP, draw.Over)
				drawBoxRadius(dst, d.Container, calcStyle, pStyle)
			default:
				box := d.Container
				if calcStyle.BackgroundColor != "" {
					borderColor := getColor(calcStyle.BackgroundColor)
					for y := box.Y1; y <= box.Y2; y++ {
						for x := box.X1; x <= box.X2; x++ {
							dst.Set(x, y, borderColor)
						}
					}
				}
				drawBoxRadius(dst, box, calcStyle, pStyle)

				borderTopRadius := getIntSize(calcStyle.BorderRadius.Top)
				borderRightRadius := getIntSize(calcStyle.BorderRadius.Right)
				borderBottomRadius := getIntSize(calcStyle.BorderRadius.Bottom)
				borderLeftRadius := getIntSize(calcStyle.BorderRadius.Left)

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
					borderWidth := getIntSize(calcStyle.BorderWidth.Top)
					borderColor := getColor(calcStyle.BorderColor.Top)
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
					borderWidth := getIntSize(calcStyle.BorderWidth.Right)
					borderColor := getColor(calcStyle.BorderColor.Right)
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
					borderWidth := getIntSize(calcStyle.BorderWidth.Bottom)
					borderColor := getColor(calcStyle.BorderColor.Bottom)
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
					borderWidth := getIntSize(calcStyle.BorderWidth.Left)
					borderColor := getColor(calcStyle.BorderColor.Left)
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
			drawChildren(dst, calcStyle, d.Children)
		} else if d.DomType == DOM_TYPE_TEXT {
			f, exist := fontMapping[calcStyle.FontFamily]
			if !exist {
				fmt.Println("Font-Family " + calcStyle.FontFamily + " not exist")
				for _, fon := range fontMapping {
					f = fon
					break
				}
			}
			fontSize := getIntSize(calcStyle.FontSize)
			col := calcStyle.Color
			if col == "" {
				col = "#000000"
			}
			fontColor := getColor(col)
			addText(f, float64(fontSize), dst, image.NewUniform(fontColor), d.TagData.(string), d.Inner.X1, d.Inner.Y1+11*fontSize/12)
		} else {
			// Comments or other document type
		}
	}
}

func addText(f *truetype.Font, size float64, dst *image.RGBA, src *image.Uniform, text string, x int, y int) {
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

func outOfCircle(x, y, radius int) bool {
	xf := float64(x) + 0.5
	yf := float64(y) + 0.5
	rf := float64(radius)
	return yf*yf+xf*xf > rf*rf
}

func drawBoxRadius(dst *image.RGBA, box Rectangle, cStyle *TagStyle, pStyle *TagStyle) {
	borderTopRadius := getIntSize(cStyle.BorderRadius.Top)
	borderRightRadius := getIntSize(cStyle.BorderRadius.Right)
	borderBottomRadius := getIntSize(cStyle.BorderRadius.Bottom)
	borderLeftRadius := getIntSize(cStyle.BorderRadius.Left)

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
		pColor := getColor(pStyle.BackgroundColor)
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
			if outOfCircle(x, y, borderTopRadius) {
				offsetX := borderTopRadius - x
				offsetY := borderTopRadius - y
				dst.Set(box.X1+offsetX, box.Y1+offsetY, col)
			}
		}
	}
	for x := 0; x <= borderRightRadius; x++ {
		for y := 0; y <= borderRightRadius; y++ {
			if outOfCircle(x, y, borderRightRadius) {
				offsetX := borderRightRadius - x
				offsetY := borderRightRadius - y
				dst.Set(box.X2-offsetX, box.Y1+offsetY, col)
			}
		}
	}
	for x := 0; x <= borderBottomRadius; x++ {
		for y := 0; y <= borderBottomRadius; y++ {
			if outOfCircle(x, y, borderBottomRadius) {
				offsetX := borderBottomRadius - x
				offsetY := borderBottomRadius - y

				dst.Set(box.X2-offsetX, box.Y2-offsetY, col)
			}
		}
	}
	for x := 0; x <= borderLeftRadius; x++ {
		for y := 0; y <= borderLeftRadius; y++ {
			if outOfCircle(x, y, borderLeftRadius) {
				offsetX := borderLeftRadius - x
				offsetY := borderLeftRadius - y
				dst.Set(box.X1+offsetX, box.Y2-offsetY, col)
			}
		}
	}
}

package img

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	// 画布宽度
	SHARE_CARD_IMG_WIDTH = 750
	// 画布长度
	SHARE_CARD_IMG_HEIGHT = 1334
)

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{A: 255}
	}
	return color.Alpha{A: 0}
}

type drawService struct {
}

func (d *drawService) AddText(fontfile string, size float64, dpi float64, dst *image.RGBA, src *image.Uniform, text string, x int, y int) {
	if x == 0 {
		x = int((SHARE_CARD_IMG_WIDTH - (float64(d.CalStrLen(text))*size*dpi/72)/3) / 2)
	}

	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

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

func (d *drawService) CalStrLen(str string) float64 {
	sl := 0.0
	rs := []rune(str)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			sl += 1.7
		} else {
			sl += 3
		}
	}
	return sl
}

func (d *drawService) GetHttpImg(imgUrl string) (image.Image, string, error) {
	resp, err := http.Get(imgUrl)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	img, fm, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return img, fm, err
}

func (d *drawService) AddMultiLineText(fontFile string, size float64, dpi float64, dst *image.RGBA, src *image.Uniform, text string, x int, y int, maxLine int, lineHeight int, appendSuffix string, paddingRight int) int {
	maxLineWidth := float64(SHARE_CARD_IMG_WIDTH-x-paddingRight) - size
	var multiTextGroupByLine []string
	var tmpStr string
	var tmpWidth float64
	for _, value := range text {
		subStr := string(value)
		strWidth := (float64(d.CalStrLen(subStr)) * size * dpi / 72) / 3
		tmpWidth += strWidth
		beforeAdd := tmpStr
		tmpStr += subStr
		if tmpWidth > maxLineWidth {
			if len(multiTextGroupByLine) >= (maxLine - 1) {
				if appendSuffix != "" {
					tmpStr = beforeAdd + appendSuffix
				}
				multiTextGroupByLine = append(multiTextGroupByLine, tmpStr)
				tmpWidth = 0
				tmpStr = ""
				break
			}
			multiTextGroupByLine = append(multiTextGroupByLine, tmpStr)
			tmpWidth = 0
			tmpStr = ""
		}
	}
	if tmpStr != "" {
		multiTextGroupByLine = append(multiTextGroupByLine, tmpStr)
	}
	rtnY := y
	for key, subStr := range multiTextGroupByLine {
		d.AddText(fontFile, size, dpi, dst, src, subStr, x, y+key*lineHeight)
		rtnY = y + key*lineHeight
	}
	return rtnY + lineHeight
}

// 画水平线
func (d *drawService) DrawHorizLine(img *image.RGBA, color color.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

// 画矩形,无填充
func (d *drawService) DrawRectangle(dst draw.Image, startX int, startY int, endX int, endY int, borderColor color.RGBA) {
	img := image.NewNRGBA(image.Rect(startX, startY, endX, endY))
	draw.Draw(dst, img.Bounds(), &image.Uniform{C: borderColor}, image.ZP, draw.Src)
	img = image.NewNRGBA(image.Rect(startX+1, startY+1, endX-1, endY-1))
	draw.Draw(dst, img.Bounds(), image.White, image.ZP, draw.Src)
}

// 画头像
func (d *drawService) DrawCardAvatar(dst draw.Image, headerUrl string, avatarSize uint, offsetX int, offsetY int) error {
	// 绘制头像
	avatarImg, _, err := d.GetHttpImg(headerUrl)
	if err != nil {
		return err
	}
	var tmpAvatarSize uint
	if avatarImg.Bounds().Max.X > avatarImg.Bounds().Max.Y {
		tmpAvatarSize = uint(avatarImg.Bounds().Max.Y)
	} else {
		tmpAvatarSize = uint(avatarImg.Bounds().Max.X)
	}
	if tmpAvatarSize > 500 {
		tmpAvatarSize = 500
	}
	avatarImg = resize.Resize(tmpAvatarSize, tmpAvatarSize, avatarImg, resize.Lanczos3)
	p := image.Point{
		X: avatarImg.Bounds().Max.X / 2,
		Y: avatarImg.Bounds().Max.Y / 2,
	}
	tmpDst := image.NewRGBA(image.Rect(0, 0, int(tmpAvatarSize), int(tmpAvatarSize)))
	draw.Draw(tmpDst, tmpDst.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	draw.DrawMask(tmpDst, tmpDst.Bounds().Add(image.Pt(0, 0)), avatarImg, image.ZP, &circle{p, int(tmpAvatarSize / 2)}, image.ZP, draw.Over)
	avatarImg = tmpDst.SubImage(image.Rect(0, 0, int(tmpAvatarSize), int(tmpAvatarSize)))
	avatarImg = resize.Resize(avatarSize, avatarSize, avatarImg, resize.Lanczos3)

	draw.Draw(dst, dst.Bounds().Add(image.Pt(offsetX, offsetY)), avatarImg, image.ZP, draw.Over)
	return nil
}

// 获取基础图片资源
func (d *drawService) GetBasicInfoDraw(coverImageUrl string, salerUserName string, saleMobile string, jobTitle string, headerUrl string, qrCodeImgUrl string, drawSellerBgOffset int) (*image.RGBA, error) {
	coverHeight := 422
	imageWidth := SHARE_CARD_IMG_WIDTH
	// 画布尺寸大小
	dst := image.NewRGBA(image.Rect(0, 0, imageWidth, SHARE_CARD_IMG_HEIGHT))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	// 绘制项目封面
	if coverImageUrl != "" {
		coverImg, fm, err := d.GetHttpImg(coverImageUrl)
		if err != nil {
			log.Fatalf("d.GetHttpImg err: %v", coverImageUrl)
			return nil, err
		}
		coverBounds := coverImg.Bounds()
		h := int(coverBounds.Dx() * coverHeight / imageWidth)
		// 等比例居中裁剪
		if h < coverBounds.Dy() {
			removeY := int((coverBounds.Dy() - h) / 2)
			coverImg, err = d.Clip(coverImg, fm, 0, removeY, coverBounds.Dx(), coverBounds.Dy()-removeY)
		} else {
			w := int(coverBounds.Dy() * imageWidth / coverHeight)
			removeX := int((coverBounds.Dx() - w) / 2)
			coverImg, err = d.Clip(coverImg, fm, removeX, 0, coverBounds.Dx()-removeX, coverBounds.Dy())
		}
		if err != nil {
			log.Fatalf("d.Clip err: %v", err)
			return nil, err
		}
		coverImg = resize.Resize(uint(imageWidth), uint(coverHeight), coverImg, resize.Lanczos3)
		draw.Draw(dst, dst.Bounds().Add(image.Pt(0, 0)), coverImg, image.ZP, draw.Over)
	}
	/*if drawSellerBgOffset > 0 && salerUserName != "" {
		sellerBgImg, _, err := d.GetHttpImg(vars.ShareCardSetting.ShareCardSellerBgImgUrl)
		if err != nil {
			log.Fatalf("d.GetHttpImg err: %v", vars.ShareCardSetting.ShareCardSellerBgImgUrl)
			return nil, err
		}
		sellerBgImg = resize.Resize(750, 184, sellerBgImg, resize.Lanczos3)
		draw.Draw(dst, dst.Bounds().Add(image.Pt(0, drawSellerBgOffset)), sellerBgImg, image.ZP, draw.Over)

		d.DrawCardAvatar(dst, headerUrl, 100, 58, drawSellerBgOffset+42)
		d.AddText(vars.ShareCardSetting.PingFangBdFontPath, 30, 72, dst, image.NewUniform(color.RGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xff}), salerUserName, 188, drawSellerBgOffset+79)
		d.AddText(vars.ShareCardSetting.PingFangFontPath, 24, 72, dst, image.NewUniform(color.RGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff}), saleMobile, 536, drawSellerBgOffset+81)
		d.AddText(vars.ShareCardSetting.PingFangFontPath, 24, 72, dst, image.NewUniform(color.RGBA{R: 0x84, G: 0x84, B: 0x84, A: 0xff}), jobTitle, 188, drawSellerBgOffset+121)
	}

	qrCodeBgImg, _, err := d.GetHttpImg(vars.ShareCardSetting.ShareCardQrcodeBgImgUrl)
	if err != nil {
		log.Fatalf("d.GetHttpImg err: %v", vars.ShareCardSetting.ShareCardQrcodeBgImgUrl)
		return nil, err
	}*/
	/*qrCodeBgImg = resize.Resize(694, 224, qrCodeBgImg, resize.Lanczos3)
	draw.Draw(dst, dst.Bounds().Add(image.Pt(28, 1082)), qrCodeBgImg, image.ZP, draw.Over)*/

	// 绘制二维码
	qrCodeImg, _, err := d.GetHttpImg(qrCodeImgUrl)
	if err != nil {
		log.Fatalf("d.GetHttpImg err: %v", qrCodeImgUrl)
		return nil, err
	}
	qrCodeImg = resize.Resize(176, 176, qrCodeImg, resize.Lanczos3)
	draw.Draw(dst, dst.Bounds().Add(image.Pt(128, 1106)), qrCodeImg, image.ZP, draw.Over)
	return dst, nil
}

// 裁剪图片
func (d *drawService) Clip(origin image.Image, format string, x0, y0, x1, y1 int) (image.Image, error) {
	switch format {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return subImg, nil
	case "png":
		switch origin.(type) {
		case *image.NRGBA:
			img := origin.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return subImg, nil
		case *image.RGBA:
			img := origin.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return subImg, nil
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return subImg, nil
	}
	return nil, errors.New("ERROR FORMAT")
}

// 给实心矩形绘制圆角
func (d *drawService) DrawRectBorderRadius(img *image.RGBA, fromX, fromY, toX, toY int, radius int, col color.Color) {
	toY = toY - 1
	toX = toX - 1
	for x := 0; x <= radius; x++ {
		for y := 0; y <= radius; y++ {
			if x*x+y*y > radius*radius {
				offsetX := radius - x
				offsetY := radius - y
				img.Set(fromX+offsetX, fromY+offsetY, col)
				img.Set(fromX+offsetX, toY-offsetY, col)
				img.Set(toX-offsetX, fromY+offsetY, col)
				img.Set(toX-offsetX, toY-offsetY, col)
			}
		}
	}
}

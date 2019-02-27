package utils

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

func GetIntPx(size string, pSize int) int {
	if size == "" {
		return 0
	}
	re := regexp.MustCompile("\\d+px")
	if re.MatchString(size) {
		ignoreUnitPx := strings.Replace(size, "px", "", 1)
		px, err := strconv.Atoi(ignoreUnitPx)
		if err != nil {
			panic(fmt.Sprintf("size err :%v", size))
		}
		return px
	}
	re = regexp.MustCompile("\\d+%")
	if re.MatchString(size) {
		sizePercent := strings.Replace(size, "%", "", 1)
		percent, err := strconv.Atoi(sizePercent)
		if err != nil {
			panic(fmt.Sprintf("size err :%v", size))
		}
		if pSize < 0 {
			panic(fmt.Sprintf("parent size err :%v", pSize))
		}
		return percent * pSize / 100
	}
	return 0
}

func GetIntSize(size string) int {
	return GetIntPx(size, 0)
}

func GetColor(colorStr string) color.Color {
	escapeColor := strings.Replace(colorStr, "#", "", 1)
	if len(escapeColor) == 3 {
		escapeColor = escapeColor[:1] + escapeColor[:1] + escapeColor[1:2] + escapeColor[1:2] + escapeColor[2:3] + escapeColor[2:3]
	} else if len(escapeColor) < 6 {
		panic(fmt.Sprintf("color err :%v", colorStr))
	}
	r, err := strconv.ParseInt(escapeColor[:2], 16, 32)
	if err != nil {
		panic(fmt.Sprintf("color err :%v", colorStr))
	}
	g, err := strconv.ParseInt(escapeColor[2:4], 16, 32)
	if err != nil {
		panic(fmt.Sprintf("color err :%v", colorStr))
	}
	b, err := strconv.ParseInt(escapeColor[4:6], 16, 32)
	if err != nil {
		panic(fmt.Sprintf("color err :%v", colorStr))
	}
	a := uint8(255)
	if len(escapeColor) == 8 {
		alp, err := strconv.ParseInt(escapeColor[6:8], 16, 32)
		if err != nil {
			panic(fmt.Sprintf("color err :%v", colorStr))
		}
		if alp > 0 {
			a = uint8(alp)
		}
	}
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

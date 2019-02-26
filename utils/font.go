package utils

func CalcCharacterPx(text string, fontSize float64) float64 {
	dpi := float64(72)
	return (float64(CalCharacterLen(text)) * fontSize * dpi / 72) / 3
}

func CalCharacterLen(str string) float64 {
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

func SplitMultiLineText(text string, size float64, maxLineWidth float64) []string {
	dpi := float64(72)
	var multiTextGroupByLine []string
	var tmpStr string
	var tmpWidth float64
	for _, value := range text {
		subStr := string(value)
		strWidth := (float64(CalCharacterLen(subStr)) * size * dpi / 72) / 3
		tmpWidth += strWidth
		tmpStr += subStr
		if tmpWidth > maxLineWidth {
			multiTextGroupByLine = append(multiTextGroupByLine, tmpStr)
			tmpWidth = 0
			tmpStr = ""
		}
	}
	if tmpStr != "" {
		multiTextGroupByLine = append(multiTextGroupByLine, tmpStr)
	}

	return multiTextGroupByLine
}

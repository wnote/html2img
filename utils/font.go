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

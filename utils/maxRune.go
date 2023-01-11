package utils

// MaxRune максимально возможной руны
func MaxRune(runeSlice []rune, num int) int {
	runeSliceLen := len(runeSlice)
	if runeSliceLen >= num {
		return num
	}

	return runeSliceLen
}

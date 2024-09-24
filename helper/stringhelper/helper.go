package stringhelper

func DelChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}

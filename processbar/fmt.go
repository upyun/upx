package processbar

import "strings"

func shortPath(s string, width int) string {
	if slen(s) <= width {
		return s
	}

	dotLen := 3
	headLen := (width - dotLen) / 2
	tailLen := width - dotLen - headLen

	st := 1
	for ; st < len(s); st++ {
		if slen(s[0:st]) > headLen {
			break
		}
	}

	ed := len(s) - 1
	for ; ed >= 0; ed-- {
		if slen(s[ed:]) > tailLen {
			break
		}
	}

	return s[0:st-1] + strings.Repeat(".", dotLen) + s[ed+1:]
}

func leftAlign(s string, width int) string {
	l := slen(s)
	for i := 0; i < width-l; i++ {
		s += " "
	}
	return s
}
func rightAlign(s string, width int) string {
	l := slen(s)
	for i := 0; i < width-l; i++ {
		s = " " + s
	}
	return s
}

func slen(s string) int {
	l, rl := len(s), len([]rune(s))
	return (l-rl)/2 + rl
}

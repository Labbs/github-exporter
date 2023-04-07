package internal

import "strings"

func SplitAndClean(str string) []string {
	s := strings.Split(str, "\n")
	r := make([]string, 0)
	for _, v := range s {
		if v != "" {
			r = append(r, v)
		}
	}
	return r
}

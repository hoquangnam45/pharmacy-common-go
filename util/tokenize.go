package util

import "strings"

func Tokenize(val string, delimiter string) []string {
	l := strings.Split(val, delimiter)
	ret := []string{}
	for _, v := range l {
		ret = append(ret, strings.Trim(v, " "))
	}
	return ret
}

package util

import "strings"

func Tokenize(val string, delimiter string) []string {
	l := strings.Split(val, delimiter)
	ret := []string{}
	for _, v := range l {
		token := strings.Trim(v, " ")
		if token != "" {
			ret = append(ret, token)
		}
	}
	return ret
}

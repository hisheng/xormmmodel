package xorm

import "strings"

func UcFirst(str string) string {
	strLen := len(str)
	if strLen == 0 {
		return ""
	} else if strLen == 1 {
		return strings.ToUpper(str)
	} else {
		return strings.ToUpper(str[:1]) + str[1:]
	}
}

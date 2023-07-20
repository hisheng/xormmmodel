package xorm

import (
	"fmt"
	"os"
	"strings"
)

// ModFilePath 获取本项目mod所有的目录
func ModFilePath() string {
	wd, _ := os.Getwd()
	part := strings.Split(wd, "/")
	path := ""
	//列出所有可行的目录
	var isModWdArr []string
	for _, p := range part {
		if p == "" {
			continue
		}
		path += "/" + p
		isModWdArr = append(isModWdArr, path)
	}
	//优先级为最近的优先筛选
	var mdFile = ""
	var modRoot = ""
	for i := len(isModWdArr) - 1; i >= 0; i-- {
		modWd := isModWdArr[i]
		mdFile = modWd + "/go.mod"
		if FileExists(mdFile) {
			modRoot = modWd
			break
		}
	}
	fmt.Println(modRoot)
	fmt.Println(mdFile)
	return modRoot
}

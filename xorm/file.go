package xorm

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// 把字符串写入到文件中 类似于php file_put_contents
func filePutContents(filename string, data string) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	f.WriteString(data)
	return nil
}

// 判断文件/文件夹是否存在 -功能类似php的file_exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func saveToFile(tableName string, tableComment string, structStrArr []string) {
	abs, _ := filepath.Abs("")
	packageName := filepath.Base(abs)
	importStr := ""
	fileStr := "package " + packageName + "\n" + importStr + "\n"
	fileStr += "// " + upperCamelCase(tableName) + " " + tableComment + "\n"
	structStr := "type " + upperCamelCase(tableName) + " struct { \n"
	for _, row := range structStrArr {
		structStr += row
	}
	structStr += "} \n"
	fileStr += structStr

	// 2 TableName
	tbStr := fmt.Sprintf("func (m *%s) TableName() string { \n return \"%s\" \n }\n", upperCamelCase(tableName), tableName)
	// fmt.Println(fileStr)
	fileStr += tbStr
	saveFile(tableName, fileStr)

}

func saveFile(tableName, fileStr string) error {
	savePath, _ := filepath.Abs("")
	filePath := savePath + "/" + tableName + ".go"
	fmt.Println("生成完成 " + filePath)
	err := filePutContents(filePath, fileStr)
	if err != nil {
		return err
	}
	cmd := exec.Command("gofmt", "-w", filePath)
	cmd.Run()
	return nil
}

// 大驼峰式命名法（upper camel case）
func upperCamelCase(s string) string {
	strArr := strings.Split(s, "_")
	upperCamelSting := ""
	for _, v := range strArr {
		upperCamelSting += UcFirst(v)
	}
	return upperCamelSting
}

type YamlFile struct {
	Data struct {
		Database struct {
			Driver string
			Source string
		}
		Mysql struct {
			Default struct {
				Driver string
				Dsn    string
			}
		}
	}
}

func ReadYamlFile(yamlFilePath string) YamlFile {
	// 1读取文件
	data, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		fmt.Println(err)
	}
	// 2解析文件
	var y YamlFile
	err = yaml.Unmarshal(data, &y)
	fmt.Println(y, err)
	return y
}

func ConfigFilePath() string {
	wd, _ := os.Getwd()
	part := strings.Split(wd, "/")
	path := ""
	arrivalProjPath := false
	for _, p := range part {
		if p == "" {
			continue
		}
		path += "/" + p
		if arrivalProjPath {
			break
		}
		if p == "api-ad.qmniu.com" {
			arrivalProjPath = true
			break
		}

		if p == "pixiu-ads-server" {
			arrivalProjPath = true
			break
		}

		if p == "pixiu-ad-backend" {
			arrivalProjPath = true
			break
		}
	}
	path += "/configs/config.yaml"
	return path
}

func XormModelFilePath() string {
	wd, _ := os.Getwd()
	part := strings.Split(wd, "/")
	path := ""
	arrivalProjPath := false
	for _, p := range part {
		if p == "" {
			continue
		}
		path += "/" + p
		if arrivalProjPath {
			break
		}
		if p == "api-ad.qmniu.com" {
			arrivalProjPath = true
			break
		}

		if p == "pixiu-ads-server" {
			arrivalProjPath = true
			break
		}

		if p == "pixiu-ad-backend" {
			arrivalProjPath = true
			break
		}
	}
	path += "/xorm_model.yaml"
	return path
}

func SaveXormModelFile(filePath, dsn string) error {
	fileStr := "driver: mysql\n"
	fileStr += fmt.Sprintf("dsn: \"%s\"\n", dsn)
	fmt.Println("生成完成 " + filePath)
	err := filePutContents(filePath, fileStr)
	if err != nil {
		return err
	}
	cmd := exec.Command("gofmt", "-w", filePath)
	cmd.Run()
	return nil
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hisheng/xormmodel/xorm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var (
	xormConfigPath string
	xormDsn        string
	xormDbTable    string
	xormDb         string
	xormTable      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xormmodel",
	Short: "xormmodel database/table",
	Long:  `xormmodel database/table`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		// 1、获取database 与 table 值
		initXormDbTable(args)
		if xormDb == "" {
			xormHelp(cmd)
			return
		}
		// 2、获取dsn值
		xormDsn = initXormDsn()
		if xormDsn == "" {
			fmt.Println("未找到dsn配置")
			return
		}
		// //3.生成struct文件
		xorm.InitStruct(xormDsn, xormTable)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&xormConfigPath, "config", "c", "", "mysql配置文件地址")
	rootCmd.Flags().StringVarP(&xormDsn, "xormDsn", "d", "", "dsn值")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".xormmodel" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".xormmodel")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func initXormDbTable(args []string) {
	if len(args) > 0 {
		xormDbTable = args[0]
		dbTableArr := strings.Split(xormDbTable, "/")
		if len(dbTableArr) > 0 {
			xormDb = dbTableArr[0]
		}
		if len(dbTableArr) > 1 {
			xormTable = dbTableArr[1]
		}
	}
}

func xormHelp(cmd *cobra.Command) {
	cmd.Usage()
	cmd.Println("\n")
	cmd.Println("\033[;36m 支持下面几种方式的cmd执行 \033[0m")
	cmd.Println("xormmodel database")
	cmd.Println("xormmodel database/table")
	cmd.Println("xormmodel database/table -c '/user/local/.../debug/config/mysql.yaml'")
	cmd.Println("xormmodel database/table -d 'user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4'")
}

func initXormDsn() string {
	// 1、 -d 获取dsn方式
	if len(xormDsn) > 0 {
		return xormDsn
	}
	// 2、-c 获取config 这个指定的优先级高于自动获取
	if len(xormConfigPath) == 0 {
		// 3、通过appDir 自动寻找 config文件
		xormConfigPath = xorm.ConfigFilePath()
	}

	dsn := getDnsFromConfig(xormConfigPath)
	if dsn != "" {
		xormModelFith := xorm.XormModelFilePath()
		if !xorm.FileExists(xormModelFith) {
			// 写一个默认配置文件
			xorm.SaveXormModelFile(xormModelFith, dsn)
		}
	}
	return dsn
}

func getDnsFromConfig(filePath string) string {
	Config := xorm.YamlFile{}
	if xorm.FileExists(filePath) {
		Config = xorm.ReadYamlFile(filePath)
	}
	dsn := Config.Data.Database.Source
	if len(dsn) == 0 {
		dsn = Config.Data.Mysql.Default.Dsn
	}
	fmt.Println(dsn)
	return dsn
}

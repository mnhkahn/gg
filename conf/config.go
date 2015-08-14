package conf

import (
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/viper"
)

var (
	AppConfig *GGConfig
	once      sync.Once
)

func init() {
	once = sync.Once{}
}

type GGConfig struct {
	sync.Once
	GOPATH            string
	GOOS              string
	AppName           string
	AppSuffix         string
	CurPath           string
	AppPath           string
	MainApplication   string
	RunDirectory      string
	RunUser           string
	LogDirectory      string
	SupervisorConf    string
	PackPaths         []string
	PackExcludePrefix []string
	PackExcludeSuffix []string
	PackExcludeRegexp []*regexp.Regexp
	PackFormat        string
}

func NewGGConfig() *GGConfig {
	if AppConfig == nil {
		once.Do(ParseConfig)
	}
	return AppConfig
}
func ParseConfig() {
	AppConfig = new(GGConfig)

	viper.SetConfigName("gg")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Fatal error config file: %s \n", err)
	}

	AppConfig.GOPATH = os.Getenv("GOPATH")
	AppConfig.GOOS = runtime.GOOS
	if v, found := syscall.Getenv("GOOS"); found {
		AppConfig.GOOS = v
	}
	if !strings.HasSuffix(AppConfig.GOPATH, "/") && !strings.HasSuffix(AppConfig.GOPATH, "\\") {
		AppConfig.GOPATH += "/"
	}

	AppConfig.CurPath, _ = os.Getwd()
	AppConfig.AppName = viper.GetString("AppName")
	if AppConfig.GOOS == "windows" {
		AppConfig.AppSuffix = ".exe"
	}
	AppConfig.AppPath = AppConfig.CurPath + "/" + AppConfig.AppName + ".tar.gz"
	AppConfig.MainApplication = viper.GetString("MainApplication")

	AppConfig.RunDirectory = viper.GetString("RunDirectory")
	AppConfig.RunUser = viper.GetString("RunUser")
	AppConfig.LogDirectory = viper.GetString("LogDirectory")
	AppConfig.SupervisorConf = viper.GetString("SupervisorConf")
	log.Println(AppConfig.PackPaths, "^^^^^^")
	AppConfig.PackPaths = append([]string{AppConfig.CurPath}, viper.GetStringSlice("PackPaths")...)
	log.Println(AppConfig.PackPaths, "^^^^^^")
	AppConfig.PackPaths = append(AppConfig.PackPaths, AppConfig.CurPath+"/"+AppConfig.AppName+AppConfig.AppSuffix)
	log.Println(AppConfig.PackPaths, "^^^^^^")
	AppConfig.PackFormat = "gzip"
	AppConfig.PackExcludePrefix = []string{".", AppConfig.AppPath, AppConfig.SupervisorConf}
	AppConfig.PackExcludeSuffix = []string{".go", ".DS_Store", ".tmp"}
}

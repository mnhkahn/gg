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
	HOME              string
	GOPATH            string
	GOOS              string
	AppName           string
	AppSuffix         string
	CurPath           string
	AppPath           string
	MainApplication   []string
	RunDirectory      string
	RunUser           string
	LogDirectory      string
	SupervisorConf    string
	PackPaths         []string
	PackExcludePrefix []string
	PackExcludeSuffix []string
	PackExcludeRegexp []*regexp.Regexp
	PackFormat        string

	IsGitPull     bool
	GitPullBranch string

	// IsNgrok   bool
	// NgrokPort string
}

func NewGGConfig() *GGConfig {
	if AppConfig == nil {
		once.Do(ParseConfig)
	}
	return AppConfig
}
func ParseConfig() {
	AppConfig = new(GGConfig)
	AppConfig.HOME = os.Getenv("HOME")
	AppConfig.GOPATH = os.Getenv("GOPATH")
	AppConfig.GOOS = runtime.GOOS
	if v, found := syscall.Getenv("GOOS"); found {
		AppConfig.GOOS = v
	}
	if !strings.HasSuffix(AppConfig.GOPATH, "/") && !strings.HasSuffix(AppConfig.GOPATH, "\\") {
		AppConfig.GOPATH += "/"
	}
	AppConfig.CurPath, _ = os.Getwd()
	if AppConfig.GOOS == "windows" {
		AppConfig.AppSuffix = ".exe"
	}
	AppConfig.GitPullBranch = "master"

	viper.SetConfigName("gg")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("There is no gg yaml config file.")
		for _, arg := range os.Args {
			if strings.HasSuffix(arg, ".go") {
				if AppConfig.AppName == "" {
					AppConfig.AppName = strings.TrimRight(arg, ".go")
				}
				AppConfig.MainApplication = append(AppConfig.MainApplication, arg)
			}
		}
	} else {
		AppConfig.AppName = viper.GetString("AppName")
		AppConfig.RunDirectory = strings.Replace(viper.GetString("RunDirectory"), "~", AppConfig.HOME, -1)
		AppConfig.RunUser = viper.GetString("RunUser")
		AppConfig.LogDirectory = strings.Replace(viper.GetString("LogDirectory"), "~", AppConfig.HOME, -1)
		AppConfig.SupervisorConf = viper.GetString("SupervisorConf")
		log.Println(AppConfig.PackPaths, "^^^^^^")
		AppConfig.PackPaths = append([]string{AppConfig.CurPath}, viper.GetStringSlice("PackPaths")...)
		log.Println(AppConfig.PackPaths, "^^^^^^")
		AppConfig.PackPaths = append(AppConfig.PackPaths, AppConfig.CurPath+"/"+AppConfig.AppName+AppConfig.AppSuffix)
		log.Println(AppConfig.PackPaths, "^^^^^^")
		AppConfig.MainApplication = viper.GetStringSlice("MainApplication")

		AppConfig.IsGitPull = viper.GetBool("Git")
		AppConfig.GitPullBranch = viper.GetString("GitPullBranch")

		// AppConfig.IsNgrok = viper.GetBool("Ngrok")
		// AppConfig.NgrokPort = viper.GetString("NgrokPort")
	}

	AppConfig.AppPath = AppConfig.CurPath + "/" + AppConfig.AppName + ".tar.gz"

	AppConfig.PackFormat = "gzip"
	AppConfig.PackExcludePrefix = []string{".", AppConfig.AppPath, AppConfig.SupervisorConf}
	AppConfig.PackExcludeSuffix = []string{".go", ".DS_Store", ".tmp"}
}

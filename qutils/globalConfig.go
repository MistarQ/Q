package utils

import (
	"Q/qlog"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

var (
	env  = flag.String("e", "prod", " -e, dev | test | prod")
	name = flag.String("m", "mmo_game", " -m, module name, for one module project it should be \".\" ")
)

type GlobalConfig struct {
	Host    string
	Port    int
	Name    string
	Version string

	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32

	LogDir        string
	LogFile       string
	LogDebugClose bool
}

var GlobalObject *GlobalConfig

func init() {

	flag.Parse()

	// default
	GlobalObject = &GlobalConfig{
		Name:             "Q",
		Version:          "default",
		Port:             8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           ".",
		LogFile:          "log",
		LogDebugClose:    false,
	}

	configPath := fmt.Sprintf("%s/config", *name)
	configFile := fmt.Sprintf("%s.yaml", *env)

	viper.SetConfigFile(configFile)
	viper.SetConfigName(*env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		qlog.Fatal("Fatal error getting config file: %s\n", err)
	}

	GlobalObject = &GlobalConfig{
		Name:             viper.GetString("Q.name"),
		Version:          viper.GetString("Q.version"),
		Port:             viper.GetInt("Q.port"),
		Host:             viper.GetString("Q.host"),
		MaxConn:          viper.GetInt("Q.maxConn"),
		MaxPackageSize:   viper.GetUint32("Q.maxPackageSize"),
		WorkerPoolSize:   viper.GetUint32("Q.workPoolSize"),
		MaxWorkerTaskLen: viper.GetUint32("Q.maxWorkerTaskLen"),
		MaxMsgChanLen:    viper.GetUint32("Q.maxMsgChanLen"),
		LogDir:           viper.GetString("log.logDir"),
		LogFile:          viper.GetString("log.logFile"),
		LogDebugClose:    viper.GetBool("log.logDebugClose"),
	}

	if GlobalObject.LogFile != "" {
		qlog.SetLogFile(GlobalObject.LogDir, GlobalObject.LogFile)
	}
	if GlobalObject.LogDebugClose == true {
		qlog.CloseDebug()
	}
}

func printConfig() {

	var topLine = `┌───────────────────────────────────────────────────┐`
	var borderLine = `│`
	var bottomLine = `└───────────────────────────────────────────────────┘`

	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/MistarQ               %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	t := reflect.TypeOf(*GlobalObject)
	v := reflect.ValueOf(*GlobalObject)
	fmt.Println("============================================")
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("%s : %v\n", t.Field(i).Name, v.Field(i).Interface())
	}
	fmt.Println("============================================")
}

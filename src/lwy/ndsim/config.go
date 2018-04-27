package ndsim

import (
	"errors"
	"fmt"

	"github.com/Unknwon/goconfig"
)

//Config struct define all config fields
type Config struct {
	UpdatePort             int
	QueryPort              int
	LogPath                string
	LogFileName            string
	LogLevel               int
	DataPath               string
	DocIDFrom              DocID
	DocIDEnd               DocID
	ForwardIndexFileName   string
	ForwardContentFileName string
	ReverseIndexFileName   string
	ReverseContentFileName string
	DictPath               string
}

//gConfig is global conf object
var gConfig = Config{
	UpdatePort:             8091,
	QueryPort:              8090,
	LogPath:                "./log",
	LogFileName:            "ndsim.log",
	LogLevel:               logLevelDebug,
	DataPath:               "./data",
	DocIDFrom:              100000000,
	DocIDEnd:               200000000,
	ForwardIndexFileName:   "forward.index",
	ForwardContentFileName: "forward.content",
	ReverseIndexFileName:   "reverse.index",
	ReverseContentFileName: "reverse.content",
	DictPath:               "D:/Go/data/dict/",
}

//LoadConfigFile load config file from configFilePath
func LoadConfigFile(configFilePath string) error {
	cfg, err := goconfig.LoadConfigFile(configFilePath)
	var errStr string
	if err != nil {
		fmt.Println(err)
		return err
	}
	{ //check for update port
		port, err := cfg.Int("global", "UpdatePort")
		if err != nil {
			fmt.Println(err)
			return err
		}
		if port <= 0 || port > 99999 {
			errStr = fmt.Sprintf("update port should in range[0~99999]")
			fmt.Println(errStr)
			return errors.New(errStr)
		}
		gConfig.UpdatePort = port
	}
	{ //check for query port
		port, err := cfg.Int("global", "QueryPort")
		if err != nil {
			fmt.Println(err)
			return err
		}
		if port <= 0 || port > 99999 {
			errStr = fmt.Sprintf("query port should in range[0~99999]")
			fmt.Println(errStr)
			return errors.New(errStr)
		}
		gConfig.QueryPort = port
	}

	{ //check for logPath , logFileName and loglevel
		logPath, err := cfg.GetValue("global", "LogPath")
		if err != nil {
			fmt.Println(err)
			return err
		}
		gConfig.LogPath = logPath

		logFileName, err := cfg.GetValue("global", "LogFileName")
		if err != nil {
			fmt.Println(err)
			return err
		}
		gConfig.LogFileName = logFileName

		logLevel, err := cfg.Int("global", "LogLevel")
		if err != nil {
			fmt.Println(err)
			return err
		}
		gConfig.LogLevel = logLevel
	}
	return nil
}

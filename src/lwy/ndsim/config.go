package ndsim

import (
	"fmt"

	"github.com/Unknwon/goconfig"
)

//Config struct define all config fields
type Config struct {
	UpdatePort int
	QueryPort  int
}

//GConfig is global conf object
var GConfig = Config{
	UpdatePort: 8091,
	QueryPort:  8090,
}

//LoadConfigFile load config file from configFilePath
func LoadConfigFile(configFilePath string) {
	cfg, err := goconfig.LoadConfigFile(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}

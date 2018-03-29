package configs

import (
	"github.com/aaa59891/mosi_demo_go/utils"
	"github.com/spf13/viper"
	"os"
	"path"
)

type config struct {
	Server   server
	Database database
}

type server struct {
	Port string
}

type database struct {
	DriveName    string
	Account      string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

var c *config = &config{}

func init() {
	dir := utils.GetProjectRoot()
	viper.SetConfigName("default")
	viper.AddConfigPath(path.Join(dir, "configs"))

	if err := viper.MergeInConfig(); err != nil {
		panic(err)
	}

	mosiGoEnv := os.Getenv("MOSI_GO")

	if len(mosiGoEnv) != 0 {
		viper.SetConfigName(mosiGoEnv)
	}

	if err := viper.MergeInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(c); err != nil {
		panic(err)
	}
}

func GetConfig() config {
	return *c
}

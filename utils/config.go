package utils

import (
	"github.com/spf13/viper"
	"log"
)
var (
	Con *Config
	Address string
	Port uint32
	ServerAddr string
	ServerPort uint32
)
type Config struct {
	V *viper.Viper
}
func InitConfig () *Config {
	Con := &Config{
		V: viper.New(),
	}
	//设置配置文件的名字
	Con.V.SetConfigName("config")
	//添加配置文件所在的路径,注意在Linux环境下%GOPATH要替换为$GOPATH
	Con.V.AddConfigPath("./")
	//设置配置文件类型
	Con.V.SetConfigType("yaml")
	if err := Con.V.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
	return Con
}

func init()  {
	//Con = InitConfig()
	//Address = GetAddress()
	//Port = GetPort()
	//ServerAddr = GetSerAddress()
	//ServerPort = GetSerPort()
	Address = "0.0.0.0"
	Port = 8989
	ServerAddr = "192.168.1.110"
	ServerPort = 8989
}

func GetAddress()string{
	return Con.V.GetString("address")
}

func GetPort()uint32{
	return Con.V.GetUint32("port")
}

func GetSerAddress()string{
	return Con.V.GetString("severaddr")
}

func GetSerPort()uint32{
	return Con.V.GetUint32("severport")
}
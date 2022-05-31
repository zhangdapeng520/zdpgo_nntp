package zdpgo_nntp

/*
@Time : 2022/5/31 10:05
@Author : 张大鹏
@File : config.go
@Software: Goland2021.3.1
@Description:
*/

type Config struct {
	Debug       bool     `yaml:"debug" json:"debug"`
	LogFilePath string   `yaml:"log_file_path" json:"log_file_path"`
	Server      HttpInfo `yaml:"server" json:"server"`
	Client      HttpInfo `yaml:"client" json:"client"`
}

type HttpInfo struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
}

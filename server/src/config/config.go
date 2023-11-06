/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2023-10-24 14:35:01
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-25 05:41:35
 * @FilePath: /1024/src/server/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)



const CONFIG_PATH = "server.yml"

type ServerConfig struct {
	Port int `yaml:"port"`
}

type CDNConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
}

type Config struct {
	Server ServerConfig
	CDN    CDNConfig
}

var config Config

func init() {
	configFile, err := os.Open(CONFIG_PATH)
	if err != nil {
		panic("Unable to locate server.yml")
	}
	data, err := io.ReadAll(configFile)
	if err != nil {
		panic(err)
	}
	// parse config
	err = yaml.UnmarshalStrict(data, &config)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return &config
}

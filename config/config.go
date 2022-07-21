package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const CONF = "\\resource\\conf.json"

var Global *Config

func init() {
	wd, _ := os.Getwd()
	Global = (&Config{}).load(wd + CONF)
}

type Config struct {
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	Protocol      string `json:"protocol"`
	MaxConn       int    `json:"maxConn"`
	MaxPacketSize int    `json:"maxPacketSize"`
	MaxPoolSize   int    `json:"maxPoolSize"`
	MaxQueLen     int    `json:"MaxQueLen"`
}

func New(ip string, port int, protoc string) *Config {
	return &Config{
		IP:       ip,
		Port:     port,
		Protocol: protoc,
	}
}

func (c *Config) load(path string) *Config {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("load config file failed: ", err)
		return nil
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return nil
	}

	return c
}

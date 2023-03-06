package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

const FILE_NAME = "config.json"

type ServerConfig struct {
	TcpAddr string `json:"tcpAddr"`
	TcpPort string `json:"tcpPort"`
	UdpAddr string `json:"udpAddr"`
	UdpPort string `json:"udpPort"`
}

type TransferConfig struct {
	SrcAddr string `json:"srcAddr"`
	SrcPort string `json:"srcPort"`
	DstAddr string `json:"dstAddr"`
	DstPort string `json:"dstPort"`
}

type ClientConfig struct {
	ServerIp   string `json:"ServerIp"`
	ServerPort string `json:"ServerPort"`
}

type Config struct {
	Server   ServerConfig   `json:"Server"`
	Transfer TransferConfig `json:"Transfer"`
	Client   ClientConfig   `json:"Client"`
}

func (cfg *Config) load() {
	cwd, _ := os.Getwd()

	file, err := os.OpenFile(filepath.Join(cwd, FILE_NAME), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(content) == 0 {
		cfg.Server = ServerConfig{}
		cfg.Transfer = TransferConfig{}
		cfg.Client = ClientConfig{}
		data, _ := json.Marshal(cfg)
		file.WriteString(string(data))
		return
	}

	json.Unmarshal([]byte(content), &cfg)
}

func (cfg *Config) save() {
	cwd, _ := os.Getwd()
	file, err := os.OpenFile(filepath.Join(cwd, FILE_NAME), os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	data, _ := json.Marshal(cfg)
	file.WriteString(string(data))
}

func (cfg *Config) GetData() Config {
	return *cfg
}

func (cfg *Config) SetData(newCfg *Config) {
	if reflect.DeepEqual(newCfg, cfg) {
		return
	}
	optionByte, _ := json.Marshal(newCfg)
	err := json.Unmarshal(optionByte, &cfg)
	if err != nil {
		return
	}
	cfg.save()
}

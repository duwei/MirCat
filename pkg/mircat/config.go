package mircat

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const FILE_NAME = "config.json"

// ServerConfig represents the configuration for the server.
type ServerConfig struct {
	// TcpAddr is the TCP IP address of the server.
	TcpAddr string `json:"tcpAddr"`
	// TcpPort is the TCP port of the server.
	TcpPort string `json:"tcpPort"`
	// UdpAddr is the UDP IP address of the server.
	UdpAddr string `json:"udpAddr"`
	// UdpPort is the UDP port of the server.
	UdpPort string `json:"udpPort"`
}

// TransferConfig represents the configuration for data transfer.
type TransferConfig struct {
	// SrcAddr is the source IP address for data transfer.
	SrcAddr string `json:"srcAddr"`
	// SrcPort is the source port for data transfer.
	SrcPort string `json:"srcPort"`
	// DstAddr is the destination IP address for data transfer.
	DstAddr string `json:"dstAddr"`
	// DstPort is the destination port for data transfer.
	DstPort string `json:"dstPort"`
}

// ClientConfig represents the configuration for the client.
type ClientConfig struct {
	// ServerIp is the IP address of the server that the client connects to.
	ServerIp string `json:"ServerIp"`
	// ServerPort is the port number of the server that the client connects to.
	ServerPort string `json:"ServerPort"`
}

// Config represents the overall configuration for the application.
type Config struct {
	// Server is the configuration for the server.
	Server ServerConfig `json:"Server"`
	// Transfer is the configuration for data transfer.
	Transfer TransferConfig `json:"Transfer"`
	// Client is the configuration for the client.
	Client ClientConfig `json:"Client"`
}

func NewConfig() *Config {
	c := &Config{}
	c.load()
	return c
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

// GetConfig retrieves the current configuration settings.
// Returns:
// - Config: the current configuration settings.
func (cfg *Config) GetConfig() Config {
	return *cfg
}

// SetConfig sets the current configuration settings.
// Parameters:
// - val (Config): the new configuration settings.
func (cfg *Config) SetConfig(newCfg *Config) (bool, error) {
	//if reflect.DeepEqual(newCfg, cfg) {
	//	return false, fmt.Errorf("same configure")
	//}
	optionByte, _ := json.Marshal(newCfg)
	err := json.Unmarshal(optionByte, &cfg)
	if err != nil {
		return false, err
	}
	return true, nil
}

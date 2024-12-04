package types

import (
	"fmt"
	"os"
)

type Config struct {
	Db  *DbConfig  `json:"db"`
	Rpc *RpcConfig `json:"rpc"`
}

type DbConfig struct {
	Path string `json:"path"`
}

type RpcConfig struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	TdRpc string `json:"td_rpc" mapstructure:"td_rpc"`
}

func DefaultConfig(idx, port, tdPort int) *Config {
	return &Config{
		Db: &DbConfig{Path: fmt.Sprintf("%s/.side-chain/%d/node_db", os.Getenv("HOME"), idx)},
		Rpc: &RpcConfig{
			Host:  "0.0.0.0",
			Port:  port,
			TdRpc: fmt.Sprintf("http://127.0.0.0:%d", tdPort),
		},
	}
}

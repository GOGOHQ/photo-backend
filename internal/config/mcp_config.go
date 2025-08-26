package config

import (
	"encoding/json"
	"os"
)

type MCPServer struct {
	Command  string   `json:"command"`
	Args     []string `json:"args"`
	Name     string   `json:"name"`
	BaseURL  string   `json:"baseUrl"`
	IsActive bool     `json:"isActive"`
}

type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

func LoadMCPConfig(path string) (*MCPConfig, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg MCPConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

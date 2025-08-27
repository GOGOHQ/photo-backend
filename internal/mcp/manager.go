package mcp

import (
	"fmt"
	"net/http"
	"strings"

	appcfg "github.com/huangqi/photo-backend/internal/config"
)

type ClientRegistry struct {
	Clients map[string]*MCPClient
}

func BuildTransportsFromMCPConfig(cfg *appcfg.MCPConfig, httpClient *http.Client) (*ClientRegistry, error) {
	if cfg == nil {
		return &ClientRegistry{Clients: map[string]*MCPClient{}}, nil
	}

	reg := &ClientRegistry{
		Clients: make(map[string]*MCPClient),
	}

	for key, srv := range cfg.MCPServers {
		if !srv.IsActive {
			continue
		}

		if srv.Command != "" {
			// 使用 stdio 传输创建 MCP 客户端
			client, err := NewStdioMCPClient(srv.Command, srv.Args, srv.Env)
			if err != nil {
				return nil, err
			}
			reg.Clients[key] = client

		} else if srv.BaseURL != "" {
			// 暂时不支持 HTTP，返回错误
			return nil, fmt.Errorf("HTTP MCP transport not yet supported")
		} else {
			continue
		}

		if srv.Name != "" {
			nameKey := strings.ToLower(strings.TrimSpace(srv.Name))
			if _, exists := reg.Clients[nameKey]; !exists {
				reg.Clients[nameKey] = reg.Clients[key]
			}
		}
	}
	return reg, nil
}

func (r *ClientRegistry) FindByKeyOrName(key string) *MCPClient {
	if r == nil {
		return nil
	}
	if t, ok := r.Clients[key]; ok {
		return t
	}
	lower := strings.ToLower(key)
	if t, ok := r.Clients[lower]; ok {
		return t
	}
	// fuzzy contains search on names/keys
	for k, t := range r.Clients {
		if strings.Contains(strings.ToLower(k), lower) {
			return t
		}
	}
	return nil
}

func (r *ClientRegistry) Close() {
	if r == nil {
		return
	}
	for k, t := range r.Clients {
		_ = t.Close()
		delete(r.Clients, k)
	}
}

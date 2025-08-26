package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/huangqi/photo-backend/internal/config"
	"github.com/huangqi/photo-backend/internal/db"
	"github.com/huangqi/photo-backend/internal/mcp"
	"github.com/huangqi/photo-backend/internal/server"
)

func main() {
	cfg := config.Load()

	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	var registry *mcp.ClientRegistry
	mcpPath := os.Getenv("MCP_CONFIG_PATH")
	if mcpPath == "" {
		if _, err := os.Stat("mcp.json"); err == nil {
			mcpPath = "mcp.json"
		}
	}
	if mcpPath != "" {
		if mcpCfg, err := config.LoadMCPConfig(mcpPath); err != nil {
			log.Printf("warn: failed to load MCP config: %v", err)
		} else {
			registry, err = mcp.BuildTransportsFromMCPConfig(mcpCfg, &http.Client{Timeout: 10 * time.Second})
			if err != nil {
				log.Fatalf("failed to build MCP clients: %v", err)
			}
		}
	}

	var xhsClient mcp.XHSClient
	if registry != nil {
		if client := registry.FindByKeyOrName("xhs"); client != nil {
			xhsClient = mcp.NewXHSClient(client)
		} else if client := registry.FindByKeyOrName("xiaohongshu"); client != nil {
			xhsClient = mcp.NewXHSClient(client)
		}
	}
	if xhsClient == nil {
		log.Fatalf("no XHS MCP client configured: set MCP_CONFIG_PATH")
	}

	var mapsClient mcp.MapsClient
	if registry != nil {
		if client := registry.FindByKeyOrName("maps"); client != nil {
			mapsClient = mcp.NewMapsClient(client)
		}
	}

	r := server.NewRouter(database, xhsClient, mapsClient)

	addr := fmt.Sprintf(":%d", cfg.Port)
	if envPort := os.Getenv("PORT"); envPort != "" {
		addr = ":" + envPort
	}
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

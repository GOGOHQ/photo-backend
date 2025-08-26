package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/huangqi/photo-backend/internal/config"
	"github.com/huangqi/photo-backend/internal/mcp"
)

func main() {
	var serverKey string
	flag.StringVar(&serverKey, "server", "xhs-local-stdio", "MCP server key to test")
	flag.Parse()

	fmt.Printf("Testing MCP server: %s\n", serverKey)

	// Load MCP config
	mcpPath := os.Getenv("MCP_CONFIG_PATH")
	if mcpPath == "" {
		if _, err := os.Stat("mcp.json"); err == nil {
			mcpPath = "mcp.json"
		}
	}
	if mcpPath == "" {
		log.Fatalf("No MCP config found, set MCP_CONFIG_PATH or ensure mcp.json exists")
	}

	fmt.Printf("Loading config from: %s\n", mcpPath)
	mcpCfg, err := config.LoadMCPConfig(mcpPath)
	if err != nil {
		log.Fatalf("Failed to load MCP config: %v", err)
	}

	// Find the specific server
	server, exists := mcpCfg.MCPServers[serverKey]
	if !exists {
		log.Fatalf("Server '%s' not found in config", serverKey)
	}

	fmt.Printf("Server config: %+v\n", server)

	if !server.IsActive {
		log.Fatalf("Server '%s' is not active", serverKey)
	}

	// Create MCP client directly
	fmt.Printf("Creating MCP client for command: %s %v\n", server.Command, server.Args)
	client, err := mcp.NewStdioMCPClient(server.Command, server.Args)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer client.Close()

	fmt.Println("✓ MCP client created successfully")

	// Test initialization
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("Initializing MCP client...")
	err = client.Initialize(ctx, "mcp-test", "1.0.0")
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Println("✓ MCP client initialized successfully")

	// Test ping
	fmt.Println("Testing ping...")
	err = client.Ping(ctx)
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}
	fmt.Println("✓ Ping successful")

	// List available tools
	fmt.Println("Listing available tools...")
	tools, err := client.ListTools(ctx)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("✓ Found %d tools:\n", len(tools.Tools))
	for i, tool := range tools.Tools {
		fmt.Printf("  %d. %s: %s\n", i+1, tool.Name, tool.Description)
	}

	// Step 1: 先进行登录
	fmt.Println("\n🔐 Step 1: 登录小红书账号...")
	loginResult, err := client.CallTool(ctx, "login", map[string]any{})
	if err != nil {
		log.Printf("⚠️  Login failed: %v", err)
		fmt.Println("   (继续测试其他功能...)")
	} else {
		fmt.Printf("✓ Login successful: %s\n", loginResult)
	}

	// Step 2: 测试搜索功能
	fmt.Println("\n🔍 Step 2: 测试搜索功能...")
	searchResult, err := client.CallTool(ctx, "search_notes", map[string]any{
		"keywords": "旅行",
		"limit":    3,
	})
	if err != nil {
		log.Printf("⚠️  Search failed: %v", err)
	} else {
		fmt.Printf("✓ Search successful: %s\n", searchResult)
	}

	// Step 3: 测试获取笔记内容（如果搜索成功）
	if searchResult != "" && searchResult != "请先登录小红书账号" {
		fmt.Println("\n📝 Step 3: 测试获取笔记内容...")
		// 这里可以从搜索结果中提取一个 URL 来测试
		// 暂时使用一个示例 URL
		noteResult, err := client.CallTool(ctx, "get_note_content", map[string]any{
			"url": "https://www.xiaohongshu.com/example",
		})
		if err != nil {
			log.Printf("⚠️  Get note content failed: %v", err)
		} else {
			fmt.Printf("✓ Get note content successful: %s\n", noteResult)
		}
	}

	fmt.Println("\n🎉 MCP server test completed successfully!")
}

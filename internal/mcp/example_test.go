package mcp

import (
	"context"
	"fmt"
	"log"
)

// ExampleNewStdioMCPClient 展示如何使用新的 stdio MCP 客户端
func ExampleNewStdioMCPClient() {
	// 创建基于 stdio 的 MCP 客户端
	client, err := NewStdioMCPClient("your-mcp-server", []string{"--config", "config.json"})
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 初始化客户端
	err = client.Initialize(ctx, "example-client", "1.0.0")
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// 发送 ping
	err = client.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	// 列出可用工具
	tools, err := client.ListTools(ctx)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Available tools: %d\n", len(tools.Tools))
	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// 调用工具
	result, err := client.CallTool(ctx, "example_tool", map[string]any{
		"param1": "value1",
		"param2": 42,
	})
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	fmt.Printf("Tool result: %s\n", result)
}

// ExampleNewXHSClient 展示如何使用新的 XHS 客户端
func ExampleNewXHSClient() {
	// 创建 MCP 客户端
	mcpClient, err := NewStdioMCPClient("xhs-mcp-server", []string{})
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// 创建 XHS 客户端
	xhsClient := NewXHSClient(mcpClient)
	defer xhsClient.Close()

	ctx := context.Background()

	// 搜索热门帖子
	posts, err := xhsClient.GetHotPosts(ctx, 5)
	if err != nil {
		log.Fatalf("Failed to get hot posts: %v", err)
	}

	fmt.Printf("Found %d hot posts\n", len(posts))
	for i, post := range posts {
		fmt.Printf("%d. %s - %s\n", i+1, post.Title, post.PostURL)
	}

	// 按关键词搜索
	posts, err = xhsClient.GetPostsByKeyword(ctx, "旅行", 10)
	if err != nil {
		log.Fatalf("Failed to search posts: %v", err)
	}

	fmt.Printf("Found %d posts for '旅行'\n", len(posts))
}

// ExampleNewMapsClient 展示如何使用新的 Maps 客户端
func ExampleNewMapsClient() {
	// 创建 MCP 客户端
	mcpClient, err := NewStdioMCPClient("maps-mcp-server", []string{})
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// 创建 Maps 客户端
	mapsClient := NewMapsClient(mcpClient)
	defer mapsClient.Close()

	ctx := context.Background()

	// 获取附近景点
	attractions, err := mapsClient.GetNearbyAttractions(ctx, 39.9042, 116.4074, 5.0)
	if err != nil {
		log.Fatalf("Failed to get nearby attractions: %v", err)
	}

	fmt.Printf("Found %d nearby attractions\n", len(attractions))
	for i, attraction := range attractions {
		fmt.Printf("%d. %s (%.4f, %.4f) - %s\n", 
			i+1, attraction.Name, attraction.Latitude, attraction.Longitude, attraction.Address)
	}
}

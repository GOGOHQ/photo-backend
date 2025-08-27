package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/huangqi/photo-backend/internal/config"
	"github.com/huangqi/photo-backend/internal/mcp"
)

func main() {
	var keyword string
	var limit int
	var skipLogin bool
	flag.StringVar(&keyword, "q", "热门", "keyword to search")
	flag.IntVar(&limit, "limit", 5, "number of results")
	flag.BoolVar(&skipLogin, "skip-login", false, "skip login step")
	flag.Parse()

	// Build clients registry from MCP_CONFIG_PATH or fallback ./mcp.json
	var registry *mcp.ClientRegistry

	mcpPath := "/Users/huangqi/code/photo-backend/mcp.json"

	if mcpCfg, err := config.LoadMCPConfig(mcpPath); err != nil {
		log.Fatalf("Failed to load MCP config: %v", err)
	} else {
		reg, err := mcp.BuildTransportsFromMCPConfig(mcpCfg, &http.Client{Timeout: 30 * time.Second})
		if err != nil {
			log.Fatalf("Failed to build MCP clients: %v", err)
		}
		registry = reg
	}

	// Resolve XHS client
	var xhsClient mcp.XHSClient
	if registry != nil {
		if client := registry.FindByKeyOrName("xhs"); client != nil {
			xhsClient = mcp.NewXHSClient(client)
		} else if client := registry.FindByKeyOrName("xiaohongshu"); client != nil {
			xhsClient = mcp.NewXHSClient(client)
		} else if client := registry.FindByKeyOrName("xhs-local-stdio"); client != nil {
			xhsClient = mcp.NewXHSClient(client)
		}
	}
	if xhsClient == nil {
		log.Fatalf("No XHS MCP client configured: set MCP_CONFIG_PATH")
	}
	defer xhsClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// 先初始化 MCP 客户端
	fmt.Println("Initializing MCP client...")
	err := xhsClient.Initialize(ctx, "photo-backend-debug", "0.1.0")
	if err != nil {
		log.Fatalf("Failed to initialize MCP client: %v", err)
	}
	fmt.Println("✓ MCP client initialized successfully")

	// Step 1: 登录小红书账号（除非跳过）
	if !skipLogin {
		fmt.Println("\n🔐 Step 1: 登录小红书账号...")
		loginResult, err := xhsClient.CallTool(ctx, "login", map[string]any{})
		if err != nil {
			log.Printf("⚠️  Login failed: %v", err)
			fmt.Println("   (继续测试搜索功能...)")
		} else {
			fmt.Printf("✓ Login successful: %s\n", loginResult)
		}
	} else {
		fmt.Println("\n⏭️  Skipping login step...")
	}

	// Step 2: 搜索笔记
	fmt.Printf("\n🔍 Step 2: 搜索关键词 '%s' (限制 %d 条结果)...\n", keyword, limit)
	posts, err := xhsClient.GetPostsByKeyword(ctx, keyword, limit)
	if err != nil {
		log.Fatalf("Search error: %v", err)
	}

	fmt.Printf("\n✓ Found %d posts:\n", len(posts))
	if len(posts) > 0 {
		for i, post := range posts {
			fmt.Printf("%d. %s\n   %s\n", i+1, post.Title, post.PostURL)
		}
	} else {
		fmt.Println("   (没有找到相关帖子，可能需要先登录)")
	}

	// Step 3: 如果找到帖子，尝试获取第一个帖子的内容
	if len(posts) > 0 && !skipLogin {
		fmt.Printf("\n📝 Step 3: 获取第一个帖子的内容...\n")
		noteResult, err := xhsClient.CallTool(ctx, "get_note_content", map[string]any{
			"url": posts[0].PostURL,
		})
		if err != nil {
			log.Printf("⚠️  Get note content failed: %v", err)
		} else {
			fmt.Printf("✓ Note content: %s\n", noteResult)
		}
	}

	fmt.Println("\n🎉 XHS debug test completed!")
}

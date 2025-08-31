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
	var testType string
	flag.StringVar(&testType, "test", "geocode", "test type: geocode, reverse-geocode, search-places, directions, weather, ip-location, traffic")
	flag.Parse()

	fmt.Printf("Testing Baidu Maps MCP server: %s\n", testType)

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

	// Find the baidu-maps server
	server, exists := mcpCfg.MCPServers["baidu-maps"]
	if !exists {
		log.Fatalf("Baidu Maps server not found in config")
	}

	fmt.Printf("Server config: %+v\n", server)

	if !server.IsActive {
		log.Fatalf("Baidu Maps server is not active")
	}

	// Create MCP client
	fmt.Printf("Creating MCP client for command: %s %v\n", server.Command, server.Args)
	client, err := mcp.NewStdioMCPClient(server.Command, server.Args, server.Env)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer client.Close()

	fmt.Println("✓ MCP client created successfully")

	// Create Baidu Maps client
	baiduMapsClient := mcp.NewBaiduMapsClient(client)

	// Test initialization
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("Initializing MCP client...")
	err = baiduMapsClient.Initialize(ctx, "baidu-maps-test", "1.0.0")
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Println("✓ MCP client initialized successfully")

	switch testType {
	case "search-places":
		testSearchPlaces(ctx, baiduMapsClient)
	case "weather":
		testWeather(ctx, baiduMapsClient)
	default:
		fmt.Println("Run specific tests with -test flag, e.g., -test search-places or -test weather")
	}
}

func testSearchPlaces(ctx context.Context, client mcp.BaiduMapsClient) {
	query := "咖啡"
	tag := "美食"
	region := "上海"
	location := "31.2304,121.4737"
	radius := "2000"
	language := "zh-CN"
	isChina := "true"

	places, err := client.SearchPlaces(ctx, query, tag, region, location, radius, language, isChina)
	if err != nil {
		log.Printf("SearchPlaces failed: %v", err)
		return
	}
	fmt.Printf("Found %d places\n", len(places))
	for i, p := range places {
		fmt.Printf("%d. %s - %s\n", i+1, p.Name, p.Address)
	}
}

func testWeather(ctx context.Context, client mcp.BaiduMapsClient) {
	location := "116.30684538228411,40.057737278172176"
	districtID := ""
	isChina := "true"
	res, err := client.GetWeather(ctx, location, districtID, isChina)
	if err != nil {
		log.Printf("GetWeather failed: %v", err)
		return
	}
	fmt.Printf("Weather ok: %+v\n", res.Result.Now)
}

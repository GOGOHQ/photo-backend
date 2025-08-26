package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPClient 使用官方 mcp-go 库的客户端包装器
type MCPClient struct {
	client client.MCPClient
}

// NewStdioMCPClient 创建基于 stdio 的 MCP 客户端
func NewStdioMCPClient(command string, args []string) (*MCPClient, error) {
	// mcp-go 的 NewStdioMCPClient 签名是 (command, env, args...)
	// 我们传递空的环境变量，然后展开 args
	client, err := client.NewStdioMCPClient(command, []string{}, args...)
	if err != nil {
		return nil, err
	}
	return &MCPClient{client: client}, nil
}

// Initialize 初始化 MCP 客户端
func (c *MCPClient) Initialize(ctx context.Context, name, version string) error {
	if c.client == nil {
		return fmt.Errorf("client not initialized")
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    name,
		Version: version,
	}

	_, err := c.client.Initialize(ctx, initRequest)
	return err
}

// Ping 发送 ping 请求
func (c *MCPClient) Ping(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("client not initialized")
	}

	err := c.client.Ping(ctx)
	return err
}

// ListTools 获取工具列表
func (c *MCPClient) ListTools(ctx context.Context) (*ToolsListResponse, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.client.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, err
	}

	response := &ToolsListResponse{
		Tools: make([]Tool, len(tools.Tools)),
	}

	for i, tool := range tools.Tools {
		response.Tools[i] = Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		}
	}

	return response, nil
}

// CallTool 调用工具
func (c *MCPClient) CallTool(ctx context.Context, name string, args map[string]any) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("client not initialized")
	}

	toolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	toolRequest.Params.Name = name
	toolRequest.Params.Arguments = args

	result, err := c.client.CallTool(ctx, toolRequest)
	if err != nil {
		return "", err
	}

	// 处理返回结果
	if len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(mcp.TextContent); ok {
			return textContent.Text, nil
		}
		// 如果不是文本内容，尝试转换为字符串
		return fmt.Sprintf("%v", result.Content[0]), nil
	}

	return "", fmt.Errorf("empty MCP tools/call result")
}

// Close 关闭客户端连接
func (c *MCPClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Tool 表示一个 MCP 工具
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	InputSchema any    `json:"input_schema,omitempty"`
}

// ToolsListResponse 工具列表响应
type ToolsListResponse struct {
	Tools []Tool `json:"tools"`
}

package mcp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type XHSPost struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Likes   int    `json:"likes"`
	Excerpt string `json:"excerpt"`
	PostURL string `json:"post_url"`
}

type XHSClient interface {
	Initialize(ctx context.Context, name, version string) error
	GetHotPosts(ctx context.Context, limit int) ([]XHSPost, error)
	GetPostsByKeyword(ctx context.Context, keyword string, limit int) ([]XHSPost, error)
	CallTool(ctx context.Context, name string, args map[string]any) (string, error)
	Close() error
}

type xhsClient struct {
	mcp *MCPClient
}

// NewXHSClient 直接使用 MCP 客户端创建 XHS 客户端
func NewXHSClient(mcpClient *MCPClient) XHSClient {
	return &xhsClient{mcp: mcpClient}
}

func (c *xhsClient) Initialize(ctx context.Context, name, version string) error {
	if c.mcp != nil {
		return c.mcp.Initialize(ctx, name, version)
	}
	return nil
}

func (c *xhsClient) GetHotPosts(ctx context.Context, limit int) ([]XHSPost, error) {
	return c.GetPostsByKeyword(ctx, "热门", limit)
}

func (c *xhsClient) GetPostsByKeyword(ctx context.Context, keyword string, limit int) ([]XHSPost, error) {
	out, err := c.mcp.CallTool(ctx, "search_notes", map[string]any{
		"keywords": keyword,
		"limit":    limit,
	})
	if err != nil {
		return nil, err
	}
	return parseSearchNotesOutput(out), nil
}

// CallTool 直接调用 MCP 工具
func (c *xhsClient) CallTool(ctx context.Context, name string, args map[string]any) (string, error) {
	if c.mcp != nil {
		return c.mcp.CallTool(ctx, name, args)
	}
	return "", fmt.Errorf("MCP client not available")
}

func (c *xhsClient) Close() error {
	if c.mcp != nil {
		return c.mcp.Close()
	}
	return nil
}

var (
	reItem = regexp.MustCompile(`(?m)^\s*\d+\.\s*(.+?)\s*\n\s*链接:\s*(\S+)\s*$`)
)

func parseSearchNotesOutput(text string) []XHSPost {
	var posts []XHSPost
	matches := reItem.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		if len(m) >= 3 {
			posts = append(posts, XHSPost{
				Title:   strings.TrimSpace(m[1]),
				PostURL: strings.TrimSpace(m[2]),
			})
		}
	}
	return posts
}

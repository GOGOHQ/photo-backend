package mcp

import (
	"context"
	"encoding/json"
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
	reItemLinkLabel    = regexp.MustCompile(`(?m)^\s*\d+\.\s*(.+?)\s*\n\s*链接:\s*(\S+)\s*$`)
	reItemTitleThenURL = regexp.MustCompile(`(?m)^\s*\d+\.\s*(.+?)\s*\n\s*(https?://\S+)\s*$`)
	reAnyXHSURL        = regexp.MustCompile(`(?m)(https?://(?:www\.)?xiaohongshu\.com/\S+)`)
	reDetailPrefer     = regexp.MustCompile(`(?i)^https?://(?:www\.)?xiaohongshu\.com/(?:explore|discovery/item)/[A-Za-z0-9]+`)
	reSearchResultID   = regexp.MustCompile(`(?i)^https?://(?:www\.)?xiaohongshu\.com/search_result/([A-Za-z0-9]+)`)
)

func parseSearchNotesOutput(text string) []XHSPost {
	var posts []XHSPost

	// 1) 标准格式：有“链接:”标签
	if matches := reItemLinkLabel.FindAllStringSubmatch(text, -1); len(matches) > 0 {
		for _, m := range matches {
			if len(m) >= 3 {
				posts = append(posts, XHSPost{Title: strings.TrimSpace(m[1]), PostURL: strings.TrimSpace(m[2])})
			}
		}
		return uniquePosts(posts)
	}

	// 2) 兼容：标题下一行直接是 URL
	if matches := reItemTitleThenURL.FindAllStringSubmatch(text, -1); len(matches) > 0 {
		for _, m := range matches {
			if len(m) >= 3 {
				posts = append(posts, XHSPost{Title: strings.TrimSpace(m[1]), PostURL: strings.TrimSpace(m[2])})
			}
		}
		return uniquePosts(posts)
	}

	// 3) 兼容 JSON 数组输出
	var any interface{}
	if json.Unmarshal([]byte(text), &any) == nil {
		switch v := any.(type) {
		case []interface{}:
			for _, elem := range v {
				if obj, ok := elem.(map[string]interface{}); ok {
					url := firstNonEmptyString(obj["url"], obj["link"], obj["post_url"], obj["href"])
					title := firstNonEmptyString(obj["title"], obj["name"], obj["desc"], obj["excerpt"])
					if url != "" {
						posts = append(posts, XHSPost{Title: title, PostURL: url})
					}
				}
			}
			if len(posts) > 0 {
				return uniquePosts(posts)
			}
		}
	}

	// 4) 兜底：抓取所有小红书 URL
	seen := map[string]struct{}{}
	for _, m := range reAnyXHSURL.FindAllStringSubmatch(text, -1) {
		if len(m) >= 2 {
			url := strings.TrimSpace(m[1])
			if url == "" {
				continue
			}
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			posts = append(posts, XHSPost{PostURL: url})
		}
	}
	return posts
}

func uniquePosts(in []XHSPost) []XHSPost {
	seen := map[string]struct{}{}
	out := make([]XHSPost, 0, len(in))
	for _, p := range in {
		url := strings.TrimSpace(p.PostURL)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		out = append(out, p)
	}
	return out
}

func firstNonEmptyString(values ...interface{}) string {
	for _, v := range values {
		if s, ok := v.(string); ok {
			s = strings.TrimSpace(s)
			if s != "" {
				return s
			}
		}
	}
	return ""
}

// NormalizeXHSURL 将链接规范化为帖子详情页链接（/explore/<id>）。
func NormalizeXHSURL(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if m := reDetailPrefer.FindStringSubmatch(s); len(m) >= 1 {
		// 已是详情页或可提取到 id
		if len(m) == 2 {
			return "https://www.xiaohongshu.com/explore/" + strings.TrimSpace(m[1])
		}
		return s
	}
	if m := reSearchResultID.FindStringSubmatch(s); len(m) >= 2 {
		return "https://www.xiaohongshu.com/explore/" + strings.TrimSpace(m[1])
	}
	if strings.Contains(s, "/search_result") {
		return ""
	}
	return s
}

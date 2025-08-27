package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huangqi/photo-backend/internal/mcp"
)

type XHSHandler struct {
	Client      mcp.XHSClient
	initialized bool
	loggedIn    bool
}

func NewXHSHandler(client mcp.XHSClient) *XHSHandler {
	return &XHSHandler{Client: client}
}

// ensureInit initializes MCP client once per process
func (h *XHSHandler) ensureInit(c *gin.Context) {
	if h.initialized {
		return
	}
	// Initialize MCP client
	_ = h.Client.Initialize(c.Request.Context(), "photo-backend-server", "1.0.0")
	h.initialized = true
	if !h.loggedIn {
		_, _ = h.Client.CallTool(c.Request.Context(), "login", map[string]any{})
		h.loggedIn = true
	}
}

// GetHot retains backward compatibility. When query param q is present, it searches.
// Response: { data: posts }
func (h *XHSHandler) GetHot(c *gin.Context) {
	h.ensureInit(c)

	skipLogin := c.DefaultQuery("skip_login", "true") == "true"
	if !skipLogin {
		// Best-effort login; ignore error but surface if needed
		_, _ = h.Client.CallTool(c.Request.Context(), "login", map[string]any{})
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if q := c.Query("q"); q != "" {
		posts, err := h.Client.GetPostsByKeyword(c.Request.Context(), q, limit)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": posts})
		return
	}
	posts, err := h.Client.GetHotPosts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h *XHSHandler) SearchLinks(c *gin.Context) {
	h.ensureInit(c)

	keyword := c.DefaultQuery("q", "热门")
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	skipLogin := c.DefaultQuery("skip_login", "false") == "true"
	if !skipLogin {
		_, _ = h.Client.CallTool(c.Request.Context(), "login", map[string]any{})
	}

	posts, err := h.Client.GetPostsByKeyword(c.Request.Context(), keyword, limit)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	links := make([]string, 0, len(posts))
	for _, p := range posts {
		if p.PostURL != "" {
			links = append(links, p.PostURL)
		}
	}
	c.JSON(http.StatusOK, gin.H{"links": links})
}

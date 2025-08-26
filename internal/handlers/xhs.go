package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huangqi/photo-backend/internal/mcp"
)

type XHSHandler struct {
	Client mcp.XHSClient
}

func NewXHSHandler(client mcp.XHSClient) *XHSHandler {
	return &XHSHandler{Client: client}
}

func (h *XHSHandler) GetHot(c *gin.Context) {
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

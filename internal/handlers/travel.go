package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huangqi/photo-backend/internal/mcp"
)

type TravelHandler struct {
	Maps mcp.MapsClient
}

func NewTravelHandler(client mcp.MapsClient) *TravelHandler {
	return &TravelHandler{Maps: client}
}

func (h *TravelHandler) GetNearby(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radiusKm, err := strconv.ParseFloat(c.DefaultQuery("radius_km", "5"), 64)
	if err != nil || radiusKm <= 0 {
		radiusKm = 5
	}
	if c.Query("lat") == "" || c.Query("lng") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng are required"})
		return
	}
	items, err := h.Maps.GetNearbyAttractions(c.Request.Context(), lat, lng, radiusKm)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

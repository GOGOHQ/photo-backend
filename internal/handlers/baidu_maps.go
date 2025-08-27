package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/huangqi/photo-backend/internal/mcp"
)

type BaiduMapsHandler struct {
	Client      mcp.BaiduMapsClient
	initialized bool
}

func NewBaiduMapsHandler(client mcp.BaiduMapsClient) *BaiduMapsHandler {
	return &BaiduMapsHandler{Client: client}
}

// ensureInit 确保客户端已初始化
func (h *BaiduMapsHandler) ensureInit(c *gin.Context) {
	if h.initialized {
		return
	}
	_ = h.Client.Initialize(c.Request.Context(), "photo-backend-server", "1.0.0")
	h.initialized = true
}

// Geocode 地理编码
func (h *BaiduMapsHandler) Geocode(c *gin.Context) {
	h.ensureInit(c)

	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	result, err := h.Client.Geocode(c.Request.Context(), address)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ReverseGeocode 逆地理编码
func (h *BaiduMapsHandler) ReverseGeocode(c *gin.Context) {
	h.ensureInit(c)

	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat parameter"})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lng parameter"})
		return
	}

	result, err := h.Client.ReverseGeocode(c.Request.Context(), lat, lng)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// SearchPlaces 搜索地点（透传 query、tag、region、location、radius、language、is_china）
func (h *BaiduMapsHandler) SearchPlaces(c *gin.Context) {
	h.ensureInit(c)

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q parameter is required"})
		return
	}
	tag := c.DefaultQuery("tag", "")
	region := c.DefaultQuery("region", "全国")
	location := c.DefaultQuery("location", "")
	radiusStr := c.DefaultQuery("radius", "")
	language := c.DefaultQuery("language", "")
	isChina := c.DefaultQuery("is_china", "true")

	var radius int
	if radiusStr != "" {
		v, err := strconv.Atoi(radiusStr)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius parameter"})
			return
		}
		radius = v
	}

	places, err := h.Client.SearchPlaces(c.Request.Context(), query, tag, region, location, radius, language, isChina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": places})
}

// GetDirections 路线规划
func (h *BaiduMapsHandler) GetDirections(c *gin.Context) {
	h.ensureInit(c)

	origin := c.Query("origin")
	destination := c.Query("destination")
	mode := c.DefaultQuery("mode", "driving")

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin and destination parameters are required"})
		return
	}

	result, err := h.Client.GetDirections(c.Request.Context(), origin, destination, mode)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetWeather 天气查询（支持 city 或 lat/lng，并透传 district_id、is_china）
func (h *BaiduMapsHandler) GetWeather(c *gin.Context) {
	h.ensureInit(c)

	city := c.Query("city")
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	districtID := c.DefaultQuery("district_id", "")
	isChina := c.DefaultQuery("is_china", "true")

	location := ""
	if latStr != "" && lngStr != "" {
		lat, err1 := strconv.ParseFloat(latStr, 64)
		lng, err2 := strconv.ParseFloat(lngStr, 64)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat or lng"})
			return
		}
		// 使用 "lat,lng" 顺序
		location = strconv.FormatFloat(lat, 'f', -1, 64) + "," + strconv.FormatFloat(lng, 'f', -1, 64)
	} else if city != "" {
		location = city
	}

	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city or (lat,lng) parameter is required"})
		return
	}

	result, err := h.Client.GetWeather(c.Request.Context(), location, districtID, isChina)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetIPLocation IP定位
func (h *BaiduMapsHandler) GetIPLocation(c *gin.Context) {
	h.ensureInit(c)

	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ip parameter is required"})
		return
	}

	result, err := h.Client.GetIPLocation(c.Request.Context(), ip)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetTraffic 路况查询
func (h *BaiduMapsHandler) GetTraffic(c *gin.Context) {
	h.ensureInit(c)

	road := c.Query("road")
	city := c.Query("city")

	if road == "" || city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "road and city parameters are required"})
		return
	}

	result, err := h.Client.GetTraffic(c.Request.Context(), road, city)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

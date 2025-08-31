package server

import (
	"github.com/gin-gonic/gin"
	"github.com/huangqi/photo-backend/internal/handlers"
	"github.com/huangqi/photo-backend/internal/mcp"
)

func NewRouter(maps mcp.MapsClient, baiduMaps mcp.BaiduMapsClient) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	//xhsHandler := handlers.NewXHSHandler(xhs)
	travelHandler := handlers.NewTravelHandler(maps)
	baiduMapsHandler := handlers.NewBaiduMapsHandler(baiduMaps)

	api := r.Group("/api")
	{
		// xhsGroup := api.Group("/xhs")
		// xhsGroup.GET("/hot", xhsHandler.GetHot)
		// xhsGroup.GET("/search", xhsHandler.SearchLinks)

		travelGroup := api.Group("/travel")
		travelGroup.GET("/nearby", travelHandler.GetNearby)

		baiduMapsGroup := api.Group("/baidu-maps")
		baiduMapsGroup.GET("/geocode", baiduMapsHandler.Geocode)
		baiduMapsGroup.GET("/reverse-geocode", baiduMapsHandler.ReverseGeocode)
		baiduMapsGroup.GET("/search-places", baiduMapsHandler.SearchPlaces)
		baiduMapsGroup.GET("/directions", baiduMapsHandler.GetDirections)
		baiduMapsGroup.GET("/weather", baiduMapsHandler.GetWeather)
		baiduMapsGroup.GET("/ip-location", baiduMapsHandler.GetIPLocation)
		baiduMapsGroup.GET("/traffic", baiduMapsHandler.GetTraffic)
	}

	return r
}

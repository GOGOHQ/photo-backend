package server

import (
	"github.com/gin-gonic/gin"
	"github.com/huangqi/photo-backend/internal/handlers"
	"github.com/huangqi/photo-backend/internal/mcp"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, xhs mcp.XHSClient, maps mcp.MapsClient) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	xhsHandler := handlers.NewXHSHandler(xhs)
	travelHandler := handlers.NewTravelHandler(maps)

	api := r.Group("/api")
	{
		xhsGroup := api.Group("/xhs")
		xhsGroup.GET("/hot", xhsHandler.GetHot)

		travelGroup := api.Group("/travel")
		travelGroup.GET("/nearby", travelHandler.GetNearby)
	}

	return r
}

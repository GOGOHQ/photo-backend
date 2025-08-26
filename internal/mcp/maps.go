package mcp

import (
	"context"
	"encoding/json"
	"fmt"
)

type Attraction struct {
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Address    string  `json:"address"`
	DistanceKm float64 `json:"distance_km"`
	PlaceID    string  `json:"place_id"`
}

type MapsClient interface {
	GetNearbyAttractions(ctx context.Context, lat, lng float64, radiusKm float64) ([]Attraction, error)
	Close() error
}

type mapsClient struct {
	mcp *MCPClient
}

func NewMapsClient(mcpClient *MCPClient) MapsClient {
	return &mapsClient{mcp: mcpClient}
}

func (c *mapsClient) GetNearbyAttractions(ctx context.Context, lat, lng float64, radiusKm float64) ([]Attraction, error) {
	params := map[string]any{
		"lat":       lat,
		"lng":       lng,
		"radius_km": radiusKm,
	}

	result, err := c.mcp.CallTool(ctx, "maps.nearby", params)
	if err != nil {
		return nil, err
	}

	// 解析返回的 JSON 结果
	var items []Attraction
	if err := json.Unmarshal([]byte(result), &items); err != nil {
		return nil, fmt.Errorf("failed to parse attractions: %w", err)
	}

	return items, nil
}

func (c *mapsClient) Close() error {
	if c.mcp != nil {
		return c.mcp.Close()
	}
	return nil
}

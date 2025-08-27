package mcp

import (
	"context"
	"encoding/json"
	"fmt"
)

// BaiduMapsClient 百度地图客户端接口
type BaiduMapsClient interface {
	Initialize(ctx context.Context, name, version string) error
	Geocode(ctx context.Context, address string) (*GeocodeResult, error)
	ReverseGeocode(ctx context.Context, lat, lng float64) (*ReverseGeocodeResult, error)
	SearchPlaces(ctx context.Context, query, tag, region, location string, radius int, language, isChina string) ([]Place, error)
	GetDirections(ctx context.Context, origin, destination string, mode string) (*DirectionsResult, error)
	GetWeather(ctx context.Context, location string, districtID string, isChina string) (*WeatherResult, error)
	GetIPLocation(ctx context.Context, ip string) (*IPLocationResult, error)
	GetTraffic(ctx context.Context, road string, city string) (*TrafficResult, error)
	Close() error
}

type baiduMapsClient struct {
	mcp *MCPClient
}

// NewBaiduMapsClient 创建百度地图客户端
func NewBaiduMapsClient(mcpClient *MCPClient) BaiduMapsClient {
	return &baiduMapsClient{mcp: mcpClient}
}

func (c *baiduMapsClient) Initialize(ctx context.Context, name, version string) error {
	if c.mcp != nil {
		return c.mcp.Initialize(ctx, name, version)
	}
	return nil
}

// Geocode 地理编码
func (c *baiduMapsClient) Geocode(ctx context.Context, address string) (*GeocodeResult, error) {
	result, err := c.mcp.CallTool(ctx, "map_geocode", map[string]any{
		"address": address,
	})
	if err != nil {
		return nil, err
	}

	var geocodeResult GeocodeResult
	if err := json.Unmarshal([]byte(result), &geocodeResult); err != nil {
		return nil, fmt.Errorf("failed to parse geocode result: %w", err)
	}
	// 检查 status 字段，非 0 视为失败
	if geocodeResult.Status != 0 {
		return nil, fmt.Errorf("geocode API error: status=%d", geocodeResult.Status)
	}

	return &geocodeResult, nil
}

// ReverseGeocode 逆地理编码（兼容 lat/lng 与 latitude/longitude）
func (c *baiduMapsClient) ReverseGeocode(ctx context.Context, lat, lng float64) (*ReverseGeocodeResult, error) {
	result, err := c.mcp.CallTool(ctx, "map_reverse_geocode", map[string]any{
		"lat":       lat,
		"lng":       lng,
		"latitude":  lat,
		"longitude": lng,
	})
	if err != nil {
		return nil, err
	}

	var reverseGeocodeResult ReverseGeocodeResult
	if err := json.Unmarshal([]byte(result), &reverseGeocodeResult); err != nil {
		return nil, fmt.Errorf("failed to parse reverse geocode result: %w", err)
	}

	return &reverseGeocodeResult, nil
}

// SearchPlaces 搜索地点，对齐 MCP Server 参数（radius 为整数）
func (c *baiduMapsClient) SearchPlaces(ctx context.Context, query, tag, region, location string, radius int, language, isChina string) ([]Place, error) {
	args := map[string]any{
		"query":    query,
		"tag":      tag,
		"region":   region,
		"location": location,
		"radius":   radius,
		"language": language,
		"is_china": isChina,
	}
	result, err := c.mcp.CallTool(ctx, "map_search_places", args)
	fmt.Printf("result: %v\n", result)
	if err != nil {
		return nil, err
	}

	var places []Place
	if err := json.Unmarshal([]byte(result), &places); err != nil {
		return nil, fmt.Errorf("failed to parse places result: %w", err)
	}

	return places, nil
}

// GetDirections 路线规划
func (c *baiduMapsClient) GetDirections(ctx context.Context, origin, destination, mode string) (*DirectionsResult, error) {
	result, err := c.mcp.CallTool(ctx, "map_directions", map[string]any{
		"origin":      origin,
		"destination": destination,
		"mode":        mode,
	})
	if err != nil {
		return nil, err
	}

	var directionsResult DirectionsResult
	if err := json.Unmarshal([]byte(result), &directionsResult); err != nil {
		return nil, fmt.Errorf("failed to parse directions result: %w", err)
	}

	return &directionsResult, nil
}

// GetWeather 天气查询（传递 location、可选 district_id 与 is_china）
func (c *baiduMapsClient) GetWeather(ctx context.Context, location string, districtID string, isChina string) (*WeatherResult, error) {
	args := map[string]any{
		"location": location,
	}
	if districtID != "" {
		args["district_id"] = districtID
	}
	if isChina != "" {
		args["is_china"] = isChina
	}

	result, err := c.mcp.CallTool(ctx, "map_weather", args)
	if err != nil {
		return nil, err
	}

	fmt.Printf("result: %v\n", result)
	var weatherResult WeatherResult
	if err := json.Unmarshal([]byte(result), &weatherResult); err != nil {
		raw := result
		if len(raw) > 200 {
			raw = raw[:200] + "..."
		}
		return nil, fmt.Errorf("failed to parse weather result: %w; raw response: %s", err, raw)
	}
	// 若包含状态码且非 0，返回错误
	if weatherResult.Status != 0 {
		return nil, fmt.Errorf("weather API error: status=%d message=%s", weatherResult.Status, weatherResult.Message)
	}

	return &weatherResult, nil
}

// GetIPLocation IP定位
func (c *baiduMapsClient) GetIPLocation(ctx context.Context, ip string) (*IPLocationResult, error) {
	result, err := c.mcp.CallTool(ctx, "map_ip_location", map[string]any{
		"ip": ip,
	})
	if err != nil {
		return nil, err
	}

	var ipLocationResult IPLocationResult
	if err := json.Unmarshal([]byte(result), &ipLocationResult); err != nil {
		return nil, fmt.Errorf("failed to parse IP location result: %w", err)
	}

	return &ipLocationResult, nil
}

// GetTraffic 路况查询（兼容 road 与 road_name）
func (c *baiduMapsClient) GetTraffic(ctx context.Context, road, city string) (*TrafficResult, error) {
	result, err := c.mcp.CallTool(ctx, "map_road_traffic", map[string]any{
		"road":      road,
		"road_name": road,
		"city":      city,
	})
	if err != nil {
		return nil, err
	}

	var trafficResult TrafficResult
	if err := json.Unmarshal([]byte(result), &trafficResult); err != nil {
		return nil, fmt.Errorf("failed to parse traffic result: %w", err)
	}

	return &trafficResult, nil
}

func (c *baiduMapsClient) Close() error {
	if c.mcp != nil {
		return c.mcp.Close()
	}
	return nil
}

// 数据结构定义

// GeocodeResult 与百度返回结构对齐
// 示例：{"status":0,"result":{"location":{"lng":116.30,"lat":40.05},"precise":1,"confidence":80,"comprehension":100,"level":"门址"}}
type GeocodeResult struct {
	Status int `json:"status"`
	Result struct {
		Location struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"location"`
		Precise       int    `json:"precise"`
		Confidence    int    `json:"confidence"`
		Comprehension int    `json:"comprehension"`
		Level         string `json:"level"`
	} `json:"result"`
}

type ReverseGeocodeResult struct {
	FormattedAddress string `json:"formatted_address"`
	UID              string `json:"uid"`
	AddressComponent struct {
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Street   string `json:"street"`
	} `json:"address_component"`
}

type Place struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UID       string  `json:"uid"`
	Type      string  `json:"type"`
}

type DirectionsResult struct {
	Distance string `json:"distance"`
	Duration string `json:"duration"`
	Route    string `json:"route"`
}

// WeatherResult 精确匹配提供的返回结构
type WeatherResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Location struct {
			Country  string `json:"country"`
			Province string `json:"province"`
			City     string `json:"city"`
			Name     string `json:"name"`
			ID       string `json:"id"`
		} `json:"location"`
		Now struct {
			Text      string  `json:"text"`
			Temp      int     `json:"temp"`
			FeelsLike int     `json:"feels_like"`
			Rh        int     `json:"rh"`
			WindClass string  `json:"wind_class"`
			WindDir   string  `json:"wind_dir"`
			Prec1h    float64 `json:"prec_1h"`
			Clouds    int     `json:"clouds"`
			Vis       int     `json:"vis"`
			Aqi       int     `json:"aqi"`
			Pm25      int     `json:"pm25"`
			Pm10      int     `json:"pm10"`
			No2       int     `json:"no2"`
			So2       int     `json:"so2"`
			O3        int     `json:"o3"`
			Co        float64 `json:"co"`
			WindAngle int     `json:"wind_angle"`
			Uvi       int     `json:"uvi"`
			Pressure  int     `json:"pressure"`
			Dpt       int     `json:"dpt"`
			Uptime    string  `json:"uptime"`
		} `json:"now"`
		Indexes []struct {
			Name   string `json:"name"`
			Brief  string `json:"brief"`
			Detail string `json:"detail"`
		} `json:"indexes"`
		Alerts []struct {
			Type  string `json:"type"`
			Level string `json:"level"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
		} `json:"alerts"`
		Forecasts []struct {
			TextDay   string `json:"text_day"`
			TextNight string `json:"text_night"`
			High      int    `json:"high"`
			Low       int    `json:"low"`
			WcDay     string `json:"wc_day"`
			WdDay     string `json:"wd_day"`
			WcNight   string `json:"wc_night"`
			WdNight   string `json:"wd_night"`
			Date      string `json:"date"`
			Week      string `json:"week"`
		} `json:"forecasts"`
		ForecastHours []struct {
			Text      string  `json:"text"`
			TempFc    int     `json:"temp_fc"`
			WindClass string  `json:"wind_class"`
			WindDir   string  `json:"wind_dir"`
			Rh        int     `json:"rh"`
			Prec1h    float64 `json:"prec_1h"`
			Clouds    int     `json:"clouds"`
			WindAngle int     `json:"wind_angle"`
			Pop       int     `json:"pop"`
			Uvi       int     `json:"uvi"`
			Pressure  int     `json:"pressure"`
			Dpt       int     `json:"dpt"`
			DataTime  string  `json:"data_time"`
		} `json:"forecast_hours"`
	} `json:"result"`
}

type IPLocationResult struct {
	IP        string  `json:"ip"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TrafficResult struct {
	RoadName string `json:"road_name"`
	Status   string `json:"status"`
	Speed    string `json:"speed"`
}

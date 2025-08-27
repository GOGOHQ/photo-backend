# Photo Backend API 文档

## 概述

Photo Backend 提供了两个主要的 MCP 服务：
1. **小红书 MCP 服务** - 用于搜索和获取小红书内容
2. **百度地图 MCP 服务** - 提供地理编码、路线规划、天气查询等功能

## 小红书 API

### 1. 搜索热门帖子

**GET** `/api/xhs/hot`

**参数：**
- `limit` (可选): 返回结果数量，默认 10
- `q` (可选): 搜索关键词，如果不提供则返回热门帖子
- `skip_login` (可选): 是否跳过登录，默认 false

**示例：**
```bash
# 搜索关键词
curl "http://localhost:8080/api/xhs/hot?q=旅行&limit=5"

# 获取热门帖子
curl "http://localhost:8080/api/xhs/hot?limit=10"

# 跳过登录
curl "http://localhost:8080/api/xhs/hot?q=美食&skip_login=true"
```

**响应：**
```json
{
  "data": [
    {
      "id": "",
      "title": "秋天穷游旅游城市推荐✔️大学生穷游必看",
      "author": "",
      "likes": 0,
      "excerpt": "",
      "post_url": "https://www.xiaohongshu.com/search_result/68a85a11000000001d0233d9"
    }
  ]
}
```

### 2. 搜索帖子链接

**GET** `/api/xhs/search`

**参数：**
- `q` (必需): 搜索关键词
- `limit` (可选): 返回结果数量，默认 5
- `skip_login` (可选): 是否跳过登录，默认 false

**示例：**
```bash
curl "http://localhost:8080/api/xhs/search?q=美食&limit=3"
```

**响应：**
```json
{
  "links": [
    "https://www.xiaohongshu.com/search_result/6777dde4000000000900d5fb",
    "https://www.xiaohongshu.com/search_result/6732c7a9000000003c01aee2"
  ]
}
```

## 百度地图 API

### 1. 地理编码

**GET** `/api/baidu-maps/geocode`

**参数：**
- `address` (必需): 地址字符串

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/geocode?address=北京市朝阳区"
```

**响应：**
```json
{
  "data": {
    "latitude": 39.9219,
    "longitude": 116.4551,
    "address": "北京市朝阳区"
  }
}
```

### 2. 逆地理编码

**GET** `/api/baidu-maps/reverse-geocode`

**参数：**
- `lat` (必需): 纬度
- `lng` (必需): 经度

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/reverse-geocode?lat=39.9042&lng=116.4074"
```

**响应：**
```json
{
  "data": {
    "formatted_address": "北京市东城区天安门",
    "uid": "123456",
    "address_component": {
      "country": "中国",
      "province": "北京市",
      "city": "北京市",
      "district": "东城区",
      "street": "天安门"
    }
  }
}
```

### 3. 搜索地点

**GET** `/api/baidu-maps/search-places`

**参数：**
- `q` (必需): 搜索关键词
- `city` (可选): 城市名称
- `limit` (可选): 返回结果数量，默认 10

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/search-places?q=餐厅&city=北京&limit=5"
```

**响应：**
```json
{
  "data": [
    {
      "name": "全聚德烤鸭店",
      "address": "北京市东城区前门大街30号",
      "latitude": 39.8994,
      "longitude": 116.3974,
      "uid": "123456",
      "type": "餐饮服务"
    }
  ]
}
```

### 4. 路线规划

**GET** `/api/baidu-maps/directions`

**参数：**
- `origin` (必需): 起点
- `destination` (必需): 终点
- `mode` (可选): 出行方式，可选值：driving, walking, cycling, transit，默认 driving

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/directions?origin=北京站&destination=天安门&mode=driving"
```

**响应：**
```json
{
  "data": {
    "distance": "3.2公里",
    "duration": "15分钟",
    "route": "北京站 → 东长安街 → 天安门"
  }
}
```

### 5. 天气查询

**GET** `/api/baidu-maps/weather`

**参数：**
- `city` (必需): 城市名称

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/weather?city=北京"
```

**响应：**
```json
{
  "data": {
    "city": "北京",
    "temperature": "15°C",
    "weather": "晴",
    "humidity": "45%",
    "wind": "东北风 3级"
  }
}
```

### 6. IP定位

**GET** `/api/baidu-maps/ip-location`

**参数：**
- `ip` (必需): IP地址

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/ip-location?ip=8.8.8.8"
```

**响应：**
```json
{
  "data": {
    "ip": "8.8.8.8",
    "city": "美国",
    "latitude": 37.4056,
    "longitude": -122.0775
  }
}
```

### 7. 路况查询

**GET** `/api/baidu-maps/traffic`

**参数：**
- `road` (必需): 道路名称
- `city` (必需): 城市名称

**示例：**
```bash
curl "http://localhost:8080/api/baidu-maps/traffic?road=长安街&city=北京"
```

**响应：**
```json
{
  "data": {
    "road_name": "长安街",
    "status": "畅通",
    "speed": "45km/h"
  }
}
```

## 配置说明

### MCP 配置文件 (mcp.json)

```json
{
  "mcpServers": {
    "xhs-local-stdio": {
      "name": "xiaohongshu stdio",
      "command": "/path/to/python3",
      "args": [
        "/path/to/xiaohongshu_mcp.py",
        "--stdio"
      ],
      "isActive": true
    },
    "baidu-maps": {
      "name": "Baidu Maps MCP",
      "command": "python",
      "args": [
        "-m",
        "mcp_server_baidu_maps",
        "--api-key",
        "YOUR_BAIDU_MAPS_API_KEY"
      ],
      "isActive": true
    }
  }
}
```

### 环境变量

- `MCP_CONFIG_PATH`: MCP 配置文件路径（可选，默认为 ./mcp.json）

## 测试工具

### 小红书测试
```bash
# 编译测试工具
go build ./cmd/tools/xhs_debug

# 测试搜索功能
MCP_CONFIG_PATH=./mcp.json ./xhs_debug -q "旅行" -limit 3

# 跳过登录测试
MCP_CONFIG_PATH=./mcp.json ./xhs_debug -q "美食" -limit 2 -skip-login
```

### 百度地图测试
```bash
# 编译测试工具
go build ./cmd/tools/baidu_maps_test

# 测试地理编码
MCP_CONFIG_PATH=./mcp.json ./baidu_maps_test -test geocode

# 测试搜索地点
MCP_CONFIG_PATH=./mcp.json ./baidu_maps_test -test search-places

# 测试路线规划
MCP_CONFIG_PATH=./mcp.json ./baidu_maps_test -test directions
```

## 错误处理

所有 API 在发生错误时都会返回以下格式：

```json
{
  "error": "错误描述信息"
}
```

常见错误：
- `400 Bad Request`: 参数错误
- `502 Bad Gateway`: MCP 服务错误

## 注意事项

1. **小红书服务**：需要先登录才能正常搜索，API 会自动处理登录流程
2. **百度地图服务**：需要有效的百度地图 API Key
3. **超时设置**：所有 API 都有合理的超时设置，避免长时间等待
4. **并发限制**：建议控制并发请求数量，避免对 MCP 服务造成压力

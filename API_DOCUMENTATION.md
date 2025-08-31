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
curl "http://192.168.1.22:8080/api/xhs/hot?q=旅行&limit=5"

# 获取热门帖子
curl "http://192.168.1.22:8080/api/xhs/hot?limit=10"

# 跳过登录
curl "http://192.168.1.22:8080/api/xhs/hot?q=美食&skip_login=true"
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
curl "http://192.168.1.22:8080/api/xhs/search?q=美食&limit=3"
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
curl "http://192.168.1.22:8080/api/baidu-maps/geocode?address=北京市朝阳区"
```

**响应：**
```json
{
  "data": {
        "status": 0,
        "result": {
            "location": {
                "lng": 116.44955872950158,
                "lat": 39.926374523079886
            },
            "precise": 0,
            "confidence": 20,
            "comprehension": 100,
            "level": "区县"
        }
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
curl "http://192.168.1.22:8080/api/baidu-maps/reverse-geocode?lat=39.9042&lng=116.4074"
```

**响应：**
```json
{
  {
    "data": {
        "status": 0,
        "result": {
            "location": {
                "lng": 116.40739999999992,
                "lat": 39.90420007788774
            },
            "formatted_address": "北京市东城区前门街道长巷二条乙5号",
            "addressComponent": {
                "country": "中国",
                "province": "北京市",
                "city": "北京市",
                "district": "东城区",
                "town": "前门街道",
                "street": "长巷二条",
                "street_number": "乙5号",
                "adcode": "110101"
            },
            "business": "前门,珠市口,大栅栏"
        }
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
curl "http://192.168.1.22:8080/api/baidu-maps/search-places?q=美食&tag=餐饮服务&region=北京市&location=39.9042,116.4074&radius=2000"
```

**响应：**
```json
{
 "status":0,"message":"ok","result_type":"poi_type","query_type":"general","results":[{"name":"福兴居家·老北京四季涮肉(前门店)","location":{"lat":39.89924060676216,"lng":116.40068947246995},"address":"煤市街博兴胡同南侧（大栅栏社区卫生服务中心斜对面）","province":"北京市","city":"北京市","area":"西城区","town":"大栅栏街道","town_code":110102013,"street_id":"1bf85b25d2f747bda1a9340b","detail":1,"uid":"1bf85b25d2f747bda1a9340b","detail_info":{"classified_poi_tag":"美食;中餐馆;火锅店;老北京火锅","distance":794,"tag":"美食;中餐厅","navi_location":{"lat":39.899279647296034,"lng":116.40097393800443},"type":"cater","detail_url":"http://api.map.baidu.com/place/detail?uid=1bf85b25d2f747bda1a9340b\u0026output=html\u0026source=placeapi_v3","price":"87.0","overall_rating":"4.0","comment_num":"7","shop_hours":"10:00-23:00","label":"老北京火锅","children":[]}}]
}
```

<!-- ### 4. 路线规划

**GET** `/api/baidu-maps/directions`

**参数：**
- `origin` (必需): 起点
- `destination` (必需): 终点
- `mode` (可选): 出行方式，可选值：driving, walking, cycling, transit，默认 driving

**示例：**
```bash
curl "http://192.168.1.22:8080/api/baidu-maps/directions?origin=北京站&destination=天安门&mode=driving"
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
``` -->

### 5. 天气查询

**GET** `/api/baidu-maps/weather`

**参数：**
- `location` (必需): 经纬度坐标：格式是{经度,维度},用逗号分隔

**示例：**
```bash
curl "http://127.0.0.1:8080/api/baidu-maps/weather?location=116.391275,39.906217"
```

**响应：**
```json
{
  
}
```

<!-- ### 6. IP定位

**GET** `/api/baidu-maps/ip-location`

**参数：**
- `ip` (必需): IP地址

**示例：**
```bash
curl "http://192.168.1.22:8080/api/baidu-maps/ip-location?ip=8.8.8.8"
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
curl "http://192.168.1.22:8080/api/baidu-maps/traffic?road=长安街&city=北京"
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
``` -->


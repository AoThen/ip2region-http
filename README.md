# ip2region-http

轻量级HTTP IP地理位置查询API，仅支持IPv4，使用Go标准库编写，无外部依赖。

## 特性

- ⚡ **超快查询**：亚毫秒级查询性能（内存数据库）
- 📦 **零依赖**：仅使用Go标准库，无外部包依赖
- 🔧 **轻量级**：~584行代码，静态编译
- 🚀 **自动下载**：启动时自动下载数据库
- 🐳 **Docker就绪**：完整的Docker和Docker Compose支持

## 快速开始

### 二进制运行

```bash
# 克隆或下载项目
cd ip2region-http

# 编译
go build -o ip2region-http .

# 运行（自动下载数据库）
./ip2region-http

# 自定义端口
SERVER_PORT=9090 ./ip2region-http
```

### Docker运行

```bash
# 构建镜像
docker build -t ip2region-http .

# 运行容器
docker run -p 8080:8080 ip2region-http
```

### Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## API文档

### 1. 查询IP位置

```http
GET /search?ip=1.2.3.4&format=full
```

**参数：**
- `ip`（必需）：IP地址
- `format`（可选）：响应格式
  - `full`（默认）：完整JSON响应
  - `text`：纯文本响应
  - `fields`：拆分字段的JSON响应

**示例：**

```bash
# 完整JSON响应
curl "http://localhost:8080/search?ip=1.2.3.4"
# {"code":0,"msg":"success","data":{"ip":"1.2.3.4","region":"中国|广东省|深圳市|阿里云"}}

# 纯文本响应
curl "http://localhost:8080/search?ip=1.2.3.4&format=text"
# 中国|广东省|深圳市|阿里云

# 拆分字段JSON响应
curl "http://localhost:8080/search?ip=1.2.3.4&format=fields"
# {"code":0,"msg":"success","data":{"ip":"1.2.3.4","country":"中国","province":"广东省","city":"深圳市","district":"","isp":"阿里云"}}
```

### 2. 健康检查

```http
GET /health
```

```bash
curl "http://localhost:8080/health"
# {"code":0,"msg":"ok","data":{"status":"healthy","db_loaded":true,"db_path":"/app/data/ip2region_n.xdb","db_size":8388608}}
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SERVER_PORT` | 8080 | HTTP服务端口 |
| `DB_PATH` | data/ip2region_n.xdb | 数据库文件路径 |
| `DB_URL` | - | 数据库下载URL（可选） |
| `DOWNLOAD_MODE` | true | 是否自动下载数据库 |

## 项目结构

```
ip2region-http/
├── main.go              # HTTP处理器和主程序
├── config.go            # 配置管理
├── download.go          # 数据库下载
├── go.mod               # Go模块（无外部依赖）
├── Dockerfile           # Docker构建
├── docker-compose.yml   # Docker Compose配置
├── test-api.sh          # API测试脚本
└── xdb/                 # 核心库（本地包）
    ├── searcher.go      # 搜索器实现
    ├── util.go          # 工具函数
    ├── version.go       # IPv4版本定义
    └── header.go        # 文件头解析
```

## Docker Compose示例

### 基础配置

```yaml
version: '3.8'

services:
  ip2region-http:
    build: .
    container_name: ip2region-http
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - SERVER_PORT=8080
      - DB_PATH=/app/data/ip2region_n.xdb
      - DOWNLOAD_MODE=true
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 使用Nginx反向代理

```yaml
version: '3.8'

services:
  ip2region-http:
    build: .
    container_name: ip2region-http
    ports:
      - "8080:8080"
    networks:
      - ip2network

  nginx:
    image: nginx:alpine
    container_name: ip2region-nginx
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - ip2region-http
    networks:
      - ip2network

networks:
  ip2network:
    driver: bridge
```

## 编译状态

✅ **完全独立**
- 无外部依赖包
- 零外部Go模块
- 仅使用Go标准库
- 所有代码自包含

## 故障排查

**数据库下载失败：**
```bash
# 手动下载数据库
mkdir -p data
wget https://github.com/hel2o/ip2region_update/releases/download/250820/ip2region_n.xdb \
  -O data/ip2region_n.xdb

# 启动时禁用自动下载
DOWNLOAD_MODE=false ./ip2region-http
```

**端口已被占用：**
```bash
# 使用自定义端口
SERVER_PORT=9090 ./ip2region-http
```

**权限问题：**
```bash
chmod +x ip2region-http
```

## 性能指标

- **启动时间**：< 1秒（数据库加载）
- **查询延迟**：< 100微秒（内存模式）
- **并发能力**：单例可安全处理数千并发
- **内存占用**：~10-50MB（取决于数据库大小）

## 许可证

Apache License 2.0

## 致谢

本项目基于 [ip2region](https://github.com/lionsoul2014/ip2region) 项目，提取核心功能并简化为独立的HTTP服务。

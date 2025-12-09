# Ecommerce Backend API

Backend API cho ứng dụng ecommerce được xây dựng bằng Go và Gin framework.

## Cấu trúc dự án

```
ecommerce-be/
├── cmd/              # Application entrypoints
├── internal/         # Private application code
├── pkg/              # Public library code
├── config/           # Configuration files
│   └── config.go     # Config loader
├── database/         # Database connection
│   └── database.go   # Database setup
├── main.go           # Application entry point
├── go.mod            # Go modules
├── .env              # Environment variables (not in git)
└── README.md         # This file
```

## Yêu cầu

- Go 1.21 hoặc cao hơn
- PostgreSQL 12 hoặc cao hơn (hoặc Docker)
- Redis (hoặc Docker)
- Docker và Docker Compose (tùy chọn - để chạy PostgreSQL và Redis)

## Cài đặt

### Cách 1: Sử dụng Docker (Khuyến nghị)

1. Clone repository:
```bash
git clone <repository-url>
cd ecommerce-be
```

2. Cài đặt dependencies:
```bash
go mod download
```

3. Tạo file `.env` và cập nhật thông tin:
```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=ecommerce_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

PORT=8080
```

4. Khởi động PostgreSQL và Redis bằng Docker:
```bash
docker-compose up -d
```

5. Chạy ứng dụng:
```bash
go run main.go
```

Server sẽ chạy tại `http://localhost:8080`

**Truy cập pgAdmin:**
- URL: `http://localhost:5050`
- Email: `admin@admin.com`
- Password: `admin`

**Kết nối PostgreSQL trong pgAdmin:**
- Host: `postgres` (tên service trong docker-compose) hoặc `localhost`
- Port: `5432` (port trong container)
- Database: `ecommerce_db`
- Username: `postgres`
- Password: `12345678`

### Cách 2: Cài đặt thủ công

1. Clone repository và cài đặt dependencies (giống bước 1-2 ở trên)

2. Cài đặt PostgreSQL và Redis trên máy

3. Tạo file `.env` với cấu hình phù hợp:
```env
DB_HOST=localhost
DB_PORT=5432  # Port của PostgreSQL đã cài
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

PORT=8080
```

4. Chạy ứng dụng:
```bash
go run main.go
```

## Quản lý Docker Services

**Khởi động services:**
```bash
docker-compose up -d
```

**Dừng services:**
```bash
docker-compose down
```

**Xem logs:**
```bash
docker-compose logs -f
```

**Xem logs của một service cụ thể:**
```bash
docker-compose logs -f postgres
docker-compose logs -f redis
```

**Dừng và xóa tất cả dữ liệu (volumes):**
```bash
docker-compose down -v
```

**Khởi động lại một service:**
```bash
docker-compose restart postgres
```

## API Endpoints

### Health Check
- `GET /health` - Kiểm tra trạng thái server

### API v1
- `GET /api/v1/` - API information

## Database Configuration

Cấu hình database trong file `.env`:

**Khi dùng Docker:**
```env
DB_HOST=localhost
DB_PORT=5433  # Port mapping trong docker-compose
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=ecommerce_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Khi cài đặt thủ công:**
```env
DB_HOST=localhost
DB_PORT=5432  # Port mặc định của PostgreSQL
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce_db

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Development

Chạy với hot reload (cần install air):
```bash
air
```

## License

MIT


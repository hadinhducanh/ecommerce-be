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
- PostgreSQL 12 hoặc cao hơn

## Cài đặt

1. Clone repository:
```bash
git clone <repository-url>
cd ecommerce-be
```

2. Cài đặt dependencies:
```bash
go mod download
```

3. Tạo file `.env` từ `.env.example` và cập nhật thông tin database:
```bash
cp .env.example .env
```

4. Chạy ứng dụng:
```bash
go run main.go
```

Server sẽ chạy tại `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Kiểm tra trạng thái server

### API v1
- `GET /api/v1/` - API information

## Database Configuration

Cấu hình database trong file `.env`:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=ecommerce_db
```

## Development

Chạy với hot reload (cần install air):
```bash
air
```

## License

MIT


# 🛍️ Ecommerce Backend API

Backend API cho ứng dụng ecommerce được xây dựng bằng Go và Gin framework.

## 📋 Mục lục

- [Giới thiệu](#giới-thiệu)
- [Công nghệ sử dụng](#công-nghệ-sử-dụng)
- [Yêu cầu hệ thống](#yêu-cầu-hệ-thống)
- [Cài đặt và chạy dự án](#cài-đặt-và-chạy-dự-án)
- [Cấu hình](#cấu-hình)
- [Quản lý Database](#quản-lý-database)
- [API Documentation](#api-documentation)
- [Troubleshooting](#troubleshooting)

## 📖 Giới thiệu

Đây là backend API cho hệ thống ecommerce, cung cấp các tính năng:

- 🔐 Xác thực người dùng (JWT)
- 👤 Quản lý người dùng
- 📦 Quản lý sản phẩm
- 🏷️ Quản lý danh mục
- 🛒 Giỏ hàng
- 💳 Đơn hàng và thanh toán
- ⭐ Đánh giá sản phẩm
- 💬 Chat
- 📧 Email service (OTP, verification)
- ☁️ Upload hình ảnh (Cloudinary)
- 🚀 Redis caching

## 🛠 Công nghệ sử dụng

- **Language:** Go 1.23+
- **Framework:** Gin
- **Database:** PostgreSQL 15
- **Cache:** Redis 7
- **ORM:** GORM
- **Authentication:** JWT
- **File Storage:** Cloudinary
- **Email:** SMTP (Gmail)
- **Container:** Docker & Docker Compose

## 💻 Yêu cầu hệ thống

### Yêu cầu bắt buộc:

- **Go:** Version 1.23.0 hoặc cao hơn ([Download Go](https://go.dev/dl/))
- **Docker Desktop:** Để chạy PostgreSQL, Redis và pgAdmin ([Download Docker](https://www.docker.com/products/docker-desktop))
- **Git:** Để clone repository

### Kiểm tra cài đặt:

```bash
# Kiểm tra Go
go version
# Kết quả mong đợi: go version go1.23.0 hoặc cao hơn

# Kiểm tra Docker
docker --version
docker-compose --version
```

## 🚀 Cài đặt và chạy dự án

### Bước 1: Clone repository

```bash
git clone <repository-url>
cd ecommerce-be
```

### Bước 2: Cài đặt Go dependencies

```bash
go mod download
```

> **Lưu ý:** Nếu gặp lỗi, chạy `go mod tidy` để làm sạch dependencies.

### Bước 3: Tạo file `.env`

Tạo file `.env` trong thư mục gốc của dự án:

```bash
# Windows PowerShell
New-Item -Path .env -ItemType File

# Mac/Linux
touch .env
```

Sau đó copy nội dung sau vào file `.env`:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=ecommerce_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
PORT=8080

# JWT Secret (tạo secret key mới cho production)
JWT_SECRET=301Ab42TpjIXhQVceE8J6Z/3z/ocytyTj0ut/Gx7Ckw=

# SMTP Configuration (Email Service)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

# Cloudinary Configuration (Image Upload)
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
```

> **⚠️ Quan trọng:**
> - Thay đổi `JWT_SECRET` cho môi trường production
> - Cấu hình SMTP với email của bạn (xem hướng dẫn bên dưới)
> - Cấu hình Cloudinary cho upload ảnh (xem hướng dẫn bên dưới)

### Bước 4: Khởi động Database với Docker

```bash
docker-compose up -d
```

Lệnh này sẽ khởi động 3 services:
- **PostgreSQL:** Database chính (port 5433)
- **Redis:** Cache server (port 6379)
- **pgAdmin:** Database management tool (port 5050)

Kiểm tra trạng thái:

```bash
docker-compose ps
```

Kết quả mong đợi:

```
NAME                  STATUS              PORTS
ecommerce-postgres    Up (healthy)        0.0.0.0:5433->5432/tcp
ecommerce-redis       Up (healthy)        0.0.0.0:6379->6379/tcp
ecommerce-pgadmin     Up                  0.0.0.0:5050->80/tcp
```

### Bước 5: Chạy ứng dụng

```bash
go run main.go
```

Nếu thành công, bạn sẽ thấy:

```
🚀 Server starting on port :8080
```

Server sẽ chạy tại: **http://localhost:8080**

### Bước 6: Kiểm tra API

Mở trình duyệt hoặc Postman và truy cập:

```
http://localhost:8080/health
```

Kết quả mong đợi:

```json
{
  "status": "OK"
}
```

## ⚙️ Cấu hình

### 📧 Cấu hình Email Service (Gmail)

1. Đăng nhập Gmail
2. Vào **Google Account Settings** → **Security**
3. Bật **2-Step Verification**
4. Tạo **App Password**:
   - Vào **Security** → **App passwords**
   - Chọn app: **Mail**
   - Chọn device: **Other (Custom name)**
   - Copy password được tạo
5. Cập nhật `.env`:

```env
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-16-digit-app-password
```

### ☁️ Cấu hình Cloudinary (Upload ảnh)

1. Tạo tài khoản tại [Cloudinary](https://cloudinary.com/)
2. Vào **Dashboard** → Copy thông tin:
   - Cloud Name
   - API Key
   - API Secret
3. Cập nhật `.env`:

```env
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
```

### 🔑 Tạo JWT Secret mới

Để bảo mật hơn, tạo JWT secret mới:

```bash
# Windows PowerShell
$bytes = New-Object byte[] 32
[Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
[Convert]::ToBase64String($bytes)

# Mac/Linux
openssl rand -base64 32
```

Copy kết quả vào `JWT_SECRET` trong file `.env`.

## 🗄 Quản lý Database

### Truy cập pgAdmin

pgAdmin là công cụ quản lý PostgreSQL database qua giao diện web.

1. Mở trình duyệt: **http://localhost:5050**
2. Đăng nhập:
   - Email: `admin@admin.com`
   - Password: `admin`

### Kết nối Database trong pgAdmin

1. Click **Add New Server**
2. Tab **General:**
   - Name: `Ecommerce DB`
3. Tab **Connection:**
   - Host: `postgres` (hoặc `localhost`)
   - Port: `5432`
   - Database: `ecommerce_db`
   - Username: `postgres`
   - Password: `12345678`
4. Click **Save**

### Xem cấu trúc Database

Sau khi ứng dụng chạy lần đầu, GORM sẽ tự động tạo các tables:

- `users` - Người dùng
- `products` - Sản phẩm
- `categories` - Danh mục
- `category_children` - Danh mục con
- `cart_items` - Giỏ hàng
- `orders` - Đơn hàng
- `order_items` - Chi tiết đơn hàng
- `reviews` - Đánh giá
- `wishlists` - Danh sách yêu thích
- `addresses` - Địa chỉ
- `payments` - Thanh toán
- `chats` - Phòng chat
- `chat_messages` - Tin nhắn

### Seed dữ liệu mẫu

Để thêm dữ liệu mẫu cho products:

```bash
go run cmd/seed/seed.go
```

### Quản lý Docker Services

**Xem logs:**

```bash
# Tất cả services
docker-compose logs -f

# Chỉ PostgreSQL
docker-compose logs -f postgres

# Chỉ Redis
docker-compose logs -f redis
```

**Dừng services:**

```bash
docker-compose down
```

**Khởi động lại:**

```bash
docker-compose restart
```

**Xóa tất cả dữ liệu (reset database):**

```bash
docker-compose down -v
```

**Backup database:**

```bash
docker exec ecommerce-postgres pg_dump -U postgres ecommerce_db > backup.sql
```

**Restore database:**

```bash
docker exec -i ecommerce-postgres psql -U postgres ecommerce_db < backup.sql
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

API sử dụng JWT Bearer Token. Thêm header:

```
Authorization: Bearer <your-token>
```

### Main Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Đăng ký tài khoản
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/verify-otp` - Xác thực OTP
- `POST /api/v1/auth/resend-otp` - Gửi lại OTP
- `POST /api/v1/auth/forgot-password` - Quên mật khẩu
- `POST /api/v1/auth/reset-password` - Reset mật khẩu

#### Users
- `GET /api/v1/users/profile` - Xem profile
- `PUT /api/v1/users/profile` - Cập nhật profile
- `POST /api/v1/users/addresses` - Thêm địa chỉ

#### Products
- `GET /api/v1/products` - Danh sách sản phẩm
- `GET /api/v1/products/:id` - Chi tiết sản phẩm
- `POST /api/v1/products` - Tạo sản phẩm (Admin)
- `PUT /api/v1/products/:id` - Cập nhật sản phẩm (Admin)
- `DELETE /api/v1/products/:id` - Xóa sản phẩm (Admin)

#### Categories
- `GET /api/v1/categories` - Danh sách danh mục
- `GET /api/v1/categories/:id` - Chi tiết danh mục
- `POST /api/v1/categories` - Tạo danh mục (Admin)

#### Cloudinary
- `POST /api/v1/cloudinary/upload` - Upload ảnh

### Chi tiết API

Xem thêm tài liệu chi tiết trong thư mục `docs/`:

- [CATEGORY_API.md](docs/CATEGORY_API.md) - API danh mục
- [CATEGORY_POSTMAN_GUIDE.md](docs/CATEGORY_POSTMAN_GUIDE.md) - Hướng dẫn test API
- [FLUTTER_CATEGORY_GUIDE.md](docs/FLUTTER_CATEGORY_GUIDE.md) - Tích hợp Flutter
- [FLUTTER_PRODUCT_GUIDE.md](docs/FLUTTER_PRODUCT_GUIDE.md) - Product API Flutter
- [FLUTTER_USER_PROFILE_GUIDE.md](docs/FLUTTER_USER_PROFILE_GUIDE.md) - User API Flutter

## 🐛 Troubleshooting

### Lỗi: "Failed to connect to database"

**Nguyên nhân:** Docker chưa chạy hoặc PostgreSQL chưa ready.

**Giải pháp:**

```bash
# Kiểm tra Docker
docker-compose ps

# Khởi động lại
docker-compose restart postgres

# Xem logs
docker-compose logs postgres
```

### Lỗi: "Redis không kết nối được"

**Nguyên nhân:** Redis chưa chạy (không ảnh hưởng nhiều, app vẫn chạy).

**Giải pháp:**

```bash
# Khởi động lại Redis
docker-compose restart redis

# Kiểm tra
docker-compose logs redis
```

### Lỗi: "port already in use"

**Nguyên nhân:** Port 8080, 5433, hoặc 6379 đã được sử dụng.

**Giải pháp:**

```bash
# Windows - Tìm process đang dùng port
netstat -ano | findstr :8080

# Kill process
taskkill /PID <process-id> /F

# Hoặc đổi port trong .env
PORT=8081
```

### Lỗi: "go: module not found"

**Giải pháp:**

```bash
go mod tidy
go mod download
```

### Reset hoàn toàn dự án

```bash
# Dừng và xóa containers + volumes
docker-compose down -v

# Xóa Go modules cache
go clean -modcache

# Cài lại
go mod download
docker-compose up -d
go run main.go
```

## 📁 Cấu trúc dự án

```
ecommerce-be/
├── cache/                 # Redis cache
│   └── redis.go
├── cmd/                   # Command line tools
│   └── seed/             # Database seeding
├── config/               # Configuration
│   └── config.go
├── database/             # Database connection
│   └── database.go
├── docs/                 # API documentation
├── dto/                  # Data Transfer Objects
├── handlers/             # HTTP handlers
├── middleware/           # Middleware (auth, cors)
├── models/              # Database models
├── routes/              # Route definitions
├── services/            # Business logic
├── utils/               # Utilities (jwt, password)
├── .env                 # Environment variables
├── docker-compose.yml   # Docker configuration
├── go.mod              # Go modules
├── main.go             # Application entry point
└── README.md           # This file
```


# Build binary
go build -o ecommerce-be.exe main.go

# Run binary
./ecommerce-be.exe
```



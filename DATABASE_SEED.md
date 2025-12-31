# 🌱 Hướng dẫn Import Dữ liệu vào Database (Seed Data)

Hướng dẫn này giúp bạn import dữ liệu mẫu vào database để có thể test và phát triển ứng dụng.

## 📋 Mục lục

- [Giới thiệu](#giới-thiệu)
- [Yêu cầu](#yêu-cầu)
- [Dữ liệu sẽ được import](#dữ-liệu-sẽ-được-import)
- [Cách chạy Seed](#cách-chạy-seed)
- [Chi tiết từng bước](#chi-tiết-từng-bước)
- [Xác minh dữ liệu](#xác-minh-dữ-liệu)
- [Xóa và Seed lại](#xóa-và-seed-lại)
- [Troubleshooting](#troubleshooting)

---

## 📖 Giới thiệu

**Seed data** là quá trình import dữ liệu mẫu vào database để:

- ✅ Test các chức năng của ứng dụng
- ✅ Có dữ liệu để phát triển frontend
- ✅ Demo sản phẩm cho khách hàng
- ✅ Có tài khoản admin để quản lý hệ thống

Dự án này có 2 file seed chính:

1. **`cmd/seed/seed.go`** - Seed Categories và Admin User
2. **`cmd/seed/products/seed.go`** - Seed Products

---

## ✅ Yêu cầu

### 1. Database đang chạy

Đảm bảo PostgreSQL đã được khởi động:

```bash
# Kiểm tra Docker services
docker-compose ps
```

Kết quả mong đợi:

```
NAME                  STATUS
ecommerce-postgres    Up (healthy)
```

Nếu chưa chạy:

```bash
docker-compose up -d
```

### 2. Server KHÔNG cần chạy

- Seed scripts chỉ cần database
- KHÔNG cần chạy `go run main.go`

### 3. File .env đã cấu hình đúng

Kiểm tra file `.env` có đầy đủ thông tin database:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=12345678
DB_NAME=ecommerce_db
```

---

## 📦 Dữ liệu sẽ được import

### 1. Seed Categories & Admin User (`cmd/seed/seed.go`)

#### Admin User
- **Email:** `admin@ecommerce.com`
- **Password:** `1`
- **Role:** `admin`
- **Status:** Đã kích hoạt, email đã xác thực

#### Categories Structure

**2 Parent Categories:**

1. 📱 **Điện thoại** (Smartphones)
   - iPhone
   - Samsung
   - Xiaomi
   - OPPO
   - Vivo
   - Realme
   - OnePlus
   - Google Pixel

2. 💻 **Laptop** (Laptops)
   - MacBook
   - Laptop Dell
   - Laptop HP
   - Laptop Lenovo
   - Laptop Asus
   - Laptop Acer
   - Laptop MSI
   - Laptop Razer

**Tổng:** 2 parent + 16 child categories = **18 categories**

### 2. Seed Products (`cmd/seed/products/seed.go`)

- **4 sản phẩm** cho mỗi child category
- Tổng: 16 categories × 4 products = **64 sản phẩm**

**Thông tin mỗi sản phẩm:**
- Tên sản phẩm (Tiếng Việt & English)
- Mô tả chi tiết
- Giá (VNĐ)
- Stock (số lượng tồn kho)
- SKU (mã sản phẩm)
- Images (URL hình ảnh)
- Category ID

**Ví dụ sản phẩm:**

| Tên | Giá | Stock | SKU |
|-----|-----|-------|-----|
| iPhone 15 Pro Max 256GB | 32,990,000 VNĐ | 50 | IPH15PM256 |
| Samsung Galaxy S24 Ultra 256GB | 29,990,000 VNĐ | 45 | SAMS24U256 |
| MacBook Pro 14" M3 Pro 512GB | 54,990,000 VNĐ | 30 | MBP14M3P512 |
| Dell XPS 13 Plus i7 16GB 512GB | 42,990,000 VNĐ | 25 | DELLXPS13P |

---

## 🚀 Cách chạy Seed

### Tóm tắt nhanh

```bash
# Bước 1: Seed Categories và Admin User (BẮT BUỘC chạy trước)
go run cmd/seed/seed.go

# Bước 2: Seed Products (chạy sau khi có categories)
go run cmd/seed/products/seed.go
```

> **⚠️ QUAN TRỌNG:**
> - **PHẢI** chạy `seed.go` trước để tạo categories
> - **SAU ĐÓ** mới chạy `products/seed.go`
> - Nếu chạy sai thứ tự, products/seed.go sẽ báo lỗi

---

## 📝 Chi tiết từng bước

### Bước 1: Seed Categories và Admin User

#### 1.1. Chạy lệnh

```bash
go run cmd/seed/seed.go
```

#### 1.2. Output mong đợi

```
👤 Starting to seed admin user...
✅ Admin user created successfully!
   📧 Email: admin@ecommerce.com
   🔑 Password: 1
   👤 Role: admin
   🆔 ID: 1

🌱 Starting to seed categories...
📦 Creating parent categories (Điện thoại, Laptop)...
✅ Created category: Điện thoại
✅ Created category: Laptop

🧹 Cleaning up old parent-child relationships...
✅ Cleaned up old relationships

📱 Creating phone child categories...
  ✓ Created/Updated: iPhone (ID: 3)
  ✓ Created/Updated: Samsung (ID: 4)
  ✓ Created/Updated: Xiaomi (ID: 5)
  ✓ Created/Updated: OPPO (ID: 6)
  ✓ Created/Updated: Vivo (ID: 7)
  ✓ Created/Updated: Realme (ID: 8)
  ✓ Created/Updated: OnePlus (ID: 9)
  ✓ Created/Updated: Google Pixel (ID: 10)

💻 Creating laptop child categories...
  ✓ Created/Updated: MacBook (ID: 11)
  ✓ Created/Updated: Laptop Dell (ID: 12)
  ✓ Created/Updated: Laptop HP (ID: 13)
  ✓ Created/Updated: Laptop Lenovo (ID: 14)
  ✓ Created/Updated: Laptop Asus (ID: 15)
  ✓ Created/Updated: Laptop Acer (ID: 16)
  ✓ Created/Updated: Laptop MSI (ID: 17)
  ✓ Created/Updated: Laptop Razer (ID: 18)

🔗 Creating parent-child relationships...
  📱 Adding 8 children to 'Điện thoại' (ID: 1)...
  💻 Adding 8 children to 'Laptop' (ID: 2)...

✅ Seeding categories completed!
📊 Summary:
  - Parent categories: 2 (Điện thoại, Laptop)
  - Phone children: 8
  - Laptop children: 8
  - Total categories: 18
```

#### 1.3. Xác minh thành công

- ✅ Thấy "✅ Seeding categories completed!"
- ✅ Tổng 18 categories được tạo
- ✅ Admin user đã được tạo

---

### Bước 2: Seed Products

#### 2.1. Chạy lệnh

```bash
go run cmd/seed/products/seed.go
```

#### 2.2. Output mong đợi

```
🌱 Starting to seed products...
📦 Fetching child categories...
✅ Found 16 child categories

📱 Creating products for category: iPhone (ID: 3)
  ✓ [1/4] Created/Updated: iPhone 15 Pro Max 256GB (ID: 1, Price: 32990000 VNĐ)
  ✓ [2/4] Created/Updated: iPhone 15 Plus 128GB (ID: 2, Price: 24990000 VNĐ)
  ✓ [3/4] Created/Updated: iPhone 14 Pro 256GB (ID: 3, Price: 26990000 VNĐ)
  ✓ [4/4] Created/Updated: iPhone 14 128GB (ID: 4, Price: 18990000 VNĐ)

📱 Creating products for category: Samsung (ID: 4)
  ✓ [1/4] Created/Updated: Samsung Galaxy S24 Ultra 256GB (ID: 5, Price: 29990000 VNĐ)
  ✓ [2/4] Created/Updated: Samsung Galaxy S24 Plus 256GB (ID: 6, Price: 24990000 VNĐ)
  ✓ [3/4] Created/Updated: Samsung Galaxy Z Fold5 512GB (ID: 7, Price: 40990000 VNĐ)
  ✓ [4/4] Created/Updated: Samsung Galaxy A54 5G 128GB (ID: 8, Price: 8990000 VNĐ)

... (tiếp tục cho tất cả categories)

✅ Seeding products completed!
📊 Summary:
  - Child categories: 16
  - Products per category: 4
  - Total products created: 64
```

#### 2.3. Xác minh thành công

- ✅ Thấy "✅ Seeding products completed!"
- ✅ Tổng 64 products được tạo
- ✅ Mỗi category có 4 products

---

## ✅ Xác minh dữ liệu

### Cách 1: Qua pgAdmin

#### 1. Truy cập pgAdmin

```
URL: http://localhost:5050
Email: admin@admin.com
Password: admin
```

#### 2. Kết nối Database

- Server: `postgres`
- Database: `ecommerce_db`
- Username: `postgres`
- Password: `12345678`

#### 3. Kiểm tra Tables

**Kiểm tra Users:**

```sql
SELECT id, email, name, role, is_active, is_email_verified 
FROM users 
WHERE role = 'admin';
```

Kết quả mong đợi:

| id | email | name | role | is_active | is_email_verified |
|----|-------|------|------|-----------|-------------------|
| 1 | admin@ecommerce.com | Administrator | admin | true | true |

**Kiểm tra Categories:**

```sql
-- Đếm tổng categories
SELECT COUNT(*) FROM categories;
-- Kết quả: 18

-- Xem parent categories
SELECT id, name, name_en FROM categories 
WHERE id IN (
  SELECT DISTINCT parent_id FROM category_children
);
-- Kết quả: Điện thoại, Laptop
```

**Kiểm tra Products:**

```sql
-- Đếm tổng products
SELECT COUNT(*) FROM products;
-- Kết quả: 64

-- Xem 10 products đầu tiên
SELECT id, name, price, stock, sku, category_id 
FROM products 
LIMIT 10;

-- Đếm products theo category
SELECT c.name, COUNT(p.id) as product_count
FROM categories c
LEFT JOIN products p ON p.category_id = c.id
GROUP BY c.id, c.name
ORDER BY product_count DESC;
```

**Kiểm tra Category Relationships:**

```sql
SELECT 
  p.name as parent_name,
  c.name as child_name
FROM category_children cc
JOIN categories p ON p.id = cc.parent_id
JOIN categories c ON c.id = cc.child_id
ORDER BY p.name, c.name;
```

### Cách 2: Qua API

#### 1. Khởi động Server

```bash
go run main.go
```

#### 2. Login với Admin

**Request:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@ecommerce.com",
    "password": "1"
  }'
```

**Response:**

```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@ecommerce.com",
    "name": "Administrator",
    "role": "admin"
  }
}
```

#### 3. Kiểm tra Categories

```bash
curl http://localhost:8080/api/v1/categories
```

#### 4. Kiểm tra Products

```bash
curl http://localhost:8080/api/v1/products
```

---

## 🔄 Xóa và Seed lại

### Tình huống: Muốn reset và import lại từ đầu

#### Cách 1: Xóa chỉ dữ liệu (giữ nguyên structure)

**Xóa qua pgAdmin:**

```sql
-- Xóa theo thứ tự (tránh foreign key constraint)
DELETE FROM products;
DELETE FROM category_children;
DELETE FROM categories;
DELETE FROM users WHERE role = 'admin';
```

**Sau đó seed lại:**

```bash
go run cmd/seed/seed.go
go run cmd/seed/products/seed.go
```

#### Cách 2: Reset toàn bộ database

**Xóa và tạo lại database:**

```bash
# Dừng và xóa tất cả (bao gồm volumes)
docker-compose down -v

# Khởi động lại
docker-compose up -d

# Đợi database ready (khoảng 10 giây)
timeout 10

# Chạy server một lần để tạo tables
go run main.go
# Nhấn Ctrl+C sau khi thấy "Server starting..."

# Seed lại dữ liệu
go run cmd/seed/seed.go
go run cmd/seed/products/seed.go
```

#### Cách 3: Seed sẽ tự động update

> **💡 Lưu ý:** 
> - Các seed scripts đã được thiết kế để **idempotent**
> - Nếu chạy lại, nó sẽ **update** thay vì tạo duplicate
> - An toàn để chạy nhiều lần

```bash
# Chạy lại an toàn, không tạo duplicate
go run cmd/seed/seed.go
go run cmd/seed/products/seed.go
```

---

## 🐛 Troubleshooting

### ❌ Lỗi: "Failed to connect to database"

**Nguyên nhân:**
- Database chưa chạy
- Cấu hình `.env` sai

**Giải pháp:**

```bash
# Kiểm tra database
docker-compose ps

# Nếu không chạy
docker-compose up -d

# Kiểm tra logs
docker-compose logs postgres

# Kiểm tra .env
cat .env
```

---

### ❌ Lỗi: "No child categories found"

**Nguyên nhân:**
- Chạy `products/seed.go` trước khi chạy `seed.go`
- Categories chưa được tạo

**Giải pháp:**

```bash
# Chạy đúng thứ tự
go run cmd/seed/seed.go       # TRƯỚC
go run cmd/seed/products/seed.go  # SAU
```

---

### ❌ Lỗi: "Duplicate key value violates unique constraint"

**Nguyên nhân:**
- Dữ liệu đã tồn tại
- Đang cố tạo duplicate

**Giải pháp:**

Seed scripts đã handle trường hợp này. Nếu vẫn gặp lỗi:

```bash
# Xóa dữ liệu cũ
# Vào pgAdmin và chạy SQL:
DELETE FROM products;
DELETE FROM category_children;
DELETE FROM categories;
DELETE FROM users WHERE email = 'admin@ecommerce.com';

# Seed lại
go run cmd/seed/seed.go
go run cmd/seed/products/seed.go
```

---

### ❌ Lỗi: "pq: relation does not exist"

**Nguyên nhân:**
- Tables chưa được tạo trong database
- Chưa chạy server lần đầu để GORM tạo tables

**Giải pháp:**

```bash
# Chạy server một lần để GORM tạo tables
go run main.go

# Đợi thấy "Server starting on port :8080"
# Sau đó Ctrl+C để dừng

# Chạy seed
go run cmd/seed/seed.go
go run cmd/seed/products/seed.go
```

---

### ⚠️ Warning: "Admin user already exists"

**Không phải lỗi!**

```
ℹ️  Admin user already exists: admin@ecommerce.com (ID: 1)
```

- Script phát hiện admin đã tồn tại
- Bỏ qua việc tạo mới
- Tiếp tục seed categories bình thường

---

## 📊 Tóm tắt Seed Data

| Loại dữ liệu | Số lượng | File seed | Thứ tự |
|--------------|----------|-----------|--------|
| Admin User | 1 | `cmd/seed/seed.go` | 1 |
| Parent Categories | 2 | `cmd/seed/seed.go` | 1 |
| Child Categories | 16 | `cmd/seed/seed.go` | 1 |
| Category Relationships | 16 | `cmd/seed/seed.go` | 1 |
| Products | 64 | `cmd/seed/products/seed.go` | 2 |

**Tổng cộng:** 1 admin + 18 categories + 64 products

---

## 🎯 Best Practices

### 1. Chạy seed khi nào?

✅ **NÊN chạy:**
- Lần đầu setup dự án
- Sau khi reset database
- Khi cần dữ liệu test mới
- Khi có thành viên mới join team

❌ **KHÔNG NÊN chạy:**
- Trên production database
- Khi đã có dữ liệu thật của khách hàng

### 2. Development vs Production

**Development:**
```bash
# Thoải mái seed và reset
go run cmd/seed/seed.go
go run cmd/seed/products/seed.go
```

**Production:**
```bash
# KHÔNG bao giờ chạy seed trên production
# Dùng migration scripts chính thức
```

### 3. Backup trước khi seed lại

```bash
# Backup database trước khi reset
docker exec ecommerce-postgres pg_dump -U postgres ecommerce_db > backup_$(date +%Y%m%d).sql
```

### 4. Custom seed data

Nếu muốn thay đổi dữ liệu seed:

1. Edit file `cmd/seed/seed.go` (categories)
2. Edit file `cmd/seed/products/seed.go` (products)
3. Chạy lại seed scripts

---

## ✅ Checklist sau khi Seed

- [ ] Database đang chạy
- [ ] Đã chạy `cmd/seed/seed.go` thành công
- [ ] Đã chạy `cmd/seed/products/seed.go` thành công
- [ ] Đã verify admin user qua pgAdmin hoặc login API
- [ ] Đã kiểm tra categories trong database
- [ ] Đã kiểm tra products trong database
- [ ] Server chạy bình thường với dữ liệu mới

---

## 📚 Tài liệu liên quan

- [README.md](../README.md) - Hướng dẫn setup dự án
- [CATEGORY_API.md](CATEGORY_API.md) - API documentation cho categories
- [DATABASE_MIGRATION.md] - Database migration guide (nếu có)

---

**🎉 Hoàn thành!** Bây giờ database đã có đầy đủ dữ liệu để bắt đầu phát triển và test!

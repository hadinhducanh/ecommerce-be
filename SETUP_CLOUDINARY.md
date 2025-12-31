# ☁️ Hướng dẫn Setup Cloudinary cho Upload Ảnh

Hướng dẫn này sẽ giúp bạn cấu hình Cloudinary để upload và quản lý hình ảnh sản phẩm trong ứng dụng.

## 📋 Mục lục

- [Giới thiệu Cloudinary](#giới-thiệu-cloudinary)
- [Bước 1: Tạo tài khoản Cloudinary](#bước-1-tạo-tài-khoản-cloudinary)
- [Bước 2: Lấy API Credentials](#bước-2-lấy-api-credentials)
- [Bước 3: Cấu hình .env](#bước-3-cấu-hình-env)
- [Bước 4: Test Upload](#bước-4-test-upload)
- [Quản lý Media](#quản-lý-media)
- [Troubleshooting](#troubleshooting)

---

## 🌟 Giới thiệu Cloudinary

**Cloudinary** là dịch vụ lưu trữ và quản lý media (ảnh, video) trên cloud.

### Tại sao dùng Cloudinary?

✅ **Free Plan hào phóng:**
- 25GB storage
- 25GB bandwidth/tháng
- 25,000 transformations/tháng

✅ **Tính năng:**
- Upload ảnh qua API
- Resize, crop, optimize ảnh tự động
- CDN toàn cầu (load ảnh nhanh)
- Quản lý media qua Dashboard

✅ **Phù hợp cho:**
- Ảnh sản phẩm ecommerce
- Avatar người dùng
- Ảnh trong reviews

---

## 📝 Bước 1: Tạo tài khoản Cloudinary

### 1.1. Truy cập Cloudinary

Mở trình duyệt và truy cập: **https://cloudinary.com/**

### 1.2. Sign Up

1. Click nút **"Sign Up"** hoặc **"Get Started for Free"**
2. Bạn có thể đăng ký bằng:
   - **Email** (khuyến nghị)
   - GitHub
   - Google Account

### 1.3. Đăng ký bằng Email

1. Chọn **"Sign up with email"**
2. Điền thông tin:
   ```
   Email: your-email@gmail.com
   Password: your-secure-password
   ```
3. Click **"Create Account"**

### 1.4. Xác thực Email

1. Kiểm tra hộp thư email
2. Tìm email từ Cloudinary
3. Click link xác thực trong email

### 1.5. Hoàn tất thiết lập tài khoản

1. Chọn **"Developer"** khi được hỏi role
2. Chọn plan: **"Free"** (đủ cho development)
3. Nhập tên công ty/dự án (tùy chọn): `Ecommerce Project`
4. Click **"Get Started"**

> **✅ Thành công!** Bạn sẽ được chuyển đến Dashboard

---

## 🔑 Bước 2: Lấy API Credentials

### 2.1. Truy cập Dashboard

Sau khi đăng nhập, bạn sẽ thấy **Dashboard** (hoặc truy cập: https://cloudinary.com/console)

### 2.2. Tìm Account Details

Ngay trên đầu Dashboard, bạn sẽ thấy mục **"Account Details"** với các thông tin:

```
Cloud Name: your-cloud-name
API Key: 123456789012345
API Secret: ****************** (click "View" để xem)
```

### 2.3. Copy Cloud Name

1. Tìm dòng **"Cloud Name"**
2. Copy giá trị (ví dụ: `dnslrwedn`)
3. Lưu vào notepad

```
Cloud Name: dnslrwedn
```

### 2.4. Copy API Key

1. Tìm dòng **"API Key"**
2. Click vào icon **Copy** hoặc select và copy
3. Lưu vào notepad

```
API Key: 942749116916526
```

### 2.5. Copy API Secret

1. Tìm dòng **"API Secret"**
2. Click **"View API Secret"** (hoặc icon con mắt)
3. Click **"Copy"** hoặc select và copy
4. Lưu vào notepad

```
API Secret: wZlZ_IVgBacQfPgOgtQEawALflc
```

> **🔒 BẢO MẬT:**
> - API Secret giống như mật khẩu
> - **KHÔNG BAO GIỜ** commit lên GitHub
> - **KHÔNG** share công khai

### 2.6. Tổng hợp thông tin

Bạn cần 3 thông tin này:

```
Cloud Name: dnslrwedn
API Key: 942749116916526
API Secret: wZlZ_IVgBacQfPgOgtQEawALflc
```

---

## ⚙️ Bước 3: Cấu hình .env

### 3.1. Mở file .env

Mở file `.env` trong thư mục gốc của dự án.

### 3.2. Cập nhật Cloudinary Configuration

Tìm mục Cloudinary và cập nhật:

```env
# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
```

### 3.3. Ví dụ cụ thể

```env
# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=dnslrwedn
CLOUDINARY_API_KEY=942749116916526
CLOUDINARY_API_SECRET=wZlZ_IVgBacQfPgOgtQEawALflc
```

**Giải thích:**

| Biến | Mô tả | Ví dụ |
|------|-------|-------|
| `CLOUDINARY_CLOUD_NAME` | Tên cloud của bạn | `dnslrwedn` |
| `CLOUDINARY_API_KEY` | Public key | `942749116916526` |
| `CLOUDINARY_API_SECRET` | Private key (giữ bí mật!) | `wZlZ_IVgBacQfPgOgtQEawALflc` |

### 3.4. Lưu file

- Lưu file `.env`
- **KHÔNG** commit lên GitHub

---

## ✅ Bước 4: Test Upload

### 4.1. Khởi động Server

```bash
go run main.go
```

### 4.2. Test Upload bằng Postman

#### Tạo Request Upload

1. Mở Postman
2. Tạo request mới:
   - Method: **POST**
   - URL: `http://localhost:8080/api/v1/cloudinary/upload`

#### Setup Headers

```
Authorization: Bearer <your-jwt-token>
```

> **Lưu ý:** Bạn cần login trước để lấy JWT token

#### Setup Body

1. Chọn tab **Body**
2. Chọn **form-data**
3. Thêm field:
   - Key: `file` (chọn type là **File**)
   - Value: Chọn file ảnh từ máy (jpg, png, etc.)

#### Send Request

Click **Send**

#### Kết quả mong đợi

```json
{
  "success": true,
  "url": "https://res.cloudinary.com/dnslrwedn/image/upload/v1234567890/ecommerce/abc123.jpg",
  "public_id": "ecommerce/abc123"
}
```

### 4.3. Test Upload bằng cURL

```bash
curl -X POST http://localhost:8080/api/v1/cloudinary/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/image.jpg"
```

### 4.4. Kiểm tra trên Cloudinary Dashboard

1. Quay lại Cloudinary Dashboard
2. Click vào **Media Library** (menu bên trái)
3. Bạn sẽ thấy ảnh vừa upload
4. Click vào ảnh để xem chi tiết

---

## 📂 Quản lý Media

### 5.1. Media Library

**Truy cập:** Dashboard → Media Library

**Chức năng:**

- ✅ Xem tất cả ảnh đã upload
- ✅ Search ảnh theo tên, tag
- ✅ Tạo folder để organize
- ✅ Xóa ảnh không dùng
- ✅ Xem thông tin chi tiết (size, format, URL)

### 5.2. Tạo Folder

1. Trong Media Library, click **"Add folder"**
2. Đặt tên folder:
   ```
   products/
   avatars/
   reviews/
   categories/
   ```

### 5.3. URL của ảnh

Sau khi upload, bạn nhận được URL dạng:

```
https://res.cloudinary.com/{cloud_name}/image/upload/{public_id}.{format}
```

**Ví dụ:**

```
https://res.cloudinary.com/dnslrwedn/image/upload/v1704012345/products/laptop-abc123.jpg
```

**Sử dụng URL này để:**
- Lưu vào database (field `image_url` của Product)
- Hiển thị ảnh trong frontend
- Share link ảnh

### 5.4. Transform ảnh (Resize, Crop)

Cloudinary cho phép transform ảnh ngay trong URL:

**Resize to 500x500:**

```
https://res.cloudinary.com/dnslrwedn/image/upload/w_500,h_500,c_fill/products/laptop.jpg
```

**Thumbnail 150x150:**

```
https://res.cloudinary.com/dnslrwedn/image/upload/w_150,h_150,c_thumb/products/laptop.jpg
```

**Parameters:**

- `w_500` - Width 500px
- `h_500` - Height 500px
- `c_fill` - Fill mode (crop để fit)
- `c_thumb` - Thumbnail mode
- `c_scale` - Scale mode (giữ tỷ lệ)
- `q_auto` - Auto quality optimization

### 5.5. Monitoring Usage

**Dashboard → Settings → Usage**

Xem:
- Storage used / 25GB
- Bandwidth used / 25GB
- Transformations used / 25,000

---

## 🐛 Troubleshooting

### ❌ Lỗi: "Invalid API credentials"

**Nguyên nhân:**

- Cloud Name, API Key, hoặc API Secret sai
- Copy thiếu ký tự

**Giải pháp:**

1. Quay lại Cloudinary Dashboard
2. Copy lại chính xác 3 giá trị:
   - Cloud Name
   - API Key
   - API Secret
3. Paste vào `.env` (không có khoảng trắng thừa)
4. Restart server: `Ctrl+C` và `go run main.go`

---

### ❌ Lỗi: "Upload failed"

**Nguyên nhân:**

- File quá lớn (max 10MB cho free plan)
- Format file không được hỗ trợ
- Network issue

**Giải pháp:**

1. Kiểm tra size file:
   ```
   Free plan: Max 10MB/file
   ```

2. Kiểm tra format được hỗ trợ:
   ```
   Images: jpg, png, gif, webp, svg
   ```

3. Compress ảnh trước khi upload:
   - Dùng TinyPNG.com
   - Photoshop "Save for Web"

---

### ❌ Lỗi: "Quota exceeded"

**Nguyên nhân:**

- Đã dùng hết quota của free plan
- Storage > 25GB
- Bandwidth > 25GB/tháng
- Transformations > 25,000/tháng

**Giải pháp:**

1. Kiểm tra usage trong Dashboard
2. Xóa ảnh không dùng
3. Nâng cấp plan (nếu cần)
4. Đợi đến tháng mới (quota reset)

---

### ❌ Ảnh không hiển thị

**Nguyên nhân:**

- URL sai
- Ảnh đã bị xóa trên Cloudinary
- CORS issue

**Giải pháp:**

1. Copy URL và mở trong browser mới
2. Kiểm tra ảnh còn tồn tại trên Media Library
3. Kiểm tra URL có đầy đủ không:
   ```
   ✅ https://res.cloudinary.com/dnslrwedn/image/upload/...
   ❌ /image/upload/... (thiếu domain)
   ```

---

## 🎨 Advanced Features

### Upload Preset (Optional)

**Upload Preset** cho phép upload không cần authentication (unsigned upload).

**Tạo Upload Preset:**

1. Dashboard → Settings → Upload
2. Scroll xuống **"Upload presets"**
3. Click **"Add upload preset"**
4. Cấu hình:
   - Preset name: `ecommerce_products`
   - Signing mode: **Unsigned**
   - Folder: `products`
5. Save

**Sử dụng:**

Frontend có thể upload trực tiếp lên Cloudinary mà không qua backend.

### Folder Organization

**Khuyến nghị cấu trúc folder:**

```
ecommerce/
├── products/
│   ├── electronics/
│   ├── fashion/
│   └── home/
├── avatars/
├── reviews/
└── categories/
```

**Lợi ích:**

- Dễ quản lý
- Dễ tìm kiếm
- Có thể set access control per folder

---

## 📊 Best Practices

### 1. Naming Convention

**Đặt tên file có ý nghĩa:**

```
✅ product-laptop-dell-xps-13.jpg
✅ avatar-user-12345.jpg
✅ category-electronics.png

❌ image1.jpg
❌ abc.png
❌ untitled.jpg
```

### 2. Optimize Images

**Trước khi upload:**

- Resize về kích thước phù hợp (max 2000x2000 cho product)
- Compress để giảm size
- Dùng format phù hợp:
  - **JPG** - Ảnh sản phẩm
  - **PNG** - Logo, icon (cần nền trong suốt)
  - **WebP** - Modern format (nhẹ hơn)

### 3. Use Transformations

**Thay vì upload nhiều size:**

Upload 1 ảnh gốc chất lượng cao, dùng URL transform:

```javascript
// Ảnh gốc
const originalUrl = "https://res.cloudinary.com/.../product.jpg"

// Thumbnail cho danh sách
const thumbnail = "https://res.cloudinary.com/.../w_300,h_300,c_fill/product.jpg"

// Full size cho chi tiết
const fullSize = "https://res.cloudinary.com/.../w_1200,q_auto/product.jpg"
```

### 4. Backup Important Images

- Export ảnh quan trọng định kỳ
- Lưu bản backup ở nơi khác
- Không dựa hoàn toàn vào một service

### 5. Clean Up Regularly

- Xóa ảnh test/demo không dùng
- Xóa ảnh sản phẩm đã ngừng bán
- Giữ storage dưới ngưỡng free plan

---

## 🔒 Security Tips

### 1. Bảo vệ API Secret

```env
# ✅ ĐÚNG - Trong .env (không commit)
CLOUDINARY_API_SECRET=wZlZ_IVgBacQfPgOgtQEawALflc

# ❌ SAI - Trong code (commit lên GitHub)
const apiSecret = "wZlZ_IVgBacQfPgOgtQEawALflc"
```

### 2. .gitignore

Đảm bảo `.env` có trong `.gitignore`:

```gitignore
.env
.env.local
.env.production
```

### 3. Separate Credentials

**Development:**

```env
CLOUDINARY_CLOUD_NAME=dev-cloud
```

**Production:**

```env
CLOUDINARY_CLOUD_NAME=prod-cloud
```

---

## 📚 Tài liệu tham khảo

- [Cloudinary Documentation](https://cloudinary.com/documentation)
- [Go SDK](https://cloudinary.com/documentation/go_integration)
- [Image Transformations](https://cloudinary.com/documentation/image_transformations)
- [Upload API](https://cloudinary.com/documentation/upload_images)

---

## ✅ Checklist

- [ ] Đã tạo tài khoản Cloudinary
- [ ] Đã copy Cloud Name
- [ ] Đã copy API Key
- [ ] Đã copy API Secret
- [ ] Đã cập nhật `.env`
- [ ] Đã test upload thành công
- [ ] Đã kiểm tra ảnh trên Media Library
- [ ] File `.env` không bị commit lên GitHub

---

**🎉 Hoàn thành!** Bây giờ ứng dụng của bạn đã có thể upload và quản lý hình ảnh!

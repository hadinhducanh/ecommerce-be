# 📧 Hướng dẫn Setup Gmail SMTP cho Email Service

Hướng dẫn này sẽ giúp bạn cấu hình Gmail để gửi email từ ứng dụng (OTP, verification, forgot password, etc.)

## 📋 Mục lục

- [Yêu cầu](#yêu-cầu)
- [Bước 1: Bật 2-Step Verification](#bước-1-bật-2-step-verification)
- [Bước 2: Tạo App Password](#bước-2-tạo-app-password)
- [Bước 3: Cấu hình .env](#bước-3-cấu-hình-env)
- [Bước 4: Test Email](#bước-4-test-email)
- [Troubleshooting](#troubleshooting)

---

## ✅ Yêu cầu

- Tài khoản Gmail (khuyến nghị tạo Gmail riêng cho dự án)
- Truy cập vào Google Account Settings

> **⚠️ Lưu ý quan trọng:**
> - **KHÔNG** sử dụng mật khẩu Gmail thông thường
> - **BẮT BUỘC** phải dùng App Password
> - App Password là mật khẩu 16 ký tự do Google tạo riêng cho ứng dụng

---

## 🔐 Bước 1: Bật 2-Step Verification

### 1.1. Truy cập Google Account

1. Mở trình duyệt và truy cập: **https://myaccount.google.com/**
2. Đăng nhập với tài khoản Gmail bạn muốn sử dụng

### 1.2. Vào Security Settings

1. Click vào **Security** (Bảo mật) ở menu bên trái
2. Hoặc truy cập trực tiếp: **https://myaccount.google.com/security**

### 1.3. Bật 2-Step Verification

1. Tìm mục **"How you sign in to Google"**
2. Click vào **"2-Step Verification"** (Xác minh 2 bước)
3. Click nút **"Get Started"** (Bắt đầu)
4. Nhập lại mật khẩu Gmail nếu được yêu cầu
5. Chọn phương thức xác minh:
   - **Text message (SMS)** - Nhận mã qua SMS (Khuyến nghị)
   - **Voice call** - Nhận mã qua cuộc gọi
   - **Authenticator app** - Dùng ứng dụng Google Authenticator
6. Nhập số điện thoại (nếu chọn SMS)
7. Nhập mã xác minh nhận được
8. Click **"Turn on"** để kích hoạt

> **✅ Xác nhận thành công:**
> Bạn sẽ thấy thông báo "2-Step Verification is on"

---

## 🔑 Bước 2: Tạo App Password

### 2.1. Truy cập App Passwords

**Cách 1: Từ Security Page**

1. Vẫn ở trang **Security**
2. Trong mục **"How you sign in to Google"**
3. Click vào **"2-Step Verification"**
4. Scroll xuống tìm **"App passwords"**
5. Click vào **"App passwords"**

**Cách 2: Truy cập trực tiếp**

- Truy cập: **https://myaccount.google.com/apppasswords**
- Đăng nhập nếu được yêu cầu

> **⚠️ Lưu ý:**
> - Nếu không thấy "App passwords", bạn chưa bật 2-Step Verification
> - Quay lại Bước 1 để bật 2-Step Verification

### 2.2. Tạo App Password mới

1. Click vào dropdown **"Select app"**
2. Chọn **"Mail"**
3. Click vào dropdown **"Select device"**
4. Chọn **"Other (Custom name)"**
5. Nhập tên: `Ecommerce Backend` hoặc tên bạn muốn
6. Click **"Generate"**

### 2.3. Lưu App Password

Google sẽ hiển thị App Password gồm **16 ký tự** (có dấu cách):

```
Ví dụ: abcd efgh ijkl mnop
```

> **🚨 CỰC KỲ QUAN TRỌNG:**
> - Copy App Password này ngay lập tức
> - Lưu vào nơi an toàn
> - Sau khi đóng popup này, bạn **KHÔNG THỂ** xem lại được
> - Nếu mất, phải tạo App Password mới

**Lưu lại theo định dạng:**

```
App Password: abcd efgh ijkl mnop
(hoặc không có dấu cách: abcdefghijklmnop)
```

---

## ⚙️ Bước 3: Cấu hình .env

### 3.1. Mở file .env

Mở file `.env` trong thư mục gốc của dự án.

### 3.2. Cập nhật SMTP Configuration

Thay đổi các giá trị sau:

```env
# SMTP Configuration for Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=abcdefghijklmnop
```

**Giải thích:**

| Biến | Giá trị | Mô tả |
|------|---------|-------|
| `SMTP_HOST` | `smtp.gmail.com` | Gmail SMTP server (giữ nguyên) |
| `SMTP_PORT` | `587` | Port TLS (giữ nguyên) |
| `SMTP_USER` | `your-email@gmail.com` | Gmail của bạn |
| `SMTP_PASS` | `abcdefghijklmnop` | App Password (16 ký tự, **KHÔNG CÓ** dấu cách) |

### 3.3. Ví dụ cụ thể

```env
# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=myapp.ecommerce@gmail.com
SMTP_PASS=xyzw1234abcd5678
```

> **⚠️ Lưu ý:**
> - Xóa **TẤT CẢ** dấu cách trong App Password
> - Nếu App Password là `abcd efgh ijkl mnop`
> - Thì `SMTP_PASS=abcdefghijklmnop`

### 3.4. Lưu file

- Lưu file `.env`
- **KHÔNG** commit file `.env` lên GitHub

---

## ✅ Bước 4: Test Email

### 4.1. Khởi động Server

```bash
go run main.go
```

### 4.2. Test gửi OTP

**Sử dụng Postman hoặc cURL:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@1234",
    "full_name": "Test User"
  }'
```

### 4.3. Kiểm tra email

1. Kiểm tra hộp thư đến của email bạn vừa đăng ký
2. Tìm email từ địa chỉ `SMTP_USER` bạn đã cấu hình
3. Email sẽ chứa mã OTP 6 số

### 4.4. Kiểm tra logs

Xem logs của server:

```
✅ Email sent successfully to: test@example.com
```

Nếu thành công, email service đã hoạt động!

---

## 🐛 Troubleshooting

### ❌ Lỗi: "535-5.7.8 Username and Password not accepted"

**Nguyên nhân:**

- Sử dụng mật khẩu Gmail thông thường thay vì App Password
- App Password sai
- Chưa bật 2-Step Verification

**Giải pháp:**

1. Kiểm tra lại 2-Step Verification đã bật chưa
2. Tạo lại App Password mới
3. Copy chính xác App Password (không có dấu cách)
4. Cập nhật lại file `.env`

---

### ❌ Lỗi: "Could not send email"

**Nguyên nhân:**

- SMTP host hoặc port sai
- Gmail chặn truy cập từ "less secure apps"

**Giải pháp:**

1. Kiểm tra lại cấu hình:
   ```env
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   ```

2. Kiểm tra 2-Step Verification đã bật

3. Xóa các App Password cũ và tạo mới:
   - Truy cập https://myaccount.google.com/apppasswords
   - Xóa các App Password cũ
   - Tạo App Password mới

---

### ❌ Lỗi: "App passwords is not available"

**Nguyên nhân:**

- 2-Step Verification chưa được bật
- Tài khoản Google Workspace có policy hạn chế

**Giải pháp:**

1. Bật 2-Step Verification (xem Bước 1)
2. Đợi vài phút sau khi bật 2-Step Verification
3. Refresh trang và thử lại
4. Nếu vẫn không được, sử dụng tài khoản Gmail cá nhân khác

---

### ❌ Email vào Spam

**Nguyên nhân:**

- Email gửi từ Gmail cá nhân thường bị đánh dấu spam

**Giải pháp tạm thời:**

1. Kiểm tra thư mục Spam
2. Đánh dấu email "Not spam"
3. Thêm địa chỉ email vào danh bạ

**Giải pháp lâu dài (Production):**

- Sử dụng dịch vụ email chuyên nghiệp:
  - SendGrid
  - AWS SES
  - Mailgun
  - Postmark

---

### ❌ Lỗi: "Timed out"

**Nguyên nhân:**

- Firewall chặn port 587
- Mạng công ty chặn SMTP

**Giải pháp:**

1. Thử port khác:
   ```env
   SMTP_PORT=465  # SSL
   ```

2. Kiểm tra firewall

3. Thử mạng khác (không phải mạng công ty)

---

## 🔒 Best Practices

### 1. Bảo mật

- ✅ Sử dụng Gmail riêng cho dự án (không dùng Gmail cá nhân)
- ✅ Không commit file `.env` lên GitHub
- ✅ Thêm `.env` vào `.gitignore`
- ✅ Sử dụng environment variables khác nhau cho dev/staging/production
- ✅ Thay đổi App Password định kỳ

### 2. Tạo Gmail riêng cho dự án

**Khuyến nghị tạo Gmail mới:**

```
Ví dụ:
- myapp.noreply@gmail.com
- myapp.notification@gmail.com
- ecommerce.backend@gmail.com
```

**Lợi ích:**

- Dễ quản lý
- Tránh lộ email cá nhân
- Có thể thu hồi quyền truy cập dễ dàng

### 3. Giới hạn gửi email

Gmail có giới hạn:

- **500 emails/ngày** cho Gmail cá nhân
- **2000 emails/ngày** cho Google Workspace

**Giải pháp nếu vượt quá:**

- Sử dụng SendGrid (100 emails/day miễn phí)
- AWS SES (pay-as-you-go)

---

## 📚 Tài liệu tham khảo

- [Google App Passwords](https://support.google.com/accounts/answer/185833)
- [Gmail SMTP Settings](https://support.google.com/mail/answer/7126229)
- [2-Step Verification](https://www.google.com/landing/2step/)

---

## ✅ Checklist

Kiểm tra lại trước khi sử dụng:

- [ ] Đã bật 2-Step Verification
- [ ] Đã tạo App Password
- [ ] Đã cập nhật `SMTP_USER` trong `.env`
- [ ] Đã cập nhật `SMTP_PASS` (không có dấu cách)
- [ ] Đã test gửi email thành công
- [ ] File `.env` không bị commit lên GitHub

---

**🎉 Hoàn thành!** Bây giờ ứng dụng của bạn đã có thể gửi email!

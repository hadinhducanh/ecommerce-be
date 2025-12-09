# Hướng dẫn Test Category API trên Postman

## Chuẩn bị

1. **Base URL**: `http://localhost:8080/api/v1` (hoặc port của bạn)
2. **Headers mặc định** (nếu cần):
   - `Content-Type: application/json`
   - `Authorization: Bearer <token>` (cho các API Admin)

---

## 1. View Danh mục Cha (Parent Categories)

### 1.1. Lấy danh sách Parent Categories (cho Dropdown Filter)

**Method:** `GET`  
**URL:** `http://localhost:8080/api/v1/categories/parents`

**Query Parameters (Optional):**
- `language`: `vi` hoặc `en` (mặc định: `vi`)
- `includeInactive`: `true` hoặc `false` (mặc định: `false`)

**Ví dụ:**
```
GET http://localhost:8080/api/v1/categories/parents?includeInactive=false&language=vi
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Điện thoại",
      "nameEn": "Smartphones",
      "parentId": null,
      "childrenIds": [3, 4, 5, 6, 7, 8, 9],
      "isActive": true,
      ...
    },
    {
      "id": 2,
      "name": "Laptop",
      "nameEn": "Laptops",
      "parentId": null,
      "childrenIds": [10, 11, 12, 13, 14, 15, 16, 17],
      "isActive": true,
      ...
    }
  ]
}
```

**Lưu ý:**
- Chỉ trả về các parent categories (root categories)
- Dùng để populate dropdown filter trên page quản lý child categories

---

### 1.2. Tìm kiếm Parent Categories (với Pagination, Filter, Search)

**Method:** `POST`  
**URL:** `http://localhost:8080/api/v1/categories/search`

**Headers:**
```
Content-Type: application/json
```

**Body (JSON):**
```json
{
  "name": null,
  "isActive": true,
  "sortBy": "createdAt",
  "sortOrder": "DESC",
  "page": 1,
  "limit": 10
}
```

**Ví dụ 1: Lấy tất cả parent categories (body rỗng)**
```json
{}
```

**Ví dụ 2: Search theo tên**
```json
{
  "name": "Điện",
  "isActive": true,
  "page": 1,
  "limit": 10
}
```

**Ví dụ 3: Filter theo status**
```json
{
  "isActive": true,
  "page": 1,
  "limit": 10
}
```

**Ví dụ 4: Lấy tất cả (active + inactive)**
```json
{
  "isActive": [true, false],
  "page": 1,
  "limit": 10
}
```

**Query Parameters (Optional):**
- `language`: `vi` hoặc `en` (mặc định: `vi`)

**Response:**
```json
{
  "success": true,
  "message": "Tìm kiếm danh mục thành công",
  "data": [
    {
      "id": 1,
      "name": "Điện thoại",
      "nameEn": "Smartphones",
      "parentId": null,
      "childrenIds": [3, 4, 5, 6, 7, 8, 9],
      "isActive": true,
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    },
    {
      "id": 2,
      "name": "Laptop",
      "nameEn": "Laptops",
      "parentId": null,
      "childrenIds": [10, 11, 12, 13, 14, 15, 16, 17],
      "isActive": true,
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    }
  ],
  "total": 2,
  "page": 1,
  "limit": 10,
  "totalPages": 1
}
```

**Lưu ý:**
- Chỉ trả về parent categories (không có child categories)
- `parentId` luôn là `null`
- `childrenIds` chứa danh sách ID các children

---

### 1.3. Lấy một Parent Category theo ID

**Method:** `GET`  
**URL:** `http://localhost:8080/api/v1/categories/1`

**Query Parameters (Optional):**
- `includeInactive`: `true` hoặc `false` (mặc định: `true`)
- `language`: `vi` hoặc `en` (mặc định: `vi`)

**Ví dụ:**
```
GET http://localhost:8080/api/v1/categories/1?includeInactive=true&language=vi
```

**Response:**
```json
{
  "success": true,
  "message": "Lấy thông tin danh mục thành công",
  "data": {
    "id": 1,
    "name": "Điện thoại",
    "nameEn": "Smartphones",
    "description": "Các dòng điện thoại thông minh...",
    "descriptionEn": "Modern smartphones...",
    "image": null,
    "isActive": true,
    "parentId": null,
    "childrenIds": [3, 4, 5, 6, 7, 8, 9],
    "createdAt": "2025-12-08T19:11:13Z",
    "updatedAt": "2025-12-08T19:11:13Z"
  }
}
```

---

## 2. View Danh mục Con (Child Categories)

### 2.1. Tìm kiếm Child Categories (với Filter theo Parent)

**Method:** `POST`  
**URL:** `http://localhost:8080/api/v1/categories/children/search`

**Headers:**
```
Content-Type: application/json
```

**Body (JSON):**

**Ví dụ 1: Lấy tất cả child categories (body rỗng)**
```json
{}
```

**Ví dụ 2: Filter theo Parent ID (quan trọng!)**
```json
{
  "parentId": 1,
  "isActive": true,
  "page": 1,
  "limit": 10
}
```

**Ví dụ 3: Filter theo Parent + Search theo tên**
```json
{
  "parentId": 1,
  "name": "iPhone",
  "isActive": true,
  "page": 1,
  "limit": 10
}
```

**Ví dụ 4: Lấy tất cả children (không filter parent)**
```json
{
  "isActive": true,
  "page": 1,
  "limit": 10
}
```

**Ví dụ 5: Filter theo status**
```json
{
  "parentId": 1,
  "isActive": false,
  "page": 1,
  "limit": 10
}
```

**Query Parameters (Optional):**
- `language`: `vi` hoặc `en` (mặc định: `vi`)

**Response:**
```json
{
  "success": true,
  "message": "Tìm kiếm danh mục con thành công",
  "data": [
    {
      "id": 3,
      "name": "iPhone",
      "nameEn": "iPhone",
      "parentId": 1,
      "childrenIds": [],
      "isActive": true,
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    },
    {
      "id": 4,
      "name": "Samsung",
      "nameEn": "Samsung",
      "parentId": 1,
      "childrenIds": [],
      "isActive": true,
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    }
  ],
  "total": 7,
  "page": 1,
  "limit": 10,
  "totalPages": 1
}
```

**Lưu ý:**
- Chỉ trả về child categories (không có parent categories)
- `parentId` luôn có giá trị (không phải `null`)
- `childrenIds` thường là mảng rỗng `[]`

---

### 2.2. Lấy danh sách Children của một Parent cụ thể

**Method:** `GET`  
**URL:** `http://localhost:8080/api/v1/categories/1/children`

**Query Parameters (Optional):**
- `includeInactive`: `true` hoặc `false` (mặc định: `true`)
- `language`: `vi` hoặc `en` (mặc định: `vi`)

**Ví dụ:**
```
GET http://localhost:8080/api/v1/categories/1/children?includeInactive=false&language=vi
```

**Response:**
```json
{
  "success": true,
  "message": "Lấy danh sách danh mục con thành công",
  "data": [
    {
      "id": 3,
      "name": "iPhone",
      "nameEn": "iPhone",
      "parentId": 1,
      "childrenIds": [],
      ...
    },
    {
      "id": 4,
      "name": "Samsung",
      "nameEn": "Samsung",
      "parentId": 1,
      "childrenIds": [],
      ...
    }
  ]
}
```

---

## 3. Workflow Test trên Postman

### Bước 1: Lấy danh sách Parent Categories

1. **Tạo request mới:**
   - Method: `GET`
   - URL: `http://localhost:8080/api/v1/categories/parents`
   - Query params: `includeInactive=false`

2. **Send request**
3. **Kiểm tra response:**
   - Phải có `parentId: null` cho tất cả items
   - Phải có `childrenIds` là mảng (có thể rỗng hoặc có giá trị)

---

### Bước 2: Lấy danh sách Child Categories (không filter)

1. **Tạo request mới:**
   - Method: `POST`
   - URL: `http://localhost:8080/api/v1/categories/children/search`
   - Headers: `Content-Type: application/json`
   - Body (raw JSON): `{}`

2. **Send request**
3. **Kiểm tra response:**
   - Tất cả items phải có `parentId` khác `null`
   - `childrenIds` thường là mảng rỗng `[]`

---

### Bước 3: Filter Child Categories theo Parent

1. **Lấy Parent ID từ Bước 1** (ví dụ: `id: 1` - "Điện thoại")

2. **Tạo request mới:**
   - Method: `POST`
   - URL: `http://localhost:8080/api/v1/categories/children/search`
   - Headers: `Content-Type: application/json`
   - Body (raw JSON):
   ```json
   {
     "parentId": 1,
     "isActive": true,
     "page": 1,
     "limit": 10
   }
   ```

3. **Send request**
4. **Kiểm tra response:**
   - Tất cả items phải có `parentId: 1`
   - Chỉ có children của "Điện thoại" (iPhone, Samsung, ...)

---

### Bước 4: So sánh Parent và Child

**Test 1: Parent Search**
```
POST /api/v1/categories/search
Body: {}
```
→ Kết quả: Chỉ có "Điện thoại", "Laptop" (parentId = null)

**Test 2: Child Search**
```
POST /api/v1/categories/children/search
Body: {}
```
→ Kết quả: Chỉ có "iPhone", "Samsung", "MacBook", ... (parentId ≠ null)

**Test 3: Child Search với Filter**
```
POST /api/v1/categories/children/search
Body: { "parentId": 1 }
```
→ Kết quả: Chỉ có "iPhone", "Samsung", ... (children của "Điện thoại")

---

## 4. Checklist Test

### ✅ Test Parent Categories

- [ ] `GET /categories/parents` → Trả về danh sách parent
- [ ] `POST /categories/search` với body rỗng → Trả về parent categories
- [ ] `POST /categories/search` với `name` → Search parent theo tên
- [ ] `POST /categories/search` với `isActive: true` → Chỉ parent active
- [ ] `GET /categories/1` → Lấy parent category theo ID
- [ ] Kiểm tra: Tất cả có `parentId: null`

### ✅ Test Child Categories

- [ ] `POST /categories/children/search` với body rỗng → Trả về tất cả children
- [ ] `POST /categories/children/search` với `parentId: 1` → Chỉ children của parent ID 1
- [ ] `POST /categories/children/search` với `parentId: 1` và `name: "iPhone"` → Search children
- [ ] `GET /categories/1/children` → Lấy children của parent ID 1
- [ ] Kiểm tra: Tất cả có `parentId` khác `null`

### ✅ Test Phân biệt

- [ ] Parent search không trả về child categories
- [ ] Child search không trả về parent categories
- [ ] Filter `parentId` trong child search hoạt động đúng

---

## 5. Ví dụ Request Collection cho Postman

### Collection: Category API Tests

#### Folder: Parent Categories
1. **Get All Parents (Simple)**
   - Method: `GET`
   - URL: `{{baseUrl}}/categories/parents?includeInactive=false`

2. **Search Parents (Pagination)**
   - Method: `POST`
   - URL: `{{baseUrl}}/categories/search`
   - Body:
   ```json
   {
     "isActive": true,
     "page": 1,
     "limit": 10
   }
   ```

3. **Get Parent by ID**
   - Method: `GET`
   - URL: `{{baseUrl}}/categories/1`

#### Folder: Child Categories
1. **Get All Children**
   - Method: `POST`
   - URL: `{{baseUrl}}/categories/children/search`
   - Body: `{}`

2. **Filter Children by Parent**
   - Method: `POST`
   - URL: `{{baseUrl}}/categories/children/search`
   - Body:
   ```json
   {
     "parentId": 1,
     "isActive": true,
     "page": 1,
     "limit": 10
   }
   ```

3. **Get Children of Parent**
   - Method: `GET`
   - URL: `{{baseUrl}}/categories/1/children`

---

## 6. Troubleshooting

### Vấn đề: Không thấy child categories

**Nguyên nhân:**
- Chưa chạy seed script
- Chưa tạo quan hệ parent-child

**Giải pháp:**
1. Chạy seed script: `go run cmd/seed/seed.go`
2. Hoặc tạo category và thêm vào parent thủ công

### Vấn đề: Parent search trả về child categories

**Nguyên nhân:**
- Logic filter chưa đúng

**Giải pháp:**
- Kiểm tra lại query: `WHERE id NOT IN (SELECT DISTINCT child_id FROM category_children)`

### Vấn đề: Child search trả về parent categories

**Nguyên nhân:**
- Logic filter chưa đúng

**Giải pháp:**
- Kiểm tra lại query: `WHERE id IN (SELECT DISTINCT child_id FROM category_children)`

---

## 7. Tips

1. **Sử dụng Environment Variables trong Postman:**
   - `baseUrl`: `http://localhost:8080/api/v1`
   - `token`: `<your_admin_token>`

2. **Lưu Parent ID từ response đầu tiên** để dùng cho các request sau

3. **Sử dụng Tests tab trong Postman** để tự động kiểm tra:
   ```javascript
   pm.test("Response has parentId null", function () {
       var jsonData = pm.response.json();
       jsonData.data.forEach(item => {
           pm.expect(item.parentId).to.be.null;
       });
   });
   ```

4. **Sử dụng Pre-request Script** để tự động lấy parent ID:
   ```javascript
   // Lưu parent ID từ request trước
   pm.environment.set("parentId", pm.response.json().data[0].id);
   ```

---

## 8. Quick Reference

| Mục đích | Method | Endpoint | Body |
|----------|--------|----------|------|
| Lấy danh sách parent (dropdown) | GET | `/categories/parents` | - |
| Tìm kiếm parent | POST | `/categories/search` | `{}` hoặc filter |
| Lấy parent theo ID | GET | `/categories/:id` | - |
| Tìm kiếm child | POST | `/categories/children/search` | `{}` hoặc filter |
| Filter child theo parent | POST | `/categories/children/search` | `{"parentId": 1}` |
| Lấy children của parent | GET | `/categories/:id/children` | - |

---

**Lưu ý:** Đảm bảo server đang chạy và database đã có dữ liệu (chạy seed script nếu cần).


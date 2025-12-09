# Category API Documentation

## Tổng quan

Category API hỗ trợ cấu trúc phân cấp (parent-child) với cấu trúc mới:
- **Category mặc định là root category** (danh mục cha) khi tạo mới
- **Quan hệ parent-child** được lưu trong bảng riêng `category_children`
- **Quản lý children** thông qua các API riêng (thêm/xóa children)

Ví dụ: "Điện thoại" là parent, "iPhone", "Samsung" là children.

## Cấu trúc dữ liệu

### Category Model

```json
{
  "id": 1,
  "name": "Điện thoại",
  "nameEn": "Smartphones",
  "description": "Mô tả tiếng Việt",
  "descriptionEn": "English description",
  "image": "https://...",
  "isActive": true,
  "parentId": null,        // ID của category cha (null nếu là root category)
  "childrenIds": [3, 4, 5], // Danh sách ID các category con
  "createdAt": "2025-12-08T19:11:13Z",
  "updatedAt": "2025-12-08T19:11:13Z"
}
```

### Quan hệ Parent-Child

- **Root Category**: Category không có parent (`parentId = null`)
  - Ví dụ: "Điện thoại", "Laptop"
  - Khi tạo category mới, mặc định là root category
  
- **Child Category**: Category có parent (`parentId = <ID của parent>`)
  - Ví dụ: "iPhone" có `parentId = 1` (ID của "Điện thoại")
  - Để tạo child, có 2 cách:
    1. **Cách 1 (Đơn giản)**: Truyền `parentId` khi tạo category → Tự động thêm vào parent
    2. **Cách 2**: Tạo category trước (mặc định là root), sau đó dùng API `POST /categories/:parentId/children` để thêm vào parent

- **Quan hệ được lưu trong bảng `category_children`**:
  - Một child chỉ có thể thuộc về một parent
  - Một parent có thể có nhiều children

## API Endpoints

### 1. Lấy danh sách Parent Categories (cho Dropdown Filter) (Public)

**Endpoint:** `GET /api/v1/categories/parents`

**Query Parameters:**
- `language`: `vi` (mặc định) hoặc `en` - Ngôn ngữ hiển thị
- `includeInactive`: `true` hoặc `false` (mặc định: `false`) - Có lấy categories inactive không

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
      "childrenIds": [3, 4, 5],
      ...
    },
    {
      "id": 2,
      "name": "Laptop",
      "nameEn": "Laptops",
      "parentId": null,
      "childrenIds": [9, 10, 11],
      ...
    }
  ]
}
```

**Mục đích:** Dùng để populate dropdown filter trên page quản lý child categories.

### 2. Lấy danh sách tất cả Child Categories (cho Dropdown Filter) (Public)

**Endpoint:** `GET /api/v1/categories/children`

**Query Parameters:**
- `language`: `vi` (mặc định) hoặc `en` - Ngôn ngữ hiển thị
- `includeInactive`: `true` hoặc `false` (mặc định: `true`) - Có lấy categories inactive không. Mặc định `true` để có thể view và kích hoạt lại category inactive

**Lưu ý:** 
- Mặc định cho phép view category inactive (không cần query parameter)
- Chỉ truyền `includeInactive=false` nếu muốn chỉ lấy category active

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 3,
      "name": "iPhone",
      "nameEn": "iPhone",
      "parentId": 1,
      "childrenIds": [],
      "isActive": true,
      ...
    },
    {
      "id": 4,
      "name": "Samsung",
      "nameEn": "Samsung",
      "parentId": 1,
      "childrenIds": [],
      "isActive": true,
      ...
    },
    {
      "id": 9,
      "name": "MacBook",
      "nameEn": "MacBook",
      "parentId": 2,
      "childrenIds": [],
      "isActive": true,
      ...
    }
  ]
}
```

**Mục đích:** Dùng để populate dropdown khi add product hoặc các form cần chọn child category. Tương tự như `GET /categories/parents` nhưng trả về tất cả child categories thay vì parent categories.

**Ví dụ sử dụng:**
```javascript
// Lấy tất cả child categories active cho dropdown
const response = await fetch('/api/v1/categories/children?includeInactive=false&language=vi');
const { data: childCategories } = await response.json();

// Sử dụng trong dropdown
<select>
  {childCategories.map(child => (
    <option key={child.id} value={child.id}>{child.name}</option>
  ))}
</select>
```

### 3. Tìm kiếm Categories (Public)

**Endpoint:** `POST /api/v1/categories/search`

**Lưu ý quan trọng:** 
- Endpoint này **CHỈ trả về các root categories (parent categories)**
- Các child categories sẽ **KHÔNG xuất hiện** trong kết quả này
- Để lấy child categories, sử dụng `POST /api/v1/categories/children/search`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Điện thoại",        // Optional: Tìm kiếm theo tên (partial match)
  "isActive": true,             // Optional: true, false, hoặc [true, false] để lấy tất cả
  "sortBy": "createdAt",        // Optional: id, name, createdAt, updatedAt
  "sortOrder": "DESC",          // Optional: ASC, DESC
  "page": 1,                    // Optional: Mặc định 1
  "limit": 10                   // Optional: Mặc định 10
}
```

**Lấy tất cả root categories:**
```json
{}  // Body rỗng hoặc null
```

**Query Parameters:**
- `language`: `vi` (mặc định) hoặc `en` - Ngôn ngữ hiển thị

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
      "description": "...",
      "descriptionEn": "...",
      "image": null,
      "isActive": true,
      "parentId": null,
      "childrenIds": [3, 4, 5],
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    },
    {
      "id": 3,
      "name": "iPhone",
      "nameEn": "iPhone",
      "description": "...",
      "descriptionEn": "...",
      "image": null,
      "isActive": true,
      "parentId": 1,
      "childrenIds": [],
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    }
  ],
  "total": 17,
  "page": 1,
  "limit": 10,
  "totalPages": 2
}
```

### 4. Tìm kiếm Category Children (Public)

**Endpoint:** `POST /api/v1/categories/children/search`

**Lưu ý quan trọng:**
- Endpoint này **CHỈ trả về các child categories**
- Các root categories (parent) sẽ **KHÔNG xuất hiện** trong kết quả này
- Để lấy root categories, sử dụng `POST /api/v1/categories/search`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "parentId": 1,                // Optional: Filter theo parent ID (exact match)
  "name": "iPhone",             // Optional: Tìm kiếm theo tên (partial match)
  "isActive": true,              // Optional: true, false, hoặc [true, false] để lấy tất cả
  "sortBy": "createdAt",        // Optional: id, name, createdAt, updatedAt
  "sortOrder": "DESC",          // Optional: ASC, DESC
  "page": 1,                    // Optional: Mặc định 1
  "limit": 10                   // Optional: Mặc định 10
}
```

**Lấy tất cả children:**
```json
{}  // Body rỗng hoặc null
```

**Query Parameters:**
- `language`: `vi` (mặc định) hoặc `en` - Ngôn ngữ hiển thị

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
      "description": "...",
      "descriptionEn": "...",
      "image": null,
      "isActive": true,
      "parentId": 1,
      "childrenIds": [],
      "createdAt": "2025-12-08T19:11:13Z",
      "updatedAt": "2025-12-08T19:11:13Z"
    },
    {
      "id": 4,
      "name": "Samsung",
      "nameEn": "Samsung",
      "description": "...",
      "descriptionEn": "...",
      "image": null,
      "isActive": true,
      "parentId": 1,
      "childrenIds": [],
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

### 5. Lấy một Category theo ID (Public)

**Endpoint:** `GET /api/v1/categories/:id`

**Query Parameters:**
- `includeInactive`: `true` hoặc `false` (mặc định: `true`) - Có lấy category inactive không. Mặc định `true` để có thể view và kích hoạt lại category inactive
- `language`: `vi` (mặc định) hoặc `en` - Ngôn ngữ hiển thị

**Lưu ý:** 
- Mặc định cho phép view category inactive (không cần query parameter)
- Chỉ truyền `includeInactive=false` nếu muốn chỉ lấy category active

**Response:**
```json
{
  "success": true,
  "message": "Lấy thông tin danh mục thành công",
  "data": {
    "id": 1,
    "name": "Điện thoại",
    "nameEn": "Smartphones",
    "description": "...",
    "descriptionEn": "...",
    "image": null,
    "isActive": true,
    "parentId": null,
    "childrenIds": [3, 4, 5],
    "createdAt": "2025-12-08T19:11:13Z",
    "updatedAt": "2025-12-08T19:11:13Z"
  }
}
```

### 6. Lấy danh sách Children của một Parent (Public)

**Endpoint:** `GET /api/v1/categories/:id/children`

**Query Parameters:**
- `includeInactive`: `true` hoặc `false` (mặc định: `true`) - Có lấy children inactive không
- `language`: `vi` (mặc định) hoặc `en` - Ngôn ngữ hiển thị

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

### 7. Tạo Category mới (Admin Only)

**Endpoint:** `POST /api/v1/categories`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Điện thoại",           // Required
  "nameEn": "Smartphones",         // Optional
  "description": "Mô tả...",      // Optional
  "descriptionEn": "Description", // Optional
  "image": "https://...",         // Optional
  "isActive": true,                // Optional, mặc định true
  "parentId": null                 // Optional: Nếu có thì tự động thêm vào parent, nếu null thì tạo như root category
}
```

**Tạo root category (parent):**
```json
{
  "name": "Điện thoại",
  "parentId": null  // hoặc không truyền
}
```

**Tạo child category (tự động thêm vào parent):**
```json
{
  "name": "iPhone",
  "parentId": 1  // ID của category "Điện thoại"
}
```

**Lưu ý quan trọng:**
- Nếu không truyền `parentId` hoặc `parentId = null`: Category mới tạo **mặc định là root category**
- Nếu truyền `parentId`: Category sẽ được tạo và **tự động thêm vào parent** (tạo quan hệ trong bảng `category_children`)
- Nếu `parentId` không tồn tại hoặc inactive: Sẽ trả về lỗi và không tạo category

**Response:**
```json
{
  "success": true,
  "message": "Tạo danh mục thành công",
  "data": {
    "id": 1,
    "name": "Điện thoại",
    "nameEn": "Smartphones",
    "description": "...",
    "descriptionEn": "...",
    "image": null,
    "isActive": true,
    "parentId": null,
    "childrenIds": [],
    "createdAt": "2025-12-08T19:11:13Z",
    "updatedAt": "2025-12-08T19:11:13Z"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "danh mục với tên này đã tồn tại"
}
```

### 8. Cập nhật Category (Partial Update - Admin Only)

**Endpoint:** `PATCH /api/v1/categories/:id`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Điện thoại mới",      // Optional
  "nameEn": "New Smartphones",   // Optional
  "description": "...",          // Optional
  "descriptionEn": "...",        // Optional
  "image": "https://...",         // Optional
  "isActive": false              // Optional
}
```

**Lưu ý:**
- Không thể đổi parent qua update
- Để thay đổi parent, sử dụng API quản lý children: `POST /categories/:newParentId/children` và `DELETE /categories/:oldParentId/children`

**Response:**
```json
{
  "success": true,
  "message": "Cập nhật danh mục thành công",
  "data": {
    "id": 1,
    "name": "Điện thoại mới",
    "nameEn": "New Smartphones",
    "description": "...",
    "descriptionEn": "...",
    "image": "https://...",
    "isActive": false,
    "parentId": null,
    "childrenIds": [3, 4, 5],
    "createdAt": "2025-12-08T19:11:13Z",
    "updatedAt": "2025-12-08T19:11:14Z"
  }
}
```

### 9. Thay thế Category (Full Update - Admin Only)

**Endpoint:** `PUT /api/v1/categories/:id`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Điện thoại",           // Required
  "nameEn": "Smartphones",        // Optional
  "description": "...",           // Optional
  "descriptionEn": "...",         // Optional
  "image": "https://...",         // Optional
  "isActive": true                // Optional
}
```

**Response:** Giống như PATCH

### 10. Xóa Category (Soft Delete - Admin Only)

**Endpoint:** `DELETE /api/v1/categories/:id`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Danh mục đã được xóa thành công"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "không thể xóa danh mục này vì có 5 sản phẩm liên quan. Vui lòng xóa hoặc di chuyển các sản phẩm sang danh mục khác trước."
}
```

### 11. Xóa vĩnh viễn Category (Hard Delete - Admin Only)

**Endpoint:** `DELETE /api/v1/categories/:id/hard`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Danh mục đã được xóa thành công"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "không thể xóa vĩnh viễn danh mục này vì có 5 sản phẩm liên quan. Vui lòng xử lý các sản phẩm trước."
}
```

### 12. Upload Image (Admin Only)

**Endpoint:** `POST /api/v1/categories/upload-image`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file`: File ảnh (required)
- `folder`: Folder trong Cloudinary (optional, mặc định: "categories")

**Response:**
```json
{
  "success": true,
  "message": "Upload ảnh thành công",
  "data": {
    "url": "https://res.cloudinary.com/...",
    "publicId": "categories/abc123"
  }
}
```

### 13. Delete Image (Admin Only)

**Endpoint:** `DELETE /api/v1/categories/delete-image`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "url": "https://res.cloudinary.com/..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Xóa ảnh thành công"
}
```

## Quản lý Category Children

### 14. Thêm Child vào Parent (Admin Only)

**Endpoint:** `POST /api/v1/categories/:id/children`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "childId": 3  // Required: ID của category con cần thêm
}
```

**Response:**
```json
{
  "success": true,
  "message": "Đã thêm danh mục con thành công"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "danh mục con này đã thuộc về một danh mục cha khác"
}
```

```json
{
  "success": false,
  "error": "danh mục không thể là parent của chính nó"
}
```

```json
{
  "success": false,
  "error": "không thể tạo circular reference: danh mục con đã là parent của danh mục cha này"
}
```

### 15. Xóa Child khỏi Parent (Admin Only)

**Endpoint:** `DELETE /api/v1/categories/:id/children`

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "childId": 3  // Required: ID của category con cần xóa
}
```

**Response:**
```json
{
  "success": true,
  "message": "Đã xóa danh mục con thành công"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "không tìm thấy quan hệ parent-child giữa danh mục 1 và 3"
}
```

## Business Rules & Logic

### 1. Cấu trúc phân cấp

- **Category mới tạo**: Mặc định là root category (không có parent)
- **Quan hệ parent-child**: Được lưu trong bảng `category_children`
- **Một child chỉ có một parent**: Unique constraint trên `child_id`
- **Một parent có nhiều children**: Không giới hạn

### 2. Validation Rules

**Khi tạo category:**
- `name` là required
- Category mới tạo mặc định là root category
- Không thể set parent khi tạo (phải dùng API quản lý children)

**Khi thêm child:**
- Parent và child phải tồn tại
- Child không thể là parent của chính nó
- Child không thể đã thuộc về parent khác
- Không cho phép circular reference

**Khi xóa category:**
- Soft delete: Không cho phép xóa nếu có products liên quan
- Hard delete: Không cho phép xóa nếu có products liên quan

### 3. Filter & Search Logic

**Phân biệt Categories và Children:**
- **`POST /categories/search`**: Chỉ trả về **root categories (parent)** - các categories không có trong `category_children` như child
- **`POST /categories/children/search`**: Chỉ trả về **child categories** - các categories có trong `category_children` như child

**Search theo name:**
- Partial match, case-insensitive
- Tìm kiếm cả tiếng Việt và tiếng Anh
- Không phân biệt dấu (normalize)

**Filter theo isActive:**
- `isActive: true` → Chỉ lấy active
- `isActive: false` → Chỉ lấy inactive
- `isActive: [true, false]` → Lấy tất cả
- Không truyền → Lấy tất cả

**Filter theo parentId (chỉ trong Children Search):**
- `parentId: 1` → Chỉ lấy children của parent có ID = 1
- Không truyền → Lấy tất cả children

**Sort:**
- `sortBy`: `id`, `name`, `createdAt`, `updatedAt`
- `sortOrder`: `ASC`, `DESC`
- Mặc định: `createdAt DESC`

**Pagination:**
- `page`: Mặc định 1
- `limit`: Mặc định 10, tối đa 1000

### 4. Multilingual Support

**Query Parameter:** `language=vi` hoặc `language=en`

- `vi` (mặc định): Hiển thị `name` và `description` (tiếng Việt)
- `en`: Nếu có `nameEn` và `descriptionEn`, sẽ hiển thị thay thế

**Ví dụ:**
```
GET /api/v1/categories/1?language=en
```

## Hướng dẫn cho Frontend - 2 Pages Quản lý

### Sử dụng cho Add Product - Dropdown Child Categories

**API sử dụng:** `GET /api/v1/categories/children`

**Ví dụ code:**
```javascript
// Lấy tất cả child categories active cho dropdown khi add product
const fetchChildCategoriesForDropdown = async () => {
  const response = await fetch('/api/v1/categories/children?includeInactive=false&language=vi', {
    method: 'GET'
  });
  
  const result = await response.json();
  return result.data; // Array of child categories
};

// Sử dụng trong form add product
const [childCategories, setChildCategories] = useState([]);

useEffect(() => {
  fetchChildCategoriesForDropdown().then(setChildCategories);
}, []);

// Render dropdown
<select name="categoryId" required>
  <option value="">Chọn danh mục</option>
  {childCategories.map(child => (
    <option key={child.id} value={child.id}>
      {child.name}
    </option>
  ))}
</select>
```

**Lưu ý:**
- Endpoint này trả về **tất cả child categories** (không phân biệt parent)
- Phù hợp cho dropdown khi add product vì product thường được gán vào child category cụ thể
- Nếu muốn filter theo parent, sử dụng `POST /api/v1/categories/children/search` với `parentId`

### Page 1: Quản lý Danh mục Cha (Parent Categories)

**API sử dụng:** `POST /api/v1/categories/search`

**Ví dụ code:**
```javascript
// Lấy danh sách parent categories với pagination
const fetchParentCategories = async (page = 1, limit = 10, filters = {}) => {
  const response = await fetch('/api/v1/categories/search', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: filters.name || null,        // Search theo tên
      isActive: filters.isActive || null, // Filter theo status
      sortBy: 'createdAt',
      sortOrder: 'DESC',
      page: page,
      limit: limit
    })
  });
  
  const result = await response.json();
  return result;
};

// Sử dụng
const parentCategories = await fetchParentCategories(1, 10, {
  isActive: true
});

// Kết quả: Chỉ có parent categories (Điện thoại, Laptop, ...)
// Mỗi item có: parentId = null, childrenIds = [3, 4, 5, ...]
```

**Phân biệt Parent:**
- `parentId === null` → Đây là parent category
- `childrenIds.length > 0` → Có children

### Page 2: Quản lý Danh mục Con (Child Categories)

**Bước 1: Lấy danh sách Parent để làm Dropdown Filter**

```javascript
// Lấy danh sách parent categories cho dropdown filter
const fetchParentsForFilter = async () => {
  const response = await fetch('/api/v1/categories/parents?includeInactive=false', {
    method: 'GET'
  });
  
  const result = await response.json();
  return result.data; // Array of parent categories
};

// Sử dụng để populate dropdown
const parentOptions = await fetchParentsForFilter();
// [{ id: 1, name: "Điện thoại" }, { id: 2, name: "Laptop" }, ...]
```

**Bước 2: Lấy danh sách Child Categories với Filter**

```javascript
// Lấy danh sách child categories với filter theo parent
const fetchChildCategories = async (page = 1, limit = 10, filters = {}) => {
  const response = await fetch('/api/v1/categories/children/search', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      parentId: filters.parentId || null,  // Filter theo parent (dropdown)
      name: filters.name || null,           // Search theo tên
      isActive: filters.isActive || null,   // Filter theo status
      sortBy: 'createdAt',
      sortOrder: 'DESC',
      page: page,
      limit: limit
    })
  });
  
  const result = await response.json();
  return result;
};

// Sử dụng
const childCategories = await fetchChildCategories(1, 10, {
  parentId: 1,      // Chỉ lấy children của "Điện thoại"
  isActive: true
});

// Kết quả: Chỉ có child categories (iPhone, Samsung, ...)
// Mỗi item có: parentId = 1, childrenIds = []
```

**Phân biệt Child:**
- `parentId !== null` → Đây là child category
- `childrenIds.length === 0` → Không có children (hoặc có thể có nếu là nested)

**Ví dụ UI Flow:**
```javascript
// Component: ChildCategoriesPage
const [selectedParentId, setSelectedParentId] = useState(null);
const [parentOptions, setParentOptions] = useState([]);
const [childCategories, setChildCategories] = useState([]);

// Load parent options cho dropdown
useEffect(() => {
  fetchParentsForFilter().then(setParentOptions);
}, []);

// Load child categories khi filter thay đổi
useEffect(() => {
  fetchChildCategories(1, 10, {
    parentId: selectedParentId,
    isActive: true
  }).then(result => {
    setChildCategories(result.data);
  });
}, [selectedParentId]);

// Render
return (
  <div>
    {/* Dropdown filter */}
    <select 
      value={selectedParentId || ''} 
      onChange={(e) => setSelectedParentId(e.target.value ? Number(e.target.value) : null)}
    >
      <option value="">Tất cả danh mục cha</option>
      {parentOptions.map(parent => (
        <option key={parent.id} value={parent.id}>{parent.name}</option>
      ))}
    </select>
    
    {/* Table hiển thị child categories */}
    <table>
      {childCategories.map(child => (
        <tr key={child.id}>
          <td>{child.name}</td>
          <td>{child.parentId ? `Parent ID: ${child.parentId}` : 'Root'}</td>
        </tr>
      ))}
    </table>
  </div>
);
```

## Ví dụ sử dụng thực tế

### 1. Tạo category và thêm children

```javascript
// Bước 1: Tạo parent category
const createParent = await fetch('/api/v1/categories', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer <token>',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: "Điện thoại",
    nameEn: "Smartphones",
    isActive: true
  })
});

const { data: parent } = await createParent.json();
console.log('Parent ID:', parent.id); // Ví dụ: 1

// Bước 2: Tạo child categories
const createChild1 = await fetch('/api/v1/categories', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer <token>',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: "iPhone",
    isActive: true
  })
});

const { data: child1 } = await createChild1.json();

// Bước 3: Thêm child vào parent
await fetch(`/api/v1/categories/${parent.id}/children`, {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer <token>',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    childId: child1.id
  })
});
```

### 2. Lấy tất cả children của một parent

```javascript
// Cách 1: Sử dụng GET endpoint
const response = await fetch('/api/v1/categories/1/children?includeInactive=false');
const { data: children } = await response.json();

// Cách 2: Sử dụng POST search với filter parentId
const response2 = await fetch('/api/v1/categories/children/search', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    parentId: 1,
    isActive: true
  })
});
const { data: children2 } = await response2.json();
```

### 3. Tìm kiếm children theo tên và parent

```javascript
const response = await fetch('/api/v1/categories/children/search', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    parentId: 1,        // Chỉ lấy children của "Điện thoại"
    name: "iPhone",     // Tìm kiếm theo tên
    isActive: true,
    page: 1,
    limit: 10
  })
});

const { data: children } = await response.json();
```

### 4. Lấy tất cả children (không filter parent)

**Cách 1: Sử dụng GET endpoint (đơn giản hơn, phù hợp cho dropdown)**
```javascript
// Lấy tất cả child categories active - phù hợp cho dropdown khi add product
const response = await fetch('/api/v1/categories/children?includeInactive=false&language=vi');
const { data: allChildren } = await response.json();

// Sử dụng trong dropdown
<select>
  {allChildren.map(child => (
    <option key={child.id} value={child.id}>{child.name}</option>
  ))}
</select>
```

**Cách 2: Sử dụng POST search endpoint (có pagination và filter)**
```javascript
const response = await fetch('/api/v1/categories/children/search', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    isActive: true
    // Không truyền parentId → lấy tất cả children
  })
});

const { data: allChildren } = await response.json();
```

### 5. Di chuyển child sang parent khác

```javascript
const oldParentId = 1;
const newParentId = 2;
const childId = 3;

// Bước 1: Xóa child khỏi parent cũ
await fetch(`/api/v1/categories/${oldParentId}/children`, {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer <token>',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    childId: childId
  })
});

// Bước 2: Thêm child vào parent mới
await fetch(`/api/v1/categories/${newParentId}/children`, {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer <token>',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    childId: childId
  })
});
```

### 6. Xóa child khỏi parent (chuyển thành root)

```javascript
const parentId = 1;
const childId = 3;

await fetch(`/api/v1/categories/${parentId}/children`, {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer <token>',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    childId: childId
  })
});

// Sau khi xóa, category có ID = 3 sẽ trở thành root category
```

### 7. Build category tree từ API

```javascript
// Bước 1: Lấy tất cả categories
const allCategoriesResponse = await fetch('/api/v1/categories/search', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    isActive: true,
    limit: 1000
  })
});

const { data: allCategories } = await allCategoriesResponse.json();

// Bước 2: Build tree ở FE
const buildTree = (categories) => {
  const map = new Map();
  const roots = [];
  
  // Tạo map
  categories.forEach(cat => {
    map.set(cat.id, { ...cat, children: [] });
  });
  
  // Build tree dựa trên parentId và childrenIds
  categories.forEach(cat => {
    if (cat.parentId === null) {
      roots.push(map.get(cat.id));
    } else {
      const parent = map.get(cat.parentId);
      if (parent) {
        parent.children.push(map.get(cat.id));
      }
    }
  });
  
  return roots;
};

const categoryTree = buildTree(allCategories);
```

## Error Codes

| Status Code | Mô tả |
|------------|-------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error, business logic error) |
| 401 | Unauthorized (thiếu hoặc token không hợp lệ) |
| 403 | Forbidden (không có quyền admin) |
| 404 | Not Found (category không tồn tại) |

## Lưu ý quan trọng

1. **Pagination mặc định**: `limit = 10` (đã thay đổi từ 50)
2. **Category mới tạo**: Mặc định là root category (không có parent)
3. **Quản lý children**: Phải sử dụng API riêng (`POST /categories/:id/children`, `DELETE /categories/:id/children`)
4. **Không thể set parent khi tạo**: Phải tạo category trước, sau đó thêm vào parent
5. **Một child chỉ có một parent**: Nếu muốn đổi parent, phải xóa khỏi parent cũ rồi thêm vào parent mới
6. **Soft delete**: Category bị xóa vẫn tồn tại trong DB nhưng `isActive = false`
7. **Hard delete**: Xóa vĩnh viễn, chỉ cho phép nếu không có products liên quan
8. **Multilingual**: Sử dụng query parameter `language` để lấy dữ liệu theo ngôn ngữ
9. **View inactive**: Mặc định cho phép view category inactive (để có thể kích hoạt lại)

## Migration Notes

- Cấu trúc mới sử dụng bảng `category_children` để lưu quan hệ parent-child
- Category không còn field `parent_id` trong bảng `categories`
- Quan hệ được quản lý thông qua bảng riêng, linh hoạt hơn và dễ mở rộng

# Flutter Search Suggestions API Guide

## 📋 Tổng quan

Tài liệu này mô tả API cho tính năng **Search Suggestions** (Gợi ý tìm kiếm) trong ứng dụng Flutter ecommerce. Tính năng này giúp cải thiện UX bằng cách cung cấp gợi ý tìm kiếm khi người dùng đang nhập từ khóa.

## 🎯 Use Cases

1. **Khi mở SearchPage (chưa nhập gì):**
   - Hiển thị lịch sử tìm kiếm (từ local storage)
   - Hoặc hiển thị popular/top search keywords từ BE

2. **Khi người dùng đang nhập (có query):**
   - Gọi API để lấy suggestions real-time dựa trên query
   - Hiển thị danh sách suggestions phù hợp

3. **Khi người dùng chọn suggestion:**
   - Thực hiện search với query đó
   - Lưu vào lịch sử search (local storage)

## 🔌 API Endpoints

### 1. Get Search Suggestions (Autocomplete)

**Endpoint:** `GET /products/search-suggestions`

**Mô tả:** Lấy danh sách gợi ý tìm kiếm dựa trên query string.

**Query Parameters:**
```typescript
{
  query: string;        // Từ khóa tìm kiếm (required)
  language: string;     // 'vi' hoặc 'en' (required)
  limit?: number;       // Số lượng suggestions (optional, default: 10, max: 20)
}
```

**Request Example:**
```http
GET /products/search-suggestions?query=điện&language=vi&limit=10
```

**Response Format:**
```json
{
  "success": true,
  "data": [
    {
      "text": "Điện thoại",
      "type": "product",        // "product" hoặc "category"
      "count": 150              // Số lượng sản phẩm/danh mục (optional)
    },
    {
      "text": "Điện tử",
      "type": "category",
      "count": 45
    },
    {
      "text": "Điện thoại Samsung",
      "type": "product",
      "count": 23
    }
  ],
  "total": 3
}
```

**Response Model:**
```dart
class SearchSuggestion {
  final String text;           // Text hiển thị
  final String type;           // "product" hoặc "category"
  final int? count;            // Số lượng (optional)
}

class SearchSuggestionsResponse {
  final bool success;
  final List<SearchSuggestion> data;
  final int total;
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error message",
  "error": "Error details"
}
```

### 2. Get Popular Search Keywords (Optional)

**Endpoint:** `GET /products/popular-searches`

**Mô tả:** Lấy danh sách từ khóa tìm kiếm phổ biến (dùng khi chưa có query).

**Query Parameters:**
```typescript
{
  language: string;     // 'vi' hoặc 'en' (required)
  limit?: number;       // Số lượng (optional, default: 10, max: 20)
}
```

**Request Example:**
```http
GET /products/popular-searches?language=vi&limit=10
```

**Response Format:**
```json
{
  "success": true,
  "data": [
    {
      "text": "Điện thoại",
      "count": 1250
    },
    {
      "text": "Laptop",
      "count": 890
    },
    {
      "text": "Tai nghe",
      "count": 650
    }
  ],
  "total": 3
}
```

## 🔄 Flow Implementation

### Flow 1: Mở SearchPage (chưa nhập)

```
1. User mở SearchPage
2. App load lịch sử search từ SharedPreferences (local)
3. Nếu có lịch sử → Hiển thị lịch sử
4. Nếu không có lịch sử → Gọi API /products/popular-searches
5. Hiển thị popular searches
```

### Flow 2: User đang nhập

```
1. User nhập text vào search field
2. App debounce 300ms
3. Nếu text rỗng → Hiển thị lịch sử/popular
4. Nếu có text → Gọi API /products/search-suggestions?query=...
5. Hiển thị suggestions từ API
```

### Flow 3: User chọn suggestion hoặc search

```
1. User chọn suggestion hoặc nhấn search
2. Thực hiện search với query đó
3. Lưu query vào lịch sử search (SharedPreferences)
4. Navigate về HomePage và hiển thị kết quả
```

## 📝 Implementation Details

### Suggestions Source Priority:

1. **Khi có query:**
   - Gọi API `/products/search-suggestions` với query
   - Suggestions có thể là:
     - Tên sản phẩm phù hợp
     - Tên danh mục phù hợp
     - Từ khóa phổ biến liên quan

2. **Khi không có query:**
   - Ưu tiên: Lịch sử search từ local (SharedPreferences)
   - Fallback: Popular searches từ API `/products/popular-searches`

### Suggestions Logic (BE):

- **Product suggestions:** Tìm trong tên sản phẩm (name, nameEn) có chứa query
- **Category suggestions:** Tìm trong tên danh mục (name, nameEn) có chứa query
- **Sorting:** 
  - Ưu tiên theo độ phù hợp (relevance)
  - Hoặc theo số lượng sản phẩm (count) nếu có
- **Limit:** Mặc định 10, tối đa 20 suggestions

### Local Storage (Flutter):

- Lưu lịch sử search vào `SharedPreferences`
- Key: `search_history`
- Format: `List<String>` - danh sách các query đã search
- Giới hạn: Tối đa 20 queries gần nhất
- Khi search thành công → Thêm vào đầu danh sách, xóa duplicate

## 🎨 UI/UX Requirements

1. **Suggestions Display:**
   - Hiển thị dạng ListTile với icon search
   - Mỗi suggestion có text và có thể có count (số lượng)
   - Tap vào suggestion → Auto-fill và search

2. **Loading State:**
   - Hiển thị CircularProgressIndicator khi đang gọi API
   - Debounce 300ms để tránh gọi API quá nhiều

3. **Empty State:**
   - Khi không có suggestions → Hiển thị message "Nhập từ khóa để tìm kiếm"
   - Khi không có lịch sử → Hiển thị popular searches

## 🔒 Security & Performance

1. **Rate Limiting:**
   - BE nên có rate limiting cho endpoint suggestions
   - Giới hạn số request mỗi phút từ một IP/user

2. **Caching:**
   - BE có thể cache popular searches (TTL: 1 giờ)
   - Flutter cache lịch sử search local

3. **Debounce:**
   - Flutter debounce 300ms trước khi gọi API
   - Tránh spam request khi user đang gõ nhanh

## 📊 Example Scenarios

### Scenario 1: User search "điện thoại"

**Request:**
```http
GET /products/search-suggestions?query=điện&language=vi&limit=10
```

**Response:**
```json
{
  "success": true,
  "data": [
    { "text": "Điện thoại", "type": "product", "count": 150 },
    { "text": "Điện tử", "type": "category", "count": 45 },
    { "text": "Điện thoại Samsung", "type": "product", "count": 23 },
    { "text": "Điện thoại iPhone", "type": "product", "count": 18 }
  ],
  "total": 4
}
```

### Scenario 2: User mở SearchPage lần đầu (chưa có lịch sử)

**Request:**
```http
GET /products/popular-searches?language=vi&limit=10
```

**Response:**
```json
{
  "success": true,
  "data": [
    { "text": "Điện thoại", "count": 1250 },
    { "text": "Laptop", "count": 890 },
    { "text": "Tai nghe", "count": 650 },
    { "text": "Chuột máy tính", "count": 420 },
    { "text": "Bàn phím", "count": 380 }
  ],
  "total": 5
}
```

## 🚀 Implementation Priority

### Phase 1 (MVP - Minimum Viable Product):
- ✅ Implement local search history (SharedPreferences)
- ✅ Hiển thị lịch sử khi mở SearchPage
- ✅ Filter suggestions từ lịch sử khi user nhập

### Phase 2 (Enhanced):
- ⏳ BE implement `/products/search-suggestions` endpoint
- ⏳ Flutter integrate API suggestions
- ⏳ Combine local history + API suggestions

### Phase 3 (Advanced):
- ⏳ BE implement `/products/popular-searches` endpoint
- ⏳ Flutter integrate popular searches
- ⏳ Smart suggestions ranking

## 📌 Notes

1. **Backward Compatibility:**
   - Nếu BE chưa có API suggestions, Flutter vẫn hoạt động với local history
   - Khi BE có API, Flutter sẽ tự động sử dụng

2. **Multi-language Support:**
   - Tất cả API đều có parameter `language` ('vi' hoặc 'en')
   - Suggestions phải match với ngôn ngữ hiện tại

3. **Error Handling:**
   - Nếu API suggestions fail → Fallback về local history
   - Nếu không có local history → Hiển thị empty state

## ❓ Questions for BE Team

1. BE có sẵn endpoint suggestions chưa? Nếu chưa, có thể implement không?
2. Suggestions nên lấy từ đâu? (Product names, Category names, hoặc cả hai?)
3. Có cần endpoint popular searches không? Hay chỉ cần suggestions khi có query?
4. Có cần rate limiting không? Nếu có, limit là bao nhiêu?
5. Suggestions có cần sort theo độ phù hợp không? Hay chỉ sort theo count?

---

**Tài liệu này được tạo để đồng bộ giữa Frontend (Flutter) và Backend team về tính năng Search Suggestions.**


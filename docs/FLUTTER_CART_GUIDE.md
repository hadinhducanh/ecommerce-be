# 🛒 Cart API - Hướng dẫn tích hợp Flutter

Tài liệu này hướng dẫn chi tiết cách sử dụng Cart API từ Flutter client để quản lý giỏ hàng.

## 📋 Mục lục

- [Tổng quan](#tổng-quan)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
- [Data Models](#data-models)
- [Flow Logic](#flow-logic)
- [Flutter Implementation](#flutter-implementation)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

---

## 🎯 Tổng quan

### Base URL
```
http://localhost:8080/api/v1
```

### Chức năng Cart API

Cart API cung cấp đầy đủ chức năng quản lý giỏ hàng:

- ✅ Thêm sản phẩm vào giỏ hàng
- ✅ Xem danh sách sản phẩm trong giỏ
- ✅ Cập nhật số lượng sản phẩm
- ✅ Xóa sản phẩm khỏi giỏ
- ✅ Xóa toàn bộ giỏ hàng
- ✅ Đếm số lượng items trong giỏ
- ✅ Tính tổng giá trị giỏ hàng

---

## 🔐 Authentication

**TẤT CẢ** các endpoint Cart đều yêu cầu JWT token.

### Header Required

```http
Authorization: Bearer <your_jwt_token>
```

### Lấy JWT Token

1. **Login** hoặc **Register** trước
2. Lưu token từ response
3. Gửi kèm mọi request đến Cart API

**Ví dụ Login:**

```dart
// Flutter code
final response = await http.post(
  Uri.parse('$baseUrl/auth/login'),
  headers: {'Content-Type': 'application/json'},
  body: jsonEncode({
    'email': 'user@example.com',
    'password': 'password123',
  }),
);

final data = jsonDecode(response.body);
final token = data['token']; // Lưu token này
```

---

## 📚 Endpoints

### 1. Thêm sản phẩm vào giỏ hàng

**Add to Cart**

```http
POST /api/v1/cart
```

#### Request Headers

```http
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body

```json
{
  "productId": 1,
  "quantity": 2
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `productId` | `uint` | ✅ | ID của sản phẩm |
| `quantity` | `int` | ✅ | Số lượng (≥ 1) |

#### Response Success (200 OK)

```json
{
  "success": true,
  "message": "Thêm sản phẩm vào giỏ hàng thành công",
  "data": {
    "id": 1,
    "quantity": 2,
    "userId": 1,
    "productId": 1,
    "product": {
      "id": 1,
      "name": "iPhone 15 Pro Max 256GB",
      "nameEn": "iPhone 15 Pro Max 256GB",
      "price": 32990000,
      "stock": 48,
      "image": "https://...",
      "categoryId": 3,
      "isActive": true
    },
    "createdAt": "2025-12-30T10:00:00Z",
    "updatedAt": "2025-12-30T10:00:00Z"
  }
}
```

#### Response Error

**Sản phẩm không tồn tại (400):**
```json
{
  "error": "sản phẩm không tồn tại hoặc không khả dụng"
}
```

**Không đủ stock (400):**
```json
{
  "error": "sản phẩm chỉ còn 5 sản phẩm trong kho"
}
```

**Chưa login (401):**
```json
{
  "error": "Unauthorized"
}
```

#### Logic đặc biệt

⚡ **Auto-merge:** Nếu sản phẩm đã có trong giỏ, số lượng sẽ được **cộng dồn** thay vì tạo item mới.

**Ví dụ:**
- Giỏ hiện tại: Product ID 1, quantity = 2
- Request: Add Product ID 1, quantity = 3
- Kết quả: Product ID 1, quantity = 5 ✅

---

### 2. Xem giỏ hàng

**Get Cart**

```http
GET /api/v1/cart
```

#### Request Headers

```http
Authorization: Bearer <token>
```

#### Response Success (200 OK)

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "quantity": 2,
        "userId": 1,
        "productId": 1,
        "product": {
          "id": 1,
          "name": "iPhone 15 Pro Max 256GB",
          "price": 32990000,
          "stock": 48,
          "image": "https://...",
          "isActive": true
        },
        "createdAt": "2025-12-30T10:00:00Z",
        "updatedAt": "2025-12-30T10:00:00Z"
      },
      {
        "id": 2,
        "quantity": 1,
        "userId": 1,
        "productId": 5,
        "product": {
          "id": 5,
          "name": "Samsung Galaxy S24 Ultra 256GB",
          "price": 29990000,
          "stock": 44,
          "image": "https://...",
          "isActive": true
        },
        "createdAt": "2025-12-30T10:05:00Z",
        "updatedAt": "2025-12-30T10:05:00Z"
      }
    ],
    "totalItems": 3,
    "totalPrice": 95970000
  }
}
```

#### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `items` | `array` | Danh sách cart items |
| `totalItems` | `int` | Tổng số lượng sản phẩm |
| `totalPrice` | `float64` | Tổng giá trị (VNĐ) |

#### Giỏ hàng rỗng

```json
{
  "success": true,
  "data": {
    "items": [],
    "totalItems": 0,
    "totalPrice": 0
  }
}
```

---

### 3. Đếm số lượng items

**Get Cart Count**

```http
GET /api/v1/cart/count
```

#### Request Headers

```http
Authorization: Bearer <token>
```

#### Response Success (200 OK)

```json
{
  "success": true,
  "count": 3
}
```

**Use case:** Hiển thị badge số lượng trên icon giỏ hàng.

---

### 4. Cập nhật số lượng

**Update Cart Item**

```http
PUT /api/v1/cart/:id
```

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `uint` | Cart Item ID (không phải Product ID) |

#### Request Headers

```http
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body

```json
{
  "quantity": 5
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `quantity` | `int` | ✅ | Số lượng mới (≥ 1) |

#### Response Success (200 OK)

```json
{
  "success": true,
  "message": "Cập nhật giỏ hàng thành công",
  "data": {
    "id": 1,
    "quantity": 5,
    "userId": 1,
    "productId": 1,
    "product": {
      "id": 1,
      "name": "iPhone 15 Pro Max 256GB",
      "price": 32990000,
      "stock": 45,
      "image": "https://..."
    },
    "createdAt": "2025-12-30T10:00:00Z",
    "updatedAt": "2025-12-30T10:30:00Z"
  }
}
```

#### Response Error

**Không đủ stock (400):**
```json
{
  "error": "sản phẩm chỉ còn 3 sản phẩm trong kho"
}
```

**Không tìm thấy (400):**
```json
{
  "error": "không tìm thấy sản phẩm trong giỏ hàng"
}
```

---

### 5. Xóa sản phẩm khỏi giỏ

**Delete Cart Item**

```http
DELETE /api/v1/cart/:id
```

#### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `uint` | Cart Item ID |

#### Request Headers

```http
Authorization: Bearer <token>
```

#### Response Success (200 OK)

```json
{
  "success": true,
  "message": "Xóa sản phẩm khỏi giỏ hàng thành công"
}
```

#### Response Error

**Không tìm thấy (400):**
```json
{
  "error": "không tìm thấy sản phẩm trong giỏ hàng"
}
```

---

### 6. Xóa toàn bộ giỏ hàng

**Clear Cart**

```http
DELETE /api/v1/cart
```

#### Request Headers

```http
Authorization: Bearer <token>
```

#### Response Success (200 OK)

```json
{
  "success": true,
  "message": "Xóa toàn bộ giỏ hàng thành công"
}
```

**Use case:** Sau khi checkout thành công.

---

## 📦 Data Models

### CartItem Model

```dart
class CartItem {
  final int id;
  final int quantity;
  final int userId;
  final int productId;
  final Product product;
  final DateTime createdAt;
  final DateTime updatedAt;

  CartItem({
    required this.id,
    required this.quantity,
    required this.userId,
    required this.productId,
    required this.product,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CartItem.fromJson(Map<String, dynamic> json) {
    return CartItem(
      id: json['id'],
      quantity: json['quantity'],
      userId: json['userId'],
      productId: json['productId'],
      product: Product.fromJson(json['product']),
      createdAt: DateTime.parse(json['createdAt']),
      updatedAt: DateTime.parse(json['updatedAt']),
    );
  }

  // Tính tổng giá của item này
  double get totalPrice => product.price * quantity;
}
```

### Cart Summary Model

```dart
class CartSummary {
  final List<CartItem> items;
  final int totalItems;
  final double totalPrice;

  CartSummary({
    required this.items,
    required this.totalItems,
    required this.totalPrice,
  });

  factory CartSummary.fromJson(Map<String, dynamic> json) {
    return CartSummary(
      items: (json['items'] as List)
          .map((item) => CartItem.fromJson(item))
          .toList(),
      totalItems: json['totalItems'],
      totalPrice: (json['totalPrice'] as num).toDouble(),
    );
  }

  bool get isEmpty => items.isEmpty;
  int get itemCount => items.length;
}
```

### Product Model (simplified)

```dart
class Product {
  final int id;
  final String name;
  final String? nameEn;
  final double price;
  final int stock;
  final String? image;
  final int categoryId;
  final bool isActive;

  Product({
    required this.id,
    required this.name,
    this.nameEn,
    required this.price,
    required this.stock,
    this.image,
    required this.categoryId,
    required this.isActive,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'],
      name: json['name'],
      nameEn: json['nameEn'],
      price: (json['price'] as num).toDouble(),
      stock: json['stock'],
      image: json['image'],
      categoryId: json['categoryId'],
      isActive: json['isActive'],
    );
  }
}
```

---

## 🔄 Flow Logic

### Flow 1: Thêm sản phẩm vào giỏ

```
User click "Add to Cart"
    ↓
Check if user logged in?
    ↓ No → Navigate to Login
    ↓ Yes
Check stock available?
    ↓ No → Show error "Out of stock"
    ↓ Yes
Send POST /cart request
    ↓
Success?
    ↓ No → Show error message
    ↓ Yes
Update local cart state
    ↓
Show success message
    ↓
Update cart badge count
```

### Flow 2: Xem giỏ hàng

```
User navigate to Cart Screen
    ↓
Show loading indicator
    ↓
Send GET /cart request
    ↓
Success?
    ↓ No → Show error screen
    ↓ Yes
Parse cart data
    ↓
Display cart items list
    ↓
Calculate & show total price
    ↓
Show "Checkout" button if not empty
```

### Flow 3: Cập nhật số lượng

```
User change quantity (+ or -)
    ↓
Validate quantity > 0
    ↓
Optimistic update UI (instant feedback)
    ↓
Send PUT /cart/:id request
    ↓
Success?
    ↓ No → Revert UI, show error
    ↓ Yes
Confirm UI update
    ↓
Recalculate total price
```

### Flow 4: Xóa sản phẩm

```
User click "Remove" button
    ↓
Show confirmation dialog (optional)
    ↓
User confirms?
    ↓ No → Cancel
    ↓ Yes
Optimistic remove from UI
    ↓
Send DELETE /cart/:id request
    ↓
Success?
    ↓ No → Revert UI, show error
    ↓ Yes
Confirm removal
    ↓
Recalculate total
    ↓
Update cart badge
```

---

## 💻 Flutter Implementation

### 1. Cart Service

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

class CartService {
  final String baseUrl = 'http://localhost:8080/api/v1';
  final String token; // JWT token

  CartService({required this.token});

  // Headers with token
  Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer $token',
  };

  // 1. Add to Cart
  Future<CartItem> addToCart(int productId, int quantity) async {
    final response = await http.post(
      Uri.parse('$baseUrl/cart'),
      headers: _headers,
      body: jsonEncode({
        'productId': productId,
        'quantity': quantity,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return CartItem.fromJson(data['data']);
    } else {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Không thể thêm vào giỏ hàng');
    }
  }

  // 2. Get Cart
  Future<CartSummary> getCart() async {
    final response = await http.get(
      Uri.parse('$baseUrl/cart'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return CartSummary.fromJson(data['data']);
    } else {
      throw Exception('Không thể tải giỏ hàng');
    }
  }

  // 3. Get Cart Count
  Future<int> getCartCount() async {
    final response = await http.get(
      Uri.parse('$baseUrl/cart/count'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data['count'];
    } else {
      return 0;
    }
  }

  // 4. Update Cart Item
  Future<CartItem> updateCartItem(int cartItemId, int quantity) async {
    final response = await http.put(
      Uri.parse('$baseUrl/cart/$cartItemId'),
      headers: _headers,
      body: jsonEncode({
        'quantity': quantity,
      }),
    );

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return CartItem.fromJson(data['data']);
    } else {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Không thể cập nhật');
    }
  }

  // 5. Delete Cart Item
  Future<void> deleteCartItem(int cartItemId) async {
    final response = await http.delete(
      Uri.parse('$baseUrl/cart/$cartItemId'),
      headers: _headers,
    );

    if (response.statusCode != 200) {
      final error = jsonDecode(response.body);
      throw Exception(error['error'] ?? 'Không thể xóa');
    }
  }

  // 6. Clear Cart
  Future<void> clearCart() async {
    final response = await http.delete(
      Uri.parse('$baseUrl/cart'),
      headers: _headers,
    );

    if (response.statusCode != 200) {
      throw Exception('Không thể xóa giỏ hàng');
    }
  }
}
```

### 2. Cart Provider (State Management)

**Sử dụng Provider pattern:**

```dart
import 'package:flutter/material.dart';

class CartProvider with ChangeNotifier {
  final CartService _cartService;
  CartSummary? _cart;
  bool _isLoading = false;
  String? _error;

  CartProvider(this._cartService);

  CartSummary? get cart => _cart;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get itemCount => _cart?.totalItems ?? 0;
  double get totalPrice => _cart?.totalPrice ?? 0;

  // Load cart
  Future<void> loadCart() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _cart = await _cartService.getCart();
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Add to cart
  Future<void> addToCart(int productId, int quantity) async {
    try {
      await _cartService.addToCart(productId, quantity);
      await loadCart(); // Reload cart
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      rethrow;
    }
  }

  // Update quantity
  Future<void> updateQuantity(int cartItemId, int quantity) async {
    // Optimistic update
    final oldCart = _cart;
    
    try {
      // Update UI immediately
      _cart = _cart?.copyWith(
        items: _cart!.items.map((item) {
          if (item.id == cartItemId) {
            return item.copyWith(quantity: quantity);
          }
          return item;
        }).toList(),
      );
      notifyListeners();

      // Send request
      await _cartService.updateCartItem(cartItemId, quantity);
      await loadCart(); // Reload to get accurate data
    } catch (e) {
      // Revert on error
      _cart = oldCart;
      _error = e.toString();
      notifyListeners();
      rethrow;
    }
  }

  // Remove item
  Future<void> removeItem(int cartItemId) async {
    final oldCart = _cart;

    try {
      // Optimistic update
      _cart = _cart?.copyWith(
        items: _cart!.items.where((item) => item.id != cartItemId).toList(),
      );
      notifyListeners();

      await _cartService.deleteCartItem(cartItemId);
      await loadCart();
    } catch (e) {
      _cart = oldCart;
      _error = e.toString();
      notifyListeners();
      rethrow;
    }
  }

  // Clear cart
  Future<void> clearCart() async {
    try {
      await _cartService.clearCart();
      _cart = CartSummary(items: [], totalItems: 0, totalPrice: 0);
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      rethrow;
    }
  }
}
```

### 3. Cart Screen UI

```dart
class CartScreen extends StatefulWidget {
  @override
  _CartScreenState createState() => _CartScreenState();
}

class _CartScreenState extends State<CartScreen> {
  @override
  void initState() {
    super.initState();
    // Load cart khi vào screen
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<CartProvider>().loadCart();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Giỏ hàng'),
        actions: [
          // Clear cart button
          Consumer<CartProvider>(
            builder: (context, cart, _) {
              if (cart.cart?.items.isEmpty ?? true) return SizedBox();
              
              return IconButton(
                icon: Icon(Icons.delete_outline),
                onPressed: () => _showClearCartDialog(context),
              );
            },
          ),
        ],
      ),
      body: Consumer<CartProvider>(
        builder: (context, cartProvider, _) {
          // Loading state
          if (cartProvider.isLoading) {
            return Center(child: CircularProgressIndicator());
          }

          // Error state
          if (cartProvider.error != null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('Có lỗi xảy ra'),
                  SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => cartProvider.loadCart(),
                    child: Text('Thử lại'),
                  ),
                ],
              ),
            );
          }

          final cart = cartProvider.cart;

          // Empty cart
          if (cart == null || cart.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.shopping_cart_outlined, size: 80, color: Colors.grey),
                  SizedBox(height: 16),
                  Text('Giỏ hàng trống'),
                  SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => Navigator.pop(context),
                    child: Text('Tiếp tục mua sắm'),
                  ),
                ],
              ),
            );
          }

          // Cart items
          return Column(
            children: [
              Expanded(
                child: ListView.builder(
                  itemCount: cart.items.length,
                  itemBuilder: (context, index) {
                    final item = cart.items[index];
                    return CartItemCard(item: item);
                  },
                ),
              ),
              // Bottom summary
              _buildCartSummary(cart),
            ],
          );
        },
      ),
    );
  }

  Widget _buildCartSummary(CartSummary cart) {
    return Container(
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black12,
            blurRadius: 4,
            offset: Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Tổng cộng:', style: TextStyle(fontSize: 16)),
                Text(
                  '${cart.totalPrice.toStringAsFixed(0)} VNĐ',
                  style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.red,
                  ),
                ),
              ],
            ),
            SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () {
                  // Navigate to checkout
                  Navigator.pushNamed(context, '/checkout');
                },
                style: ElevatedButton.styleFrom(
                  padding: EdgeInsets.symmetric(vertical: 16),
                ),
                child: Text(
                  'Thanh toán (${cart.totalItems} sản phẩm)',
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showClearCartDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Xóa giỏ hàng'),
        content: Text('Bạn có chắc muốn xóa toàn bộ giỏ hàng?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Hủy'),
          ),
          TextButton(
            onPressed: () {
              context.read<CartProvider>().clearCart();
              Navigator.pop(context);
            },
            child: Text('Xóa', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }
}
```

### 4. Cart Item Card Widget

```dart
class CartItemCard extends StatelessWidget {
  final CartItem item;

  const CartItemCard({required this.item});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Padding(
        padding: EdgeInsets.all(12),
        child: Row(
          children: [
            // Product image
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: Image.network(
                item.product.image ?? 'https://via.placeholder.com/80',
                width: 80,
                height: 80,
                fit: BoxFit.cover,
              ),
            ),
            SizedBox(width: 12),
            
            // Product info
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    item.product.name,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  SizedBox(height: 4),
                  Text(
                    '${item.product.price.toStringAsFixed(0)} VNĐ',
                    style: TextStyle(
                      fontSize: 14,
                      color: Colors.red,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  SizedBox(height: 8),
                  
                  // Quantity controls
                  Row(
                    children: [
                      // Decrease button
                      _buildQuantityButton(
                        context,
                        icon: Icons.remove,
                        onPressed: item.quantity > 1
                            ? () => _updateQuantity(context, item.quantity - 1)
                            : null,
                      ),
                      
                      // Quantity display
                      Container(
                        padding: EdgeInsets.symmetric(horizontal: 16),
                        child: Text(
                          '${item.quantity}',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      
                      // Increase button
                      _buildQuantityButton(
                        context,
                        icon: Icons.add,
                        onPressed: item.quantity < item.product.stock
                            ? () => _updateQuantity(context, item.quantity + 1)
                            : null,
                      ),
                      
                      Spacer(),
                      
                      // Delete button
                      IconButton(
                        icon: Icon(Icons.delete_outline, color: Colors.red),
                        onPressed: () => _deleteItem(context),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildQuantityButton(
    BuildContext context, {
    required IconData icon,
    VoidCallback? onPressed,
  }) {
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey[300]!),
        borderRadius: BorderRadius.circular(4),
      ),
      child: InkWell(
        onTap: onPressed,
        child: Padding(
          padding: EdgeInsets.all(4),
          child: Icon(
            icon,
            size: 20,
            color: onPressed != null ? Colors.black : Colors.grey,
          ),
        ),
      ),
    );
  }

  void _updateQuantity(BuildContext context, int newQuantity) {
    context.read<CartProvider>().updateQuantity(item.id, newQuantity).catchError((e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString())),
      );
    });
  }

  void _deleteItem(BuildContext context) {
    context.read<CartProvider>().removeItem(item.id).catchError((e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(e.toString())),
      );
    });
  }
}
```

### 5. Add to Cart Button (Product Screen)

```dart
class AddToCartButton extends StatelessWidget {
  final Product product;

  const AddToCartButton({required this.product});

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: product.stock > 0
          ? () => _addToCart(context)
          : null,
      icon: Icon(Icons.shopping_cart),
      label: Text(
        product.stock > 0 ? 'Thêm vào giỏ' : 'Hết hàng',
      ),
      style: ElevatedButton.styleFrom(
        padding: EdgeInsets.symmetric(vertical: 12, horizontal: 24),
      ),
    );
  }

  void _addToCart(BuildContext context) async {
    try {
      await context.read<CartProvider>().addToCart(product.id, 1);
      
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Đã thêm vào giỏ hàng'),
          action: SnackBarAction(
            label: 'Xem giỏ hàng',
            onPressed: () {
              Navigator.pushNamed(context, '/cart');
            },
          ),
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(e.toString()),
          backgroundColor: Colors.red,
        ),
      );
    }
  }
}
```

### 6. Cart Badge (App Bar)

```dart
class CartBadge extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Stack(
        children: [
          Icon(Icons.shopping_cart),
          
          // Badge with count
          Consumer<CartProvider>(
            builder: (context, cart, _) {
              final count = cart.itemCount;
              
              if (count == 0) return SizedBox();
              
              return Positioned(
                right: 0,
                top: 0,
                child: Container(
                  padding: EdgeInsets.all(2),
                  decoration: BoxDecoration(
                    color: Colors.red,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  constraints: BoxConstraints(
                    minWidth: 16,
                    minHeight: 16,
                  ),
                  child: Text(
                    count > 99 ? '99+' : '$count',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ),
              );
            },
          ),
        ],
      ),
      onPressed: () {
        Navigator.pushNamed(context, '/cart');
      },
    );
  }
}
```

---

## ⚠️ Error Handling

### Common Errors

| Status | Error Message | Nguyên nhân | Xử lý |
|--------|--------------|-------------|-------|
| 401 | Unauthorized | Token không hợp lệ/hết hạn | Logout → Login lại |
| 400 | sản phẩm không tồn tại | Product ID sai | Hiển thị lỗi |
| 400 | sản phẩm chỉ còn X trong kho | Không đủ stock | Giới hạn quantity |
| 400 | không tìm thấy sản phẩm trong giỏ | Cart item đã bị xóa | Reload cart |
| 500 | Internal Server Error | Lỗi server | Thử lại sau |

### Error Handling Strategy

```dart
Future<void> handleCartError(BuildContext context, dynamic error) async {
  String message = 'Có lỗi xảy ra';

  if (error.toString().contains('Unauthorized')) {
    message = 'Phiên đăng nhập hết hạn';
    // Logout and redirect to login
    await context.read<AuthProvider>().logout();
    Navigator.pushReplacementNamed(context, '/login');
    return;
  } else if (error.toString().contains('không đủ')) {
    message = 'Sản phẩm không đủ số lượng trong kho';
  } else if (error.toString().contains('không tồn tại')) {
    message = 'Sản phẩm không còn tồn tại';
    // Reload cart to remove invalid items
    context.read<CartProvider>().loadCart();
  }

  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Text(message),
      backgroundColor: Colors.red,
      action: SnackBarAction(
        label: 'OK',
        textColor: Colors.white,
        onPressed: () {},
      ),
    ),
  );
}
```

---

## 🎯 Best Practices

### 1. State Management

✅ **Sử dụng Provider/Riverpod/Bloc**
- Centralized state
- Reactive UI updates
- Easy to test

❌ **Tránh setState() trực tiếp cho cart logic**

### 2. Optimistic Updates

✅ **Update UI ngay lập tức**
```dart
// Update UI first
_cart = updatedCart;
notifyListeners();

// Then send request
await _cartService.updateItem();

// Revert if failed
catch (e) {
  _cart = oldCart;
  notifyListeners();
}
```

### 3. Token Management

✅ **Lưu token secure**
```dart
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

final storage = FlutterSecureStorage();
await storage.write(key: 'jwt_token', value: token);
```

❌ **Không lưu trong SharedPreferences (insecure)**

### 4. Loading States

✅ **Hiển thị loading indicator**
```dart
if (isLoading) {
  return Center(child: CircularProgressIndicator());
}
```

### 5. Network Error Handling

✅ **Retry mechanism**
```dart
Future<T> retry<T>(Future<T> Function() fn, {int maxAttempts = 3}) async {
  int attempt = 0;
  while (true) {
    try {
      return await fn();
    } catch (e) {
      attempt++;
      if (attempt >= maxAttempts) rethrow;
      await Future.delayed(Duration(seconds: attempt));
    }
  }
}
```

### 6. Cache Strategy

✅ **Cache cart data locally**
```dart
// Save to local storage
await storage.write(key: 'cached_cart', value: jsonEncode(cart.toJson()));

// Load from cache first, then refresh
final cachedData = await storage.read(key: 'cached_cart');
if (cachedData != null) {
  _cart = CartSummary.fromJson(jsonDecode(cachedData));
  notifyListeners();
}
// Then fetch from server
await loadCart();
```

### 7. Validation

✅ **Validate before sending request**
```dart
if (quantity < 1) {
  throw Exception('Số lượng phải >= 1');
}
if (quantity > product.stock) {
  throw Exception('Không đủ hàng trong kho');
}
```

### 8. Debouncing

✅ **Debounce quantity updates**
```dart
import 'package:rxdart/rxdart.dart';

final _quantitySubject = BehaviorSubject<int>();

_quantitySubject
  .debounceTime(Duration(milliseconds: 500))
  .listen((quantity) {
    updateCartItem(cartItemId, quantity);
  });
```

### 9. Error Messages

✅ **User-friendly messages**
```dart
try {
  await addToCart();
} catch (e) {
  String userMessage = _getUserFriendlyMessage(e);
  showError(userMessage);
}

String _getUserFriendlyMessage(dynamic error) {
  if (error.toString().contains('stock')) {
    return 'Sản phẩm không đủ trong kho';
  }
  return 'Có lỗi xảy ra. Vui lòng thử lại';
}
```

### 10. Testing

✅ **Write unit tests**
```dart
test('Add to cart increases item count', () async {
  final cartService = MockCartService();
  final provider = CartProvider(cartService);
  
  await provider.addToCart(1, 2);
  
  expect(provider.itemCount, 2);
});
```

---

## 📱 Complete Flow Example

### Scenario: User mua 2 iPhone 15

```dart
// 1. User ở Product Detail Screen
class ProductDetailScreen extends StatelessWidget {
  final Product product; // iPhone 15
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          // Product info...
          
          ElevatedButton(
            onPressed: () async {
              // 2. Click "Thêm vào giỏ"
              try {
                await context.read<CartProvider>().addToCart(
                  product.id,  // productId = 1
                  2,          // quantity = 2
                );
                
                // 3. Success → Show snackbar
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Đã thêm vào giỏ')),
                );
                
                // 4. Cart badge tự động update từ 0 → 2
                
              } catch (e) {
                // 5. Error → Show error
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(e.toString()),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            child: Text('Thêm vào giỏ'),
          ),
        ],
      ),
    );
  }
}

// 6. User click vào Cart icon
// → Navigate to CartScreen
// → CartScreen tự động load cart từ server
// → Hiển thị 2 iPhone 15 với tổng giá

// 7. User thay đổi quantity từ 2 → 3
// → Optimistic update UI ngay
// → Gửi PUT request
// → Thành công → Confirm UI
// → Lỗi → Revert về 2

// 8. User click Checkout
// → Navigate to CheckoutScreen với cart data
```

---

## 🔍 Debugging Tips

### 1. Check API Response

```dart
print('Response status: ${response.statusCode}');
print('Response body: ${response.body}');
```

### 2. Verify Token

```dart
print('Token: ${token.substring(0, 20)}...'); // First 20 chars
```

### 3. Network Inspector

- Sử dụng Charles Proxy / Proxyman
- Xem request/response thực tế

### 4. Error Logs

```dart
try {
  await cartService.addToCart();
} catch (e, stackTrace) {
  print('Error: $e');
  print('Stack trace: $stackTrace');
}
```

---

## ✅ Checklist Integration

- [ ] Đã tạo CartService với tất cả methods
- [ ] Đã setup CartProvider/State management
- [ ] Đã implement Cart Screen UI
- [ ] Đã implement Cart Item Card widget
- [ ] Đã thêm Cart Badge vào AppBar
- [ ] Đã implement Add to Cart button
- [ ] Đã handle errors properly
- [ ] Đã test add to cart flow
- [ ] Đã test update quantity flow
- [ ] Đã test delete item flow
- [ ] Đã test clear cart flow
- [ ] Đã test với network offline
- [ ] Đã test với invalid token
- [ ] Đã optimize performance (debouncing, caching)

---

## 📞 Support

Nếu gặp vấn đề khi tích hợp:

1. Check API response trong network inspector
2. Verify JWT token còn hạn
3. Check product stock > 0
4. Review error messages từ server
5. Test với Postman trước khi implement Flutter

---

**🎉 Hoàn thành!** Giờ bạn có thể tích hợp đầy đủ Cart functionality vào Flutter app!

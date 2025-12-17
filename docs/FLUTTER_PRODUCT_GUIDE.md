# Flutter Product API Guide (Customer & Guest)

## Tổng quan

Guide này hướng dẫn cách tích hợp Product API vào Flutter app dành cho **customer và guest**, tập trung vào logic và API calls. Không bao gồm UI implementation.

**Lưu ý:** Guide này chỉ bao gồm các **Public endpoints** (không yêu cầu authentication):
- ✅ Tìm kiếm products
- ✅ Lấy thông tin product theo ID
- ❌ Không bao gồm các endpoints admin-only (create, update, delete, upload image)

## Cấu trúc dữ liệu

### Product Model

```dart
class Product {
  final int id;
  final String name;
  final String? nameEn;
  final String? description;
  final String? descriptionEn;
  final double price;
  final int stock;
  final String? image;
  final List<String> images;
  final int sold;
  final double rating;
  final int reviewCount;
  final bool isActive;
  final String? sku;
  final int categoryId;
  final Category? category;
  final DateTime createdAt;
  final DateTime updatedAt;

  Product({
    required this.id,
    required this.name,
    this.nameEn,
    this.description,
    this.descriptionEn,
    required this.price,
    required this.stock,
    this.image,
    required this.images,
    required this.sold,
    required this.rating,
    required this.reviewCount,
    required this.isActive,
    this.sku,
    required this.categoryId,
    this.category,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'] as int,
      name: json['name'] as String,
      nameEn: json['nameEn'] as String?,
      description: json['description'] as String?,
      descriptionEn: json['descriptionEn'] as String?,
      price: (json['price'] as num).toDouble(),
      stock: json['stock'] as int,
      image: json['image'] as String?,
      images: (json['images'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList() ?? [],
      sold: json['sold'] as int,
      rating: (json['rating'] as num).toDouble(),
      reviewCount: json['reviewCount'] as int,
      isActive: json['isActive'] as bool,
      sku: json['sku'] as String?,
      categoryId: json['categoryId'] as int,
      category: json['category'] != null
          ? Category.fromJson(json['category'] as Map<String, dynamic>)
          : null,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'nameEn': nameEn,
      'description': description,
      'descriptionEn': descriptionEn,
      'price': price,
      'stock': stock,
      'image': image,
      'images': images,
      'sold': sold,
      'rating': rating,
      'reviewCount': reviewCount,
      'isActive': isActive,
      'sku': sku,
      'categoryId': categoryId,
      'category': category?.toJson(),
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  // Helper methods
  bool get isInStock => stock > 0;
  bool get isOutOfStock => stock == 0;
  bool get hasImages => images.isNotEmpty || image != null;
  
  // Get display name based on language
  String getDisplayName(String language) {
    if (language == 'en' && nameEn != null && nameEn!.isNotEmpty) {
      return nameEn!;
    }
    return name;
  }
  
  // Get display description based on language
  String? getDisplayDescription(String language) {
    if (language == 'en' && descriptionEn != null && descriptionEn!.isNotEmpty) {
      return descriptionEn;
    }
    return description;
  }

  // Get primary image (first image in images array, or fallback to image field)
  String? get primaryImage {
    if (images.isNotEmpty) return images.first;
    return image;
  }

  // Get all images (combine image and images)
  List<String> get allImages {
    final all = <String>[];
    if (image != null) all.add(image!);
    all.addAll(images);
    return all;
  }

  // Format price
  String getFormattedPrice({String locale = 'vi_VN'}) {
    return NumberFormat.currency(
      locale: locale,
      symbol: '₫',
      decimalDigits: 0,
    ).format(price);
  }
}
```

**Note:** Import `intl` package for `NumberFormat`:
```dart
import 'package:intl/intl.dart';
```

### Response Models

```dart
class ProductListResponse {
  final bool success;
  final String? message;
  final List<Product> data;

  ProductListResponse({
    required this.success,
    this.message,
    required this.data,
  });

  factory ProductListResponse.fromJson(Map<String, dynamic> json) {
    return ProductListResponse(
      success: json['success'] as bool,
      message: json['message'] as String?,
      data: (json['data'] as List<dynamic>)
          .map((e) => Product.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}

class ProductPaginationResponse {
  final bool success;
  final String? message;
  final List<Product> data;
  final int total;
  final int page;
  final int limit;
  final int totalPages;

  ProductPaginationResponse({
    required this.success,
    this.message,
    required this.data,
    required this.total,
    required this.page,
    required this.limit,
    required this.totalPages,
  });

  factory ProductPaginationResponse.fromJson(Map<String, dynamic> json) {
    return ProductPaginationResponse(
      success: json['success'] as bool,
      message: json['message'] as String?,
      data: (json['data'] as List<dynamic>)
          .map((e) => Product.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      page: json['page'] as int,
      limit: json['limit'] as int,
      totalPages: json['totalPages'] as int,
    );
  }

  bool get hasMore => page < totalPages;
  bool get isLastPage => page >= totalPages;
}

class ProductDetailResponse {
  final bool success;
  final String? message;
  final Product data;

  ProductDetailResponse({
    required this.success,
    this.message,
    required this.data,
  });

  factory ProductDetailResponse.fromJson(Map<String, dynamic> json) {
    return ProductDetailResponse(
      success: json['success'] as bool,
      message: json['message'] as String?,
      data: Product.fromJson(json['data'] as Map<String, dynamic>),
    );
  }
}
```

### Request Models

```dart
class SearchProductRequest {
  final String? name;
  final int? categoryId;
  final int? parentCategoryId;
  final dynamic isActive; // bool, List<bool>, or null
  final double? minPrice;
  final double? maxPrice;
  final bool? inStock;
  final String? sortBy; // id, name, price, stock, createdAt, updatedAt
  final String? sortOrder; // ASC, DESC
  final int? page;
  final int? limit; // Max 1000

  SearchProductRequest({
    this.name,
    this.categoryId,
    this.parentCategoryId,
    this.isActive,
    this.minPrice,
    this.maxPrice,
    this.inStock,
    this.sortBy,
    this.sortOrder,
    this.page,
    this.limit,
  });

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{};
    if (name != null) map['name'] = name;
    if (categoryId != null) map['categoryId'] = categoryId;
    if (parentCategoryId != null) map['parentCategoryId'] = parentCategoryId;
    if (isActive != null) {
      if (isActive is bool) {
        map['isActive'] = isActive;
      } else if (isActive is List<bool>) {
        map['isActive'] = isActive;
      }
    }
    if (minPrice != null) map['minPrice'] = minPrice;
    if (maxPrice != null) map['maxPrice'] = maxPrice;
    if (inStock != null) map['inStock'] = inStock;
    if (sortBy != null) map['sortBy'] = sortBy;
    if (sortOrder != null) map['sortOrder'] = sortOrder;
    if (page != null) map['page'] = page;
    if (limit != null) map['limit'] = limit;
    return map;
  }

  SearchProductRequest copyWith({
    String? name,
    int? categoryId,
    int? parentCategoryId,
    dynamic isActive,
    double? minPrice,
    double? maxPrice,
    bool? inStock,
    String? sortBy,
    String? sortOrder,
    int? page,
    int? limit,
  }) {
    return SearchProductRequest(
      name: name ?? this.name,
      categoryId: categoryId ?? this.categoryId,
      parentCategoryId: parentCategoryId ?? this.parentCategoryId,
      isActive: isActive ?? this.isActive,
      minPrice: minPrice ?? this.minPrice,
      maxPrice: maxPrice ?? this.maxPrice,
      inStock: inStock ?? this.inStock,
      sortBy: sortBy ?? this.sortBy,
      sortOrder: sortOrder ?? this.sortOrder,
      page: page ?? this.page,
      limit: limit ?? this.limit,
    );
  }
}
```

## API Service

### Base API Configuration

```dart
class ApiConfig {
  static const String baseUrl = 'https://your-api-domain.com/api/v1';
  static const Duration timeout = Duration(seconds: 30);
  
  // Language preference (can be stored in SharedPreferences)
  static String getLanguage() => 'vi'; // or 'en'
}
```

### HTTP Client Setup

```dart
import 'package:http/http.dart' as http;
import 'dart:convert';

class ProductApiService {
  final String baseUrl;
  final http.Client client;
  
  ProductApiService({
    required this.baseUrl,
    http.Client? client,
  }) : client = client ?? http.Client();

  // Helper method to get headers
  Map<String, String> _getHeaders() {
    return {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
  }

  // Helper method to handle response
  Map<String, dynamic> _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      return json.decode(response.body) as Map<String, dynamic>;
    } else {
      final error = json.decode(response.body) as Map<String, dynamic>;
      throw ApiException(
        message: error['error'] as String? ?? 'Unknown error',
        statusCode: response.statusCode,
      );
    }
  }
}
```

### Custom Exception

```dart
class ApiException implements Exception {
  final String message;
  final int? statusCode;

  ApiException({
    required this.message,
    this.statusCode,
  });

  @override
  String toString() => 'ApiException: $message (Status: $statusCode)';
}
```

## API Methods

### 1. Tìm kiếm Products (Public)

```dart
Future<ProductPaginationResponse> searchProducts({
  SearchProductRequest? request,
  String language = 'vi',
}) async {
  final uri = Uri.parse('$baseUrl/products/search')
      .replace(queryParameters: {'language': language});

  final body = request?.toJson() ?? <String, dynamic>{};
  
  final response = await client.post(
    uri,
    headers: _getHeaders(),
    body: json.encode(body),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  return ProductPaginationResponse.fromJson(json);
}
```

**Lưu ý:**
- Body có thể null/empty → sẽ lấy tất cả products
- Hỗ trợ filter theo `categoryId` hoặc `parentCategoryId`
- Nếu có cả `categoryId` và `parentCategoryId`, hệ thống ưu tiên `categoryId`
- `isActive` có thể là `bool`, `List<bool>`, hoặc `null`

### 2. Lấy một Product theo ID (Public)

```dart
Future<Product> getProductById({
  required int id,
  String language = 'vi',
  bool includeInactive = false,
}) async {
  final uri = Uri.parse('$baseUrl/products/$id')
      .replace(queryParameters: {
    'language': language,
    'includeInactive': includeInactive.toString(),
  });

  final response = await client.get(
    uri,
    headers: _getHeaders(),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final productResponse = ProductDetailResponse.fromJson(json);
  return productResponse.data;
}
```

## Business Logic Helpers

### Product Helper

```dart
class ProductHelper {
  // Filter products by stock status
  static List<Product> filterByStockStatus(
    List<Product> products, {
    bool? inStock,
  }) {
    if (inStock == null) return products;
    return products.where((product) {
      if (inStock) {
        return product.isInStock;
      } else {
        return product.isOutOfStock;
      }
    }).toList();
  }

  // Filter products by price range
  static List<Product> filterByPriceRange(
    List<Product> products, {
    double? minPrice,
    double? maxPrice,
  }) {
    return products.where((product) {
      if (minPrice != null && product.price < minPrice) return false;
      if (maxPrice != null && product.price > maxPrice) return false;
      return true;
    }).toList();
  }

  // Filter products by category
  static List<Product> filterByCategory(
    List<Product> products, {
    int? categoryId,
  }) {
    if (categoryId == null) return products;
    return products.where((product) => product.categoryId == categoryId).toList();
  }

  // Filter products by active status
  static List<Product> filterByActiveStatus(
    List<Product> products, {
    bool? isActive,
  }) {
    if (isActive == null) return products;
    return products.where((product) => product.isActive == isActive).toList();
  }

  // Sort products
  static List<Product> sortProducts(
    List<Product> products, {
    String sortBy = 'createdAt',
    String sortOrder = 'DESC',
  }) {
    final sorted = List<Product>.from(products);
    
    sorted.sort((a, b) {
      int comparison = 0;
      
      switch (sortBy) {
        case 'id':
          comparison = a.id.compareTo(b.id);
          break;
        case 'name':
          comparison = a.name.compareTo(b.name);
          break;
        case 'price':
          comparison = a.price.compareTo(b.price);
          break;
        case 'stock':
          comparison = a.stock.compareTo(b.stock);
          break;
        case 'createdAt':
          comparison = a.createdAt.compareTo(b.createdAt);
          break;
        case 'updatedAt':
          comparison = a.updatedAt.compareTo(b.updatedAt);
          break;
        default:
          comparison = a.createdAt.compareTo(b.createdAt);
      }
      
      return sortOrder == 'ASC' ? comparison : -comparison;
    });
    
    return sorted;
  }

  // Group products by category
  static Map<int, List<Product>> groupByCategory(List<Product> products) {
    final grouped = <int, List<Product>>{};
    for (var product in products) {
      if (!grouped.containsKey(product.categoryId)) {
        grouped[product.categoryId] = [];
      }
      grouped[product.categoryId]!.add(product);
    }
    return grouped;
  }

  // Get products with low stock (below threshold)
  static List<Product> getLowStockProducts(
    List<Product> products, {
    int threshold = 10,
  }) {
    return products.where((product) => product.stock > 0 && product.stock <= threshold).toList();
  }

  // Get out of stock products
  static List<Product> getOutOfStockProducts(List<Product> products) {
    return products.where((product) => product.isOutOfStock).toList();
  }

  // Calculate total value of products
  static double calculateTotalValue(List<Product> products) {
    return products.fold(0.0, (sum, product) => sum + (product.price * product.stock));
  }
}
```

## State Management Patterns

### Using Provider/ChangeNotifier

```dart
class ProductProvider extends ChangeNotifier {
  final ProductApiService _apiService;

  List<Product> _products = [];
  Product? _selectedProduct;
  bool _isLoading = false;
  String? _error;
  int _total = 0;
  int _page = 1;
  int _totalPages = 0;
  bool _hasMore = false;

  ProductProvider({
    required ProductApiService apiService,
  }) : _apiService = apiService;

  // Getters
  List<Product> get products => _products;
  Product? get selectedProduct => _selectedProduct;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get total => _total;
  int get page => _page;
  int get totalPages => _totalPages;
  bool get hasMore => _hasMore;

  // Search products
  Future<void> searchProducts({
    SearchProductRequest? request,
    String language = 'vi',
    bool loadMore = false,
  }) async {
    if (!loadMore) {
      _isLoading = true;
      _error = null;
      _page = 1;
      _products = [];
      notifyListeners();
    }

    try {
      final searchRequest = request?.copyWith(
        page: loadMore ? _page + 1 : 1,
      ) ?? SearchProductRequest(page: loadMore ? _page + 1 : 1);

      final response = await _apiService.searchProducts(
        request: searchRequest,
        language: language,
      );

      if (loadMore) {
        _products.addAll(response.data);
      } else {
        _products = response.data;
      }

      _total = response.total;
      _page = response.page;
      _totalPages = response.totalPages;
      _hasMore = response.hasMore;
      _error = null;
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Load more products
  Future<void> loadMoreProducts({
    SearchProductRequest? request,
    String language = 'vi',
  }) async {
    if (_hasMore && !_isLoading) {
      await searchProducts(request: request, language: language, loadMore: true);
    }
  }

  // Get product by ID
  Future<void> getProductById({
    required int id,
    String language = 'vi',
    bool includeInactive = false,
  }) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _selectedProduct = await _apiService.getProductById(
        id: id,
        language: language,
        includeInactive: includeInactive,
      );
      _error = null;
    } catch (e) {
      _error = e.toString();
      _selectedProduct = null;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Select product
  void selectProduct(Product? product) {
    _selectedProduct = product;
    notifyListeners();
  }

  // Clear products
  void clearProducts() {
    _products = [];
    _selectedProduct = null;
    _error = null;
    _page = 1;
    _totalPages = 0;
    _hasMore = false;
    notifyListeners();
  }
}
```

### Using Riverpod

```dart
final productApiServiceProvider = Provider<ProductApiService>((ref) {
  final baseUrl = ref.watch(apiConfigProvider).baseUrl;
  return ProductApiService(baseUrl: baseUrl);
});

final productSearchProvider = StateNotifierProvider<ProductSearchNotifier, ProductSearchState>((ref) {
  final apiService = ref.watch(productApiServiceProvider);
  return ProductSearchNotifier(apiService: apiService);
});

class ProductSearchState {
  final List<Product> products;
  final bool isLoading;
  final String? error;
  final int page;
  final int totalPages;
  final bool hasMore;
  final SearchProductRequest? currentRequest;

  ProductSearchState({
    required this.products,
    required this.isLoading,
    this.error,
    required this.page,
    required this.totalPages,
    required this.hasMore,
    this.currentRequest,
  });

  ProductSearchState copyWith({
    List<Product>? products,
    bool? isLoading,
    String? error,
    int? page,
    int? totalPages,
    bool? hasMore,
    SearchProductRequest? currentRequest,
  }) {
    return ProductSearchState(
      products: products ?? this.products,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      totalPages: totalPages ?? this.totalPages,
      hasMore: hasMore ?? this.hasMore,
      currentRequest: currentRequest ?? this.currentRequest,
    );
  }
}

class ProductSearchNotifier extends StateNotifier<ProductSearchState> {
  final ProductApiService _apiService;
  String _language = 'vi';

  ProductSearchNotifier({required ProductApiService apiService})
      : _apiService = apiService,
        super(ProductSearchState(
          products: [],
          isLoading: false,
          page: 1,
          totalPages: 0,
          hasMore: false,
        ));

  Future<void> search({
    SearchProductRequest? request,
    String language = 'vi',
    bool loadMore = false,
  }) async {
    _language = language;

    if (!loadMore) {
      state = state.copyWith(
        products: [],
        isLoading: true,
        error: null,
        page: 1,
        currentRequest: request,
      );
    } else {
      state = state.copyWith(isLoading: true, error: null);
    }

    try {
      final searchRequest = (request ?? SearchProductRequest()).copyWith(
        page: loadMore ? state.page + 1 : 1,
      );

      final response = await _apiService.searchProducts(
        request: searchRequest,
        language: language,
      );

      final newProducts = loadMore
          ? [...state.products, ...response.data]
          : response.data;

      state = state.copyWith(
        products: newProducts,
        isLoading: false,
        page: response.page,
        totalPages: response.totalPages,
        hasMore: response.hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  Future<void> loadMore() {
    if (state.hasMore && !state.isLoading) {
      return search(
        request: state.currentRequest,
        language: _language,
        loadMore: true,
      );
    }
    return Future.value();
  }
}

final productDetailProvider = FutureProvider.family<Product, int>((ref, productId) async {
  final apiService = ref.watch(productApiServiceProvider);
  final language = ref.watch(languageProvider);
  return await apiService.getProductById(
    id: productId,
    language: language,
  );
});
```

## Error Handling

### Error Handler

```dart
class ProductErrorHandler {
  static String getErrorMessage(dynamic error) {
    if (error is ApiException) {
      switch (error.statusCode) {
        case 400:
          return 'Dữ liệu không hợp lệ: ${error.message}';
        case 401:
          return 'Bạn cần đăng nhập để thực hiện thao tác này';
        case 403:
          return 'Bạn không có quyền thực hiện thao tác này';
        case 404:
          return 'Không tìm thấy sản phẩm';
        case 500:
          return 'Lỗi server: ${error.message}';
        default:
          return error.message;
      }
    } else if (error is TimeoutException) {
      return 'Kết nối quá lâu. Vui lòng thử lại';
    } else if (error is SocketException) {
      return 'Không có kết nối mạng. Vui lòng kiểm tra kết nối';
    } else {
      return 'Đã xảy ra lỗi: ${error.toString()}';
    }
  }

  static bool isNetworkError(dynamic error) {
    return error is SocketException || error is TimeoutException;
  }

  static bool isAuthError(dynamic error) {
    return error is ApiException && 
           (error.statusCode == 401 || error.statusCode == 403);
  }
}
```

## Caching Strategy

### Simple Cache Implementation

```dart
class ProductCache {
  static final Map<String, CachedData> _cache = {};
  static const Duration defaultCacheDuration = Duration(minutes: 5);
  static const Duration productDetailCacheDuration = Duration(minutes: 10);

  static T? get<T>(String key) {
    final cached = _cache[key];
    if (cached == null) return null;
    if (cached.isExpired) {
      _cache.remove(key);
      return null;
    }
    return cached.data as T?;
  }

  static void set<T>(String key, T data, {Duration? duration}) {
    _cache[key] = CachedData(
      data: data,
      expiresAt: DateTime.now().add(duration ?? defaultCacheDuration),
    );
  }

  static void clear() {
    _cache.clear();
  }

  static void remove(String key) {
    _cache.remove(key);
  }

  // Cache product detail
  static Product? getProductDetail(int productId) {
    return get<Product>('product_detail_$productId');
  }

  static void setProductDetail(int productId, Product product) {
    set('product_detail_$productId', product, duration: productDetailCacheDuration);
  }

  // Cache search results
  static ProductPaginationResponse? getSearchResults(String cacheKey) {
    return get<ProductPaginationResponse>(cacheKey);
  }

  static void setSearchResults(String cacheKey, ProductPaginationResponse response) {
    set(cacheKey, response);
  }

  // Generate cache key from search request
  static String generateSearchCacheKey(SearchProductRequest request, String language) {
    final key = StringBuffer('product_search_$language');
    if (request.name != null) key.write('_name_${request.name}');
    if (request.categoryId != null) key.write('_cat_${request.categoryId}');
    if (request.parentCategoryId != null) key.write('_parent_${request.parentCategoryId}');
    if (request.minPrice != null) key.write('_min_${request.minPrice}');
    if (request.maxPrice != null) key.write('_max_${request.maxPrice}');
    if (request.inStock != null) key.write('_stock_${request.inStock}');
    if (request.page != null) key.write('_page_${request.page}');
    if (request.limit != null) key.write('_limit_${request.limit}');
    return key.toString();
  }
}

class CachedData {
  final dynamic data;
  final DateTime expiresAt;

  CachedData({required this.data, required this.expiresAt});

  bool get isExpired => DateTime.now().isAfter(expiresAt);
}

// Usage in API Service
Future<ProductPaginationResponse> searchProductsWithCache({
  SearchProductRequest? request,
  String language = 'vi',
  bool forceRefresh = false,
}) async {
  final cacheKey = ProductCache.generateSearchCacheKey(
    request ?? SearchProductRequest(),
    language,
  );

  if (!forceRefresh) {
    final cached = ProductCache.getSearchResults(cacheKey);
    if (cached != null) return cached;
  }

  final response = await searchProducts(
    request: request,
    language: language,
  );

  ProductCache.setSearchResults(cacheKey, response);
  return response;
}
```

## Usage Examples

### Example 1: Search Products with Filters

```dart
class ProductSearchService {
  final ProductApiService _apiService;

  ProductSearchService(this._apiService);

  Future<ProductPaginationResponse> search({
    String? name,
    int? categoryId,
    int? parentCategoryId,
    double? minPrice,
    double? maxPrice,
    bool? inStock,
    bool? isActive,
    int page = 1,
    int limit = 20,
    String language = 'vi',
  }) async {
    try {
      final request = SearchProductRequest(
        name: name,
        categoryId: categoryId,
        parentCategoryId: parentCategoryId,
        minPrice: minPrice,
        maxPrice: maxPrice,
        inStock: inStock,
        isActive: isActive,
        page: page,
        limit: limit,
        sortBy: 'createdAt',
        sortOrder: 'DESC',
      );

      return await _apiService.searchProducts(
        request: request,
        language: language,
      );
    } catch (e) {
      throw Exception('Failed to search products: ${ProductErrorHandler.getErrorMessage(e)}');
    }
  }
}
```

### Example 2: Get Products by Parent Category

```dart
class ProductByCategoryService {
  final ProductApiService _apiService;

  ProductByCategoryService(this._apiService);

  Future<ProductPaginationResponse> getProductsByParentCategory({
    required int parentCategoryId,
    String language = 'vi',
    int page = 1,
    int limit = 20,
  }) async {
    try {
      final request = SearchProductRequest(
        parentCategoryId: parentCategoryId,
        isActive: true,
        page: page,
        limit: limit,
        sortBy: 'createdAt',
        sortOrder: 'DESC',
      );

      return await _apiService.searchProducts(
        request: request,
        language: language,
      );
    } catch (e) {
      throw Exception('Failed to get products: ${ProductErrorHandler.getErrorMessage(e)}');
    }
  }
}
```


## Best Practices

1. **Error Handling**: Luôn wrap API calls trong try-catch và xử lý lỗi phù hợp
2. **Loading States**: Quản lý loading state để tránh multiple requests
3. **Caching**: Cache product details và search results để giảm API calls
4. **Pagination**: Luôn sử dụng pagination cho danh sách lớn, implement load more
5. **Language Support**: Luôn truyền language parameter để hỗ trợ đa ngôn ngữ
6. **Network Retry**: Implement retry logic cho network errors
7. **Stock Management**: Kiểm tra stock trước khi thêm vào cart
8. **Price Formatting**: Sử dụng NumberFormat để format giá tiền đúng locale
9. **Filter Optimization**: Sử dụng filter trên server thay vì filter trên client khi có thể
10. **Image Loading**: Sử dụng image caching và lazy loading cho product images

## Dependencies

Thêm vào `pubspec.yaml`:

```yaml
dependencies:
  http: ^1.1.0
  intl: ^0.18.0  # For number formatting
  shared_preferences: ^2.2.0  # For storing language preference, etc.
```

## Notes

- Tất cả API methods đều hỗ trợ language parameter
- Tất cả endpoints trong guide này đều là **Public** (không yêu cầu authentication)
- Search endpoint hỗ trợ filter theo `categoryId` hoặc `parentCategoryId`
- Nếu có cả `categoryId` và `parentCategoryId`, hệ thống ưu tiên `categoryId`
- `isActive` có thể là `bool`, `List<bool>`, hoặc `null` (lấy tất cả)
- Product có thể có nhiều images (array) hoặc một image (field)
- Mặc định chỉ lấy products active, có thể dùng `includeInactive=true` để lấy cả inactive


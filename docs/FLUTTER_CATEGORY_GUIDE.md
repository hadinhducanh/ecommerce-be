# Flutter Category API Guide

## Tổng quan

Guide này hướng dẫn cách tích hợp Category API vào Flutter app, tập trung vào logic và API calls. Không bao gồm UI implementation.

## Cấu trúc dữ liệu

### Category Model

```dart
class Category {
  final int id;
  final String name;
  final String? nameEn;
  final String? description;
  final String? descriptionEn;
  final String? image;
  final bool isActive;
  final int? parentId;
  final List<int> childrenIds;
  final DateTime createdAt;
  final DateTime updatedAt;

  Category({
    required this.id,
    required this.name,
    this.nameEn,
    this.description,
    this.descriptionEn,
    this.image,
    required this.isActive,
    this.parentId,
    required this.childrenIds,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'] as int,
      name: json['name'] as String,
      nameEn: json['nameEn'] as String?,
      description: json['description'] as String?,
      descriptionEn: json['descriptionEn'] as String?,
      image: json['image'] as String?,
      isActive: json['isActive'] as bool,
      parentId: json['parentId'] as int?,
      childrenIds: (json['childrenIds'] as List<dynamic>?)
          ?.map((e) => e as int)
          .toList() ?? [],
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
      'image': image,
      'isActive': isActive,
      'parentId': parentId,
      'childrenIds': childrenIds,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  // Helper methods
  bool get isRoot => parentId == null;
  bool get isChild => parentId != null;
  bool get hasChildren => childrenIds.isNotEmpty;
  
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
}
```

### Response Models

```dart
class CategoryListResponse {
  final bool success;
  final List<Category> data;

  CategoryListResponse({
    required this.success,
    required this.data,
  });

  factory CategoryListResponse.fromJson(Map<String, dynamic> json) {
    return CategoryListResponse(
      success: json['success'] as bool,
      data: (json['data'] as List<dynamic>)
          .map((e) => Category.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }
}

class CategoryPaginationResponse {
  final bool success;
  final String? message;
  final List<Category> data;
  final int total;
  final int page;
  final int limit;
  final int totalPages;

  CategoryPaginationResponse({
    required this.success,
    this.message,
    required this.data,
    required this.total,
    required this.page,
    required this.limit,
    required this.totalPages,
  });

  factory CategoryPaginationResponse.fromJson(Map<String, dynamic> json) {
    return CategoryPaginationResponse(
      success: json['success'] as bool,
      message: json['message'] as String?,
      data: (json['data'] as List<dynamic>)
          .map((e) => Category.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      page: json['page'] as int,
      limit: json['limit'] as int,
      totalPages: json['totalPages'] as int,
    );
  }
}

class CategoryDetailResponse {
  final bool success;
  final String? message;
  final Category data;

  CategoryDetailResponse({
    required this.success,
    this.message,
    required this.data,
  });

  factory CategoryDetailResponse.fromJson(Map<String, dynamic> json) {
    return CategoryDetailResponse(
      success: json['success'] as bool,
      message: json['message'] as String?,
      data: Category.fromJson(json['data'] as Map<String, dynamic>),
    );
  }
}
```

### Request Models

```dart
class SearchCategoryRequest {
  final String? name;
  final bool? isActive;
  final List<bool>? isActiveArray;
  final String? sortBy;
  final String? sortOrder;
  final int? page;
  final int? limit;

  SearchCategoryRequest({
    this.name,
    this.isActive,
    this.isActiveArray,
    this.sortBy,
    this.sortOrder,
    this.page,
    this.limit,
  });

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{};
    if (name != null) map['name'] = name;
    if (isActive != null) {
      map['isActive'] = isActive;
    } else if (isActiveArray != null) {
      map['isActive'] = isActiveArray;
    }
    if (sortBy != null) map['sortBy'] = sortBy;
    if (sortOrder != null) map['sortOrder'] = sortOrder;
    if (page != null) map['page'] = page;
    if (limit != null) map['limit'] = limit;
    return map;
  }
}

class SearchCategoryChildRequest {
  final int? parentId;
  final String? name;
  final bool? isActive;
  final List<bool>? isActiveArray;
  final String? sortBy;
  final String? sortOrder;
  final int? page;
  final int? limit;

  SearchCategoryChildRequest({
    this.parentId,
    this.name,
    this.isActive,
    this.isActiveArray,
    this.sortBy,
    this.sortOrder,
    this.page,
    this.limit,
  });

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{};
    if (parentId != null) map['parentId'] = parentId;
    if (name != null) map['name'] = name;
    if (isActive != null) {
      map['isActive'] = isActive;
    } else if (isActiveArray != null) {
      map['isActive'] = isActiveArray;
    }
    if (sortBy != null) map['sortBy'] = sortBy;
    if (sortOrder != null) map['sortOrder'] = sortOrder;
    if (page != null) map['page'] = page;
    if (limit != null) map['limit'] = limit;
    return map;
  }
}

class CreateCategoryRequest {
  final String name;
  final String? nameEn;
  final String? description;
  final String? descriptionEn;
  final String? image;
  final bool? isActive;
  final int? parentId;

  CreateCategoryRequest({
    required this.name,
    this.nameEn,
    this.description,
    this.descriptionEn,
    this.image,
    this.isActive,
    this.parentId,
  });

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{
      'name': name,
    };
    if (nameEn != null) map['nameEn'] = nameEn;
    if (description != null) map['description'] = description;
    if (descriptionEn != null) map['descriptionEn'] = descriptionEn;
    if (image != null) map['image'] = image;
    if (isActive != null) map['isActive'] = isActive;
    if (parentId != null) map['parentId'] = parentId;
    return map;
  }
}

class UpdateCategoryRequest {
  final String? name;
  final String? nameEn;
  final String? description;
  final String? descriptionEn;
  final String? image;
  final bool? isActive;

  UpdateCategoryRequest({
    this.name,
    this.nameEn,
    this.description,
    this.descriptionEn,
    this.image,
    this.isActive,
  });

  Map<String, dynamic> toJson() {
    final map = <String, dynamic>{};
    if (name != null) map['name'] = name;
    if (nameEn != null) map['nameEn'] = nameEn;
    if (description != null) map['description'] = description;
    if (descriptionEn != null) map['descriptionEn'] = descriptionEn;
    if (image != null) map['image'] = image;
    if (isActive != null) map['isActive'] = isActive;
    return map;
  }
}

class AddChildRequest {
  final int childId;

  AddChildRequest({required this.childId});

  Map<String, dynamic> toJson() => {'childId': childId};
}

class RemoveChildRequest {
  final int childId;

  RemoveChildRequest({required this.childId});

  Map<String, dynamic> toJson() => {'childId': childId};
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

class CategoryApiService {
  final String baseUrl;
  final http.Client client;
  
  CategoryApiService({
    required this.baseUrl,
    http.Client? client,
  }) : client = client ?? http.Client();

  // Helper method to get headers
  Map<String, String> _getHeaders({String? token}) {
    final headers = <String, String>{
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    if (token != null) {
      headers['Authorization'] = 'Bearer $token';
    }
    return headers;
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

### 1. Lấy danh sách Parent Categories (cho Dropdown)

```dart
Future<List<Category>> getParentCategories({
  String language = 'vi',
  bool includeInactive = false,
  String? token,
}) async {
  final uri = Uri.parse('$baseUrl/categories/parents')
      .replace(queryParameters: {
    'language': language,
    'includeInactive': includeInactive.toString(),
  });

  final response = await client.get(
    uri,
    headers: _getHeaders(token: token),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryListResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 2. Lấy danh sách tất cả Child Categories (cho Dropdown)

```dart
Future<List<Category>> getAllChildren({
  String language = 'vi',
  bool includeInactive = true,
  String? token,
}) async {
  final uri = Uri.parse('$baseUrl/categories/children')
      .replace(queryParameters: {
    'language': language,
    'includeInactive': includeInactive.toString(),
  });

  final response = await client.get(
    uri,
    headers: _getHeaders(token: token),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryListResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 3. Tìm kiếm Categories (Parent Categories)

```dart
Future<CategoryPaginationResponse> searchCategories({
  SearchCategoryRequest? request,
  String language = 'vi',
  String? token,
}) async {
  final uri = Uri.parse('$baseUrl/categories/search')
      .replace(queryParameters: {'language': language});

  final body = request?.toJson() ?? <String, dynamic>{};
  
  final response = await client.post(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(body),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  return CategoryPaginationResponse.fromJson(json);
}
```

### 4. Tìm kiếm Category Children

```dart
Future<CategoryPaginationResponse> searchChildren({
  SearchCategoryChildRequest? request,
  String language = 'vi',
  String? token,
}) async {
  final uri = Uri.parse('$baseUrl/categories/children/search')
      .replace(queryParameters: {'language': language});

  final body = request?.toJson() ?? <String, dynamic>{};
  
  final response = await client.post(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(body),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  return CategoryPaginationResponse.fromJson(json);
}
```

### 5. Lấy một Category theo ID

```dart
Future<Category> getCategoryById({
  required int id,
  String language = 'vi',
  bool includeInactive = true,
  String? token,
}) async {
  final uri = Uri.parse('$baseUrl/categories/$id')
      .replace(queryParameters: {
    'language': language,
    'includeInactive': includeInactive.toString(),
  });

  final response = await client.get(
    uri,
    headers: _getHeaders(token: token),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryDetailResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 6. Lấy danh sách Children của một Parent

```dart
Future<List<Category>> getChildrenByParentId({
  required int parentId,
  String language = 'vi',
  bool includeInactive = true,
  String? token,
}) async {
  final uri = Uri.parse('$baseUrl/categories/$parentId/children')
      .replace(queryParameters: {
    'language': language,
    'includeInactive': includeInactive.toString(),
  });

  final response = await client.get(
    uri,
    headers: _getHeaders(token: token),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryListResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 7. Tạo Category mới (Admin Only)

```dart
Future<Category> createCategory({
  required CreateCategoryRequest request,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories');

  final response = await client.post(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(request.toJson()),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryDetailResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 8. Cập nhật Category (Partial Update - Admin Only)

```dart
Future<Category> updateCategory({
  required int id,
  required UpdateCategoryRequest request,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories/$id');

  final response = await client.patch(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(request.toJson()),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryDetailResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 9. Thay thế Category (Full Update - Admin Only)

```dart
Future<Category> replaceCategory({
  required int id,
  required CreateCategoryRequest request,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories/$id');

  final response = await client.put(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(request.toJson()),
  ).timeout(ApiConfig.timeout);

  final json = _handleResponse(response);
  final categoryResponse = CategoryDetailResponse.fromJson(json);
  return categoryResponse.data;
}
```

### 10. Xóa Category (Soft Delete - Admin Only)

```dart
Future<void> deleteCategory({
  required int id,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories/$id');

  final response = await client.delete(
    uri,
    headers: _getHeaders(token: token),
  ).timeout(ApiConfig.timeout);

  _handleResponse(response);
}
```

### 11. Xóa vĩnh viễn Category (Hard Delete - Admin Only)

```dart
Future<void> hardDeleteCategory({
  required int id,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories/$id/hard');

  final response = await client.delete(
    uri,
    headers: _getHeaders(token: token),
  ).timeout(ApiConfig.timeout);

  _handleResponse(response);
}
```

### 12. Thêm Child vào Parent (Admin Only)

```dart
Future<void> addChildToParent({
  required int parentId,
  required int childId,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories/$parentId/children');

  final response = await client.post(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(AddChildRequest(childId: childId).toJson()),
  ).timeout(ApiConfig.timeout);

  _handleResponse(response);
}
```

### 13. Xóa Child khỏi Parent (Admin Only)

```dart
Future<void> removeChildFromParent({
  required int parentId,
  required int childId,
  String? token,
}) async {
  if (token == null) {
    throw ApiException(message: 'Authentication required', statusCode: 401);
  }

  final uri = Uri.parse('$baseUrl/categories/$parentId/children');

  final response = await client.delete(
    uri,
    headers: _getHeaders(token: token),
    body: json.encode(RemoveChildRequest(childId: childId).toJson()),
  ).timeout(ApiConfig.timeout);

  _handleResponse(response);
}
```

## Business Logic Helpers

### Category Tree Builder

```dart
class CategoryTree {
  final Category category;
  final List<CategoryTree> children;

  CategoryTree({
    required this.category,
    required this.children,
  });
}

class CategoryHelper {
  // Build category tree from flat list
  static List<CategoryTree> buildTree(List<Category> categories) {
    final categoryMap = <int, CategoryTree>{};
    final roots = <CategoryTree>[];

    // Create tree nodes
    for (var category in categories) {
      categoryMap[category.id] = CategoryTree(
        category: category,
        children: [],
      );
    }

    // Build tree structure
    for (var category in categories) {
      final node = categoryMap[category.id]!;
      if (category.parentId == null) {
        roots.add(node);
      } else {
        final parent = categoryMap[category.parentId];
        if (parent != null) {
          parent.children.add(node);
        }
      }
    }

    return roots;
  }

  // Get all descendants of a category
  static List<Category> getAllDescendants(
    Category category,
    List<Category> allCategories,
  ) {
    final descendants = <Category>[];
    final queue = <Category>[category];

    while (queue.isNotEmpty) {
      final current = queue.removeAt(0);
      final children = allCategories
          .where((c) => current.childrenIds.contains(c.id))
          .toList();
      descendants.addAll(children);
      queue.addAll(children);
    }

    return descendants;
  }

  // Get all ancestors of a category
  static List<Category> getAllAncestors(
    Category category,
    List<Category> allCategories,
  ) {
    final ancestors = <Category>[];
    Category? current = category;

    while (current != null && current.parentId != null) {
      final parent = allCategories.firstWhere(
        (c) => c.id == current!.parentId,
        orElse: () => throw Exception('Parent not found'),
      );
      ancestors.insert(0, parent);
      current = parent;
    }

    return ancestors;
  }

  // Check if category is ancestor of another
  static bool isAncestorOf(Category ancestor, Category descendant, List<Category> allCategories) {
    final ancestors = getAllAncestors(descendant, allCategories);
    return ancestors.any((a) => a.id == ancestor.id);
  }

  // Filter categories by type
  static List<Category> filterByType(
    List<Category> categories, {
    bool? isRoot,
    bool? isChild,
    bool? hasChildren,
    bool? isActive,
  }) {
    return categories.where((category) {
      if (isRoot != null && category.isRoot != isRoot) return false;
      if (isChild != null && category.isChild != isChild) return false;
      if (hasChildren != null && category.hasChildren != hasChildren) return false;
      if (isActive != null && category.isActive != isActive) return false;
      return true;
    }).toList();
  }
}
```

## State Management Patterns

### Using Provider/ChangeNotifier

```dart
class CategoryProvider extends ChangeNotifier {
  final CategoryApiService _apiService;
  final String? _token;

  List<Category> _parentCategories = [];
  List<Category> _childCategories = [];
  Category? _selectedCategory;
  bool _isLoading = false;
  String? _error;

  CategoryProvider({
    required CategoryApiService apiService,
    String? token,
  })  : _apiService = apiService,
        _token = token;

  // Getters
  List<Category> get parentCategories => _parentCategories;
  List<Category> get childCategories => _childCategories;
  Category? get selectedCategory => _selectedCategory;
  bool get isLoading => _isLoading;
  String? get error => _error;

  // Load parent categories
  Future<void> loadParentCategories({
    String language = 'vi',
    bool includeInactive = false,
  }) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _parentCategories = await _apiService.getParentCategories(
        language: language,
        includeInactive: includeInactive,
        token: _token,
      );
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Load child categories
  Future<void> loadChildCategories({
    String language = 'vi',
    bool includeInactive = true,
  }) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _childCategories = await _apiService.getAllChildren(
        language: language,
        includeInactive: includeInactive,
        token: _token,
      );
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Search categories
  Future<CategoryPaginationResponse> searchCategories({
    SearchCategoryRequest? request,
    String language = 'vi',
  }) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _apiService.searchCategories(
        request: request,
        language: language,
        token: _token,
      );
      return response;
    } catch (e) {
      _error = e.toString();
      rethrow;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  // Select category
  void selectCategory(Category? category) {
    _selectedCategory = category;
    notifyListeners();
  }
}
```

### Using Riverpod

```dart
final categoryApiServiceProvider = Provider<CategoryApiService>((ref) {
  final baseUrl = ref.watch(apiConfigProvider).baseUrl;
  return CategoryApiService(baseUrl: baseUrl);
});

final parentCategoriesProvider = FutureProvider<List<Category>>((ref) async {
  final apiService = ref.watch(categoryApiServiceProvider);
  final language = ref.watch(languageProvider);
  return await apiService.getParentCategories(language: language);
});

final childCategoriesProvider = FutureProvider<List<Category>>((ref) async {
  final apiService = ref.watch(categoryApiServiceProvider);
  final language = ref.watch(languageProvider);
  return await apiService.getAllChildren(language: language);
});

final categorySearchProvider = StateNotifierProvider<CategorySearchNotifier, CategorySearchState>((ref) {
  final apiService = ref.watch(categoryApiServiceProvider);
  return CategorySearchNotifier(apiService: apiService);
});

class CategorySearchState {
  final List<Category> categories;
  final bool isLoading;
  final String? error;
  final int page;
  final int totalPages;
  final bool hasMore;

  CategorySearchState({
    required this.categories,
    required this.isLoading,
    this.error,
    required this.page,
    required this.totalPages,
    required this.hasMore,
  });

  CategorySearchState copyWith({
    List<Category>? categories,
    bool? isLoading,
    String? error,
    int? page,
    int? totalPages,
    bool? hasMore,
  }) {
    return CategorySearchState(
      categories: categories ?? this.categories,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      totalPages: totalPages ?? this.totalPages,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

class CategorySearchNotifier extends StateNotifier<CategorySearchState> {
  final CategoryApiService _apiService;
  String _language = 'vi';
  SearchCategoryRequest? _currentRequest;

  CategorySearchNotifier({required CategoryApiService apiService})
      : _apiService = apiService,
        super(CategorySearchState(
          categories: [],
          isLoading: false,
          page: 1,
          totalPages: 0,
          hasMore: false,
        ));

  Future<void> search({
    SearchCategoryRequest? request,
    String language = 'vi',
    bool loadMore = false,
  }) async {
    _language = language;
    _currentRequest = request;

    if (!loadMore) {
      state = state.copyWith(
        categories: [],
        isLoading: true,
        error: null,
        page: 1,
      );
    } else {
      state = state.copyWith(isLoading: true, error: null);
    }

    try {
      final searchRequest = (request ?? SearchCategoryRequest()).copyWith(
        page: loadMore ? state.page + 1 : 1,
      );

      final response = await _apiService.searchCategories(
        request: searchRequest,
        language: language,
      );

      final newCategories = loadMore
          ? [...state.categories, ...response.data]
          : response.data;

      state = state.copyWith(
        categories: newCategories,
        isLoading: false,
        page: response.page,
        totalPages: response.totalPages,
        hasMore: response.page < response.totalPages,
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
        request: _currentRequest,
        language: _language,
        loadMore: true,
      );
    }
    return Future.value();
  }
}
```

## Error Handling

### Error Handler

```dart
class CategoryErrorHandler {
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
          return 'Không tìm thấy danh mục';
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
class CategoryCache {
  static final Map<String, CachedData> _cache = {};
  static const Duration defaultCacheDuration = Duration(minutes: 5);

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
}

class CachedData {
  final dynamic data;
  final DateTime expiresAt;

  CachedData({required this.data, required this.expiresAt});

  bool get isExpired => DateTime.now().isAfter(expiresAt);
}

// Usage in API Service
Future<List<Category>> getParentCategoriesWithCache({
  String language = 'vi',
  bool includeInactive = false,
  String? token,
  bool forceRefresh = false,
}) async {
  final cacheKey = 'parent_categories_${language}_$includeInactive';
  
  if (!forceRefresh) {
    final cached = CategoryCache.get<List<Category>>(cacheKey);
    if (cached != null) return cached;
  }

  final categories = await getParentCategories(
    language: language,
    includeInactive: includeInactive,
    token: token,
  );

  CategoryCache.set(cacheKey, categories);
  return categories;
}
```

## Usage Examples

### Example 1: Load Parent Categories for Dropdown

```dart
class CategoryDropdownService {
  final CategoryApiService _apiService;

  CategoryDropdownService(this._apiService);

  Future<List<Category>> loadParentOptions({
    String language = 'vi',
    String? token,
  }) async {
    try {
      final categories = await _apiService.getParentCategories(
        language: language,
        includeInactive: false,
        token: token,
      );
      return categories;
    } catch (e) {
      throw Exception('Failed to load parent categories: ${CategoryErrorHandler.getErrorMessage(e)}');
    }
  }
}
```

### Example 2: Search Categories with Pagination

```dart
class CategorySearchService {
  final CategoryApiService _apiService;

  CategorySearchService(this._apiService);

  Future<CategoryPaginationResponse> search({
    String? name,
    bool? isActive,
    int page = 1,
    int limit = 10,
    String language = 'vi',
    String? token,
  }) async {
    try {
      final request = SearchCategoryRequest(
        name: name,
        isActive: isActive,
        page: page,
        limit: limit,
        sortBy: 'createdAt',
        sortOrder: 'DESC',
      );

      return await _apiService.searchCategories(
        request: request,
        language: language,
        token: token,
      );
    } catch (e) {
      throw Exception('Failed to search categories: ${CategoryErrorHandler.getErrorMessage(e)}');
    }
  }
}
```

### Example 3: Build Category Tree

```dart
class CategoryTreeService {
  final CategoryApiService _apiService;

  CategoryTreeService(this._apiService);

  Future<List<CategoryTree>> buildFullTree({
    String language = 'vi',
    String? token,
  }) async {
    try {
      // Load all categories
      final allCategoriesResponse = await _apiService.searchCategories(
        request: SearchCategoryRequest(limit: 1000),
        language: language,
        token: token,
      );

      // Build tree
      return CategoryHelper.buildTree(allCategoriesResponse.data);
    } catch (e) {
      throw Exception('Failed to build category tree: ${CategoryErrorHandler.getErrorMessage(e)}');
    }
  }
}
```

### Example 4: Filter Products by Parent Category

```dart
// This requires Product API service (not included in this guide)
// But shows how to use parentCategoryId filter

class ProductFilterService {
  final ProductApiService _productApiService;

  ProductFilterService(this._productApiService);

  Future<ProductPaginationResponse> getProductsByParentCategory({
    required int parentCategoryId,
    String language = 'vi',
    int page = 1,
    int limit = 20,
  }) async {
    final request = SearchProductRequest(
      parentCategoryId: parentCategoryId,
      isActive: true,
      page: page,
      limit: limit,
    );

    return await _productApiService.searchProducts(
      request: request,
      language: language,
    );
  }
}
```

## Best Practices

1. **Error Handling**: Luôn wrap API calls trong try-catch và xử lý lỗi phù hợp
2. **Loading States**: Quản lý loading state để tránh multiple requests
3. **Caching**: Cache dữ liệu ít thay đổi (như parent categories) để giảm API calls
4. **Pagination**: Luôn sử dụng pagination cho danh sách lớn
5. **Language Support**: Luôn truyền language parameter để hỗ trợ đa ngôn ngữ
6. **Token Management**: Lưu token an toàn và refresh khi cần
7. **Network Retry**: Implement retry logic cho network errors
8. **Offline Support**: Cache dữ liệu quan trọng để hỗ trợ offline mode

## Dependencies

Thêm vào `pubspec.yaml`:

```yaml
dependencies:
  http: ^1.1.0
  shared_preferences: ^2.2.0  # For storing language preference, token, etc.
```

## Notes

- Tất cả API methods đều hỗ trợ language parameter
- Admin-only endpoints yêu cầu authentication token
- Parent categories là các categories không có parentId (root categories)
- Child categories là các categories có parentId
- Một category có thể có nhiều children nhưng chỉ có một parent


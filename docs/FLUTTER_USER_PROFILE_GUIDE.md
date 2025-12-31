# Flutter User Profile Management Guide

## Mục Đích

Guide này hướng dẫn các Flutter developer cách thao tác với tính năng chỉnh sửa thông tin cá nhân của người dùng. Người dùng chỉ có thể chỉnh sửa thông tin của chính mình, khác với admin có thể quản lý nhiều người dùng.

---

## Thông Tin Chung

| Thông Tin | Chi Tiết |
|-----------|----------|
| **Auth Required** | ✅ Bắt buộc (JWT Token) |
| **Base URL** | `/api/v1/users` |
| **Scope** | Chỉ quản lý thông tin của chính người dùng |
| **Headers** | `Authorization: Bearer {jwt_token}` |

---

## 1. Lấy Thông Tin Profile

### Endpoint
```
GET /api/v1/users/profile
```

### Mô Tả
Lấy thông tin profile đầy đủ của người dùng hiện tại (chính người dùng đó).

### Headers
```
Authorization: Bearer {jwt_token}
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "Lấy thông tin cá nhân thành công",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "Nguyễn Văn A",
    "role": "customer",
    "phone": "0123456789",
    "avatar": "https://example.com/avatar.jpg",
    "address": "123 Đường ABC, Quận 1, TP.HCM",
    "gender": "male",
    "isEmailVerified": true,
    "isActive": true,
    "isFirstLogin": false,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-12-29T15:45:30Z"
  }
}
```

### Các Trường Thông Tin
| Field | Kiểu | Mô Tả | Có thể chỉnh sửa? |
|-------|------|-------|------------------|
| `id` | uint | ID người dùng | ❌ Không |
| `email` | string | Email đăng nhập | ❌ Không |
| `name` | string | Tên đầy đủ | ✅ Có |
| `role` | string | Vai trò (customer/admin) | ❌ Không |
| `phone` | string | Số điện thoại | ✅ Có |
| `avatar` | string | URL ảnh đại diện | ✅ Có |
| `address` | string | Địa chỉ | ✅ Có |
| `gender` | string | Giới tính (male/female/other) | ✅ Có |
| `isEmailVerified` | boolean | Xác thực email | ❌ Không |
| `isActive` | boolean | Tài khoản hoạt động | ❌ Không |
| `isFirstLogin` | boolean | Lần đăng nhập đầu tiên | ❌ Không |
| `createdAt` | string | Ngày tạo (ISO 8601) | ❌ Không |
| `updatedAt` | string | Lần cập nhật cuối (ISO 8601) | ❌ Không |

### Error Cases
```json
{
  "success": false,
  "error": "không tìm thấy người dùng"
}
```

### Flutter Implementation Example
```dart
class UserProfileService {
  final dio = Dio();
  final storage = GetStorage();

  Future<UserProfile> getProfile() async {
    try {
      final token = storage.read('jwt_token');
      
      final response = await dio.get(
        'http://localhost:8080/api/v1/users/profile',
        options: Options(
          headers: {'Authorization': 'Bearer $token'},
        ),
      );

      if (response.statusCode == 200) {
        return UserProfile.fromJson(response.data['data']);
      }
    } catch (e) {
      throw Exception('Không thể lấy thông tin profile: $e');
    }
  }
}

// Model
class UserProfile {
  final int id;
  final String email;
  final String name;
  final String role;
  final String? phone;
  final String? avatar;
  final String? address;
  final String? gender;
  final bool isEmailVerified;
  final bool isActive;
  final bool isFirstLogin;
  final DateTime createdAt;
  final DateTime updatedAt;

  UserProfile({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
    this.phone,
    this.avatar,
    this.address,
    this.gender,
    required this.isEmailVerified,
    required this.isActive,
    required this.isFirstLogin,
    required this.createdAt,
    required this.updatedAt,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      id: json['id'] as int,
      email: json['email'] as String,
      name: json['name'] as String,
      role: json['role'] as String,
      phone: json['phone'] as String?,
      avatar: json['avatar'] as String?,
      address: json['address'] as String?,
      gender: json['gender'] as String?,
      isEmailVerified: json['isEmailVerified'] as bool? ?? false,
      isActive: json['isActive'] as bool? ?? true,
      isFirstLogin: json['isFirstLogin'] as bool? ?? false,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}
```

---

## 2. Cập Nhật Thông Tin Profile

### Endpoint
```
PATCH /api/v1/users/profile
```

### Mô Tả
Cập nhật thông tin cá nhân của người dùng hiện tại. Người dùng chỉ có thể cập nhật các trường được phép.

### Headers
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

### Request Body
**Các trường là tùy chọn (optional)** - chỉ gửi những trường cần cập nhật:

```json
{
  "name": "Nguyễn Văn A",
  "phone": "0123456789",
  "avatar": "https://example.com/new-avatar.jpg",
  "address": "456 Đường XYZ, Quận 2, TP.HCM",
  "gender": "male"
}
```

### Validation Rules
| Field | Validation |
|-------|-----------|
| `name` | Không bắt buộc, kiểu string |
| `phone` | Không bắt buộc, kiểu string |
| `avatar` | Không bắt buộc, kiểu string (URL hợp lệ) |
| `address` | Không bắt buộc, kiểu string |
| `gender` | Không bắt buộc, phải là: `male`, `female`, hoặc `other` |

### Response (200 OK)
```json
{
  "success": true,
  "message": "Cập nhật thông tin cá nhân thành công",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "Nguyễn Văn A",
    "role": "customer",
    "phone": "0123456789",
    "avatar": "https://example.com/new-avatar.jpg",
    "address": "456 Đường XYZ, Quận 2, TP.HCM",
    "gender": "male",
    "isEmailVerified": true,
    "isActive": true,
    "isFirstLogin": false,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-12-29T16:00:00Z"
  }
}
```

### Error Cases

#### Invalid Gender Value
```json
{
  "success": false,
  "error": "Dữ liệu không hợp lệ",
  "details": "Key: 'UpdateProfileRequest.Gender' Error:Field validation for 'Gender' failed on the 'oneof' tag"
}
```

#### User Not Found (Hiếm gặp)
```json
{
  "success": false,
  "error": "không tìm thấy người dùng"
}
```

### Flutter Implementation Example
```dart
class UserProfileService {
  final dio = Dio();
  final storage = GetStorage();

  Future<UserProfile> updateProfile({
    String? name,
    String? phone,
    String? avatar,
    String? address,
    String? gender,
  }) async {
    try {
      final token = storage.read('jwt_token');
      
      final payload = <String, dynamic>{};
      
      if (name != null) payload['name'] = name;
      if (phone != null) payload['phone'] = phone;
      if (avatar != null) payload['avatar'] = avatar;
      if (address != null) payload['address'] = address;
      if (gender != null) payload['gender'] = gender;

      final response = await dio.patch(
        'http://localhost:8080/api/v1/users/profile',
        data: payload,
        options: Options(
          headers: {'Authorization': 'Bearer $token'},
        ),
      );

      if (response.statusCode == 200) {
        return UserProfile.fromJson(response.data['data']);
      }
    } catch (e) {
      throw Exception('Không thể cập nhật profile: $e');
    }
  }
}

// Controller example
class EditProfileController extends GetxController {
  final userService = UserProfileService();
  
  final nameController = TextEditingController();
  final phoneController = TextEditingController();
  final addressController = TextEditingController();
  final genderController = TextEditingController();
  
  var selectedGender = Rx<String?>(null);
  var isLoading = false.obs;

  @override
  void onInit() {
    super.onInit();
    loadProfile();
  }

  void loadProfile() async {
    try {
      final profile = await userService.getProfile();
      nameController.text = profile.name;
      phoneController.text = profile.phone ?? '';
      addressController.text = profile.address ?? '';
      selectedGender.value = profile.gender;
    } catch (e) {
      Get.snackbar('Lỗi', 'Không thể tải profile');
    }
  }

  void saveProfile() async {
    try {
      isLoading.value = true;
      
      await userService.updateProfile(
        name: nameController.text.isNotEmpty ? nameController.text : null,
        phone: phoneController.text.isNotEmpty ? phoneController.text : null,
        address: addressController.text.isNotEmpty ? addressController.text : null,
        gender: selectedGender.value,
      );

      Get.snackbar('Thành công', 'Cập nhật profile thành công');
    } catch (e) {
      Get.snackbar('Lỗi', e.toString());
    } finally {
      isLoading.value = false;
    }
  }
}
```

---

## 3. Đổi Mật Khẩu

### Endpoint
```
PATCH /api/v1/users/change-password
```

### Mô Tả
Đổi mật khẩu của người dùng hiện tại. Yêu cầu mật khẩu cũ và mật khẩu mới.

### Headers
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

### Request Body
```json
{
  "oldPassword": "oldPassword123",
  "newPassword": "newPassword123"
}
```

### Validation Rules
| Field | Validation |
|-------|-----------|
| `oldPassword` | Bắt buộc, tối thiểu 6 ký tự |
| `newPassword` | Bắt buộc, tối thiểu 6 ký tự, phải khác mật khẩu cũ |

### Response (200 OK)
```json
{
  "success": true,
  "message": "Đổi mật khẩu thành công"
}
```

### Error Cases

#### Mật khẩu cũ sai
```json
{
  "success": false,
  "error": "mật khẩu cũ không chính xác"
}
```

#### Mật khẩu cũ = mật khẩu mới
```json
{
  "success": false,
  "error": "mật khẩu mới phải khác mật khẩu cũ"
}
```

#### Invalid Input
```json
{
  "success": false,
  "error": "Dữ liệu không hợp lệ",
  "details": "Key: 'ChangePasswordRequest.NewPassword' Error:Field validation for 'NewPassword' failed on the 'min' tag"
}
```

### Flutter Implementation Example
```dart
class PasswordService {
  final dio = Dio();
  final storage = GetStorage();

  Future<bool> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    try {
      final token = storage.read('jwt_token');
      
      final response = await dio.patch(
        'http://localhost:8080/api/v1/users/change-password',
        data: {
          'oldPassword': oldPassword,
          'newPassword': newPassword,
        },
        options: Options(
          headers: {'Authorization': 'Bearer $token'},
        ),
      );

      if (response.statusCode == 200) {
        return response.data['success'] ?? false;
      }
    } catch (e) {
      throw Exception('Không thể đổi mật khẩu: $e');
    }
  }
}

// Controller example
class ChangePasswordController extends GetxController {
  final passwordService = PasswordService();
  
  final oldPasswordController = TextEditingController();
  final newPasswordController = TextEditingController();
  final confirmPasswordController = TextEditingController();
  
  var isLoading = false.obs;
  var oldPasswordVisible = false.obs;
  var newPasswordVisible = false.obs;
  var confirmPasswordVisible = false.obs;

  void changePassword() async {
    // Validate
    if (oldPasswordController.text.isEmpty) {
      Get.snackbar('Lỗi', 'Vui lòng nhập mật khẩu cũ');
      return;
    }

    if (newPasswordController.text.isEmpty) {
      Get.snackbar('Lỗi', 'Vui lòng nhập mật khẩu mới');
      return;
    }

    if (newPasswordController.text.length < 6) {
      Get.snackbar('Lỗi', 'Mật khẩu mới phải có tối thiểu 6 ký tự');
      return;
    }

    if (newPasswordController.text != confirmPasswordController.text) {
      Get.snackbar('Lỗi', 'Xác nhận mật khẩu không khớp');
      return;
    }

    if (oldPasswordController.text == newPasswordController.text) {
      Get.snackbar('Lỗi', 'Mật khẩu mới phải khác mật khẩu cũ');
      return;
    }

    try {
      isLoading.value = true;
      
      await passwordService.changePassword(
        oldPassword: oldPasswordController.text,
        newPassword: newPasswordController.text,
      );

      Get.snackbar('Thành công', 'Đổi mật khẩu thành công');
      
      // Clear fields
      oldPasswordController.clear();
      newPasswordController.clear();
      confirmPasswordController.clear();
      
      // Optionally navigate to login or profile screen
      Get.back();
    } catch (e) {
      Get.snackbar('Lỗi', e.toString());
    } finally {
      isLoading.value = false;
    }
  }
}
```

---

## 4. Upload Avatar

### Endpoint
```
POST /api/v1/users/upload-avatar
```

### Mô Tả
Upload ảnh đại diện mới cho người dùng hiện tại. Backend sẽ tự động xử lý upload lên Cloudinary và trả về URL.

### Headers
```
Authorization: Bearer {jwt_token}
Content-Type: multipart/form-data
```

### Request Parameters
| Parameter | Kiểu | Bắt buộc | Mô Tả |
|-----------|------|---------|-------|
| `avatar` | File | ✅ Có | Tệp ảnh (JPG, PNG, GIF, WebP) |

### File Constraints
- **Kích thước tối đa**: Phụ thuộc cấu hình Cloudinary (thường ~5-10MB)
- **Định dạng**: JPG, PNG, GIF, WebP
- **Tỷ lệ**: Không hạn chế

### Response (200 OK)
```json
{
  "success": true,
  "message": "Upload ảnh đại diện thành công",
  "data": {
    "avatarUrl": "https://res.cloudinary.com/example/image/upload/v123456/avatar.jpg"
  }
}
```

### Error Cases

#### No file provided
```json
{
  "success": false,
  "error": "Vui lòng chọn tệp ảnh"
}
```

#### Invalid file type
```json
{
  "success": false,
  "error": "Định dạng tệp không hợp lệ"
}
```

#### Upload failed
```json
{
  "success": false,
  "error": "Không thể upload ảnh"
}
```

### Flutter Implementation Example
```dart
class AvatarService {
  final dio = Dio();
  final storage = GetStorage();

  Future<String> uploadAvatar(File imageFile) async {
    try {
      final token = storage.read('jwt_token');
      
      final formData = FormData.fromMap({
        'avatar': await MultipartFile.fromFile(
          imageFile.path,
          filename: 'avatar_${DateTime.now().millisecondsSinceEpoch}.jpg',
        ),
      });

      final response = await dio.post(
        'http://localhost:8080/api/v1/users/upload-avatar',
        data: formData,
        options: Options(
          headers: {'Authorization': 'Bearer $token'},
        ),
      );

      if (response.statusCode == 200) {
        return response.data['data']['avatarUrl'] as String;
      }
    } catch (e) {
      throw Exception('Không thể upload avatar: $e');
    }
  }
}

// Controller example
class AvatarController extends GetxController {
  final avatarService = AvatarService();
  final userProfileService = UserProfileService();
  
  var avatarUrl = Rx<String?>(null);
  var isUploading = false.obs;

  void pickAndUploadAvatar() async {
    try {
      final picker = ImagePicker();
      final image = await picker.pickImage(source: ImageSource.gallery);

      if (image == null) return;

      isUploading.value = true;

      final file = File(image.path);
      
      // Optionally compress image
      final compressedFile = await _compressImage(file);
      
      final url = await avatarService.uploadAvatar(compressedFile);
      avatarUrl.value = url;

      Get.snackbar('Thành công', 'Upload avatar thành công');
    } catch (e) {
      Get.snackbar('Lỗi', e.toString());
    } finally {
      isUploading.value = false;
    }
  }

  Future<File> _compressImage(File file) async {
    final result = await FlutterImageCompress.compressAndGetFile(
      file.absolute.path,
      "${file.parent.path}/compressed_${DateTime.now().millisecondsSinceEpoch}.jpg",
      quality: 80,
    );
    return File(result?.path ?? file.path);
  }
}
```

---

## 5. Xóa Avatar

### Endpoint
```
DELETE /api/v1/users/delete-avatar
```

### Mô Tả
Xóa ảnh đại diện hiện tại của người dùng.

### Headers
```
Authorization: Bearer {jwt_token}
```

### Response (200 OK)
```json
{
  "success": true,
  "message": "Xóa ảnh đại diện thành công"
}
```

### Error Cases

#### No avatar to delete
```json
{
  "success": false,
  "error": "Người dùng không có ảnh đại diện"
}
```

#### Delete failed
```json
{
  "success": false,
  "error": "Không thể xóa ảnh đại diện"
}
```

### Flutter Implementation Example
```dart
class AvatarService {
  final dio = Dio();
  final storage = GetStorage();

  Future<bool> deleteAvatar() async {
    try {
      final token = storage.read('jwt_token');
      
      final response = await dio.delete(
        'http://localhost:8080/api/v1/users/delete-avatar',
        options: Options(
          headers: {'Authorization': 'Bearer $token'},
        ),
      );

      if (response.statusCode == 200) {
        return response.data['success'] ?? false;
      }
    } catch (e) {
      throw Exception('Không thể xóa avatar: $e');
    }
  }
}

// In controller
void deleteAvatar() async {
  try {
    final confirm = await Get.dialog<bool>(
      AlertDialog(
        title: const Text('Xóa ảnh đại diện'),
        content: const Text('Bạn có chắc muốn xóa ảnh đại diện?'),
        actions: [
          TextButton(
            onPressed: () => Get.back(result: false),
            child: const Text('Hủy'),
          ),
          TextButton(
            onPressed: () => Get.back(result: true),
            child: const Text('Xóa'),
          ),
        ],
      ),
    ) ?? false;

    if (!confirm) return;

    isUploading.value = true;
    await avatarService.deleteAvatar();
    avatarUrl.value = null;

    Get.snackbar('Thành công', 'Xóa avatar thành công');
  } catch (e) {
    Get.snackbar('Lỗi', e.toString());
  } finally {
    isUploading.value = false;
  }
}
```

---

## Flow Diagram

```
┌─────────────────┐
│  Open Profile   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│ GET /users/profile      │ ◄─── Lấy thông tin hiện tại
│ (AuthMiddleware)        │
└────────┬────────────────┘
         │
         ├─────────────────────────────────────────┐
         │                                         │
         ▼                                         ▼
   ┌──────────────┐                    ┌──────────────────┐
   │ Edit Info    │                    │ Change Avatar    │
   │ (Text Fields)│                    │ (Image Picker)   │
   └──────┬───────┘                    └────────┬─────────┘
          │                                     │
          ▼                                     ▼
   ┌─────────────────────┐             ┌──────────────────────┐
   │PATCH /users/profile │             │POST /users/upload-   │
   │(Send delta changes) │             │avatar (Multipart)    │
   └──────┬──────────────┘             └────────┬─────────────┘
          │                                     │
          └─────────────────┬───────────────────┘
                            │
                            ▼
                    ┌──────────────────┐
                    │ Success Response │
                    │ Reload Profile   │
                    └──────────────────┘
```

---

## Best Practices

### 1. **Authentication**
```dart
// ✅ Luôn kiểm tra token trước khi gửi request
final token = storage.read('jwt_token');
if (token == null || token.isEmpty) {
  // Redirect to login
  Get.offAllNamed('/login');
  return;
}
```

### 2. **Error Handling**
```dart
// ✅ Xử lý các loại lỗi khác nhau
try {
  await updateProfile(...);
} on DioException catch (e) {
  if (e.response?.statusCode == 401) {
    // Token expired, redirect to login
    Get.offAllNamed('/login');
  } else if (e.response?.statusCode == 400) {
    // Validation error
    final errorMsg = e.response?.data['error'] ?? 'Lỗi không xác định';
    Get.snackbar('Lỗi', errorMsg);
  }
} catch (e) {
  // Handle unexpected errors
  Get.snackbar('Lỗi', 'Đã xảy ra lỗi: $e');
}
```

### 3. **Validation Trước Khi Gửi**
```dart
// ✅ Validate local trước, rồi gửi server
bool validateProfile() {
  if (nameController.text.isEmpty) {
    Get.snackbar('Lỗi', 'Tên không thể trống');
    return false;
  }

  if (selectedGender.value != null && 
      !['male', 'female', 'other'].contains(selectedGender.value)) {
    Get.snackbar('Lỗi', 'Giới tính không hợp lệ');
    return false;
  }

  return true;
}
```

### 4. **Loading State**
```dart
// ✅ Hiển thị loading indicator khi đang xử lý
Obx(() => isLoading.value
  ? Center(child: CircularProgressIndicator())
  : ElevatedButton(
      onPressed: saveProfile,
      child: Text('Lưu'),
    )
)
```

### 5. **Image Optimization**
```dart
// ✅ Nén ảnh trước khi upload để tiết kiệm bandwidth
Future<File> _compressImage(File file) async {
  final result = await FlutterImageCompress.compressAndGetFile(
    file.absolute.path,
    "${file.parent.path}/compressed_${DateTime.now().millisecondsSinceEpoch}.jpg",
    quality: 75,
  );
  return File(result?.path ?? file.path);
}
```

### 6. **Cache Management**
```dart
// ✅ Cache profile data để giảm API calls
class UserProfileService {
  UserProfile? _cachedProfile;
  
  Future<UserProfile> getProfile({bool forceRefresh = false}) async {
    if (!forceRefresh && _cachedProfile != null) {
      return _cachedProfile!;
    }
    
    final response = await dio.get('/users/profile', ...);
    _cachedProfile = UserProfile.fromJson(response.data['data']);
    return _cachedProfile!;
  }
  
  void invalidateCache() {
    _cachedProfile = null;
  }
}
```

### 7. **Timeout Handling**
```dart
// ✅ Set timeout phù hợp
final response = await dio.patch(
  url,
  data: payload,
  options: Options(
    headers: {'Authorization': 'Bearer $token'},
    connectTimeout: Duration(seconds: 10),
    receiveTimeout: Duration(seconds: 10),
  ),
);
```

---

## Common Issues & Solutions

### Issue 1: "Unauthorized" Error (401)
**Nguyên nhân**: Token đã hết hạn hoặc không hợp lệ
```dart
// Solution: Refresh token hoặc redirect to login
if (e.response?.statusCode == 401) {
  // Try to refresh token
  final newToken = await authService.refreshToken();
  if (newToken != null) {
    // Retry the request
  } else {
    // Token invalid, logout
    await logout();
    Get.offAllNamed('/login');
  }
}
```

### Issue 2: Timeout During Large File Upload
**Nguyên nhân**: File quá lớn hoặc network chậm
```dart
// Solution: Increase timeout, compress image, show progress
final result = await FlutterImageCompress.compressAndGetFile(
  file.path,
  compressedPath,
  quality: 60, // Reduce quality
);

final response = await dio.post(
  url,
  data: formData,
  options: Options(connectTimeout: Duration(seconds: 30)),
  onSendProgress: (sent, total) {
    uploadProgress.value = (sent / total * 100).toInt();
  },
);
```

### Issue 3: "Invalid Gender Value" Validation Error
**Nguyên nhân**: Gửi giá trị không trong list: `male`, `female`, `other`
```dart
// Solution: Validate before sending
if (selectedGender.value != null && 
    !['male', 'female', 'other'].contains(selectedGender.value)) {
  Get.snackbar('Lỗi', 'Giới tính không hợp lệ');
  return;
}
```

### Issue 4: Old Password Incorrect
**Nguyên nhân**: User nhập sai mật khẩu cũ
```dart
// Solution: Inform user clearly
try {
  await passwordService.changePassword(...);
} on DioException catch (e) {
  if (e.response?.data['error'].contains('mật khẩu cũ')) {
    Get.snackbar('Lỗi', 'Mật khẩu cũ không chính xác');
  }
}
```

---

## Request/Response Examples

### Complete Profile Update Flow

#### Step 1: Fetch Current Profile
```bash
curl -X GET "http://localhost:8080/api/v1/users/profile" \
  -H "Authorization: Bearer eyJhbGc..."
```

#### Step 2: Update Profile
```bash
curl -X PATCH "http://localhost:8080/api/v1/users/profile" \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nguyễn Văn A Mới",
    "phone": "0987654321",
    "gender": "female",
    "address": "789 Đường DEF, Quận 3, TP.HCM"
  }'
```

Response:
```json
{
  "success": true,
  "message": "Cập nhật thông tin cá nhân thành công",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "Nguyễn Văn A Mới",
    "role": "customer",
    "phone": "0987654321",
    "avatar": "https://example.com/avatar.jpg",
    "address": "789 Đường DEF, Quận 3, TP.HCM",
    "gender": "female",
    "isEmailVerified": true,
    "isActive": true,
    "isFirstLogin": false,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-12-29T16:30:00Z"
  }
}
```

---

## Summary Table

| Action | Method | Endpoint | Auth | Yêu Cầu |
|--------|--------|----------|------|---------|
| Xem Profile | GET | `/users/profile` | ✅ | JWT Token |
| Cập Nhật Info | PATCH | `/users/profile` | ✅ | JWT Token, Valid fields |
| Đổi Mật Khẩu | PATCH | `/users/change-password` | ✅ | JWT Token, Old & New Password |
| Upload Avatar | POST | `/users/upload-avatar` | ✅ | JWT Token, Image File |
| Xóa Avatar | DELETE | `/users/delete-avatar` | ✅ | JWT Token |

---

## Kết Luận

- **Người dùng có thể**: Chỉnh sửa thông tin cá nhân của chính mình (tên, số điện thoại, địa chỉ, giới tính)
- **Người dùng không thể**: Thay đổi email, role, hoặc các trường được bảo vệ
- **Luôn bắt buộc**: Token JWT trong header Authorization
- **Validation**: Xử lý cả client-side và server-side

Với guide này, team Flutter có thể xây dựng feature chỉnh sửa profile hoàn chỉnh và an toàn! 🎉

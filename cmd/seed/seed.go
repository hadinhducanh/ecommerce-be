package main

import (
	"log"

	"ecommerce-be/config"
	"ecommerce-be/database"
	"ecommerce-be/models"
	"ecommerce-be/utils"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// Seed Admin User
	log.Println("👤 Starting to seed admin user...")
	seedAdminUser()

	log.Println("🌱 Starting to seed categories...")

	// Bước 1: Tạo hoặc lấy 2 parent categories (chỉ có 2 parent: Điện thoại và Laptop)
	log.Println("📦 Creating parent categories (Điện thoại, Laptop)...")
	phoneParent := createOrUpdateCategory(models.Category{
		Name:          "Điện thoại",
		NameEn:        stringPtr("Smartphones"),
		Description:   stringPtr("Các dòng điện thoại thông minh với công nghệ hiện đại, hiệu năng mạnh mẽ"),
		DescriptionEn: stringPtr("Modern smartphones with advanced technology and powerful performance"),
		IsActive:      true,
	})

	laptopParent := createOrUpdateCategory(models.Category{
		Name:          "Laptop",
		NameEn:        stringPtr("Laptops"),
		Description:   stringPtr("Các dòng laptop đa dạng phục vụ nhu cầu học tập, làm việc và giải trí"),
		DescriptionEn: stringPtr("Various laptop models for study, work, and entertainment needs"),
		IsActive:      true,
	})

	// Bước 2: Xóa tất cả quan hệ parent-child cũ để đảm bảo cấu trúc sạch
	log.Println("🧹 Cleaning up old parent-child relationships...")
	if err := database.DB.Exec("DELETE FROM category_children").Error; err != nil {
		log.Printf("⚠️  Warning: Failed to clean up old relationships: %v", err)
	} else {
		log.Println("✅ Cleaned up old relationships")
	}

	// Danh sách child categories của Điện thoại
	phoneCategories := []models.Category{
		{
			Name:          "iPhone",
			NameEn:        stringPtr("iPhone"),
			Description:   stringPtr("Điện thoại iPhone của Apple với hệ điều hành iOS"),
			DescriptionEn: stringPtr("Apple iPhone smartphones with iOS operating system"),
			IsActive:      true,
		},
		{
			Name:          "Samsung",
			NameEn:        stringPtr("Samsung"),
			Description:   stringPtr("Điện thoại Samsung Galaxy với màn hình AMOLED và camera chất lượng cao"),
			DescriptionEn: stringPtr("Samsung Galaxy smartphones with AMOLED displays and high-quality cameras"),
			IsActive:      true,
		},
		{
			Name:          "Xiaomi",
			NameEn:        stringPtr("Xiaomi"),
			Description:   stringPtr("Điện thoại Xiaomi với giá cả hợp lý và cấu hình mạnh"),
			DescriptionEn: stringPtr("Xiaomi smartphones with affordable prices and powerful specifications"),
			IsActive:      true,
		},
		{
			Name:          "OPPO",
			NameEn:        stringPtr("OPPO"),
			Description:   stringPtr("Điện thoại OPPO với camera selfie xuất sắc và thiết kế đẹp mắt"),
			DescriptionEn: stringPtr("OPPO smartphones with excellent selfie cameras and beautiful design"),
			IsActive:      true,
		},
		{
			Name:          "Vivo",
			NameEn:        stringPtr("Vivo"),
			Description:   stringPtr("Điện thoại Vivo với công nghệ camera AI và pin dung lượng lớn"),
			DescriptionEn: stringPtr("Vivo smartphones with AI camera technology and large battery capacity"),
			IsActive:      true,
		},
		{
			Name:          "Realme",
			NameEn:        stringPtr("Realme"),
			Description:   stringPtr("Điện thoại Realme với hiệu năng tốt và giá cả phải chăng"),
			DescriptionEn: stringPtr("Realme smartphones with good performance and affordable prices"),
			IsActive:      true,
		},
		{
			Name:          "OnePlus",
			NameEn:        stringPtr("OnePlus"),
			Description:   stringPtr("Điện thoại OnePlus với hiệu năng flagship và sạc nhanh"),
			DescriptionEn: stringPtr("OnePlus smartphones with flagship performance and fast charging"),
			IsActive:      true,
		},
	}

	// Danh sách child categories của Laptop
	laptopCategories := []models.Category{
		{
			Name:          "MacBook",
			NameEn:        stringPtr("MacBook"),
			Description:   stringPtr("Laptop MacBook của Apple với chip M-series, hiệu năng vượt trội"),
			DescriptionEn: stringPtr("Apple MacBook laptops with M-series chips and outstanding performance"),
			IsActive:      true,
		},
		{
			Name:          "Laptop Dell",
			NameEn:        stringPtr("Dell Laptops"),
			Description:   stringPtr("Laptop Dell với độ bền cao, phù hợp cho doanh nghiệp và cá nhân"),
			DescriptionEn: stringPtr("Dell laptops with high durability, suitable for businesses and individuals"),
			IsActive:      true,
		},
		{
			Name:          "Laptop HP",
			NameEn:        stringPtr("HP Laptops"),
			Description:   stringPtr("Laptop HP với thiết kế đẹp, hiệu năng ổn định"),
			DescriptionEn: stringPtr("HP laptops with beautiful design and stable performance"),
			IsActive:      true,
		},
		{
			Name:          "Laptop Lenovo",
			NameEn:        stringPtr("Lenovo Laptops"),
			Description:   stringPtr("Laptop Lenovo ThinkPad và IdeaPad với bàn phím tốt, bền bỉ"),
			DescriptionEn: stringPtr("Lenovo ThinkPad and IdeaPad laptops with good keyboards and durability"),
			IsActive:      true,
		},
		{
			Name:          "Laptop Asus",
			NameEn:        stringPtr("Asus Laptops"),
			Description:   stringPtr("Laptop Asus với card đồ họa mạnh, phù hợp gaming và đồ họa"),
			DescriptionEn: stringPtr("Asus laptops with powerful graphics cards, suitable for gaming and graphics"),
			IsActive:      true,
		},
		{
			Name:          "Laptop Acer",
			NameEn:        stringPtr("Acer Laptops"),
			Description:   stringPtr("Laptop Acer với giá cả hợp lý, cấu hình đa dạng"),
			DescriptionEn: stringPtr("Acer laptops with affordable prices and diverse configurations"),
			IsActive:      true,
		},
		{
			Name:          "Laptop MSI",
			NameEn:        stringPtr("MSI Laptops"),
			Description:   stringPtr("Laptop MSI chuyên gaming với hiệu năng cao và tản nhiệt tốt"),
			DescriptionEn: stringPtr("MSI gaming laptops with high performance and good cooling"),
			IsActive:      true,
		},
		{
			Name:          "Laptop Razer",
			NameEn:        stringPtr("Razer Laptops"),
			Description:   stringPtr("Laptop Razer Blade với thiết kế premium, hiệu năng gaming mạnh"),
			DescriptionEn: stringPtr("Razer Blade laptops with premium design and strong gaming performance"),
			IsActive:      true,
		},
	}

	// Bước 3: Tạo tất cả child categories (tất cả đều là root categories khi tạo, sau đó sẽ được gán parent)
	log.Println("📱 Creating phone child categories...")
	var phoneChildCategories []models.Category
	for _, category := range phoneCategories {
		created := createOrUpdateCategory(category)
		phoneChildCategories = append(phoneChildCategories, created)
		log.Printf("  ✓ Created/Updated: %s (ID: %d)", created.Name, created.ID)
	}

	log.Println("💻 Creating laptop child categories...")
	var laptopChildCategories []models.Category
	for _, category := range laptopCategories {
		created := createOrUpdateCategory(category)
		laptopChildCategories = append(laptopChildCategories, created)
		log.Printf("  ✓ Created/Updated: %s (ID: %d)", created.Name, created.ID)
	}

	// Bước 4: Tạo quan hệ parent-child trong bảng category_children
	log.Println("🔗 Creating parent-child relationships...")

	// Thêm children cho Điện thoại
	log.Printf("  📱 Adding %d children to 'Điện thoại' (ID: %d)...", len(phoneChildCategories), phoneParent.ID)
	for _, child := range phoneChildCategories {
		createOrUpdateCategoryChild(phoneParent.ID, child.ID)
	}

	// Thêm children cho Laptop
	log.Printf("  💻 Adding %d children to 'Laptop' (ID: %d)...", len(laptopChildCategories), laptopParent.ID)
	for _, child := range laptopChildCategories {
		createOrUpdateCategoryChild(laptopParent.ID, child.ID)
	}

	log.Println("✅ Seeding categories completed!")
	log.Printf("📊 Summary:")
	log.Printf("  - Parent categories: 2 (Điện thoại, Laptop)")
	log.Printf("  - Phone children: %d", len(phoneChildCategories))
	log.Printf("  - Laptop children: %d", len(laptopChildCategories))
	log.Printf("  - Total categories: %d", 2+len(phoneChildCategories)+len(laptopChildCategories))
}

// createOrUpdateCategory tạo hoặc cập nhật category (mặc định là root category)
func createOrUpdateCategory(category models.Category) models.Category {
	var existingCategory models.Category
	result := database.DB.Where("name = ?", category.Name).First(&existingCategory)

	if result.Error != nil {
		// Category chưa tồn tại → tạo mới
		if err := database.DB.Create(&category).Error; err != nil {
			log.Printf("❌ Failed to create category %s: %v", category.Name, err)
			return category
		}
		log.Printf("✅ Created category: %s", category.Name)
		// Lấy lại category vừa tạo để có ID
		database.DB.Where("name = ?", category.Name).First(&category)
		return category
	} else {
		// Category đã tồn tại → cập nhật
		existingCategory.NameEn = category.NameEn
		existingCategory.Description = category.Description
		existingCategory.DescriptionEn = category.DescriptionEn
		existingCategory.IsActive = category.IsActive

		if err := database.DB.Save(&existingCategory).Error; err != nil {
			log.Printf("❌ Failed to update category %s: %v", category.Name, err)
			return existingCategory
		}
		log.Printf("🔄 Updated category: %s", category.Name)
		return existingCategory
	}
}

// createOrUpdateCategoryChild tạo hoặc cập nhật quan hệ parent-child
func createOrUpdateCategoryChild(parentID, childID uint) {
	var existingRelation models.CategoryChild
	result := database.DB.Where("parent_id = ? AND child_id = ?", parentID, childID).First(&existingRelation)

	if result.Error != nil {
		// Quan hệ chưa tồn tại → tạo mới
		relation := models.CategoryChild{
			ParentID: parentID,
			ChildID:  childID,
		}
		if err := database.DB.Create(&relation).Error; err != nil {
			log.Printf("❌ Failed to create parent-child relation %d -> %d: %v", parentID, childID, err)
		} else {
			log.Printf("✅ Created parent-child relation: %d -> %d", parentID, childID)
		}
	} else {
		// Quan hệ đã tồn tại
		log.Printf("ℹ️  Parent-child relation already exists: %d -> %d", parentID, childID)
	}
}

// Helper function để tạo string pointer
func stringPtr(s string) *string {
	return &s
}

// seedAdminUser tạo tài khoản admin mặc định
func seedAdminUser() {
	adminEmail := "admin@ecommerce.com"
	
	// Kiểm tra xem admin đã tồn tại chưa
	var existingAdmin models.User
	result := database.DB.Where("email = ?", adminEmail).First(&existingAdmin)
	
	if result.Error == nil {
		log.Printf("ℹ️  Admin user already exists: %s (ID: %d)", adminEmail, existingAdmin.ID)
		return
	}
	
	// Hash mật khẩu
	hashedPassword, err := utils.HashPassword("1")
	if err != nil {
		log.Printf("❌ Failed to hash password: %v", err)
		return
	}
	
	// Tạo admin user
	admin := models.User{
		Email:           adminEmail,
		Password:        hashedPassword,
		Name:            "Administrator",
		Role:            "admin",
		IsActive:        true,
		IsEmailVerified: true,
		IsFirstLogin:    false,
		Phone:           stringPtr("0123456789"),
		Gender:          stringPtr("other"),
	}
	
	if err := database.DB.Create(&admin).Error; err != nil {
		log.Printf("❌ Failed to create admin user: %v", err)
		return
	}
	
	log.Printf("✅ Admin user created successfully!")
	log.Printf("   📧 Email: %s", adminEmail)
	log.Printf("   🔑 Password: 1")
	log.Printf("   👤 Role: admin")
	log.Printf("   🆔 ID: %d", admin.ID)
}

package main

import (
	"errors"
	"fmt"
	"log"

	"ecommerce-be/config"
	"ecommerce-be/database"
	"ecommerce-be/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
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

	log.Println("🌱 Starting to seed products...")

	// Bước 1: Lấy tất cả child categories
	log.Println("📦 Fetching child categories...")
	var childCategories []models.Category
	var relations []models.CategoryChild
	
	// Lấy tất cả quan hệ parent-child
	if err := database.DB.Find(&relations).Error; err != nil {
		log.Fatalf("❌ Failed to fetch category relations: %v", err)
	}

	if len(relations) == 0 {
		log.Println("⚠️  No child categories found. Please run seed.go first to create categories.")
		return
	}

	// Lấy danh sách ID các children
	childIDs := make([]uint, len(relations))
	for i, rel := range relations {
		childIDs[i] = rel.ChildID
	}

	// Lấy thông tin các child categories
	if err := database.DB.Where("id IN ?", childIDs).Find(&childCategories).Error; err != nil {
		log.Fatalf("❌ Failed to fetch child categories: %v", err)
	}

	log.Printf("✅ Found %d child categories", len(childCategories))

	// Bước 2: Tạo products cho mỗi child category (4 products mỗi category)
	totalProducts := 0
	for _, category := range childCategories {
		log.Printf("\n📱 Creating products for category: %s (ID: %d)", category.Name, category.ID)
		
		products := generateProductsForCategory(category)
		
		for i, product := range products {
			created := createOrUpdateProduct(product)
			totalProducts++
			log.Printf("  ✓ [%d/4] Created/Updated: %s (ID: %d, Price: %.0f VNĐ)", i+1, created.Name, created.ID, created.Price)
		}
	}

	log.Println("\n✅ Seeding products completed!")
	log.Printf("📊 Summary:")
	log.Printf("  - Child categories: %d", len(childCategories))
	log.Printf("  - Products per category: 4")
	log.Printf("  - Total products created: %d", totalProducts)
}

// generateProductsForCategory tạo danh sách 4 products cho một category
func generateProductsForCategory(category models.Category) []models.Product {
	categoryName := category.Name
	products := []models.Product{}

	// Tạo products dựa trên tên category
	switch categoryName {
	case "iPhone":
		products = []models.Product{
			{
				Name:          "iPhone 15 Pro Max 256GB",
				NameEn:        stringPtrProduct("iPhone 15 Pro Max 256GB"),
				Description:   stringPtrProduct("iPhone 15 Pro Max với chip A17 Pro, camera 48MP, màn hình Super Retina XDR 6.7 inch"),
				DescriptionEn: stringPtrProduct("iPhone 15 Pro Max with A17 Pro chip, 48MP camera, 6.7 inch Super Retina XDR display"),
				Price:         32990000,
				Stock:         50,
				SKU:           stringPtrProduct("IPH15PM256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/iphone15pm1.jpg", "https://example.com/iphone15pm2.jpg"},
			},
			{
				Name:          "iPhone 15 Pro 128GB",
				NameEn:        stringPtrProduct("iPhone 15 Pro 128GB"),
				Description:   stringPtrProduct("iPhone 15 Pro với chip A17 Pro, camera 48MP, màn hình Super Retina XDR 6.1 inch"),
				DescriptionEn: stringPtrProduct("iPhone 15 Pro with A17 Pro chip, 48MP camera, 6.1 inch Super Retina XDR display"),
				Price:         26990000,
				Stock:         75,
				SKU:           stringPtrProduct("IPH15P128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/iphone15p1.jpg", "https://example.com/iphone15p2.jpg"},
			},
			{
				Name:          "iPhone 14 128GB",
				NameEn:        stringPtrProduct("iPhone 14 128GB"),
				Description:   stringPtrProduct("iPhone 14 với chip A15 Bionic, camera kép 12MP, màn hình Super Retina XDR 6.1 inch"),
				DescriptionEn: stringPtrProduct("iPhone 14 with A15 Bionic chip, dual 12MP camera, 6.1 inch Super Retina XDR display"),
				Price:         19990000,
				Stock:         100,
				SKU:           stringPtrProduct("IPH14128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/iphone141.jpg"},
			},
			{
				Name:          "iPhone 13 128GB",
				NameEn:        stringPtrProduct("iPhone 13 128GB"),
				Description:   stringPtrProduct("iPhone 13 với chip A15 Bionic, camera kép 12MP, màn hình Super Retina XDR 6.1 inch"),
				DescriptionEn: stringPtrProduct("iPhone 13 with A15 Bionic chip, dual 12MP camera, 6.1 inch Super Retina XDR display"),
				Price:         15990000,
				Stock:         80,
				SKU:           stringPtrProduct("IPH13128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/iphone131.jpg"},
			},
		}

	case "Samsung":
		products = []models.Product{
			{
				Name:          "Samsung Galaxy S24 Ultra 256GB",
				NameEn:        stringPtrProduct("Samsung Galaxy S24 Ultra 256GB"),
				Description:   stringPtrProduct("Galaxy S24 Ultra với chip Snapdragon 8 Gen 3, camera 200MP, bút S Pen, màn hình Dynamic AMOLED 2X 6.8 inch"),
				DescriptionEn: stringPtrProduct("Galaxy S24 Ultra with Snapdragon 8 Gen 3 chip, 200MP camera, S Pen, 6.8 inch Dynamic AMOLED 2X display"),
				Price:         28990000,
				Stock:         60,
				SKU:           stringPtrProduct("SGS24U256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/s24u1.jpg"},
			},
			{
				Name:          "Samsung Galaxy S23 128GB",
				NameEn:        stringPtrProduct("Samsung Galaxy S23 128GB"),
				Description:   stringPtrProduct("Galaxy S23 với chip Snapdragon 8 Gen 2, camera 50MP, màn hình Dynamic AMOLED 2X 6.1 inch"),
				DescriptionEn: stringPtrProduct("Galaxy S23 with Snapdragon 8 Gen 2 chip, 50MP camera, 6.1 inch Dynamic AMOLED 2X display"),
				Price:         17990000,
				Stock:         90,
				SKU:           stringPtrProduct("SGS23128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/s231.jpg"},
			},
			{
				Name:          "Samsung Galaxy A54 128GB",
				NameEn:        stringPtrProduct("Samsung Galaxy A54 128GB"),
				Description:   stringPtrProduct("Galaxy A54 với chip Exynos 1380, camera 50MP, màn hình Super AMOLED 6.4 inch"),
				DescriptionEn: stringPtrProduct("Galaxy A54 with Exynos 1380 chip, 50MP camera, 6.4 inch Super AMOLED display"),
				Price:         8990000,
				Stock:         120,
				SKU:           stringPtrProduct("SGA54128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/a541.jpg"},
			},
			{
				Name:          "Samsung Galaxy Z Fold5 256GB",
				NameEn:        stringPtrProduct("Samsung Galaxy Z Fold5 256GB"),
				Description:   stringPtrProduct("Galaxy Z Fold5 màn hình gập với chip Snapdragon 8 Gen 2, camera 50MP, màn hình chính 7.6 inch"),
				DescriptionEn: stringPtrProduct("Galaxy Z Fold5 foldable with Snapdragon 8 Gen 2 chip, 50MP camera, 7.6 inch main display"),
				Price:         39990000,
				Stock:         30,
				SKU:           stringPtrProduct("SGZF5256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/zfold51.jpg"},
			},
		}

	case "Xiaomi":
		products = []models.Product{
			{
				Name:          "Xiaomi 14 Pro 256GB",
				NameEn:        stringPtrProduct("Xiaomi 14 Pro 256GB"),
				Description:   stringPtrProduct("Xiaomi 14 Pro với chip Snapdragon 8 Gen 3, camera Leica 50MP, màn hình AMOLED 6.73 inch"),
				DescriptionEn: stringPtrProduct("Xiaomi 14 Pro with Snapdragon 8 Gen 3 chip, Leica 50MP camera, 6.73 inch AMOLED display"),
				Price:         19990000,
				Stock:         70,
				SKU:           stringPtrProduct("XM14P256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/xm14p1.jpg"},
			},
			{
				Name:          "Xiaomi 13T 256GB",
				NameEn:        stringPtrProduct("Xiaomi 13T 256GB"),
				Description:   stringPtrProduct("Xiaomi 13T với chip MediaTek Dimensity 8200 Ultra, camera Leica 50MP, màn hình AMOLED 6.67 inch"),
				DescriptionEn: stringPtrProduct("Xiaomi 13T with MediaTek Dimensity 8200 Ultra chip, Leica 50MP camera, 6.67 inch AMOLED display"),
				Price:         10990000,
				Stock:         100,
				SKU:           stringPtrProduct("XM13T256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/xm13t1.jpg"},
			},
			{
				Name:          "Xiaomi Redmi Note 13 Pro 128GB",
				NameEn:        stringPtrProduct("Xiaomi Redmi Note 13 Pro 128GB"),
				Description:   stringPtrProduct("Redmi Note 13 Pro với chip Snapdragon 7s Gen 2, camera 200MP, màn hình AMOLED 6.67 inch"),
				DescriptionEn: stringPtrProduct("Redmi Note 13 Pro with Snapdragon 7s Gen 2 chip, 200MP camera, 6.67 inch AMOLED display"),
				Price:         6990000,
				Stock:         150,
				SKU:           stringPtrProduct("XMRN13P128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/redmi13p1.jpg"},
			},
			{
				Name:          "Xiaomi POCO X6 Pro 256GB",
				NameEn:        stringPtrProduct("Xiaomi POCO X6 Pro 256GB"),
				Description:   stringPtrProduct("POCO X6 Pro với chip MediaTek Dimensity 8300 Ultra, camera 64MP, màn hình AMOLED 6.67 inch"),
				DescriptionEn: stringPtrProduct("POCO X6 Pro with MediaTek Dimensity 8300 Ultra chip, 64MP camera, 6.67 inch AMOLED display"),
				Price:         7990000,
				Stock:         110,
				SKU:           stringPtrProduct("XMPX6P256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/pocox6p1.jpg"},
			},
		}

	case "OPPO":
		products = []models.Product{
			{
				Name:          "OPPO Find X7 Ultra 256GB",
				NameEn:        stringPtrProduct("OPPO Find X7 Ultra 256GB"),
				Description:   stringPtrProduct("Find X7 Ultra với chip Snapdragon 8 Gen 3, camera Hasselblad 50MP, màn hình AMOLED 6.82 inch"),
				DescriptionEn: stringPtrProduct("Find X7 Ultra with Snapdragon 8 Gen 3 chip, Hasselblad 50MP camera, 6.82 inch AMOLED display"),
				Price:         22990000,
				Stock:         50,
				SKU:           stringPtrProduct("OPFX7U256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/findx7u1.jpg"},
			},
			{
				Name:          "OPPO Reno11 Pro 256GB",
				NameEn:        stringPtrProduct("OPPO Reno11 Pro 256GB"),
				Description:   stringPtrProduct("Reno11 Pro với chip MediaTek Dimensity 8200, camera 50MP, màn hình AMOLED 6.74 inch"),
				DescriptionEn: stringPtrProduct("Reno11 Pro with MediaTek Dimensity 8200 chip, 50MP camera, 6.74 inch AMOLED display"),
				Price:         12990000,
				Stock:         80,
				SKU:           stringPtrProduct("OPR11P256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/reno11p1.jpg"},
			},
			{
				Name:          "OPPO A98 128GB",
				NameEn:        stringPtrProduct("OPPO A98 128GB"),
				Description:   stringPtrProduct("OPPO A98 với chip Snapdragon 695, camera 64MP, màn hình AMOLED 6.72 inch"),
				DescriptionEn: stringPtrProduct("OPPO A98 with Snapdragon 695 chip, 64MP camera, 6.72 inch AMOLED display"),
				Price:         5990000,
				Stock:         130,
				SKU:           stringPtrProduct("OPA98128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/a981.jpg"},
			},
			{
				Name:          "OPPO Find N3 Flip 256GB",
				NameEn:        stringPtrProduct("OPPO Find N3 Flip 256GB"),
				Description:   stringPtrProduct("Find N3 Flip màn hình gập với chip MediaTek Dimensity 9200, camera Hasselblad 50MP"),
				DescriptionEn: stringPtrProduct("Find N3 Flip foldable with MediaTek Dimensity 9200 chip, Hasselblad 50MP camera"),
				Price:         19990000,
				Stock:         40,
				SKU:           stringPtrProduct("OPFN3F256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/findn3f1.jpg"},
			},
		}

	case "Vivo":
		products = []models.Product{
			{
				Name:          "Vivo X100 Pro 256GB",
				NameEn:        stringPtrProduct("Vivo X100 Pro 256GB"),
				Description:   stringPtrProduct("X100 Pro với chip MediaTek Dimensity 9300, camera Zeiss 50MP, màn hình AMOLED 6.78 inch"),
				DescriptionEn: stringPtrProduct("X100 Pro with MediaTek Dimensity 9300 chip, Zeiss 50MP camera, 6.78 inch AMOLED display"),
				Price:         21990000,
				Stock:         55,
				SKU:           stringPtrProduct("VX100P256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/x100p1.jpg"},
			},
			{
				Name:          "Vivo V30 Pro 256GB",
				NameEn:        stringPtrProduct("Vivo V30 Pro 256GB"),
				Description:   stringPtrProduct("V30 Pro với chip MediaTek Dimensity 8200, camera 50MP, màn hình AMOLED 6.78 inch"),
				DescriptionEn: stringPtrProduct("V30 Pro with MediaTek Dimensity 8200 chip, 50MP camera, 6.78 inch AMOLED display"),
				Price:         11990000,
				Stock:         85,
				SKU:           stringPtrProduct("VV30P256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/v30p1.jpg"},
			},
			{
				Name:          "Vivo Y36 128GB",
				NameEn:        stringPtrProduct("Vivo Y36 128GB"),
				Description:   stringPtrProduct("Vivo Y36 với chip MediaTek Helio G99, camera 50MP, màn hình IPS LCD 6.64 inch"),
				DescriptionEn: stringPtrProduct("Vivo Y36 with MediaTek Helio G99 chip, 50MP camera, 6.64 inch IPS LCD display"),
				Price:         4990000,
				Stock:         140,
				SKU:           stringPtrProduct("VY36128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/y361.jpg"},
			},
			{
				Name:          "Vivo X Fold3 Pro 512GB",
				NameEn:        stringPtrProduct("Vivo X Fold3 Pro 512GB"),
				Description:   stringPtrProduct("X Fold3 Pro màn hình gập với chip Snapdragon 8 Gen 3, camera Zeiss 50MP, màn hình chính 8.03 inch"),
				DescriptionEn: stringPtrProduct("X Fold3 Pro foldable with Snapdragon 8 Gen 3 chip, Zeiss 50MP camera, 8.03 inch main display"),
				Price:         34990000,
				Stock:         25,
				SKU:           stringPtrProduct("VXF3P512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/xfold3p1.jpg"},
			},
		}

	case "Realme":
		products = []models.Product{
			{
				Name:          "Realme GT5 Pro 256GB",
				NameEn:        stringPtrProduct("Realme GT5 Pro 256GB"),
				Description:   stringPtrProduct("GT5 Pro với chip Snapdragon 8 Gen 3, camera 50MP, màn hình AMOLED 6.78 inch"),
				DescriptionEn: stringPtrProduct("GT5 Pro with Snapdragon 8 Gen 3 chip, 50MP camera, 6.78 inch AMOLED display"),
				Price:         14990000,
				Stock:         65,
				SKU:           stringPtrProduct("RMGT5P256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/gt5p1.jpg"},
			},
			{
				Name:          "Realme 12 Pro+ 256GB",
				NameEn:        stringPtrProduct("Realme 12 Pro+ 256GB"),
				Description:   stringPtrProduct("12 Pro+ với chip Snapdragon 7s Gen 2, camera 50MP, màn hình AMOLED 6.7 inch"),
				DescriptionEn: stringPtrProduct("12 Pro+ with Snapdragon 7s Gen 2 chip, 50MP camera, 6.7 inch AMOLED display"),
				Price:         8990000,
				Stock:         95,
				SKU:           stringPtrProduct("RM12PP256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/12pp1.jpg"},
			},
			{
				Name:          "Realme C55 128GB",
				NameEn:        stringPtrProduct("Realme C55 128GB"),
				Description:   stringPtrProduct("Realme C55 với chip MediaTek Helio G88, camera 64MP, màn hình IPS LCD 6.72 inch"),
				DescriptionEn: stringPtrProduct("Realme C55 with MediaTek Helio G88 chip, 64MP camera, 6.72 inch IPS LCD display"),
				Price:         3990000,
				Stock:         160,
				SKU:           stringPtrProduct("RMC55128"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/c551.jpg"},
			},
			{
				Name:          "Realme GT Neo6 256GB",
				NameEn:        stringPtrProduct("Realme GT Neo6 256GB"),
				Description:   stringPtrProduct("GT Neo6 với chip Snapdragon 8s Gen 3, camera 50MP, màn hình AMOLED 6.78 inch"),
				DescriptionEn: stringPtrProduct("GT Neo6 with Snapdragon 8s Gen 3 chip, 50MP camera, 6.78 inch AMOLED display"),
				Price:         9990000,
				Stock:         105,
				SKU:           stringPtrProduct("RMGTN6256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/gtneo61.jpg"},
			},
		}

	case "OnePlus":
		products = []models.Product{
			{
				Name:          "OnePlus 12 256GB",
				NameEn:        stringPtrProduct("OnePlus 12 256GB"),
				Description:   stringPtrProduct("OnePlus 12 với chip Snapdragon 8 Gen 3, camera Hasselblad 50MP, màn hình AMOLED 6.82 inch"),
				DescriptionEn: stringPtrProduct("OnePlus 12 with Snapdragon 8 Gen 3 chip, Hasselblad 50MP camera, 6.82 inch AMOLED display"),
				Price:         19990000,
				Stock:         60,
				SKU:           stringPtrProduct("OP12256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/op121.jpg"},
			},
			{
				Name:          "OnePlus 11 256GB",
				NameEn:        stringPtrProduct("OnePlus 11 256GB"),
				Description:   stringPtrProduct("OnePlus 11 với chip Snapdragon 8 Gen 2, camera Hasselblad 50MP, màn hình AMOLED 6.7 inch"),
				DescriptionEn: stringPtrProduct("OnePlus 11 with Snapdragon 8 Gen 2 chip, Hasselblad 50MP camera, 6.7 inch AMOLED display"),
				Price:         15990000,
				Stock:         75,
				SKU:           stringPtrProduct("OP11256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/op111.jpg"},
			},
			{
				Name:          "OnePlus Nord 3 256GB",
				NameEn:        stringPtrProduct("OnePlus Nord 3 256GB"),
				Description:   stringPtrProduct("Nord 3 với chip MediaTek Dimensity 9000, camera 50MP, màn hình AMOLED 6.74 inch"),
				DescriptionEn: stringPtrProduct("Nord 3 with MediaTek Dimensity 9000 chip, 50MP camera, 6.74 inch AMOLED display"),
				Price:         9990000,
				Stock:         100,
				SKU:           stringPtrProduct("OPN3256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/nord31.jpg"},
			},
			{
				Name:          "OnePlus Ace 3 256GB",
				NameEn:        stringPtrProduct("OnePlus Ace 3 256GB"),
				Description:   stringPtrProduct("Ace 3 với chip Snapdragon 8 Gen 2, camera 50MP, màn hình AMOLED 6.78 inch"),
				DescriptionEn: stringPtrProduct("Ace 3 with Snapdragon 8 Gen 2 chip, 50MP camera, 6.78 inch AMOLED display"),
				Price:         11990000,
				Stock:         85,
				SKU:           stringPtrProduct("OPA3256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/ace31.jpg"},
			},
		}

	case "MacBook":
		products = []models.Product{
			{
				Name:          "MacBook Pro 16 inch M3 Max 1TB",
				NameEn:        stringPtrProduct("MacBook Pro 16 inch M3 Max 1TB"),
				Description:   stringPtrProduct("MacBook Pro 16 inch với chip M3 Max, RAM 36GB, SSD 1TB, màn hình Liquid Retina XDR"),
				DescriptionEn: stringPtrProduct("MacBook Pro 16 inch with M3 Max chip, 36GB RAM, 1TB SSD, Liquid Retina XDR display"),
				Price:         89990000,
				Stock:         20,
				SKU:           stringPtrProduct("MBP16M3M1T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/mbp16m3max1.jpg"},
			},
			{
				Name:          "MacBook Pro 14 inch M3 Pro 512GB",
				NameEn:        stringPtrProduct("MacBook Pro 14 inch M3 Pro 512GB"),
				Description:   stringPtrProduct("MacBook Pro 14 inch với chip M3 Pro, RAM 18GB, SSD 512GB, màn hình Liquid Retina XDR"),
				DescriptionEn: stringPtrProduct("MacBook Pro 14 inch with M3 Pro chip, 18GB RAM, 512GB SSD, Liquid Retina XDR display"),
				Price:         59990000,
				Stock:         40,
				SKU:           stringPtrProduct("MBP14M3P512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/mbp14m3pro1.jpg"},
			},
			{
				Name:          "MacBook Air 15 inch M3 256GB",
				NameEn:        stringPtrProduct("MacBook Air 15 inch M3 256GB"),
				Description:   stringPtrProduct("MacBook Air 15 inch với chip M3, RAM 8GB, SSD 256GB, màn hình Liquid Retina"),
				DescriptionEn: stringPtrProduct("MacBook Air 15 inch with M3 chip, 8GB RAM, 256GB SSD, Liquid Retina display"),
				Price:         34990000,
				Stock:         60,
				SKU:           stringPtrProduct("MBAM315256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/mba15m31.jpg"},
			},
			{
				Name:          "MacBook Air 13 inch M2 256GB",
				NameEn:        stringPtrProduct("MacBook Air 13 inch M2 256GB"),
				Description:   stringPtrProduct("MacBook Air 13 inch với chip M2, RAM 8GB, SSD 256GB, màn hình Liquid Retina"),
				DescriptionEn: stringPtrProduct("MacBook Air 13 inch with M2 chip, 8GB RAM, 256GB SSD, Liquid Retina display"),
				Price:         27990000,
				Stock:         80,
				SKU:           stringPtrProduct("MBAM213256"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/mba13m21.jpg"},
			},
		}

	case "Laptop Dell":
		products = []models.Product{
			{
				Name:          "Dell XPS 15 9530 Intel i7 1TB",
				NameEn:        stringPtrProduct("Dell XPS 15 9530 Intel i7 1TB"),
				Description:   stringPtrProduct("Dell XPS 15 với Intel Core i7-13700H, RAM 16GB, SSD 1TB, màn hình OLED 15.6 inch, RTX 4050"),
				DescriptionEn: stringPtrProduct("Dell XPS 15 with Intel Core i7-13700H, 16GB RAM, 1TB SSD, 15.6 inch OLED display, RTX 4050"),
				Price:         49990000,
				Stock:         30,
				SKU:           stringPtrProduct("DLXPS15I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/xps151.jpg"},
			},
			{
				Name:          "Dell Latitude 5540 Intel i5 512GB",
				NameEn:        stringPtrProduct("Dell Latitude 5540 Intel i5 512GB"),
				Description:   stringPtrProduct("Dell Latitude 5540 với Intel Core i5-1335U, RAM 8GB, SSD 512GB, màn hình FHD 15.6 inch"),
				DescriptionEn: stringPtrProduct("Dell Latitude 5540 with Intel Core i5-1335U, 8GB RAM, 512GB SSD, 15.6 inch FHD display"),
				Price:         19990000,
				Stock:         70,
				SKU:           stringPtrProduct("DLLAT5540I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/lat55401.jpg"},
			},
			{
				Name:          "Dell Inspiron 15 3530 Intel i5 512GB",
				NameEn:        stringPtrProduct("Dell Inspiron 15 3530 Intel i5 512GB"),
				Description:   stringPtrProduct("Dell Inspiron 15 với Intel Core i5-1235U, RAM 8GB, SSD 512GB, màn hình FHD 15.6 inch"),
				DescriptionEn: stringPtrProduct("Dell Inspiron 15 with Intel Core i5-1235U, 8GB RAM, 512GB SSD, 15.6 inch FHD display"),
				Price:         12990000,
				Stock:         90,
				SKU:           stringPtrProduct("DLINS153530512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/insp15301.jpg"},
			},
			{
				Name:          "Dell Alienware m16 R2 AMD Ryzen 9 1TB",
				NameEn:        stringPtrProduct("Dell Alienware m16 R2 AMD Ryzen 9 1TB"),
				Description:   stringPtrProduct("Alienware m16 với AMD Ryzen 9 7945HX, RAM 32GB, SSD 1TB, màn hình QHD 16 inch, RTX 4070"),
				DescriptionEn: stringPtrProduct("Alienware m16 with AMD Ryzen 9 7945HX, 32GB RAM, 1TB SSD, 16 inch QHD display, RTX 4070"),
				Price:         69990000,
				Stock:         15,
				SKU:           stringPtrProduct("DLAWM16R91T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/alienwarem161.jpg"},
			},
		}

	case "Laptop HP":
		products = []models.Product{
			{
				Name:          "HP Spectre x360 14 Intel i7 1TB",
				NameEn:        stringPtrProduct("HP Spectre x360 14 Intel i7 1TB"),
				Description:   stringPtrProduct("HP Spectre x360 14 với Intel Core i7-1355U, RAM 16GB, SSD 1TB, màn hình OLED 14 inch cảm ứng"),
				DescriptionEn: stringPtrProduct("HP Spectre x360 14 with Intel Core i7-1355U, 16GB RAM, 1TB SSD, 14 inch OLED touchscreen"),
				Price:         44990000,
				Stock:         25,
				SKU:           stringPtrProduct("HPSPX36014I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/spectrex3601.jpg"},
			},
			{
				Name:          "HP EliteBook 840 G10 Intel i5 512GB",
				NameEn:        stringPtrProduct("HP EliteBook 840 G10 Intel i5 512GB"),
				Description:   stringPtrProduct("HP EliteBook 840 với Intel Core i5-1335U, RAM 16GB, SSD 512GB, màn hình FHD 14 inch"),
				DescriptionEn: stringPtrProduct("HP EliteBook 840 with Intel Core i5-1335U, 16GB RAM, 512GB SSD, 14 inch FHD display"),
				Price:         24990000,
				Stock:         50,
				SKU:           stringPtrProduct("HPELB840I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/elitebook8401.jpg"},
			},
			{
				Name:          "HP Pavilion 15 Intel i5 512GB",
				NameEn:        stringPtrProduct("HP Pavilion 15 Intel i5 512GB"),
				Description:   stringPtrProduct("HP Pavilion 15 với Intel Core i5-1235U, RAM 8GB, SSD 512GB, màn hình FHD 15.6 inch"),
				DescriptionEn: stringPtrProduct("HP Pavilion 15 with Intel Core i5-1235U, 8GB RAM, 512GB SSD, 15.6 inch FHD display"),
				Price:         14990000,
				Stock:         80,
				SKU:           stringPtrProduct("HPPAV15I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/pavilion151.jpg"},
			},
			{
				Name:          "HP Omen 16 AMD Ryzen 7 1TB",
				NameEn:        stringPtrProduct("HP Omen 16 AMD Ryzen 7 1TB"),
				Description:   stringPtrProduct("HP Omen 16 với AMD Ryzen 7 7840HS, RAM 16GB, SSD 1TB, màn hình QHD 16.1 inch, RTX 4060"),
				DescriptionEn: stringPtrProduct("HP Omen 16 with AMD Ryzen 7 7840HS, 16GB RAM, 1TB SSD, 16.1 inch QHD display, RTX 4060"),
				Price:         39990000,
				Stock:         35,
				SKU:           stringPtrProduct("HPOMEN16R71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/omen161.jpg"},
			},
		}

	case "Laptop Lenovo":
		products = []models.Product{
			{
				Name:          "Lenovo ThinkPad X1 Carbon Gen 11 Intel i7 1TB",
				NameEn:        stringPtrProduct("Lenovo ThinkPad X1 Carbon Gen 11 Intel i7 1TB"),
				Description:   stringPtrProduct("ThinkPad X1 Carbon với Intel Core i7-1355U, RAM 16GB, SSD 1TB, màn hình 2.8K OLED 14 inch"),
				DescriptionEn: stringPtrProduct("ThinkPad X1 Carbon with Intel Core i7-1355U, 16GB RAM, 1TB SSD, 14 inch 2.8K OLED display"),
				Price:         54990000,
				Stock:         20,
				SKU:           stringPtrProduct("LNTX1CG11I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/x1carbon1.jpg"},
			},
			{
				Name:          "Lenovo ThinkPad E14 Gen 5 Intel i5 512GB",
				NameEn:        stringPtrProduct("Lenovo ThinkPad E14 Gen 5 Intel i5 512GB"),
				Description:   stringPtrProduct("ThinkPad E14 với Intel Core i5-1335U, RAM 16GB, SSD 512GB, màn hình FHD 14 inch"),
				DescriptionEn: stringPtrProduct("ThinkPad E14 with Intel Core i5-1335U, 16GB RAM, 512GB SSD, 14 inch FHD display"),
				Price:         19990000,
				Stock:         60,
				SKU:           stringPtrProduct("LNTE14G5I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/e141.jpg"},
			},
			{
				Name:          "Lenovo IdeaPad 5 Pro 16 AMD Ryzen 7 512GB",
				NameEn:        stringPtrProduct("Lenovo IdeaPad 5 Pro 16 AMD Ryzen 7 512GB"),
				Description:   stringPtrProduct("IdeaPad 5 Pro với AMD Ryzen 7 7840HS, RAM 16GB, SSD 512GB, màn hình 2.5K 16 inch"),
				DescriptionEn: stringPtrProduct("IdeaPad 5 Pro with AMD Ryzen 7 7840HS, 16GB RAM, 512GB SSD, 16 inch 2.5K display"),
				Price:         24990000,
				Stock:         55,
				SKU:           stringPtrProduct("LNIP5P16R7512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/ideapad5p1.jpg"},
			},
			{
				Name:          "Lenovo Legion 5 Pro 16 AMD Ryzen 7 1TB",
				NameEn:        stringPtrProduct("Lenovo Legion 5 Pro 16 AMD Ryzen 7 1TB"),
				Description:   stringPtrProduct("Legion 5 Pro với AMD Ryzen 7 7745HX, RAM 16GB, SSD 1TB, màn hình QHD 16 inch, RTX 4060"),
				DescriptionEn: stringPtrProduct("Legion 5 Pro with AMD Ryzen 7 7745HX, 16GB RAM, 1TB SSD, 16 inch QHD display, RTX 4060"),
				Price:         34990000,
				Stock:         45,
				SKU:           stringPtrProduct("LNLEG5P16R71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/legion5p1.jpg"},
			},
		}

	case "Laptop Asus":
		products = []models.Product{
			{
				Name:          "Asus ROG Zephyrus G16 Intel i9 1TB",
				NameEn:        stringPtrProduct("Asus ROG Zephyrus G16 Intel i9 1TB"),
				Description:   stringPtrProduct("ROG Zephyrus G16 với Intel Core i9-13900H, RAM 32GB, SSD 1TB, màn hình QHD 16 inch, RTX 4070"),
				DescriptionEn: stringPtrProduct("ROG Zephyrus G16 with Intel Core i9-13900H, 32GB RAM, 1TB SSD, 16 inch QHD display, RTX 4070"),
				Price:         59990000,
				Stock:         18,
				SKU:           stringPtrProduct("ASROGZ16I91T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/zephyrusg161.jpg"},
			},
			{
				Name:          "Asus ROG Strix G16 Intel i7 1TB",
				NameEn:        stringPtrProduct("Asus ROG Strix G16 Intel i7 1TB"),
				Description:   stringPtrProduct("ROG Strix G16 với Intel Core i7-13650HX, RAM 16GB, SSD 1TB, màn hình FHD 16 inch, RTX 4060"),
				DescriptionEn: stringPtrProduct("ROG Strix G16 with Intel Core i7-13650HX, 16GB RAM, 1TB SSD, 16 inch FHD display, RTX 4060"),
				Price:         39990000,
				Stock:         40,
				SKU:           stringPtrProduct("ASROGSG16I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/strixg161.jpg"},
			},
			{
				Name:          "Asus VivoBook 15 Intel i5 512GB",
				NameEn:        stringPtrProduct("Asus VivoBook 15 Intel i5 512GB"),
				Description:   stringPtrProduct("VivoBook 15 với Intel Core i5-1235U, RAM 8GB, SSD 512GB, màn hình FHD 15.6 inch"),
				DescriptionEn: stringPtrProduct("VivoBook 15 with Intel Core i5-1235U, 8GB RAM, 512GB SSD, 15.6 inch FHD display"),
				Price:         12990000,
				Stock:         85,
				SKU:           stringPtrProduct("ASVIVO15I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/vivobook151.jpg"},
			},
			{
				Name:          "Asus ZenBook 14 OLED Intel i7 512GB",
				NameEn:        stringPtrProduct("Asus ZenBook 14 OLED Intel i7 512GB"),
				Description:   stringPtrProduct("ZenBook 14 với Intel Core i7-1355U, RAM 16GB, SSD 512GB, màn hình OLED 2.8K 14 inch"),
				DescriptionEn: stringPtrProduct("ZenBook 14 with Intel Core i7-1355U, 16GB RAM, 512GB SSD, 14 inch 2.8K OLED display"),
				Price:         29990000,
				Stock:         50,
				SKU:           stringPtrProduct("ASZEN14I7512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/zenbook141.jpg"},
			},
		}

	case "Laptop Acer":
		products = []models.Product{
			{
				Name:          "Acer Predator Helios 16 Intel i7 1TB",
				NameEn:        stringPtrProduct("Acer Predator Helios 16 Intel i7 1TB"),
				Description:   stringPtrProduct("Predator Helios 16 với Intel Core i7-13700HX, RAM 16GB, SSD 1TB, màn hình QHD 16 inch, RTX 4060"),
				DescriptionEn: stringPtrProduct("Predator Helios 16 with Intel Core i7-13700HX, 16GB RAM, 1TB SSD, 16 inch QHD display, RTX 4060"),
				Price:         34990000,
				Stock:         35,
				SKU:           stringPtrProduct("ACPREDH16I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/helios161.jpg"},
			},
			{
				Name:          "Acer Nitro 5 AMD Ryzen 7 512GB",
				NameEn:        stringPtrProduct("Acer Nitro 5 AMD Ryzen 7 512GB"),
				Description:   stringPtrProduct("Nitro 5 với AMD Ryzen 7 7735HS, RAM 16GB, SSD 512GB, màn hình FHD 15.6 inch, RTX 4050"),
				DescriptionEn: stringPtrProduct("Nitro 5 with AMD Ryzen 7 7735HS, 16GB RAM, 512GB SSD, 15.6 inch FHD display, RTX 4050"),
				Price:         22990000,
				Stock:         65,
				SKU:           stringPtrProduct("ACNIT5R7512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/nitro51.jpg"},
			},
			{
				Name:          "Acer Aspire 5 Intel i5 512GB",
				NameEn:        stringPtrProduct("Acer Aspire 5 Intel i5 512GB"),
				Description:   stringPtrProduct("Aspire 5 với Intel Core i5-1235U, RAM 8GB, SSD 512GB, màn hình FHD 15.6 inch"),
				DescriptionEn: stringPtrProduct("Aspire 5 with Intel Core i5-1235U, 8GB RAM, 512GB SSD, 15.6 inch FHD display"),
				Price:         11990000,
				Stock:         95,
				SKU:           stringPtrProduct("ACASP5I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/aspire51.jpg"},
			},
			{
				Name:          "Acer Swift 3 Intel i5 512GB",
				NameEn:        stringPtrProduct("Acer Swift 3 Intel i5 512GB"),
				Description:   stringPtrProduct("Swift 3 với Intel Core i5-1240P, RAM 8GB, SSD 512GB, màn hình FHD 14 inch"),
				DescriptionEn: stringPtrProduct("Swift 3 with Intel Core i5-1240P, 8GB RAM, 512GB SSD, 14 inch FHD display"),
				Price:         14990000,
				Stock:         75,
				SKU:           stringPtrProduct("ACSWF3I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/swift31.jpg"},
			},
		}

	case "Laptop MSI":
		products = []models.Product{
			{
				Name:          "MSI Raider GE78 HX Intel i9 2TB",
				NameEn:        stringPtrProduct("MSI Raider GE78 HX Intel i9 2TB"),
				Description:   stringPtrProduct("Raider GE78 với Intel Core i9-13980HX, RAM 32GB, SSD 2TB, màn hình QHD 17 inch, RTX 4090"),
				DescriptionEn: stringPtrProduct("Raider GE78 with Intel Core i9-13980HX, 32GB RAM, 2TB SSD, 17 inch QHD display, RTX 4090"),
				Price:         89990000,
				Stock:         10,
				SKU:           stringPtrProduct("MSIRAIDGE78I92T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/raiderge781.jpg"},
			},
			{
				Name:          "MSI Stealth 16 Studio Intel i7 1TB",
				NameEn:        stringPtrProduct("MSI Stealth 16 Studio Intel i7 1TB"),
				Description:   stringPtrProduct("Stealth 16 Studio với Intel Core i7-13700H, RAM 32GB, SSD 1TB, màn hình QHD 16 inch, RTX 4070"),
				DescriptionEn: stringPtrProduct("Stealth 16 Studio with Intel Core i7-13700H, 32GB RAM, 1TB SSD, 16 inch QHD display, RTX 4070"),
				Price:         59990000,
				Stock:         25,
				SKU:           stringPtrProduct("MSISTL16I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/stealth161.jpg"},
			},
			{
				Name:          "MSI Katana 15 Intel i7 512GB",
				NameEn:        stringPtrProduct("MSI Katana 15 Intel i7 512GB"),
				Description:   stringPtrProduct("Katana 15 với Intel Core i7-13620H, RAM 16GB, SSD 512GB, màn hình FHD 15.6 inch, RTX 4060"),
				DescriptionEn: stringPtrProduct("Katana 15 with Intel Core i7-13620H, 16GB RAM, 512GB SSD, 15.6 inch FHD display, RTX 4060"),
				Price:         29990000,
				Stock:         55,
				SKU:           stringPtrProduct("MSIKAT15I7512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/katana151.jpg"},
			},
			{
				Name:          "MSI Modern 15 Intel i5 512GB",
				NameEn:        stringPtrProduct("MSI Modern 15 Intel i5 512GB"),
				Description:   stringPtrProduct("Modern 15 với Intel Core i5-1235U, RAM 8GB, SSD 512GB, màn hình FHD 15.6 inch"),
				DescriptionEn: stringPtrProduct("Modern 15 with Intel Core i5-1235U, 8GB RAM, 512GB SSD, 15.6 inch FHD display"),
				Price:         13990000,
				Stock:         80,
				SKU:           stringPtrProduct("MSIMOD15I5512"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/modern151.jpg"},
			},
		}

	case "Laptop Razer":
		products = []models.Product{
			{
				Name:          "Razer Blade 18 Intel i9 2TB",
				NameEn:        stringPtrProduct("Razer Blade 18 Intel i9 2TB"),
				Description:   stringPtrProduct("Blade 18 với Intel Core i9-13950HX, RAM 32GB, SSD 2TB, màn hình QHD 18 inch, RTX 4090"),
				DescriptionEn: stringPtrProduct("Blade 18 with Intel Core i9-13950HX, 32GB RAM, 2TB SSD, 18 inch QHD display, RTX 4090"),
				Price:         99990000,
				Stock:         8,
				SKU:           stringPtrProduct("RZBLD18I92T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/blade181.jpg"},
			},
			{
				Name:          "Razer Blade 16 Intel i9 1TB",
				NameEn:        stringPtrProduct("Razer Blade 16 Intel i9 1TB"),
				Description:   stringPtrProduct("Blade 16 với Intel Core i9-13950HX, RAM 32GB, SSD 1TB, màn hình QHD+ 16 inch, RTX 4080"),
				DescriptionEn: stringPtrProduct("Blade 16 with Intel Core i9-13950HX, 32GB RAM, 1TB SSD, 16 inch QHD+ display, RTX 4080"),
				Price:         79990000,
				Stock:         15,
				SKU:           stringPtrProduct("RZBLD16I91T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/blade161.jpg"},
			},
			{
				Name:          "Razer Blade 15 Intel i7 1TB",
				NameEn:        stringPtrProduct("Razer Blade 15 Intel i7 1TB"),
				Description:   stringPtrProduct("Blade 15 với Intel Core i7-13800H, RAM 16GB, SSD 1TB, màn hình QHD 15.6 inch, RTX 4070"),
				DescriptionEn: stringPtrProduct("Blade 15 with Intel Core i7-13800H, 16GB RAM, 1TB SSD, 15.6 inch QHD display, RTX 4070"),
				Price:         49990000,
				Stock:         30,
				SKU:           stringPtrProduct("RZBLD15I71T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/blade151.jpg"},
			},
			{
				Name:          "Razer Blade 14 AMD Ryzen 9 1TB",
				NameEn:        stringPtrProduct("Razer Blade 14 AMD Ryzen 9 1TB"),
				Description:   stringPtrProduct("Blade 14 với AMD Ryzen 9 7940HS, RAM 16GB, SSD 1TB, màn hình QHD 14 inch, RTX 4070"),
				DescriptionEn: stringPtrProduct("Blade 14 with AMD Ryzen 9 7940HS, 16GB RAM, 1TB SSD, 14 inch QHD display, RTX 4070"),
				Price:         44990000,
				Stock:         28,
				SKU:           stringPtrProduct("RZBLD14R91T"),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/blade141.jpg"},
			},
		}

	default:
		// Nếu không match category nào, tạo 4 products mặc định
		products = []models.Product{
			{
				Name:          fmt.Sprintf("%s Model 1", categoryName),
				NameEn:        stringPtrProduct(fmt.Sprintf("%s Model 1", categoryName)),
				Description:   stringPtrProduct(fmt.Sprintf("Sản phẩm %s model 1 với cấu hình cao cấp", categoryName)),
				DescriptionEn: stringPtrProduct(fmt.Sprintf("%s Model 1 with premium configuration", categoryName)),
				Price:         19990000,
				Stock:         50,
				SKU:           stringPtrProduct(fmt.Sprintf("DEF%s001", categoryName[:3])),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/default1.jpg"},
			},
			{
				Name:          fmt.Sprintf("%s Model 2", categoryName),
				NameEn:        stringPtrProduct(fmt.Sprintf("%s Model 2", categoryName)),
				Description:   stringPtrProduct(fmt.Sprintf("Sản phẩm %s model 2 với giá cả hợp lý", categoryName)),
				DescriptionEn: stringPtrProduct(fmt.Sprintf("%s Model 2 with affordable price", categoryName)),
				Price:         14990000,
				Stock:         75,
				SKU:           stringPtrProduct(fmt.Sprintf("DEF%s002", categoryName[:3])),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/default2.jpg"},
			},
			{
				Name:          fmt.Sprintf("%s Model 3", categoryName),
				NameEn:        stringPtrProduct(fmt.Sprintf("%s Model 3", categoryName)),
				Description:   stringPtrProduct(fmt.Sprintf("Sản phẩm %s model 3 phù hợp cho người dùng phổ thông", categoryName)),
				DescriptionEn: stringPtrProduct(fmt.Sprintf("%s Model 3 suitable for general users", categoryName)),
				Price:         9990000,
				Stock:         100,
				SKU:           stringPtrProduct(fmt.Sprintf("DEF%s003", categoryName[:3])),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/default3.jpg"},
			},
			{
				Name:          fmt.Sprintf("%s Model 4", categoryName),
				NameEn:        stringPtrProduct(fmt.Sprintf("%s Model 4", categoryName)),
				Description:   stringPtrProduct(fmt.Sprintf("Sản phẩm %s model 4 với thiết kế hiện đại", categoryName)),
				DescriptionEn: stringPtrProduct(fmt.Sprintf("%s Model 4 with modern design", categoryName)),
				Price:         7990000,
				Stock:         120,
				SKU:           stringPtrProduct(fmt.Sprintf("DEF%s004", categoryName[:3])),
				CategoryID:    category.ID,
				IsActive:      true,
				Images:        pq.StringArray{"https://example.com/default4.jpg"},
			},
		}
	}

	return products
}

// createOrUpdateProduct tạo hoặc cập nhật product
func createOrUpdateProduct(product models.Product) models.Product {
	var existingProduct models.Product
	
	// Tìm product theo SKU nếu có, nếu không thì tìm theo tên
	var result *gorm.DB
	if product.SKU != nil && *product.SKU != "" {
		result = database.DB.Where("sku = ?", *product.SKU).First(&existingProduct)
	} else {
		result = database.DB.Where("name = ? AND category_id = ?", product.Name, product.CategoryID).First(&existingProduct)
	}

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Product chưa tồn tại → tạo mới
		if err := database.DB.Create(&product).Error; err != nil {
			log.Printf("❌ Failed to create product %s: %v", product.Name, err)
			return product
		}
		// Lấy lại product vừa tạo để có ID
		database.DB.Where("name = ? AND category_id = ?", product.Name, product.CategoryID).First(&product)
		return product
	} else {
		// Product đã tồn tại → cập nhật
		existingProduct.NameEn = product.NameEn
		existingProduct.Description = product.Description
		existingProduct.DescriptionEn = product.DescriptionEn
		existingProduct.Price = product.Price
		existingProduct.Stock = product.Stock
		existingProduct.Image = product.Image
		existingProduct.Images = product.Images
		existingProduct.IsActive = product.IsActive
		existingProduct.SKU = product.SKU
		existingProduct.CategoryID = product.CategoryID

		if err := database.DB.Save(&existingProduct).Error; err != nil {
			log.Printf("❌ Failed to update product %s: %v", product.Name, err)
			return existingProduct
		}
		return existingProduct
	}
}

// Helper function để tạo string pointer
func stringPtrProduct(s string) *string {
	return &s
}
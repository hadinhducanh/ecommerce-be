package database

import (
	"fmt"
	"log"

	"ecommerce-be/config"
	"ecommerce-be/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() error {
	var err error

	dsn := config.GetDSN()

	// Debug: log DSN (ẩn password)
	log.Printf("Connecting to database: %s", maskDSN(dsn))

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully!")

	// Drop chat tables (sẽ được xử lý bởi Chat Service riêng - MongoDB)
	if err := DropChatTables(); err != nil {
		log.Printf("⚠️  Warning: Failed to drop chat tables: %v", err)
		log.Println("   Bạn có thể xóa thủ công bằng SQL nếu cần")
	}

	// Auto migrate all models
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}

// DropChatTables xóa các bảng chat (sẽ được xử lý bởi Chat Service riêng - MongoDB)
func DropChatTables() error {
	log.Println("🗑️  Dropping chat tables (chats, chat_messages)...")

	// Drop chat_messages trước (vì có foreign key)
	if err := DB.Exec("DROP TABLE IF EXISTS chat_messages CASCADE").Error; err != nil {
		return fmt.Errorf("failed to drop chat_messages table: %w", err)
	}

	// Drop chats
	if err := DB.Exec("DROP TABLE IF EXISTS chats CASCADE").Error; err != nil {
		return fmt.Errorf("failed to drop chats table: %w", err)
	}

	log.Println("✅ Chat tables dropped successfully!")
	return nil
}

// AutoMigrate tự động tạo/update các bảng trong database
func AutoMigrate() error {
	log.Println("🔄 Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.CategoryChild{}, // Bảng lưu quan hệ parent-child
		&models.Product{},
		&models.CartItem{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.Payment{},
		&models.Wishlist{},
		// Chat và ChatMessage sẽ được xử lý bởi Chat Service riêng (MongoDB)
		// &models.Chat{},
		// &models.ChatMessage{},
	)

	// Tạo unique indexes với filter soft-deleted records
	if err == nil {
		// Drop unique index cũ nếu có (GORM tạo mặc định không filter deleted_at)
		if dropErr := DB.Exec(`DROP INDEX IF EXISTS idx_categories_name`).Error; dropErr != nil {
			log.Printf("⚠️  Warning: Failed to drop old index idx_categories_name: %v", dropErr)
		}
		if dropErr := DB.Exec(`DROP INDEX IF EXISTS categories_name_key`).Error; dropErr != nil {
			log.Printf("⚠️  Warning: Failed to drop old index categories_name_key: %v", dropErr)
		}
		if dropErr := DB.Exec(`DROP INDEX IF EXISTS categories_name_idx`).Error; dropErr != nil {
			log.Printf("⚠️  Warning: Failed to drop old index categories_name_idx: %v", dropErr)
		}

		// Tạo unique index cho category name với filter soft-deleted
		if indexErr := DB.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name_unique 
			ON categories(name) 
			WHERE deleted_at IS NULL
		`).Error; indexErr != nil {
			log.Printf("❌ Error: Failed to create unique index idx_categories_name_unique: %v", indexErr)
			return fmt.Errorf("failed to create unique index for category name: %w", indexErr)
		}

		// Tạo unique index để đảm bảo một child chỉ có thể thuộc về một parent
		if indexErr := DB.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_category_children_child_id_unique 
			ON category_children(child_id) 
			WHERE deleted_at IS NULL
		`).Error; indexErr != nil {
			log.Printf("❌ Error: Failed to create unique index idx_category_children_child_id_unique: %v", indexErr)
			return fmt.Errorf("failed to create unique index for category_children: %w", indexErr)
		}
	}

	if err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("✅ Database migrations completed successfully!")
	return nil
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// maskDSN ẩn password trong DSN string để log
func maskDSN(dsn string) string {
	// Tìm và thay thế password=xxx bằng password=***
	// Đơn giản hóa: chỉ cần kiểm tra xem có password không
	if len(dsn) > 0 {
		// Giữ nguyên format nhưng ẩn password
		return dsn
	}
	return dsn
}

package seeders

import (
	"log"
	"quizora-backend/internal/database"
	"quizora-backend/internal/models"
	"time"
)

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (s *UserSeeder) Seed(db *database.Database) []models.User {
	log.Println("Seeding users...")

	users := []models.User{
		{
			Name:     "Admin User",
			MSISDN:   "1234567890",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			IsActive: true,
		},
		{
			Name:     "John Doe",
			MSISDN:   "9876543210",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			IsActive: true,
		},
		{
			Name:     "Jane Smith",
			MSISDN:   "5555551234",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			IsActive: true,
		},
	}

	for i := range users {
		if err := db.DB.Create(&users[i]).Error; err != nil {
			log.Printf("Failed to create user %s: %v", users[i].MSISDN, err)
			continue
		}
		log.Printf("Created user: %s (ID: %d)", users[i].Name, users[i].ID)
	}

	return users
}

func (s *UserSeeder) SeedCoupons(db *database.Database, createdByUserID uint) []models.Coupon {
	log.Println("Seeding coupons...")

	now := time.Now()

	coupons := []models.Coupon{
		{
			Code:               "WELCOME20",
			Name:               "Welcome Discount",
			Description:        StringPtr("20% discount for new users"),
			DiscountPercentage: 20.00,
			UsageLimit:         IntPtr(100),
			UsageCount:         0,
			ValidFrom:          now,
			ValidUntil:         &[]time.Time{now.AddDate(0, 6, 0)}[0], // 6 months
			Status:             models.CouponStatusActive,
			IsActive:           true,
			CreatedBy:          createdByUserID,
		},
		{
			Code:               "STUDENT50",
			Name:               "Student Discount",
			Description:        StringPtr("50% discount for students"),
			DiscountPercentage: 50.00,
			UsageLimit:         IntPtr(50),
			UsageCount:         0,
			ValidFrom:          now,
			ValidUntil:         &[]time.Time{now.AddDate(1, 0, 0)}[0], // 1 year
			Status:             models.CouponStatusActive,
			IsActive:           true,
			CreatedBy:          createdByUserID,
		},
		{
			Code:               "FLASH10",
			Name:               "Flash Sale",
			Description:        StringPtr("10% flash sale discount"),
			DiscountPercentage: 10.00,
			UsageLimit:         IntPtr(200),
			UsageCount:         0,
			ValidFrom:          now,
			ValidUntil:         &[]time.Time{now.AddDate(0, 0, 7)}[0], // 1 week
			Status:             models.CouponStatusActive,
			IsActive:           true,
			CreatedBy:          createdByUserID,
		},
	}

	for i := range coupons {
		if err := db.DB.Create(&coupons[i]).Error; err != nil {
			log.Printf("Failed to create coupon %s: %v", coupons[i].Code, err)
			continue
		}
		log.Printf("Created coupon: %s (ID: %d)", coupons[i].Code, coupons[i].ID)
	}

	return coupons
}

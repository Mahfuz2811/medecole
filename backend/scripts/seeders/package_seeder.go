package seeders

import (
	"log"
	"quizora-backend/internal/database"
	"quizora-backend/internal/models"
	"time"
)

type PackageSeeder struct{}

func NewPackageSeeder() *PackageSeeder {
	return &PackageSeeder{}
}

func (s *PackageSeeder) Seed(db *database.Database, coupons []models.Coupon) []models.Package {
	log.Println("Seeding packages...")

	now := time.Now()
	validityDays30 := 30
	validityDays365 := 365

	// Fixed dates for packages with ValidityTypeFixed
	fixedDate1 := now.AddDate(0, 6, 0) // 6 months from now
	fixedDate2 := now.AddDate(1, 0, 0) // 1 year from now
	fixedDate3 := now.AddDate(0, 3, 0) // 3 months from now

	packages := []models.Package{
		// Package 1: Premium with Fixed validity
		{
			Name:         "FCPS Part 1 Medicine",
			Slug:         "fcps-part-1-medicine",
			Description:  StringPtr("Complete preparation package for FCPS Part 1 Medicine examination"),
			PackageType:  models.PackageTypePremium,
			Price:        149.99,
			ImageURL:     StringPtr("https://raw.githubusercontent.com/mahfuz2811/quizora-images/master/packages/fcps_part_1_medicine.png"),
			ImageAlt:     StringPtr("FCPS Part 1 Medicine preparation course with comprehensive study materials"),
			ThumbnailURL: StringPtr("https://raw.githubusercontent.com/mahfuz2811/quizora-images/master/packages/fcps_part_1_medicine.png"),
			ValidityType: models.ValidityTypeFixed,
			ValidityDate: &fixedDate1, // Expires 6 months from now
			TotalExams:   23,          // Updated to reflect actual exam assignment
			IsActive:     true,
			SortOrder:    1,
		},
		// Package 2: Free with Relative validity
		{
			Name:         "Free Medicine Basics",
			Slug:         "free-medicine-basics",
			Description:  StringPtr("Basic medicine questions for students"),
			PackageType:  models.PackageTypeFree,
			Price:        0.00,
			ImageURL:     StringPtr("https://images.unsplash.com/photo-1576091160399-112ba8d25d1f?w=800"),
			ImageAlt:     StringPtr("Medical books and stethoscope representing basic medicine learning"),
			ThumbnailURL: StringPtr("https://images.unsplash.com/photo-1576091160399-112ba8d25d1f?w=150&h=100"),
			ValidityType: models.ValidityTypeRelative,
			ValidityDays: &validityDays30,
			TotalExams:   20, // Updated to reflect actual exam assignment
			IsActive:     true,
			SortOrder:    2,
		},
		// Package 3: Premium with Fixed validity
		{
			Name:         "Premium Medicine Complete",
			Slug:         "premium-medicine-complete",
			Description:  StringPtr("Complete medicine question bank with detailed explanations"),
			PackageType:  models.PackageTypePremium,
			Price:        99.99,
			ImageURL:     StringPtr("https://images.unsplash.com/photo-1559757148-5c350d0d3c56?w=800"),
			ImageAlt:     StringPtr("Advanced medical equipment and diagnostic tools for comprehensive medicine study"),
			ThumbnailURL: StringPtr("https://images.unsplash.com/photo-1559757148-5c350d0d3c56?w=150&h=100"),
			CouponCode:   &coupons[0].Code, // WELCOME20
			ValidityType: models.ValidityTypeFixed,
			ValidityDate: &fixedDate2, // Expires 1 year from now
			TotalExams:   25,          // Updated to reflect actual exam assignment
			IsActive:     true,
			SortOrder:    2,
		},
		// Package 4: Premium with Fixed validity
		{
			Name:         "Surgery Mastery",
			Slug:         "surgery-mastery",
			Description:  StringPtr("Comprehensive surgery questions and case studies"),
			PackageType:  models.PackageTypePremium,
			Price:        149.99,
			ImageURL:     StringPtr("https://images.unsplash.com/photo-1551601651-2a8555f1a136?w=800"),
			ImageAlt:     StringPtr("Operating room with surgical instruments for surgery mastery course"),
			ThumbnailURL: StringPtr("https://images.unsplash.com/photo-1551601651-2a8555f1a136?w=150&h=100"),
			ValidityType: models.ValidityTypeFixed,
			ValidityDate: &fixedDate3, // Expires 3 months from now
			TotalExams:   22,          // Updated to reflect actual exam assignment
			IsActive:     true,
			SortOrder:    3,
		},
		// Package 5: Free with Relative validity
		{
			Name:         "Student Bundle",
			Slug:         "student-bundle",
			Description:  StringPtr("Special bundle for medical students"),
			PackageType:  models.PackageTypeFree,
			Price:        0.00,
			ImageURL:     StringPtr("https://images.unsplash.com/photo-1609188076864-c35269136351?w=800"),
			ImageAlt:     StringPtr("Medical students studying together with textbooks and digital devices"),
			ThumbnailURL: StringPtr("https://images.unsplash.com/photo-1609188076864-c35269136351?w=150&h=100"),
			ValidityType: models.ValidityTypeRelative,
			ValidityDays: &validityDays365,
			TotalExams:   23, // Updated to reflect actual exam assignment
			IsActive:     true,
			SortOrder:    4,
		},
		// Package 6: Free with Fixed validity
		{
			Name:         "Free Pediatrics Essentials",
			Slug:         "free-pediatrics-essentials",
			Description:  StringPtr("Essential pediatrics questions for medical students and residents"),
			PackageType:  models.PackageTypeFree,
			Price:        0.00,
			ImageURL:     StringPtr("https://images.unsplash.com/photo-1576091160550-2173dba999ef?w=800"),
			ImageAlt:     StringPtr("Pediatric medical care and child health essentials"),
			ThumbnailURL: StringPtr("https://images.unsplash.com/photo-1576091160550-2173dba999ef?w=150&h=100"),
			ValidityType: models.ValidityTypeFixed,
			ValidityDate: &fixedDate1, // Expires 6 months from now
			TotalExams:   22,          // Updated to reflect actual exam assignment
			IsActive:     true,
			SortOrder:    5,
		},
	}

	// Set image metadata for each package
	for i := range packages {
		if packages[i].ImageURL != nil {
			metadata := models.ImageMetadata{
				Width:    800,
				Height:   533,
				FileSize: 95000 + int64(i*10000), // Simulate different file sizes
				Format:   "JPEG",
			}
			packages[i].SetImageMetadata(metadata)
		}

		if err := db.DB.Create(&packages[i]).Error; err != nil {
			log.Printf("Failed to create package %s: %v", packages[i].Name, err)
			continue
		}
		log.Printf("Created package: %s (ID: %d)", packages[i].Name, packages[i].ID)
	}

	return packages
}

func (s *PackageSeeder) SeedEnrollments(db *database.Database, users []models.User, packages []models.Package, coupons []models.Coupon) {
	log.Println("Seeding user enrollments...")

	now := time.Now()

	enrollments := []models.UserPackageEnrollment{
		// User 1 (Admin) - Free package
		{
			UserID:              users[0].ID,
			PackageID:           packages[0].ID,
			EnrollmentType:      models.EnrollmentTypeFull,
			EnrolledAt:          now,
			ExpiresAt:           &[]time.Time{now.AddDate(0, 1, 0)}[0],
			IsTrialUsed:         false,
			EnrolledPackageType: packages[0].PackageType,
			EnrolledPrice:       packages[0].Price,
			PaymentStatus:       models.PaymentStatusFree,
			IsActive:            true,
		},
		// User 1 (Admin) - Premium package with coupon
		{
			UserID:              users[0].ID,
			PackageID:           packages[1].ID,
			EnrollmentType:      models.EnrollmentTypeFull,
			EnrolledAt:          now,
			ExpiresAt:           &[]time.Time{now.AddDate(0, 3, 0)}[0],
			IsTrialUsed:         false,
			EnrolledPackageType: packages[1].PackageType,
			EnrolledPrice:       packages[1].Price,
			PaymentStatus:       models.PaymentStatusPaid,
			PaymentAmount:       Float64Ptr(79.99), // After 20% discount
			PaymentDate:         &now,
			CouponID:            &coupons[0].ID,
			CouponCode:          &coupons[0].Code,
			OriginalPrice:       &packages[1].Price,
			DiscountPercentage:  &coupons[0].DiscountPercentage,
			DiscountAmount:      Float64Ptr(20.00),
			FinalPrice:          Float64Ptr(79.99),
			IsActive:            true,
		},
		// User 2 (John) - Trial enrollment
		{
			UserID:              users[1].ID,
			PackageID:           packages[1].ID,
			EnrollmentType:      models.EnrollmentTypeTrial,
			EnrolledAt:          now,
			ExpiresAt:           &[]time.Time{now.AddDate(0, 3, 0)}[0],
			IsTrialUsed:         true,
			TrialExpiresAt:      &[]time.Time{now.AddDate(0, 0, 7)}[0], // 7 days trial
			EnrolledPackageType: packages[1].PackageType,
			EnrolledPrice:       packages[1].Price,
			PaymentStatus:       models.PaymentStatusFree,
			IsActive:            true,
		},
		// User 3 (Jane) - Student bundle with student discount
		{
			UserID:              users[2].ID,
			PackageID:           packages[3].ID,
			EnrollmentType:      models.EnrollmentTypeFull,
			EnrolledAt:          now,
			ExpiresAt:           &[]time.Time{now.AddDate(1, 0, 0)}[0],
			IsTrialUsed:         false,
			EnrolledPackageType: packages[3].PackageType,
			EnrolledPrice:       packages[3].Price,
			PaymentStatus:       models.PaymentStatusPaid,
			PaymentAmount:       Float64Ptr(99.99), // After 50% discount
			PaymentDate:         &now,
			CouponID:            &coupons[1].ID,
			CouponCode:          &coupons[1].Code,
			OriginalPrice:       &packages[3].Price,
			DiscountPercentage:  &coupons[1].DiscountPercentage,
			DiscountAmount:      Float64Ptr(100.00),
			FinalPrice:          Float64Ptr(99.99),
			IsActive:            true,
		},
	}

	for i := range enrollments {
		if err := db.DB.Create(&enrollments[i]).Error; err != nil {
			log.Printf("Failed to create enrollment: %v", err)
			continue
		}
	}

	log.Printf("Created %d user enrollments", len(enrollments))

	// Create coupon usage records for paid enrollments with coupons
	s.seedCouponUsageRecords(db, enrollments, coupons)
}

func (s *PackageSeeder) seedCouponUsageRecords(db *database.Database, enrollments []models.UserPackageEnrollment, coupons []models.Coupon) {
	log.Println("Seeding coupon usage records...")

	now := time.Now()

	// Find enrollments that used coupons
	for _, enrollment := range enrollments {
		if enrollment.CouponID != nil && enrollment.PaymentStatus == models.PaymentStatusPaid {
			usage := models.CouponUsage{
				CouponID:           *enrollment.CouponID,
				UserID:             enrollment.UserID,
				EnrollmentID:       enrollment.ID,
				PackageID:          enrollment.PackageID,
				OriginalPrice:      *enrollment.OriginalPrice,
				DiscountPercentage: *enrollment.DiscountPercentage,
				DiscountAmount:     *enrollment.DiscountAmount,
				FinalPrice:         *enrollment.FinalPrice,
				CouponCode:         *enrollment.CouponCode,
				UsedAt:             now,
			}

			if err := db.DB.Create(&usage).Error; err != nil {
				log.Printf("Failed to create coupon usage record: %v", err)
				continue
			}
		}
	}

	log.Printf("Created coupon usage records")
}

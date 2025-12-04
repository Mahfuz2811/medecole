package seeders

import (
	"log"
	"quizora-backend/internal/database"
	"quizora-backend/internal/models"
)

type DatabaseSeeder struct {
	db *database.Database
}

func NewDatabaseSeeder(db *database.Database) *DatabaseSeeder {
	return &DatabaseSeeder{db: db}
}

func (s *DatabaseSeeder) Run() {
	log.Println("Starting data seeding...")

	// Clear existing data
	s.clearExistingData()

	// Initialize seeders
	userSeeder := NewUserSeeder()
	packageSeeder := NewPackageSeeder()
	examSeeder := NewExamSeeder()
	questionSeeder := NewQuestionSeeder()

	// Seed data in order of dependencies
	users := userSeeder.Seed(s.db)
	coupons := userSeeder.SeedCoupons(s.db, users[0].ID) // First user as admin

	// Create subjects, systems, and questions
	questionSeeder.Seed(s.db)

	// Create packages
	packages := packageSeeder.Seed(s.db, coupons)

	// Create exams
	exams := examSeeder.Seed(s.db, packages)

	// Create package-exam relationships
	examSeeder.SeedPackageExamRelationships(s.db, packages, exams)

	// Create user enrollments
	packageSeeder.SeedEnrollments(s.db, users, packages, coupons)

	// Create user exam attempts (realistic exam scenarios)
	attempts := examSeeder.SeedExamAttempts(s.db, users, exams)

	// Create user question answers (background analytics data)
	examSeeder.SeedQuestionAnswers(s.db, attempts, exams)

	// Print statistics
	s.printStatistics()

	log.Println("Data seeding completed successfully!")
}

func (s *DatabaseSeeder) clearExistingData() {
	log.Println("Clearing existing data...")

	// Delete in reverse order due to foreign key constraints
	s.db.DB.Exec("DELETE FROM user_question_answers")
	s.db.DB.Exec("DELETE FROM user_exam_attempts")
	s.db.DB.Exec("DELETE FROM coupon_usages")
	s.db.DB.Exec("DELETE FROM user_package_enrollments")
	s.db.DB.Exec("DELETE FROM package_exams")
	s.db.DB.Exec("DELETE FROM exams")
	s.db.DB.Exec("DELETE FROM packages")
	s.db.DB.Exec("DELETE FROM coupons")
	s.db.DB.Exec("DELETE FROM questions")
	s.db.DB.Exec("DELETE FROM systems")
	s.db.DB.Exec("DELETE FROM subjects")
	s.db.DB.Exec("DELETE FROM users")

	// Reset auto increment
	s.db.DB.Exec("ALTER TABLE user_question_answers AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE user_exam_attempts AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE coupon_usages AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE user_package_enrollments AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE package_exams AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE exams AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE packages AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE coupons AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE questions AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE systems AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE subjects AUTO_INCREMENT = 1")
	s.db.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func (s *DatabaseSeeder) printStatistics() {
	var subjectCount, systemCount, questionCount, packageCount, examCount, enrollmentCount, couponCount, attemptCount, questionAnswerCount int64

	s.db.DB.Model(&models.Subject{}).Count(&subjectCount)
	s.db.DB.Model(&models.System{}).Count(&systemCount)
	s.db.DB.Model(&models.Question{}).Count(&questionCount)
	s.db.DB.Model(&models.Package{}).Count(&packageCount)
	s.db.DB.Model(&models.Exam{}).Count(&examCount)
	s.db.DB.Model(&models.UserPackageEnrollment{}).Count(&enrollmentCount)
	s.db.DB.Model(&models.Coupon{}).Count(&couponCount)
	s.db.DB.Model(&models.UserExamAttempt{}).Count(&attemptCount)
	s.db.DB.Model(&models.UserQuestionAnswer{}).Count(&questionAnswerCount)

	log.Printf("\n=== SEEDING STATISTICS ===")
	log.Printf("Subjects created: %d", subjectCount)
	log.Printf("Systems created: %d", systemCount)
	log.Printf("Questions created: %d", questionCount)
	log.Printf("Packages created: %d", packageCount)
	log.Printf("Exams created: %d", examCount)
	log.Printf("Enrollments created: %d", enrollmentCount)
	log.Printf("Coupons created: %d", couponCount)
	log.Printf("Exam attempts created: %d", attemptCount)
	log.Printf("Question answers created: %d", questionAnswerCount)
	log.Printf("========================")
}

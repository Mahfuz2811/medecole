package seeders

import (
	"encoding/json"
	"fmt"
	"log"
	"quizora-backend/internal/database"
	"quizora-backend/internal/models"
	"strings"
	"time"
)

type ExamSeeder struct{}

func NewExamSeeder() *ExamSeeder {
	return &ExamSeeder{}
}

func (s *ExamSeeder) Seed(db *database.Database, packages []models.Package) []models.Exam {
	log.Println("Seeding exams...")

	now := time.Now()

	// Sample questions data for embedding in exams
	sampleQuestions := []map[string]interface{}{
		{
			"id":            1,
			"question_text": "A 45-year-old man presents with chest pain. What is the most likely diagnosis?",
			"question_type": "SBA",
			"options": map[string]interface{}{
				"a": map[string]interface{}{"text": "Myocardial infarction", "is_correct": true},
				"b": map[string]interface{}{"text": "Pneumonia", "is_correct": false},
				"c": map[string]interface{}{"text": "Anxiety", "is_correct": false},
				"d": map[string]interface{}{"text": "GERD", "is_correct": false},
			},
			"explanation":      "Classic presentation of MI",
			"difficulty_level": "MEDIUM",
		},
		{
			"id":            2,
			"question_text": "True or False: Hypertension is a major risk factor for stroke",
			"question_type": "TRUE_FALSE",
			"options": map[string]interface{}{
				"a": map[string]interface{}{"text": "True", "is_correct": true},
				"b": map[string]interface{}{"text": "False", "is_correct": false},
			},
			"explanation":      "Hypertension is indeed a major modifiable risk factor for stroke",
			"difficulty_level": "EASY",
		},
	}

	questionsJSON, _ := json.Marshal(sampleQuestions)

	// Generate 100+ exams with different combinations
	exams := s.generateExams(questionsJSON, now)

	for i := range exams {
		if err := db.DB.Create(&exams[i]).Error; err != nil {
			log.Printf("Failed to create exam %s: %v", exams[i].Title, err)
			continue
		}
		log.Printf("Created exam: %s (ID: %d)", exams[i].Title, exams[i].ID)
	}

	return exams
}

func (s *ExamSeeder) generateExams(questionsJSON []byte, now time.Time) []models.Exam {
	var exams []models.Exam

	// Medical subjects and topics for realistic exam names
	subjects := []string{
		"Cardiology", "Respiratory", "Gastroenterology", "Neurology", "Endocrinology",
		"Nephrology", "Hematology", "Oncology", "Infectious Disease", "Rheumatology",
		"Dermatology", "Psychiatry", "Pediatrics", "Surgery", "Orthopedics",
		"Ophthalmology", "ENT", "Obstetrics", "Gynecology", "Anesthesia",
		"Emergency Medicine", "Radiology", "Pathology", "Pharmacology", "Anatomy",
	}

	examTypes := []models.ExamType{
		models.ExamTypeDaily, models.ExamTypeMock, models.ExamTypeReview, models.ExamTypeFinal,
	}

	examStatuses := []models.ExamStatus{
		models.ExamStatusDraft, models.ExamStatusScheduled, models.ExamStatusActive, models.ExamStatusCompleted,
	}

	// Generate exams with different combinations
	examCounter := 0
	for _, subject := range subjects {
		for _, examType := range examTypes {
			for _, status := range examStatuses {
				examCounter++
				if examCounter > 120 { // Generate 120 exams total
					return exams
				}

				exam := s.createExamVariant(subject, examType, status, examCounter, questionsJSON, now)
				exams = append(exams, exam)
			}
		}
	}

	return exams
}

func (s *ExamSeeder) createExamVariant(subject string, examType models.ExamType, status models.ExamStatus, counter int, questionsJSON []byte, now time.Time) models.Exam {
	// Generate realistic exam configurations based on type
	var totalQuestions int
	var durationMinutes int
	var passingScore float64
	var maxAttempts int

	switch examType {
	case models.ExamTypeDaily:
		totalQuestions = 10 + (counter % 10)      // 10-19 questions
		durationMinutes = 15 + (counter % 15)     // 15-29 minutes
		passingScore = 60.0 + float64(counter%20) // 60-79%
		maxAttempts = 3
	case models.ExamTypeMock:
		totalQuestions = 50 + (counter % 50)      // 50-99 questions
		durationMinutes = 120 + (counter % 60)    // 120-179 minutes
		passingScore = 70.0 + float64(counter%15) // 70-84%
		maxAttempts = 2
	case models.ExamTypeReview:
		totalQuestions = 20 + (counter % 30)      // 20-49 questions
		durationMinutes = 45 + (counter % 45)     // 45-89 minutes
		passingScore = 65.0 + float64(counter%20) // 65-84%
		maxAttempts = 5
	case models.ExamTypeFinal:
		totalQuestions = 100 + (counter % 100)    // 100-199 questions
		durationMinutes = 180 + (counter % 120)   // 180-299 minutes
		passingScore = 75.0 + float64(counter%15) // 75-89%
		maxAttempts = 1
	}

	// Generate realistic scheduling based on status
	var scheduledStartDate *time.Time
	var scheduledEndDate *time.Time

	switch status {
	case models.ExamStatusDraft:
		// No scheduling for drafts
	case models.ExamStatusScheduled:
		// Future dates
		futureStart := now.AddDate(0, 0, 1+(counter%30))       // 1-30 days from now
		futureEnd := futureStart.AddDate(0, 0, 7+(counter%14)) // 7-20 days after start
		scheduledStartDate = &futureStart
		scheduledEndDate = &futureEnd
	case models.ExamStatusActive:
		// Current or recent past dates
		activeStart := now.AddDate(0, 0, -(counter % 7)) // Up to 7 days ago
		activeEnd := now.AddDate(0, 0, 7+(counter%7))    // Up to 14 days from now
		scheduledStartDate = &activeStart
		scheduledEndDate = &activeEnd
	case models.ExamStatusCompleted:
		// Past dates
		pastStart := now.AddDate(0, 0, -30-(counter%60))  // 30-89 days ago
		pastEnd := pastStart.AddDate(0, 0, 7+(counter%7)) // 7-13 days after start
		scheduledStartDate = &pastStart
		scheduledEndDate = &pastEnd
	}

	// Generate slug
	slug := fmt.Sprintf("%s-%s-%s-%d",
		strings.ToLower(strings.ReplaceAll(subject, " ", "-")),
		strings.ToLower(string(examType)),
		strings.ToLower(string(status)),
		counter)

	// Generate realistic analytics based on status
	var attemptCount int
	var completedAttemptCount int
	var averageScore *float64
	var passRate *float64
	var lastAttemptAt *time.Time

	if status == models.ExamStatusActive || status == models.ExamStatusCompleted {
		attemptCount = 10 + (counter % 50)                                                      // 10-59 attempts
		completedAttemptCount = int(float64(attemptCount) * (0.6 + 0.3*float64(counter%10)/10)) // 60-90% completion rate

		avgScore := passingScore + float64(counter%20) - 10 // Vary around passing score
		if avgScore < 0 {
			avgScore = passingScore * 0.8
		}
		averageScore = &avgScore

		passRateVal := float64(completedAttemptCount) * (0.4 + 0.4*float64(counter%10)/10) / float64(completedAttemptCount) * 100
		passRate = &passRateVal

		lastAttempt := now.Add(-time.Hour * time.Duration(counter%72)) // Within last 3 days
		lastAttemptAt = &lastAttempt
	}

	return models.Exam{
		Title:       fmt.Sprintf("%s %s Exam", subject, s.capitalizeFirst(string(examType))),
		Slug:        slug,
		Description: StringPtr(fmt.Sprintf("Comprehensive %s examination focusing on %s concepts and clinical applications", strings.ToLower(string(examType)), strings.ToLower(subject))),

		ExamType:        examType,
		TotalQuestions:  totalQuestions,
		DurationMinutes: durationMinutes,
		PassingScore:    passingScore,
		MaxAttempts:     maxAttempts,

		QuestionsData: string(questionsJSON),

		ScheduledStartDate: scheduledStartDate,
		ScheduledEndDate:   scheduledEndDate,

		Instructions: StringPtr(fmt.Sprintf("This is a %s %s exam. Read each question carefully and select the best answer. You have %d minutes to complete %d questions.", strings.ToLower(string(examType)), strings.ToLower(subject), durationMinutes, totalQuestions)),

		AttemptCount:          attemptCount,
		CompletedAttemptCount: completedAttemptCount,
		AverageScore:          averageScore,
		PassRate:              passRate,
		LastAttemptAt:         lastAttemptAt,

		Status:   status,
		IsActive: status != models.ExamStatusDraft,
	}
}

func (s *ExamSeeder) capitalizeFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	lower := strings.ToLower(str)
	return strings.ToUpper(string(lower[0])) + lower[1:]
}

func (s *ExamSeeder) SeedPackageExamRelationships(db *database.Database, packages []models.Package, exams []models.Exam) {
	log.Println("Seeding package-exam relationships...")

	var relationships []models.PackageExam

	// Distribute exams across packages based on package type and exam characteristics
	for i, pkg := range packages {
		var selectedExams []models.Exam

		switch i {
		case 0: // FCPS Part 1 Medicine (Premium Fixed) - Advanced medical exams
			selectedExams = s.selectExamsByType(exams, []models.ExamType{models.ExamTypeMock, models.ExamTypeFinal}, []models.ExamStatus{models.ExamStatusActive, models.ExamStatusScheduled}, 25)
		case 1: // Free Medicine Basics (Free Relative) - Basic daily exams
			selectedExams = s.selectExamsByType(exams, []models.ExamType{models.ExamTypeDaily, models.ExamTypeReview}, []models.ExamStatus{models.ExamStatusActive}, 15)
		case 2: // Premium Medicine Complete (Premium Fixed) - Comprehensive mix
			selectedExams = s.selectExamsByType(exams, []models.ExamType{models.ExamTypeMock, models.ExamTypeReview, models.ExamTypeFinal}, []models.ExamStatus{models.ExamStatusActive, models.ExamStatusScheduled}, 30)
		case 3: // Surgery Mastery (Premium Fixed) - Surgery focused
			selectedExams = s.selectSurgeryExams(exams, 20)
		case 4: // Student Bundle (Free Relative) - Mix of all types
			selectedExams = s.selectExamsByType(exams, []models.ExamType{models.ExamTypeDaily, models.ExamTypeReview, models.ExamTypeMock}, []models.ExamStatus{models.ExamStatusActive, models.ExamStatusCompleted}, 35)
		case 5: // Free Pediatrics Essentials (Free Fixed) - Pediatrics focused
			selectedExams = s.selectPediatricsExams(exams, 10)
		}

		// Create relationships for selected exams
		for j, exam := range selectedExams {
			relationships = append(relationships, models.PackageExam{
				PackageID: pkg.ID,
				ExamID:    exam.ID,
				SortOrder: j + 1,
			})
		}
	}

	// Create the relationships in database
	for _, rel := range relationships {
		if err := db.DB.Create(&rel).Error; err != nil {
			log.Printf("Failed to create package-exam relationship: %v", err)
			continue
		}
	}

	log.Printf("Created %d package-exam relationships", len(relationships))
}

// Helper methods for exam selection
func (s *ExamSeeder) selectExamsByType(exams []models.Exam, examTypes []models.ExamType, statuses []models.ExamStatus, limit int) []models.Exam {
	var selected []models.Exam
	count := 0

	for _, exam := range exams {
		if count >= limit {
			break
		}

		// Check if exam type matches
		typeMatch := false
		for _, et := range examTypes {
			if exam.ExamType == et {
				typeMatch = true
				break
			}
		}

		// Check if status matches
		statusMatch := false
		for _, st := range statuses {
			if exam.Status == st {
				statusMatch = true
				break
			}
		}

		if typeMatch && statusMatch {
			selected = append(selected, exam)
			count++
		}
	}

	return selected
}

func (s *ExamSeeder) selectSurgeryExams(exams []models.Exam, limit int) []models.Exam {
	var selected []models.Exam
	count := 0

	for _, exam := range exams {
		if count >= limit {
			break
		}

		// Look for surgery-related subjects in the title
		title := strings.ToLower(exam.Title)
		if strings.Contains(title, "surgery") || strings.Contains(title, "orthopedics") || strings.Contains(title, "anesthesia") {
			selected = append(selected, exam)
			count++
		}
	}

	// If not enough surgery-specific exams, fill with other exams
	if count < limit {
		for _, exam := range exams {
			if count >= limit {
				break
			}

			// Skip already selected exams
			alreadySelected := false
			for _, sel := range selected {
				if sel.ID == exam.ID {
					alreadySelected = true
					break
				}
			}

			if !alreadySelected && exam.Status == models.ExamStatusActive {
				selected = append(selected, exam)
				count++
			}
		}
	}

	return selected
}

func (s *ExamSeeder) selectPediatricsExams(exams []models.Exam, limit int) []models.Exam {
	var selected []models.Exam
	count := 0

	for _, exam := range exams {
		if count >= limit {
			break
		}

		// Look for pediatrics-related subjects in the title
		title := strings.ToLower(exam.Title)
		if strings.Contains(title, "pediatrics") || strings.Contains(title, "pediatric") {
			selected = append(selected, exam)
			count++
		}
	}

	// If not enough pediatrics-specific exams, fill with basic exams
	if count < limit {
		for _, exam := range exams {
			if count >= limit {
				break
			}

			// Skip already selected exams
			alreadySelected := false
			for _, sel := range selected {
				if sel.ID == exam.ID {
					alreadySelected = true
					break
				}
			}

			if !alreadySelected && (exam.ExamType == models.ExamTypeDaily || exam.ExamType == models.ExamTypeReview) && exam.Status == models.ExamStatusActive {
				selected = append(selected, exam)
				count++
			}
		}
	}

	return selected
}

func (s *ExamSeeder) SeedExamAttempts(db *database.Database, users []models.User, exams []models.Exam) []models.UserExamAttempt {
	log.Println("Seeding exam attempts...")

	now := time.Now()

	// Sample answers data (JSON format)
	completedAnswers := `[
		{
			"question_id": 1,
			"question_index": 0,
			"selected_options": ["a"],
			"time_spent": 45,
			"answered_at": "2024-01-15T10:15:30Z",
			"change_count": 2
		},
		{
			"question_id": 2,
			"question_index": 1,
			"selected_options": ["a"],
			"time_spent": 30,
			"answered_at": "2024-01-15T10:16:00Z",
			"change_count": 1
		}
	]`

	inProgressAnswers := `[
		{
			"question_id": 1,
			"question_index": 0,
			"selected_options": ["b"],
			"time_spent": 60,
			"answered_at": "2024-01-15T11:15:30Z",
			"change_count": 1
		}
	]`

	attempts := []models.UserExamAttempt{
		// User 1 (Admin) - Completed exam with good score
		{
			UserID:           users[0].ID,
			ExamID:           exams[0].ID, // Basic Cardiology Assessment
			PackageID:        1,           // Free Medicine Basics package
			AttemptNumber:    1,
			Status:           models.AttemptStatusCompleted,
			StartedAt:        now.Add(-2 * time.Hour),                  // Started 2 hours ago
			CompletedAt:      &[]time.Time{now.Add(-1 * time.Hour)}[0], // Completed 1 hour ago
			TimeLimitSeconds: 3600,                                     // 60 minutes = 3600 seconds
			ActualTimeSpent:  3420,                                     // 57 minutes (efficient)
			AnswersData:      completedAnswers,
			TotalQuestions:   20,
			PassingScore:     70.0,
			IsScored:         true,
			Score:            Float64Ptr(85.5),
			CorrectAnswers:   IntPtr(17),
			IsPassed:         BoolPtr(true),
		},
		// User 1 (Admin) - Auto-submitted exam (time expired)
		{
			UserID:           users[0].ID,
			ExamID:           exams[1].ID, // Respiratory System Quiz
			PackageID:        1,           // Free Medicine Basics package
			AttemptNumber:    1,
			Status:           models.AttemptStatusAutoSubmitted,
			StartedAt:        now.Add(-4 * time.Hour),                                 // Started 4 hours ago
			CompletedAt:      &[]time.Time{now.Add(-2*time.Hour + 30*time.Minute)}[0], // Auto-submitted when time expired
			TimeLimitSeconds: 5400,                                                    // 90 minutes = 5400 seconds
			ActualTimeSpent:  5400,                                                    // Full time used (auto-submitted)
			AnswersData:      completedAnswers,
			TotalQuestions:   30,
			PassingScore:     75.0,
			IsScored:         true,
			Score:            Float64Ptr(72.3),
			CorrectAnswers:   IntPtr(22),
			IsPassed:         BoolPtr(false), // Below passing score
		},
		// User 2 (John) - Currently in progress
		{
			UserID:           users[1].ID,
			ExamID:           exams[0].ID, // Basic Cardiology Assessment
			PackageID:        1,           // Free Medicine Basics package
			AttemptNumber:    1,
			Status:           models.AttemptStatusStarted,
			StartedAt:        now.Add(-30 * time.Minute), // Started 30 minutes ago
			CompletedAt:      nil,                        // Still in progress
			TimeLimitSeconds: 3600,                       // 60 minutes = 3600 seconds
			ActualTimeSpent:  0,                          // Will be calculated on completion
			AnswersData:      inProgressAnswers,          // Partial answers from Redis
			TotalQuestions:   20,
			PassingScore:     70.0,
			IsScored:         false,
			Score:            nil,
			CorrectAnswers:   nil,
			IsPassed:         nil,
		},
		// User 3 (Jane) - Completed surgery exam with excellent score
		{
			UserID:           users[2].ID,
			ExamID:           exams[2].ID, // Surgery Fundamentals
			PackageID:        1,           // Free Medicine Basics package
			AttemptNumber:    1,
			Status:           models.AttemptStatusCompleted,
			StartedAt:        now.Add(-6 * time.Hour),                  // Started 6 hours ago
			CompletedAt:      &[]time.Time{now.Add(-4 * time.Hour)}[0], // Completed 4 hours ago
			TimeLimitSeconds: 7200,                                     // 120 minutes = 7200 seconds
			ActualTimeSpent:  6900,                                     // 115 minutes (used most of the time)
			AnswersData:      completedAnswers,
			TotalQuestions:   50,
			PassingScore:     80.0,
			IsScored:         true,
			Score:            Float64Ptr(92.0),
			CorrectAnswers:   IntPtr(46),
			IsPassed:         BoolPtr(true),
		},
		// User 2 (John) - Failed attempt (abandoned early)
		{
			UserID:           users[1].ID,
			ExamID:           exams[3].ID, // Advanced Medicine Challenge
			PackageID:        1,           // Free Medicine Basics package
			AttemptNumber:    1,
			Status:           models.AttemptStatusAbandoned,
			StartedAt:        now.Add(-24 * time.Hour),                  // Started yesterday
			CompletedAt:      &[]time.Time{now.Add(-23 * time.Hour)}[0], // Abandoned after 1 hour
			TimeLimitSeconds: 10800,                                     // 180 minutes = 10800 seconds
			ActualTimeSpent:  3600,                                      // Only 1 hour spent
			AnswersData:      inProgressAnswers,                         // Minimal answers
			TotalQuestions:   100,
			PassingScore:     85.0,
			IsScored:         true,
			Score:            Float64Ptr(15.0), // Very low score
			CorrectAnswers:   IntPtr(15),
			IsPassed:         BoolPtr(false),
		},
	}

	for i := range attempts {
		if err := db.DB.Create(&attempts[i]).Error; err != nil {
			log.Printf("Failed to create exam attempt: %v", err)
			continue
		}
	}

	log.Printf("Created %d exam attempts", len(attempts))
	return attempts
}

func (s *ExamSeeder) SeedQuestionAnswers(db *database.Database, attempts []models.UserExamAttempt, exams []models.Exam) {
	log.Println("Seeding question answers...")

	now := time.Now()

	// Create question answers only for completed/scored attempts
	for _, attempt := range attempts {
		if !attempt.IsScored {
			continue // Skip in-progress attempts
		}

		// Create sample question answers for this attempt
		questionAnswers := s.generateQuestionAnswersForAttempt(attempt, now)

		for _, qa := range questionAnswers {
			if err := db.DB.Create(&qa).Error; err != nil {
				log.Printf("Failed to create question answer: %v", err)
				continue
			}
		}

		log.Printf("Created %d question answers for attempt ID: %d", len(questionAnswers), attempt.ID)
	}
}

func (s *ExamSeeder) generateQuestionAnswersForAttempt(attempt models.UserExamAttempt, baseTime time.Time) []models.UserQuestionAnswer {
	answers := []models.UserQuestionAnswer{}

	// Generate realistic question answers based on the attempt's performance
	questionsToGenerate := 5 // Generate 5 sample questions per attempt for demonstration

	for i := 0; i < questionsToGenerate; i++ {
		// Simulate realistic answer patterns based on attempt score
		isCorrect := false
		partialScore := 0.0
		if attempt.Score != nil {
			// Higher scoring attempts have more correct answers
			correctnessThreshold := *attempt.Score / 100.0
			// Add some randomness but bias toward the overall score
			if float64(i) < correctnessThreshold*float64(questionsToGenerate) {
				isCorrect = true
				partialScore = 1.0
			} else {
				// Some partial credit for wrong answers
				partialScore = 0.1
			}
		}

		// Simulate time spent (varies by question)
		timeSpent := 30 + (i * 15) // 30-90 seconds per question

		// Simulate change count (more changes for difficult questions)
		changeCount := 0
		if i > questionsToGenerate/2 {
			changeCount = 1 + (i % 3) // More changes for later questions
		}

		qa := models.UserQuestionAnswer{
			AttemptID:       attempt.ID,
			UserID:          attempt.UserID,
			ExamID:          attempt.ExamID,
			QuestionID:      uint(i + 1), // Sequential question IDs
			QuestionType:    models.QuestionTypeSBA,
			QuestionText:    fmt.Sprintf("Sample question %d for %s", i+1, s.getExamTypeFromID(attempt.ExamID)),
			DifficultyLevel: s.getDifficultyForQuestionIndex(i),
			QuestionIndex:   i,
			SelectedOptions: fmt.Sprintf(`["%c"]`, 'a'+i%4), // Rotate through a,b,c,d
			CorrectOptions:  `["a"]`,                        // Assume 'a' is correct for simplicity
			IsCorrect:       isCorrect,
			PartialScore:    partialScore,
			MaxScore:        1.0,
			TimeSpent:       timeSpent,
			AnsweredAt:      baseTime.Add(time.Duration(i*60) * time.Second), // 1 minute apart
			IsSkipped:       false,
			ChangeCount:     changeCount,
			IsLastAnswer:    true,
		}

		answers = append(answers, qa)
	}

	return answers
}

func (s *ExamSeeder) getExamTypeFromID(examID uint) string {
	examTypes := map[uint]string{
		1: "Cardiology",
		2: "Respiratory",
		3: "Surgery",
		4: "Advanced Medicine",
	}
	if examType, exists := examTypes[examID]; exists {
		return examType
	}
	return "General"
}

func (s *ExamSeeder) getDifficultyForQuestionIndex(index int) models.DifficultyLevel {
	// First few questions are easy, middle are medium, last are hard
	if index < 2 {
		return models.DifficultyEasy
	} else if index < 4 {
		return models.DifficultyMedium
	}
	return models.DifficultyHard
}

package service

import (
	"encoding/json"
	"fmt"
	"quizora-backend/internal/dto"
	"quizora-backend/internal/logger"
	"quizora-backend/internal/mapper"
	"quizora-backend/internal/models"
	"quizora-backend/internal/repository"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ExamService handles business logic for exam operations
type ExamService interface {
	GetPackageExamsBySlug(packageSlug string, userID uint) (dto.ExamListResponse, error)
	GetExamMetaBySlug(examSlug string) (dto.ExamMetaResponse, error)
	StartExam(packageSlug string, examSlug string, userID uint, deviceInfo map[string]string) (dto.StartExamResponse, error)
	GetSession(sessionID string, userID uint) (dto.ExamSessionResponse, error)
	SyncSession(sessionID string, userID uint, answers []dto.UserAnswerSync) (dto.SyncSessionResponse, error)
	SubmitExam(sessionID string, userID uint) (dto.SubmitExamResponse, error)
	GetExamBySlug(examSlug string) (*models.Exam, error)
	GetUserAttemptsByExam(userID uint, examID uint) ([]models.UserExamAttempt, error)
	GetExamResultsBySession(sessionID string, userID uint) (interface{}, error)
}

// examService implements ExamService
type examService struct {
	examRepo       repository.ExamRepository
	enrollmentRepo repository.EnrollmentRepository
	examMapper     mapper.ExamMapper
}

// NewExamService creates a new exam service
func NewExamService(examRepo repository.ExamRepository, enrollmentRepo repository.EnrollmentRepository, examMapper mapper.ExamMapper) ExamService {
	return &examService{
		examRepo:       examRepo,
		enrollmentRepo: enrollmentRepo,
		examMapper:     examMapper,
	}
}

// GetPackageExamsBySlug retrieves package data along with all exams for a specific package by slug
func (s *examService) GetPackageExamsBySlug(packageSlug string, userID uint) (dto.ExamListResponse, error) {
	// Get package with exams from repository
	packageWithExams, err := s.examRepo.GetPackageWithExamsBySlug(packageSlug, userID)
	if err != nil {
		return dto.ExamListResponse{}, err
	}

	// Convert to response format including package data
	response := s.examMapper.ToExamListResponseWithPackage(packageWithExams.Package, packageWithExams.Exams)

	return response, nil
}

// GetExamMetaBySlug retrieves only the exam metadata without questions or user data
func (s *examService) GetExamMetaBySlug(examSlug string) (dto.ExamMetaResponse, error) {
	// Get exam details from repository
	exam, err := s.examRepo.GetExamBySlug(examSlug)
	if err != nil {
		return dto.ExamMetaResponse{}, err
	}

	// Convert to response format
	response := s.examMapper.ToExamMetaResponse(*exam)

	return response, nil
}

// StartExam initializes a new exam session for a user
// Phase 1 & 2: Optimized with reduced DB calls and enrollment validation
func (s *examService) StartExam(packageSlug string, examSlug string, userID uint, deviceInfo map[string]string) (dto.StartExamResponse, error) {
	// 1. Get exam details (1 DB call)
	exam, err := s.examRepo.GetExamBySlug(examSlug)
	if err != nil {
		return dto.StartExamResponse{}, err
	}

	// 2. Get package details by slug to get package ID (1 DB call)
	packageData, err := s.examRepo.GetPackageWithExamsBySlug(packageSlug, userID)
	if err != nil {
		return dto.StartExamResponse{}, fmt.Errorf("failed to get package: %w", err)
	}
	packageID := packageData.Package.ID

	// 3. Validate user enrollment (1 DB call)
	enrolled, err := s.enrollmentRepo.IsUserEnrolledInPackage(userID, packageID)
	if err != nil {
		return dto.StartExamResponse{}, fmt.Errorf("failed to check enrollment: %w", err)
	}
	if !enrolled {
		return dto.StartExamResponse{}, repository.ErrNotEnrolledInPackage
	}

	// 4. Verify exam belongs to the specified package
	examBelongsToPackage := false
	for _, examData := range packageData.Exams {
		if examData.ID == exam.ID {
			examBelongsToPackage = true
			break
		}
	}
	if !examBelongsToPackage {
		return dto.StartExamResponse{}, fmt.Errorf("exam does not belong to the specified package")
	}

	// 5. Check for existing attempt in THIS package context (1 DB call)
	existingAttempt, err := s.examRepo.GetUserAttemptForExamInPackage(userID, exam.ID, packageID)
	if err != nil {
		return dto.StartExamResponse{}, err
	}

	if existingAttempt != nil {
		// Check if already completed (single attempt rule per package)
		if existingAttempt.IsCompleted() {
			return dto.StartExamResponse{}, repository.ErrExamAlreadySubmitted
		}
		// Return active session
		return dto.StartExamResponse{
			SessionID: existingAttempt.GetSessionKey(),
			AttemptID: existingAttempt.ID,
			ExamMeta:  s.examMapper.ToExamMetaResponse(*exam),
		}, nil
	}

	// 6. Create new attempt - reuse exam data (1 DB call)
	attempt, err := s.examRepo.CreateExamAttemptWithExam(userID, exam, packageID, deviceInfo)
	if err != nil {
		return dto.StartExamResponse{}, err
	}

	// Convert to response
	response := dto.StartExamResponse{
		SessionID: attempt.GetSessionKey(),
		AttemptID: attempt.ID,
		ExamMeta:  s.examMapper.ToExamMetaResponse(*exam),
	}

	return response, nil
}

// GetSession retrieves exam session data including exam content and session state
func (s *examService) GetSession(sessionID string, userID uint) (dto.ExamSessionResponse, error) {
	// Initialize logger with service context
	log := logger.WithService("ExamService").WithFields(logrus.Fields{
		"operation":  "GetSession",
		"session_id": sessionID,
		"user_id":    userID,
	})

	// Get active session data from repository (STARTED status only)
	sessionData, err := s.examRepo.GetActiveSessionByID(sessionID)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve session data from repository")
		return dto.ExamSessionResponse{}, err
	}

	// Verify the session belongs to the requesting user
	if sessionData.Attempt.UserID != userID {
		log.WithFields(logrus.Fields{
			"session_owner_id":   sessionData.Attempt.UserID,
			"requesting_user_id": userID,
		}).Warn("Session ownership verification failed - user attempting to access another user's session")
		return dto.ExamSessionResponse{}, fmt.Errorf("session does not belong to user")
	}

	// Check if session is still valid (not expired)
	if sessionData.Attempt.IsTimeExpired() {
		timeRemaining := sessionData.Attempt.GetRemainingTimeSeconds()
		log.WithFields(logrus.Fields{
			"time_remaining_seconds": timeRemaining,
			"is_expired":             timeRemaining <= 0,
			"started_at":             sessionData.Attempt.StartedAt.Format(time.RFC3339),
			"duration_minutes":       sessionData.Exam.DurationMinutes,
		}).Warn("Session has expired - time limit exceeded")
		return dto.ExamSessionResponse{}, fmt.Errorf("session has expired")
	}

	// Retrieve saved answers from cache
	savedAnswers, err := s.examRepo.GetSessionAnswers(sessionID)
	if err != nil {
		// Log warning but continue - session restoration is not critical
		// User can still take the exam, just won't have previous answers restored
		log.WithError(err).Warn("Failed to retrieve saved answers from cache - continuing without restored answers")
		savedAnswers = make(map[uint]string)
	}

	// Convert saved answers to response format
	var savedAnswersResponse []dto.UserAnswerResponse
	for questionID, selectedOption := range savedAnswers {
		savedAnswersResponse = append(savedAnswersResponse, dto.UserAnswerResponse{
			QuestionID:     questionID,
			SelectedOption: selectedOption,
		})
	}

	// Convert to response format with saved answers
	response := s.examMapper.ToExamSessionResponseWithAnswers(sessionData.Attempt, sessionData.Exam, savedAnswersResponse)

	return response, nil
}

// SyncSession syncs user answers for an active exam session
func (s *examService) SyncSession(sessionID string, userID uint, answers []dto.UserAnswerSync) (dto.SyncSessionResponse, error) {
	// Get active session data from repository to validate session and ownership (STARTED status only)
	sessionData, err := s.examRepo.GetActiveSessionByID(sessionID)
	if err != nil {
		return dto.SyncSessionResponse{}, err
	}

	// Verify the session belongs to the requesting user
	if sessionData.Attempt.UserID != userID {
		return dto.SyncSessionResponse{}, fmt.Errorf("session does not belong to user")
	}

	// Check if session is still valid (not expired)
	if sessionData.Attempt.IsTimeExpired() {
		return dto.SyncSessionResponse{}, fmt.Errorf("session has expired")
	}

	// Check if session is still in progress
	if !sessionData.Attempt.IsInProgress() {
		return dto.SyncSessionResponse{}, fmt.Errorf("session is not in progress")
	}

	// Convert answers to map format for repository
	answersMap := make(map[uint]string)
	for _, answer := range answers {
		answersMap[answer.QuestionID] = answer.SelectedOption
	}

	// Sync answers to cache
	err = s.examRepo.SyncSessionAnswers(sessionID, answersMap)
	if err != nil {
		return dto.SyncSessionResponse{}, fmt.Errorf("failed to sync answers: %w", err)
	}

	// Create response (time remaining calculated from original session start time)
	response := dto.SyncSessionResponse{
		Success:       true,
		SyncedCount:   len(answers),
		LastSyncAt:    sessionData.Attempt.StartedAt.Format("2006-01-02T15:04:05Z"), // Use session start time for now
		TimeRemaining: sessionData.Attempt.GetRemainingTimeSeconds(),
	}

	return response, nil
}

// GetExamBySlug retrieves exam details by slug
func (s *examService) GetExamBySlug(examSlug string) (*models.Exam, error) {
	return s.examRepo.GetExamBySlug(examSlug)
}

// GetUserAttemptsByExam retrieves all user attempts for a specific exam
func (s *examService) GetUserAttemptsByExam(userID uint, examID uint) ([]models.UserExamAttempt, error) {
	return s.examRepo.GetUserAttemptsByExam(userID, examID)
}

// SubmitExam finalizes an exam session and calculates results
func (s *examService) SubmitExam(sessionID string, userID uint) (dto.SubmitExamResponse, error) {
	// Get active session data from repository to validate session and ownership (STARTED status only)
	sessionData, err := s.examRepo.GetActiveSessionByID(sessionID)
	if err != nil {
		return dto.SubmitExamResponse{}, err
	}

	// Verify the session belongs to the requesting user
	if sessionData.Attempt.UserID != userID {
		return dto.SubmitExamResponse{}, fmt.Errorf("session does not belong to user")
	}

	// Check if session is already completed/submitted
	if sessionData.Attempt.IsCompleted() {
		return dto.SubmitExamResponse{}, fmt.Errorf("exam already submitted")
	}

	// Retrieve saved answers from cache
	savedAnswers, err := s.examRepo.GetSessionAnswers(sessionID)
	if err != nil {
		// Continue even if answers not found - user might submit with no answers
		savedAnswers = make(map[uint]string)
	}

	// Parse questions from the exam's JSON data
	var examQuestions []map[string]interface{} // Change to handle raw JSON with explanation
	err = json.Unmarshal([]byte(sessionData.Exam.QuestionsData), &examQuestions)
	if err != nil {
		return dto.SubmitExamResponse{}, fmt.Errorf("failed to parse exam questions: %w", err)
	}

	// Calculate score and build detailed answer data
	totalQuestions := len(examQuestions)
	correctAnswers := 0
	submissionTime := time.Now()

	// Structure to store detailed answers for AnswersData field
	type OptionDetail struct {
		Key       string `json:"key"`
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	}

	type AnswerDetail struct {
		QuestionID    uint           `json:"question_id"`
		QuestionText  string         `json:"question_text"`
		QuestionType  string         `json:"question_type"`
		UserAnswer    []string       `json:"user_answer"`
		CorrectAnswer interface{}    `json:"correct_answer"`
		IsCorrect     bool           `json:"is_correct"`
		PointsEarned  float64        `json:"points_earned"`
		MaxPoints     int            `json:"max_points"`
		Explanation   string         `json:"explanation"`
		Options       []OptionDetail `json:"options"`
	}

	type AnswersDataStructure struct {
		SubmissionTimestamp string         `json:"submission_timestamp"`
		Answers             []AnswerDetail `json:"answers"`
		ExamSnapshot        struct {
			TotalQuestions  int     `json:"total_questions"`
			PassingScore    float64 `json:"passing_score"`
			DurationMinutes int     `json:"duration_minutes"`
		} `json:"exam_snapshot"`
	}

	answersData := AnswersDataStructure{
		SubmissionTimestamp: submissionTime.Format("2006-01-02T15:04:05Z"),
		Answers:             make([]AnswerDetail, 0),
	}

	// Set exam snapshot
	answersData.ExamSnapshot.TotalQuestions = totalQuestions
	answersData.ExamSnapshot.PassingScore = sessionData.Exam.PassingScore
	answersData.ExamSnapshot.DurationMinutes = sessionData.Exam.DurationMinutes

	// Process each question for scoring and detailed storage
	var totalPointsEarned float64 = 0
	var totalMaxPoints float64 = 0

	for _, questionData := range examQuestions {
		// Extract question fields from map
		questionID := uint(questionData["id"].(float64))
		questionText := questionData["question_text"].(string)
		questionType := questionData["question_type"].(string)
		options := questionData["options"].(map[string]interface{})
		points := int(questionData["points"].(float64))
		explanation := ""
		if exp, exists := questionData["explanation"]; exists && exp != nil {
			explanation = exp.(string)
		}

		// Build options array for answer detail
		optionDetails := make([]OptionDetail, 0, len(options))
		for key, value := range options {
			if optData, ok := value.(map[string]interface{}); ok {
				text := ""
				if t, ok := optData["text"].(string); ok {
					text = t
				}
				isCorrect := false
				if c, ok := optData["is_correct"].(bool); ok {
					isCorrect = c
				}
				optionDetails = append(optionDetails, OptionDetail{
					Key:       key,
					Text:      text,
					IsCorrect: isCorrect,
				})
			}
		}

		userAnswer, exists := savedAnswers[questionID]

		// Parse user answer (handle both single and multiple selections)
		var userAnswerArray []string
		if exists && userAnswer != "" {
			// Try to parse as JSON array first (for TRUE_FALSE multiple selections)
			var parsedArray []string
			if err := json.Unmarshal([]byte(userAnswer), &parsedArray); err == nil {
				userAnswerArray = parsedArray
			} else {
				// Single answer (for SBA questions)
				userAnswerArray = []string{userAnswer}
			}
		}

		// Find correct answers and check if user answer is correct
		var correctAnswerArray []string
		var correctAnswersMap map[string]bool // for TRUE_FALSE questions
		isCorrect := false
		pointsEarned := 0.0
		maxPoints := float64(points)
		totalMaxPoints += maxPoints

		if questionType == string(models.QuestionTypeSBA) {
			// Parse options to find correct answer
			for key, value := range options {
				if optData, ok := value.(map[string]interface{}); ok {
					if isCorrectOption, exists := optData["is_correct"].(bool); exists && isCorrectOption {
						correctAnswerArray = append(correctAnswerArray, key)
						// Check if user's answer matches
						if len(userAnswerArray) == 1 && userAnswerArray[0] == key {
							isCorrect = true
							pointsEarned = maxPoints
							correctAnswers++
						}
					}
				}
			}
		} else if questionType == string(models.QuestionTypeTrueFalse) {
			// For TRUE_FALSE questions, extract correct answers
			correctAnswersMap = make(map[string]bool) // key -> true/false (expected answer)
			for key, value := range options {
				if optData, ok := value.(map[string]interface{}); ok {
					if isCorrectOption, exists := optData["is_correct"].(bool); exists {
						correctAnswersMap[key] = isCorrectOption
						// For backward compatibility, also populate the old array format
						if isCorrectOption {
							correctAnswerArray = append(correctAnswerArray, key)
						}
					}
				}
			}

			// Parse user's TRUE_FALSE answers in format ["a:true", "b:false", "c:true"]
			userTrueFalseAnswers := make(map[string]bool)

			for _, userAns := range userAnswerArray {
				parts := strings.Split(userAns, ":")
				if len(parts) == 2 {
					optionKey := parts[0]
					answerValue := parts[1] == "true"
					userTrueFalseAnswers[optionKey] = answerValue
				}
			}

			// Calculate score based on correct true/false answers
			if len(userTrueFalseAnswers) > 0 && len(correctAnswersMap) > 0 {
				correctMatches := 0
				totalOptions := len(correctAnswersMap)

				// Count how many of user's answers are correct
				for optionKey, expectedAnswer := range correctAnswersMap {
					if userAnswer, exists := userTrueFalseAnswers[optionKey]; exists {
						if userAnswer == expectedAnswer {
							correctMatches++
						}
					}
				}

				// Calculate score based on the proportion of correct answers
				if correctMatches == totalOptions {
					// Perfect answer: all options marked correctly
					isCorrect = true
					pointsEarned = maxPoints
					correctAnswers++
				} else if correctMatches > 0 {
					// Partial credit: (correct matches / total options) * max points
					pointsEarned = (float64(correctMatches) / float64(totalOptions)) * maxPoints
				}
			}
		}

		totalPointsEarned += pointsEarned

		// Build answer detail with proper correct answer
		var correctAnswerInterface interface{}
		if questionType == string(models.QuestionTypeSBA) {
			// For SBA, return single correct answer
			if len(correctAnswerArray) > 0 {
				correctAnswerInterface = correctAnswerArray[0]
			}
		} else {
			// For TRUE_FALSE, return map of correct true/false answers
			if questionType == string(models.QuestionTypeTrueFalse) && correctAnswersMap != nil {
				correctAnswerInterface = correctAnswersMap
			} else {
				// Legacy fallback for other question types
				correctAnswerInterface = correctAnswerArray
			}
		}

		// Create detailed answer record
		answerDetail := AnswerDetail{
			QuestionID:    questionID,
			QuestionText:  questionText,
			QuestionType:  questionType,
			UserAnswer:    userAnswerArray,
			CorrectAnswer: correctAnswerInterface,
			IsCorrect:     isCorrect,
			PointsEarned:  pointsEarned,
			MaxPoints:     int(maxPoints),
			Explanation:   explanation, // Store explanation for this question
			Options:       optionDetails,
		}

		answersData.Answers = append(answersData.Answers, answerDetail)
	}

	// Calculate score based on points earned
	score := totalPointsEarned

	// Determine if passed
	passed := score >= sessionData.Exam.PassingScore

	// Convert answers data to JSON string for storage
	answersDataJSON, err := json.Marshal(answersData)
	if err != nil {
		return dto.SubmitExamResponse{}, fmt.Errorf("failed to marshal answers data: %w", err)
	}

	// Update attempt status to completed with answers data
	err = s.examRepo.CompleteExamAttemptWithAnswers(sessionData.Attempt.ID, score, passed, string(answersDataJSON), correctAnswers)
	if err != nil {
		return dto.SubmitExamResponse{}, fmt.Errorf("failed to complete exam attempt: %w", err)
	}

	// Calculate time taken
	timeTaken := sessionData.Attempt.GetTimeSpentSeconds()

	// Create response
	response := dto.SubmitExamResponse{
		SessionID:        sessionID,
		Score:            score,
		Passed:           passed,
		TotalQuestions:   totalQuestions,
		CorrectAnswers:   correctAnswers,
		TimeTakenSeconds: timeTaken,
		SubmittedAt:      submissionTime.Format("2006-01-02T15:04:05Z"),
	}

	return response, nil
}

// GetExamResultsBySession returns raw exam attempt data by session ID for frontend processing
func (s *examService) GetExamResultsBySession(sessionID string, userID uint) (interface{}, error) {
	// Initialize logger with service context
	log := logger.WithService("ExamService").WithFields(logrus.Fields{
		"operation":  "GetExamResultsBySession",
		"session_id": sessionID,
		"user_id":    userID,
	})

	// Find attempt by session_id and user_id
	attempt, err := s.examRepo.GetAttemptBySessionAndUser(sessionID, userID)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve attempt by session and user from repository")
		return nil, fmt.Errorf("failed to get attempt by session: %w", err)
	}

	// Add attempt details to logger context
	log = log.WithFields(logrus.Fields{
		"attempt_id":     attempt.ID,
		"attempt_status": attempt.Status,
		"exam_id":        attempt.ExamID,
	})

	// Verify the attempt is in a final state (completed, auto-submitted, or abandoned)
	validStatuses := []models.AttemptStatus{
		models.AttemptStatusCompleted,
		models.AttemptStatusAutoSubmitted,
		models.AttemptStatusAbandoned,
	}

	isValidStatus := false
	for _, validStatus := range validStatuses {
		if attempt.Status == validStatus {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		log.WithFields(logrus.Fields{
			"valid_statuses": []string{
				string(models.AttemptStatusCompleted),
				string(models.AttemptStatusAutoSubmitted),
				string(models.AttemptStatusAbandoned),
			},
			"actual_status": attempt.Status,
		}).Warn("Attempt is not in a final state - cannot retrieve results")
		return nil, fmt.Errorf("exam attempt is not in a final state (completed, auto-submitted, or abandoned)")
	}

	// Parse the raw answers_data JSON to extract just the answers array
	var answersDataJson map[string]interface{}
	var answersArray interface{}

	if attempt.AnswersData != "" && attempt.AnswersData != "{}" {
		err = json.Unmarshal([]byte(attempt.AnswersData), &answersDataJson)
		if err != nil {
			log.WithError(err).WithField("answers_data_sample", attempt.AnswersData[:min(100, len(attempt.AnswersData))]).Error("Failed to parse answers data JSON")
			return nil, fmt.Errorf("failed to parse answers data: %w", err)
		}

		// Extract just the answers array from the nested JSON
		if answersDataJson != nil {
			if answers, ok := answersDataJson["answers"]; ok {
				answersArray = answers
				log.Debug("Successfully extracted answers from stored answers data")
			}
		}
	} else {
		// No answers data found - need to reconstruct from exam questions for abandoned/incomplete exams
		log.Warn("No answers data found - reconstructing question details from exam data")

		// Get the completed session data to access exam questions (final states only)
		sessionData, err := s.examRepo.GetCompletedSessionByID(sessionID)
		if err != nil {
			log.WithError(err).Error("Failed to retrieve completed session data for answer reconstruction")
			return nil, fmt.Errorf("failed to get completed session data: %w", err)
		}

		// Parse the exam questions to create answer details with correct answers
		reconstructedAnswers, err := s.reconstructAnswerDetails(sessionData.Exam.QuestionsData, sessionID)
		if err != nil {
			log.WithError(err).Error("Failed to reconstruct answer details from exam questions")
			return nil, fmt.Errorf("failed to reconstruct answer details: %w", err)
		}

		answersArray = reconstructedAnswers
		log.WithField("reconstructed_questions", len(reconstructedAnswers)).Info("Successfully reconstructed answer details from exam questions")
	}

	// Create clean response structure with authoritative data from database
	response := map[string]interface{}{
		"answers": answersArray,
		"exam_snapshot": map[string]interface{}{
			"duration_minutes":  attempt.TimeLimitSeconds / 60, // Convert seconds to minutes
			"passing_score":     attempt.PassingScore,
			"total_questions":   attempt.TotalQuestions,
			"actual_time_spent": attempt.ActualTimeSpent, // in seconds
			"score":             attempt.Score,           // final calculated score
			"correct_answers":   attempt.CorrectAnswers,  // final calculated correct answers
			"is_passed":         attempt.IsPassed,        // final pass/fail status
		},
		"submission_timestamp": attempt.CompletedAt, // authoritative completion time
	}

	return response, nil
}

// reconstructAnswerDetails creates answer details from exam questions when no answers data exists
// This is used for abandoned or incomplete exams where detailed scoring wasn't performed
func (s *examService) reconstructAnswerDetails(questionsData string, sessionID string) ([]map[string]interface{}, error) {
	// Parse questions from the exam's JSON data
	var examQuestions []map[string]interface{}
	err := json.Unmarshal([]byte(questionsData), &examQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse exam questions: %w", err)
	}

	// Get any saved answers from cache (user might have answered some questions)
	savedAnswers, err := s.examRepo.GetSessionAnswers(sessionID)
	if err != nil {
		// Continue without saved answers - user didn't answer anything
		savedAnswers = make(map[uint]string)
	}

	var reconstructedAnswers []map[string]interface{}

	// Process each question to create answer details
	for _, questionData := range examQuestions {
		// Extract question fields from map
		questionID := uint(questionData["id"].(float64))
		questionText := questionData["question_text"].(string)
		questionType := questionData["question_type"].(string)
		options := questionData["options"].(map[string]interface{})
		points := int(questionData["points"].(float64))
		explanation := ""
		if exp, exists := questionData["explanation"]; exists && exp != nil {
			explanation = exp.(string)
		}

		// Build options array for answer detail
		type OptionDetail struct {
			Key       string `json:"key"`
			Text      string `json:"text"`
			IsCorrect bool   `json:"is_correct"`
		}

		optionDetails := make([]OptionDetail, 0, len(options))
		for key, value := range options {
			if optData, ok := value.(map[string]interface{}); ok {
				text := ""
				if t, ok := optData["text"].(string); ok {
					text = t
				}
				isCorrect := false
				if c, ok := optData["is_correct"].(bool); ok {
					isCorrect = c
				}
				optionDetails = append(optionDetails, OptionDetail{
					Key:       key,
					Text:      text,
					IsCorrect: isCorrect,
				})
			}
		}

		// Get user's answer if they provided one
		userAnswer, hasAnswer := savedAnswers[questionID]
		var userAnswerArray []string
		if hasAnswer && userAnswer != "" {
			// Try to parse as JSON array first (for TRUE_FALSE multiple selections)
			var parsedArray []string
			if err := json.Unmarshal([]byte(userAnswer), &parsedArray); err == nil {
				userAnswerArray = parsedArray
			} else {
				// Single answer (for SBA questions)
				userAnswerArray = []string{userAnswer}
			}
		}

		// Find correct answers
		var correctAnswerInterface interface{}
		if questionType == string(models.QuestionTypeSBA) {
			// For SBA, find the single correct answer
			for key, value := range options {
				if optData, ok := value.(map[string]interface{}); ok {
					if isCorrectOption, exists := optData["is_correct"].(bool); exists && isCorrectOption {
						correctAnswerInterface = key
						break
					}
				}
			}
		} else if questionType == string(models.QuestionTypeTrueFalse) {
			// For TRUE_FALSE, create map of correct true/false answers
			correctAnswersMap := make(map[string]bool)
			for key, value := range options {
				if optData, ok := value.(map[string]interface{}); ok {
					if isCorrectOption, exists := optData["is_correct"].(bool); exists {
						correctAnswersMap[key] = isCorrectOption
					}
				}
			}
			correctAnswerInterface = correctAnswersMap
		}

		// Calculate if user's answer is correct (basic scoring for display)
		isCorrect := false
		pointsEarned := 0.0
		if hasAnswer {
			if questionType == string(models.QuestionTypeSBA) {
				if len(userAnswerArray) == 1 && correctAnswerInterface != nil {
					if userAnswerArray[0] == correctAnswerInterface.(string) {
						isCorrect = true
						pointsEarned = float64(points)
					}
				}
			} else if questionType == string(models.QuestionTypeTrueFalse) {
				if correctAnswersMap, ok := correctAnswerInterface.(map[string]bool); ok {
					// Parse user's TRUE_FALSE answers
					userTrueFalseAnswers := make(map[string]bool)
					for _, userAns := range userAnswerArray {
						parts := strings.Split(userAns, ":")
						if len(parts) == 2 {
							optionKey := parts[0]
							answerValue := parts[1] == "true"
							userTrueFalseAnswers[optionKey] = answerValue
						}
					}

					// Check if all answers match
					if len(userTrueFalseAnswers) == len(correctAnswersMap) {
						allCorrect := true
						for key, expected := range correctAnswersMap {
							if userAnswer, exists := userTrueFalseAnswers[key]; !exists || userAnswer != expected {
								allCorrect = false
								break
							}
						}
						if allCorrect {
							isCorrect = true
							pointsEarned = float64(points)
						}
					}
				}
			}
		}

		// Create answer detail
		answerDetail := map[string]interface{}{
			"question_id":    questionID,
			"question_text":  questionText,
			"question_type":  questionType,
			"user_answer":    userAnswerArray,
			"correct_answer": correctAnswerInterface,
			"is_correct":     isCorrect,
			"points_earned":  pointsEarned,
			"max_points":     points,
			"explanation":    explanation,
			"options":        optionDetails,
		}

		reconstructedAnswers = append(reconstructedAnswers, answerDetail)
	}

	return reconstructedAnswers, nil
}

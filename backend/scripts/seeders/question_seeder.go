package seeders

import (
	"encoding/json"
	"log"
	"github.com/Mahfuz2811/medecole/backend/internal/database"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
)

type QuestionSeeder struct{}

func NewQuestionSeeder() *QuestionSeeder {
	return &QuestionSeeder{}
}

func (s *QuestionSeeder) Seed(db *database.Database) {
	log.Println("Seeding subjects, systems, and questions...")

	subjects := s.getSampleSubjects()

	for subjectIndex, subjectData := range subjects {
		// Create subject
		subject := models.Subject{
			Name:        subjectData.Name,
			Slug:        subjectData.Slug,
			Description: &subjectData.Description,
			SortOrder:   subjectIndex + 1,
			IsActive:    true,
		}

		if err := db.DB.Create(&subject).Error; err != nil {
			log.Printf("Failed to create subject %s: %v", subject.Name, err)
			continue
		}

		log.Printf("Created subject: %s (ID: %d)", subject.Name, subject.ID)

		// Create systems for this subject
		for systemIndex, systemData := range subjectData.Systems {
			system := models.System{
				SubjectID:   subject.ID,
				Name:        systemData.Name,
				Slug:        systemData.Slug,
				Description: &systemData.Description,
				SortOrder:   systemIndex + 1,
				IsActive:    true,
			}

			if err := db.DB.Create(&system).Error; err != nil {
				log.Printf("Failed to create system %s: %v", system.Name, err)
				continue
			}

			log.Printf("  Created system: %s (ID: %d)", system.Name, system.ID)

			// Create questions for this system
			for _, questionData := range systemData.Questions {
				optionsJSON, _ := json.Marshal(questionData.Options)
				tagsJSON, _ := json.Marshal(questionData.Tags)

				question := models.Question{
					SystemID:        system.ID,
					QuestionText:    questionData.QuestionText,
					QuestionType:    questionData.QuestionType,
					DifficultyLevel: questionData.DifficultyLevel,
					Options:         string(optionsJSON),
					Explanation:     &questionData.Explanation,
					Reference:       &questionData.Reference,
					Tags:            string(tagsJSON),
					IsActive:        true,
				}

				if err := db.DB.Create(&question).Error; err != nil {
					log.Printf("Failed to create question: %v", err)
					continue
				}
			}

			log.Printf("    Created %d questions for system: %s", len(systemData.Questions), system.Name)
		}
	}
}

func (s *QuestionSeeder) getSampleSubjects() []SampleSubject {
	return []SampleSubject{
		{
			Name:        "Medicine",
			Slug:        "medicine",
			Description: "Internal Medicine and related medical specialties",
			Systems: []SampleSystem{
				{
					Name:        "Cardiovascular System",
					Slug:        "cardiovascular",
					Description: "Heart, blood vessels, and circulation",
					Questions:   s.getCardiovascularQuestions(),
				},
				{
					Name:        "Respiratory System",
					Slug:        "respiratory",
					Description: "Lungs, airways, and breathing",
					Questions:   s.getRespiratoryQuestions(),
				},
			},
		},
		{
			Name:        "Surgery",
			Slug:        "surgery",
			Description: "Surgical procedures and operative medicine",
			Systems: []SampleSystem{
				{
					Name:        "General Surgery",
					Slug:        "general-surgery",
					Description: "Common surgical procedures and techniques",
					Questions:   s.getGeneralSurgeryQuestions(),
				},
			},
		},
	}
}

// Question generators for different systems
func (s *QuestionSeeder) getCardiovascularQuestions() []SampleQuestion {
	return []SampleQuestion{
		{
			QuestionType: models.QuestionTypeSBA,
			QuestionText: "A 45-year-old man presents with chest pain, sweating, and shortness of breath. ECG shows ST-elevation in leads II, III, and aVF. What is the most likely diagnosis?",
			Options: map[string]interface{}{
				"a": map[string]interface{}{"text": "Anterior STEMI", "is_correct": false},
				"b": map[string]interface{}{"text": "Inferior STEMI", "is_correct": true},
				"c": map[string]interface{}{"text": "Posterior STEMI", "is_correct": false},
				"d": map[string]interface{}{"text": "Unstable angina", "is_correct": false},
				"e": map[string]interface{}{"text": "Pericarditis", "is_correct": false},
			},
			Explanation:     "ST-elevation in leads II, III, and aVF indicates inferior wall myocardial infarction, as these leads face the inferior wall of the heart.",
			DifficultyLevel: models.DifficultyMedium,
			Reference:       "Harrison's Principles of Internal Medicine, 21st Edition",
			Tags:            []string{"cardiology", "myocardial_infarction", "ECG", "emergency"},
		},
		{
			QuestionType: models.QuestionTypeTrueFalse,
			QuestionText: "Atrial fibrillation always requires immediate cardioversion.",
			Options: map[string]interface{}{
				"a": map[string]interface{}{"text": "Atrial fibrillation always requires immediate cardioversion", "is_correct": false},
				"b": map[string]interface{}{"text": "Rate control is preferred in stable patients", "is_correct": true},
				"c": map[string]interface{}{"text": "Cardioversion is only for unstable patients", "is_correct": true},
				"d": map[string]interface{}{"text": "Anticoagulation is always contraindicated", "is_correct": false},
				"e": map[string]interface{}{"text": "Duration of AF determines management", "is_correct": true},
			},
			Explanation:     "False. Atrial fibrillation management depends on hemodynamic stability, duration, and patient factors. Rate control and anticoagulation may be preferred over immediate cardioversion.",
			DifficultyLevel: models.DifficultyEasy,
			Reference:       "AHA/ACC/HRS Guidelines for Management of Atrial Fibrillation",
			Tags:            []string{"cardiology", "atrial_fibrillation", "treatment"},
		},
		{
			QuestionType: models.QuestionTypeSBA,
			QuestionText: "Which medication is first-line treatment for heart failure with reduced ejection fraction?",
			Options: map[string]interface{}{
				"a": map[string]interface{}{"text": "Digoxin", "is_correct": false},
				"b": map[string]interface{}{"text": "ACE inhibitor", "is_correct": true},
				"c": map[string]interface{}{"text": "Calcium channel blocker", "is_correct": false},
				"d": map[string]interface{}{"text": "Beta-blocker", "is_correct": false},
				"e": map[string]interface{}{"text": "Diuretic", "is_correct": false},
			},
			Explanation:     "ACE inhibitors are first-line therapy for HFrEF as they improve survival and reduce hospitalizations by blocking the renin-angiotensin system.",
			DifficultyLevel: models.DifficultyMedium,
			Reference:       "AHA/ACC Heart Failure Guidelines",
			Tags:            []string{"cardiology", "heart_failure", "pharmacology"},
		},
	}
}

func (s *QuestionSeeder) getRespiratoryQuestions() []SampleQuestion {
	return []SampleQuestion{
		{
			QuestionType: models.QuestionTypeSBA,
			QuestionText: "A 28-year-old smoker presents with productive cough, fever, and consolidation on chest X-ray. What is the most likely organism?",
			Options: map[string]interface{}{
				"a": map[string]interface{}{"text": "Streptococcus pneumoniae", "is_correct": true},
				"b": map[string]interface{}{"text": "Haemophilus influenzae", "is_correct": false},
				"c": map[string]interface{}{"text": "Mycoplasma pneumoniae", "is_correct": false},
				"d": map[string]interface{}{"text": "Legionella pneumophila", "is_correct": false},
				"e": map[string]interface{}{"text": "Staphylococcus aureus", "is_correct": false},
			},
			Explanation:     "Streptococcus pneumoniae is the most common cause of community-acquired pneumonia, especially in smokers.",
			DifficultyLevel: models.DifficultyMedium,
			Reference:       "IDSA Guidelines for Community-Acquired Pneumonia",
			Tags:            []string{"respiratory", "pneumonia", "infectious_disease"},
		},
		{
			QuestionType: models.QuestionTypeTrueFalse,
			QuestionText: "COPD is characterized by reversible airflow obstruction.",
			Options: map[string]interface{}{
				"a": map[string]interface{}{"text": "COPD is characterized by reversible airflow obstruction", "is_correct": false},
				"b": map[string]interface{}{"text": "Airflow obstruction is permanent in COPD", "is_correct": true},
				"c": map[string]interface{}{"text": "Bronchodilators completely reverse obstruction", "is_correct": false},
				"d": map[string]interface{}{"text": "COPD has some reversible components", "is_correct": true},
				"e": map[string]interface{}{"text": "Asthma has reversible airflow obstruction", "is_correct": true},
			},
			Explanation:     "False. COPD is characterized by irreversible or poorly reversible airflow obstruction, unlike asthma which has reversible obstruction.",
			DifficultyLevel: models.DifficultyEasy,
			Reference:       "GOLD Guidelines for COPD",
			Tags:            []string{"respiratory", "COPD", "pathophysiology"},
		},
	}
}

func (s *QuestionSeeder) getGeneralSurgeryQuestions() []SampleQuestion {
	return []SampleQuestion{
		{
			QuestionType: models.QuestionTypeSBA,
			QuestionText: "A 25-year-old man presents with sudden onset right lower quadrant pain, nausea, and fever. What is the most likely diagnosis?",
			Options: map[string]interface{}{
				"a": map[string]interface{}{"text": "Appendicitis", "is_correct": true},
				"b": map[string]interface{}{"text": "Cholecystitis", "is_correct": false},
				"c": map[string]interface{}{"text": "Diverticulitis", "is_correct": false},
				"d": map[string]interface{}{"text": "Kidney stone", "is_correct": false},
				"e": map[string]interface{}{"text": "Hernia", "is_correct": false},
			},
			Explanation:     "Right lower quadrant pain with nausea and fever in a young adult is classic for acute appendicitis.",
			DifficultyLevel: models.DifficultyEasy,
			Reference:       "Schwartz's Principles of Surgery",
			Tags:            []string{"surgery", "appendicitis", "acute_abdomen"},
		},
	}
}

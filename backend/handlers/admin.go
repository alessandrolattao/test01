package handlers

import (
	"net/http"
	"os"
	"sort"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB        *gorm.DB
	JWTSecret string
	AIService *services.AIService
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := h.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(admin.ID, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AdminHandler) ListCandidates(c *gin.Context) {
	var candidates []models.Candidate
	if err := h.DB.Order("total_score DESC").Find(&candidates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch candidates"})
		return
	}
	c.JSON(http.StatusOK, candidates)
}

func (h *AdminHandler) GetCandidate(c *gin.Context) {
	id := c.Param("id")

	var candidate models.Candidate
	err := h.DB.
		Preload("CandidateAnswers").
		Preload("CandidateAnswers.Question").
		Preload("CandidateAnswers.Question.Answers", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&candidate, "id = ?", id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	// Sort by question order
	sort.Slice(candidate.CandidateAnswers, func(i, j int) bool {
		return candidate.CandidateAnswers[i].Question.SortOrder < candidate.CandidateAnswers[j].Question.SortOrder
	})

	// Build answers in the format frontend expects
	type answerInfo struct {
		ID    string `json:"id"`
		Text  string `json:"text"`
		Score int    `json:"score"`
	}
	type questionAnswer struct {
		QuestionID       string       `json:"question_id"`
		QuestionText     string       `json:"question_text"`
		SelectedAnswerID string       `json:"selected_answer_id"`
		Score            int          `json:"score"`
		AllAnswers       []answerInfo `json:"all_answers"`
	}

	var answers []questionAnswer
	for _, ca := range candidate.CandidateAnswers {
		qa := questionAnswer{
			QuestionID:       ca.QuestionID,
			QuestionText:     ca.Question.Text,
			SelectedAnswerID: ca.AnswerID,
			Score:            ca.Score,
		}
		for _, a := range ca.Question.Answers {
			qa.AllAnswers = append(qa.AllAnswers, answerInfo{
				ID:    a.ID,
				Text:  a.Text,
				Score: a.Score,
			})
		}
		answers = append(answers, qa)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               candidate.ID,
		"first_name":       candidate.FirstName,
		"last_name":        candidate.LastName,
		"email":            candidate.Email,
		"questionnaire_id": candidate.QuestionnaireID,
		"total_score":      candidate.TotalScore,
		"audio_path":       candidate.AudioPath,
		"completed":        candidate.Completed,
		"created_at":       candidate.CreatedAt,
		"answers":          answers,
		"transcript":       candidate.Transcript,
		"ai_analysis":      candidate.AIAnalysis,
		"ai_score":         candidate.AIScore,
		"analysis_status":  candidate.AnalysisStatus,
	})
}

func (h *AdminHandler) GetCandidateAudio(c *gin.Context) {
	id := c.Param("id")

	var candidate models.Candidate
	if err := h.DB.First(&candidate, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	if candidate.AudioPath == nil || *candidate.AudioPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no audio recording found"})
		return
	}

	data, err := os.ReadFile(*candidate.AudioPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read audio file"})
		return
	}
	c.Data(http.StatusOK, "audio/webm", data)
}

func (h *AdminHandler) ListQuestionnaires(c *gin.Context) {
	var questionnaires []models.Questionnaire
	if err := h.DB.Preload("Questions").Order("version DESC").Find(&questionnaires).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch questionnaires"})
		return
	}

	type questionnaireItem struct {
		ID             string `json:"id"`
		Version        int    `json:"version"`
		IsActive       bool   `json:"is_active"`
		CreatedAt      string `json:"created_at"`
		QuestionsCount int    `json:"questions_count"`
	}

	var result []questionnaireItem
	for _, q := range questionnaires {
		result = append(result, questionnaireItem{
			ID:             q.ID,
			Version:        q.Version,
			IsActive:       q.IsActive,
			CreatedAt:      q.CreatedAt.Format("2006-01-02T15:04:05Z"),
			QuestionsCount: len(q.Questions),
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminHandler) GetQuestionnaire(c *gin.Context) {
	id := c.Param("id")

	var questionnaire models.Questionnaire
	err := h.DB.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&questionnaire, "id = ?", id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "questionnaire not found"})
		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

type createQuestionnaireRequest struct {
	Questions []createQuestionRequest `json:"questions" binding:"required,dive"`
}

type createQuestionRequest struct {
	Text      string                `json:"text" binding:"required"`
	SortOrder int                   `json:"sort_order"`
	Answers   []createAnswerRequest `json:"answers" binding:"required,dive"`
}

type createAnswerRequest struct {
	Text      string `json:"text" binding:"required"`
	Score     int    `json:"score"`
	SortOrder int    `json:"sort_order"`
}

func (h *AdminHandler) CreateQuestionnaire(c *gin.Context) {
	var req createQuestionnaireRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var questionnaire models.Questionnaire

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Deactivate all existing questionnaires
		if err := tx.Model(&models.Questionnaire{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
			return err
		}

		// Create new questionnaire
		questionnaire = models.Questionnaire{IsActive: true}
		if err := tx.Create(&questionnaire).Error; err != nil {
			return err
		}

		// Create questions and answers
		for _, qReq := range req.Questions {
			question := models.Question{
				QuestionnaireID: questionnaire.ID,
				Text:            qReq.Text,
				SortOrder:       qReq.SortOrder,
			}
			if err := tx.Create(&question).Error; err != nil {
				return err
			}

			for _, aReq := range qReq.Answers {
				answer := models.Answer{
					QuestionID: question.ID,
					Text:       aReq.Text,
					Score:      aReq.Score,
					SortOrder:  aReq.SortOrder,
				}
				if err := tx.Create(&answer).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create questionnaire"})
		return
	}

	// Reload with associations
	h.DB.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&questionnaire, "id = ?", questionnaire.ID)

	c.JSON(http.StatusCreated, questionnaire)
}

func (h *AdminHandler) ReanalyzeCandidate(c *gin.Context) {
	id := c.Param("id")

	if h.AIService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI service not configured"})
		return
	}

	if err := h.AIService.ReanalyzeCandidate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reanalysis started"})
}

package handlers

import (
	"net/http"

	"recruitment-platform/middleware"
	"recruitment-platform/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB        *gorm.DB
	JWTSecret string
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
		Preload("CandidateAnswers.Question").
		Preload("CandidateAnswers.Answer").
		First(&candidate, "id = ?", id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	c.JSON(http.StatusOK, candidate)
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

	c.File(*candidate.AudioPath)
}

func (h *AdminHandler) ListQuestionnaires(c *gin.Context) {
	var questionnaires []models.Questionnaire
	if err := h.DB.Order("version DESC").Find(&questionnaires).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch questionnaires"})
		return
	}
	c.JSON(http.StatusOK, questionnaires)
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

package handlers

import (
	"net/http"

	"recruitment-platform/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionnaireHandler struct {
	DB *gorm.DB
}

// publicQuestion is the public view of a question (answers without scores)
type publicQuestion struct {
	ID              string                `json:"id"`
	QuestionnaireID string                `json:"questionnaire_id"`
	Text            string                `json:"text"`
	SortOrder       int                   `json:"sort_order"`
	Answers         []models.PublicAnswer `json:"answers"`
}

// publicQuestionnaire is the public view (no scores exposed)
type publicQuestionnaire struct {
	ID        string           `json:"id"`
	Version   int              `json:"version"`
	Questions []publicQuestion `json:"questions"`
}

// GetActive returns the active questionnaire with questions and answers, but without scores
func (h *QuestionnaireHandler) GetActive(c *gin.Context) {
	var questionnaire models.Questionnaire
	err := h.DB.
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Questions.Answers", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("is_active = ?", true).
		First(&questionnaire).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active questionnaire found"})
		return
	}

	// Strip scores from answers
	result := publicQuestionnaire{
		ID:      questionnaire.ID,
		Version: questionnaire.Version,
	}
	for _, q := range questionnaire.Questions {
		pq := publicQuestion{
			ID:              q.ID,
			QuestionnaireID: q.QuestionnaireID,
			Text:            q.Text,
			SortOrder:       q.SortOrder,
		}
		for _, a := range q.Answers {
			pq.Answers = append(pq.Answers, models.PublicAnswer{
				ID:         a.ID,
				QuestionID: a.QuestionID,
				Text:       a.Text,
				SortOrder:  a.SortOrder,
			})
		}
		result.Questions = append(result.Questions, pq)
	}

	c.JSON(http.StatusOK, result)
}

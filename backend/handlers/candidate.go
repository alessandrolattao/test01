package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"recruitment-platform/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CandidateHandler struct {
	DB *gorm.DB
}

type registerRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

func (h *CandidateHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidate := models.Candidate{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	if err := h.DB.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	c.JSON(http.StatusCreated, candidate)
}

type submitAnswersRequest struct {
	Answers []answerEntry `json:"answers" binding:"required,dive"`
}

type answerEntry struct {
	QuestionID string `json:"question_id" binding:"required"`
	AnswerID   string `json:"answer_id" binding:"required"`
}

func (h *CandidateHandler) SubmitAnswers(c *gin.Context) {
	id := c.Param("id")

	var req submitAnswersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var candidate models.Candidate
	if err := h.DB.First(&candidate, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	// Get the active questionnaire
	var questionnaire models.Questionnaire
	if err := h.DB.Where("is_active = ?", true).First(&questionnaire).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active questionnaire found"})
		return
	}

	totalScore := 0
	var candidateAnswers []models.CandidateAnswer

	// Look up each answer's score server-side
	for _, entry := range req.Answers {
		var answer models.Answer
		if err := h.DB.First(&answer, "id = ? AND question_id = ?", entry.AnswerID, entry.QuestionID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid answer %s for question %s", entry.AnswerID, entry.QuestionID)})
			return
		}

		totalScore += answer.Score
		candidateAnswers = append(candidateAnswers, models.CandidateAnswer{
			CandidateID: id,
			QuestionID:  entry.QuestionID,
			AnswerID:    entry.AnswerID,
			Score:       answer.Score,
		})
	}

	// Save everything in a transaction
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		for i := range candidateAnswers {
			if err := tx.Create(&candidateAnswers[i]).Error; err != nil {
				return err
			}
		}

		return tx.Model(&candidate).Updates(map[string]interface{}{
			"questionnaire_id": questionnaire.ID,
			"total_score":      totalScore,
		}).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save answers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_score": totalScore,
		"message":     "answers submitted successfully",
	})
}

func (h *CandidateHandler) UploadAudio(c *gin.Context) {
	id := c.Param("id")

	var candidate models.Candidate
	if err := h.DB.First(&candidate, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "candidate not found"})
		return
	}

	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio file is required"})
		return
	}

	// Ensure upload directory exists
	uploadDir := "uploads/audio"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	filename := id + ".webm"
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save audio file"})
		return
	}

	// Update candidate record
	h.DB.Model(&candidate).Updates(map[string]interface{}{
		"audio_path": savePath,
		"completed":  true,
	})

	c.JSON(http.StatusOK, gin.H{"message": "audio uploaded successfully"})
}

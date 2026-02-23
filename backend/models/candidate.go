package models

import "time"

type Candidate struct {
	ID               string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName        string            `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName         string            `gorm:"type:varchar(100);not null" json:"last_name"`
	Email            string            `gorm:"uniqueIndex;not null" json:"email"`
	QuestionnaireID  *string           `gorm:"type:uuid" json:"questionnaire_id,omitempty"`
	TotalScore       int               `gorm:"default:0" json:"total_score"`
	AudioPath        *string           `gorm:"type:varchar(500)" json:"audio_path,omitempty"`
	Completed        bool              `gorm:"default:false" json:"completed"`
	CreatedAt        time.Time         `json:"created_at"`
	CandidateAnswers []CandidateAnswer `gorm:"foreignKey:CandidateID" json:"candidate_answers,omitempty"`
}

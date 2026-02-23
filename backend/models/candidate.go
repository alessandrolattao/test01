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
	Transcript       *string           `gorm:"type:text" json:"transcript"`
	AIAnalysis       *string           `gorm:"type:text;column:ai_analysis" json:"ai_analysis"`
	AIScore          *int              `gorm:"column:ai_score" json:"ai_score"`
	AnalysisStatus   string            `gorm:"type:varchar(20);default:pending" json:"analysis_status"`
	CreatedAt        time.Time         `json:"created_at"`
	CandidateAnswers []CandidateAnswer `gorm:"foreignKey:CandidateID" json:"candidate_answers,omitempty"`
}

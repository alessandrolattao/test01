package models

import "time"

type Questionnaire struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Version   int        `json:"version"`
	IsActive  bool       `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	Questions []Question `gorm:"foreignKey:QuestionnaireID" json:"questions,omitempty"`
}

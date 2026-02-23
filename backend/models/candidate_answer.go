package models

type CandidateAnswer struct {
	ID          string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CandidateID string   `gorm:"type:uuid;not null" json:"candidate_id"`
	QuestionID  string   `gorm:"type:uuid;not null" json:"question_id"`
	AnswerID    string   `gorm:"type:uuid;not null" json:"answer_id"`
	Score       int      `gorm:"not null" json:"score"`
	Question    Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	Answer      Answer   `gorm:"foreignKey:AnswerID" json:"answer,omitempty"`
}

package models

type Answer struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	QuestionID string `gorm:"type:uuid;not null" json:"question_id"`
	Text       string `gorm:"type:text;not null" json:"text"`
	Score      int    `gorm:"not null;default:0" json:"score"`
	SortOrder  int    `gorm:"not null" json:"sort_order"`
}

// PublicAnswer hides the score field from candidates
type PublicAnswer struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	Text       string `json:"text"`
	SortOrder  int    `json:"sort_order"`
}

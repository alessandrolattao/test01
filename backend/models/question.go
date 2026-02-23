package models

type Question struct {
	ID               string   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	QuestionnaireID  string   `gorm:"type:uuid;not null" json:"questionnaire_id"`
	Text             string   `gorm:"type:text;not null" json:"text"`
	SortOrder        int      `gorm:"not null" json:"sort_order"`
	Answers          []Answer `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}

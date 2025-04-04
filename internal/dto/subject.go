package dto

//type Subject struct {
//	ID             int    `json:"id"`
//	Comment        string `json:"comment,omitempty"`
//	DateCreate     string `json:"date_create"`
//	Description    string `json:"description,omitempty"`
//	LastTimeUpdate string `json:"last_time_update"`
//	Name           string `json:"name"`
//	Type           string `json:"type,omitempty"`
//	ParentID       int    `json:"parent_id"`
//	QuestionCount  int    `json:"question_count"`
//}

type Subject struct {
	ID                    int64  `json:"id,string"`
	Name                  string `json:"name"`
	Description           string `json:"description,omitempty"`
	CreatorUserId         int64  `json:"creator_user_id,string"`
	Active                bool   `json:"active"`
	ParentId              int64  `json:"parent_id,string"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	QuestionCount         int    `json:"question_count"`
	ApprovedQuestionCount int    `json:"approved_question_count"`
}

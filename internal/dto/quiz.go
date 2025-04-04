package dto

type InputQuiz struct {
	Name          string     `json:"name,omitempty"`
	Description   string     `json:"description,omitempty"`
	CreatorUserID int64      `json:"creator_user_id,string,omitempty"`
	QuestionIDs   Int64Array `json:"question_ids,omitempty"`
}

type Quiz struct {
	ID          int64      `json:"id,string,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Creator     User       `json:"creator,omitempty"`
	QuestionIDs Int64Array `json:"question_ids,omitempty"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

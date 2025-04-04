package dto

// input types

type InputQuestion struct {
	Text            string                 `json:"text,omitempty"`
	Code            string                 `json:"code,omitempty"`
	Variants        map[string]interface{} `json:"variants"`
	Answer          map[string]interface{} `json:"answer"`
	TypeID          int64                  `json:"type_id,string,omitempty"`
	StatusName      QuestionStatusName     `json:"status_name,omitempty"`
	SubjectID       int64                  `json:"subject_id,string,omitempty"`
	CreatorUserID   int64                  `json:"creator_user_id,string,omitempty"`
	ModeratorUserID int64                  `json:"moderator_user_id,string,omitempty"`
	ModeratedAt     string                 `json:"moderated_at,omitempty"`
}

// internal/output types //TODO

type Question struct {
	ID          int64                  `json:"id,string,omitempty"`
	Text        string                 `json:"text,omitempty"`
	Code        string                 `json:"code,omitempty"`
	Variants    map[string]interface{} `json:"variants"`
	Answer      map[string]interface{} `json:"answer"`
	Type        QuestionType           `json:"type,omitempty"`
	Status      QuestionStatusName     `json:"status,omitempty"`
	SubjectID   int64                  `json:"subject_id,string,omitempty"`
	SubjectName string                 `json:"subject_name,omitempty"`
	Creator     User                   `json:"creator,omitempty"`
	Moderator   User                   `json:"moderator,omitempty"`
	ModeratedAt string                 `json:"moderated_at,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
}

type QuestionTypeName string

const (
	QuestionTypeNameTest       = "Тест"
	QuestionTypeNameComparison = "Сопоставление"
	QuestionTypeNameText       = "Текст"
)

type QuestionType struct {
	ID               int    `json:"id,string"`
	QuestionTypeName string `json:"name"`
}

type QuestionStatus struct {
	ID   int    `json:"id,string"`
	Name string `json:"name"`
}

type QuestionStatusName string

const (
	QuestionStatusNameCreated  = "Создан"
	QuestionStatusNameApproved = "Одобрен"
	QuestionStatusNameDeclined = "Отклонен"
)

//type QuestionAnswer struct {
//	QuestionID int         `json:"question_id"`
//	Answer     interface{} `json:"answer"`
//}
//
//type QuestionChoice struct {
//	QuestionID   int         `json:"question_id"`
//	QuestionText string      `json:"question_text"`
//	QuestionCode string      `json:"question_code"`
//	Variant      interface{} `json:"variant"`
//}

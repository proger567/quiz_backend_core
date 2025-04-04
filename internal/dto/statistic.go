package dto

type Statistic struct {
	TestsCount                        int                    `json:"tests_count"`
	SubjectCount                      int                    `json:"subject_count"`
	QuestionCount                     int                    `json:"question_count"`
	QuestionToModerateCount           int                    `json:"question_to_moderate_count"`
	TestCountCreatedByCurrentUser     int                    `json:"test_count_created_by_current_user"`
	QuestionCountCreatedByCurrentUser int                    `json:"question_count_created_by_current_user"`
	QuestionCountsBySubject           []SubjectStatisticItem `json:"question_counts_by_subject,omitempty"`
}

type SubjectStatisticItem struct {
	SubjectId     int64  `json:"subject_id,string,omitempty"`
	SubjectName   string `json:"subject_name,omitempty"`
	QuestionCount int    `json:"question_count,omitempty"`
}

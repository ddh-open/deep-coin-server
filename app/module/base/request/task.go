package request

// TaskReceiver 任务
type TaskReceiver struct {
	TaskID     int    `json:"taskID"`
	UserID     string `json:"userID,omitempty"`
	UserName   string `json:"username,omitempty"`
	Pass       string `json:"pass,omitempty"`
	Company    string `json:"company,omitempty"`
	ProcInstID int    `json:"procInstID,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Candidate  string `json:"candidate,omitempty"`
}

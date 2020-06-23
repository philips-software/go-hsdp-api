package iron

import "time"

type TasksServices struct {
	client    *Client
	projectID string
}

type Task struct {
	ID            string     `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	ProjectID     string     `json:"project_id"`
	CodeID        string     `json:"code_id"`
	CodeHistoryID string     `json:"code_history_id"`
	Status        string     `json:"status"`
	Msg           string     `json:"msg"`
	CodeName      string     `json:"code_name"`
	CodeRev       string     `json:"code_rev"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	Timeout       int        `json:"timeout"`
	Payload       string     `json:"payload"`
	ScheduleID    string     `json:"schedule_id"`
	MessageID     string     `json:"message_id"`
	Cluster       string     `json:"cluster"`
	Duration      int        `json:"duration,omitempty"`
	LogSize       int        `json:"log_size,omitempty"`
}

func (t *TasksServices) GetTasks() (*[]Task, *Response, error) {
	req, err := t.client.NewRequest("GET", "projects/"+t.projectID+"/tasks", nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var tasks struct {
		Tasks []Task `json:"tasks"`
	}
	resp, err := t.client.Do(req, &tasks)
	return &tasks.Tasks, resp, err
}

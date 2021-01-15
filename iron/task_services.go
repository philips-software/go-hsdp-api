package iron

import (
	"time"
)

type TasksServices struct {
	client    *Client
	projectID string
}

type Task struct {
	ID            string     `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	ProjectID     string     `json:"project_id,omitempty"`
	CodeID        string     `json:"code_id,omitempty"`
	CodeHistoryID string     `json:"code_history_id,omitempty"`
	Status        string     `json:"status,omitempty"`
	Msg           string     `json:"msg,omitempty"`
	CodeName      string     `json:"code_name,omitempty"`
	CodeRev       string     `json:"code_rev,omitempty"`
	StartTime     *time.Time `json:"start_time,omitempty"`
	EndTime       *time.Time `json:"end_time,omitempty"`
	Timeout       int        `json:"timeout,omitempty"`
	Payload       string     `json:"payload,omitempty"`
	ScheduleID    string     `json:"schedule_id,omitempty"`
	MessageID     string     `json:"message_id,omitempty"`
	Cluster       string     `json:"cluster,omitempty"`
	Duration      int        `json:"duration,omitempty"`
	LogSize       int        `json:"log_size,omitempty"`
}

// GetTasks gets the tasks of the project
func (t *TasksServices) GetTasks() (*[]Task, *Response, error) {
	req, err := t.client.newRequest(
		"GET",
		t.client.Path("projects", t.projectID, "tasks"),
		nil,
		nil)
	if err != nil {
		return nil, nil, err
	}
	var tasks struct {
		Tasks []Task `json:"tasks"`
	}
	resp, err := t.client.Do(req, &tasks)
	return &tasks.Tasks, resp, err
}

// GetTask gets info on a single task
func (t *TasksServices) GetTask(taskID string) (*Task, *Response, error) {
	req, err := t.client.newRequest(
		"GET",
		t.client.Path("projects", t.projectID, "tasks", taskID),
		nil,
		nil)
	if err != nil {
		return nil, nil, err
	}
	var task Task
	resp, err := t.client.Do(req, &task)
	return &task, resp, err
}

// QueueTask queues a single task for execution
func (t *TasksServices) QueueTask(task Task) (*Task, *Response, error) {
	taskList := []Task{task}
	tasks, resp, err := t.QueueTasks(taskList)
	if err != nil {
		return nil, resp, err
	}
	if len(*tasks) == 0 {
		return nil, resp, err
	}
	return &(*tasks)[0], resp, err
}

// QueueTasks queues one or more tasks for execution
func (t *TasksServices) QueueTasks(tasks []Task) (*[]Task, *Response, error) {
	var queueRequest struct {
		Tasks []Task `json:"tasks"`
	}
	queueRequest.Tasks = tasks

	req, err := t.client.newRequest(
		"POST",
		t.client.Path("projects", t.projectID, "tasks"),
		&queueRequest,
		nil)
	if err != nil {
		return nil, nil, err
	}
	var queueResponse struct {
		Tasks []Task `json:"tasks"`
	}
	resp, err := t.client.Do(req, &queueResponse)
	return &queueResponse.Tasks, resp, err
}

// CancelTask cancels the given task
func (t *TasksServices) CancelTask(taskID string) (bool, *Response, error) {
	req, err := t.client.newRequest(
		"POST",
		t.client.Path("projects", t.projectID, "tasks", taskID, "cancel"),
		nil,
		nil)
	if err != nil {
		return false, nil, err
	}
	var cancelResponse struct {
		Message string `json:"msg"`
	}
	resp, err := t.client.Do(req, &cancelResponse)
	if cancelResponse.Message != "Cancelled" {
		return false, resp, err
	}
	return true, resp, nil
}

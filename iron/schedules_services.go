package iron

import "time"

type SchedulesServices struct {
	client    *Client
	projectID string
}

type Schedule struct {
	ID          string     `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	ProjectID   string     `json:"project_id,omitempty"`
	Status      string     `json:"status,omitempty"`
	CodeName    string     `json:"code_name"`
	StartAt     *time.Time `json:"start_at,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	NextStart   *time.Time `json:"next_start,omitempty"`
	LastRunTime *time.Time `json:"last_run_time,omitempty"`
	Timeout     int        `json:"timeout"`
	RunTimes    int        `json:"run_times,omitempty"`
	RunEvery    int        `json:"run_every,omitempty"`
	Cluster     string     `json:"cluster,omitempty"`
	Payload     string     `json:"payload,omitempty"`
}

type pageOptions struct {
	Page    *int `url:"page,omitempty"`
	PerPage *int `url:"per_page,omitempty"`
}

// CreateSchedules creates one or more schedules
func (s *SchedulesServices) CreateSchedules(schedules []Schedule) (*[]Schedule, *Response, error) {
	var createSchedules struct {
		Schedules []Schedule `json:"schedules"`
	}
	createSchedules.Schedules = schedules

	path := s.client.Path("projects", s.projectID, "schedules")
	req, err := s.client.newRequest(
		"POST",
		path,
		&createSchedules,
		nil)
	if err != nil {
		return nil, nil, err
	}
	var schedulesResponse struct {
		Schedules []Schedule `json:"schedules"`
	}
	resp, err := s.client.do(req, &schedulesResponse)
	return &schedulesResponse.Schedules, resp, err
}

// CreateSchedule creates a schedule
func (s *SchedulesServices) CreateSchedule(schedule Schedule) (*Schedule, *Response, error) {
	schedules, resp, err := s.CreateSchedules([]Schedule{schedule})
	if err != nil {
		return nil, resp, err
	}
	if len(*schedules) == 0 {
		return nil, resp, err
	}
	return &(*schedules)[0], resp, err
}

// GetSchedules gets the schedules of the project
func (s *SchedulesServices) GetSchedules() (*[]Schedule, *Response, error) {
	var schedules struct {
		Schedules []Schedule `json:"schedules"`
	}
	page := 0
	perPage := 100

	path := s.client.Path("projects", s.projectID, "schedules")
	req, err := s.client.newRequest(
		"GET",
		path,
		pageOptions{
			PerPage: &perPage,
			Page:    &page,
		},
		nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.do(req, &schedules)
	return &schedules.Schedules, resp, err
}

// GetSchedulesWithCode gets schedules which use code
func (s *SchedulesServices) GetSchedulesWithCode(codeName string) (*[]Schedule, *Response, error) {
	schedules, resp, err := s.GetSchedules()
	if err != nil {
		return nil, resp, err
	}
	filtered := make([]Schedule, 0)
	for _, schedule := range *schedules {
		if schedule.CodeName == codeName {
			filtered = append(filtered, schedule)
		}
	}
	return &filtered, resp, err
}

// GetSchedule gets info on a schedule
func (s *SchedulesServices) GetSchedule(scheduleID string) (*Schedule, *Response, error) {
	path := s.client.Path("projects", s.projectID, "schedules", scheduleID)

	page := 0
	perPage := 100
	req, err := s.client.newRequest(
		"GET",
		path,
		pageOptions{
			PerPage: &perPage,
			Page:    &page,
		},
		nil)
	if err != nil {
		return nil, nil, err
	}
	var schedule Schedule
	resp, err := s.client.do(req, &schedule)
	return &schedule, resp, err
}

// CancelSchedule cancels a schedule
func (s *SchedulesServices) CancelSchedule(scheduleID string) (bool, *Response, error) {
	path := s.client.Path("projects", s.projectID, "schedules", scheduleID, "cancel")
	req, err := s.client.newRequest(
		"POST",
		path,
		nil,
		nil)
	if err != nil {
		return false, nil, err
	}
	var cancelResponse struct {
		Message string `json:"msg"`
	}
	resp, err := s.client.do(req, &cancelResponse)
	if cancelResponse.Message != "Cancelled" {
		return false, resp, err
	}
	return true, resp, nil
}

package digdag

import (
	"errors"
	"strconv"
	"net/http"
	"net/url"
	"github.com/hashicorp/errwrap"
	uuid "github.com/satori/go.uuid"
)

type attempts struct {
	Attempts []Attempt `json:"attempts"`
}

// Attempt is the struct for digdag attempt
type Attempt struct {
	ID      string `json:"id"`
	Index   int    `json:"index"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Workflow struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"workflow"`
	SessionID        string      `json:"sessionId"`
	SessionUUID      string      `json:"sessionUuid"`
	SessionTime      string      `json:"sessionTime"`
	RetryAttemptName interface{} `json:"retryAttemptName"`
	Done             bool        `json:"done"`
	Success          bool        `json:"success"`
	CancelRequested  bool        `json:"cancelRequested"`
	Params           interface{} `json:"params"`
	CreatedAt        string      `json:"createdAt"`
	FinishedAt       string      `json:"finishedAt"`
}

// CreateAttempt is struct for create a new attempt
type CreateAttempt struct {
	Attempt
	WorkflowID       string                 `json:"workflowId"`
	SessionTime      string                 `json:"sessionTime"`
	RetryAttemptName string                 `json:"retryAttemptName,omitempty"`
	Params           map[string]interface{} `json:"params"` // TODO: set the optional params.
}

type tasks struct {
	Tasks []Task `json:"tasks"`
}

// Task is struct for attempts task result
type Task struct {
	ID           string        `json:"id"`
	FullName     string        `json:"fullName"`
	ParentID     interface{}   `json:"parentId"`
	Config       interface{}   `json:"config"`
	Upstreams    []interface{} `json:"upstreams"`
	State        string        `json:"state"`
	ExportParams interface{}   `json:"exportParams"`
	StoreParams  interface{}   `json:"storeParams"`
	StateParams  interface{}   `json:"stateParams"`
	UpdatedAt    string        `json:"updatedAt"`
	RetryAt      interface{}   `json:"retryAt"`
	StartedAt    interface{}   `json:"startedAt"`
	IsGroup      bool          `json:"isGroup"`
}

// NewCreateAttempt to create a new CreateAttempt struct
func NewCreateAttempt(workflowID, sessionTime, retryAttemptName string) *CreateAttempt {
	ca := new(CreateAttempt)
	ca.WorkflowID = workflowID
	ca.SessionTime = sessionTime
	ca.RetryAttemptName = retryAttemptName
	// TODO: set the optional params.
	ca.Params = map[string]interface{}{}

	return ca
}

// GetAttempts get attempts response
func (c *Client) GetAttempts(includeRetried bool) ([]Attempt, error) {
	spath := "/api/attempts"

	params := url.Values{}
	params.Set("project", c.ProjectName)
	params.Set("workflow", c.WorkflowName)
	params.Set("include_retried", strconv.FormatBool(includeRetried))

	var attempts *attempts
	err := c.doReq(http.MethodGet, spath, params, &attempts)
	if err != nil {
		return nil, err
	}

	// If any attempts not found
	if len(attempts.Attempts) == 0 {
		err := errors.New("attempts does not exist at `" + c.WorkflowName + "` workflow")
		return nil, err
	}

	return attempts.Attempts, err
}

// GetAttemptIDs to get attemptID from sessionTime
func (c *Client) GetAttemptIDs() (attemptIDs []string, err error) {
	attempts, err := c.GetAttempts(true)
	if err != nil {
		return nil, err
	}

	for k := range attempts {
		sessionTime := attempts[k].SessionTime

		if sessionTime == c.SessionTime {
			attemptIDs = append(attemptIDs, attempts[k].ID)
		}
	}

	// If any sessionTime not found
	return attemptIDs, err
}

// GetTasks to get tasks list
func (c *Client) GetTasks(attemptID string) ([]Task, error) {
	spath := "/api/attempts/" + attemptID + "/tasks"

	var tasks *tasks
	err := c.doReq(http.MethodGet, spath, nil, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks.Tasks, err
}

// GetTaskResult to get task result
func (c *Client) GetTaskResult(attemptIDs []string, taskName string) (*Task, error) {

	for _, attemptID := range attemptIDs {
		tasks, err := c.GetTasks(attemptID)

		for k := range tasks {
			if tasks[k].FullName == taskName {
				state := tasks[k].State
				if state == "success" {
					return &tasks[k], nil
				}

				err = errors.New("task `" + taskName + "` state is " + state)
				return nil, err
			}
		}
	}

	err := errors.New("task `" + taskName + "` result not found")
	return nil, err
}

// CreateNewAttempt to create a new attempt
func (c *Client) CreateNewAttempt(workflowID, date string, retry bool) (attempt *CreateAttempt, done bool, err error) {
	spath := "/api/attempts"

	ca := NewCreateAttempt(workflowID, c.SessionTime, "")

	// Retry workflow
	if retry == true {
		// TODO: add retry attempt name (optional)
		generatedUUID := uuid.NewV4()
		textID, err := generatedUUID.MarshalText()
		if err != nil {
			return nil, done, err
		}
		ca.RetryAttemptName = string(textID)
	}

	// Create new attempt
	err = c.doReq(http.MethodPut, spath, nil, &ca)
	if err != nil {
		// if already session exist
		if errwrap.Contains(err, "409 Conflict") {
			done := true
			return nil, done, err
		}
		return nil, done, err
	}

	return ca, done, err
}

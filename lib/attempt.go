package digdag

import (
	"errors"

	"strings"

	"net/http"
	"net/url"

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

// PutAttempt is struct for create a new attemp
type PutAttempt struct {
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

// NewPutAttempt to create a new PutAttempt struct
func NewPutAttempt(workflowID, sessionTime, retryAttemptName string) *PutAttempt {
	pa := new(PutAttempt)
	pa.WorkflowID = workflowID
	pa.SessionTime = sessionTime
	pa.RetryAttemptName = retryAttemptName
	// TODO: set the optional params.
	pa.Params = map[string]interface{}{}

	return pa
}

// GetAttempts get atemmpts reponse
func (c *Client) GetAttempts() ([]Attempt, error) {
	spath := "/api/attempts"

	params := url.Values{}
	params.Set("project", c.ProjectName)
	params.Set("workflow", c.WorkflowName)

	var attempts *attempts
	err := c.doReq(http.MethodGet, spath, params, &attempts)
	if err != nil {
		return nil, err
	}

	return attempts.Attempts, err
}

// GetLatestAttemptID to get a latest attemptID from sessionDate
func (c *Client) GetLatestAttemptID() (attemptID string, err error) {
	attempts, err := c.GetAttempts()

	// If any attempts not found
	if len(attempts) == 0 {
		err := errors.New("attempts does not exist at `" + c.WorkflowName + "` workflow")
		return attemptID, err
	}

	// c.SessionTime to date like this
	date := c.SessionTime[0:13]

	for k := range attempts {
		sessionTime := attempts[k].SessionTime

		if strings.Contains(sessionTime, date) {
			attemptID = attempts[k].ID
			return attemptID, err
		}
	}

	// If any sesssionTime not found
	err = errors.New("input session " + date + " not found")
	return attemptID, err
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
func (c *Client) GetTaskResult(attemptID, taskName string) (*Task, error) {
	tasks, err := c.GetTasks(attemptID)

	for k := range tasks {
		if tasks[k].FullName == taskName {
			state := tasks[k].State
			if state == "success" {
				return &tasks[k], nil
			}

			err = errors.New("task " + taskName + " state is " + state)
			return nil, err
		}
	}

	return nil, err
}

// CreateNewAttempt to create a new attempt
func (c *Client) CreateNewAttempt(workflowID, date string, retry bool) (attempt *PutAttempt, done bool, err error) {
	spath := "/api/attempts"

	pa := NewPutAttempt(workflowID, c.SessionTime, "")

	// Retry workflow
	if retry == true {
		// TODO: add retry attempt name (optional)
		generatedUUID := uuid.NewV4()
		textID, err := generatedUUID.MarshalText()
		if err != nil {
			return nil, done, err
		}
		pa.RetryAttemptName = string(textID)
	}

	// Create new attempt
	err = c.doReq(http.MethodPut, spath, nil, &pa)
	if err != nil {
		return nil, done, err
	}

	// If alredy attempt done
	if pa.Attempt.Done == true {
		return nil, true, err
	}

	return pa, done, err
}

package digdag

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

type sessions struct {
	Sessions []Session `json:"sessions"`
}

// Session is the struct for digdag session
type Session struct {
	ID      string `json:"id"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"project"`
	Workflow struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"workflow"`
	SessionUUID string    `json:"sessionUuid"`
	SessionTime time.Time `json:"sessionTime"`
	LastAttempt struct {
		ID               string      `json:"id"`
		RetryAttemptName interface{} `json:"retryAttemptName"`
		Done             bool        `json:"done"`
		Success          bool        `json:"success"`
		CancelRequested  bool        `json:"cancelRequested"`
		Params           struct {
		} `json:"params"`
		CreatedAt  time.Time `json:"createdAt"`
		FinishedAt time.Time `json:"finishedAt"`
	}
}

// GetProjectWorkflowSessions to get sessions by projectID and workflow
func (c *Client) GetProjectWorkflowSessions(projectID, workflowName string) ([]Session, error) {
	spath := "/api/projects/" + projectID + "/sessions"

	params := url.Values{}
	params.Set("workflow", workflowName)

	var sessions *sessions
	err := c.doReq(http.MethodGet, spath, params, &sessions)
	if err != nil {
		return nil, err
	}

	// if any sessions not found
	if len(sessions.Sessions) == 0 {
		return nil, errors.New("session not found")
	}

	return sessions.Sessions, err
}

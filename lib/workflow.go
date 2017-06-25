package digdag

import (
	"net/http"
	"net/url"
)

type workflows struct {
	Workflows []Workflow `json:"workflows"`
}

// Workflow is struct for digdag workflow
type Workflow struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Project `json:"project"`
}

// GetWorkflowID to get a workflowID from a projectID
func (c *Client) GetWorkflowID(projectID string) (workflowID string, err error) {
	spath := "/api/projects/" + projectID + "/workflows"

	params := url.Values{}
	params.Set("workflow", c.WorkflowName)

	var workflows *workflows
	err = c.doReq(http.MethodGet, spath, params, &workflows)
	if err != nil {
		return "", err
	}

	workflowID = workflows.Workflows[0].ID

	return workflowID, err
}

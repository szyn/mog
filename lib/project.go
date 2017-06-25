package digdag

import (
	"net/http"
	"net/url"
)

// projects is struct for recevied json
type projects struct {
	Projects []Project `json:"projects"`
}

// Project is struct for digdag project
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetProjectIDByName to get project ID by project name
func (c *Client) GetProjectIDByName() (projectID string, err error) {
	spath := "/api/projects"

	params := url.Values{}
	params.Set("name", c.ProjectName)

	var projects *projects
	err = c.doReq(http.MethodGet, spath, params, &projects)
	if err != nil {
		return "", err
	}

	projectID = projects.Projects[0].ID
	return projectID, nil
}

package util

import (
	"time"

	digdag "github.com/szyn/digdag-go-client"
)

// FetchLocation is to set timezone
func FetchLocation(c *digdag.Client, projectName, workflowName string) (*time.Location, error) {
	project, err := c.GetProject(projectName)
	if err != nil {
		return nil, err
	}

	workflow, err := c.GetWorkflow(project.ID, workflowName)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(workflow.Timezone)
	if err != nil {
		return nil, err
	}

	// Set workflow's timezone
	time.Local = loc

	return loc, nil
}

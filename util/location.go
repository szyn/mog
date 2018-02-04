package util

import (
	"time"

	digdag "github.com/szyn/digdag-go-client"
)

// SetLocation is to set timezone
func SetLocation(c *digdag.Client, projectName, workflowName string) error {
	project, err := c.GetProject(projectName)
	if err != nil {
		return err
	}

	workflow, err := c.GetWorkflow(project.ID, workflowName)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(workflow.Timezone)
	if err != nil {
		return err
	}

	time.Local = loc
	return nil
}

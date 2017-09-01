package main

import (
	"time"

	"errors"
	"fmt"

	"strconv"

	digdag "github.com/szyn/mog/lib"
	"github.com/szyn/mog/logger"
	"github.com/urfave/cli"
)

// Commands is the avalible commands
var Commands = []cli.Command{
	commandStatus,
	commandStart,
	commandRetry,
	commandPollingStatus,
}

var commandStatus = cli.Command{
	Name:      "status",
	Aliases:   []string{"s"},
	Usage:     "Show a status of the task",
	ArgsUsage: "<taskName>",
	Flags:     commonFlag,
	Action:    status,
}

var commandPollingStatus = cli.Command{
	Name:    "polling",
	Aliases: []string{"p"},
	Usage:   "Poll to get a status of the task",
	Subcommands: []cli.Command{
		{
			Name:         "status",
			Usage:        "Poll to get a status of the task",
			ArgsUsage:    "<taskName>",
			Flags:        pollingStatusFlag,
			Action:       pollingStatus,
			OnUsageError: CustomOnUsageError,
		},
		{
			Name:         "trigger",
			Usage:        "Poll to get a status of the task",
			ArgsUsage:    "<taskName>",
			Flags:        pollingStatusFlag,
			Action:       pollingStatus,
			OnUsageError: CustomOnUsageError,
		},
	},
}

var commandStart = cli.Command{
	Name:   "start",
	Usage:  "Start a new session attempt of a workflow",
	Flags:  commonFlag,
	Action: newAttempt,
}

var commandRetry = cli.Command{
	Name:    "retry",
	Aliases: []string{"r"},
	Usage:   "Retry a session",
	Flags:   retryAttemptFlag,
	Action:  newAttempt,
}

// NewClientFromContext
func newClientFromContext(c *cli.Context) *digdag.Client {
	project := c.String("project")
	workflow := c.String("workflow")
	session := c.String("session")

	if workflow == "" {
		err := errors.New("--workflow option")
		if c.App.OnUsageError != nil {
			c.App.OnUsageError(c, err, false)
		} else {
			c.Command.OnUsageError(c, err, true)
		}
	}

	ssl := c.Bool("ssl")
	host := c.GlobalString("host")
	port := strconv.Itoa(c.GlobalInt("port"))

	scheme := "http:"
	if ssl == true {
		scheme = "https:"
	}
	url := scheme + "//" + host + ":" + port

	client, err := digdag.NewClient(url, project, workflow, session, false)
	logger.DieIf(err)

	return client
}

func status(c *cli.Context) error {
	client := newClientFromContext(c)

	task := c.Args().Get(0)
	if task == "" {
		logger.DieIf(errors.New("<taskName> is requied"))
	}
	logger.Log("task: " + task)

	attemptID, err := client.GetLatestAttemptID()
	logger.DieIf(err)

	result, err := client.GetTaskResult(attemptID, task)
	logger.DieIf(err)

	if result == nil {
		logger.DieIf(errors.New("result not found"))
	}

	fmt.Println(prettyPrintJSON(result))

	return nil
}

func pollingStatus(c *cli.Context) error {
	task := c.Args().Get(0)

	if task == "" {
		logger.DieIf(errors.New("<taskName> is requied"))
	}
	logger.Log("task: " + task)

	result := getResult(c)
	resultJSON := prettyPrintJSON(result)

	if resultJSON == "" {
		logger.DieIf(errors.New("result not found"))
	}

	if c.Command.Name == "status" {
		fmt.Println(resultJSON)
	}

	return nil
}

func newAttempt(c *cli.Context) error {
	client := newClientFromContext(c)

	projectID, err := client.GetProjectIDByName()
	logger.DieIf(err)
	logger.Log("projectID: " + projectID)

	workflowID, err := client.GetWorkflowID(projectID)
	logger.DieIf(err)
	logger.Log("workflowID: " + workflowID)

	var retry bool
	retry = c.BoolT("retry")
	logger.Log("retry: " + strconv.FormatBool(retry))

	result, done, err := client.CreateNewAttempt(workflowID, client.SessionTime, retry)
	logger.DieIf(err)

	if done == true {
		msg1 := "A session for the requested session_time already exists.\n"
		msg2 := "`mog retry` to run the session again for the same session_time."
		err := errors.New(msg1 + msg2)
		logger.DieIf(err)
	}

	// Print JSON Response
	fmt.Println(prettyPrintJSON(result))

	return nil
}

func getResult(c *cli.Context) *digdag.Task {
	client := newClientFromContext(c)

	maxTime := c.Int("max-waittime")
	interval := c.Int("interval")
	ticker := time.Tick(time.Duration(interval) * time.Second)
	timeout := time.After(time.Duration(maxTime) * time.Second)

	for {
		select {
		case <-timeout:
			logger.DieIf(fmt.Errorf("wait time exceeded limit at %d sec", maxTime))
		case <-ticker:
			attemptID, err := client.GetLatestAttemptID()
			if err != nil {
				logger.Info(err.Error())
				logger.Info(fmt.Sprintf("state is not success. waiting %d sec for retry...", interval))
				continue
			}

			task := c.Args().Get(0)
			result, err := client.GetTaskResult(attemptID, task)
			if err != nil {
				logger.Info(err.Error())
				logger.Info(fmt.Sprintf("state is not success. waiting %d sec for retry...", interval))
				continue
			}
			if result != nil {
				return result
			}
		}
	}
}

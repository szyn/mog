package main

import (
	"regexp"
	"time"

	"errors"
	"fmt"

	"strconv"

	digdag "github.com/szyn/digdag-go-client"
	"github.com/szyn/mog-kai/util"
	"github.com/szyn/mog/logger"
	"github.com/szyn/mog/util"
	"github.com/urfave/cli"
)

const (
	dailyTimeFormat  = "2006-01-02T00:00:00-07:00"
	hourlyTimeFormat = "2006-01-02T15:00:00-07:00"
	nowTimeFormat    = "2006-01-02T15:04:05-07:00"
)

// Commands is the avalible commands
var Commands = []cli.Command{
	commandStatus,
	commandStart,
	commandRetry,
	commandPollingStatus,
	commandLog,
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

var commandLog = cli.Command{
	Name:    "log",
	Aliases: []string{"l"},
	Usage:   "Show logs of a session",
	Flags:   commonFlag,
	Action:  showLogs,
}

// NewClientFromContext
func newClientFromContext(c *cli.Context) *digdag.Client {
	ssl := c.Bool("ssl")
	host := c.GlobalString("host")
	port := strconv.Itoa(c.GlobalInt("port"))

	scheme := "http:"
	if ssl == true {
		scheme = "https:"
	}
	url := scheme + "//" + host + ":" + port

	client, err := digdag.NewClient(url, c.GlobalBool("verbose"))
	logger.DieIf(err)

	return client
}

func status(c *cli.Context) error {
	client := newClientFromContext(c)
	projectName := c.String("project")
	workflowName := c.String("workflow")

	err := util.SetLocation(client, projectName, workflowName)
	logger.DieIf(err)

	targetSession, err := convertSession(c.String("session"))
	logger.DieIf(err)

	task := c.Args().Get(0)
	if task == "" {
		logger.DieIf(errors.New("<taskName> is requied"))
	}
	logger.Log("task: " + task)

	attemptIDs, err := client.GetAttemptIDs(projectName, workflowName, targetSession)
	logger.DieIf(err)

	result, err := client.GetTaskResult(attemptIDs, task)
	logger.DieIf(err)

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

	if c.Command.Name == "status" {
		fmt.Println(resultJSON)
	}

	return nil
}

func newAttempt(c *cli.Context) error {
	client := newClientFromContext(c)

	project, err := client.GetProject(c.String("project"))
	logger.DieIf(err)
	logger.Log("projectID: " + project.ID)

	workflow, err := client.GetWorkflow(project.ID, c.String("workflow"))
	logger.DieIf(err)
	logger.Log("workflowID: " + workflow.ID)

	var retry bool
	retry = c.BoolT("retry")
	logger.Log("retry: " + strconv.FormatBool(retry))

	err = util.SetLocation(client, project.Name, workflow.Name)
	logger.DieIf(err)

	targetSession, err := convertSession(c.String("session"))
	logger.DieIf(err)

	result, done, err := client.CreateNewAttempt(workflow.ID, targetSession, []string{}, retry)
	if done == true {
		msg1 := "A session for the requested session_time already exists.\n"
		msg2 := "`mog retry` to run the session again for the same session_time."
		err := errors.New(msg1 + msg2)
		logger.DieIf(err)
	}
	logger.DieIf(err)

	// Print JSON Response
	fmt.Println(prettyPrintJSON(result))

	return nil
}

func showLogs(c *cli.Context) error {
	client := newClientFromContext(c)

	task := c.Args().Get(0)
	if task == "" {
		logger.DieIf(errors.New("<taskName> is requied"))
	}
	logger.Log("task: " + task)

	project, err := client.GetProject(c.String("project"))
	logger.DieIf(err)
	logger.Log("projectID: " + project.ID)

	sessions, err := client.GetProjectWorkflowSessions(project.ID, c.String("workflow"))
	logger.DieIf(err)

	lastAttemptID := sessions[0].LastAttempt.ID

	logFile, err := client.GetLogFileResult(lastAttemptID, task)
	logger.DieIf(err)

	logText, err := client.GetLogText(lastAttemptID, logFile.FileName)
	logger.DieIf(err)

	fmt.Println(logText)

	return nil
}

func getResult(c *cli.Context) *digdag.Task {
	client := newClientFromContext(c)

	maxTime := c.Int("max-waittime")
	interval := c.Int("interval")
	ticker := time.Tick(time.Duration(interval) * time.Second)
	timeout := time.After(time.Duration(maxTime) * time.Second)
	projectName := c.String("project")
	workflowName := c.String("workflow")

	err := util.SetLocation(client, projectName, workflowName)
	logger.DieIf(err)

	targetSession, err := convertSession(c.String("session"))
	logger.DieIf(err)

	for {
		select {
		case <-timeout:
			logger.DieIf(fmt.Errorf("wait time exceeded limit at %d sec", maxTime))
		case <-ticker:
			attemptIDs, err := client.GetAttemptIDs(projectName, workflowName, targetSession)
			if err != nil {
				logger.Info(err.Error())
				logger.Info(fmt.Sprintf("state is not success. waiting %d sec for retry...", interval))
				continue
			}

			task := c.Args().Get(0)
			result, err := client.GetTaskResult(attemptIDs, task)
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

func convertSession(session string) (string, error) {
	var sessionTime string

	var daily = regexp.MustCompile(`^[0-9]{4}-[01][0-9]-[0-3][0-9]$`)
	var hourly = regexp.MustCompile(`^[0-9]{4}-[01][0-9]-[0-3][0-9]T[0-9]{2}:[0-9]{2}:[0-9]{2}$`)
	var now = regexp.MustCompile(`^[0-9]{4}-[01][0-9]-[0-3][0-9]T[0-9]{2}:[0-9]{2}:[0-9]{2}(\+|-)[0-9]{2}:[0-9]{2}$`)

	switch {
	case daily.MatchString(session):
		t, err := time.Parse("2006-01-02", session)
		if err != nil {
			return "", err
		}
		sessionTime = t.Format(dailyTimeFormat)
	case hourly.MatchString(session):
		t, err := time.Parse("2006-01-02T15:00:00", session)
		if err != nil {
			return "", err
		}
		sessionTime = t.Format(hourlyTimeFormat)
	case now.MatchString(session):
		t, err := time.Parse(nowTimeFormat, session)
		if err != nil {
			return "", err
		}
		sessionTime = t.Format(nowTimeFormat)
	case session == "daily":
		sessionTime = time.Now().Format(dailyTimeFormat)
	case session == "hourly":
		sessionTime = time.Now().Format(hourlyTimeFormat)
	case session == "now":
		sessionTime = time.Now().Format(nowTimeFormat)
	default:
		return "", fmt.Errorf("Failed to parse input session")
	}

	return sessionTime, nil
}

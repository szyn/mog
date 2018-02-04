package main

import (
	"time"

	"github.com/urfave/cli"
)

var commonFlag = []cli.Flag{
	cli.StringFlag{
		Name:  "project, p",
		Value: "default",
		Usage: "project name",
	},
	cli.StringFlag{
		Name:  "workflow, w",
		Usage: "workflow name",
	},
	cli.StringFlag{
		Name:  "session, s",
		Value: time.Now().Format(dailyTimeFormat),
		Usage: "set session_time to this time",
	},
}

var pollingFlag = []cli.Flag{
	cli.IntFlag{
		Name:  "max-waittime",
		Usage: "wating time(sec)",
		Value: 3600,
	},
	cli.IntFlag{
		Name:  "interval",
		Usage: "polling interval(sec)",
		Value: 60,
	},
}

// TODO: add retry attempt name
var retryFlag = []cli.Flag{
	cli.BoolTFlag{
		Name:  "retry, r",
		Usage: "retry attempts",
	},
}

// pollingStatusFlag for polling status
var pollingStatusFlag = append(commonFlag, pollingFlag...)

// retryAttemptFlag for retry command
var retryAttemptFlag = append(commonFlag, retryFlag...)

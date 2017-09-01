package digdag

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"runtime"

	"github.com/franela/goreq"

	"time"
)

const (
	version          = "0.1" //client-version
	dateTimeFormat   = "2006-01-02"
	dailyTimeFormat  = "2006-01-02T00:00:00-07:00"
	hourlyTimeFormat = "2006-01-02T15:00:00-07:00"
	nowTimeFormat    = "2006-01-02T15:04:05-07:00"
)

// Client api client for digdag
type Client struct {
	URL *url.URL
	http.Client

	Verbose bool

	ProjectName  string
	WorkflowName string
	SessionTime  string
	Date         string
}

// userAgent
var userAgent = fmt.Sprintf("DigdagGoClient/%s (%s)", version, runtime.Version())

// NewClient return new client for digdag
func NewClient(urlStr, project, workflow, session string, verbose bool) (*Client, error) {

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, errors.New("failed to parse url")
	}

	if len(workflow) == 0 {
		return nil, errors.New("missing workflow")
	}
	if len(project) == 0 {
		return nil, errors.New("missing project")
	}
	if len(session) == 0 {
		return nil, errors.New("missing session")
	}

	client := new(Client)
	client.URL = parsedURL
	client.ProjectName = project
	client.WorkflowName = workflow
	client.Verbose = verbose

	r := regexp.MustCompile(`^[0-9]{4}-[01][0-9]-[0-3][0-9]$`).Match([]byte(session))
	if r == true {
		inputSession, err := time.Parse(dateTimeFormat, session)
		if err != nil {
			return nil, err
		}
		client.SessionTime = inputSession.Format(dailyTimeFormat)
		return client, err
	}

	switch session {
	case "daily":
		client.SessionTime = time.Now().Format(dailyTimeFormat)
	case "hourly":
		client.SessionTime = time.Now().Format(hourlyTimeFormat)
	case "now":
		client.SessionTime = time.Now().Format(nowTimeFormat)
	default: // default is dailyTimeFormat
		client.SessionTime = time.Now().Format(dailyTimeFormat)
	}

	return client, err
}

//
func (c *Client) doReq(method, spath string, params, res interface{}) error {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := goreq.Request{
		Method:      method,
		Uri:         u.String(),
		QueryString: params,
		ContentType: "application/json",
		UserAgent:   userAgent,
		Body:        res,
		// ShowDebug:   true,
	}.Do()
	if err != nil {
		return err
	}

	return req.Body.FromJsonTo(&res)
}

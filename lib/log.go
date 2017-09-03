package digdag

import (
	"net/http"

	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type logfiles struct {
	Files []LogFile `json:"files"`
}

type LogFile struct {
	FileName string      `json:"fileName"`
	FileSize int         `json:"fileSize"`
	TaskName string      `json:"taskName"`
	FileTime string      `json:"fileTime"`
	AgentId  string      `json:"agentId"`
	Direct   interface{} `json:"direct"`
}

// GetLogFiles to get logfile list
func (c *Client) GetLogFiles(attemptID string) ([]LogFile, error) {
	spath := "/api/logs/" + attemptID + "/files"

	var logfiles *logfiles
	err := c.doReq(http.MethodGet, spath, nil, &logfiles)
	if err != nil {
		return nil, err
	}

	return logfiles.Files, err
}

// GetLogFileResult to get logfile result
func (c *Client) GetLogFileResult(attemptID, taskName string) (*LogFile, error) {
	logfiles, err := c.GetLogFiles(attemptID)

	for l := range logfiles {
		if logfiles[l].TaskName == taskName {
			return &logfiles[l], nil
		}
	}

	return nil, err
}

// GetLogText to get logtext
func (c *Client) GetLogText(attemptID string, fileName string) (string, error) {
	spath := "/api/logs/" + attemptID + "/files/" + fileName

	gztext, err := c.doRawReq(http.MethodGet, spath, nil)
	if err != nil {
		return "", err
	}

	gr, err := gzip.NewReader(bytes.NewBufferString(gztext))
	defer gr.Close()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(gr)
	return string(data), err
}

package digdag

import (
	"errors"
	"net/http"

	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type logFiles struct {
	Files []LogFile `json:"files"`
}

// LogFile is struct for digdag task log files
type LogFile struct {
	FileName string      `json:"fileName"`
	FileSize int         `json:"fileSize"`
	TaskName string      `json:"taskName"`
	FileTime string      `json:"fileTime"`
	AgentID  string      `json:"agentId"`
	Direct   interface{} `json:"direct"`
}

// GetLogFiles to get logfile list
func (c *Client) GetLogFiles(attemptID string) ([]LogFile, error) {
	spath := "/api/logs/" + attemptID + "/files"

	var logFiles *logFiles
	err := c.doReq(http.MethodGet, spath, nil, &logFiles)
	if err != nil {
		return nil, err
	}

	// if any logFiles not found
	if len(logFiles.Files) == 0 {
		return nil, errors.New("task log not found")
	}

	return logFiles.Files, err
}

// GetLogFileResult to get logfile result
func (c *Client) GetLogFileResult(attemptID, taskName string) (*LogFile, error) {
	logFiles, err := c.GetLogFiles(attemptID)
	if err != nil {
		return nil, err
	}

	for l := range logFiles {
		if logFiles[l].TaskName == taskName {
			return &logFiles[l], nil
		}
		err = errors.New("task log `" + taskName + "` not found")
	}

	return nil, err
}

// GetLogText to get logtext
func (c *Client) GetLogText(attemptID, fileName string) (string, error) {
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

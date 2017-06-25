package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWorkflowID(t *testing.T) {

	testProjectID := "1"
	testWorkflowID := "18"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/projects/"+testProjectID+"/workflows" {
			t.Errorf("request URL should be /api/projects/%v/workflows but : %v", testProjectID, req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/workflows.json`)
		if err != nil {
			t.Error("unexpected error: ", err)
		}

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSONFile))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test", "test", "2017-06-14", false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	workflowID, err := client.GetWorkflowID(testProjectID)
	if err != nil {
		t.Error("err shoud be nil but: ", err)
	}

	if workflowID != testWorkflowID {
		t.Errorf("got %v, want %v", workflowID, testWorkflowID)
	}
}

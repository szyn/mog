package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProjectWorkflowSessions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/projects/1/sessions" {
			t.Error("request URL should be /api/projects/1/sessions but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/sessions.json`)
		if err != nil {
			t.Error("unexpected error: ", err)
		}

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSONFile))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test", "test", "2017-06-24", false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	sessions, err := client.GetProjectWorkflowSessions("1", "test")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("result should be two: %d", len(sessions))
	}
	if sessions[0].ID != "2" {
		t.Fatalf("want %v but %v", "2", sessions[0].ID)
	}
}

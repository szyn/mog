package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProjectIDByName(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/projects" {
			t.Error("request URL should be /api/projects but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/projects.json`)
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

	projectID, err := client.GetProjectIDByName()
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if projectID != "1" {
		t.Fatalf("want %v but %v", "1", projectID)
	}
}

package digdag

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAttempts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts" {
			t.Error("request URL should be /api/attempts but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/attempts.json`)
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

	attempts, err := client.GetAttempts(true)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if len(attempts) != 3 {
		t.Fatalf("result should be two: %d", len(attempts))
	}
	if attempts[0].ID != "27" {
		t.Fatalf("want %v but %v", "27", attempts[0].ID)
	}
}

func TestGetTasks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts/27/tasks" {
			t.Error("request URL should be /api/attempts/27/tasks but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/tasks.json`)
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

	result, err := client.GetTaskResult([]string{"27"}, "+test+setup")
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if result.State != "success" {
		t.Fatalf("want %v but %v", "success", result.State)
	}

	result, err = client.GetTaskResult([]string{"27"}, "+test+failed")
	if err == nil {
		t.Fatalf("should be fail: %v", err)
	}
}

func TestCreateNewAttempt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got '%s'", req.Method)
		}
		if req.URL.Path != "/api/attempts" {
			t.Error("request URL should be /api/attempts but :", req.URL.Path)
		}

		respJSONFile, err := ioutil.ReadFile(`testdata/new_attempt.json`)
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

	result, done, err := client.CreateNewAttempt("2", "2017-06-24", false)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if done != false {
		t.Fatalf("want %v but %v", true, result.Done)
	}

	result, done, err = client.CreateNewAttempt("2", "2017-06-24", true)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}
	if done != false {
		t.Fatalf("want %v but %v", true, result.Done)
	}

}

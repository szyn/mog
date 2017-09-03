package main

import (
	"testing"

	"github.com/rendon/testcli"
)

func goRun(arg ...string) *testcli.Cmd {
	c := testcli.Command("go", append([]string{"run", "format.go", "flag.go", "mog.go", "command.go"}, arg...)...)
	return c
}

func TestHelp(t *testing.T) {
	c := goRun("--help")
	c.Run()
	if !c.Success() {
		t.Fatalf("Expected to succeed, but failed with error: %s", c.Error())
	}

}

func TestVersion(t *testing.T) {
	version := "v0.1.1"

	c := goRun("--version")
	c.Run()
	if !c.Success() {
		t.Fatalf("Expected to succeed, but failed with error: %s", c.Error())
	}

	if !c.StdoutContains(version) {
		t.Fatalf("Expected %q to contain %q", c.Stdout(), version)
	}

}

func TestStatusFailed(t *testing.T) {
	c := goRun("status")
	c.Run()
	if !c.Failure() {
		t.Fatalf("Expected to fail, but succeed: %s", c.Error())
	}
}

func TestStatus(t *testing.T) {

	c := goRun("status", "-s", "2017-06-24", "-w", "test", "+test+setup")
	c.Run()
	if !c.Failure() {
		t.Fatalf("Expected to fail, but succeed: %s", c.Error())
	}

}

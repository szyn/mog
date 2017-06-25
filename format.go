package main

import (
	"encoding/json"
	"strings"

	"github.com/szyn/mog/logger"
)

func prettyPrintJSON(res interface{}) string {
	raw, err := json.MarshalIndent(res, "", "    ")
	logger.DieIf(err)
	return replaceAngreBrackets(string(raw))
}

func replaceAngreBrackets(s string) string {
	s = strings.Replace(s, "\\u003c", "<", -1)
	return strings.Replace(s, "\\u003e", ">", -1)
}

mog
---
[![GitHub release](https://img.shields.io/github/release/szyn/mog.svg?style=flat-square)](https://github.com/szyn/mog/releases/latest)
[![Circle CI](https://img.shields.io/circleci/project/github/szyn/mog.svg?style=flat-square)](https://circleci.com/gh/szyn/mog)
[![Language](https://img.shields.io/badge/language-go-brightgreen.svg?style=flat-square)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/szyn/mog)

mog - A CLI Tool for Digdag.

# Description
mog is a command-line interface tool for the Digdag.  
mog output format is JSON, so it can be filtered with a JSON processor such as jq.

## What's Digdag?
Digdag is an open source Workload Automation System (http://www.digdag.io)

# Installation

```
$ curl https://raw.githubusercontent.com/szyn/mog/master/_tool/get | sh
```

Note:  
Get the latest release of mog.   
The script puts it with Go binaries at `/usr/local/bin`.   
Also you can get binaries at https://github.com/szyn/mog/releases

# Usage

`mog --help` show help.

```console
$ mog --help
NAME:
   mog - A CLI Tool for Digdag

USAGE:
   mog [global options] command [command options] [arguments...]

VERSION:
   v0.1.0

COMMANDS:
     status, s  Show a status of the task
     start,     Start a new session attempt of a workflow
     retry, r   Retry a session
     polling, p  Poll to get a status of the task
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value, -H value  digdag host or ip addr (default: "localhost")
   --port value, -P value  digdag port number (default: 65432)
   --ssl https             make https request
   --verbose               verbose output
   --help, -h              show help
   --version, -v           print the version
```

## Expamles

#### Check a status of the task
Use `mog status`

```console
e.g.
#############################################
# host: localhost (default)
# project: default (default)
# workflow: digdag-sample
# sessionDate: 2017-05-02
# taskName: +digdag-sample+disp_current_date
#############################################

$ mog -H localhost status -w digdag-sample --session 2017-05-02 +digdag-sample+disp_current_date
{
  "config": {
    "echo>": "${moment(session_time).utc().format('YYYY-MM-DD HH:mm:ss Z')}"
  },
  "exportParams": {},
  "fullName": "+digdag-sample+disp_current_date",
  "id": "9",
  "isGroup": false,
  "parentId": "7",
  "retryAt": null,
  "startedAt": "2017-05-02T06:34:05Z",
  "state": "success",
  "stateParams": {},
  "storeParams": {},
  "updatedAt": "2017-05-02T06:34:08Z",
  "upstreams": [
    "8"
  ]
}
```

See also `mog status --help`

#### Polling to check success state of the task
Use `mog polling status`

#### Start a workflow (experimental)
Use `mog start`

```console
e.g.
#############################################
# host: localhost (default)
# project: default (default)
# workflow: digdag-sample
# sessionDate: 2017-05-02
#############################################

$ mog -H localhost start -w digdag-sample --session 2017-05-02
{
  "cancelRequested": false,
  "createdAt": "2017-05-02T06:34:03Z",
  "done": false,
  "finishedAt": null,
  "id": "3",
  "index": 1,
  "params": {},
  "project": {
    "id": "1",
    "name": "default"
  },
  "retryAttemptName": null,
  "sessionId": "3",
  "sessionTime": "2017-05-02T00:00:00+00:00",
  "sessionUuid": "b88a9653-9a34-4763-aa8c-5de213f4826a",
  "success": false,
  "workflow": {
    "id": "4",
    "name": "digdag-sample"
  }
}

```

See also `mog start --help`

#### Retry a workflow (experimental)
Use `mog retry`

```console
e.g.
#############################################
# host: localhost (default)
# project: default (default)
# workflow: digdag-sample
# sessionDate: 2017-05-02
#############################################

$ mog -H localhost retry -w digdag-sample --session 2017-05-02
{
  "cancelRequested": false,
  "createdAt": "2017-05-02T06:34:03Z",
  "done": false,
  "finishedAt": null,
  "id": "4",
  "index": 1,
  "params": {},
  "project": {
    "id": "1",
    "name": "default"
  },
  "retryAttemptName": "47453eb1-a07e-4e8f-a7a5-27d399d0852d",
  "sessionId": "3",
  "sessionTime": "2017-05-02T00:00:00+00:00",
  "sessionUuid": "b88a9653-9a34-4763-aa8c-5de213f4826a",
  "success": false,
  "workflow": {
    "id": "4",
    "name": "digdag-sample"
  }
}

```

See also `mog retry --help`

## Licence

[Apache License 2.0](LICENSE)

## Author

[szyn](https://github.com/szyn)

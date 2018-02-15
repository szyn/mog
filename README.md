mog
---
[![GitHub release](https://img.shields.io/github/release/szyn/mog.svg?style=flat-square)](https://github.com/szyn/mog/releases/latest)
[![Circle CI](https://img.shields.io/circleci/project/github/szyn/mog.svg?style=flat-square)](https://circleci.com/gh/szyn/mog)
[![Language](https://img.shields.io/badge/language-go-brightgreen.svg?style=flat-square)](https://golang.org/)
[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/szyn/mog)
[![Docker Pulls](https://img.shields.io/docker/pulls/szyn/mog.svg?style=flat-square)](https://hub.docker.com/r/szyn/mog/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)

mog - A CLI Tool for Digdag.

# Description
mog is a command-line interface tool for the Digdag.  
mog output format is JSON, so it can be filtered with a JSON processor such as jq.

## What's Digdag?
Digdag is an open source Workload Automation System (https://www.digdag.io)

# Installation

You can download the binary from the [releases](https://github.com/szyn/mog/releases) page.

e.g. os: `linux`, arch: `amd64`  
Download to `/usr/local/bin`
```console 
$ curl -L https://github.com/szyn/mog/releases/download/v0.1.6/mog_linux_amd64.tar.gz | tar zx -C /usr/local/bin
```

### macOS

You can use Homebrew:  

```console
$ brew tap szyn/mog
$ brew install mog
```

### Docker
You can also to use docker image: 

```console
$ docker run --rm szyn/mog:v0.1.6
```

### go get
...Or you can install via go get:

```
$ go get -u github.com/szyn/mog
```

# Usage

`mog --help` show help.

```console
$ mog --help
NAME:
   mog - A CLI Tool for Digdag

USAGE:
   mog [global options] command [command options] [arguments...]

VERSION:
   v0.1.6

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

## Examples

#### Check a status of the task
Use `mog status`

```console
e.g.
#############################################
# host: localhost (default)
# project: sample
# workflow: digdag-sample
# sessionTime: 2017-10-08T15:00:00+00:00 (UTC)
# taskName: +digdag-sample+disp_current_date
#############################################

$ mog -H localhost status -p sample -w digdag-sample --session 2017-10-08T15:00:00+00:00 +digdag-sample+disp_current_date
{
    "id": "41",
    "fullName": "+digdag-sample+disp_current_date",
    "parentId": "39",
    "config": {
        "echo>": "${moment(session_time).utc().format('YYYY-MM-DD HH:mm:ss Z')}"
    },
    "upstreams": [
        "40"
    ],
    "state": "success",
    "exportParams": {},
    "storeParams": {},
    "stateParams": {},
    "updatedAt": "2017-10-09T14:50:26Z",
    "retryAt": null,
    "startedAt": "2017-10-09T14:50:26Z",
    "isGroup": false
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
# project: sample
# workflow: digdag-sample
# sessionTime: 2017-10-08T15:00:00+00:00 (UTC)
#############################################

$ mog -H localhost start -p sample -w digdag-sample --session 2017-10-09
{
    "id": "5",
    "index": 1,
    "project": {
        "id": "2",
        "name": "sample"
    },
    "workflow": {
        "name": "digdag-sample",
        "id": "3"
    },
    "sessionId": "3",
    "sessionUuid": "948a9083-095c-4eea-b910-d63763006de7",
    "done": false,
    "success": false,
    "cancelRequested": false,
    "createdAt": "2017-10-09T14:50:03Z",
    "finishedAt": "",
    "workflowId": "3",
    "sessionTime": "2017-10-08T15:00:00+00:00",
    "params": {}
}
```

See also `mog start --help`

#### Retry a workflow (experimental)
Use `mog retry`

```console
e.g.
#############################################
# host: localhost (default)
# project: sample
# workflow: digdag-sample
# sessionTime: 2017-10-08T15:00:00+00:00 (UTC)
#############################################

$ mog -H localhost retry -p sample -w digdag-sample --session 2017-10-09
{
    "id": "6",
    "index": 2,
    "project": {
        "id": "2",
        "name": "sample"
    },
    "workflow": {
        "name": "digdag-sample",
        "id": "3"
    },
    "sessionId": "3",
    "sessionUuid": "948a9083-095c-4eea-b910-d63763006de7",
    "done": false,
    "success": false,
    "cancelRequested": false,
    "createdAt": "2017-10-09T14:50:26Z",
    "finishedAt": "",
    "workflowId": "3",
    "sessionTime": "2017-10-08T15:00:00+00:00",
    "retryAttemptName": "f01529fd-fc2c-4f77-b6c5-f484321e2001",
    "params": {}
}
```

See also `mog retry --help`

## Licence

[Apache License 2.0](LICENSE)

## Author

[szyn](https://github.com/szyn)

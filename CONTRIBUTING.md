# How to Contribute

1. Fork it ( http://github.com/szyn/mog )
1. Create your feature branch (git checkout -b my-awesome-feature)
1. Commit your changes (git commit -am 'Add awesome feature')
1. Rebase your local changes against the master branch
1. Push to the branch (git push origin my-awesome-feature)
1. Run test suite with the `go test $(glide novendor)` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request

---

# Development

## Requirement
- Docker 1.13 or later

## build

```
$ docker build -t mog-dev .
$ docker run --rm -v ${PWD}:/go/src/github.com/szyn/mog mog-dev
```
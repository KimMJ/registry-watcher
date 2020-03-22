module github.com/kimmj/registry-watcher/src

go 1.13

replace github.com/kimmj/registry-watcher => ../

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/goharbor/harbor/src v0.0.0-20200321042307-e1a1e4d1217d
	github.com/robfig/cron/v3 v3.0.0 // indirect
)

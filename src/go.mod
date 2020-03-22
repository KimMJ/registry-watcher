module github.com/kimmj/registry-watcher/src

go 1.13

replace github.com/kimmj/registry-watcher => ../

require (
	github.com/gin-gonic/gin v1.5.0
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/robfig/cron/v3 v3.0.0
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/sys v0.0.0-20190916202348-b4ddaad3f8a3 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

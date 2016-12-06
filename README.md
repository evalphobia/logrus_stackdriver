logrus_stackdriver
====

[![Build Status](https://travis-ci.org/evalphobia/logrus_stackdriver.svg?branch=master)](https://travis-ci.org/evalphobia/logrus_stackdriver) [![Coverage Status](https://coveralls.io/repos/evalphobia/logrus_stackdriver/badge.svg?branch=master&service=github)](https://coveralls.io/github/evalphobia/logrus_stackdriver?branch=master) [![codecov](https://codecov.io/gh/evalphobia/logrus_stackdriver/branch/master/graph/badge.svg)](https://codecov.io/gh/evalphobia/logrus_stackdriver)
 [![GoDoc](https://godoc.org/github.com/evalphobia/logrus_stackdriver?status.svg)](https://godoc.org/github.com/evalphobia/logrus_stackdriver)


# Google Stackdriver logging Hook for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Usage

```go
import (
    "github.com/Sirupsen/logrus"
    "github.com/evalphobia/google-api-go-wrapper/config"
    "github.com/evalphobia/logrus_stackdriver"
)

func main() {
    hook, err := logrus_stackdriver.NewWithConfig("project_id", "test_log", config.Config{
        Email:      "xxx@xxx.iam.gserviceaccount.com",
        PrivateKey: "-----BEGIN PRIVATE KEY-----\nXXX\n-----END PRIVATE KEY-----\n",
    })

    // set custom fire level
    hook.SetLevels([]logrus.Level{
        logrus.PanicLevel,
        logrus.ErrorLevel,
        logrus.WarnLevel,
    })

    // ignore field
    hook.AddIgnore("context")

    // add custome filter
    hook.AddFilter("error", logrus_stackdriver.FilterError)


    // send log with logrus
    logger := logrus.New()
    logger.Hooks.Add(hook)
    logger.WithFields(f).Error("my_message") // send log data to Google Stackdriver logging API
}
```


## Special fields

Some logrus fields have a special meaning in this hook.

| Field Name | Description |
|:--|:--|
|`message`|if `message` is not set, entry.Message is added to log data in "message" field. |
|`log_name`|`log_name` is a custom log name. If not set, `defaultLogName` is used as log name.|
|`http_request`|`http_request` is *http.Request for detailed http logging.|
|`http_response`|`http_response` is *http.Response for detailed http logging.|

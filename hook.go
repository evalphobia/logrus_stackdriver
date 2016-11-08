package logrus_stackdriver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/stackdriver/logging"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

// StackdriverHook is logrus hook for Google Stackdriver.
type StackdriverHook struct {
	client *logging.Logger

	defaultLogName string
	commonLabels   map[string]string
	async          bool
	levels         []logrus.Level
	ignoreFields   map[string]struct{}
	filters        map[string]func(interface{}) interface{}
}

// New returns initialized logrus hook for Stackdriver.
func New(projectID string, logName string) (*StackdriverHook, error) {
	return NewWithConfig(projectID, logName, config.Config{})
}

// NewWithConfig returns initialized logrus hook for Stackdriver.
func NewWithConfig(projectID string, logName string, conf config.Config) (*StackdriverHook, error) {
	logger, err := logging.NewLogger(conf, projectID)
	if err != nil {
		return nil, err
	}

	return &StackdriverHook{
		client:         logger,
		defaultLogName: logName,
		levels:         defaultLevels,
		ignoreFields:   make(map[string]struct{}),
		filters:        make(map[string]func(interface{}) interface{}),
	}, nil
}

// Levels returns logging level to fire this hook.
func (h *StackdriverHook) Levels() []logrus.Level {
	return h.levels
}

// SetLevels sets logging level to fire this hook.
func (h *StackdriverHook) SetLevels(levels []logrus.Level) {
	h.levels = levels
}

// SetLabels sets logging level to fire this hook.
func (h *StackdriverHook) SetLabels(labels map[string]string) {
	h.commonLabels = labels
}

// Async sets async flag and send log asynchroniously.
// If use this option, Fire() does not return error.
func (h *StackdriverHook) Async() {
	h.async = true
}

// AddIgnore adds field name to ignore.
func (h *StackdriverHook) AddIgnore(name string) {
	h.ignoreFields[name] = struct{}{}
}

// AddFilter adds a custom filter function.
func (h *StackdriverHook) AddFilter(name string, fn func(interface{}) interface{}) {
	h.filters[name] = fn
}

// Fire is invoked by logrus and sends log to kinesis.
func (h *StackdriverHook) Fire(entry *logrus.Entry) error {
	if !h.async {
		return h.fire(entry)
	}

	// send log asynchroniously and return no error.
	go h.fire(entry)
	return nil
}

// Fire is invoked by logrus and sends log to kinesis.
func (h *StackdriverHook) fire(entry *logrus.Entry) error {
	req := getRequest(entry.Data)
	resp := getResponse(entry.Data)

	return h.client.Write(logging.WriteData{
		Labels:   h.commonLabels,
		Severity: getSeverity(entry.Level),
		LogName:  h.getLogName(entry.Data),
		Data:     h.getData(entry),
		Request:  req,
		Response: resp,
		Resource: &logging.Resource{
			Type: "global",
		},
	})
}

func (h *StackdriverHook) getLogName(data logrus.Fields) string {
	if name, ok := data["log_name"].(string); ok {
		return name
	}
	return h.defaultLogName
}

func (h *StackdriverHook) getData(entry *logrus.Entry) map[string]interface{} {
	if _, ok := entry.Data["message"]; !ok {
		entry.Data["message"] = entry.Message
	}

	data := make(logrus.Fields)
	for k, v := range entry.Data {
		if _, ok := h.ignoreFields[k]; ok {
			continue
		}
		if fn, ok := h.filters[k]; ok {
			v = fn(v) // apply custom filter
		} else {
			v = formatData(v) // use default formatter
		}
		data[k] = v
	}
	return data
}

// formatData returns value as a suitable format.
func formatData(value interface{}) (formatted interface{}) {
	switch value := value.(type) {
	case json.Marshaler:
		return value
	case error:
		return value.Error()
	case fmt.Stringer:
		return value.String()
	default:
		return value
	}
}

func getRequest(data logrus.Fields) *http.Request {
	if req, ok := data["http_request"].(*http.Request); ok {
		delete(data, "http_request")
		return req
	}
	return nil
}

func getResponse(data logrus.Fields) *http.Response {
	if resp, ok := data["http_response"].(*http.Response); ok {
		delete(data, "http_response")
		return resp
	}
	return nil
}

func getSeverity(level logrus.Level) logging.Severity {
	switch level {
	case logrus.DebugLevel:
		return logging.SeverityDebug
	case logrus.InfoLevel:
		return logging.SeverityInfo
	case logrus.WarnLevel:
		return logging.SeverityWarning
	case logrus.ErrorLevel:
		return logging.SeverityError
	case logrus.PanicLevel:
		return logging.SeverityCritical
	case logrus.FatalLevel:
		return logging.SeverityAlert
	default:
		return logging.SeverityDefault
	}
}

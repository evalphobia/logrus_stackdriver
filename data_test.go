package logrus_stackdriver

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		fieldSize int
	}{
		{0},   // empty fileds
		{1},   // "0"
		{2},   // "0", "1"
		{9},   // "0", "1", "2" ... "8"
		{100}, // "0", "1", "2" ... "99"
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		fields := logrus.Fields{}
		for i, max := 0, tt.fieldSize; i < max; i++ {
			fields[strconv.Itoa(i)] = struct{}{}
		}

		df := dataField{
			data: fields,
		}
		assert.Equal(tt.fieldSize, df.len(), "dataField.Len() should equal fieldSize", target)
	}
}

func TestIsOmit(t *testing.T) {
	assert := assert.New(t)

	omitList := map[string]struct{}{
		"key_1": struct{}{},
		"key_2": struct{}{},
		"key_3": struct{}{},
		"key_4": struct{}{},
	}

	tests := []struct {
		key      string
		expected bool
	}{
		{"key_1", true},
		{"key_2", true},
		{"key_3", true},
		{"key_4", true},
		{"not_key", false},
		{"foo", false},
		{"bar", false},
		{"_key_1", false},
		{"key_1_", false},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		df := dataField{
			omitList: omitList,
		}
		assert.Equal(tt.expected, df.isOmit(tt.key), target)
	}
}

func TestGetLogName(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key         string
		value       interface{}
		expected    bool
		description string
	}{
		{"log_name", "test_log_name", true, "valid server name"},
		{"log_name", "", true, "valid server name"},
		{"not_log_name", "test_log_name", false, "invalid key"},
		{"log_name", 1, false, "invalid value type"},
		{"log_name", true, false, "invalid value type"},
		{"log_name", struct{}{}, false, "invalid value type"},
	}

	const defaultLogName = "default_log_name"
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		fields := logrus.Fields{}
		fields[tt.key] = tt.value
		entry := &logrus.Entry{
			Data: fields,
		}

		df := newDataFieldFromEntry(defaultLogName, entry)
		logName := df.getLogName()
		if tt.expected {
			assert.Equal(tt.value, logName, target)
		} else {
			assert.Equal(defaultLogName, logName, target)
		}
	}
}

func TestGetRequest(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key         string
		value       interface{}
		expected    bool
		description string
	}{
		{"http_request", &http.Request{}, true, "valid http_request"},
		{"not_http_request", &http.Request{}, false, "invalid key"},
		{"http_request", http.Request{}, false, "invalid value type"},
		{"http_request", "test_http_request", false, "invalid value type"},
		{"http_request", 1, false, "invalid value type"},
		{"http_request", true, false, "invalid value type"},
		{"http_request", struct{}{}, false, "invalid value type"},
	}

	const defaultLogName = "default_log_name"
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		fields := logrus.Fields{}
		fields[tt.key] = tt.value
		entry := &logrus.Entry{
			Data: fields,
		}

		df := newDataFieldFromEntry(defaultLogName, entry)
		req := df.getRequest()
		if tt.expected {
			assert.Equal(tt.value, req, target)
			assert.True(df.isOmit("http_request"), "`http_request` should be in omitList")
		} else {
			assert.Nil(req, target)
			assert.False(df.isOmit("http_request"), "`http_request` should not be in omitList")
		}
	}
}

func TestGetResponse(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key         string
		value       interface{}
		expected    bool
		description string
	}{
		{"http_response", &http.Response{}, true, "valid http_response"},
		{"not_http_response", &http.Response{}, false, "invalid key"},
		{"http_response", http.Response{}, false, "invalid value type"},
		{"http_response", "test_http_response", false, "invalid value type"},
		{"http_response", 1, false, "invalid value type"},
		{"http_response", true, false, "invalid value type"},
		{"http_response", struct{}{}, false, "invalid value type"},
	}

	const defaultLogName = "default_log_name"
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)

		fields := logrus.Fields{}
		fields[tt.key] = tt.value
		entry := &logrus.Entry{
			Data: fields,
		}

		df := newDataFieldFromEntry(defaultLogName, entry)
		resp := df.getResponse()
		if tt.expected {
			assert.Equal(tt.value, resp, target)
			assert.True(df.isOmit("http_response"), "`http_response` should be in omitList")
		} else {
			assert.Nil(resp, target)
			assert.False(df.isOmit("http_response"), "`http_response` should not be in omitList")
		}
	}
}

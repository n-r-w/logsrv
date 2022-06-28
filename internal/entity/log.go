// Package entity ...
package entity

import (
	"encoding/json"
	"time"
)

// LogRecord - запись в журнале
type LogRecord struct {
	ID          uint64            `json:"id"`
	RecordTime  time.Time         `json:"recordTime"`
	LogTime     time.Time         `json:"logTime,omitempty"`
	Service     string            `json:"service,omitempty"`
	Source      string            `json:"source,omitempty"`
	Category    string            `json:"category,omitempty"`
	Level       string            `json:"level,omitempty"`
	Session     string            `json:"session,omitempty"`
	Info        string            `json:"info,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	Url         string            `json:"url,omitempty"`
	HttpType    string            `json:"httpType,omitempty"`
	HttpCode    int               `json:"httpCode,omitempty"`
	ErrorCode   int               `json:"errorCode,omitempty"`
	HttpHeaders map[string]string `json:"httpHeaders,omitempty"`
	Body        json.RawMessage   `json:"body,omitempty"`
}

// IsEmpty ...
func (l *LogRecord) IsEmpty() bool {
	return l.RecordTime.IsZero()
}

// Validate ...
func (l *LogRecord) Validate() error {
	return nil
}

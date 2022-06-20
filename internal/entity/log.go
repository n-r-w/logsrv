// Package entity ...
package entity

import (
	"time"
)

// LogRecord - запись в журнале
type LogRecord struct {
	ID          uint64            `json:"id"`
	RecordTime  time.Time         `json:"recordTime"`
	LogTime     time.Time         `json:"logTime"`
	Service     string            `json:"service"`
	Source      string            `json:"source"`
	Category    string            `json:"category"`
	Level       string            `json:"level"`
	Session     string            `json:"session"`
	Info        string            `json:"info"`
	Properties  map[string]string `json:"properties"`
	Url         string            `json:"url"`
	HttpType    string            `json:"httpType"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	Body        interface{}       `json:"body"`
}

// IsEmpty ...
func (l *LogRecord) IsEmpty() bool {
	return l.RecordTime.IsZero()
}

// Validate ...
func (l *LogRecord) Validate() error {
	return nil
}

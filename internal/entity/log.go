// Package entity ...
package entity

import (
	"time"
)

// LogRecord - запись в журнале
type LogRecord struct {
	ID          uint64            `json:"id"`
	Time        time.Time         `json:"time"`
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
	return l.Time.IsZero()
}

// Validate ...
func (l *LogRecord) Validate() error {
	return nil
}

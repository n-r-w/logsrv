package entity

import (
	"encoding/json"
	"time"
)

type SearchRequest struct {
	Criteria []SearchCriteria `json:"criteria"`
	AndOp    bool             `json:"and"`
}

// SearchCriteria - критерии поиска
type SearchCriteria struct {
	AndOp       bool                       `json:"and"`
	From        time.Time                  `json:"from"`
	To          time.Time                  `json:"to"`
	Service     string                     `json:"service"`
	Source      string                     `json:"source"`
	Category    string                     `json:"category"`
	Level       string                     `json:"level"`
	Session     string                     `json:"session"`
	Info        string                     `json:"info"`
	Properties  map[string]string          `json:"properties"`
	Url         string                     `json:"url"`
	HttpType    string                     `json:"httpType"`
	HttpCode    int                        `json:"httpCode"`
	ErrorCode   int                        `json:"errorCode"`
	HttpHeaders map[string]string          `json:"httpHeaders"`
	BodyValues  map[string]json.RawMessage `json:"bodyValues"`
	Body        json.RawMessage            `json:"body"`
}

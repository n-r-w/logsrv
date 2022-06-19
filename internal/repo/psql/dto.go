package psql

import "github.com/jackc/pgtype"

type LogRecord struct {
	ID          uint64            `json:"id"`
	Time        pgtype.Timestamp  `json:"time"`
	Service     pgtype.Text       `json:"service"`
	Source      pgtype.Text       `json:"source"`
	Category    pgtype.Text       `json:"category"`
	Level       pgtype.Text       `json:"level"`
	Session     pgtype.Text       `json:"session"`
	Url         pgtype.Text       `json:"url"`
	HttpType    pgtype.Text       `json:"httpType"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	JsonBody    interface{}       `json:"jsonBody"`
}

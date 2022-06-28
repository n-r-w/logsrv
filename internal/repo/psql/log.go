// Package psql Содержит реализацию интерфейса репозитория логов для postgresql
package psql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/n-r-w/eno"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/entity"
	"github.com/n-r-w/nerr"
	"github.com/n-r-w/postgres"
	"github.com/n-r-w/sqlb"
	"github.com/n-r-w/sqlq"
	"github.com/n-r-w/tools"
)

type LogRepo struct {
	logger lg.Logger
	*postgres.Postgres
	config *config.Config
}

func NewLog(pg *postgres.Postgres, logger lg.Logger, config *config.Config) *LogRepo {
	return &LogRepo{
		logger:   logger,
		Postgres: pg,
		config:   config,
	}
}

func (p *LogRepo) Size() int {
	return int(p.Pool.Stat().TotalConns())
}

func (p *LogRepo) Insert(records []entity.LogRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(p.config.DbWriteTimeout))
	defer cancel()

	tx := sqlq.NewTx(p.Pool, ctx)
	if err := tx.Begin(); err != nil {
		return err
	}
	defer tx.Rollback()

	var valuesSql []string

	for _, lr := range records {
		if err := lr.Validate(); err != nil {
			return err
		}

		bnd := sqlb.NewBinder("(:log_time, :service, :source, :category, :level, :session, :info, :url, :http_type, :http_code, :error_code, :json_body)", "InsertLogs")
		if err := bnd.BindValues(map[string]interface{}{
			"log_time":   lr.LogTime,
			"service":    sqlb.VNull(lr.Service),
			"source":     sqlb.VNull(lr.Source),
			"category":   sqlb.VNull(lr.Category),
			"level":      sqlb.VNull(lr.Level),
			"session":    sqlb.VNull(lr.Session),
			"info":       sqlb.VNull(lr.Info),
			"url":        sqlb.VNull(lr.Url),
			"http_type":  sqlb.VNull(lr.HttpType),
			"http_code":  sqlb.VNull(lr.HttpCode),
			"error_code": sqlb.VNull(lr.ErrorCode),
			"json_body":  sqlb.VNull(lr.Body),
		}); err != nil {
			return err
		}

		sql, err := bnd.Sql()
		if err != nil {
			return err
		}

		valuesSql = append(valuesSql, sql)
	}

	sqlText := fmt.Sprintf(
		"INSERT INTO logs (log_time, service, source, category, level, session, info, url, http_type, http_code, error_code, json_body) VALUES %s RETURNING id", strings.Join(valuesSql, ","))

	q, err := sqlq.SelectTx(tx, sqlText)
	if err != nil {
		fmt.Println(sqlText)
		return nerr.New(err, tools.SimplifyString(sqlText))
	}
	defer q.Close()

	var ids []uint64
	for q.Next() {
		ids = append(ids, q.UInt64("id"))
	}

	if len(ids) != len(records) {
		return nerr.New(eno.ErrInternal, "len(ids) != len(records)")
	}

	// заголовки http
	var headersSql []string
	for i, id := range ids {
		for key, value := range records[i].HttpHeaders {
			bnd := sqlb.NewBinder("(:record_id, :header_name, :header_value)", "InsertHeaders")
			if err := bnd.BindValues(map[string]interface{}{
				"record_id":    id,
				"header_name":  key,
				"header_value": value,
			}); err != nil {
				return err
			}

			s, err := bnd.Sql()
			if err != nil {
				return err
			}
			headersSql = append(headersSql, s)
		}
	}
	if len(headersSql) > 0 {
		sqlText = fmt.Sprintf("INSERT INTO http_headers (record_id, header_name, header_value) VALUES %s", strings.Join(headersSql, ","))
		if err = sqlq.ExecTx(tx, sqlText); err != nil {
			return nerr.New(err, tools.SimplifyString(sqlText))
		}
	}

	// свойства
	var propsSql []string
	for i, id := range ids {
		for key, value := range records[i].Properties {
			bnd := sqlb.NewBinder("(:record_id, :p_name, :p_value)", "InsertProps")
			if err := bnd.BindValues(map[string]interface{}{
				"record_id": id,
				"p_name":    key,
				"p_value":   value,
			}); err != nil {
				return err
			}

			s, err := bnd.Sql()
			if err != nil {
				return err
			}
			propsSql = append(propsSql, s)
		}
	}
	if len(propsSql) > 0 {
		sqlText = fmt.Sprintf("INSERT INTO properties (record_id, p_name, p_value) VALUES %s", strings.Join(propsSql, ","))
		if err = sqlq.ExecTx(tx, sqlText); err != nil {
			return nerr.New(err, tools.SimplifyString(sqlText))
		}
	}

	return tx.Commit()
}

func (p *LogRepo) Find(request entity.SearchRequest) (records []entity.LogRecord, limited bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(p.config.DbReadTimeout))
	defer cancel()

	if len(request.Criteria) == 0 {
		return []entity.LogRecord{}, false, nil
	}

	var opGlobal string
	if request.AndOp {
		opGlobal = " AND "
	} else {
		opGlobal = " OR "
	}

	var whereSql []string
	for _, c := range request.Criteria {
		var critSql []string

		var op string
		if c.AndOp {
			op = " AND "
		} else {
			op = " OR "
		}

		if !c.From.IsZero() || !c.To.IsZero() {
			var timeSql []string
			if !c.From.IsZero() {
				if v, err := sqlb.ToSql(c.From); err != nil {
					return nil, false, err
				} else {
					timeSql = append(timeSql, fmt.Sprintf("l.record_time >= %s", v))
				}
			}
			if !c.To.IsZero() {
				if v, err := sqlb.ToSql(c.To); err != nil {
					return nil, false, err
				} else {
					timeSql = append(timeSql, fmt.Sprintf("l.record_time <= %s", v))
				}
			}
			critSql = append(critSql, fmt.Sprintf("(%s)", strings.Join(timeSql, " AND ")))
		}

		if len(c.Service) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.service='%s'", c.Service))
		}
		if len(c.Source) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.source='%s'", c.Source))
		}
		if len(c.Category) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.category='%s'", c.Category))
		}
		if len(c.Level) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.level='%s'", c.Level))
		}
		if len(c.Session) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.session='%s'", c.Session))
		}
		if len(c.Url) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.url~'%s'", c.Url)) // регулярка
		}
		if len(c.Info) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.info~'%s'", c.Info)) // регулярка
		}
		if len(c.HttpType) > 0 {
			critSql = append(critSql, fmt.Sprintf("l.http_type='%s'", c.HttpType))
		}
		if c.HttpCode > 0 {
			critSql = append(critSql, fmt.Sprintf("l.http_code=%d", c.HttpCode))
		}
		if c.ErrorCode > 0 {
			critSql = append(critSql, fmt.Sprintf("l.error_code=%d", c.ErrorCode))
		}

		for key, value := range c.BodyValues {
			v, err := sqlb.ToSql(value, sqlb.JsonPath)
			if err != nil {
				return nil, false, err
			}
			critSql = append(critSql, fmt.Sprintf(`jsonb_path_exists(l.json_body, '$.** ? (@.%s == %s)')`, key, v))
		}

		if c.Body != nil && len(c.Body) > 0 {
			s, err := sqlb.ToSql(c.Body, sqlb.Json)
			if err != nil {
				return nil, false, err
			}
			if len(s) > 0 {
				sqlb.ToSql(c.Body)
				critSql = append(critSql, fmt.Sprintf(`json_body @> %s`, s))
			}
		}

		if len(c.HttpHeaders) > 0 {
			var headersSql []string
			for key, value := range c.HttpHeaders {
				headersSql = append(headersSql, fmt.Sprintf("(h.header_name = '%s' and h.header_value = '%s')", key, value))
			}

			critSql = append(critSql, fmt.Sprintf("exists(SELECT * FROM http_headers h WHERE h.record_id = l.id AND (%s))",
				strings.Join(headersSql, op)))
		}

		if len(c.Properties) > 0 {
			var propsSql []string
			for key, value := range c.Properties {
				propsSql = append(propsSql, fmt.Sprintf("(h.p_name = '%s' and h.p_value = '%s')", key, value))
			}

			critSql = append(critSql, fmt.Sprintf("exists(SELECT * FROM properties h WHERE h.record_id = l.id AND (%s))",
				strings.Join(propsSql, op)))
		}

		if len(critSql) > 0 {
			whereSql = append(whereSql, fmt.Sprintf("(%s)", strings.Join(critSql, op)))
		}
	}

	if len(whereSql) == 0 {
		return nil, false, nerr.New(eno.ErrBadRequest, "no criteria")
	}

	sql := fmt.Sprintf(
		`select id, record_time, log_time, service, source, category, level, session, info, url, http_type, http_code, error_code, json_body,
			(ARRAY(
			SELECT header_name
			FROM http_headers h1
			WHERE h1.record_id = l.id
			)) as header_names,
			(ARRAY(
			SELECT header_value
			FROM http_headers h2
			WHERE h2.record_id = l.id
			)) as header_values,
			(ARRAY(
			SELECT p_name
			FROM properties p1
			WHERE p1.record_id = l.id
			)) as p_names,
			(ARRAY(
			SELECT p_value
			FROM properties p2
			WHERE p2.record_id = l.id
			)) as p_values
		from logs l where %s order by log_time`,
		strings.Join(whereSql, opGlobal))

	q, err := sqlq.Select(p.Pool, ctx, sql)
	if err != nil {
		return nil, false, nerr.New(err, tools.SimplifyString(sql))
	}
	defer q.Close()

	var recs []entity.LogRecord

	var rowCount uint64
	limited = false

	for q.Next() {
		body := q.Json("json_body")
		// пробуем его запарсить, иначе это вылезет потом и мы не сможем отправить ответ
		if _, err := json.Marshal(body); err != nil {
			body = json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
			p.logger.Debug("%s", string(body))
		}

		lr := entity.LogRecord{
			ID:          q.UInt64("id"),
			RecordTime:  q.Time("record_time"),
			LogTime:     q.Time("log_time"),
			Service:     q.String("service"),
			Source:      q.String("source"),
			Category:    q.String("category"),
			Level:       q.String("level"),
			Session:     q.String("session"),
			Info:        q.String("info"),
			Properties:  map[string]string{},
			Url:         q.String("url"),
			HttpType:    q.String("http_type"),
			HttpCode:    q.Int("http_code"),
			ErrorCode:   q.Int("error_code"),
			HttpHeaders: map[string]string{},
			Body:        body,
		}

		headerNames := q.StringArray("header_names")
		headerValues := q.StringArray("header_values")
		if len(headerNames) != len(headerValues) {
			return nil, false, nerr.New(eno.ErrInternal, "invalid headers count")
		}

		for i, h := range headerNames {
			lr.HttpHeaders[h] = headerValues[i]
		}

		pNames := q.StringArray("p_names")
		pValues := q.StringArray("p_values")
		if len(pNames) != len(pValues) {
			return nil, false, nerr.New(eno.ErrInternal, "invalid headers count")
		}

		for i, h := range pNames {
			lr.Properties[h] = pValues[i]
		}

		rowCount++
		if rowCount > uint64(p.config.MaxLogRecordsResult) {
			limited = true

			break
		}

		if rowCount > uint64(p.config.MaxLogRecordsResult) {
			err := fmt.Errorf("too many records, max %d", p.config.MaxLogRecordsResult)

			return nil, false, err
		}

		recs = append(recs, lr)
	}

	return recs, limited, nil
}

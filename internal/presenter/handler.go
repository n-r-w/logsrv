package presenter

import (
	"encoding/json"
	"net/http"

	"github.com/n-r-w/httprouter"
	"github.com/n-r-w/logsrv/internal/entity"
	"github.com/n-r-w/nerr"
)

// Добавить в лог
func (p *Presenter) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []entity.LogRecord
		// парсим входящий json
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p.controller.RespondError(w, http.StatusBadRequest, nerr.New(err))

			return
		}

		if err := p.logRepo.Insert(req); err != nil {
			p.controller.RespondError(w, http.StatusForbidden, err)

			return
		}

		p.controller.RespondData(w, http.StatusCreated, "application/json; charset=utf-8", nil)
	}
}

// Получить записи из лога
func (p *Presenter) search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := entity.SearchRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			p.controller.RespondError(w, http.StatusBadRequest, err)

			return
		}

		records, _, err := p.logRepo.Find(req)
		if err != nil {
			p.controller.RespondError(w, http.StatusInternalServerError, err)

			return
		}

		if len(records) == 0 {
			p.controller.RespondData(w, http.StatusOK, "application/json; charset=utf-8", nil)

			return
		}

		// отдаем с gzip сжатием если клиент это желает
		p.controller.RespondCompressed(w, r, http.StatusOK, httprouter.CompressionGzip, "application/json; charset=utf-8", &records)
	}
}

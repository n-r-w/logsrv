package rest

import (
	"net/http"

	"github.com/n-r-w/eno"
	"github.com/n-r-w/httprouter"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/presenter"
	"github.com/n-r-w/nerr"
)

type Service struct {
	controller httprouter.Router
	logRepo    presenter.LogInterface
	config     *config.Config
	presenter  *presenter.Service
}

// New Инициализация маршрутов
func New(router httprouter.Router, logRepo presenter.LogInterface, config *config.Config, presenter *presenter.Service) *Service {
	p := &Service{
		controller: router,
		logRepo:    logRepo,
		config:     config,
		presenter:  presenter,
	}

	// устанавливаем middleware для проверки валидности сессии
	router.AddMiddleware("/api", p.authenticateUser)

	// добавить запись в лог
	router.AddRoute("/api", "/add", p.add(), "POST")
	// поиск записей в логе
	router.AddRoute("/api", "/search", p.search(), "POST")

	return p
}

// Аутентификация пользователя
func (p *Service) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Authorization")
		if p.presenter.CheckRights(token, true, true) != nil {
			p.controller.RespondError(w, http.StatusUnauthorized, nerr.New(eno.ErrNoAccess))
			return
		}

		next.ServeHTTP(w, r)
	})
}

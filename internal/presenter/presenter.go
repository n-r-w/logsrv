package presenter

import (
	"net/http"

	"github.com/n-r-w/eno"
	"github.com/n-r-w/httprouter"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/nerr"
)

type Presenter struct {
	controller httprouter.Router
	logRepo    LogInterface
	config     *config.Config

	tokens map[string]bool // список токенов доступа
}

// New Инициализация маршрутов
func New(router httprouter.Router, logRepo LogInterface, config *config.Config) (*Presenter, error) {
	p := &Presenter{
		controller: router,
		logRepo:    logRepo,
		config:     config,
		tokens:     map[string]bool{},
	}

	if len(config.Tokens) == 0 {
		return nil, nerr.New("no access tokens")
	}

	for _, v := range config.Tokens {
		p.tokens[v] = true
	}

	// устанавливаем middleware для проверки валидности сессии
	router.AddMiddleware("/api", p.authenticateUser)

	// добавить запись в лог
	router.AddRoute("/api", "/add", p.add(), "POST")
	// поиск записей в логе
	router.AddRoute("/api", "/search", p.search(), "POST")

	return p, nil
}

// Аутентификация пользователя
func (p *Presenter) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Authorization")
		if _, ok := p.tokens[token]; !ok {
			p.controller.RespondError(w, http.StatusUnauthorized, nerr.New(eno.ErrNoAccess))
			return
		}

		next.ServeHTTP(w, r)
	})
}

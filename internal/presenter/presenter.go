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

	tokens      map[string]bool // список всех токенов
	tokensRead  map[string]bool // список токенов доступа на чтение
	tokensWrite map[string]bool // список токенов доступа на запись
}

// New Инициализация маршрутов
func New(router httprouter.Router, logRepo LogInterface, config *config.Config) (*Presenter, error) {
	p := &Presenter{
		controller:  router,
		logRepo:     logRepo,
		config:      config,
		tokens:      map[string]bool{},
		tokensRead:  map[string]bool{},
		tokensWrite: map[string]bool{},
	}

	if len(config.TokensRead) == 0 {
		return nil, nerr.New("no access read tokens")
	}
	if len(config.TokensWrite) == 0 {
		return nil, nerr.New("no access write tokens")
	}

	// инициализация хранилища токенов
	for _, v := range config.TokensRead {
		p.tokensRead[v] = true
		p.tokens[v] = true
	}
	for _, v := range config.TokensWrite {
		p.tokensWrite[v] = true
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

// Проверка прав
func (p *Presenter) checkRights(r *http.Request, writeAccess bool) error {
	token := r.Header.Get("X-Authorization")
	if writeAccess {
		if _, ok := p.tokensWrite[token]; !ok {
			return nerr.New(eno.ErrNoAccess, "no write access")
		}
	} else {
		if _, ok := p.tokensRead[token]; !ok {
			return nerr.New(eno.ErrNoAccess, "no read access")
		}
	}
	return nil
}

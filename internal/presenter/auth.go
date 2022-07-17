package presenter

import (
	"github.com/n-r-w/eno"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/nerr"
)

type Service struct {
	cfg         *config.Config
	tokensRead  map[string]bool // список токенов доступа на чтение
	tokensWrite map[string]bool // список токенов доступа на запись
}

func New(cfg *config.Config) *Service {
	s := &Service{
		cfg:         cfg,
		tokensRead:  map[string]bool{},
		tokensWrite: map[string]bool{},
	}

	// инициализация хранилища токенов
	for _, v := range cfg.TokensRead {
		s.tokensRead[v] = true
	}
	for _, v := range cfg.TokensWrite {
		s.tokensWrite[v] = true
	}

	return s
}

// CheckRights проверка прав
func (s *Service) CheckRights(token string, writeAccess bool, readAccess bool) error {
	if _, ok := s.tokensWrite[token]; ok && writeAccess {
		return nil
	}

	if _, ok := s.tokensRead[token]; ok && readAccess {
		return nil
	}

	return nerr.New(eno.ErrNoAccess)
}

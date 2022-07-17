// Package wbuf - группировка запросов в пакеты, чтобы отправлять в БД не по одной записи
package wbuf

import (
	"runtime"

	"github.com/n-r-w/aworker"
	"github.com/n-r-w/eno"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/entity"
	"github.com/n-r-w/logsrv/internal/presenter"
	"github.com/n-r-w/nerr"
	"golang.org/x/time/rate"
)

// WBufInterface - алиас для корректной работы google wire
type WBufInterface presenter.LogInterface

type Service struct {
	log            lg.Logger
	dbRepo         WBufInterface
	aw             *aworker.AWorker
	requestLimiter *rate.Limiter
}

func New(dbRepo WBufInterface, log lg.Logger, config *config.Config) *Service {
	w := &Service{
		log:            log,
		dbRepo:         dbRepo,
		requestLimiter: rate.NewLimiter(rate.Limit(config.RateLimit), config.RateLimitBurst),
	}

	const packetSize = 100

	wcount := 0
	if config.MaxDbSessions > runtime.NumCPU() {
		wcount = runtime.NumCPU()
	} else {
		wcount = config.MaxDbSessions
	}
	w.aw = aworker.NewAWorker(packetSize*runtime.NumCPU(), packetSize, wcount, w.worker, w.processError)
	return w
}

func (w *Service) worker(messages []any) error {
	var records []entity.LogRecord

	for _, message := range messages {
		if m, ok := message.([]entity.LogRecord); !ok {
			panic("internal error")
		} else {
			records = append(records, m...)
		}
	}

	// w.log.Info("%d, %d", len(records), w.BufferSize())
	return w.dbRepo.Insert(records)
}

func (w *Service) Start() {
	w.aw.Start()
}

func (w *Service) Stop() {
	w.aw.Stop()
}

func (w *Service) BufferSize() int {
	return w.aw.QueueSize()
}

func (w *Service) processError(err error) {
	if err != nil {
		w.log.Error("write buffer error: %v", err)
	}
}

// Size - реализация интерфейса usecase.LogInterface
func (w *Service) Size() int {
	return w.BufferSize()
}

// Insert - реализация интерфейса usecase.LogInterface
func (w *Service) Insert(records []entity.LogRecord) error {
	// Защита от DDOS и в целом от перегрузки сервера БД запросами
	if !w.requestLimiter.Allow() {
		return nerr.New("wbuf: too many requests", eno.ErrTooManyRequests)
	}

	w.aw.SendMessage(records)
	return nil
}

// Find - реализация интерфейса usecase.LogInterface для его подмены
func (w *Service) Find(request entity.SearchRequest) (records []entity.LogRecord, limited bool, err error) {
	// просто пересылаем запрос
	return w.dbRepo.Find(request)
}

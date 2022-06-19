// Package wdispatch - асинхронная параллельная запись в БД
package wdispatch

import (
	"time"

	"github.com/gammazero/workerpool"
	"github.com/n-r-w/eno"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/entity"
	presentation "github.com/n-r-w/logsrv/internal/presenter"
	"github.com/n-r-w/nerr"
	"golang.org/x/time/rate"
)

type Dispatcher struct {
	log     lg.Logger
	dbRepo  presentation.LogInterface
	limiter *rate.Limiter
	pool    *workerpool.WorkerPool

	logLimiter *rate.Limiter
}

func New(dbRepo presentation.LogInterface, log lg.Logger, config *config.Config) *Dispatcher {
	d := &Dispatcher{
		log:        log,
		dbRepo:     dbRepo,
		limiter:    rate.NewLimiter(rate.Limit(config.RateLimit), config.RateLimitBurst),
		pool:       workerpool.New(config.MaxDbSessions),
		logLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}

	return d
}

// Size - реализация интерфейса usecase.LogInterface
func (d *Dispatcher) Size() int {
	return d.pool.WaitingQueueSize()
}

// Insert - реализация интерфейса usecase.LogInterface
func (d *Dispatcher) Insert(records []entity.LogRecord) error {
	// Защита от DDOS и в целом от перегрузки сервера БД запросами
	if !d.limiter.Allow() {
		return nerr.New("wdisp: too many requests", eno.ErrTooManyRequests)
	}
	// Контроль за размером очереди пула задач. Если дать ему бескотрольно расти, то можно остаться без свободных ресурсов
	// Фактически тут мы искусственно увеличиваем время отклика входящих запросов при переполнении очереди задач
	for {
		if d.pool.WaitingQueueSize() > d.pool.Size()*2 {
			if d.logLimiter.Allow() {
				d.log.Warn("wdisp: request slowing down")
			}
			time.Sleep(time.Millisecond)
		} else {
			break
		}
	}
	// Отправляем задачу на асинхронную обработку
	d.pool.Submit(func() {
		if err := d.dbRepo.Insert(records); err != nil {
			d.log.Error("worker error: %v", err)
		}
	})

	return nil
}

// Find - реализация интерфейса usecase.LogInterface для его подмены
func (d *Dispatcher) Find(request entity.SearchRequest) (records []entity.LogRecord, limited bool, err error) {
	// просто пересылаем запрос
	return d.dbRepo.Find(request)
}

func (d *Dispatcher) Start() {
}

func (d *Dispatcher) Stop() {
	d.log.Info("wdispatcher stoping...")
	d.pool.StopWait()
	d.log.Info("wdispatcher stopped OK")
}

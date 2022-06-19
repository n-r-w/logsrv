// Package wbuf - запись в БД путем накопления логов в буфер
package wbuf

import (
	"sync"
	"time"

	"github.com/n-r-w/eno"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/logsrv/internal/config"
	"github.com/n-r-w/logsrv/internal/entity"
	presentation "github.com/n-r-w/logsrv/internal/presenter"
	"github.com/n-r-w/nerr"
	"golang.org/x/time/rate"
)

// WBufInterface - алиас для корректной работы google wire
type WBufInterface presentation.LogInterface

type WBuf struct {
	log     lg.Logger
	dbRepo  WBufInterface
	records []entity.LogRecord
	mutex   sync.Mutex

	active   bool
	info     chan bool
	quit     chan struct{}
	quitWait sync.WaitGroup

	requestLimiter *rate.Limiter
	normalSize     int // размер буфера до которого будет идти накопление
	maxSize        int // критический размер буфера

	logLimiter *rate.Limiter
}

func New(dbRepo WBufInterface, log lg.Logger, config *config.Config) *WBuf {
	maxSize := 1000
	return &WBuf{
		log:            log,
		dbRepo:         dbRepo,
		records:        []entity.LogRecord{},
		mutex:          sync.Mutex{},
		active:         false,
		info:           make(chan bool, maxSize),
		quit:           make(chan struct{}),
		quitWait:       sync.WaitGroup{},
		requestLimiter: rate.NewLimiter(rate.Limit(config.RateLimit), config.RateLimitBurst),
		normalSize:     maxSize / 3,
		maxSize:        maxSize,
		logLimiter:     rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (w *WBuf) Start() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.active = true
	w.quitWait.Add(1)
	go w.worker()
}

func (w *WBuf) Stop() {
	w.log.Info("wbuffer stoping...")
	defer w.log.Info("wbuffer stopped OK")
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.active {
		return
	}

	close(w.info)
	w.quit <- struct{}{}
	w.quitWait.Wait() // если пойти дальше без этого, то можем получить креш из-за того, что воркер не успеет выйти из горутины
	close(w.quit)
}

func (w *WBuf) BufferSize() int {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return len(w.records)
}

func (w *WBuf) worker() {
	defer w.quitWait.Done()
	for {
		select {
		case ok := <-w.info: // мы приедем сюда при закрытии канала с false, поэтому надо реагировать только на true
			if ok {
				// выпиливаем записи из буфера и передаем их на обработку
				w.mutex.Lock()
				rcopy := make([]entity.LogRecord, len(w.records))
				copy(rcopy, w.records)
				w.records = nil
				w.mutex.Unlock()

				w.processError(w.processBuffer(rcopy))
			}

		case <-w.quit:
			// не лочим мьютекс, т.к. все уже залочено при отправке в канал
			if len(w.records) > 0 {
				w.processError(w.processBuffer(w.records))
				w.records = nil
			}
			return
		}
	}
}

func (w *WBuf) processBuffer(records []entity.LogRecord) error {
	if len(records) == 0 {
		return nil
	}

	return w.dbRepo.Insert(records)
}

func (w *WBuf) processError(err error) {
	if err != nil {
		w.log.Error("write buffer error: %v", err)
	}
}

// Size - реализация интерфейса usecase.LogInterface
func (w *WBuf) Size() int {
	return w.BufferSize()
}

// Insert - реализация интерфейса usecase.LogInterface
func (w *WBuf) Insert(records []entity.LogRecord) error {
	// Защита от DDOS и в целом от перегрузки сервера БД запросами
	if !w.requestLimiter.Allow() {
		return nerr.New("wbuf: too many requests", eno.ErrTooManyRequests)
	}

	if w.logLimiter.Allow() {
		size := w.BufferSize()
		if size > w.normalSize {
			w.log.Warn("buffer size: %d", size)
		}
	}

	// увеличиваем время отклика входящих запросов при переполнении буфера
	for {
		if w.BufferSize() >= w.maxSize {
			if w.logLimiter.Allow() {
				w.log.Warn("wbuf: request slowing down")
			}
			time.Sleep(time.Millisecond)
		} else {
			break
		}
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.active {
		return nerr.New(eno.ErrNotReady)
	}

	w.records = append(w.records, records...)
	if len(w.records) > w.normalSize {
		w.info <- true
	}

	return nil
}

// Find - реализация интерфейса usecase.LogInterface для его подмены
func (w *WBuf) Find(request entity.SearchRequest) (records []entity.LogRecord, limited bool, err error) {
	// просто пересылаем запрос
	return w.dbRepo.Find(request)
}

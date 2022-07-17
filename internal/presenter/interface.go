// Package presenter ...
package presenter

import (
	"github.com/n-r-w/logsrv/internal/entity"
)

// LogInterface Интерфейс работы с логами
// Реализуется в каталоге repo всеми тремя, находящимися там пакетами: psql, wbuf, wdispatch
// Контекст не прокидывается в обработчики, т.к. они выполняются асинхронно
type LogInterface interface {
	// Добавить запись в БД
	Insert(records []entity.LogRecord) error
	// Выборка записей по критериям
	Find(request entity.SearchRequest) (records []entity.LogRecord, limited bool, err error)
	// размер буфера, очереди и т.п. в зависимости от реализации интерфейса
	Size() int
}

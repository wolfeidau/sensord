package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next SensorService) SensorService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	SensorService
}

func (mw logmw) Record(id, value, stype string) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "record",
			"id", id,
			"value", value,
			"stype", stype,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.SensorService.Record(id, value, stype)
	return
}

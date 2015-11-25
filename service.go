package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/dancannon/gorethink"
)

// Reading reading used to persist in rethindb
type Reading struct {
	ID    string
	Value float64
	SType string
	TS    time.Time
}

// SensorService small sensor service for recording readings
type SensorService interface {
	Record(string, string, string) error
}

type sensorService struct {
	session *gorethink.Session
	dbName  string
}

// NewSensorService new user store
func NewSensorService(session *gorethink.Session) SensorService {
	return &sensorService{session, "sensor_data"}
}

func (s *sensorService) Record(id, value, stype string) error {

	f, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return ErrInvalid
	}

	r := &Reading{id, f, stype, time.Now()}
	if _, err = gorethink.Table("sensor_data").Insert(r).RunWrite(s.session); err != nil {
		return ErrInternal
	}
	return nil
}

// ErrInvalid is returned when a reading fails validation
var ErrInvalid = errors.New("invalid reading")

// ErrInternal is returned when an internal error occurs
var ErrInternal = errors.New("internal server error")

// ServiceMiddleware is a chainable behavior modifier for SensorService.
type ServiceMiddleware func(SensorService) SensorService

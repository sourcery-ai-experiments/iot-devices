package logging

import (
	log "github.com/kloudlite/operator/pkg/logging"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(err error, format string, args ...interface{})
}

type Options log.Options

type logger struct {
	log log.Logger
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) Errorf(err error, format string, args ...interface{}) {
	l.log.Errorf(err, format, args...)
}

func New(options *Options) (Logger, error) {
	log, err := log.New(&log.Options{
		Name:        options.Name,
		Dev:         options.Dev,
		CallerTrace: options.CallerTrace,
	})
	if err != nil {
		return nil, err
	}

	return &logger{
		log: log,
	}, nil
}

package tg

import "github.com/TheonAegor/go-framework/pkg/logger"

type Options struct {
	log logger.Loggerer
}

type Option func(*Options)

func WithLog(log logger.Loggerer) Option {
	return func(options *Options) {
		options.log = log
	}
}

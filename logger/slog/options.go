package slog

import "go.unistack.org/micro/v3/logger"

type sourceKey struct{}

func WithSourceKey(v string) logger.Option {
	return logger.SetOption(sourceKey{}, v)
}

type timeKey struct{}

func WithTimeKey(v string) logger.Option {
	return logger.SetOption(timeKey{}, v)
}

type messageKey struct{}

func WithMessageKey(v string) logger.Option {
	return logger.SetOption(messageKey{}, v)
}

type levelKey struct{}

func WithLevelKey(v string) logger.Option {
	return logger.SetOption(levelKey{}, v)
}

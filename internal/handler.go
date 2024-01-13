package internal

import "github.com/TheonAegor/gobog/pkg/tg"

type Handler struct {
	TgClient *tg.TgClient
	Cfg      Config

	Opts *Options
}

func NewHandler(tgClient *tg.TgClient, cfg Config, opts ...Option) *Handler {
	h := Handler{TgClient: tgClient, Cfg: cfg, Opts: new(Options)}

	for _, o := range opts {
		o(h.Opts)
	}

	return &h
}

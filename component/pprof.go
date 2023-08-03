package component

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/log"
)

var _ Component = &pprof{}

type pprof struct {
	Base
}

func NewPProf() *pprof {
	return &pprof{}
}

func (*pprof) Name() string {
	return "pprof"
}

func (*pprof) Start() {
	if addr := conf.GetString("app.pprof.addr"); addr != "" {
		go func() {
			log.Debug("pprof addr:", addr)
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				log.Errorf("pprof server start failed: %v", err)
			}
		}()
	}
}

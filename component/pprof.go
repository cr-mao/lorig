package component

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/utils/xcall"
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
	xcall.Call(func() {
		addr := conf.GetString("app.pprof.addr", "0.0.0.0:13000")
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Errorf("pprof server start failed: %v", err)
		}
	})
}

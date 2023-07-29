package xcall

import (
	"runtime"

	"github.com/cr-mao/lorig/log"
)

func Call(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Error(err)
			default:
				log.Errorf("panic error: %v", err)
			}
		}
	}()

	fn()
}

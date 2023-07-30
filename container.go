package lorig

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/cr-mao/lorig/component"
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/log"
	"github.com/cr-mao/lorig/utils/xfile"
)

type Container struct {
	sig        chan os.Signal
	components []component.Component
}

// NewContainer 创建一个容器
func NewContainer() *Container {
	return &Container{sig: make(chan os.Signal)}
}

// Add 添加组件
func (c *Container) Add(components ...component.Component) {
	c.components = append(c.components, components...)
}

// Serve 启动容器
func (c *Container) Serve() {
	log.Debug(fmt.Sprintf("Welcome to the lorig framework %s, Learn more at %s", Version, Website))

	for _, comp := range c.components {
		comp.Init()
	}
	for _, comp := range c.components {
		comp.Start()
	}
	c.doSavePID()
	switch runtime.GOOS {
	case `windows`:
		signal.Notify(c.sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	default:
		signal.Notify(c.sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	}
	sig := <-c.sig

	log.Warnf("process got signal %v, container will close", sig)
	signal.Stop(c.sig)
	for _, comp := range c.components {
		comp.Destroy()
	}
}

func (c *Container) doSavePID() {
	filename := conf.GetString("app.pidPath", "server.pid")
	err := xfile.WriteFile(filename, []byte(strconv.Itoa(syscall.Getpid())))
	if err != nil {
		log.Fatalf("pid save failed: %v", err)
	}
}

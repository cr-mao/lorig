package async_op

import (
	"github.com/cr-mao/lorig/log"
)

// 理解为其中的一个线程 + 队列
type worker struct {
	taskQ chan func() // LinkedBlockingQueue<Function> taskQ
}

// 处理异步过程
func (w *worker) process(asyncOp func(), continueWith func()) {
	// w *worker 这个就相当于 this
	//
	if nil == asyncOp {
		log.Error("异步操作为空")
		return
	}

	if nil == w.taskQ {
		log.Error("任务队列尚未初始化")
		return
	}

	w.taskQ <- func() { // taskQ.offer(new Function())
		// 执行异步操作
		asyncOp()

		if nil != continueWith {
			continueWith()
			//main_thread.Process(continueWith)
		}
	}
}

// 循环执行任务
func (w *worker) loopExecTask() {
	if nil == w.taskQ {
		log.Error("任务队列尚未初始化")
		return
	}

	for {
		task := <-w.taskQ

		if nil != task {
			task()
		}
	}
}

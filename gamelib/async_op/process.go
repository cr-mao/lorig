package async_op

import "sync"

var workerArray = [2048]*worker{}

// 初始化工人用的锁
var initWorkerLocker = &sync.Mutex{}

// Process 处理异步过程,
// asyncOp 异步函数, 将被放到一个新协程里去执行...
// continueWith 则是回到主线程继续执行的函数
func Process(bindId int, asyncOp func(), continueWith func()) {
	if nil == asyncOp {
		return
	}

	// 根据 bindId 获取一个工人
	currWorker := getCurrWorker(bindId)

	if nil != currWorker {
		currWorker.process(asyncOp, continueWith)
	}
}

// 根据 bindId 获取一个工人,
// bindId 会经过一次余除运算得到工人的索引
func getCurrWorker(bindId int) *worker {
	if bindId < 0 {
		bindId = -bindId
	}

	workerIndex := bindId % len(workerArray)
	currWorker := workerArray[workerIndex]

	if nil != currWorker {
		return currWorker
	}

	// 加锁
	initWorkerLocker.Lock()
	defer initWorkerLocker.Unlock()

	// 获取当前工人
	currWorker = workerArray[workerIndex]

	// 二次判断
	if nil != currWorker {
		return currWorker
	}

	// 创建一个新的工人
	currWorker = &worker{
		taskQ: make(chan func(), 2048),
	}

	// 保存到工人数组
	workerArray[workerIndex] = currWorker
	// 循环执行任务
	go currWorker.loopExecTask()

	return currWorker
}

package async_op

import (
	"fmt"
	"sync/atomic"
	"time"
)

// AsyncBizResult 异步业务结果
type AsyncBizResult struct {
	// 已返回对象
	returnedObj interface{}
	// 完成回调函数
	completeFunc func()

	// 是否有返回对象
	hasReturnedObj int32
	// 是否有完成回调函数
	hasCompleteFunc int32
	// 是否已经调用过回调函数
	completeFuncHasAlreadyBeenCalled int32
}

// GetReturnedObj 获取已返回对象
func (bizResult *AsyncBizResult) GetReturnedObj() interface{} {
	return bizResult.returnedObj
}

// SetReturnedObj 设置已返回对象
func (bizResult *AsyncBizResult) SetReturnedObj(val interface{}) {
	if atomic.CompareAndSwapInt32(&bizResult.hasReturnedObj, 0, 1) {
		fmt.Println("SetReturnedObj", time.Now().UnixMilli())
		bizResult.returnedObj = val
		bizResult.doComplete()
	}
}

// OnComplete 完成回调函数
func (bizResult *AsyncBizResult) OnComplete(val func()) {
	if atomic.CompareAndSwapInt32(&bizResult.hasCompleteFunc, 0, 1) {
		fmt.Println("OnComplete", time.Now().UnixMilli())
		bizResult.completeFunc = val

		if 1 == bizResult.hasReturnedObj {
			bizResult.doComplete()
		}
	}
}

// DoComplete 执行完成回调函数
func (bizResult *AsyncBizResult) doComplete() {
	//fmt.Printf("%T", bizResult.completeFunc)
	fmt.Println("doComplete", time.Now().UnixMilli())

	if nil == bizResult.completeFunc {
		return
	}

	fmt.Println(1111)
	if atomic.CompareAndSwapInt32(&bizResult.completeFuncHasAlreadyBeenCalled, 0, 1) {
		// 扔到主线程里去执行
		bizResult.completeFunc()
		//main_thread.Process(bizResult.completeFunc)
	}
}

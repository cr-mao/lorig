/**
*  针对 delayFunc.go 做单元测试，主要测试延迟函数结构体是否正常使用
 */
package timer

import (
	"fmt"
	"testing"
)

func SayHello(message ...interface{}) {
	fmt.Println(message[0].(string), " ", message[1].(string))
}

func TestDelayfunc(t *testing.T) {
	t.SkipNow()
	df := NewDelayFunc(SayHello, []interface{}{"hello", "zinx!"})
	fmt.Println("df.String() = ", df.String())
	df.Call()
}

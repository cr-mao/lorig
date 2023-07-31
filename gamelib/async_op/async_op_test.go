package async_op_test

import (
	"fmt"
	"github.com/cr-mao/lorig/gamelib/async_op"
	"testing"
	"time"
)

// LoginByPasswordAsync 根据用户名称和密码进行登录
func LoginByPasswordAsync(userName string, password string) *async_op.AsyncBizResult {
	if len(userName) <= 0 ||
		len(password) <= 0 {
		return nil
	}

	bizResult := &async_op.AsyncBizResult{}
	bindId := async_op.StrToWorkId(userName)
	async_op.Process(bindId, func() {
		// 模拟获得用户数据
		bizResult.SetReturnedObj("user_result")
	}, nil)

	return bizResult
}

func TestAsyncBizResult(t *testing.T) {
	// 根据用户名称和密码登录
	res := LoginByPasswordAsync("crmao", "123456")
	time.Sleep(time.Second * 3)

	res.OnComplete(func() {
		returnedObj := res.GetReturnedObj()
		if nil == returnedObj {
			fmt.Println("err nil ")
			return
		}
		user := returnedObj.(string)
		fmt.Println(111, user)
	})

	time.Sleep(time.Second * 5)
}

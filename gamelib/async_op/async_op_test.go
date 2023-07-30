package async_op

import (
	"fmt"
	"testing"
	"time"
)

// LoginByPasswordAsync 根据用户名称和密码进行登录
func LoginByPasswordAsync(userName string, password string) *AsyncBizResult {
	if len(userName) <= 0 ||
		len(password) <= 0 {
		return nil
	}

	bizResult := &AsyncBizResult{}
	bindId := StrToWorkId(userName)
	Process(bindId, func() {
		// 模拟获得用户数据
		bizResult.SetReturnedObj("user_result")
	}, nil)

	return bizResult
}

func TestAsyncBizResult(t *testing.T) {
	// 根据用户名称和密码登录
	res := LoginByPasswordAsync("crmao", "123456")
	time.Sleep(time.Second * 10)

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

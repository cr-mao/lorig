package lazy_save

import (
	"fmt"
	"testing"
	"time"
)

type User struct {
	UserId   int64 //用户id
	Blood    int64 //血量
	Capacity int64 //能量
}

type UserLazySaveObj struct {
	*User
}

func (user *UserLazySaveObj) GetLsoId() string {
	return fmt.Sprintf("userdata_%v", user.UserId)
}

func (user *UserLazySaveObj) SaveOrUpdate() {

	fmt.Println("执行存库操作，可能是异步,当前时间戳:", time.Now().Unix())
}

/*
=== RUN   TestLazySave
begin 1690676226
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=记录延迟保存数据, lsoId = userdata_1
INFO msg=执行延迟保存, lsoId = userdata_1
执行存库操作，可能是异步,当前时间戳: 1690676257
*/
func TestLazySave(t *testing.T) {
	// 模拟高频场景
	times := 0
	fmt.Println("begin", time.Now().Unix())
	for {
		time.Sleep(1 * time.Second)
		var user = &User{
			1,
			10000,
			20000,
		}
		// 延迟保存用户数据
		lazySaveObj := &UserLazySaveObj{
			User: user,
		}
		// 20秒后才会更改
		SaveOrUpdate(lazySaveObj)
		times++
		if times > 10 {
			break
		}
	}

	select {}

}

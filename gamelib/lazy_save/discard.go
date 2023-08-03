/**
User: cr-mao
Date: 2023/8/3 09:01
Email: crmao@qq.com
Desc: discard.go
*/
package lazy_save

// 取消
func Discard(lso LazySaveObj) {
	if lso == nil {
		return
	}
	lazySaveRecordMap.Delete(lso.GetLsoId())
}

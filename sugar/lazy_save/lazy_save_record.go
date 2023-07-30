package lazy_save

// 延迟保存记录
type lazySaveRecord struct {
	// 延迟保存对象
	objRef LazySaveObj
	// 最后修改时间
	lastUpdateTime int64
}

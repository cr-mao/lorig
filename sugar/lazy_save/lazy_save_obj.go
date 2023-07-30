package lazy_save

// LazySaveObj 延迟保存对象
type LazySaveObj interface {
	//
	// GetLsoId 获取 延迟保存对象 Id
	GetLsoId() string

	//
	// SaveOrUpdate 保存或更新
	SaveOrUpdate()
}

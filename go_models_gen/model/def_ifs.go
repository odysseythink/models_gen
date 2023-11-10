package model

// IModel Implement the interface to acquire database information and initialize it.实现接口获取数据库信息获取并初始化
type IModel interface {
	GenModel(username, password, host, dbname string, port uint) ([]TabInfo, error)
	GetAllTabName() map[string]string
	GetTabNamesPage(page, limit uint32, filter string) (map[string]string, uint32)
	GetTabSql(tbname string) string
	GetTabInfo(tbname string) *TabInfo
}

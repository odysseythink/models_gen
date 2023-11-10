package genmysql

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"demo.com/szlanyou_demo/models_gen/go_models_gen/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mlib.com/mlog"
)

// MySQLModel mysql model from IModel
var MySQLModel mysqlModel

type mysqlModel struct {
	db       *gorm.DB
	username string
	password string
	host     string
	dbname   string
	port     uint
}

func NewModel(username, password, host, dbname string, port uint) model.IModel {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&interpolateParams=True", username, password, host, port, dbname)
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		mlog.Errorf("gorm open dsn(%s) failed:%v", dsn, err)
		return nil
	}

	return &mysqlModel{
		db:       orm,
		username: username,
		password: password,
		host:     host,
		dbname:   dbname,
		port:     port,
	}
}

// GenModel get model.DBInfo info.获取数据库相关属性
func (m *mysqlModel) GetAllTabName() map[string]string {
	if m.db == nil || m.username == "" || m.password == "" || m.host == "" || m.dbname == "" || m.port == 0 {
		mlog.Warningf("don't have db instance")
		return nil
	}

	return m.getTables(m.db, m.dbname)
}

// getTables Get columns and comments.获取表列及注释
func (m *mysqlModel) GetTabNamesPage(page, limit uint32, filter string) (map[string]string, uint32) {
	if m.db == nil || m.username == "" || m.password == "" || m.host == "" || m.dbname == "" || m.port == 0 {
		mlog.Warningf("don't have db instance")
		return nil, 0
	}
	if limit == 0 {
		mlog.Warningf("invalid arg")
		return nil, 0
	}
	var rows *sql.Rows
	var err error
	if filter != "" {
		filter = "%" + filter + "%"
		rows, err = m.db.Raw("select COUNT(TABLE_NAME) as total from information_schema.tables where table_type='BASE TABLE' and TABLE_SCHEMA=? and TABLE_NAME LIKE ?;", m.dbname, filter).Rows()
	} else {
		rows, err = m.db.Raw("select COUNT(TABLE_NAME) as total from information_schema.tables where table_type='BASE TABLE' and TABLE_SCHEMA=?;", m.dbname).Rows()
	}

	if err != nil {
		mlog.Warningf("get total tab failed:%v", err)
		return nil, 0
	}
	var total uint32
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			mlog.Warningf("rows scan failed:%v", err)
			return nil, 0
		}
	}

	if total == 0 {
		mlog.Warningf("no tables found:%v", err)
		return nil, 0
	}

	tbDesc := make(map[string]string)
	if page < 1 {
		page = 1
	}

	if (page-1)*limit >= total {
		mlog.Warningf("beyond range")
		return nil, 0
	}
	offset := (page - 1) * limit

	if filter != "" {
		filter = "%" + filter + "%"
		rows, err = m.db.Raw("select TABLE_NAME,TABLE_COMMENT as total from information_schema.tables where table_type='BASE TABLE' and TABLE_SCHEMA=? and TABLE_NAME LIKE ? LIMIT ?,?;", m.dbname, filter, offset, limit).Rows()
	} else {
		rows, err = m.db.Raw("select TABLE_NAME,TABLE_COMMENT as total from information_schema.tables where table_type='BASE TABLE' and TABLE_SCHEMA=? LIMIT ?,?;", m.dbname, offset, limit).Rows()
	}
	if err != nil {
		mlog.Warningf("get tabes page failed:%v", err)
		return nil, 0
	}

	for rows.Next() {
		var table string
		var desc string
		rows.Scan(&table, &desc)
		tbDesc[table] = desc
	}
	rows.Close()
	return tbDesc, total
}

func (m *mysqlModel) GetTabSql(tbname string) string {
	if m.db == nil || m.username == "" || m.password == "" || m.host == "" || m.dbname == "" || m.port == 0 {
		mlog.Warningf("don't have db instance")
		return ""
	}
	if tbname == "" {
		mlog.Warning("invalid arg")
		return ""
	}
	var table, CreateTable string
	// Get create SQL statements.获取创建sql语句
	rows, err := m.db.Raw("show create table " + tbname).Rows()
	//defer rows.Close()
	if err == nil {
		if rows.Next() {
			rows.Scan(&table, &CreateTable)
		}
	}
	rows.Close()

	return CreateTable
}

func (m *mysqlModel) GetTabInfo(tbname string) *model.TabInfo {
	if m.db == nil || m.username == "" || m.password == "" || m.host == "" || m.dbname == "" || m.port == 0 {
		mlog.Warningf("don't have db instance")
		return nil
	}
	if tbname == "" {
		mlog.Warning("invalid arg")
		return nil
	}
	tab := &model.TabInfo{
		BaseInfo: model.BaseInfo{
			Name: tbname,
		},
	}
	// Get table annotations.获取表注释
	rows, err := m.db.Raw("SELECT TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema=? AND TABLE_NAME=?", m.dbname, tbname).Rows()
	if err == nil {
		if rows.Next() {
			var desc string
			rows.Scan(&desc)
			tab.BaseInfo.Notes = desc
		}
	}
	rows.Close()

	// build element.构造元素
	tab.Em = m.getTableElement(m.db, m.dbname, tbname, false)
	// --------end

	return tab
}

// GenModel get model.DBInfo info.获取数据库相关属性
func (m *mysqlModel) GenModel(username, password, host, dbname string, port uint) ([]model.TabInfo, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&interpolateParams=True", username, password, host, port, dbname)
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		mlog.Errorf("gorm open dsn(%s) failed:%v", dsn, err)
		return nil, fmt.Errorf("gorm open dsn(%s) failed:%v", dsn, err)
	}

	return m.getPackageInfo(orm, dbname, true)
}

// getTables Get columns and comments.获取表列及注释
func (m *mysqlModel) BackupTable(backfilename, tbname string, maxrecords uint) error {
	if m.db == nil || m.username == "" || m.password == "" || m.host == "" || m.dbname == "" || m.port == 0 {
		mlog.Warningf("don't have db instance")
		return fmt.Errorf("don't have db instance")
	}
	if backfilename == "" || tbname == "" {
		mlog.Warningf("invalid arg")
		return fmt.Errorf("invalid arg")
	}
	if maxrecords == 0 {
		maxrecords = 50
	}
	var rows *sql.Rows
	var err error
	rows, err = m.db.Raw(fmt.Sprintf("select COUNT(*) as total from %s;", tbname)).Rows()

	if err != nil {
		mlog.Warningf("get total table(%s) failed:%v", tbname, err)
		return fmt.Errorf("get total table(%s) failed:%v", tbname, err)
	}
	var total uint
	if rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			mlog.Warningf("rows scan failed:%v", err)
			return fmt.Errorf("rows scan failed:%v", err)
		}
	}

	if total == 0 {
		mlog.Warningf("no tables found:%v", err)
		return fmt.Errorf("no tables found:%v", err)
	}

	tbDesc := make(map[string]string)

	var page uint = 0
	var limit uint = 50
	filter := ""
	if page*limit > total {
		if total%limit > 0 {
			page = total / limit
		} else {
			page = total/limit - 1
		}
		page = total/limit - 1
	}

	if filter != "" {
		filter = "%" + filter + "%"
		rows, err = m.db.Raw("select TABLE_NAME,TABLE_COMMENT as total from information_schema.tables where table_type='BASE TABLE' and TABLE_SCHEMA=? and TABLE_NAME LIKE ? LIMIT ?,?;", m.dbname, filter, page, limit).Rows()
	} else {
		rows, err = m.db.Raw("select TABLE_NAME,TABLE_COMMENT as total from information_schema.tables where table_type='BASE TABLE' and TABLE_SCHEMA=? LIMIT ?,?;", m.dbname, page, limit).Rows()
	}
	if err != nil {
		mlog.Warningf("get tabes page failed:%v", err)
		return nil
	}

	for rows.Next() {
		var table string
		var desc string
		rows.Scan(&table, &desc)
		tbDesc[table] = desc
	}
	rows.Close()
	return nil
}

func (m *mysqlModel) getPackageInfo(orm *gorm.DB, dbname string, isOutSQL bool) ([]model.TabInfo, error) {
	tabInfoList := make([]model.TabInfo, 0)
	tabls := m.getTables(orm, dbname) // get table and notes
	// if m := config.GetTableList(); len(m) > 0 {
	// 	// 制定了表之后
	// 	newTabls := make(map[string]string)
	// 	for t := range m {
	// 		if notes, ok := tabls[t]; ok {
	// 			newTabls[t] = notes
	// 		} else {
	// 			fmt.Printf("table: %s not found in db\n", t)
	// 		}
	// 	}
	// 	tabls = newTabls
	// }
	for tabName, notes := range tabls {
		if tabName == "" {
			continue
		}
		mlog.Infof("tabName(%s)", tabName)
		var tab model.TabInfo
		tab.Name = tabName
		tab.Notes = notes

		if isOutSQL {
			// Get create SQL statements.获取创建sql语句
			rows, err := orm.Raw("show create table " + tabName).Rows()
			//defer rows.Close()
			if err == nil {
				if rows.Next() {
					var table, CreateTable string
					rows.Scan(&table, &CreateTable)
					tab.SQLBuildStr = CreateTable
				}
			}
			rows.Close()
			// ----------end
		}

		// build element.构造元素
		tab.Em = m.getTableElement(orm, dbname, tabName, false)
		// --------end

		tabInfoList = append(tabInfoList, tab)
	}
	// sort tables
	sort.Slice(tabInfoList, func(i, j int) bool {
		return tabInfoList[i].Name < tabInfoList[j].Name
	})
	return tabInfoList, nil
}

// getTableElement Get table columns and comments.获取表列及注释
func (m *mysqlModel) getTableElement(orm *gorm.DB, dbname, tab string, isForeignKey bool) (el []model.ColumnsInfo) {
	if orm == nil || dbname == "" || tab == "" {
		mlog.Warningf("invalid arg dbname(%s) tab(%s)", dbname, tab)
		return nil
	}

	keyNameCount := make(map[string]int)
	KeyColumnMp := make(map[string][]keys)
	// get keys
	var Keys []keys
	orm.Raw("show keys from " + tab).Scan(&Keys)
	for _, v := range Keys {
		keyNameCount[v.KeyName]++
		KeyColumnMp[v.ColumnName] = append(KeyColumnMp[v.ColumnName], v)
	}
	// ----------end

	var list []genColumns
	// Get table annotations.获取表注释
	orm.Raw("show FULL COLUMNS from " + tab).Scan(&list)
	// filter gorm.Model.过滤 gorm.Model
	if filterModel(&list) {
		el = append(el, model.ColumnsInfo{
			Type: "gorm.Model",
		})
	}
	// -----------------end

	// ForeignKey
	var foreignKeyList []genForeignKey
	if isForeignKey {
		sql := fmt.Sprintf(`select table_schema as table_schema,table_name as table_name,column_name as column_name,referenced_table_schema as referenced_table_schema,referenced_table_name as referenced_table_name,referenced_column_name as referenced_column_name
		from INFORMATION_SCHEMA.KEY_COLUMN_USAGE where table_schema = '%v' AND REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_NAME = '%v'`, dbname, tab)
		orm.Raw(sql).Scan(&foreignKeyList)
	}
	// ------------------end

	for _, v := range list {
		var tmp model.ColumnsInfo
		tmp.Name = v.Field
		tmp.Type = v.Type
		FixNotes(&tmp, v.Desc) // 分析表注释

		if v.Default != nil {
			if *v.Default == "" {
				tmp.Gormt = "default:''"
			} else {
				tmp.Gormt = fmt.Sprintf("default:%s", *v.Default)
			}
		}

		// keys
		if keylist, ok := KeyColumnMp[v.Field]; ok { // maybe have index or key
			for _, v := range keylist {
				if v.NonUnique == 0 { // primary or unique
					if strings.EqualFold(v.KeyName, "PRIMARY") { // PRI Set primary key.设置主键
						tmp.Index = append(tmp.Index, model.KList{
							Key:     model.ColumnsKeyPrimary,
							Multi:   (keyNameCount[v.KeyName] > 1),
							KeyType: v.IndexType,
						})
					} else { // unique
						if keyNameCount[v.KeyName] > 1 {
							tmp.Index = append(tmp.Index, model.KList{
								Key:     model.ColumnsKeyUniqueIndex,
								Multi:   (keyNameCount[v.KeyName] > 1),
								KeyName: v.KeyName,
								KeyType: v.IndexType,
							})
						} else { // unique index key.唯一复合索引
							tmp.Index = append(tmp.Index, model.KList{
								Key:     model.ColumnsKeyUnique,
								Multi:   (keyNameCount[v.KeyName] > 1),
								KeyName: v.KeyName,
								KeyType: v.IndexType,
							})
						}
					}
				} else { // mut
					tmp.Index = append(tmp.Index, model.KList{
						Key:     model.ColumnsKeyIndex,
						Multi:   true,
						KeyName: v.KeyName,
						KeyType: v.IndexType,
					})
				}
			}
		}

		tmp.IsNull = strings.EqualFold(v.Null, "YES")

		// ForeignKey
		fixForeignKey(foreignKeyList, tmp.Name, &tmp.ForeignKeyList)
		// -----------------end
		el = append(el, tmp)
	}
	return
}

// getTables Get columns and comments.获取表列及注释
func (m *mysqlModel) getTables(orm *gorm.DB, dbname string) map[string]string {
	tbDesc := make(map[string]string)

	// Get column names.获取列名
	var tables []string

	rows, err := orm.Raw("show tables").Rows()
	if err != nil {
		return tbDesc
	}

	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
		tbDesc[table] = ""
	}
	rows.Close()

	// Get table annotations.获取表注释
	rows1, err := orm.Raw("SELECT TABLE_NAME,TABLE_COMMENT FROM information_schema.TABLES WHERE table_schema= '" + dbname + "'").Rows()
	if err != nil {
		return tbDesc
	}

	for rows1.Next() {
		var table, desc string
		rows1.Scan(&table, &desc)
		tbDesc[table] = desc
	}
	rows1.Close()

	return tbDesc
}

func assemblyTable(name string) string {
	return "`" + name + "`"
}

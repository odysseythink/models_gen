package model

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// getCamelName Big Hump or Capital Letter.大驼峰或者首字母大写
func getCamelName(name string) string {

	return name
}

// TitleCase title case.首字母大写
func TitleCase(name string) string {
	vv := []rune(name)
	if len(vv) > 0 {
		if bool(vv[0] >= 'a' && vv[0] <= 'z') { // title case.首字母大写
			vv[0] -= 32
		}
	}

	return string(vv)
}

// CapLowercase 小写.且兼容 golint 驼峰命名规则
func CapLowercase(name string) string { // IDAPIID == > idAPIID
	list := strings.Split(name, "_")
	if len(list) == 0 {
		return ""
	}

	re := list[0] + name[len(list[0]):]

	return FilterKeywords(re)
}

// GetTablePrefixName 获取带表前缀名字的tablename
func GetTablePrefixName(name string) string { //
	tablePrefix := ""
	//如果设置了表前缀
	if tablePrefix != "" {
		return fmt.Sprintf("%v.%v", tablePrefix, name)
	}

	return name
}

func FilterKeywords(src string) string {
	return src
}

// getTypeName Type acquisition filtering.类型获取过滤
func getTypeName(name string, isNull bool) string {
	// 优先匹配自定义类型

	selfDefineTypeMqlDicMap := make(map[string]string)
	if v, ok := selfDefineTypeMqlDicMap[name]; ok {
		return fixNullToPorint(v, isNull)
	}

	// Precise matching first.先精确匹配
	if v, ok := TypeMysqlDicMp[name]; ok {
		return fixNullToPorint(v, isNull)
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for _, l := range TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, name); ok {
			return fixNullToPorint(l.Value, isNull)
		}
	}

	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add on (https://github.com/xxjwxc/gormt/blob/master/data/view/cnf/def.go)", name))
}

// 过滤null point 类型
func fixNullToPorint(name string, isNull bool) string {
	IsNullToPoint := false
	if isNull && IsNullToPoint {
		if strings.HasPrefix(name, "uint") {
			return "*" + name
		}
		if strings.HasPrefix(name, "int") {
			return "*" + name
		}
		if strings.HasPrefix(name, "float") {
			return "*" + name
		}
	}

	return name
}

func getUninStr(left, middle, right string) string {
	re := left
	if len(right) > 0 {
		re = left + middle + right
	}
	return re
}

func getGormModelElement() []EmInfo {
	var result []EmInfo
	result = append(result, EmInfo{
		IsMulti:       false,
		Notes:         "Primary key",
		Type:          "int64", // Type.类型标记
		ColName:       "id",
		ColStructName: "ID",
	})
	result = append(result, EmInfo{
		IsMulti:       false,
		Notes:         "created time",
		Type:          "*time.Time", // Type.类型标记
		ColName:       "created_at",
		ColStructName: "CreatedAt",
	})

	result = append(result, EmInfo{
		IsMulti:       false,
		Notes:         "updated at",
		Type:          "*time.Time", // Type.类型标记
		ColName:       "updated_at",
		ColStructName: "UpdatedAt",
	})

	result = append(result, EmInfo{
		IsMulti:       false,
		Notes:         "deleted time",
		Type:          "*time.Time", // Type.类型标记
		ColName:       "deleted_at",
		ColStructName: "DeletedAt",
	})
	return result
}

func buildFList(list *[]FList, key ColumnsKey, keyName, tp, colName string) {
	for i := 0; i < len(*list); i++ {
		if (*list)[i].KeyName == keyName {
			(*list)[i].Kem = append((*list)[i].Kem, FEm{
				Type:          tp,
				ColName:       colName,
				ColStructName: getCamelName(colName),
			})
			return
		}
	}
	// 没有 添加一个
	flist := FList{
		Key:     key,
		KeyName: keyName,
		Kem:     make([]FEm, 0),
	}
	flist.Kem = append(flist.Kem, FEm{
		Type:          tp,
		ColName:       colName,
		ColStructName: getCamelName(colName),
	})
	*list = append(*list, flist)

}

// GenPreloadList 生成list
func GenPreloadList(list []PreloadInfo, multi bool) string {
	if len(list) > 0 {
		tmpl, err := template.New("gen_preload").Parse(GetGenPreloadTpl(multi))
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		tmpl.Execute(&buf, list)

		return buf.String()
	}

	return ""
}

// GenFListIndex 生成list status(1:获取函数名,2:获取参数列表,3:获取sql case,4:值列表)
func GenFListIndex(info FList, status int) string {
	switch status {
	case 1: // 1:获取函数名
		{
			return widthFunctionName(info)
		}
	case 2: // 2:获取参数列表
		{
			var strs []string
			for _, v := range info.Kem {
				strs = append(strs, fmt.Sprintf("%v %v ", CapLowercase(v.ColStructName), v.Type))
			}
			return strings.Join(strs, ",")
		}
	case 3: // 3:获取sql case,
		{
			var strs []string
			for _, v := range info.Kem {
				strs = append(strs, fmt.Sprintf("%v = ?", v.ColName)) //wxw 2021.2.26 16:53
			}
			return strings.Join(strs, " AND ")
		}
	case 4: // 4:值列表
		{
			var strs []string
			for _, v := range info.Kem {
				strs = append(strs, CapLowercase(v.ColStructName))
			}
			return strings.Join(strs, " , ")
		}
	}

	return ""
}

func widthFunctionName(info FList) string {
	switch info.Key {
	// case ColumnsKeyDefault:
	case ColumnsKeyPrimary: // primary key.主键
		return "FetchByPrimaryKey"
	case ColumnsKeyUnique: // unique key.唯一索引
		return "FetchUniqueBy" + getCamelName(info.KeyName)
	case ColumnsKeyIndex: // index key.复合索引
		return "FetchIndexBy" + getCamelName(info.KeyName)
	case ColumnsKeyUniqueIndex: // unique index key.唯一复合索引
		return "FetchUniqueIndexBy" + getCamelName(info.KeyName)
	}

	return ""
}

func fixNotes(str string) string { // 注释
	return strings.Replace(str, "\n", "\n//", -1)
}

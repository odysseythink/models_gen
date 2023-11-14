package model

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"demo.com/szlanyou_demo/models_gen/go_models_gen/hump"
	"mlib.com/mlog"
)

type _Model struct {
	info *DBInfo
	pkg  *GenPackage
}

// Generate build code string.生成代码
func GenerateTabModelCode(m IModel, tbname, tbprefix string) string {
	if m == nil || tbname == "" {
		mlog.Warning("invalid arg")
		return ""
	}
	info := m.GetTabInfo(tbname)
	if info == nil {
		mlog.Warningf("can't get tabinfo by tbname(%s)", tbname)
		return ""
	}
	var sct GenStruct

	sct.SetTableName(info.Name)

	structName := ""
	if tbprefix != "" && strings.HasPrefix(tbname, tbprefix) {
		structName = tbname[len(tbprefix):]
	} else {
		structName = tbname
	}
	structName = hump.BigHumpName(structName)
	//如果设置了表前缀
	// if tablePrefix != "" {
	// 	tab.Name = strings.TrimLeft(tab.Name, tablePrefix)
	// }
	var pa PrintAtom
	pa.Add("package", "models")
	sct.SetStructName(structName) // Big hump.大驼峰
	sct.SetNotes(info.BaseInfo.Notes)
	sct.AddElement(getGenElement(info.Em)...) // build element.构造元素
	sct.SetCreatTableStr(info.SQLBuildStr)

	for _, v1 := range sct.Generates() {
		pa.Add(v1)
	}

	// add table name func
	for _, v1 := range sct.GenerateTableName() {
		pa.Add(v1)
	}

	output := ""
	for _, v := range pa.Generates() {
		output += v
		output += "\n"
	}
	return output
}

// Generate build code string.生成代码
func Generate(info *DBInfo) (out []GenOutInfo, m _Model) {
	m = _Model{
		info: info,
	}

	// struct
	var stt GenOutInfo
	stt.FileCtx = m.generate()
	stt.FileName = info.DbName + ".go"

	OutFileName := ""
	if name := OutFileName; len(name) > 0 {
		stt.FileName = name + ".go"
	}

	out = append(out, stt)
	// ------end

	// gen function
	IsOutFunc := false
	if IsOutFunc {
		out = append(out, m.generateFunc()...)
	}
	// -------------- end
	return
}

// genTableElement Get table columns and comments.获取表列及注释
func getGenElement(cols []ColumnsInfo) (el []GenElement) {
	_tagGorm := ""
	if _tagGorm == "" {
		_tagGorm = "gorm"
	}
	_tagJSON := ""
	if _tagJSON == "" {
		_tagJSON = "json"
	}

	for _, v := range cols {
		var tmp GenElement
		var isPK bool
		if strings.EqualFold(v.Type, "gorm.Model") { // gorm model
			tmp.SetType(v.Type) //
		} else {
			tmp.SetName(getCamelName(v.Name))
			tmp.SetNotes(v.Notes)
			tmp.SetType(getTypeName(v.Type, v.IsNull))
			for _, v1 := range v.Index {
				switch v1.Key {
				// case ColumnsKeyDefault:
				case ColumnsKeyPrimary: // primary key.主键
					tmp.AddTag(_tagGorm, "primaryKey")
					isPK = true
				case ColumnsKeyUnique: // unique key.唯一索引
					tmp.AddTag(_tagGorm, "unique")
				case ColumnsKeyIndex: // index key.复合索引
					if v1.KeyType == "FULLTEXT" {
						tmp.AddTag(_tagGorm, getUninStr("index", ":", v1.KeyName)+",class:FULLTEXT")
					} else {
						tmp.AddTag(_tagGorm, getUninStr("index", ":", v1.KeyName))
					}
				case ColumnsKeyUniqueIndex: // unique index key.唯一复合索引
					tmp.AddTag(_tagGorm, getUninStr("uniqueIndex", ":", v1.KeyName))
				}
			}
		}

		if len(v.Name) > 0 {
			// not simple output
			Simple := false
			if !Simple {
				tmp.AddTag(_tagGorm, "column:"+v.Name)
				tmp.AddTag(_tagGorm, "type:"+v.Type)
				if !v.IsNull {
					tmp.AddTag(_tagGorm, "not null")
				}
			}
			// default tag
			if len(v.Gormt) > 0 {
				tmp.AddTag(_tagGorm, v.Gormt)
			}

			// json tag
			IsWebTagPkHidden := false
			if isPK && IsWebTagPkHidden {
				tmp.AddTag(_tagJSON, "-")
			} else {
				tmp.AddTag(_tagJSON, v.Name)
			}

		}

		tmp.ColumnName = v.Name // 列名
		el = append(el, tmp)
	}

	return
}

// GetPackage gen struct on table
func (m *_Model) GetPackage() GenPackage {
	if m.pkg == nil {
		var pkg GenPackage
		pkg.SetPackage(m.info.PackageName) //package name

		// tablePrefix := config.GetTablePrefix()

		for _, tab := range m.info.TabList {
			var sct GenStruct

			sct.SetTableName(tab.Name)

			//如果设置了表前缀
			// if tablePrefix != "" {
			// 	tab.Name = strings.TrimLeft(tab.Name, tablePrefix)
			// }

			sct.SetStructName(getCamelName(tab.Name)) // Big hump.大驼峰
			sct.SetNotes(tab.Notes)
			sct.AddElement(m.genTableElement(tab.Em)...) // build element.构造元素
			sct.SetCreatTableStr(tab.SQLBuildStr)
			pkg.AddStruct(sct)
		}
		m.pkg = &pkg
	}

	return *m.pkg
}

func (m *_Model) generate() string {
	m.pkg = nil
	m.GetPackage()
	return m.pkg.Generate()
}

// genTableElement Get table columns and comments.获取表列及注释
func (m *_Model) genTableElement(cols []ColumnsInfo) (el []GenElement) {
	_tagGorm := ""
	_tagJSON := ""

	for _, v := range cols {
		var tmp GenElement
		var isPK bool
		if strings.EqualFold(v.Type, "gorm.Model") { // gorm model
			tmp.SetType(v.Type) //
		} else {
			tmp.SetName(getCamelName(v.Name))
			tmp.SetNotes(v.Notes)
			tmp.SetType(getTypeName(v.Type, v.IsNull))
			for _, v1 := range v.Index {
				switch v1.Key {
				// case ColumnsKeyDefault:
				case ColumnsKeyPrimary: // primary key.主键
					tmp.AddTag(_tagGorm, "primaryKey")
					isPK = true
				case ColumnsKeyUnique: // unique key.唯一索引
					tmp.AddTag(_tagGorm, "unique")
				case ColumnsKeyIndex: // index key.复合索引
					if v1.KeyType == "FULLTEXT" {
						tmp.AddTag(_tagGorm, getUninStr("index", ":", v1.KeyName)+",class:FULLTEXT")
					} else {
						tmp.AddTag(_tagGorm, getUninStr("index", ":", v1.KeyName))
					}
				case ColumnsKeyUniqueIndex: // unique index key.唯一复合索引
					tmp.AddTag(_tagGorm, getUninStr("uniqueIndex", ":", v1.KeyName))
				}
			}
		}

		if len(v.Name) > 0 {
			// not simple output
			Simple := false
			if !Simple {
				tmp.AddTag(_tagGorm, "column:"+v.Name)
				tmp.AddTag(_tagGorm, "type:"+v.Type)
				if !v.IsNull {
					tmp.AddTag(_tagGorm, "not null")
				}
			}
			// default tag
			if len(v.Gormt) > 0 {
				tmp.AddTag(_tagGorm, v.Gormt)
			}

			// json tag
			IsWEBTag := false
			if IsWEBTag {
				IsWebTagPkHidden := false
				if isPK && IsWebTagPkHidden {
					tmp.AddTag(_tagJSON, "-")
				} else {
					WebTagType := 0
					if WebTagType == 0 {
						tmp.AddTag(_tagJSON, v.Name)
					} else {
						tmp.AddTag(_tagJSON, v.Name)
					}
				}
			}

		}

		tmp.ColumnName = v.Name // 列名
		el = append(el, tmp)

		// ForeignKey
		IsForeignKey := false
		if IsForeignKey && len(v.ForeignKeyList) > 0 {
			fklist := m.genForeignKey(v)
			el = append(el, fklist...)
		}
		// -----------end
	}

	return
}

// genForeignKey Get information about foreign key of table column.获取表列外键相关信息
func (m *_Model) genForeignKey(col ColumnsInfo) (fklist []GenElement) {
	_tagGorm := ""
	_tagJSON := ""

	for _, v := range col.ForeignKeyList {
		isMulti, isFind, notes := m.getColumnsKeyMulti(v.TableName, v.ColumnName)
		if isFind {
			var tmp GenElement
			tmp.SetNotes(notes)
			if isMulti {
				tmp.SetName(getCamelName(v.TableName) + "List")
				tmp.SetType("[]" + getCamelName(v.TableName))
			} else {
				tmp.SetName(getCamelName(v.TableName))
				tmp.SetType(getCamelName(v.TableName))
			}

			tmp.AddTag(_tagGorm, "joinForeignKey:"+col.Name) // association_foreignkey
			tmp.AddTag(_tagGorm, "foreignKey:"+v.ColumnName)

			// json tag
			IsWEBTag := false
			if IsWEBTag {
				WebTagType := 0
				if WebTagType == 0 {
					tmp.AddTag(_tagJSON, v.TableName+"List")
				} else {
					tmp.AddTag(_tagJSON, v.TableName+"_list")
				}
			}

			fklist = append(fklist, tmp)
		}
	}

	return
}

func (m *_Model) getColumnsKeyMulti(tableName, col string) (isMulti bool, isFind bool, notes string) {
	var haveGomod bool
	for _, v := range m.info.TabList {
		if strings.EqualFold(v.Name, tableName) {
			for _, v1 := range v.Em {
				if strings.EqualFold(v1.Name, col) {
					for _, v2 := range v1.Index {
						switch v2.Key {
						case ColumnsKeyPrimary, ColumnsKeyUnique, ColumnsKeyUniqueIndex: // primary key unique key . 主键，唯一索引
							{
								if !v2.Multi { // 唯一索引
									return false, true, v.Notes
								}
							}
							// case ColumnsKeyIndex: // index key. 复合索引
							// 	{
							// 		isMulti = true
							// 	}
						}
					}
					return true, true, v.Notes
				} else if strings.EqualFold(v1.Type, "gorm.Model") {
					haveGomod = true
					notes = v.Notes
				}
			}
			break
		}
	}

	// default gorm.Model
	if haveGomod {
		if strings.EqualFold(col, "id") {
			return false, true, notes
		}

		if strings.EqualFold(col, "created_at") ||
			strings.EqualFold(col, "updated_at") ||
			strings.EqualFold(col, "deleted_at") {
			return true, true, notes
		}
	}

	return false, false, ""
	// -----------------end
}

// ///////////////////////// func
func (m *_Model) generateFunc() (genOut []GenOutInfo) {
	// getn base
	tmpl, err := template.New("gen_base").Parse(GetGenBaseTpl())
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, m.info)
	genOut = append(genOut, GenOutInfo{
		FileName: "gen.base.go",
		FileCtx:  buf.String(),
	})
	//tools.WriteFile(outDir+"gen_router.go", []string{buf.String()}, true)
	// -------end------

	for _, tab := range m.info.TabList {
		var pkg GenPackage
		pkg.SetPackage(m.info.PackageName) //package name
		pkg.AddImport(`"fmt"`)
		pkg.AddImport(`"context"`) // 添加import信息
		pkg.AddImport(EImportsHead["gorm.Model"])

		// wxw 2021.2.26 17:17
		var data funDef
		data.TableName = tab.Name
		// tablePrefix := config.GetTablePrefix()
		// //如果设置了表前缀
		// if tablePrefix != "" {
		// 	tab.Name = strings.TrimLeft(tab.Name, tablePrefix)
		// }
		data.StructName = getCamelName(tab.Name)

		var primary, unique, uniqueIndex, index []FList
		for _, el := range tab.Em {
			if strings.EqualFold(el.Type, "gorm.Model") {
				data.Em = append(data.Em, getGormModelElement()...)
				pkg.AddImport(`"time"`)
				buildFList(&primary, ColumnsKeyPrimary, "", "int64", "id")
			} else {
				typeName := getTypeName(el.Type, el.IsNull)
				isMulti := (len(el.Index) == 0)
				isUniquePrimary := false
				for _, v1 := range el.Index {
					if v1.Multi {
						isMulti = v1.Multi
					}

					switch v1.Key {
					// case ColumnsKeyDefault:
					case ColumnsKeyPrimary: // primary key.主键
						isUniquePrimary = !v1.Multi
						buildFList(&primary, ColumnsKeyPrimary, v1.KeyName, typeName, el.Name)
					case ColumnsKeyUnique: // unique key.唯一索引
						buildFList(&unique, ColumnsKeyUnique, v1.KeyName, typeName, el.Name)
					case ColumnsKeyIndex: // index key.复合索引
						buildFList(&index, ColumnsKeyIndex, v1.KeyName, typeName, el.Name)
					case ColumnsKeyUniqueIndex: // unique index key.唯一复合索引
						buildFList(&uniqueIndex, ColumnsKeyUniqueIndex, v1.KeyName, typeName, el.Name)
					}
				}

				if isMulti && isUniquePrimary { // 主键唯一
					isMulti = false
				}

				data.Em = append(data.Em, EmInfo{
					IsMulti:       isMulti,
					Notes:         fixNotes(el.Notes),
					Type:          typeName, // Type.类型标记
					ColName:       el.Name,
					ColNameEx:     fmt.Sprintf("`%v`", el.Name),
					ColStructName: getCamelName(el.Name),
				})
				if v2, ok := EImportsHead[typeName]; ok {
					if len(v2) > 0 {
						pkg.AddImport(v2)
					}
				}
			}

			// 外键列表
			for _, v := range el.ForeignKeyList {
				isMulti, isFind, notes := m.getColumnsKeyMulti(v.TableName, v.ColumnName)
				if isFind {
					var info PreloadInfo
					info.IsMulti = isMulti
					info.Notes = fixNotes(notes)
					info.ForeignkeyTableName = v.TableName
					info.ForeignkeyCol = v.ColumnName
					info.ForeignkeyStructName = getCamelName(v.TableName)
					info.ColName = el.Name
					info.ColStructName = getCamelName(el.Name)
					data.PreloadList = append(data.PreloadList, info)
				}
			}
			// ---------end--
		}

		data.Primay = append(data.Primay, primary...)
		data.Primay = append(data.Primay, unique...)
		data.Primay = append(data.Primay, uniqueIndex...)
		data.Index = append(data.Index, index...)
		tmpl, err := template.New("gen_logic").
			Funcs(template.FuncMap{"GenPreloadList": GenPreloadList, "GenFListIndex": GenFListIndex, "CapLowercase": CapLowercase, "GetTablePrefixName": GetTablePrefixName}).
			Parse(GetGenLogicTpl())
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		tmpl.Execute(&buf, data)

		pkg.AddFuncStr(buf.String())
		genOut = append(genOut, GenOutInfo{
			FileName: fmt.Sprintf(m.info.DbName+".gen.%v.go", tab.Name),
			FileCtx:  pkg.Generate(),
		})
	}

	return
}

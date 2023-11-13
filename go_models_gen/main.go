package main

/*
#include "c_2_go_interface.h"
void C_Callback(int cmd,void* p,int l);;
*/
import "C"
import (
	"fmt"
	"unsafe"

	"demo.com/szlanyou_demo/models_gen/go_models_gen/model"
	"demo.com/szlanyou_demo/models_gen/go_models_gen/model/genmysql"
	"demo.com/szlanyou_demo/models_gen/go_models_gen/pbapi"
	"github.com/golang/protobuf/proto"
	"mlib.com/mlog"
)

// go build -buildmode=c-archive -o libmodels_gen.a
func parseReqFromCPointer(req unsafe.Pointer, reqlen C.int, pbmsg proto.Message) error {
	reqarr := unsafe.Slice((*C.char)(req), reqlen)
	reqbuf := make([]byte, reqlen)
	for iLoop := 0; iLoop < int(reqlen); iLoop++ {
		reqbuf[iLoop] = byte(reqarr[iLoop])
	}
	err := proto.Unmarshal(reqbuf, pbmsg)
	if err != nil {
		mlog.Errorf("proto.Marshal failed")
		return err
	}
	return nil
}

func packRspToCPointer(pbmsg proto.Message, l *C.int) (rsp unsafe.Pointer) {
	data, err := proto.Marshal(pbmsg)
	if err != nil {
		mlog.Errorf("proto.Marshal failed")
		*l = C.int(0)
		return nil
	}

	p := C.malloc(C.size_t(len(data)))
	arr := unsafe.Slice((*C.char)(p), len(data))
	for k, v := range data {
		arr[k] = C.char(v)
	}
	*l = C.int(len(data))
	return p
}

func callNoticeCallback(cmd uint32, pbmsg proto.Message) {
	var l C.int = 0
	data, err := proto.Marshal(pbmsg)
	if err != nil {
		mlog.Errorf("proto.Marshal failed")
		return
	}

	p := C.malloc(C.size_t(len(data)))
	arr := unsafe.Slice((*C.char)(p), len(data))
	for k, v := range data {
		arr[k] = C.char(v)
	}
	l = C.int(len(data))
	C.C_Callback(C.int(cmd), p, l)
}

//export OpenDb
func OpenDb(req unsafe.Pointer, reqlen C.int, l *C.int) unsafe.Pointer {
	myreq := pbapi.PK_OPEN_DB_REQ{}
	myrsp := pbapi.PK_OPEN_DB_RSP{}
	err := parseReqFromCPointer(req, reqlen, &myreq)
	if err != nil {
		mlog.Errorf("parse req failed:%v", err)
		myrsp.Errmsg = fmt.Sprintf("parse req failed:%v", err)
		return packRspToCPointer(&myrsp, l)
	}
	mlog.Infof("OpenDb call %#v", myreq)

	if myreq.Host == "" || myreq.Port == 0 || myreq.Username == "" || myreq.Passwd == "" || myreq.Dbname == "" {
		mlog.Warningf("invalid arg")
		myrsp.Errmsg = "invalid arg"
		return packRspToCPointer(&myrsp, l)
	}
	mDbModel = genmysql.NewModel(myreq.Username, myreq.Passwd, myreq.Host, myreq.Dbname, uint(myreq.Port))
	if mDbModel == nil {
		mlog.Warningf("create mysql model failed")
		myrsp.Errmsg = "create mysql model failed"
		return packRspToCPointer(&myrsp, l)
	}

	return packRspToCPointer(&myrsp, l)
}

//export GetAllTabNames
func GetAllTabNames(req unsafe.Pointer, reqlen C.int, l *C.int) unsafe.Pointer {
	myreq := pbapi.PK_GET_ALL_TABNAMES_REQ{}
	myrsp := pbapi.PK_GET_ALL_TABNAMES_RSP{}
	err := parseReqFromCPointer(req, reqlen, &myreq)
	if err != nil {
		mlog.Errorf("parse req failed:%v", err)
		myrsp.Errmsg = fmt.Sprintf("parse req failed:%v", err)
		return packRspToCPointer(&myrsp, l)
	}
	mlog.Infof("GetAllTabNames call %#v", myreq)

	if mDbModel == nil {
		mlog.Errorf("no db model instance")
		myrsp.Errmsg = "no db model instance"
		return packRspToCPointer(&myrsp, l)
	} else {
		tbnames := mDbModel.GetAllTabName()
		if len(tbnames) > 0 {
			myrsp.Names = make([]string, 0)
			myrsp.Deses = make([]string, 0)
			for k, v := range tbnames {
				myrsp.Names = append(myrsp.Names, k)
				myrsp.Deses = append(myrsp.Deses, v)
			}
		}
		return packRspToCPointer(&myrsp, l)
	}
}

//export GetTabNamesPage
func GetTabNamesPage(req unsafe.Pointer, reqlen C.int, l *C.int) unsafe.Pointer {
	myreq := pbapi.PK_GET_TABNAMES_PAGE_REQ{}
	myrsp := pbapi.PK_GET_TABNAMES_PAGE_RSP{}
	err := parseReqFromCPointer(req, reqlen, &myreq)
	if err != nil {
		mlog.Errorf("parse req failed:%v", err)
		myrsp.Errmsg = fmt.Sprintf("parse req failed:%v", err)
		return packRspToCPointer(&myrsp, l)
	}
	mlog.Infof("GetTabNamesPage call %#v", myreq)

	if myreq.Limit == 0 {
		mlog.Error("invalid arg")
		myrsp.Errmsg = "invalid arg"
		return packRspToCPointer(&myrsp, l)
	}
	if mDbModel == nil {
		mlog.Errorf("no db model instance")
		myrsp.Errmsg = "no db model instance"
		return packRspToCPointer(&myrsp, l)
	} else {
		tbnames, total := mDbModel.GetTabNamesPage(myreq.Page, myreq.Limit, myreq.Filter)
		if len(tbnames) > 0 {
			myrsp.Names = make([]string, 0)
			myrsp.Deses = make([]string, 0)
			for k, v := range tbnames {
				myrsp.Names = append(myrsp.Names, k)
				myrsp.Deses = append(myrsp.Deses, v)
			}
		}
		myrsp.Total = total
		return packRspToCPointer(&myrsp, l)
	}
}

//export GetTabSql
func GetTabSql(req unsafe.Pointer, reqlen C.int, l *C.int) unsafe.Pointer {
	myreq := pbapi.PK_GET_TAB_SQL_REQ{}
	myrsp := pbapi.PK_GET_TAB_SQL_RSP{}
	err := parseReqFromCPointer(req, reqlen, &myreq)
	if err != nil {
		mlog.Errorf("parse req failed:%v", err)
		myrsp.Errmsg = fmt.Sprintf("parse req failed:%v", err)
		return packRspToCPointer(&myrsp, l)
	}
	mlog.Infof("GetTabSql call %#v", myreq)

	if myreq.Tabname == "" {
		mlog.Error("invalid arg")
		myrsp.Errmsg = "invalid arg"
		return packRspToCPointer(&myrsp, l)
	}
	if mDbModel == nil {
		mlog.Errorf("no db model instance")
		myrsp.Errmsg = "no db model instance"
		return packRspToCPointer(&myrsp, l)
	} else {
		myrsp.Sql = mDbModel.GetTabSql(myreq.Tabname)
		return packRspToCPointer(&myrsp, l)
	}
}

//export GetTabModelCode
func GetTabModelCode(req unsafe.Pointer, reqlen C.int, l *C.int) unsafe.Pointer {
	myreq := pbapi.PK_GET_TAB_MODEL_CODE_REQ{}
	myrsp := pbapi.PK_GET_TAB_MODEL_CODE_RSP{}
	err := parseReqFromCPointer(req, reqlen, &myreq)
	if err != nil {
		mlog.Errorf("parse req failed:%v", err)
		myrsp.Errmsg = fmt.Sprintf("parse req failed:%v", err)
		return packRspToCPointer(&myrsp, l)
	}
	mlog.Infof("GetTabModelCode call %#v", myreq)

	if myreq.Tabname == "" {
		mlog.Error("invalid arg")
		myrsp.Errmsg = "invalid arg"
		return packRspToCPointer(&myrsp, l)
	}
	if mDbModel == nil {
		mlog.Errorf("no db model instance")
		myrsp.Errmsg = "no db model instance"
		return packRspToCPointer(&myrsp, l)
	} else {
		myrsp.Code = model.GenerateTabModelCode(mDbModel, myreq.Tabname, myreq.Prefix)
		return packRspToCPointer(&myrsp, l)
	}
}

//export SetLogDir
func SetLogDir(req unsafe.Pointer, reqlen C.int, l *C.int) unsafe.Pointer {
	myreq := pbapi.PK_SET_LOG_DIR_REQ{}
	myrsp := pbapi.PK_SET_LOG_DIR_RSP{}
	err := parseReqFromCPointer(req, reqlen, &myreq)
	if err != nil {
		mlog.Errorf("parse req failed:%v", err)
		myrsp.Errmsg = fmt.Sprintf("parse req failed:%v", err)
		return packRspToCPointer(&myrsp, l)
	}
	if myreq.Dir == "" {
		mlog.Flush()
	} else {
		mlog.SetLogDir(myreq.Dir)
		mlog.Info("log dir=", myreq.Dir)
	}
	return packRspToCPointer(&myrsp, l)
}

var mDbModel model.IModel = nil

func main() {}

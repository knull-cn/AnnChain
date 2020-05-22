package server

import (
	"context"
	"path"

	"github.com/dappledger/AnnChain/bcstore/proto"
	"github.com/dappledger/AnnChain/bcstore/types"
	log "github.com/sirupsen/logrus"
)

type StoreService struct {
	dbtype  types.StoreType
	datadir string
	servers map[string]Server
}

func NewStoreService(datadir string, dbtype types.StoreType) *StoreService {
	log.Tracef("NewStoreService(%s;%v)",datadir,dbtype)
	if datadir == "" {
		datadir = "~/.data"
	}
	switch dbtype {
	case types.ST_GOLevelDB, types.ST_MemDB, types.ST_CLevelDB:
	default:
		dbtype = types.ST_GOLevelDB
	}
	return &StoreService{
		dbtype,
		datadir,
		make(map[string]Server, 8),
	}
}

func proto2types(kv *proto.KeyValue) types.KeyValue {
	return types.KeyValue{
		Key:   kv.GetKey(),
		Value: kv.GetValue(),
	}
}

func types2proto(from types.KeyValue, to *proto.KeyValue) {
	to.Key = from.Key
	to.Value = from.Value
}

func (ss *StoreService) GetServer(name string) (svr Server, err error) {
	var ok bool
	svr, ok = ss.servers[name]
	if !ok {
		svr, err = NewServer(path.Join(ss.datadir, name), ss.dbtype)
		if err == nil {
			ss.servers[name] = svr
		}
		log.Tracef("GetServer(%s;%s) new err=%v.",ss.datadir,name,err)
	}
	return
}

func (ss *StoreService) Get(ctx context.Context, req *proto.ReqGet) (*proto.RespGet, error) {
	log.Tracef("Get.req=%x(%s).",req.GetKey(),string(req.GetKey()))
	var resp = proto.RespGet{
		Base: &proto.RespBase{},
	}
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.Get : GetServer error:%s", err.Error())
		resp.Base.ErrorCode = int64(types.ErrDatabase)
		resp.Base.Error = "ss.Get : GetServer error:" + err.Error()
		return &resp, nil
	}
	kv, err := is.Get(ctx, req.GetKey())
	if err != nil {
		log.Errorf("ss.Get (%x) error:%s.", req.GetKey(), err.Error())
		resp.Base.ErrorCode = int64(types.ErrDatabase)
		resp.Base.Error = "ss.Get error:" + err.Error()
		return &resp, nil
	}
	resp.Kv = &proto.KeyValue{
		Key:   kv.Key,
		Value: kv.Value,
	}
	return &resp, nil
}

func (ss *StoreService) Set(ctx context.Context, req *proto.ReqSet) (*proto.RespBase, error) {
	log.Tracef("Set.req=%x(%s).",req.Kv.GetKey(),string(req.Kv.GetKey()))
	var resp proto.RespBase
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.Set : GetServer error:%s", err.Error())
		resp.ErrorCode = int64(types.ErrDatabase)
		resp.Error = "ss.Set : GetServer error:" + err.Error()
		return &resp, nil
	}
	err = is.Set(ctx, proto2types(req.GetKv()))
	if err != nil {
		log.Errorf("ss.Set (%x) error:%s.", req.GetKv(), err.Error())
		resp.ErrorCode = int64(types.ErrDatabase)
		resp.Error = "ss.Set error:" + err.Error()
		return &resp, nil
	}
	return &resp, nil
}

func (ss *StoreService) GetByPrefix(ctx context.Context, req *proto.ReqGetByPrefix) (*proto.RespGetByPrefix, error) {
	log.Tracef("GetByPrefix.req=%x(%s).",req.GetPrefix(),string(req.GetPrefix()))
	var resp = proto.RespGetByPrefix{
		Base: &proto.RespBase{},
	}
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.GetByPrefix : GetServer error:%s", err.Error())
		resp.GetBase().ErrorCode = int64(types.ErrDatabase)
		resp.GetBase().Error = "ss.GetByPrefix : GetServer error:" + err.Error()
		return &resp, nil
	}
	kvs, err := is.GetByPrefix(ctx, req.GetPrefix(), req.GetLastKey(), req.GetLimit())
	if err != nil {
		log.Errorf("ss.GetByPrefix (%x;%x;%d) error:%s.", req.GetPrefix(), req.GetLastKey(), req.GetLimit(), err.Error())
		resp.Base.ErrorCode = int64(types.ErrDatabase)
		resp.Base.Error = "ss.GetByPrefix error:" + err.Error()
		return &resp, nil
	}
	var list = make([]*proto.KeyValue, len(kvs))
	for i := 0; i < len(kvs); i++ {
		list[i] = new(proto.KeyValue)
		types2proto(kvs[i], list[i])
	}
	resp.List = list
	return &resp, nil
}
func (ss *StoreService) Batch(ctx context.Context, req *proto.ReqBatch) (*proto.RespBase, error) {
	log.Tracef("Batch.req= %d'dels;%d's sets.",len(req.GetDels()),len(req.GetSets()))
	var resp proto.RespBase
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.Batch : GetServer error:%s", err.Error())
		resp.ErrorCode = int64(types.ErrDatabase)
		resp.Error = "ss.Batch : GetServer error:" + err.Error()
		return &resp, nil
	}
	var dels = make([]types.Key, len(req.Dels))
	var sets = make([]types.KeyValue, len(req.Sets))
	for i := 0; i < len(req.Dels); i++ {
		dels[i] = req.Dels[i]
	}
	for i := 0; i < len(req.Sets); i++ {
		sets[i].Key = req.Sets[i].GetKey()
		sets[i].Value = req.Sets[i].GetValue()
	}
	err = is.Batch(ctx, dels, sets)
	if err != nil {
		log.Errorf("ss.Batch <len(dels)=%d;len(sets)=%d) error:%s.", len(dels), len(sets), err.Error())
		resp.ErrorCode = int64(types.ErrDatabase)
		resp.Error = "ss.Batch error:" + err.Error()
		return &resp, nil
	}
	return &resp, nil
}
func (ss *StoreService) IsExist(ctx context.Context, req *proto.ReqIsExist) (*proto.RespIsExist, error) {
	log.Tracef("IsExist.req=%x(%s).",req.GetKey(),string(req.GetKey()))
	var resp = proto.RespIsExist{
		Base: &proto.RespBase{},
	}
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.IsExist : GetServer error:%s", err.Error())
		resp.GetBase().ErrorCode = int64(types.ErrDatabase)
		resp.GetBase().Error = "ss.IsExist : GetServer error:" + err.Error()
		return &resp, nil
	}
	ok, err := is.Has(ctx, req.GetKey())
	if err != nil {
		log.Errorf("ss.IsExist (%x) error:%s.", req.GetKey(), err.Error())
		resp.GetBase().ErrorCode = int64(types.ErrDatabase)
		resp.GetBase().Error = "ss.IsExist error:" + err.Error()
		return &resp, nil
	}
	resp.IsExist = ok
	return &resp, nil
}
func (ss *StoreService) GetProperty(ctx context.Context, req *proto.ReqGetProperty) (*proto.RespGetProperty, error) {
	log.Tracef("GetProperty.req=%d(%v).",req.GetProperty(),req.GetProperty())
	var resp = proto.RespGetProperty{
		Base: &proto.RespBase{},
	}
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.GetProperty : GetServer error:%s", err.Error())
		resp.GetBase().ErrorCode = int64(types.ErrDatabase)
		resp.GetBase().Error = "ss.GetProperty : GetServer error:" + err.Error()
		return &resp, nil
	}
	value, err := is.GetProperty(ctx, req.GetProperty())
	if err != nil {
		log.Errorf("ss.GetProperty (%v) error:%s.", req.GetProperty(), err.Error())
		resp.GetBase().ErrorCode = int64(types.ErrDatabase)
		resp.GetBase().Error = "ss.GetProperty error:" + err.Error()
		return &resp, nil
	}
	resp.Ret = value
	return &resp, nil
}

func (ss *StoreService) Delete(ctx context.Context, req *proto.ReqDelete) (*proto.RespBase, error) {
	log.Tracef("Delete.req=%x(%s).",req.GetKey(),string(req.GetKey()))
	var resp proto.RespBase
	is, err := ss.GetServer(req.GetDbName())
	if err != nil {
		log.Errorf("ss.Delete : GetServer error:%s", err.Error())
		resp.ErrorCode = int64(types.ErrDatabase)
		resp.Error = "ss.Delete : GetServer error:" + err.Error()
		return &resp, nil
	}
	err = is.Delete(ctx, req.GetKey())
	if err != nil {
		log.Errorf("ss.Delete(%x) error:%s.", req.GetKey(), err.Error())
		resp.ErrorCode = int64(types.ErrDatabase)
		resp.Error = "ss.Delete error:" + err.Error()
		return &resp, nil
	}
	return &resp, nil
}

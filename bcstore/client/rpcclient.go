package client

import (
    `context`
    `fmt`

    `google.golang.org/grpc`

    `github.com/dappledger/AnnChain/bcstore/proto`
    `github.com/dappledger/AnnChain/bcstore/types`
    log "github.com/sirupsen/logrus"
)

type RpcCLient struct {
    cli proto.BCStoreClient
    dbName string
}

func newRpcCLient(dbname string,addrs []string)(DBStore,error){
    log.Infof("newRpcClient(%s;%v)",dbname,addrs)
    cc,err:=grpc.Dial(addrs[0],grpc.WithInsecure())
    if err != nil {
        return nil,err
    }
    cli := proto.NewBCStoreClient(cc)
    return &RpcCLient{cli,dbname},nil
}


func (rc *RpcCLient)Get(key types.Key)(types.KeyValue,error){
    in := proto.ReqGet{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Key:                  key,
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    var kv types.KeyValue
    resp,err:=rc.cli.Get(context.TODO(),&in)
    if err != nil {
        return kv,err
    }
    kv.Key = resp.GetKv().GetKey()
    kv.Value = resp.GetKv().GetValue()
    return kv,nil
}

func (rc *RpcCLient)Set(kv types.KeyValue)error{
    in := proto.ReqSet{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Kv:                   &proto.KeyValue{Key:kv.Key,Value:kv.Value},
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    resp,err := rc.cli.Set(context.TODO(),&in)
    if err != nil {
        return err
    }
    if resp.ErrorCode != int64(types.ErrOK){
        return fmt.Errorf(resp.GetError())
    }
    return nil
}

//rpc GetByPrefix(ReqGetByPrefix)returns(RespGetByPrefix){}
func (rc *RpcCLient)GetByPrefix(prefix, lastKey types.Key, limit int64)([]types.KeyValue,error){
    if limit <= 0{
        limit = types.DefPageLimit
    }
    in := &proto.ReqGetByPrefix{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Prefix:               prefix,
        LastKey:              lastKey,
        Limit:                limit,
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    resp,err:=rc.cli.GetByPrefix(context.TODO(),in)
    if err != nil {
        return nil,err
    }
    if resp.GetBase().ErrorCode != int64(types.ErrOK){
        return nil,fmt.Errorf(resp.GetBase().GetError())
    }
    var ret = make([]types.KeyValue,len(resp.List))
    for i:=0;i<len(resp.List);i++{
        ret[i].Key = resp.List[i].GetKey()
        ret[i].Value = resp.List[i].GetValue()
    }
    return ret,nil
}
//rpc Batch(ReqBatch)returns(RespBase){}
func (rc *RpcCLient)Batch(dels []types.Key,sets []types.KeyValue)error{
    in := &proto.ReqBatch{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Sets:                 make([]*proto.KeyValue,len(sets)),
        Dels:                 make([][]byte,len(dels)),
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    for i:=0;i<len(dels);i++{
        in.Dels[i] = dels[i]
    }
    for i:=0;i<len(sets);i++{
        in.Sets[i] = &proto.KeyValue{}
        in.Sets[i].Key = sets[i].Key
        in.Sets[i].Value = sets[i].Value
    }
    resp,err := rc.cli.Batch(context.TODO(),in)
    if err != nil {
        return err
    }
    if resp.ErrorCode != int64(types.ErrOK){
        return fmt.Errorf(resp.GetError())
    }
    return nil
}
//rpc IsExist(ReqIsExist)returns(RespIsExist){}
func (rc *RpcCLient)Has(key types.Key)(bool,error){
    in := &proto.ReqIsExist{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Key:                  key,
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    var ok bool
    resp,err := rc.cli.IsExist(context.TODO(),in)
    if err != nil {
        return ok,err
    }
    if resp.GetBase().ErrorCode != int64(types.ErrOK){
        return ok,fmt.Errorf(resp.GetBase().GetError())
    }
    return resp.GetIsExist(),nil
}
//rpc GetProperty(ReqGetProperty)returns(RespGetProperty){}
func (rc *RpcCLient)GetProperty(property proto.PropertyType) (string, error){
    in := &proto.ReqGetProperty{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Property:             property,
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    resp,err := rc.cli.GetProperty(context.TODO(),in)
    if err != nil {
        return "",err
    }
    if resp.GetBase().ErrorCode != int64(types.ErrOK){
        return "",fmt.Errorf(resp.GetBase().GetError())
    }
    return resp.GetRet(),nil
}

func (rc *RpcCLient)Delete(key types.Key) error{
    in := &proto.ReqDelete{
        Base:                 &proto.ReqBase{},
        DbName:               rc.dbName,
        Key:                  key,
        XXX_NoUnkeyedLiteral: struct{}{},
        XXX_unrecognized:     nil,
        XXX_sizecache:        0,
    }
    resp,err:= rc.cli.Delete(context.TODO(),in)
    if err != nil {
        return err
    }
    if resp.ErrorCode != int64(types.ErrOK){
        return fmt.Errorf(resp.GetError())
    }
    return nil
}
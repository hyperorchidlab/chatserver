package db

import (
	"encoding/json"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/hdb"
	"sync"
	"log"
	"github.com/kprc/chat-protocol/groupid"
	"github.com/btcsuite/btcutil/base58"
	"crypto/sha256"
	"strconv"
	"fmt"
)

type GroupMsgHDB struct {
	hdb.HistoryDBIntf
	dbLock sync.Mutex
	cursor *db.DBCusor
}

const (
	MemHistoryCount = 1000
)

var (
	gmhdb     *GroupMsgHDB
	gmhdbLock sync.Mutex
)

type GroupMsg struct {
	AesKey string `json:"aes_key"`
	Msg    string `json:"msg"`
}

type LabelGroupMsg struct {
	AesKey string `json:"aes_key"`
	Msg    string `json:"msg"`
	Cnt    int    `json:"cnt"`
}

func newChatGroupMsgDb() *GroupMsgHDB {
	cfg := config.GetCSC()
	db := hdb.New(MemHistoryCount, cfg.GetGroupMsgDbPath()).Load()

	gmdb := &GroupMsgHDB{HistoryDBIntf: db}

	return gmdb
}

func GetGMsgDb() *GroupMsgHDB {
	if gmhdb == nil {
		gmhdbLock.Lock()
		defer gmhdbLock.Unlock()

		if gmhdb == nil {
			gmhdb = newChatGroupMsgDb()
		}
	}

	return gmhdb
}

func (gmdb *GroupMsgHDB) Insert(id string, keyHash string, cipherTxt string) {

	grpKeyDb:=GetChatGrpKeysDb()

	keys:=grpKeyDb.Find(keyHash)
	if keys == nil{
		log.Println("No group key in db , drop the message")
		return
	}

	gm := &GroupMsg{AesKey: keyHash, Msg: cipherTxt}

	j, _ := json.Marshal(*gm)

	gmdb.dbLock.Lock()
	defer gmdb.dbLock.Unlock()

	idx,_:=gmdb.HistoryDBIntf.Insert(id, string(j))

	for i:=0;i<len(keys.PubKeys);i++{
		pk:=keys.PubKeys[i]
		u:=groupid.GrpID(id).ToBytes()
		u = append(u,base58.Decode(pk)...)

		hash:=sha256.Sum256(u)

		gmdb.HistoryDBIntf.Insert(base58.Encode(hash[:]),strconv.Itoa(idx))
	}

}

func (gmdb *GroupMsgHDB)FindMsg2(gid groupid.GrpID, pk string, begin, n int) (msgs []*LabelGroupMsg) {
	gmdb.dbLock.Lock()
	defer gmdb.dbLock.Unlock()

	if _,err := gmdb.HistoryDBIntf.FindBlock(gid.String()); err!=nil{
		return nil
	}

	u:=groupid.GrpID(gid).ToBytes()
	u=append(u,base58.Decode(pk)...)
	hash:=sha256.Sum256(u)
	fid:=base58.Encode(hash[:])

	if _,err:=gmdb.HistoryDBIntf.FindBlock(fid);err!=nil{
		return nil
	}

	r,err:=gmdb.HistoryDBIntf.Find(fid,begin,n)
	if err != nil || len(r) == 0{
		return nil
	}

	



	return nil


}

type Section struct {
	begin, to int
}

func (s *Section)String() string {
	return  fmt.Sprintf("begin: %-8d To: %-8d",s.begin,s.to)
}

func Discrete2Section(discretes []int) []Section  {

	if len(discretes) == 0{
		return nil
	}

	prev:=discretes[0]

	var secs []Section

	sec := Section{begin:prev}

	for i:=1;i<len(discretes);i++{
		if prev + 1 == discretes[i]{
			prev ++
			continue
		}
		sec.to = prev
		prev = discretes[i]

		secs = append(secs,sec)

		sec = Section{begin:prev}
	}

	sec.to = prev

	secs = append(secs,sec)

	return secs
}


func (gmdb *GroupMsgHDB) FindMsg(id string, begin, n int) (msgs []*LabelGroupMsg) {
	gmdb.dbLock.Lock()
	defer gmdb.dbLock.Unlock()

	if _, err := gmdb.HistoryDBIntf.FindBlock(id); err != nil {
		return nil
	}

	r, err := gmdb.HistoryDBIntf.Find(id, begin, n)
	if err != nil || len(r) == 0 {
		return nil
	}

	for i := 0; i < len(r); i++ {
		v := r[i]

		gm := &GroupMsg{}

		json.Unmarshal([]byte(v.V), gm)
		lgm := &LabelGroupMsg{}
		lgm.Msg = gm.Msg
		lgm.AesKey = gm.AesKey
		lgm.Cnt = v.Cnt

		msgs = append(msgs, lgm)

	}

	return
}



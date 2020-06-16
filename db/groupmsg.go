package db

import (
	"encoding/json"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/hdb"
	"sync"
)

type GroupMsgHDB struct {
	hdb.HistoryDBIntf
	dbLock sync.Mutex
	cursor *db.DBCusor
}

const (
	MemHistoryCount = 200
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

	gm := &GroupMsg{AesKey: keyHash, Msg: cipherTxt}

	j, _ := json.Marshal(*gm)

	gmdb.dbLock.Lock()
	defer gmdb.dbLock.Unlock()

	gmdb.HistoryDBIntf.Insert(id, string(j))

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

		_ := json.Unmarshal([]byte(v.V), gm)
		lgm := &LabelGroupMsg{}
		lgm.Msg = gm.Msg
		lgm.AesKey = gm.AesKey
		lgm.Cnt = v.Cnt

		msgs = append(msgs, lgm)

	}

	return
}

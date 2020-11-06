package db

import (
	"encoding/json"
	"github.com/hyperorchidlab/chatserver/config"
	"github.com/hyperorchidlab/chatserver/db/hdb"
	"sync"
)

type P2PMsgHDB struct {
	hdb.HistoryDBIntf
	dbLock sync.Mutex
	cursor *DBCusor
}

const (
	MemHistoryP2PCount = 200
)

var (
	p2phdb    *P2PMsgHDB
	p2pdbLock sync.Mutex
)

type P2PMsg struct {
	AesKey string `json:"aes_key"`
	Msg    string `json:"msg"`
}

type LabelP2pMsg struct {
	AesKey string `json:"aes_key"`
	Msg    string `json:"msg"`
	Cnt    int    `json:"cnt"`
}

func newChatP2PMsgDb() *P2PMsgHDB {
	cfg := config.GetCSC()
	db := hdb.New(MemHistoryP2PCount, cfg.GetP2PMsgDbPath()).Load()

	pdb := &P2PMsgHDB{HistoryDBIntf: db}

	return pdb
}

func GetP2PMsgDb() *P2PMsgHDB {
	if p2phdb == nil {
		p2pdbLock.Lock()
		defer p2pdbLock.Unlock()

		if p2phdb == nil {
			p2phdb = newChatP2PMsgDb()
		}
	}

	return p2phdb
}

func (pdb *P2PMsgHDB) Insert(id string, pk string, cipherTxt string) {

	pm := &P2PMsg{AesKey: pk, Msg: cipherTxt}

	j, _ := json.Marshal(*pm)

	pdb.dbLock.Lock()
	defer pdb.dbLock.Unlock()

	pdb.HistoryDBIntf.Insert(id, string(j))
}

func (pdb *P2PMsgHDB) FindMsg(id string, begin, n int) (msgs []*LabelP2pMsg) {
	pdb.dbLock.Lock()
	defer pdb.dbLock.Unlock()

	if _, err := pdb.HistoryDBIntf.FindBlock(id); err != nil {
		return nil
	}

	r, err := pdb.HistoryDBIntf.Find(id, begin, n)
	if err != nil || len(r) == 0 {
		return nil
	}

	for i := 0; i < len(r); i++ {
		v := r[i]

		pm := &P2PMsg{}

		json.Unmarshal([]byte(v.V), pm)

		lm := &LabelP2pMsg{}
		lm.Msg = pm.Msg
		lm.AesKey = pm.AesKey
		lm.Cnt = v.Cnt

		msgs = append(msgs, lm)

	}

	return
}

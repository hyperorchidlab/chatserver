package db

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chat-protocol/groupid"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/hdb"
	"log"
	"strconv"
	"sync"
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
	AesKey string              `json:"aes_key"`
	Msg    string              `json:"msg"`
	Speek  address.ChatAddress `json:"seepk"`
}

type LabelGroupMsg struct {
	AesKey string              `json:"aes_key"`
	Msg    string              `json:"msg"`
	Speek  address.ChatAddress `json:"speek"`
	Cnt    int                 `json:"cnt"`
	UCnt   int 				   `json:"u_cnt"`
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

func (gmdb *GroupMsgHDB) Insert(gid groupid.GrpID, keyHash string, speek address.ChatAddress, cipherTxt string) {

	grpKeyDb := GetChatGrpKeysDb()

	keys := grpKeyDb.Find(keyHash)
	if keys == nil {
		log.Println("No group key in db , drop the message")
		return
	}

	gm := &GroupMsg{AesKey: keyHash, Msg: cipherTxt, Speek: speek}

	j, _ := json.Marshal(*gm)

	gmdb.dbLock.Lock()
	defer gmdb.dbLock.Unlock()

	idx, _ := gmdb.HistoryDBIntf.Insert(gid.String(), string(j))

	log.Println("gid:", gid.String())

	for i := 0; i < len(keys.PubKeys); i++ {
		pk := keys.PubKeys[i]
		var ku []byte
		u := gid.ToBytes()
		ku = append(ku, u...)
		pku := base58.Decode(pk)
		ku = append(ku, pku...)

		hash := sha256.Sum256(ku)
		k := base58.Encode(hash[:])
		log.Println(base58.Encode(u), pk, k, idx)
		gmdb.HistoryDBIntf.Insert(k, strconv.Itoa(idx))
	}

}

type SecIdx struct {
	Idx int
	VIdx int
}

func (gmdb *GroupMsgHDB) FindMsg2(gid groupid.GrpID, pk string, begin, n int) (msgs []*LabelGroupMsg) {
	gmdb.dbLock.Lock()
	defer gmdb.dbLock.Unlock()

	if _, err := gmdb.HistoryDBIntf.FindBlock(gid.String()); err != nil {
		log.Println(gid.String(), err)
		return nil
	}

	u := gid.ToBytes()
	var ku []byte
	ku = append(ku, u...)
	pku := address.ChatAddress(pk).GetBytes()
	ku = append(ku, pku...)
	hash := sha256.Sum256(ku)
	fid := base58.Encode(hash[:])

	log.Println("gid", gid.String())

	log.Println(base58.Encode(ku), base58.Encode(pku), fid)

	if _, err := gmdb.HistoryDBIntf.FindBlock(fid); err != nil {
		log.Println(fid, err)
		return nil
	}

	log.Println("--->",fid,begin,n)

	r, err := gmdb.HistoryDBIntf.Find(fid, begin, n)
	if err != nil || len(r) == 0 {
		log.Println(fid, err, len(r))
		return nil
	}
	var dc []SecIdx
	for i := 0; i < len(r); i++ {
		cnt, _ := strconv.Atoi(r[i].V)
		dc = append(dc, SecIdx{Idx: r[i].Cnt,VIdx: cnt})
	}
	ses := Discrete2Section(dc)

	log.Println(ses)

	return gmdb.findMsgBySecs(gid.String(), ses,dc)
}

func findIdx(vidx int, dc []SecIdx) int  {
	for i:=0;i<len(dc);i++{
		if dc[i].VIdx == vidx{
			return dc[i].Idx
		}
	}

	return 0
}

func (gmdb *GroupMsgHDB) findMsgBySecs(id string, ses []Section, dc []SecIdx) (msgs []*LabelGroupMsg) {

	for i := 0; i < len(ses); i++ {
		r, err := gmdb.HistoryDBIntf.Find(id, ses[i].begin.VIdx, (ses[i].to.VIdx-ses[i].begin.VIdx + 1))
		if err != nil || len(r) == 0 {
			continue
		}

		for j := 0; j < len(r); j++ {
			v := r[j]

			gm := &GroupMsg{}

			json.Unmarshal([]byte(v.V), gm)
			lgm := &LabelGroupMsg{}
			lgm.Msg = gm.Msg
			lgm.AesKey = gm.AesKey
			lgm.Cnt = v.Cnt
			lgm.Speek = gm.Speek
			lgm.UCnt = findIdx(v.Cnt,dc)

			msgs = append(msgs, lgm)
		}
	}

	return
}

type Section struct {
	begin, to SecIdx
}

func (s *Section) String() string {
	return fmt.Sprintf("begin: %-8d To: %-8d", s.begin, s.to)
}

func Discrete2Section(discretes []SecIdx) []Section {

	if len(discretes) == 0 {
		return nil
	}

	prev := discretes[0].VIdx
	prevV := discretes[0]

	var secs []Section

	sec := Section{begin: prevV}

	for i := 1; i < len(discretes); i++ {
		if prev+1 == discretes[i].VIdx {
			prev++
			prevV = discretes[i]
			continue
		}
		sec.to = prevV
		prev = discretes[i].VIdx
		prevV = discretes[i]

		secs = append(secs, sec)

		sec = Section{begin: prevV}
	}

	sec.to = prevV

	secs = append(secs, sec)

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

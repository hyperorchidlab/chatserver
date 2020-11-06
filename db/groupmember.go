package db

import (
	"encoding/json"
	"github.com/hyperorchidlab/chatserver/app/cmdcommon"
	"github.com/hyperorchidlab/chatserver/config"
	"github.com/pkg/errors"
	"sync"
)

type ChatGroupMemberDB struct {
	NbsDbInter
	dbLock sync.Mutex
	cursor *DBCusor
}

var (
	cgmStore     *ChatGroupMemberDB
	cgmStoreLock sync.Mutex
)

type GroupMember struct {
	GrpID      string   `json:"-"`
	Owner      string   `json:"owr"`
	Members    []string `json:"mbrs"`
	CreateTime int64    `json:"ct"`
	UpdateTime int64    `json:"ut"`
}

func newChatGroupMemberDB() *ChatGroupMemberDB {
	cfg := config.GetCSC()
	db := NewFileDb(cfg.GetGrpMbrsDbPath()).Load()

	return &ChatGroupMemberDB{NbsDbInter: db}
}

func GetChatGrpMbrsDB() *ChatGroupMemberDB {
	if cgmStore == nil {
		cgmStoreLock.Lock()
		defer cgmStoreLock.Unlock()

		if cgmStore == nil {
			cgmStore = newChatGroupMemberDB()
		}

	}

	return cgmStore
}

func (cgm *ChatGroupMemberDB) Insert(grpId string, owner string) error {
	cgm.dbLock.Lock()
	defer cgm.dbLock.Unlock()

	if _, err := cgm.NbsDbInter.Find(grpId); err == nil {
		return errors.New("group id is in db")
	}
	now := cmdcommon.GetNowMsTime()
	gm := &GroupMember{}
	gm.Owner = owner
	gm.Members = append(gm.Members, owner)
	gm.CreateTime = now
	gm.UpdateTime = now

	if v, err := json.Marshal(*gm); err != nil {
		return err
	} else {
		return cgm.NbsDbInter.Insert(grpId, string(v))
	}
}

func (cgm *ChatGroupMemberDB) AddMember(grpId string, mbr string) error {
	cgm.dbLock.Lock()
	defer cgm.dbLock.Unlock()

	if v, err := cgm.NbsDbInter.Find(grpId); err != nil {
		return err
	} else {
		gm := &GroupMember{}

		err = json.Unmarshal([]byte(v), gm)
		if err != nil {
			return err
		}

		for i := 0; i < len(gm.Members); i++ {
			if mbr == gm.Members[i] {
				return nil
			}
		}

		gm.Members = append(gm.Members, mbr)
		gm.UpdateTime = cmdcommon.GetNowMsTime()

		var vv []byte

		if vv, err = json.Marshal(*gm); err != nil {
			return err
		}

		cgm.NbsDbInter.Update(grpId, string(vv))

	}

	return nil
}

func (cgm *ChatGroupMemberDB) DelMember(grpId string, mbr string) error {
	cgm.dbLock.Lock()
	defer cgm.dbLock.Unlock()

	if v, err := cgm.NbsDbInter.Find(grpId); err != nil {
		return err
	} else {
		gm := &GroupMember{}

		err = json.Unmarshal([]byte(v), gm)
		if err != nil {
			return err
		}

		var dflag bool
		for i := 0; i < len(gm.Members); i++ {
			if mbr == gm.Members[i] {
				if i != len(gm.Members)-1 {
					gm.Members[i] = gm.Members[len(gm.Members)-1]
				}
				gm.Members = gm.Members[:len(gm.Members)-1]
				dflag = true
				break
			}
		}

		if !dflag {
			return nil
		}

		gm.UpdateTime = cmdcommon.GetNowMsTime()

		var vv []byte

		if vv, err = json.Marshal(*gm); err != nil {
			return err
		}

		cgm.NbsDbInter.Update(grpId, string(vv))

	}

	return nil

}

func (cgm *ChatGroupMemberDB) DelGroupMember(grpId string) error {
	cgm.dbLock.Lock()
	defer cgm.dbLock.Unlock()

	if v, err := cgm.NbsDbInter.Find(grpId); err != nil {
		return nil
	} else {
		gm := &GroupMember{}

		err = json.Unmarshal([]byte(v), gm)
		if err != nil {
			return err
		}

		if len(gm.Members) > 0 {
			return errors.New("group have members")
		}

		cgm.NbsDbInter.Delete(grpId)

	}

	return nil
}

func (cgm *ChatGroupMemberDB) Find(grpId string) (*GroupMember, error) {
	cgm.dbLock.Lock()
	defer cgm.dbLock.Unlock()

	if vs, err := cgm.NbsDbInter.Find(grpId); err != nil {
		return nil, err
	} else {
		f := &GroupMember{}
		err = json.Unmarshal([]byte(vs), f)
		if err != nil {
			return nil, err
		} else {
			f.GrpID = grpId

			return f, nil
		}
	}

}

func (cgm *ChatGroupMemberDB) Iterator() {

	cgm.dbLock.Lock()
	defer cgm.dbLock.Unlock()

	cgm.cursor = cgm.NbsDbInter.DBIterator()
}

func (cgm *ChatGroupMemberDB) Next() (key string, meta *GroupMember, r1 error) {
	if cgm.cursor == nil {
		return
	}
	cgm.dbLock.Lock()
	//cgm.dbLock.Unlock()
	k, v := cgm.cursor.Next()
	if k == "" {
		cgm.dbLock.Unlock()
		return "", nil, nil
	}
	cgm.dbLock.Unlock()
	meta = &GroupMember{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return "", nil, err
	}
	meta.GrpID = k

	key = k

	return

}

package db

import (
	"encoding/json"
	"errors"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type ChatGroupsDB struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	cgStore     *ChatGroupsDB
	cgStoreLock sync.Mutex
)

type Group struct {
	Alias      string `json:"as"`
	GrpId      string `json:"-"`
	RefCnt     int    `json:"rcnt"`
	Owner      string `json:"owner"`
	CreateTime int64  `json:"ct"`
	UpdateTime int64  `json:"ut"`
}

func newChatGroupDb() *ChatGroupsDB {
	cfg := config.GetCSC()
	db := db.NewFileDb(cfg.GetGroupsDbPath()).Load()

	return &ChatGroupsDB{NbsDbInter: db}
}

func GetChatGroupsDB() *ChatGroupsDB {
	if cgStore == nil {
		cgStoreLock.Lock()
		defer cgStoreLock.Unlock()

		if cgStore == nil {
			cgStore = newChatGroupDb()
		}

	}
	return cgStore
}

func (cg *ChatGroupsDB) Insert(grpID string, alias, owner string) error {
	cg.dbLock.Lock()
	defer cg.dbLock.Unlock()

	if _, err := cg.NbsDbInter.Find(grpID); err == nil {
		return errors.New("group id is in db")
	}
	now := tools.GetNowMsTime()
	g := &Group{}
	g.Alias = alias
	g.CreateTime = now
	g.UpdateTime = now
	g.Owner = owner
	g.RefCnt = 0

	if v, err := json.Marshal(*g); err != nil {
		return err
	} else {
		return cg.NbsDbInter.Insert(grpID, string(v))
	}

}

func (cg *ChatGroupsDB) UpdateAlias(grpId string, alias string) error {
	cg.dbLock.Lock()
	defer cg.dbLock.Unlock()

	if vs, err := cg.NbsDbInter.Find(grpId); err != nil {
		return err
	} else {
		g := &Group{}
		err = json.Unmarshal([]byte(vs), g)
		if err != nil {
			return err
		} else {
			g.GrpId = grpId
		}

		now := tools.GetNowMsTime()
		g.Alias = alias

		g.UpdateTime = now

		if v, err := json.Marshal(*g); err != nil {
			return err
		} else {
			cg.NbsDbInter.Update(grpId, string(v))
		}
	}

	return nil
}

func (cg *ChatGroupsDB) IncRefer(grpId string) error {
	cg.dbLock.Lock()
	defer cg.dbLock.Unlock()
	if vs, err := cg.NbsDbInter.Find(grpId); err != nil {
		return err
	} else {
		g := &Group{}
		err = json.Unmarshal([]byte(vs), g)
		if err != nil {
			return err
		}
		g.RefCnt++

		if v, err := json.Marshal(*g); err != nil {
			return err
		} else {
			cg.NbsDbInter.Update(grpId, string(v))
		}
	}

	return nil

}

func (cg *ChatGroupsDB) DecRefer(grpId string) error {
	cg.dbLock.Lock()
	defer cg.dbLock.Unlock()
	if vs, err := cg.NbsDbInter.Find(grpId); err != nil {
		return err
	} else {
		g := &Group{}
		err = json.Unmarshal([]byte(vs), g)
		if err != nil {
			return err
		}
		g.RefCnt--

		if g.RefCnt <= 0 {
			cg.NbsDbInter.Delete(grpId)
			return nil
		} else {
			if v, err := json.Marshal(*g); err != nil {
				return err
			} else {
				cg.NbsDbInter.Update(grpId, string(v))
			}
		}
	}

	return nil

}

func (cg *ChatGroupsDB) Find(grpId string) (*Group, error) {
	cg.dbLock.Unlock()
	defer cg.dbLock.Unlock()

	if vs, err := cg.NbsDbInter.Find(grpId); err != nil {
		return nil, err
	} else {
		f := &Group{}
		err = json.Unmarshal([]byte(vs), f)
		if err != nil {
			return nil, err
		} else {
			f.GrpId = grpId

			return f, nil
		}
	}

}

func (s *ChatGroupsDB) Iterator() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.cursor = s.NbsDbInter.DBIterator()
}

func (s *ChatGroupsDB) Next() (key string, meta *Group, r1 error) {
	if s.cursor == nil {
		return
	}
	s.dbLock.Lock()
	//s.dbLock.Unlock()
	k, v := s.cursor.Next()
	if k == "" {
		s.dbLock.Unlock()
		return "", nil, nil
	}
	s.dbLock.Unlock()
	meta = &Group{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return "", nil, err
	}
	meta.GrpId = k

	key = k

	return

}

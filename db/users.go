package db

import (
	"encoding/json"
	"errors"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"time"
)

type ChatUsersDB struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cusor  *db.DBCusor
}

var (
	cuStore     *ChatUsersDB
	cuStoreLock sync.Mutex
)

type ChatUser struct {
	Alias      string `json:"as"`
	PubKey     string `json:"-"`
	CreateTime int64  `json:"ct"`
	UpdateTime int64  `json:"ut"`
	ExpireTime int64  `json:"et"`
}

func newChatUserDb() *ChatUsersDB {
	cfg := config.GetCSC()
	db := db.NewFileDb(cfg.GetUsersDbPath()).Load()

	return &ChatUsersDB{NbsDbInter: db}
}

func GetChatUserDB() *ChatUsersDB {
	if cuStore == nil {
		cuStoreLock.Lock()
		defer cuStoreLock.Unlock()

		if cuStore == nil {
			cuStore = newChatUserDb()
		}
	}

	return cuStore
}

func (s *ChatUsersDB) Insert(alias string, pubkey string, tv int64) error {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if _, err := s.NbsDbInter.Find(pubkey); err == nil {
		return errors.New("insert Failed, pubKey is in db")
	}

	now := tools.GetNowMsTime()
	cu := &ChatUser{}
	cu.Alias = alias
	cu.PubKey = pubkey
	cu.CreateTime = now
	cu.UpdateTime = now

	nowtm := time.Now().AddDate(0, int(tv), 0)

	cu.ExpireTime = nowtm.UnixNano() / 1e6

	if v, err := json.Marshal(*cu); err != nil {
		return err
	} else {
		return s.NbsDbInter.Insert(pubkey, string(v))
	}
}

func (s *ChatUsersDB) Update(alias string, pubkey string, tv int64) error {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(pubkey); err != nil {
		return err
	} else {
		cu := &ChatUser{}
		err = json.Unmarshal([]byte(vs), cu)
		if err != nil {
			return err
		} else {
			cu.PubKey = pubkey
		}

		now := tools.GetNowMsTime()
		cu.Alias = alias
		cu.PubKey = pubkey
		//cu.CreateTime = now
		cu.UpdateTime = now

		if now > cu.ExpireTime {
			cu.ExpireTime = (time.Now().AddDate(0, int(tv), 0).UnixNano()) / 1e6
		} else {
			sec := cu.ExpireTime / 1000
			nsec := (cu.ExpireTime - sec*1000) * 1e6
			cu.ExpireTime = (time.Unix(sec, nsec).AddDate(0, int(tv), 0).UnixNano()) / 1e6
		}

		if v, err := json.Marshal(*cu); err != nil {
			return err
		} else {
			s.NbsDbInter.Update(pubkey, string(v))
		}
	}

	return nil

}

func (s *ChatUsersDB) UpdateExpireTime(pubkey string, tv int64) error {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(pubkey); err != nil {
		return err
	} else {
		cu := &ChatUser{}
		err = json.Unmarshal([]byte(vs), cu)
		if err != nil {
			return err
		} else {
			cu.PubKey = pubkey
		}

		now := tools.GetNowMsTime()
		//cu.PubKey = pubkey
		//cu.CreateTime = now
		cu.UpdateTime = now

		if now > cu.ExpireTime {
			cu.ExpireTime = now + tv
		} else {
			cu.ExpireTime += tv
		}

		if v, err := json.Marshal(*cu); err != nil {
			return err
		} else {
			s.NbsDbInter.Update(pubkey, string(v))
		}
	}
	return nil
}

func (s *ChatUsersDB) UpdateAlias(pubkey string, alias string) error {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(pubkey); err != nil {
		return err
	} else {
		cu := &ChatUser{}
		err = json.Unmarshal([]byte(vs), cu)
		if err != nil {
			return err
		} else {
			cu.PubKey = pubkey
		}

		now := tools.GetNowMsTime()
		//cu.PubKey = pubkey
		//cu.CreateTime = now
		cu.UpdateTime = now
		cu.Alias = alias
		//
		//if now > cu.ExpireTinme{
		//	cu.ExpireTinme = now+tv
		//}else{
		//	cu.ExpireTinme += tv
		//}

		if v, err := json.Marshal(*cu); err != nil {
			return err
		} else {
			s.NbsDbInter.Update(pubkey, string(v))
		}
	}
	return nil
}

func (s *ChatUsersDB) Find(pk string) (*ChatUser, error) {
	s.dbLock.Unlock()
	defer s.dbLock.Unlock()

	if vs, err := s.NbsDbInter.Find(pk); err != nil {
		return nil, err
	} else {
		f := &ChatUser{}
		err = json.Unmarshal([]byte(vs), f)
		if err != nil {
			return nil, err
		} else {
			f.PubKey = pk

			return f, nil
		}
	}

}

func (s *ChatUsersDB) Remove(pk string) {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Delete(pk)
}

func (s *ChatUsersDB) Save() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Save()
}

func (s *ChatUsersDB) Iterator() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.cusor = s.NbsDbInter.DBIterator()
}

func (s *ChatUsersDB) Next() (key string, meta *ChatUser, r1 error) {
	if s.cusor == nil {
		return
	}
	s.dbLock.Lock()
	//s.dbLock.Unlock()
	k, v := s.cusor.Next()
	if k == "" {
		s.dbLock.Unlock()
		return "", nil, nil
	}
	s.dbLock.Unlock()
	meta = &ChatUser{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return "", nil, err
	}
	meta.PubKey = k

	key = k

	return

}

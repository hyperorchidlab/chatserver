package db

import (
	"encoding/json"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/nbsnetwork/db"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"errors"
)

type ChatFriendsDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}

var (
	cfStore     *ChatFriendsDb
	cfStoreLock sync.Mutex
)

type Friend struct {
	OwnerPk string `json:"-"`
	PubKey  string `json:"pk"`
	AddTime int64  `json:"at"`
	Agree   bool   `json:"agr"`
}

type ChatFriends struct {
	Owner string	`json:"-"`
	Friends []Friend `json:"fs"`
	Groups  []string `json:"gs"`
}

func newChatFriendsDB() *ChatFriendsDb {
	cfg := config.GetCSC()

	db := db.NewFileDb(cfg.GetFriendsDbPath()).Load()

	return &ChatFriendsDb{NbsDbInter: db}
}

func GetChatFriendsDB() *ChatFriendsDb {
	if cfStore == nil {
		cfStoreLock.Lock()
		defer cfStoreLock.Unlock()

		if cfStore == nil {
			cfStore = newChatFriendsDB()
		}
	}

	return cfStore
}

func (cf *ChatFriendsDb) AddFriend(ownerPk string, friendPK string) error {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	cfs := &ChatFriends{}

	if vs, err := cf.NbsDbInter.Find(ownerPk); err == nil {
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return err
		}

	}
	cfs.Owner = ownerPk

	for i := 0; i < len(cfs.Friends); i++ {
		if friendPK == cfs.Friends[i].PubKey {
			return nil
		}
	}

	f := &Friend{}
	f.PubKey = friendPK
	f.AddTime = tools.GetNowMsTime()

	cfs.Friends = append(cfs.Friends, *f)

	if v, err := json.Marshal(cfs); err != nil {
		return err
	} else {
		cf.Update(ownerPk, string(v))
	}

	return nil
}

func (cf *ChatFriendsDb)AgreeFriend(ownerPk string,friendPk string,agree bool) error  {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	cfs := &ChatFriends{}

	if vs, err := cf.NbsDbInter.Find(ownerPk); err == nil {
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return err
		}

	}
	cfs.Owner = ownerPk

	var i int
	for i = 0; i < len(cfs.Friends); i++{
		if friendPk == cfs.Friends[i].PubKey{
			cfs.Friends[i].Agree = agree
		}
	}

	if i == len(cfs.Friends){
		f := &Friend{}
		f.PubKey = friendPk
		f.AddTime = tools.GetNowMsTime()
		f.Agree = agree

		cfs.Friends = append(cfs.Friends,*f)
	}

	if v,err:=json.Marshal(cfs); err!=nil{
		return err
	}else{
		cf.NbsDbInter.Update(ownerPk,string(v))
	}

	return nil

}


func (cf *ChatFriendsDb) DelFriend(ownerPK string, friendPK string) error {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	var cfs *ChatFriends

	if vs, err := cf.NbsDbInter.Find(ownerPK); err == nil {
		cfs = &ChatFriends{}
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return err
		}
		for i := 0; i < len(cfs.Friends); i++ {
			if cfs.Friends[i].PubKey == friendPK {
				cfs.Friends = append(cfs.Friends[:i], cfs.Friends[i+1:]...)
				return nil
			}
		}

	}

	return nil

}

func (cf *ChatFriendsDb) AddGroup(ownerPk string, group string) error {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	var cfs *ChatFriends

	if vs, err := cf.NbsDbInter.Find(ownerPk); err == nil {
		cfs = &ChatFriends{}
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return err
		}
	} else {
		cfs = &ChatFriends{}
	}

	for i := 0; i < len(cfs.Groups); i++ {
		if group == cfs.Groups[i] {
			return nil
		}
	}

	cfs.Groups = append(cfs.Groups, group)

	if v, err := json.Marshal(cfs); err != nil {
		return err
	} else {
		cf.Update(ownerPk, string(v))
	}

	return nil

}

func (cf *ChatFriendsDb) DelGroup(ownerPK string, group string) error {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	var cfs *ChatFriends

	if vs, err := cf.NbsDbInter.Find(ownerPK); err == nil {
		cfs = &ChatFriends{}
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return err
		}
		for i := 0; i < len(cfs.Groups); i++ {
			if cfs.Groups[i] == group {
				cfs.Groups = append(cfs.Groups[:i], cfs.Groups[i+1:]...)
				return nil
			}
		}

	}

	return nil

}

func (cf *ChatFriendsDb) Find(ownerPk string) (*ChatFriends, error) {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	if vs, err := cf.NbsDbInter.Find(ownerPk); err != nil {
		return nil, err
	} else {
		cfs := &ChatFriends{}
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return nil, err
		}
		cfs.Owner = ownerPk
		return cfs, nil
	}
}

func (cf *ChatFriendsDb)FindFriend(ownerPk string, friend string) (*Friend,error)  {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	if vs, err := cf.NbsDbInter.Find(ownerPk); err != nil {
		return nil, err
	} else {
		cfs := &ChatFriends{}
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return nil, err
		}
		//cfs.Owner = ownerPk
		for i:=0;i<len(cfs.Friends);i++{
			if cfs.Friends[i].PubKey == friend{
				return &cfs.Friends[i],err
			}
		}
	}

	return nil,errors.New("Not Found")
}

func (cf *ChatFriendsDb)FindGroup(ownerPk string,grpId string) (string,error)  {
	cf.dbLock.Lock()
	defer cf.dbLock.Unlock()

	if vs, err := cf.NbsDbInter.Find(ownerPk); err != nil {
		return "", err
	} else {
		cfs := &ChatFriends{}
		if err = json.Unmarshal([]byte(vs), cfs); err != nil {
			return "", err
		}
		//cfs.Owner = ownerPk
		for i:=0;i<len(cfs.Groups);i++{
			if cfs.Groups[i] == grpId{
				return cfs.Groups[i],err
			}
		}
	}

	return "",errors.New("Not Found")
}


func (s *ChatFriendsDb) Save() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.NbsDbInter.Save()
}

func (s *ChatFriendsDb) Iterator() {

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	s.cursor = s.NbsDbInter.DBIterator()
}

func (s *ChatFriendsDb) Next() (key string, meta *ChatFriends, r1 error) {
	if s.cursor == nil {
		return
	}
	s.dbLock.Lock()
	s.dbLock.Unlock()
	k, v := s.cursor.Next()
	if k == "" {
		s.dbLock.Unlock()
		return "", nil, nil
	}
	s.dbLock.Unlock()
	meta = &ChatFriends{}

	if err := json.Unmarshal([]byte(v), meta); err != nil {
		return "", nil, err
	}

	key = k

	return

}

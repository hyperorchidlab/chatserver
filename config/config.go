package config

import (
	"crypto/ed25519"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"os"
	"path"
	"sync"
)

const (
	CS_HomeDir      = ".chatserver"
	CS_CFG_FileName = "chatserver.json"
)

type ChatServerConfig struct {
	MgtHttpPort int `json:"mgthttpport"`

	CmdListenPort string `json:"cmdlistenport"`

	DBPath          string `json:"dbpath"`
	UsersDBFile     string `json:"usersdbfile"`
	FriendsDBFile   string `json:"friendsdbfile"`
	GroupsDBFile    string `json:"groupsdbfile"`
	GrpMemberDBFile string `json:"grpmemberdbfile"`

	ChatListenPort int `json:"chatport"`
	//ChatMgmtPort   int			`json:"chatmgmtport"`
	KeyFile string `json:"keyfile"`

	PrivKey ed25519.PrivateKey `json:"-"`
	PubKey  ed25519.PublicKey  `json:"-"`
}

var (
	cscfgInst     *ChatServerConfig
	cscfgInstLock sync.Mutex
)

func (bc *ChatServerConfig) InitCfg() *ChatServerConfig {
	bc.MgtHttpPort = 50818
	bc.CmdListenPort = "127.0.0.1:59527"
	bc.DBPath = "/db"
	bc.UsersDBFile = "users.db"
	bc.FriendsDBFile = "friends.db"
	bc.GroupsDBFile = "groups.db"
	bc.GrpMemberDBFile = "grpm.db"
	//bc.ChatMgmtPort = 39527
	bc.ChatListenPort = 39527
	bc.KeyFile = "chat_server.key"

	return bc
}

func (bc *ChatServerConfig) Load() *ChatServerConfig {
	if !tools.FileExists(GetCSCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetCSCFGFile())
	if err != nil {
		log.Println("load file failed", err)
		return nil
	}

	//bc1:=&BASDConfig{}

	err = json.Unmarshal(jbytes, bc)
	if err != nil {
		log.Println("load configuration unmarshal failed", err)
		return nil
	}

	return bc

}

func newSSCCfg() *ChatServerConfig {

	bc := &ChatServerConfig{}

	bc.InitCfg()

	return bc
}

func GetCSC() *ChatServerConfig {
	if cscfgInst == nil {
		cscfgInstLock.Lock()
		defer cscfgInstLock.Unlock()
		if cscfgInst == nil {
			cscfgInst = newSSCCfg()
		}
	}

	return cscfgInst
}

func PreLoad() *ChatServerConfig {
	bc := &ChatServerConfig{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *ChatServerConfig {
	bc := &ChatServerConfig{}

	bc.InitCfg()

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	cscfgInstLock.Lock()
	defer cscfgInstLock.Unlock()
	cscfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *ChatServerConfig) *ChatServerConfig) *ChatServerConfig {
	cscfgInstLock.Lock()
	defer cscfgInstLock.Unlock()

	lbc := newSSCCfg().Load()

	if lbc != nil {
		cscfgInst = lbc
	} else {
		lbc = newSSCCfg()
	}

	cscfgInst = initfromcmd(lbc)

	return cscfgInst
}

func GetCSCHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, CS_HomeDir)
}

func GetCSCFGFile() string {
	return path.Join(GetCSCHomeDir(), CS_CFG_FileName)
}

func (bc *ChatServerConfig) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetCSCHomeDir()) {
		os.MkdirAll(GetCSCHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetCSCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func (bc *ChatServerConfig) getDbPath() string {
	dbpath := path.Join(GetCSCHomeDir(), bc.DBPath)

	if tools.FileExists(dbpath) {
		return dbpath
	} else {
		os.MkdirAll(dbpath, 0755)
	}

	return dbpath
}

func (bc *ChatServerConfig) GetUsersDbPath() string {
	return path.Join(bc.getDbPath(), bc.UsersDBFile)
}

func (bc *ChatServerConfig) GetFriendsDbPath() string {
	return path.Join(bc.getDbPath(), bc.FriendsDBFile)
}

func (bc *ChatServerConfig) GetGroupsDbPath() string {
	return path.Join(bc.getDbPath(), bc.GroupsDBFile)
}

func (bc *ChatServerConfig) GetKeyPath() string {
	return path.Join(GetCSCHomeDir(), bc.KeyFile)
}

func (bc *ChatServerConfig) GetGrpMbrsDbPath() string {
	return path.Join(bc.getDbPath(), bc.GrpMemberDBFile)
}

func IsInitialized() bool {
	if tools.FileExists(GetCSCFGFile()) {
		return true
	}

	return false
}

package api

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatserver/chatcrypt"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/db"
	"github.com/kprc/nbsnetwork/tools"
	"io/ioutil"
	"net/http"
)

type MessageDispatch struct {
}

func (uc *MessageDispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	req := &protocol.UserCommand{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	if !ValidateSig(&req.SP) || req.Op < protocol.AddFriend || req.Op > protocol.QuitGroup {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	var reply *protocol.UCReply

	switch req.Op {
	case protocol.AddFriend:
		reply = AddFriend(req)
	case protocol.DelFriend:
		reply = DelFriend(req)
	case protocol.ChgGroup:
		reply = ChangeGroup(req)
	case protocol.AddGroup:
		reply = AddGroup(req)
	case protocol.DelGroup:
		reply = DelGroup(req)
	case protocol.JoinGroup:
		reply = JoinGroup(req)
	case protocol.QuitGroup:
		reply = QuitGroup(req)
	}

	var bresp []byte

	bresp, err = json.Marshal(*reply)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	w.WriteHeader(200)
	w.Write(bresp)

}

func ValidateSig(sp *protocol.SignPack) bool {
	if !address.ChatAddress(sp.SignText.CPubKey).IsValid() || !address.ChatAddress(sp.SignText.SPubKey).IsValid() {
		return false
	}
	cfg := config.GetCSC()
	if address.ToAddress(cfg.PubKey).String() != sp.SignText.SPubKey {
		return false
	}
	userdb := db.GetChatUserDB()
	u, err := userdb.Find(sp.SignText.CPubKey)
	if err != nil {
		return false
	}

	if u.Alias != sp.SignText.AliasName || u.ExpireTinme != sp.SignText.ExpireTime {
		return false
	}

	if u.ExpireTinme < tools.GetNowMsTime() {
		return false
	}

	var data []byte
	data, err = (&(sp.SignText)).ForSig()
	if err != nil {
		return false
	}

	if !chatcrypt.Verify(cfg.PubKey, data, base58.Decode(sp.Sign)) {
		return false
	}

	return true
}

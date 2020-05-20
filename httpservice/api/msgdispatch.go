package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/kprc/chat-protocol/protocol"
	"encoding/json"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/db"
	"github.com/kprc/chatserver/ed25519"
	"github.com/btcsuite/btcutil/base58"
)

type MessageDispatch struct {

}


func (uc *MessageDispatch)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
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

	req:=&protocol.UserCommand{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	if !ValidateSig(&req.SP) || req.Op<protocol.AddFriend || req.Op>protocol.QuitGroup{
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
	case protocol.AddGroup:
	case protocol.DelGroup:
	case protocol.JoinGroup:
	case protocol.QuitGroup:

	}

	var bresp []byte

	bresp, err = json.Marshal(*reply)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	//log.Println("response:",string(bresp))
	w.WriteHeader(200)
	w.Write(bresp)

}


func ValidateSig(sp *protocol.SignPack) bool {
	if !address.ChatAddress(sp.SignText.CPubKey).IsValid() || !address.ChatAddress(sp.SignText.SPubKey).IsValid(){
		return false
	}
	cfg:=config.GetCSC()
	if address.ToAddress(cfg.PubKey).String() != sp.SignText.SPubKey{
		return false
	}
	userdb:=db.GetChatUserDB()
	u,err:=userdb.Find(sp.SignText.CPubKey)
	if err!=nil{
		return false
	}

	if u.Alias != sp.SignText.AliasName || u.ExpireTinme != sp.SignText.ExpireTime{
		return false
	}

	var data []byte
	data,err = (&(sp.SignText)).ForSig()
	if err!=nil{
		return false
	}

	if !ed25519.Verify(cfg.PubKey,data,base58.Decode(sp.Sign)){
		return false
	}

	return true
}

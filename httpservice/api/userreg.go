package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperorchidlab/chat-protocol/address"
	"github.com/hyperorchidlab/chat-protocol/protocol"
	"github.com/hyperorchidlab/chatserver/chatcrypt"
	"github.com/hyperorchidlab/chatserver/config"
	"github.com/hyperorchidlab/chatserver/db"
)

type UserRegister struct {
}

//const (
//	ChallegeSize int = 16
//)
//
//func newChallege() []byte {
//	challenge := make([]byte, ChallegeSize)
//	for{
//		n, err := io.ReadFull(rand.Reader, challenge)
//		if err != nil || n != ChallegeSize{
//			continue
//		}
//	}
//	return challenge
//}

func (ur *UserRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := &protocol.UserRegReq{}

	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}
	addr := address.ChatAddress(req.CPubKey)
	if !addr.IsValid() {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	resp := &protocol.UserRegResp{}

	userdb := db.GetChatUserDB()

	_, err = userdb.Find(req.CPubKey)
	if err != nil {
		err = userdb.Insert(req.AliasName, req.CPubKey, req.TimeInterval)
	} else {
		err = userdb.Update(req.AliasName, req.CPubKey, req.TimeInterval)
	}

	if err != nil {
		resp.ErrCode = 1
	} else {

		var user *db.ChatUser
		if user, err = userdb.Find(req.CPubKey); err != nil {
			resp.ErrCode = 1
		} else {
			resp.SP.SignText.AliasName = req.AliasName
			resp.SP.SignText.CPubKey = req.CPubKey
			resp.SP.SignText.ExpireTime = user.ExpireTime
			resp.SP.SignText.SPubKey = address.ToAddress(config.GetCSC().PubKey).String()

			signtxt, _ := json.Marshal(resp.SP.SignText)

			resp.SP.Sign = base58.Encode(chatcrypt.Sign(config.GetCSC().PrivKey, signtxt))

		}

	}

	var bresp []byte

	bresp, err = json.Marshal(*resp)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{}")
		return
	}

	//log.Println("response:",string(bresp))
	w.WriteHeader(200)
	w.Write(bresp)

}

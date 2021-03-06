package api

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperorchidlab/chat-protocol/address"
	"github.com/hyperorchidlab/chat-protocol/protocol"
	"github.com/hyperorchidlab/chatserver/db"
	"log"
)

func StoreGroupMsg(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req      *protocol.GroupMsgStoreReq
		err      error
		plaintxt []byte
	)

	cm := &CipherMachine{}

	plaintxt, err = cm.Decrypt(uc)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	req = &protocol.GroupMsgStoreReq{}

	err = json.Unmarshal(plaintxt, &req.GMsg)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	gdb := db.GetGMsgDb()

	gdb.Insert(req.GMsg.Gid, req.GMsg.AesHash, address.ChatAddress(uc.SP.SignText.CPubKey), req.GMsg.Msg)

	return reply

}

func FetchGroupMsg(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req      *protocol.GMsgFetchReq
		err      error
		plaintxt []byte
	)

	cm := &CipherMachine{}

	plaintxt, err = cm.Decrypt(uc)
	if err != nil {
		log.Println(err, 1)
		reply.ResultCode = 1
		return reply
	}

	req = &protocol.GMsgFetchReq{}

	err = json.Unmarshal(plaintxt, &req.GMsg)
	if err != nil {
		log.Println(err, 2)
		reply.ResultCode = 1
		return reply
	}

	log.Println(req.GMsg)

	gdb := db.GetGMsgDb()
	ms := gdb.FindMsg2(req.GMsg.Gid, uc.SP.SignText.CPubKey, req.GMsg.Begin, req.GMsg.Count)

	if len(ms) == 0 {
		log.Println("ms is none")
		reply.ResultCode = 1
		return reply
	}

	resp := &protocol.GMsgFetchResp{}

	resp.GMsg.Gid = req.GMsg.Gid

	for i := 0; i < len(ms); i++ {
		m := ms[i]

		lm := protocol.LGroupMsg{}
		lm.Cnt = m.Cnt
		lm.Msg = m.Msg
		lm.AesHash = m.AesKey
		lm.Speek = m.Speek
		lm.UCnt = m.UCnt

		resp.GMsg.LM = append(resp.GMsg.LM, lm)

	}

	var ciphertxt []byte

	data, _ := json.Marshal(resp.GMsg)

	ciphertxt, err = cm.Encrpt(uc, string(data))
	if err != nil {
		log.Println(err, 3)
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(ciphertxt)

	return reply
}

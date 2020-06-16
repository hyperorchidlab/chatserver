package api

import (
	"encoding/json"
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatserver/db"
	"github.com/mr-tron/base58"
)

func StoreP2pMsg(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req      *protocol.P2pMsgStoreReq
		err      error
		plaintxt []byte
	)

	cm := &CipherMachine{}

	plaintxt, err = cm.Decrypt(uc)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	req = &protocol.P2pMsgStoreReq{}

	err = json.Unmarshal(plaintxt, &req.Msg)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	pdb := db.GetP2PMsgDb()

	pdb.Insert(req.Msg.PeerPk, req.Msg.MyPk, req.Msg.Msg)

	return reply
}

func FetchP2pMs(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req      *protocol.P2pMsgFetchReq
		err      error
		plaintxt []byte
	)

	cm := &CipherMachine{}

	plaintxt, err = cm.Decrypt(uc)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	req = &protocol.P2pMsgFetchReq{}

	err = json.Unmarshal(plaintxt, &req.Msg)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	pdb := db.GetP2PMsgDb()
	ms := pdb.FindMsg(uc.SP.SignText.CPubKey, req.Msg.Begin, req.Msg.Count)

	if len(ms) == 0 {
		reply.ResultCode = 1
		return reply
	}

	resp := &protocol.P2pMsgFetchResp{}

	for i := 0; i < len(ms); i++ {
		m := ms[i]

		lm := protocol.LP2pMsg{}
		lm.Cnt = m.Cnt
		lm.Msg = m.Msg
		lm.PeerPk = m.AesKey

		resp.Msg = append(resp.Msg, lm)

	}

	var ciphertxt []byte

	data, _ := json.Marshal(resp.Msg)

	ciphertxt, err = cm.Encrpt(uc, string(data))
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	uc.CipherTxt = base58.Encode(ciphertxt)

	return reply
}

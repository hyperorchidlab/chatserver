package api

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatserver/chatcrypt"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/db"
)

func ListFriendSummary(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	fdb := db.GetChatFriendsDB()

	var (
		err error
		cf  *db.ChatFriends
	)

	cf, err = fdb.Find(uc.SP.SignText.CPubKey)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	lfsr := &protocol.ListFriendSummaryResp{}

	lfsr.FS.FriendUpdateTime = cf.UpdateTime

	var (
		key, ciphertxt []byte
	)
	cfg := config.GetCSC()
	key, err = chatcrypt.GenerateAesKey(address.ChatAddress(uc.SP.SignText.CPubKey).ToPubKey(), cfg.PrivKey)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	data, _ := json.Marshal(lfsr.FS)
	ciphertxt, err = chatcrypt.Encrypt(key, data)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(ciphertxt)

	return reply

}

func ListGroupSummary(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req *protocol.ListGroupSummaryReq
		err error
	)

	if req, err = DecryptGroupSummaryReq(uc); err != nil {
		reply.ResultCode = 1
		return reply
	}

	resp := &protocol.ListGroupSummaryResp{}

	for i := 0; i < len(req.GS.GroupId); i++ {
		gid := req.GS.GroupId[i]
		gmt := &protocol.GroupModTime{}
		gdb := db.GetChatGroupsDB()

		g, _ := gdb.Find(gid.String())
		if g != nil {
			gmt.GroupUpdate = g.UpdateTime
		}

		gmdb := db.GetChatGrpMbrsDB()
		gm, _ := gmdb.Find(gid.String())
		if gm != nil {
			gmt.GrpMemberUpdate = gm.UpdateTime
		}

		resp.GSR.GMT = append(resp.GSR.GMT, *gmt)
	}

	if len(resp.GSR.GMT) == 0 {
		reply.ResultCode = 1
		return reply
	}

	var (
		key, ciphertxt []byte
	)
	cfg := config.GetCSC()
	key, err = chatcrypt.GenerateAesKey(address.ChatAddress(uc.SP.SignText.CPubKey).ToPubKey(), cfg.PrivKey)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	data, _ := json.Marshal(resp.GSR)

	ciphertxt, err = chatcrypt.Encrypt(key, data)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(ciphertxt)

	return reply

}

func DecryptGroupSummaryReq(uc *protocol.UserCommand) (lgs *protocol.ListGroupSummaryReq, err error) {
	cfg := config.GetCSC()
	cryptbytes := base58.Decode(uc.CipherTxt)

	var (
		key, plainbytes []byte
	)
	key, err = chatcrypt.GenerateAesKey(address.ChatAddress(uc.SP.SignText.CPubKey).ToPubKey(), cfg.PrivKey)
	if err != nil {
		return nil, err
	}

	plainbytes, err = chatcrypt.Decrypt(key, cryptbytes)
	if err != nil {
		return nil, err
	}

	lgs = &protocol.ListGroupSummaryReq{}

	err = json.Unmarshal(plainbytes, &lgs.GS)
	if err != nil {
		return nil, err
	}

	return

}

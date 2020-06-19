package api

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chat-protocol/groupid"
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatserver/chatcrypt"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/db"
	"github.com/kprc/nbsnetwork/tools"
	"log"
)

func JoinGroup(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt

	var (
		req        *protocol.GroupMemberReq
		err        error
		plainBytes []byte
	)

	cm := CipherMachine{}
	plainBytes, err = cm.Decrypt(uc)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	req = &protocol.GroupMemberReq{}
	err = json.Unmarshal(plainBytes, &req.GMD)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	gdb := db.GetChatGroupsDB()
	var g *db.Group
	g, err = gdb.Find(req.GMD.GroupID)
	if err != nil {
		reply.ResultCode = 2
		return reply
	}

	if g.Owner != uc.SP.SignText.CPubKey {
		reply.ResultCode = 3
		return reply
	}

	fdb := db.GetChatFriendsDB()

	_, err = fdb.FindGroup(uc.SP.SignText.CPubKey, req.GMD.GroupID)
	if err != nil {
		reply.ResultCode = 2
		return reply
	}

	_, err = fdb.FindGroup(req.GMD.Friend, req.GMD.GroupID)
	if err == nil {
		reply.ResultCode = 4
		return reply
	}

	err = fdb.AddGroup(req.GMD.Friend, req.GMD.GroupID)
	if err != nil {
		reply.ResultCode = 5
		return reply
	}

	gdb.IncRefer(req.GMD.GroupID)

	gmdb := db.GetChatGrpMbrsDB()

	gmdb.AddMember(req.GMD.GroupID, req.GMD.Friend)

	resp := protocol.GroupMemberResp{}
	gmai := &resp.GMAI

	gmai.GID = groupid.GrpID(req.GMD.GroupID)
	gmai.FriendAddr = address.ChatAddress(req.GMD.Friend)
	udb := db.GetChatUserDB()
	u, _ := udb.Find(req.GMD.Friend)
	if u != nil {
		gmai.FriendName = u.Alias
	}
	gmai.JoinTime = tools.GetNowMsTime()

	f1, _ := fdb.FindFriend(uc.SP.SignText.CPubKey, req.GMD.Friend)
	f2, _ := fdb.FindFriend(req.GMD.Friend, uc.SP.SignText.CPubKey)
	b1 := false
	b2 := false
	if f1 != nil {
		b1 = f1.Agree
	}
	if f2 != nil {
		b2 = f2.Agree
	}

	gmai.Agree = getAgree(b1, b2)

	var j []byte
	var cipherbytes []byte

	j, err = json.Marshal(resp.GMAI)
	if err != nil {
		log.Println(err)
		return reply
	}

	cipherbytes, err = cm.Encrpt(uc, string(j))
	if err != nil {
		log.Println(err)
		return reply
	}

	reply.CipherTxt = base58.Encode(cipherbytes)

	return reply
}

func QuitGroup(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt

	var (
		req *protocol.GroupMemberReq
		err error
	)

	if req, err = DecryptGroupMbrDesc(uc); err != nil {
		reply.ResultCode = 1
		return reply
	}

	gdb := db.GetChatGroupsDB()
	var g *db.Group
	g, err = gdb.Find(req.GMD.GroupID)
	if err != nil {
		reply.ResultCode = 2
		return reply
	}

	if g.Owner != uc.SP.SignText.CPubKey {
		reply.ResultCode = 3
		return reply
	}

	fdb := db.GetChatFriendsDB()

	_, err = fdb.FindGroup(uc.SP.SignText.CPubKey, req.GMD.GroupID)
	if err != nil {
		reply.ResultCode = 2
		return reply
	}

	_, err = fdb.FindGroup(req.GMD.Friend, req.GMD.GroupID)
	if err != nil {
		reply.ResultCode = 4
		return reply
	}

	fdb.DelGroup(req.GMD.Friend, req.GMD.GroupID)
	gdb.DecRefer(req.GMD.GroupID)

	gmdb := db.GetChatGrpMbrsDB()
	gmdb.DelMember(req.GMD.GroupID, req.GMD.Friend)

	return reply

}

func ListGroupMbrs(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req *protocol.ListGroupMbrsReq
		err error
	)
	if req, err = DecryptListGrpMbrs(uc); err != nil {
		reply.ResultCode = 1
		return reply
	}

	//check user in the group id
	fdb := db.GetChatFriendsDB()
	_, err = fdb.FindGroup(uc.SP.SignText.CPubKey, req.LG.GroupId)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}
	gmdb := db.GetChatGrpMbrsDB()
	var gm *db.GroupMember
	gm, err = gmdb.Find(req.LG.GroupId)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	gml := &protocol.GroupMbrDetailsList{}

	for i := 0; i < len(gm.Members); i++ {
		m := gm.Members[i]
		mbr := &protocol.GMember{}
		mbr.PubKey = m
		var u *db.ChatUser
		u, err = db.GetChatUserDB().Find(m)
		if err != nil {
			continue
		}
		mbr.Alias = u.Alias
		mbr.ExpireTime = u.ExpireTime

		var ag, bg bool

		var f *db.Friend
		f, err = db.GetChatFriendsDB().FindFriend(uc.SP.SignText.CPubKey, m)
		if f != nil {
			ag = f.Agree
		}
		f, err = db.GetChatFriendsDB().FindFriend(m, uc.SP.SignText.CPubKey)
		if f != nil {
			bg = f.Agree
		}

		mbr.Agree = getAgree(ag, bg)
		gml.FD = append(gml.FD, *mbr)
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

	data, _ := json.Marshal(*gml)

	ciphertxt, err = chatcrypt.Encrypt(key, data)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(ciphertxt)

	return reply
}

func DecryptListGrpMbrs(uc *protocol.UserCommand) (gd *protocol.ListGroupMbrsReq, err error) {
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

	gd = &protocol.ListGroupMbrsReq{}

	err = json.Unmarshal(plainbytes, &gd.LG)
	if err != nil {
		return nil, err
	}

	return gd, nil
}

func DecryptGroupMbrDesc(uc *protocol.UserCommand) (gd *protocol.GroupMemberReq, err error) {
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

	gd = &protocol.GroupMemberReq{}

	err = json.Unmarshal(plainbytes, &gd.GMD)
	if err != nil {
		return nil, err
	}

	return gd, nil
}

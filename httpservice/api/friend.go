package api

import (
	"github.com/hyperorchidlab/chat-protocol/protocol"

	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperorchidlab/chat-protocol/address"
	"github.com/hyperorchidlab/chatserver/chatcrypt"
	"github.com/hyperorchidlab/chatserver/config"
	"github.com/hyperorchidlab/chatserver/db"
)

func AddFriend(uc *protocol.UserCommand) *protocol.UCReply {

	reply := &protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt

	var (
		req                 *protocol.FriendReq
		err                 error
		plaintxt, cippertxt []byte
	)

	cm := &CipherMachine{}

	plaintxt, err = cm.Decrypt(uc)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	req = &protocol.FriendReq{}
	err = json.Unmarshal(plaintxt, &req.FD)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	udb := db.GetChatUserDB()
	var u *db.ChatUser
	u, err = udb.Find(req.FD.PeerPubKey)
	if err != nil {
		reply.ResultCode = 2
		return reply
	}

	fdb := db.GetChatFriendsDB()
	var agree int

	_, err = fdb.FindFriend(uc.SP.SignText.CPubKey, req.FD.PeerPubKey)
	if err != nil {
		agree = 1
		fdb.AgreeFriend(uc.SP.SignText.CPubKey, req.FD.PeerPubKey, true)
		fdb.AgreeFriend(req.FD.PeerPubKey, uc.SP.SignText.CPubKey, false)
	} else {
		agree = 3
		fdb.AgreeFriend(uc.SP.SignText.CPubKey, req.FD.PeerPubKey, true)
	}

	var f *db.Friend
	f, err = fdb.FindFriend(uc.SP.SignText.CPubKey, req.FD.PeerPubKey)
	if err != nil {
		reply.ResultCode = 3
		return reply
	}

	resp := &protocol.FriendAddResp{}
	resp.FAI.Addr = address.ChatAddress(req.FD.PeerPubKey)
	resp.FAI.AliasName = u.Alias
	resp.FAI.Agree = agree
	resp.FAI.AddTime = f.AddTime

	var j []byte
	j, err = json.Marshal(resp.FAI)
	if err != nil {
		reply.ResultCode = 4
		return reply
	}

	cippertxt, err = cm.Encrpt(uc, string(j))

	reply.CipherTxt = base58.Encode(cippertxt)

	return reply
}

func DelFriend(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt

	var (
		req *protocol.FriendReq
		err error
	)

	if req, err = DecryptFriendDesc(uc); err != nil {
		reply.ResultCode = 1
		return reply
	}

	fdb := db.GetChatFriendsDB()

	_, err = fdb.Find(uc.SP.SignText.CPubKey)
	if err == nil {
		fdb.DelFriend(uc.SP.SignText.CPubKey, req.FD.PeerPubKey)
		fdb.DelFriend(req.FD.PeerPubKey, uc.SP.SignText.CPubKey)
	}

	return reply
}

func getAgree(own bool, friend bool) int {
	if own {
		if friend {
			return 3
		} else {
			return 1
		}
	} else {
		if friend {
			return 2
		} else {
			return 0
		}
	}

}

func ListFriends(uc *protocol.UserCommand) *protocol.UCReply {
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

	if len(cf.Friends) == 0 && len(cf.Groups) == 0 {
		reply.ResultCode = 1
		return reply
	}

	fl := protocol.FriendList{}

	fl.UpdateTime = cf.UpdateTime

	for i := 0; i < len(cf.Friends); i++ {
		f := cf.Friends[i]
		fd := &protocol.FriendDetails{}
		//fd.ExpireTime
		fd.PubKey = f.PubKey
		//fd.Alias = f
		fd.AddTime = f.AddTime
		udb := db.GetChatUserDB()
		var u *db.ChatUser
		u, err = udb.Find(f.PubKey)
		if err == nil {
			fd.ExpireTime = u.ExpireTime
			fd.Alias = u.Alias
		}
		var peerf *db.Friend
		peerf, err = fdb.FindFriend(fd.PubKey, cf.Owner)
		if err == nil {
			fd.Agree = getAgree(f.Agree, peerf.Agree)
		}

		fl.FD = append(fl.FD, *fd)
	}

	for i := 0; i < len(cf.Groups); i++ {
		gid := cf.Groups[i]

		gd := protocol.GroupDetails{}
		gd.GrpId = gid
		gdb := db.GetChatGroupsDB()
		var (
			g *db.Group
		)
		g, err = gdb.Find(gid)
		if err == nil {
			gd.Alias = g.Alias
			gd.CreateTime = g.CreateTime
			gd.IsOwner = false
			if cf.Owner == g.Owner {
				gd.IsOwner = true
			}
		}

		gmbrdb := db.GetChatGrpMbrsDB()
		var mbr *db.GroupMember
		mbr, err = gmbrdb.Find(gid)
		if err == nil {
			gd.MembrsCnt = len(mbr.Members)
		}

		fl.GD = append(fl.GD, gd)

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

	data, _ := json.Marshal(fl)

	ciphertxt, err = chatcrypt.Encrypt(key, data)
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(ciphertxt)

	return reply

}

func DecryptFriendDesc(uc *protocol.UserCommand) (fr *protocol.FriendReq, err error) {

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

	fr = &protocol.FriendReq{}

	err = json.Unmarshal(plainbytes, &fr.FD)
	if err != nil {
		return nil, err
	}

	return fr, nil
}

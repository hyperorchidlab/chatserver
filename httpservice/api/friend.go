package api

import (
	"github.com/kprc/chat-protocol/protocol"

	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chatserver/config"
	"encoding/json"
	"github.com/kprc/chatserver/db"
	"github.com/kprc/chatserver/chatcrypt"
)

func AddFriend(uc *protocol.UserCommand) *protocol.UCReply {

	reply := &protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt

	var (
		req *protocol.FriendReq
		err error
		)

	if req,err = DecryptFriendDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}

	fdb:=db.GetChatFriendsDB()

	_,err = fdb.FindFriend(uc.SP.SignText.CPubKey,req.FD.PeerPubKey)
	if err!=nil{
		fdb.AgreeFriend(uc.SP.SignText.CPubKey,req.FD.PeerPubKey,true)
		fdb.AgreeFriend(req.FD.PeerPubKey,uc.SP.SignText.CPubKey,false)
	}else{
		fdb.AgreeFriend(uc.SP.SignText.CPubKey,req.FD.PeerPubKey,true)
	}

	return reply
}

func DelFriend(uc *protocol.UserCommand) *protocol.UCReply  {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt

	var (
		req *protocol.FriendReq
		err error
	)

	if req,err = DecryptFriendDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}

	fdb:=db.GetChatFriendsDB()

	_,err = fdb.Find(uc.SP.SignText.CPubKey)
	if err==nil{
		fdb.DelFriend(uc.SP.SignText.CPubKey,req.FD.PeerPubKey)
		fdb.DelFriend(req.FD.PeerPubKey,uc.SP.SignText.CPubKey)
	}

	return reply
}


func DecryptFriendDesc(uc *protocol.UserCommand) (fr *protocol.FriendReq,err error) {

	cfg:=config.GetCSC()

	cryptbytes := base58.Decode(uc.CipherTxt)

	var (
		key,plainbytes []byte
	)
	key,err = chatcrypt.GenerateAesKey(address.ChatAddress(uc.SP.SignText.CPubKey).ToPubKey(),cfg.PrivKey)
	if err!=nil{
		return nil,err
	}


	plainbytes,err = chatcrypt.Decrypt(key,cryptbytes)
	if err!=nil{
		return nil,err
	}

	fr = &protocol.FriendReq{}

	err=json.Unmarshal(plainbytes,&fr.FD)
	if err!=nil{
		return nil,err
	}

	return fr,nil
}


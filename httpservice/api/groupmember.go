package api

import (
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatserver/db"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/ed25519"
	"github.com/kprc/chat-protocol/address"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
)

func JoinGroup(uc *protocol.UserCommand) *protocol.UCReply  {
	reply:=&protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt


	var (
		req *protocol.GroupMemberReq
		err error
	)

	if req,err = DecryptGroupMbrDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}

	gdb:=db.GetChatGroupsDB()
	var g *db.Group
	g,err=gdb.Find(req.GMD.GroupID)
	if err!=nil{
		reply.ResultCode = 2
		return reply
	}

	if g.Owner != uc.SP.SignText.CPubKey{
		reply.ResultCode = 3
		return  reply
	}

	fdb:=db.GetChatFriendsDB()

	//var gs string

	_,err = fdb.FindGroup(uc.SP.SignText.CPubKey,req.GMD.GroupID)
	if err!=nil{
		reply.ResultCode = 2
		return reply
	}

	_,err = fdb.FindGroup(req.GMD.Friend,req.GMD.GroupID)
	if err==nil{
		reply.ResultCode = 4
		return reply
	}

	err = fdb.AddGroup(req.GMD.Friend,req.GMD.GroupID)
	if err!=nil{
		reply.ResultCode = 5
		return reply
	}

	gdb.IncRefer(req.GMD.GroupID)

	gmdb:=db.GetChatGrpMbrsDB()

	gmdb.AddMember(req.GMD.GroupID,req.GMD.Friend)

	return reply
}

func QuitGroup(uc *protocol.UserCommand) *protocol.UCReply  {
	reply:=&protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt


	var (
		req *protocol.GroupMemberReq
		err error
	)

	if req,err = DecryptGroupMbrDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}

	gdb:=db.GetChatGroupsDB()
	var g *db.Group
	g,err=gdb.Find(req.GMD.GroupID)
	if err!=nil{
		reply.ResultCode = 2
		return reply
	}

	if g.Owner != uc.SP.SignText.CPubKey{
		reply.ResultCode = 3
		return  reply
	}

	fdb:=db.GetChatFriendsDB()

	_,err = fdb.FindGroup(uc.SP.SignText.CPubKey,req.GMD.GroupID)
	if err!=nil{
		reply.ResultCode = 2
		return reply
	}

	_,err = fdb.FindGroup(req.GMD.Friend,req.GMD.GroupID)
	if err!=nil{
		reply.ResultCode = 4
		return reply
	}

	fdb.DelGroup(req.GMD.Friend,req.GMD.GroupID)
	gdb.DecRefer(req.GMD.GroupID)

	gmdb:=db.GetChatGrpMbrsDB()
	gmdb.DelMember(req.GMD.GroupID,req.GMD.Friend)

	return reply

}




func DecryptGroupMbrDesc(uc *protocol.UserCommand) (gd *protocol.GroupMemberReq,err error) {
	cfg:=config.GetCSC()

	cryptbytes := base58.Decode(uc.CipherTxt)

	var (
		key,plainbytes []byte
	)
	key,err = ed25519.GenerateAesKey(address.ChatAddress(uc.SP.SignText.CPubKey).ToPubKey(),cfg.PrivKey)
	if err!=nil{
		return nil,err
	}


	plainbytes,err = ed25519.Decrypt(key,cryptbytes)
	if err!=nil{
		return nil,err
	}

	gd = &protocol.GroupMemberReq{}

	err=json.Unmarshal(plainbytes,&gd.GMD)
	if err!=nil{
		return nil,err
	}

	return gd,nil
}
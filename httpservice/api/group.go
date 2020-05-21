package api

import (
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chat-protocol/address"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chatserver/db"
	"github.com/kprc/chatserver/chatcrypt"
)

func AddGroup(uc *protocol.UserCommand)  *protocol.UCReply{
	reply:=&protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt


	var (
		req *protocol.GroupReq
		err error
	)

	if req,err = DecryptGroupDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}

	gdb:=db.GetChatGroupsDB()
	//var g *db.Group
	if _,err = gdb.Find(req.GD.GroupID);err!=nil{
		fdb:=db.GetChatFriendsDB()
		err = fdb.AddGroup(uc.SP.SignText.CPubKey,req.GD.GroupID)
		if err!=nil{
			reply.ResultCode = 1
		}else{
			err = gdb.Insert(req.GD.GroupID,req.GD.GroupAlias,uc.SP.SignText.CPubKey)
			if err!=nil{
				fdb.DelGroup(uc.SP.SignText.CPubKey,req.GD.GroupID)
				reply.ResultCode = 1
			}
			gdb.IncRefer(req.GD.GroupID)
			gmdb:=db.GetChatGrpMbrsDB()
			gmdb.Insert(req.GD.GroupID,uc.SP.SignText.CPubKey)
			gmdb.AddMember(req.GD.GroupID,uc.SP.SignText.CPubKey)
		}
	}else{
		reply.ResultCode = 2  //exists
		return reply
	}

	return reply

}

func DelGroup(uc *protocol.UserCommand)  *protocol.UCReply  {
	reply:=&protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt


	var (
		req *protocol.GroupReq
		err error
	)

	if req,err = DecryptGroupDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}

	gdb:=db.GetChatGroupsDB()
	var g *db.Group
	if g,err = gdb.Find(req.GD.GroupID);err!=nil {
		return reply
	}

	if g.RefCnt > 1{
		reply.ResultCode = 1
		return reply
	}
	if g.Owner != uc.SP.SignText.CPubKey{
		reply.ResultCode = 2
		return reply
	}

	fdb:=db.GetChatFriendsDB()
	fdb.DelGroup(uc.SP.SignText.CPubKey,req.GD.GroupID)
	gmdb:=db.GetChatGrpMbrsDB()
	gmdb.DelGroupMember(req.GD.GroupID)
	gdb.DecRefer(req.GD.GroupID)

	return reply

}

func ChangeGroup(uc *protocol.UserCommand) *protocol.UCReply  {
	reply:=&protocol.UCReply{}
	reply.OP = uc.Op
	reply.CipherTxt = uc.CipherTxt


	var (
		req *protocol.GroupReq
		err error
	)

	if req,err = DecryptGroupDesc(uc);err!=nil{
		reply.ResultCode = 1
		return reply
	}
	gdb:=db.GetChatGroupsDB()
	var g *db.Group
	if g,err = gdb.Find(req.GD.GroupID);err!=nil {
		reply.ResultCode = 1
		return reply
	}

	if g.Owner != uc.SP.SignText.CPubKey{
		reply.ResultCode = 2
		return reply
	}

	err = gdb.UpdateAlias(req.GD.GroupID,req.GD.GroupAlias)
	if err!=nil{
		reply.ResultCode = 3
		return reply
	}

	return reply

}



func DecryptGroupDesc(uc *protocol.UserCommand) (gd *protocol.GroupReq,err error) {
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

	gd = &protocol.GroupReq{}

	err=json.Unmarshal(plainbytes,&gd.GD)
	if err!=nil{
		return nil,err
	}

	return gd,nil
}
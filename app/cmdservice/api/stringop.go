package api

import (
	"context"
	"github.com/kprc/chatserver/app/cmdcommon"
	"github.com/kprc/chatserver/app/cmdpb"

	"github.com/kprc/chatserver/chatcrypt"
	"github.com/kprc/chat-protocol/address"

	"github.com/kprc/chatserver/config"
)

type CmdStringOPSrv struct {
}

func (cso *CmdStringOPSrv) StringOpDo(cxt context.Context, so *cmdpb.StringOP) (*cmdpb.DefaultResp, error) {
	msg := ""
	switch so.Op {
	case cmdcommon.CMD_ACCOUNT_CREATE:
		msg = createAccount(so.Param)
	case cmdcommon.CMD_ACCOUNT_LOAD:
		msg = loadAccount(so.Param)
	//case cmdcommon.CMD_DOMAIN:
	//	msg = GetRecords(so.Param)
	//case cmdcommon.CMD_DEAL:
	//	msg = GetDeal(so.Param)
	//case cmdcommon.CMD_ORDER:
	//	msg = GetOrder(so.Param)
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}

func createAccount(passwd string) string  {
	err := chatcrypt.GenEd25519KeyAndSave(passwd)
	if err!=nil{
		return "create account failed"
	}

	chatcrypt.LoadKey(passwd)

	addr:=address.ToAddress(config.GetCSC().PubKey).String()

	return "Address: "+ addr
}


func loadAccount(passwd string) string  {

	chatcrypt.LoadKey(passwd)

	addr:=address.ToAddress(config.GetCSC().PubKey).String()


	return "load account success! \r\nAddress: "+ addr
}

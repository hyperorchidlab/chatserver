package api

import (
	"context"
	"github.com/kprc/chatserver/app/cmdcommon"
	"github.com/kprc/chatserver/app/cmdpb"

	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chatserver/chatcrypt"

	"fmt"
	"github.com/kprc/chat-protocol/groupid"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/db"
	"strconv"
	"time"
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
	case cmdcommon.CMD_LIST_USER:
		msg = showUser(so.Param)
	case cmdcommon.CMD_LIST_GROUP:
		msg = showGroup(so.Param)
	case cmdcommon.CMD_LIST_GRPMBR:
		msg = showGrpMbr(so.Param)
	case cmdcommon.CMD_LIST_FRIEND:
		msg = showFriend(so.Param)
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}

func createAccount(passwd string) string {
	err := chatcrypt.GenEd25519KeyAndSave(passwd)
	if err != nil {
		return "create account failed"
	}

	chatcrypt.LoadKey(passwd)

	addr := address.ToAddress(config.GetCSC().PubKey).String()

	return "Address: " + addr
}

func loadAccount(passwd string) string {

	chatcrypt.LoadKey(passwd)

	addr := address.ToAddress(config.GetCSC().PubKey).String()

	return "load account success! \r\nAddress: " + addr
}

func int64time2string(t int64) string {
	tm := time.Unix(t/1000, 0)
	return tm.Format("2006-01-02/15:04:05")
}

func showUserDetail(u *db.ChatUser) string {
	msg := fmt.Sprintf("%-48s", u.PubKey)
	msg += fmt.Sprintf("%-24s", u.Alias)
	msg += fmt.Sprintf("%-22s", int64time2string(u.CreateTime))
	msg += fmt.Sprintf("%-22s", int64time2string(u.UpdateTime))
	msg += fmt.Sprintf("%-22s", int64time2string(u.ExpireTime))

	return msg
}

func showUser(pk string) string {
	msg := ""

	if pk != "" {
		if !address.ChatAddress(pk).IsValid() {
			msg = "user address not correct"
			return msg
		}

		udb := db.GetChatUserDB()
		if u, err := udb.Find(pk); err != nil {
			msg = "user not found in db"
		} else {
			msg = showUserDetail(u)
		}

		return msg
	}

	udb := db.GetChatUserDB()
	udb.Iterator()

	for {
		_, u, _ := udb.Next()
		if u == nil {
			break
		}
		if msg != "" {
			msg += "\r\n"
		}
		msg += showUserDetail(u)
	}
	if len(msg) == 0 {
		msg = "no user in db"
	}

	return msg

}

func showFriendDetail(f *db.Friend) string {
	msg := fmt.Sprintf("%-48s", f.PubKey)
	msg += fmt.Sprintf("%-21s", int64time2string(f.AddTime))
	msg += fmt.Sprintf("%-8v", f.Agree)

	return msg
}

func showFriendsDetails(f *db.ChatFriends) string {
	msg := fmt.Sprintf("%-48s", f.Owner)
	msg += "\r\nFriends:\r\n"
	for i := 0; i < len(f.Friends); i++ {
		if i > 0 {
			msg += "\r\n"
		}
		msg += showFriendDetail(&(f.Friends[i]))
	}
	msg += "\r\nGroups:\r\n"
	for i := 0; i < len(f.Groups); i++ {
		if i > 0 {
			msg += "\r\n"
		}
		msg += fmt.Sprintf("%-48s", f.Groups[i])
	}

	return msg
}

func showFriend(ownerPk string) string {
	msg := ""
	if ownerPk != "" {
		if !address.ChatAddress(ownerPk).IsValid() {
			msg = "address not correct"
			return msg
		}
		fdb := db.GetChatFriendsDB()
		if f, err := fdb.Find(ownerPk); err != nil {
			msg = "friend not found in db"
		} else {
			msg = showFriendsDetails(f)
		}

		return msg
	}
	fdb := db.GetChatFriendsDB()
	fdb.Iterator()
	for {
		_, f, _ := fdb.Next()
		if f == nil {
			break
		}
		if msg != "" {
			msg += "\r\n"
		}
		msg += showFriendsDetails(f)
	}

	if len(msg) == 0 {
		msg = "no user in db"
	}

	return msg

}

func showGroupDetails(g *db.Group) string {
	msg := fmt.Sprintf("%-48s", g.GrpId)
	msg += fmt.Sprintf("%-24s", g.Alias)
	msg += fmt.Sprintf("%-10s", strconv.Itoa(g.RefCnt))
	msg += fmt.Sprintf("%-22s", int64time2string(g.CreateTime))
	msg += fmt.Sprintf("%-22s", int64time2string(g.UpdateTime))
	msg += "\r\n        "
	msg += fmt.Sprintf("%-48s", g.Owner)

	return msg
}

func showGroup(grpId string) string {
	msg := ""
	if grpId != "" {
		if !groupid.GrpID(grpId).IsValid() {
			msg = "not a correct group id"
			return msg
		}
		gdb := db.GetChatGroupsDB()
		if g, err := gdb.Find(grpId); err != nil {
			msg = "group not found in db"
		} else {
			msg = showGroupDetails(g)
		}
		return msg
	}
	gdb := db.GetChatGroupsDB()
	gdb.Iterator()

	for {
		_, g, _ := gdb.Next()
		if g == nil {
			break
		}
		if msg != "" {
			msg += "\r\n"
		}

		msg += showGroupDetails(g)

	}

	if len(msg) == 0 {
		msg = "no group in db"
	}

	return msg
}

func showGrpMbrDetails(gm *db.GroupMember) string {
	msg := fmt.Sprintf("%-48s", gm.GrpID)
	msg += fmt.Sprintf("%-48s", gm.Owner)
	msg += fmt.Sprintf("%-22s", int64time2string(gm.CreateTime))
	msg += fmt.Sprintf("%-22s", int64time2string(gm.UpdateTime))
	msg += "\r\n members:"
	for i := 0; i < len(gm.Members); i++ {
		msg += "\r\n"
		msg += fmt.Sprintf("%-48s", gm.Members[i])
	}

	return msg
}

func showGrpMbr(grpId string) string {
	msg := ""
	if grpId != "" {
		if !groupid.GrpID(grpId).IsValid() {
			msg = "not a correct group id"
			return msg
		}
		gmdb := db.GetChatGrpMbrsDB()
		if g, err := gmdb.Find(grpId); err != nil {
			msg = "group not found in group member db"
		} else {
			msg = showGrpMbrDetails(g)
		}
		return msg
	}

	gmdb := db.GetChatGrpMbrsDB()
	gmdb.Iterator()

	for {
		_, g, _ := gmdb.Next()
		if g == nil {
			break
		}
		if msg != "" {
			msg += "\r\n"
		}

		msg += showGrpMbrDetails(g)

	}

	if len(msg) == 0 {
		msg = "no group id in group member db"
	}

	return msg

}

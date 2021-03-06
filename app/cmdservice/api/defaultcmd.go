package api

import (
	"context"

	"time"

	"encoding/json"

	"github.com/hyperorchidlab/chatserver/config"

	"github.com/hyperorchidlab/chat-protocol/address"
	"github.com/hyperorchidlab/chatserver/app/cmdcommon"
	"github.com/hyperorchidlab/chatserver/app/cmdpb"
	"github.com/hyperorchidlab/chatserver/httpservice"
	"sync"
)

type CmdDefaultServer struct {
	Stop func()
}

func (cds *CmdDefaultServer) DefaultCmdDo(ctx context.Context,
	request *cmdpb.DefaultRequest) (*cmdpb.DefaultResp, error) {
	if request.Reqid == cmdcommon.CMD_STOP {
		return cds.stop()
	}

	if request.Reqid == cmdcommon.CMD_CONFIG_SHOW {
		return cds.configShow()
	}

	if request.Reqid == cmdcommon.CMD_PK_SHOW {
		return cds.accountShow()
	}

	if request.Reqid == cmdcommon.CMD_RUN {
		return cds.serverRun()
	}

	resp := &cmdpb.DefaultResp{}

	resp.Message = "no cmd found"

	return resp, nil
}

func (cds *CmdDefaultServer) stop() (*cmdpb.DefaultResp, error) {

	go func() {
		time.Sleep(time.Second * 2)
		cds.Stop()
	}()
	resp := &cmdpb.DefaultResp{}
	resp.Message = "server stoped"
	return resp, nil
}

func encapResp(msg string) *cmdpb.DefaultResp {
	resp := &cmdpb.DefaultResp{}
	resp.Message = msg

	return resp
}

func (cds *CmdDefaultServer) configShow() (*cmdpb.DefaultResp, error) {
	cfg := config.GetCSC()

	bapc, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return encapResp("Internal error"), nil
	}

	return encapResp(string(bapc)), nil
}

func (cds *CmdDefaultServer) accountShow() (*cmdpb.DefaultResp, error) {
	cfg := config.GetCSC()

	msg := "please create account"

	if cfg.PubKey != nil {
		msg = address.ToAddress(cfg.PubKey).String()
	}

	return encapResp(msg), nil
}

var (
	runningFlag     bool
	runningFlagLock sync.Mutex
)

func (cds *CmdDefaultServer) serverRun() (*cmdpb.DefaultResp, error) {
	if config.GetCSC().PubKey == nil || config.GetCSC().PrivKey == nil {
		return encapResp("chat server need account"), nil
	}

	if !runningFlag {
		runningFlagLock.Lock()
		defer runningFlagLock.Unlock()

		if !runningFlag {
			go httpservice.StartWebDaemon()
			runningFlag = true
		}
	}

	return encapResp("Server running"), nil
}

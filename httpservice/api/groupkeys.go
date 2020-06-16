package api

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/chat-protocol/address"
	"github.com/kprc/chat-protocol/protocol"
	"github.com/kprc/chatclient/chatcrypt"
	"github.com/kprc/chatserver/config"
	"github.com/kprc/chatserver/db"
)

type CipherMachine struct {
	key []byte
}

func StoreGroupKeys(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}

	reply.OP = uc.Op

	var (
		req *protocol.GroupKeyStoreReq
		err error
	)

	cm := &CipherMachine{}

	if req, err = cm.DecryptStoreGroupKeys(uc); err != nil {
		reply.ResultCode = 1
		return reply
	}

	gksdb := db.GetChatGrpKeysDb()

	gk := gksdb.Insert2(req.GKs.GroupKeys, req.GKs.PubKeys)

	var (
		cipherBytes []byte
	)

	resp := &protocol.GroupKeyStoreResp{}
	resp.GKI.IndexKey = gk

	data, _ := json.Marshal(resp.GKI)

	cipherBytes, err = cm.Encrpt(uc, string(data))
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(cipherBytes)

	return reply
}

func FetchGroupKeys(uc *protocol.UserCommand) *protocol.UCReply {
	reply := &protocol.UCReply{}
	reply.OP = uc.Op

	var (
		req *protocol.GroupKeyFetchReq
		err error
	)

	cm := &CipherMachine{}

	if req, err = cm.DecryptFetchGrpKeyIdx(uc); err != nil {
		reply.ResultCode = 1
		return reply
	}

	gksdb := db.GetChatGrpKeysDb()
	gk := gksdb.Find(req.GKI.IndexKey)
	if gk == nil {
		reply.ResultCode = 1
		return reply
	}

	resp := &protocol.GroupKeyFetchResp{}
	resp.GKs.PubKeys = gk.PubKeys
	resp.GKs.GroupKeys = gk.GroupKeys

	data, _ := json.Marshal(resp.GKs)

	var (
		cipherBytes []byte
	)

	cipherBytes, err = cm.Encrpt(uc, string(data))
	if err != nil {
		reply.ResultCode = 1
		return reply
	}

	reply.CipherTxt = base58.Encode(cipherBytes)

	return reply
}

func (cm *CipherMachine) GenKey(uc *protocol.UserCommand) error {
	if cm.key == nil {
		cfg := config.GetCSC()

		key, err := chatcrypt.GenerateAesKey(address.ChatAddress(uc.SP.SignText.CPubKey).ToPubKey(), cfg.PrivKey)
		if err != nil {
			return err
		}
		cm.key = key
	}

	return nil
}

func (cm *CipherMachine) Encrpt(uc *protocol.UserCommand, data string) (cipherBytes []byte, err error) {
	if err = cm.GenKey(uc); err != nil {
		return nil, err
	}

	return chatcrypt.Encrypt(cm.key, []byte(data))
}

func (cm *CipherMachine) Decrypt(uc *protocol.UserCommand) (plainBytes []byte, err error) {
	if err = cm.GenKey(uc); err != nil {
		return nil, err
	}

	cryptbytes := base58.Decode(uc.CipherTxt)

	plainBytes, err = chatcrypt.Decrypt(cm.key, cryptbytes)
	if err != nil {
		return nil, err
	}

	return
}

func (cm *CipherMachine) DecryptStoreGroupKeys(uc *protocol.UserCommand) (req *protocol.GroupKeyStoreReq, err error) {

	if err = cm.GenKey(uc); err != nil {
		return nil, err
	}

	var (
		plainbytes []byte
	)

	cryptbytes := base58.Decode(uc.CipherTxt)

	plainbytes, err = chatcrypt.Decrypt(cm.key, cryptbytes)
	if err != nil {
		return nil, err
	}

	req = &protocol.GroupKeyStoreReq{}

	err = json.Unmarshal(plainbytes, &req.GKs)
	if err != nil {
		return nil, err
	}

	return

}

func (cm *CipherMachine) DecryptFetchGrpKeyIdx(uc *protocol.UserCommand) (req *protocol.GroupKeyFetchReq, err error) {
	if err = cm.GenKey(uc); err != nil {
		return nil, err
	}

	var (
		plainbytes []byte
	)

	cryptbytes := base58.Decode(uc.CipherTxt)

	plainbytes, err = chatcrypt.Decrypt(cm.key, cryptbytes)
	if err != nil {
		return nil, err
	}

	req = &protocol.GroupKeyFetchReq{}
	err = json.Unmarshal(plainbytes, &req.GKI)
	if err != nil {
		return nil, err
	}

	return

}

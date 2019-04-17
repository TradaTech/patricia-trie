package state

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// DumpAccount ...
type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

// Dump ...
type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

// RawDump ...
func (manager *Manager) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", manager.db.Root()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(manager.db.AccountTrie().NodeIterator(nil))
	for it.Next() {
		addr := manager.db.AccountTrie().GetKey(it.Key)
		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		addressString := new(string)
		rlp.DecodeBytes(addr, addressString)
		account := DumpAccount{
			Balance:  data.Balance.String(),
			Nonce:    data.Nonce,
			Root:     common.Bytes2Hex(data.Root[:]),
			CodeHash: common.Bytes2Hex(data.CodeHash),
			Code:     *manager.GetCode(*addressString),
			Storage:  make(map[string]string),
		}
		storageTrie, _ := manager.db.StorageTrie(*addressString, data.Root)
		storageIt := trie.NewIterator(storageTrie.NodeIterator(nil))
		for storageIt.Next() {
			keyByte := manager.db.AccountTrie().GetKey(storageIt.Key)
			keyString := new(string)
			valueString := new(string)
			rlp.DecodeBytes(keyByte, keyString)
			rlp.DecodeBytes(storageIt.Value, valueString)
			account.Storage[*keyString] = *valueString
		}
		dump.Accounts[*addressString] = account
	}
	return dump
}

// Dump ...
func (manager *Manager) Dump() []byte {
	json, err := json.MarshalIndent(manager.RawDump(), "", "    ")
	if err != nil {
		manager.logger.Error("Cannot dump state", "err", err)
	}

	return json
}

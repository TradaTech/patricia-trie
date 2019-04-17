package database

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

// StorageManager ...
type StorageManager struct {
	trieDB      *trie.Database
	storageTrie *trie.SecureTrie
	trieMutex   sync.RWMutex
}

// Account ...
type Account struct {
	Nonce    uint64      `json:"nonce"`
	Balance  *big.Int    `json:"balance"`
	Root     common.Hash `json:"root"`      // for storage
	CodeHash []byte      `json:"code_hash"` // for code
}

// Update ...
func (manager *StorageManager) Update(key, value []byte) {
	manager.trieMutex.Lock()
	defer manager.trieMutex.Unlock()

	manager.storageTrie.Update(key, value)
	manager.storageTrie.Commit(nil)
	manager.trieDB.Commit(manager.storageTrie.Hash(), true)
}

// Get ...
func (manager *StorageManager) Get(key []byte) []byte {
	manager.trieMutex.RLock()
	defer manager.trieMutex.RUnlock()

	return manager.storageTrie.Get(key)
}

// Delete ...
func (manager *StorageManager) Delete(key []byte) {
	manager.trieMutex.Lock()
	defer manager.trieMutex.Unlock()

	manager.storageTrie.Delete(key)
	manager.storageTrie.Commit(nil)
	manager.trieDB.Commit(manager.storageTrie.Hash(), true)
}

// Root ...
func (manager *StorageManager) Root() common.Hash {
	return manager.storageTrie.Hash()
}

package database

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	log "github.com/inconshreveable/log15"
)

var rootKey = []byte("root")

// Database ...
type Database struct {
	logger       log.Logger
	db           ethdb.Database
	trieDB       *trie.Database
	trieMutex    sync.RWMutex
	accountTrie  *trie.SecureTrie
	accountMutex sync.RWMutex
}

// NewDatabase ...
func NewDatabase() *Database {
	logger := log.New("module", "database")

	db, err := rawdb.NewLevelDBDatabase("./chaindata", 256, 256, "eth/db/chaindata/")
	if err != nil {
		logger.Error("Cannot open level db", "err", err)
	}

	trieDB := trie.NewDatabaseWithCache(db, 256)
	accountDB, err := trie.NewSecure(common.Hash{}, trieDB)
	if err != nil {
		logger.Error("Cannot create trie", "err", err)
	}

	root, err := db.Get(rootKey)
	if root != nil && err == nil {
		accountDB, err = trie.NewSecure(common.BytesToHash(root), trieDB)
	}

	return &Database{
		logger:       logger,
		db:           db,
		trieDB:       trieDB,
		trieMutex:    sync.RWMutex{},
		accountTrie:  accountDB,
		accountMutex: sync.RWMutex{},
	}
}

// Commit ...
func (db *Database) Commit() {
	db.trieMutex.Lock()
	defer db.trieMutex.Unlock()

	root := db.accountTrie.Hash()
	db.trieDB.Commit(root, true)
	if err := db.db.Put(rootKey, root[:]); err != nil {
		db.logger.Error("Cannot put root hash", "err", err)
	}
}

// Reset ...
func (db *Database) Reset() {
	db.trieMutex.Lock()
	defer db.trieMutex.Unlock()

	root, err := db.db.Get(rootKey)
	if err != nil {
		db.logger.Error("Cannot get root", "err", err)
		return
	}
	if accountTrie, err := trie.NewSecure(common.BytesToHash(root), db.trieDB); err == nil {
		db.accountTrie = accountTrie
	} else {
		db.logger.Error("Cannot reset trie", "err", err)
	}
}

// Update ...
func (db *Database) Update(key, value []byte) {
	db.accountMutex.Lock()
	defer db.accountMutex.Unlock()

	db.accountTrie.Update(key, value)
	db.accountTrie.Commit(nil)
}

// Get ...
func (db *Database) Get(key []byte) []byte {
	db.accountMutex.RLock()
	defer db.accountMutex.RUnlock()

	return db.accountTrie.Get(key)
}

// Delete ...
func (db *Database) Delete(key []byte) {
	db.accountMutex.Lock()
	defer db.accountMutex.Unlock()

	db.accountTrie.Delete(key)
	db.accountTrie.Commit(nil)
}

// Root ...
func (db *Database) Root() common.Hash {
	return db.accountTrie.Hash()
}

// TrieDB ...
func (db *Database) TrieDB() *trie.Database { return db.trieDB }

// AccountTrie ...
func (db *Database) AccountTrie() *trie.SecureTrie { return db.accountTrie }

// StorageTrie ...
func (db *Database) StorageTrie(address string, root common.Hash) (*trie.SecureTrie, error) {
	return trie.NewSecure(root, db.trieDB)
}

// StorageManager ...
func (db *Database) StorageManager(address string, root common.Hash) (*StorageManager, error) {
	storageTrie, err := db.StorageTrie(address, root)
	if err != nil {
		return nil, err
	}
	return &StorageManager{db.trieDB, storageTrie, sync.RWMutex{}}, nil
}

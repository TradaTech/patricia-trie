package state

import (
	"math/big"

	"github.com/TradaTech/patricia-trie/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/inconshreveable/log15"
)

// Manager ...
type Manager struct {
	db     *database.Database
	logger log.Logger
}

// Account ...
type Account struct {
	Nonce    uint64      `json:"nonce"`
	Balance  *big.Int    `json:"balance"`
	Root     common.Hash `json:"root"`      // for storage
	CodeHash []byte      `json:"code_hash"` // for code
}

// NewManager ...
func NewManager(db *database.Database) *Manager {
	return &Manager{db, log.New("module", "statemanager")}
}

// Get ...
func (manager *Manager) Get(address string) *Account {
	addressByte, _ := rlp.EncodeToBytes(address)
	stateObjectByte := manager.db.Get(addressByte)
	if stateObjectByte == nil {
		return &Account{
			Nonce:    0,
			Balance:  new(big.Int),
			Root:     common.Hash{},
			CodeHash: []byte{},
		}
	}
	stateObject := new(Account)
	rlp.DecodeBytes(stateObjectByte, stateObject)
	return stateObject
}

// Set ...
func (manager *Manager) Set(address string, stateObject *Account) {
	addressByte, _ := rlp.EncodeToBytes(address)
	stateObjectByte, _ := rlp.EncodeToBytes(stateObject)
	manager.db.Update(addressByte, stateObjectByte)
}

// Root ...
func (manager *Manager) Root() common.Hash {
	return manager.db.Root()
}

// AddBalance ...
func (manager *Manager) AddBalance(address string, amount *big.Int) {
	stateObject := manager.Get(address)
	stateObject.Balance.Add(stateObject.Balance, amount)
	manager.Set(address, stateObject)
}

// SubBalance ...
func (manager *Manager) SubBalance(address string, amount *big.Int) {
	stateObject := manager.Get(address)
	stateObject.Balance.Sub(stateObject.Balance, amount)
	manager.Set(address, stateObject)
}

// SetBalance ...
func (manager *Manager) SetBalance(address string, amount *big.Int) {
	stateObject := manager.Get(address)
	stateObject.Balance = amount
	manager.Set(address, stateObject)
}

// GetBalance ...
func (manager *Manager) GetBalance(address string) *big.Int {
	stateObject := manager.Get(address)
	return stateObject.Balance
}

// SetCode ...
func (manager *Manager) SetCode(address, code string) {
	stateObject := manager.Get(address)
	codeByte, _ := rlp.EncodeToBytes(code)
	codeHash := crypto.Keccak256(codeByte)
	manager.db.TrieDB().InsertBlob(common.BytesToHash(codeHash), codeByte)
	manager.db.TrieDB().Commit(common.BytesToHash(codeHash), true)
	stateObject.CodeHash = codeHash
	manager.Set(address, stateObject)
}

// GetCode ...
func (manager *Manager) GetCode(address string) *string {
	stateObject := manager.Get(address)
	codeByte, err := manager.db.TrieDB().Node(common.BytesToHash(stateObject.CodeHash))
	if err != nil {
		// manager.logger.Error("Cannot get code", "err", err)
		return new(string)
	}
	code := new(string)
	rlp.DecodeBytes(codeByte, code)
	return code
}

// SetAccountRoot ...
func (manager *Manager) SetAccountRoot(address string, root common.Hash) {
	stateObject := manager.Get(address)
	stateObject.Root = root
	manager.Set(address, stateObject)
}

// GetAccountRoot ...
func (manager *Manager) GetAccountRoot(address string) common.Hash {
	stateObject := manager.Get(address)
	return stateObject.Root
}

// GetState ...
func (manager *Manager) GetState(address string, key string) (*string, error) {
	stateObject := manager.Get(address)
	storageManager, err := manager.db.StorageManager(address, stateObject.Root)
	if err != nil {
		return nil, err
	}
	keyByte, _ := rlp.EncodeToBytes(key)
	storageByte := storageManager.Get(keyByte)
	if storageByte == nil {
		return new(string), nil
	}
	result := new(string)
	err = rlp.DecodeBytes(storageByte, result)
	return result, err
}

// SetState ...
func (manager *Manager) SetState(address string, key string, value string) error {
	stateObject := manager.Get(address)
	storageManager, err := manager.db.StorageManager(address, stateObject.Root)
	if err != nil {
		return err
	}
	keyByte, _ := rlp.EncodeToBytes(key)
	valueByte, _ := rlp.EncodeToBytes(value)
	storageManager.Update(keyByte, valueByte)
	root := storageManager.Root()
	manager.SetAccountRoot(address, root)
	return err
}

package ws

import (
	"math/big"

	"github.com/TradaTech/patricia-trie/database"
	"github.com/TradaTech/patricia-trie/state"
	"github.com/ethereum/go-ethereum/common"
)

// Root ...
type Root struct {
	Hash common.Hash `json:"hash"`
}

// StateAPI ...
type StateAPI struct {
	db           *database.Database
	stateManager *state.Manager
}

// NewStateAPI ...
func NewStateAPI(db *database.Database) *StateAPI {
	return &StateAPI{db, state.NewManager(db)}
}

// Account ...
type Account struct {
	Nonce    uint64      `json:"nonce"`
	Balance  *big.Int    `json:"balance"`
	Root     common.Hash `json:"root"`      // for storage
	CodeHash []byte      `json:"code_hash"` // for code
}

// Root return the root of trie
// Example:
// {
//   "jsonrpc": "2.0",
// 	 "method": "state_root",
// 	 "params": [],
// 	 "id": 2
// }
func (node *StateAPI) Root() (*Root, error) {
	root := node.db.Root()
	return &Root{Hash: root}, nil
}

// Commit commit changes in memory to disk
// Example:
// {
//   "jsonrpc": "2.0",
// 	 "method": "state_commit",
// 	 "params": [],
// 	 "id": 2
// }
func (node *StateAPI) Commit() (*Root, error) {
	node.db.Commit()
	root := node.db.Root()
	return &Root{Hash: root}, nil
}

// Reset revert changes to latest commited state
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_reset",
// 	 "params": [],
// 	 "id": 2
// }
func (node *StateAPI) Reset() (*Root, error) {
	node.db.Reset()
	root := node.db.Root()
	return &Root{Hash: root}, nil
}

// Update update state object
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_update",
// 	 "params": ["tea_1", {"balance": 100}],
// 	 "id": 2
// }
func (node *StateAPI) Update(address string, stateObject *state.Account) (*Root, error) {
	node.stateManager.Set(address, stateObject)
	return &Root{Hash: node.stateManager.Root()}, nil
}

// Get get the state object
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_get",
// 	 "params": ["tea_1"],
// 	 "id": 2
// }
func (node *StateAPI) Get(address string) (*state.Account, error) {
	return node.stateManager.Get(address), nil
}

// AddBalance add balance of an address
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_addBalance",
// 	 "params": ["tea_1", 100],
// 	 "id": 2
// }
func (node *StateAPI) AddBalance(address string, amount *big.Int) (*Root, error) {
	node.stateManager.AddBalance(address, amount)
	return &Root{Hash: node.stateManager.Root()}, nil
}

// SubBalance sub balance of an address
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_subBalance",
// 	 "params": ["tea_1", 100],
// 	 "id": 2
// }
func (node *StateAPI) SubBalance(address string, amount *big.Int) (*Root, error) {
	node.stateManager.SubBalance(address, amount)
	return &Root{Hash: node.stateManager.Root()}, nil
}

// SetBalance set balance of an address
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_setBalance",
// 	 "params": ["tea_1", 100],
// 	 "id": 2
// }
func (node *StateAPI) SetBalance(address string, amount *big.Int) (*Root, error) {
	node.stateManager.SetBalance(address, amount)
	return &Root{Hash: node.stateManager.Root()}, nil
}

// GetBalance get balance of an address
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_getBalance",
// 	 "params": ["tea_1"],
// 	 "id": 2
// }
func (node *StateAPI) GetBalance(address string) (*big.Int, error) {
	balance := node.stateManager.GetBalance(address)
	return balance, nil
}

// Code ...
type Code struct {
	Code string `json:"code"`
}

// SetCode set code of an address
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_setCode",
// 	 "params": ["tea_1", "code"],
// 	 "id": 2
// }
func (node *StateAPI) SetCode(address, code string) (*Root, error) {
	node.stateManager.SetCode(address, code)
	return &Root{Hash: node.stateManager.Root()}, nil
}

// GetCode get code of an address
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_getCode",
// 	 "params": ["tea_1"],
// 	 "id": 2
// }
func (node *StateAPI) GetCode(address string) (*Code, error) {
	code := node.stateManager.GetCode(address)
	return &Code{*code}, nil
}

// Storage ...
type Storage struct {
	Value string `json:"value"`
}

// GetState get state of an address (smart contract)
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_getState",
// 	 "params": ["tea_1", "value"],
// 	 "id": 2
// }
func (node *StateAPI) GetState(address, key string) (*Storage, error) {
	storage, err := node.stateManager.GetState(address, key)
	if err != nil {
		return nil, err
	}
	return &Storage{*storage}, nil
}

// SetState ...
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_setState",
// 	 "params": ["tea_1", "value", "1"],
// 	 "id": 2
// }
func (node *StateAPI) SetState(address, key, value string) (*Root, error) {
	err := node.stateManager.SetState(address, key, value)
	if err != nil {
		return nil, err
	}
	return &Root{Hash: node.stateManager.Root()}, nil
}

// Dump dump the state
// Example:
// {
// 	 "jsonrpc": "2.0",
// 	 "method": "state_dump",
// 	 "params": [],
// 	 "id": 2
// }
func (node *StateAPI) Dump() (state.Dump, error) {
	dump := node.stateManager.RawDump()
	return dump, nil
}

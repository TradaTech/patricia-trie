# patricia-trie
IceTea state manager base on patricia trie

## Connection
Using JSON RPC over websocket

## Command
```js
go run main.go
```

## API
#### func (*StateAPI) AddBalance

```go
func (node *StateAPI) AddBalance(address string, amount *big.Int) (*Root, error)
```
AddBalance add balance of an address Example: {

    "jsonrpc": "2.0",
    "method": "state_addBalance",
    "params": ["tea_1", 100],
    "id": 2

}

#### func (*StateAPI) Commit

```go
func (node *StateAPI) Commit() (*Root, error)
```
Commit commit changes in memory to disk Example: {

      "jsonrpc": "2.0",
    	 "method": "state_commit",
    	 "params": [],
    	 "id": 2

}

#### func (*StateAPI) Dump

```go
func (node *StateAPI) Dump() (state.Dump, error)
```
Dump dump the state Example: {

    "jsonrpc": "2.0",
    "method": "state_dump",
    "params": [],
    "id": 2

}

#### func (*StateAPI) Get

```go
func (node *StateAPI) Get(address string) (*state.Account, error)
```
Get get the state object Example: {

    "jsonrpc": "2.0",
    "method": "state_get",
    "params": ["tea_1"],
    "id": 2

}

#### func (*StateAPI) GetBalance

```go
func (node *StateAPI) GetBalance(address string) (*big.Int, error)
```
GetBalance get balance of an address Example: {

    "jsonrpc": "2.0",
    "method": "state_getBalance",
    "params": ["tea_1"],
    "id": 2

}

#### func (*StateAPI) GetCode

```go
func (node *StateAPI) GetCode(address string) (*Code, error)
```
GetCode get code of an address Example: {

    "jsonrpc": "2.0",
    "method": "state_getCode",
    "params": ["tea_1"],
    "id": 2

}

#### func (*StateAPI) GetState

```go
func (node *StateAPI) GetState(address, key string) (*Storage, error)
```
GetState get state of an address (smart contract) Example: {

    "jsonrpc": "2.0",
    "method": "state_getState",
    "params": ["tea_1", "value"],
    "id": 2

}

#### func (*StateAPI) Reset

```go
func (node *StateAPI) Reset() (*Root, error)
```
Reset revert changes to latest commited state Example: {

    "jsonrpc": "2.0",
    "method": "state_reset",
    "params": [],
    "id": 2

}

#### func (*StateAPI) Root

```go
func (node *StateAPI) Root() (*Root, error)
```
Root return the root of trie Example: {

      "jsonrpc": "2.0",
    	 "method": "state_root",
    	 "params": [],
    	 "id": 2

}

#### func (*StateAPI) SetBalance

```go
func (node *StateAPI) SetBalance(address string, amount *big.Int) (*Root, error)
```
SetBalance set balance of an address Example: {

    "jsonrpc": "2.0",
    "method": "state_setBalance",
    "params": ["tea_1", 100],
    "id": 2

}

#### func (*StateAPI) SetCode

```go
func (node *StateAPI) SetCode(address, code string) (*Root, error)
```
SetCode set code of an address Example: {

    "jsonrpc": "2.0",
    "method": "state_setCode",
    "params": ["tea_1", "code"],
    "id": 2

}

#### func (*StateAPI) SetState

```go
func (node *StateAPI) SetState(address, key, value string) (*Root, error)
```
SetState ... Example: {

    "jsonrpc": "2.0",
    "method": "state_setState",
    "params": ["tea_1", "value", "1"],
    "id": 2

}

#### func (*StateAPI) SubBalance

```go
func (node *StateAPI) SubBalance(address string, amount *big.Int) (*Root, error)
```
SubBalance sub balance of an address Example: {

    "jsonrpc": "2.0",
    "method": "state_subBalance",
    "params": ["tea_1", 100],
    "id": 2

}

#### func (*StateAPI) Update

```go
func (node *StateAPI) Update(address string, stateObject *state.Account) (*Root, error)
```
Update update state object Example: {

    "jsonrpc": "2.0",
    "method": "state_update",
    "params": ["tea_1", {"balance": 100}],
    "id": 2

}
package ws

import (
	"fmt"
	"net"

	"github.com/TradaTech/patricia-trie/database"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/inconshreveable/log15"
)

// WS ...
type WS struct {
	handler *rpc.Server
	logger  log.Logger
}

// NewWs ...
func NewWs(db *database.Database) *WS {
	var apis = []rpc.API{
		rpc.API{
			Namespace: "state",
			Version:   "1.0",
			Service:   NewStateAPI(db),
			Public:    true,
		},
	}

	handler := rpc.NewServer()
	logger := log.New("module", "ws")
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			logger.Error("Cannot register api service", "err", err)
			return nil
		}
	}
	return &WS{handler, logger}
}

// Start ...
func (ws *WS) Start(host string, port uint64) {
	var (
		listener net.Listener
		err      error
		endpoint = fmt.Sprintf("%s:%d", host, port)
	)
	if listener, err = net.Listen("tcp", endpoint); err != nil {
		ws.logger.Error("Cannot start the listener", "err", err)
		return
	}
	ws.logger.Info("Start ws successfully", "endpoint", endpoint)
	rpc.NewWSServer([]string{"*"}, ws.handler).Serve(listener)
}

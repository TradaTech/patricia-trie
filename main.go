package main

import (
	"github.com/TradaTech/patricia-trie/database"
	"github.com/TradaTech/patricia-trie/ws"
)

func main() {
	db := database.NewDatabase()
	ws := ws.NewWs(db)
	ws.Start("127.0.0.1", 7000)
}

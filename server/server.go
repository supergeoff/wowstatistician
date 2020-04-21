package server

import "github.com/dgraph-io/badger/v2"

type Context struct {
	Db *badger.DB
}

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dgraph-io/badger/v2"
)

var (
	Db *badger.DB
)

type DefaultController struct {
	beego.Controller
}

func (this *DefaultController) Get() {
	this.Render()
}

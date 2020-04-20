package controllers

import (
	"wowstatistician/helpers/databases"

	"github.com/astaxie/beego"
)

type StatsController struct {
	beego.Controller
}

func (this *StatsController) GetStats() {
	dbname := this.Ctx.Input.Param(":dbname")
	stats, err := databases.GetStatsFromDb(Db, dbname)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	this.Data["json"] = stats
	this.ServeJSON()
}

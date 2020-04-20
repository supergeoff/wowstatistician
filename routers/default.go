package routers

import (
	"wowstatistician/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.DefaultController{})
	beego.Router("/stats/:dbname:string", &controllers.StatsController{}, "get:GetStats")
}

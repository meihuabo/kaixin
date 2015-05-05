package routers

import (
	"github.com/astaxie/beego"
	"kaixin/controllers"
	"kaixin/controllers/admin"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/admin/kefu", &admin.AdminKefuController{})
	beego.Router("/admin/getAccessToken", &admin.AdminKefuController{}, "get:GetAccessToken")
	beego.Router("/admin/addKefu", &admin.AdminKefuController{}, "post:AddKefu")
	beego.Router("/admin/getcallbackip", &admin.AdminKefuController{}, "get:GetCallbackIp")
	beego.Router("/admin/sendCustomMessage", &admin.AdminKefuController{}, "post:SendCustomMessage")
}

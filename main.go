package main

import (
	_ "ark-api/routers"
	_ "ark-api/seeders"
	os "os"

	beego "github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == os.Getenv("RUN_MODE") {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

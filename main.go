package main

import (
	_"github.com/joho/godotenv/autoload"
	_ "ark-api/routers"
	"github.com/astaxie/beego"
	"os"
)


func main() {
	if beego.BConfig.RunMode == os.Getenv("RUN_MODE") {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

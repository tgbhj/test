package main

import (
	"context"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-bindata/go-bindata"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris/v12"
	"log"
	"net/http"
	"time"
)

type Infos struct {
	gorm.Model
	Company string
	Name    string
	Phone   string
	Email   string
	Msg     string
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Failed to connect database")
	}
	defer db.Close()

	if db.HasTable(&Infos{}) {
		db.AutoMigrate(&Infos{})
	} else {
		db.CreateTable(&Infos{})
	}

	//fs := assetfs.AssetFS{
	//	Asset:     bindata.Asset,
	//	AssetDir:  bindata.AssetDir,
	//	AssetInfo: bindata.AssetInfo,
	//}
	http.Handle("/*", http.FileServer(&fs))

	app.RegisterView(iris.HTML("./assets", ".html").Reload(true).Binary(Asset, AssetNames))

	app.Get("/*", func(ctx iris.Context) {
		ctx.View("assets/index.html") // 渲染模板文件： ./assets/index.html
	})

	app.Get("/info", func(ctx iris.Context) {
		var infos []Infos
		db.Order("ID desc").Find(&infos)
		ctx.ReadJSON(&infos)
		ctx.JSON(iris.Map{
			"code": 20000,
			"msg":  "Success",
			"cb":   infos,
		})
	})

	app.Post("/info", func(ctx iris.Context) {
		var infos Infos
		ctx.ReadJSON(&infos)
		db.Create(&infos)
		ctx.JSON(iris.Map{
			"code": 20000,
			"msg":  "Success",
			"cb":   nil,
		})
	})

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		log.Print("\nServer Closing\n")
		app.Shutdown(ctx) // 关闭所有主机
	})

	app.Listen(":9000", iris.WithoutInterruptHandler)
}

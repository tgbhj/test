package main

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris/v12"
	"log"
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
	} else {
		log.Print("Connect database Success")
	}
	defer db.Close()

	if db.HasTable(&Infos{}) {
		db.AutoMigrate(&Infos{})
	} else {
		db.CreateTable(&Infos{})
	}

	app.HandleDir("/", "./build", iris.DirOptions{
		IndexName:  "/index.html",
		Gzip:       false,
		ShowList:   false,
		Asset:      Asset,
		AssetNames: AssetNames,
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
			"cb":   infos,
		})
	})

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		log.Print("\nServer Closing\n")
		app.Shutdown(ctx)
	})

	app.Run(iris.Addr(":9000"), iris.WithoutInterruptHandler)
}

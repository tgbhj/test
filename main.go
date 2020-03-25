package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Info struct {
	Company string        `bson:"company"`
	Name    string        `bson:"name"`
	Phone   string        `bson:"phone"`
	Email   string        `bson:"email"`
	Msg     string        `bson:"msg"`
	Date    time.Time     `bson:"date"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Print("DB Connect Failed：", err)
	} else {
		log.Print("DB Connect Success\n")
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Print("DB Check Failed：", err)
	} else {
		log.Print("DB Check Passed\n")
	}

	collection := client.Database("adinno").Collection("infos")

	app := iris.New()

	app.Get("/info", func(ctx iris.Context) {
		var results []*Info
		cur, err := collection.Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"date": -1}))
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			// 创建一个值，将单个文档解码为该值
			var elem Info
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// 完成后关闭游标
		cur.Close(context.TODO())
		ctx.JSON(iris.Map{
			"code": 20000,
			"cb":   results,
		})
	})

	app.Post("/info", func(ctx iris.Context) {
		info := new(Info)
		info.Date = time.Now()
		err := ctx.ReadJSON(info)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

		insertResult, err := collection.InsertOne(context.TODO(), info)
		if err != nil {
			log.Fatal(err)
		} else {
			ctx.JSON(iris.Map{
				"code": 20000,
				"cb":   insertResult.InsertedID,
			})
		}
	})

	app.HandleDir("/", "./build")

	app.RegisterView(iris.HTML("./build", ".html"))

	app.Get("/*", func(ctx iris.Context) {
		// 渲染模板文件： ./build/index.html
		ctx.View("index.html")
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
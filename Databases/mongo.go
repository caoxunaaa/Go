package Databases

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDb *mongo.Client

func ConnectToDB(uri, name string, timeout time.Duration, num uint64) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(num)
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}
	return client.Database(name), nil
}

func InitMongoDb() {
	var err error
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://172.20.3.12:27017")

	// 连接到MongoDB
	MongoDb, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = MongoDb.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func CloseMongoDb() {
	err := MongoDb.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

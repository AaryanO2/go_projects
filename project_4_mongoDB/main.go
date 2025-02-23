package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/AaryanO2/go_projects/project_4_mongoDB/controller"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()
	uc := controller.NewUserControlller(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe("0.0.0.0:9010", r)

}

func getSession() *mongo.Client {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@mongodb:27017", os.Getenv("mongo_user"), os.Getenv("mongo_password")))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	return client
}

package main

import (
	"context"
	"log"
	"net"
	app "github.com/rahullenkala/activityapp"
	pb "github.com/rahullenkala/activityapp/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("Could not listen on port: %v", err)
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	db := &app.DataBase{
		Client: client,
		DbName: "test",
	}
	app := &app.ActivityApp{
		Db: db,
	}
	pb.RegisterActivityAppServiceServer(s, app)
	log.Println("Server Started")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

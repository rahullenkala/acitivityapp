package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	app "github.com/rahullenkala/activityapp/pkg"
	pb "github.com/rahullenkala/activityapp/proto"

	"google.golang.org/grpc"
)

const (
	dbURL   = "mongodb://localhost:27017"
	address = "127.0.0.1:50052"
)

func main() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Could not listen on port: %v", err)
	}
	//Create  new grpc server instance
	s := grpc.NewServer()
	//Create a new database instance
	db := app.NewDataBase(dbURL, "activityapp")
	//Create Activity app instance
	usrApp := app.NewApp(db)
	//Reister the app
	pb.RegisterActivityAppServiceServer(s, usrApp)
	//Start the server
	go s.Serve(listen)
	log.Println("Server Started")
	fmt.Println("Type exit to stop the server")
	choice := bufio.NewReader(os.Stdin)
	text, _ := choice.ReadString('\n')
	if strings.EqualFold(text, "exit") {
		s.Stop()
	}
}

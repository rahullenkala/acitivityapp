package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	app "github.com/rahullenkala/activityapp"
	pb "github.com/rahullenkala/activityapp/proto"

	"google.golang.org/grpc"
)

const (
	dbURL = "mongodb://localhost:27017"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("Could not listen on port: %v", err)
	}

	s := grpc.NewServer()

	db := app.NewDataBase(dbURL, "test")
	usrApp := app.NewApp(db)
	pb.RegisterActivityAppServiceServer(s, usrApp)
	log.Println("Server Started")
	go s.Serve(listen)

	fmt.Println("Type exit to stop the server")
	choice := bufio.NewReader(os.Stdin)
	text, _ := choice.ReadString('\n')
	if strings.EqualFold(text, "exit") {
		s.Stop()
	}
}

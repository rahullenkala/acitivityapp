/*
This is client implementation for the Activity App,
This implementation supports following methods
1)CreateUser
2)CreateActivity
3)GetUsers


Full documentation of the application can be found at the following link
https://github.com/rahullenkala/activityapp/blob/master/README.md
*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/rahullenkala/activityapp/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	appService := pb.NewActivityAppServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	fmt.Println("Welcome to Activityapp client demo....!!!!")
	fmt.Println("Please select from the following options")
	fmt.Print(" 1:Create User\n 2:Create Activity\n 3:GetUsers\n 4:Exit\n")
	fmt.Print()
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	loop := true
	for loop {

		switch choice {

		case "1\n":
			fmt.Print("\nEnter the name: ")
			reader := bufio.NewReader(os.Stdin)
			name, _ := reader.ReadString('\n')
			name = strings.Trim(name, "\n")

			fmt.Print("Enter the email ID: ")
			emailReader := bufio.NewReader(os.Stdin)
			i, _ := emailReader.ReadString('\n')
			emailID := strings.Trim(i, "\n")

			fmt.Print("Enter the Phone No: ")
			id := bufio.NewReader(os.Stdin)
			phNo, _ := id.ReadString('\n')
			Phone := strings.Trim(phNo, "\n")

			response, err := appService.CreateUser(ctx, &pb.User{
				Name:  name,
				Email: emailID,
				Phone: Phone,
			})
			if err != nil {
				log.Println("Error", err)
			}
			log.Println("User Created", response)

		case "2\n":
			fmt.Print("Enter the Phone NO: ")
			id := bufio.NewReader(os.Stdin)
			phNo, _ := id.ReadString('\n')
			Phone := strings.Trim(phNo, "\n")

			fmt.Print("\nSelect the Activity Type:")
			fmt.Print("1:PLAY \n 2:SLEEP \n 3:EAT \n 4:READ \n")
			reader := bufio.NewReader(os.Stdin)
			optionStr, _ := reader.ReadString('\n')
			option, err := strconv.ParseInt(strings.Trim(optionStr, "\n"), 10, 32)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("\n Enter the duration in Sec: ")
			secReader := bufio.NewReader(os.Stdin)
			secStr, _ := secReader.ReadString('\n')
			duration, parseerr := strconv.ParseUint(strings.Trim(secStr, "\n"), 10, 64)
			if parseerr != nil {
				log.Fatal(err)
			}
			act := &pb.Activity{
				Type:     pb.ActivityType(option),
				Duration: duration,
			}
			
			response, creationerr := appService.CreateActivity(ctx, &pb.CreateActivityRequest{
				Activity: act,
				Phone:    Phone,
			})
			if creationerr != nil {
				log.Println("Error", err)
			}
			log.Println("Activity Created", response)

		case "3\n":
			
			stream, streamErr := appService.GetUsers(ctx, &pb.Empty{})
			if err != nil {
				log.Fatal(streamErr)
			}
			for {
				user, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				log.Println(user)
			}

		case "4\n":
			conn.Close()
			loop = false
			break

		}
		fmt.Println("Enter your choice")
		newReader := bufio.NewReader(os.Stdin)
		choice, _ = newReader.ReadString('\n')

	}

}

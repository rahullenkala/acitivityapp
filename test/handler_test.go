package activityapp

/*
Only few test cases are written.
TODO:ALL TESTCASES need to implemented
Before running these test cases make sure the collection are indexed ans set unique .

Index fields for Acitvity are {Phone,ActivityType}
Index fields for user     are {Phone}
*/
import (
	"context"
	"log"
	"testing"
	"time"

	pb "github.com/rahullenkala/activityapp/proto"
	"google.golang.org/grpc"
)

func Test_CreateUser_EmptyPhone(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewActivityAppServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, usrerr := c.CreateUser(ctx, &pb.User{Name: "Rahul", Phone: "", Email: "mymail@gmail.com"})

	if usrerr == nil {
		t.Errorf("Test failed %v", usrerr)
	}

}

func Test_CreateUser_EmptyName(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewActivityAppServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, usrerr := c.CreateUser(ctx, &pb.User{Name: "", Phone: "1234567890", Email: "mymail@gmail.com"})
	if usrerr == nil {
		t.Errorf("Test failed %v", usrerr)
	}
}

func Test_CreateActivity_EmptyPhone(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewActivityAppServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, acterr := c.CreateActivity(ctx, &pb.CreateActivityRequest{Phone: "", Activity: &pb.Activity{Type: pb.ActivityType_EAT, Duration: 1234}})

	if acterr == nil {
		t.Errorf("Test failed %v", acterr)
	}

}

func Test_CreateUser_DuplicatePhone(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewActivityAppServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, usrerr := c.CreateUser(ctx, &pb.User{Name: "User1", Phone: "1234567890", Email: "mymail@gmail.com"})

	if usrerr != nil {
		t.Errorf("Test failed %v", usrerr)
	}
	_, usrerr1 := c.CreateUser(ctx, &pb.User{Name: "User1", Phone: "1234567890", Email: "mymail@gmail.com"})
	if usrerr1 == nil {
		t.Errorf("Test failed %v", usrerr1)
	}
}

func Test_CreateActivity_DuplicateActivity(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewActivityAppServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, acterr := c.CreateActivity(ctx, &pb.CreateActivityRequest{Phone: "1234567890", Activity: &pb.Activity{Type: pb.ActivityType_EAT, Duration: 1234}})

	if acterr != nil {
		t.Errorf("Test failed %v", acterr)
	}
	_, acterr1 := c.CreateActivity(ctx, &pb.CreateActivityRequest{Phone: "1234567890", Activity: &pb.Activity{Type: pb.ActivityType_EAT, Duration: 1234}})

	if acterr1 == nil {
		t.Errorf("Test failed %v", acterr1)
	}

}

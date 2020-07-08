package activityapp

import (
	"context"
	"errors"
	fmt "fmt"
	"log"
	"time"

	pb "github.com/rahullenkala/activityapp/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//CreateUser creates a user and adds the record to db
func (app *ActivityApp) CreateUser(ctx context.Context, usr *pb.User) (*pb.Response, error) {

	userCollection := app.Db.Client.Database(app.Db.DbName).Collection("user")
	insertResult, err := userCollection.InsertOne(ctx, usr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	res := &pb.Response{

		Message: fmt.Sprintf("%v", insertResult.InsertedID),
	}

	return res, err
}

//CreateActivity adds an acitivity for the user
func (app *ActivityApp) CreateActivity(ctx context.Context, creq *pb.CreateActivityRequest) (*pb.Response, error) {
	userCollection := app.Db.Client.Database(app.Db.DbName).Collection("user")
	activityCollection := app.Db.Client.Database(app.Db.DbName).Collection("activity")
	usrPhone := creq.GetPhone()
	if err := userCollection.FindOne(ctx, bson.M{"phone": usrPhone}).Err(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	var timestmp int64
	if creq.Activity.Timestamp != 0 {
		timestmp = creq.Activity.Timestamp
	} else {
		timestmp = time.Now().UTC().Unix()
	}
	actRecord := ActivityRecord{
		Phone:     usrPhone,
		Timestamp: timestmp,
		Type:      creq.Activity.Type.String(),
		Status:    creq.Activity.Status,
		Duration:  creq.Activity.Duration,
	}
	insertResult, err := activityCollection.InsertOne(ctx, actRecord)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	res := &pb.Response{

		Message: fmt.Sprintf("%v", insertResult.InsertedID),
	}
	return res, err
}

//UpdateActivity updates the attributes of an activity associated with user
func (app *ActivityApp) UpdateActivity(ctx context.Context, uReq *pb.UpdateActivityRequest) (*pb.Response, error) {
	activityCollection := app.Db.Client.Database(app.Db.DbName).Collection("activity")
	reqUnixTime := uReq.Time
	reqTime := time.Unix(reqUnixTime, 0).UTC()
	log.Println(reqTime.Location())
	year, month, day := reqTime.Date()
	log.Println(reqTime.Date())
	strt := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(year, month, day, 23, 59, 59, 0, time.UTC).Unix()

	quer := bson.D{
		{"phone", uReq.Phone},
		{"type", uReq.Activity.Type.String()},
		{"timestamp", bson.D{
			{"$gte", strt},
			{"$lte", end},
		}},
	}
	var update primitive.D
	switch uReq.Parameter {
	case pb.UpdateParam_STATUS:
		update = bson.D{
			{"$set", bson.D{
				{"status", uReq.Activity.Status},
			}},
		}
	case pb.UpdateParam_DURATION:
		update = bson.D{
			{"$set", bson.D{
				{"duration", uReq.Activity.Duration},
			}},
		}
	case pb.UpdateParam_BOTH:
		update = bson.D{
			{"$set", bson.D{
				{"status", uReq.Activity.Status},
				{"duration", uReq.Activity.Duration},
			}},
		}

	}
	log.Println(strt, end, reqTime.Unix())
	updateErr := activityCollection.FindOneAndUpdate(ctx, quer, update).Err()
	if updateErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, updateErr.Error())
	}
	response := &pb.Response{

		Message: "Update Succesfull",
	}
	return response, nil

}

//GetActivityStatus returns the result of executing  isDone,isValid methods of an activity
func (app *ActivityApp) GetActivityStatus(ctx context.Context, sReq *pb.ActivityStatusRequest) (*pb.ActivityStatusResponse, error) {

	activityCollection := app.Db.Client.Database(app.Db.DbName).Collection("activity")

	startTime, endTime := getStartAndEndOfDay(sReq.Time)
	query := bson.D{
		{"phone", sReq.Phone},
		{"type", sReq.Activitytype.String()},
		{"timestamp", bson.D{
			{"$gte", startTime},
			{"$lte", endTime},
		}},
	}
	queryResp := activityCollection.FindOne(ctx, query)
	if queryResp.Err() != nil {
		return nil, status.Errorf(codes.InvalidArgument, queryResp.Err().Error())
	}
	activityObj := getAcitivityObject(queryResp, sReq.Activitytype.String())

	log.Println(activityObj)
	result, err := executeActivityMethod(activityObj, sReq.Method)
	if err != nil {
		return nil, err
	}
	return &pb.ActivityStatusResponse{

		Status: result,
	}, nil
}

//GetUserActivities returns all the activities associated with the user on the specified date
func (app *ActivityApp) GetUserActivities(ctx context.Context, uReq *pb.UserActivityRequest) (*pb.UserActivityResponse, error) {
	activityCollection := app.Db.Client.Database(app.Db.DbName).Collection("activity")

	startTime, endTime := getStartAndEndOfDay(uReq.Time)
	query := bson.D{
		{"phone", uReq.Phone},
		{"timestamp", bson.D{
			{"$gte", startTime},
			{"$lte", endTime},
		}},
	}

	cursor, queryErr := activityCollection.Find(ctx, query)
	if queryErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, queryErr.Error())
	}
	var resp pb.UserActivityResponse

	for cursor.Next(ctx) {
		var record ActivityRecord

		if decodeErr := cursor.Decode(&record); decodeErr != nil {
			log.Fatal("Decoding error", decodeErr)
		}

		activity := &pb.Activity{
			Status:    record.Status,
			Timestamp: record.Timestamp,
			Type:      pb.ActivityType(pb.ActivityType_value[record.Type]),
			Duration:  record.Duration,
		}

		resp.Activities = append(resp.Activities, activity)
	}
	return &resp, nil
}

//GetUsers returns all the users registered with this application
func (app *ActivityApp) GetUsers(msg *pb.Empty, stream pb.ActivityAppService_GetUsersServer) error {
	ctx := context.TODO()
	userCollection := app.Db.Client.Database(app.Db.DbName).Collection("user")
	cursor, err := userCollection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		var record pb.User
		if decodeErr := cursor.Decode(&record); decodeErr != nil {
			log.Fatal(decodeErr)
		}
		if streamerr := stream.Send(&record); streamerr != nil {
			return streamerr
		}
	}
	return nil
}
func getStartAndEndOfDay(reqTime int64) (int64, int64) {
	unixTime := time.Unix(reqTime, 0).UTC()
	year, month, day := unixTime.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(year, month, day, 23, 59, 59, 0, time.UTC).Unix()
	return start, end
}
func executeActivityMethod(act UserActivity, method pb.StatusMethod) (bool, error) {

	switch method {
	case pb.StatusMethod_DONE:
		return act.isDone(), nil
	case pb.StatusMethod_VALID:
		return act.isValid(), nil
	default:
		return false, errors.New("Invalid method")
	}

}
func getAcitivityObject(res *mongo.SingleResult, acType string) UserActivity {
	var activity UserActivity
	var record ActivityRecord
	res.Decode(&record)
	switch acType {
	case "EAT":
		activity = eat{
			Timestamp: record.Timestamp,
			Duration:  record.Duration,
			Status:    record.Status,
		}

	case "SLEEP":
		activity = sleep{
			Timestamp: record.Timestamp,
			Duration:  record.Duration,
			Status:    record.Status,
		}

	case "PLAY":
		activity = play{
			Timestamp: record.Timestamp,
			Duration:  record.Duration,
			Status:    record.Status,
		}

	case "READ":
		activity = read{
			Timestamp: record.Timestamp,
			Duration:  record.Duration,
			Status:    record.Status,
		}

	}

	return activity
}

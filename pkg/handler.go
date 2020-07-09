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

	if inputErr := validPhone(usr.Phone); inputErr != nil {
		return nil, inputErr
	}
	if usr.Name == "" {
		return nil, &RequestError{ErrCode: 200, ErrMessage: "Invalid name"}
	}
	userCollection := app.db.client.Database(app.db.dbName).Collection("user")

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

	if inputErr := validPhone(creq.Phone); inputErr != nil {
		return nil, inputErr
	}

	var timestmp int64
	//create the database collection objects
	userCollection := app.db.client.Database(app.db.dbName).Collection("user")
	activityCollection := app.db.client.Database(app.db.dbName).Collection("activity")

	if err := userCollection.FindOne(ctx, bson.M{"phone": creq.GetPhone()}).Err(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	//Check whether request have a timestamp if not create one
	if creq.Activity.Timestamp != 0 {
		timestmp = creq.Activity.Timestamp
	} else {
		timestmp = time.Now().UTC().Unix()
	}

	//Create an activity record for the given details
	actRecord := ActivityRecord{
		Phone:     creq.GetPhone(),
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
	if inputErr := validPhone(uReq.Phone); inputErr != nil {
		return nil, inputErr
	}
	activityCollection := app.db.client.Database(app.db.dbName).Collection("activity")
	//Generate the Beginning and ending timestamps of the day
	startTime, endTime := getStartAndEndOfDay(uReq.Time)
	quer := bson.D{
		{"phone", uReq.Phone},
		{"type", uReq.Activity.Type.String()},
		{"timestamp", bson.D{
			{"$gte", startTime},
			{"$lte", endTime},
		}},
	}
	var update primitive.D
	//Identify the attribute that need to be modified
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
	if inputErr := validPhone(sReq.Phone); inputErr != nil {
		return nil, inputErr
	}
	activityCollection := app.db.client.Database(app.db.dbName).Collection("activity")

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
	//Create an activity instance from the given details
	activityObj := getActivityObject(queryResp, sReq.Activitytype.String())

	//execute the specified method on the activity
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
	var query primitive.D
	var resp pb.UserActivityResponse

	//We can still enhance the error model by using grpc errdetails package
	if inputErr := validPhone(uReq.Phone); inputErr != nil {
		return nil, inputErr
	}

	activityCollection := app.db.client.Database(app.db.dbName).Collection("activity")

	startTime, endTime := getStartAndEndOfDay(uReq.Time)

	if uReq.Batch {
		query = bson.D{
			{"phone", uReq.Phone},
			{"timestamp", bson.D{
				{"$gte", startTime},
				{"$lte", endTime},
			}},
		}
	} else {
		query = bson.D{
			{"phone", uReq.Phone},
			{"type", uReq.Type.String()},
			{"timestamp", bson.D{
				{"$gte", startTime},
				{"$lte", endTime},
			}},
		}
	}
	cursor, queryErr := activityCollection.Find(ctx, query)
	if queryErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, queryErr.Error())
	}

	//Iterate over all the queried records and generate the response
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
	userCollection := app.db.client.Database(app.db.dbName).Collection("user")
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
func validPhone(phone string) error {
	if phone == "" || len(phone) != 10 {
		return &RequestError{ErrCode: 100, ErrMessage: "Invalid User Phone"}
	}
	return nil
}
func getActivityObject(res *mongo.SingleResult, acType string) UserActivity {
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

# Activity App
Activity app is a golang based application used to track the daily activites of a user.This app supports following functionalities:

 - <a  href="#User"><code>Create a user</code></a>
 - <a  href="#Activity"><code>Create an activity for user</code></a>
 - <a  href="#Update"><code>Update the activity attributes</code></a> 
 - <a  href="#Status"><code>Retrieve the status of an activity</code></a>
 - <a  href="#Activities"><code>Retrieve the activities of a user</code></a>
 - <a  href="#Users"><code>Retrieve all the user</code></a>
 
  
  
This implementation at present supports only four activities namely ***PLAY***,***SLEEP***,***READ**,**EAT***



## Environment
This application uses following tech stack 
Language: [Golang](https://golang.org/doc/install) 
Database:  	&nbsp; [MongoDB](https://godoc.org/go.mongodb.org/mongo-driver/mongo) 
Communication: [gRPC](https://grpc.io/docs/languages/go/quickstart/)  
Message Formats: [ProtoBuffer](https://developers.google.com/protocol-buffers)

 
## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to install Go on your computer. The application requires **Go 1.13 or above** and **mongoDB v4.2.1**
Make sure  you [install](https://docs.mongodb.com/manual/installation/) and run mongo service before continuing.
After installing go  run the following commands to start using this application

	# clone the repo from github  
    git clone github.com/rahullenkala/activityapp
  
    cd activityapp
    
    # navigate to server repo 
    cd activityapp/cmd/server
    
    # run the server...!
    go run server.go
   
You have successfully started gRPC server listening on **"127.0.0.1:50052"**

 Let  us start the client program.
 
Note:-Client implements only few methods of the activity app service.

    # navigate to the client repo
    cd activityapp/cmd/server
    
    # run the client...!!!
    go run client.go

## API Doc
Now let us take a look at each rpc call and their usage, supported by this application.

<a  name="User"></a>
### CreateUser(User) returns (Response)
This call is used to create a new user in the application, this call returns an error if an invalid user-name or phone-number are provided.
Following are the definitions of User and Response 

    message User {
	    string name=1;
	    string email=2;
	    string phone=3; #Phone is used as a primary key in db
	 }
	 
	message Response{
	    string message=1; 
    }
    
    #Example 
    response, err := appService.CreateUser(ctx, &pb.User{
				Name:  name,
				Email: emailID,
				Phone: Phone,
			})

<a name="Activity"></a>
### CreateActivity(CreateActivityRequest) returns (Response)
This call is used to create and assign an activity to a user, this call return an error in the following cases.

 - Invalid User
 - Redundant Activity on the same day 

Activites can be created in advance by specifying the timestamp of the date.
Following are the definitions of  Activity and CreateActivityRequest
    
    
    message Activity {
	    ActivityType type=1;
	    int64        timestamp=2;
	    bool         status=3;
	    uint64       duration=4;
    }
    
    message CreateActivityRequest{
	    Activity Activity=1;
	    string   Phone=2;   
    }
    
    #Example   
    act := &pb.Activity{
				Type:     pb.ActivityType(option),
				Duration: duration,
			}
			
    response, creationerr := appService.CreateActivity(ctx, &pb.CreateActivityRequest{
				Activity: act,
				Phone:    Phone,
			})

<a  name="Update"></a>
### UpdateActivity(UpdateActivityRequest) returns (Response)
This call is used to update the attributes of an activity,users can only update the status and the duration of the activity.
Following are  the definitions of UpdateActivityRequest

    message UpdateActivityRequest{
	    Activity Activity=1;
	    string Phone=2;
	    int64  Time=3; #Timestamp to identify the activity based on the date 
	    UpdateParam Parameter=4;
    }
  
    #Example 
    response,err:=appService.UpdateActivity(ctx,&pb.UpdateActivityRequest{
				Phone:"1234567890",
				Time:12345678,
				Parameter: pb.UpdateParam_STATUS,	
			})
    
<a  name="Status"></a>
### GetActivityStatus(ActivityStatusRequest) returns (ActivityStatusResponse)
This call is used to retrieve the status of the activity by calling two in-built methods.

 - isDone()
 - isValid()
 
Following are the definitions for ActivityStatus Request/Response.

    message ActivityStatusRequest{
	    string       Phone=1;
	    ActivityType Activitytype=2;
	    StatusMethod Method=3;
	    int64        Time=4;
    }
    message ActivityStatusResponse{
	    bool  status=1; 
    }
    #Example
    request,err:=appService.GetActivityStatus(ctx,&pb.ActivityStatusRequest{
				Phone:"",
				Method:pb.StatusMethod_DONE,
				Activitytype: pb.ActivityType_EAT,
			})

<a  name="Activities"></a>
### GetUserActivities(UserActivityRequest) returns (UserActivityResponse)
This call is used to retrieve activity/activities of a user on a particular date.
Following are the definitions of UserActivity Request/Response

    message UserActivityRequest{

	    string Phone=1;
	    int64  Time=2;
	    ActivityType Type=3;
	    bool     batch=4;
    }
    
    message UserActivityResponse{   
	    repeated Activity activities=1;
    }
    #Example
    response,err:=appService.GetUserActivities(ctx,&pb.UserActivityRequest{
				Phone: "0987654321",
				Time: 1234567,
				Type: pb.ActivityType_EAT,
				Batch: false,
			})
			

<a  name="Users"></a>
### GetUsers(Empty) returns (stream  User)
This call is used to retrieve all the users registered with this application,this call returns a uni-directional gRPC stream. 
Data can be retrieved from this stream.

    stream, streamErr := appService.GetUsers(ctx, &pb.Empty{})
			if err != nil {
				log.Fatal(streamErr)
			}

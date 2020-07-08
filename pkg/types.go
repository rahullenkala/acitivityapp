package activityapp

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

//ActivityApp implements all the methods defined by  acitivity app service
type ActivityApp struct {
	db *DataBase
}

//DataBase defines basic attributes and instance of the db used by Acitivity app
type DataBase struct {
	client *mongo.Client
	dbName string
}

//UserData represents a user record in the database
type UserData struct {
	Name  string
	Email string
	Phone string
}

//ActivityRecord represents a activity record in the database
type ActivityRecord struct {
	Phone     string
	Status    bool
	Type      string
	Duration  uint64
	Timestamp int64
}

//UserActivity defines an iterfaces for various activites supported by this application
type UserActivity interface {
	isValid() bool
	isDone() bool
}

type play struct {
	Timestamp int64
	Duration  uint64
	Status    bool
}

func (p play) isValid() bool {
	playHours := float32(p.Duration / 3600)
	if playHours >= 1 && playHours <= 3 {
		return true
	}
	return false
}
func (p play) isDone() bool {
	return p.Status
}

type eat struct {
	Timestamp int64
	Duration  uint64
	Status    bool
}

func (e eat) isDone() bool {
	return e.Status
}
func (e eat) isValid() bool {
	eatHours := float32(e.Duration / 3600)
	if eatHours >= 0.25 && eatHours <= 1 {
		return true
	}
	return false
}

type sleep struct {
	Timestamp int64
	Duration  uint64
	Status    bool
}

func (s sleep) isDone() bool {
	return s.Status
}
func (s sleep) isValid() bool {
	sleepHours := float32(s.Duration / 3600)
	if sleepHours >= 6 && sleepHours <= 8 {
		return true
	}
	return false
}

type read struct {
	Timestamp int64
	Duration  uint64
	Status    bool
}

func (r read) isDone() bool {
	return r.Status
}
func (r read) isValid() bool {
	readHours := float32(r.Duration / 3600)
	if readHours >= 0.5 && readHours <= 8 {
		return true
	}
	return false
}

//NewDataBase return a new db instance
func NewDataBase(dbURL string, dbName string) *DataBase {

	clientOptions := options.Client().ApplyURI(dbURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	_, indexerr := client.Database(dbName).Collection("user").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"phone", bsonx.Int32(1)}, {"type", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)
	if indexerr != nil {
		log.Fatal(indexerr)
	}
	db := &DataBase{
		client: client,
		dbName: dbName,
	}
	return db
}

//NewApp return a new Activityapp instance
func NewApp(db *DataBase) *ActivityApp {
	return &ActivityApp{
		db: db,
	}
}

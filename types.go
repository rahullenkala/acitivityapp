package activityapp

import "go.mongodb.org/mongo-driver/mongo"

//ActivityApp ...
type ActivityApp struct {
	Db *DataBase
}

//DataBase ...
type DataBase struct {
	Client *mongo.Client
	DbName string
}

type UserData struct {
	Name  string
	Email string
	Phone string
}

type ActivityRecord struct {
	Phone     string
	Status    bool
	Type      string
	Duration  uint64
	Timestamp int64
}

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

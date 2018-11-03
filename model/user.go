package model

import "github.com/globalsign/mgo/bson"

type User struct {
	Id bson.ObjectId
	Username string
	Password string
	Role string
}

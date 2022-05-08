package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	Openid    string        `json:"openid" bson:"openid"`
	Name      string        `json:"name" bson:"name"`
	Pass      string        `json:"-" bson:"pass"`
	Tel       string        `json:"tel" bson:"tel"`
	Email     string        `json:"email" bson:"email"`
	Gender    int8          `json:"gender" bson:"gender"`
	Created   int64         `json:"created" bson:"created"`
	Last      int64         `json:"last" bson:"last"`
	ShareDocs []string      `json:"-" bson:"share_docs"`
}

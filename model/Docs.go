package model

import "gopkg.in/mgo.v2/bson"

type Docs struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	Title       string        `json:"title" bson:"title"`
	Content     string        `json:"-" bson:"content"`
	Creator     string        `json:"creator" bson:"creator"`
	Openid      string        `json:"openid" bson:"-"`
	CreatorName string        `json:"creator_name" bson:"-"`
	Created     int64         `json:"created" bson:"created"`
	Last        int64         `json:"last" bson:"last"`
	Share       []string      `json:"-" bson:"share"`
}

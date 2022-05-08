package db

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongo *mgo.Session

type mongoDB struct {
	DB *mgo.Database
}

var MongoDB mongoDB
var TestCollection *mgo.Collection

func (t *mongoDB) Open() error {
	var err error
	mongo, err = mgo.Dial("")
	if err != nil {
		return err
	}
	MongoDB.DB = mongo.DB("MetaDocs")
	UsersCollection = MongoDB.DB.C("users")
	DocsCollection = MongoDB.DB.C("docs")
	TestCollection = MongoDB.DB.C("test")
	return err
}

func (t *mongoDB) Close() {
	mongo.Close()
}

type Person struct {
	Name  string
	Phone string
}

func (t *mongoDB) Test() {
	mongo.SetMode(mgo.Monotonic, true)

	c := mongo.DB("test").C("people")
	err := c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}

func TestInsert(d any) error {
	c := mongo.DB("test").C("test")
	err := c.Insert(&d)
	return err
}

func TestFinds() []any {
	c := mongo.DB("test").C("test")
	var d any
	var datas []any
	iter := c.Find(nil).Skip(0).Limit(100).Iter()
	for iter.Next(&d) {
		datas = append(datas, d)
	}
	return datas
}

func TestFind(id string) (any, error) {
	c := mongo.DB("test").C("test")
	var data any
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(&data)
	return data, err
}

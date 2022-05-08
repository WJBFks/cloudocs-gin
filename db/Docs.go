package db

import (
	"Cloudocs/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type docs struct{}

var Docs docs

var DocsCollection *mgo.Collection

func (t *docs) Insert(Docs model.Docs) error {
	err := DocsCollection.Insert(&Docs)
	return err
}

func (t *docs) FindId(id string) (model.Docs, error) {
	Docs := model.Docs{}
	objectId := bson.ObjectIdHex(id)
	err := DocsCollection.FindId(objectId).One(&Docs)
	return Docs, err
}

func (t *docs) UpdateTitle(id string, title string) error {
	objectIdnew := bson.NewObjectId()
	objectId := bson.ObjectIdHex(id)
	err := DocsCollection.Update(
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{"title": title, "last": objectIdnew.Time().Unix()}},
	)
	return err
}

func (t *docs) UpdateContent(id string, content string) error {
	objectIdnew := bson.NewObjectId()
	objectId := bson.ObjectIdHex(id)
	err := DocsCollection.Update(
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{"content": content, "last": objectIdnew.Time().Unix()}},
	)
	return err
}

func (t *docs) Find(key string, value string) (model.Docs, error) {
	Docs := model.Docs{}
	err := DocsCollection.Find(bson.M{key: value}).One(&Docs)
	return Docs, err
}

func (t *docs) Finds(value string, skip, limit int) []model.Docs {
	doc := model.Docs{}
	var Docs []model.Docs
	iter := DocsCollection.Find(bson.M{"creator": value}).Skip(skip).Limit(limit).Iter()
	for iter.Next(&doc) {
		Docs = append(Docs, doc)
	}
	return Docs
}

func (t *docs) Delete(id string) error {
	objectId := bson.ObjectIdHex(id)
	err := DocsCollection.Remove(bson.M{"_id": objectId})
	return err
}

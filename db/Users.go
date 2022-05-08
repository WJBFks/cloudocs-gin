package db

import (
	"Cloudocs/model"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type users struct{}

var Users users

var UsersCollection *mgo.Collection

func (t *users) Insert(User model.User) error {
	err := UsersCollection.Insert(&User)
	return err
}

func (t *users) Update(User model.User) error {
	err := UsersCollection.UpdateId(User.Id, User)
	return err
}

func (t *users) FindId(id string) (model.User, error) {
	User := model.User{}
	objectId := bson.ObjectIdHex(id)
	err := UsersCollection.FindId(objectId).One(&User)
	return User, err
}

func (t *users) Find(key string, value string) (model.User, error) {
	User := model.User{}
	err := UsersCollection.Find(bson.M{key: value}).One(&User)
	return User, err
}

func (t *users) FindAll(skip int, limit int) []model.User {
	User := model.User{}
	var Users []model.User
	iter := UsersCollection.Find(nil).Skip(skip).Limit(limit).Iter()
	for iter.Next(&User) {
		Users = append(Users, User)
	}
	return Users
}

func (t *users) UpdateLast(id string) error {
	objectIdnew := bson.NewObjectId()
	objectId := bson.ObjectIdHex(id)
	err := UsersCollection.Update(
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{"last": objectIdnew.Time().Unix()}},
	)
	return err
}

func (t *users) AddShareDocs(objectId bson.ObjectId, doc string) error {
	User := model.User{}
	Doc := model.Docs{}
	docObjectId := bson.ObjectIdHex(doc)
	err := UsersCollection.FindId(objectId).One(&User)
	err = DocsCollection.FindId(docObjectId).One(&Doc)
	if User.ShareDocs == nil {
		User.ShareDocs = make([]string, 0)
		User.ShareDocs = append(User.ShareDocs, doc)
	} else {
		for _, value := range User.ShareDocs {
			if value == doc {
				return nil
			}
		}
		User.ShareDocs = append(User.ShareDocs, doc)
	}
	err = UsersCollection.Update(
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{"share_docs": User.ShareDocs}},
	)

	if Doc.Share == nil {
		Doc.Share = make([]string, 0)
		Doc.Share = append(Doc.Share, User.Tel)
	} else {
		Doc.Share = append(Doc.Share, User.Tel)
	}
	err = DocsCollection.Update(
		bson.M{"_id": docObjectId},
		bson.M{"$set": bson.M{"share": Doc.Share}},
	)
	return err
}

func (t *users) DelShareDocs(doc string, tel string) error {
	Doc := model.Docs{}
	docObjectId := bson.ObjectIdHex(doc)
	User, err := Users.Find("tel", tel)
	err = DocsCollection.FindId(docObjectId).One(&Doc)
	ShareDocs := make([]string, 0)
	if User.ShareDocs != nil {
		for _, value := range User.ShareDocs {
			if value != doc {
				ShareDocs = append(ShareDocs, value)
			}
		}
	}
	err = UsersCollection.Update(
		bson.M{"_id": User.Id},
		bson.M{"$set": bson.M{"share_docs": ShareDocs}},
	)
	Share := make([]string, 0)
	if Doc.Share != nil {
		for _, value := range Doc.Share {
			if value != tel {
				Share = append(Share, value)
			}
		}
	}
	err = DocsCollection.Update(
		bson.M{"_id": docObjectId},
		bson.M{"$set": bson.M{"share": Share}},
	)
	return err
}

func (t *users) UpdateShareDocs(objectId bson.ObjectId, docs []string) error {
	err := UsersCollection.Update(
		bson.M{"_id": objectId},
		bson.M{"$set": bson.M{"share_docs": docs}},
	)
	return err
}

func (t *users) Test() {
	User := model.User{
		Name:  "WJBFks",
		Pass:  "************",
		Tel:   "12300000000",
		Email: "123456@123.com",
	}
	err := Users.Insert(User)
	if err != nil {
		return
	}
	uis := Users.FindAll(0, 10)
	fmt.Println(uis)
}

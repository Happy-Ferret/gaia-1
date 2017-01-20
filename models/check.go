package models

import (
	"github.com/notyim/gaia/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"strings"
	//"time"
)

type Check struct {
	// identification information
	ID   bson.ObjectId `bson:"_id"`
	URI  string        `bson:"uri"`
	Type string        `bson:"type"`
}

func NewCheck(uri, checkType string) *Check {
	return &Check{
		ID:   bson.NewObjectId(),
		URI:  uri,
		Type: checkType,
	}
}

func (s *Check) FindByID(id bson.ObjectId, db *mgo.Database) error {
	return s.coll(db).FindId(id).One(s)
}

func (*Check) coll(db *mgo.Database) *mgo.Collection {
	return db.C("service")
}

type Checks []Check

func (s *Checks) All() error {
	return mongo.Query(func(session *mgo.Database) error {
		session.C("checks").Find(nil).Sort("_id").All(s)

		return nil
	})
}

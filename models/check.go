package models

import (
	"github.com/notyim/gaia/db/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
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

type Checks *[]Check

//func (c *Checks) All() error {
//	return mongo.Query(func(session *mgo.Database) error {
//		var checks []Check
//		session.C("checks").Find(nil).Sort("_id").All(&checks)
//		c = checks
//		log.Println(c)
//		return nil
//	})
//}

func AllChecks(c *[]Check) error {
	return mongo.Query(func(session *mgo.Database) error {
		session.C("checks").Find(nil).Sort("_id").All(c)
		return nil
	})
}

func FindChecksAfter(c *[]Check, id bson.ObjectId) error {
	return mongo.Query(func(session *mgo.Database) error {
		session.C("checks").Find(bson.M{"_id": bson.M{"$gt": id}}).Sort("_id").All(c)
		return nil
	})
}
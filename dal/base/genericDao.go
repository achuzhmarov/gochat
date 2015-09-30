package base

import (
	"gochat/model"

	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var GenericDao = sessionWrapper{initSession()}
/*var DB_HOST = os.Getenv("MONGO_PORT_27017_TCP_ADDR")
var DB_PORT = os.Getenv("MONGO_PORT_27017_TCP_PORT")
var DB_URL = DB_HOST + ":" + DB_PORT*/
var DB_URL = os.Getenv("MONGO_URL")

const DB_NAME = "gochat"

type sessionWrapper struct {
	Session *mgo.Session
}

func initSession() *mgo.Session {
	session, err := mgo.Dial(DB_URL)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return session
}

func (s *sessionWrapper) Upsert(
	model model.BaseModel, collectionName string) error {

	session := s.Session.Copy()
	defer session.Close()

	collection := session.DB(DB_NAME).C(collectionName)

	if model.GetId().Valid() {
		return collection.UpdateId(model.GetId(), model)
	} else {
		model.SetId(bson.NewObjectId())
		return collection.Insert(model)
	}
}

func (s *sessionWrapper) Update(collectionName string, selector interface{}, update interface{}) error {
	session := s.Session.Copy()
	defer session.Close()

	collection := session.DB(DB_NAME).C(collectionName)

	return collection.Update(selector, update)
}

func (s *sessionWrapper) Delete(
	model model.BaseModel, collectionName string) error {

	session := s.Session.Copy()
	defer session.Close()

	collection := session.DB(DB_NAME).C(collectionName)

	return collection.RemoveId(model.GetId())
}

func (s *sessionWrapper) SearchCollection(
	collection string,
	search func(*mgo.Collection) error) error {

	session := s.Session.Copy()
	defer session.Close()
	c := session.DB(DB_NAME).C(collection)
	return search(c)
}

func (s *sessionWrapper) ClearCollection(collection string) error {
	session := s.Session.Copy()
	defer session.Close()
	c := session.DB(DB_NAME).C(collection)
	return c.DropCollection()
}

func (s *sessionWrapper) GetAllCollection(collection string) {

}

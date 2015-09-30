package model

import (
	"gopkg.in/mgo.v2/bson"
)

type BaseModel interface {
	GetId() bson.ObjectId
	SetId(bson.ObjectId)
}

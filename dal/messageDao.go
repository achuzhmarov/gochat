package dal

import (
	"gochat/dal/base"
	"gochat/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var MessageDaoInstance = MessageDao{"messages"}

type MessageDao struct {
	ColName string
}

func (m *MessageDao) Upsert(message *model.Message) (*model.Message, error) {
	err := base.GenericDao.Upsert(message, m.ColName)
	return message, err
}

func (m *MessageDao) Delete(message *model.Message) error {
	return base.GenericDao.Delete(message, m.ColName)
}

func (m *MessageDao) GetById(messageId bson.ObjectId) (*model.Message, error) {
	query := bson.M{"_id": messageId}
	messages, err := m.search(query, 0, 0)
	if err != nil {
		return nil, err
	} else {
		return (*messages)[0], err
	}
}

func (m *MessageDao) GetAllByTalk(talkId bson.ObjectId) (*model.Messages, error) {
	query := bson.M{"_talk_id": talkId}
	return m.search(query, 0, 0)
}

func (m *MessageDao) GetAll() (*model.Messages, error) {
	query := bson.M{}
	return m.search(query, 0, 0)
}

func (m *MessageDao) search(
	q interface{}, skip int, limit int) (*model.Messages, error) {

	results := &model.Messages{}

	query := func(c *mgo.Collection) error {
		if limit > 0 {
			return c.Find(q).Skip(skip).Limit(limit).All(results)
		} else {
			return c.Find(q).Skip(skip).All(results)
		}
	}

	err := base.GenericDao.SearchCollection(m.ColName, query)

	return results, err
}

func (m *MessageDao) ClearCollection() {
	base.GenericDao.ClearCollection(m.ColName)
}

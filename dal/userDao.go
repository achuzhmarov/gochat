package dal

import (
	"gochat/dal/base"
	"gochat/model"

	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
)

var UserDaoInstance = UserDao{"users"}

type UserDao struct {
	ColName string
}

type UserNotFoundError struct {
	UserLogin string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user %s not found", e.UserLogin)
}

type UserTalkNotFoundError struct {
	UserLogin string
	TalkId    bson.ObjectId
}

func (e *UserTalkNotFoundError) Error() string {
	return fmt.Sprintf("user talk %s not found", e.UserLogin+"-"+string(e.TalkId))
}

type MoreThanOneUserError struct {
	UserLogin string
}

func (e *MoreThanOneUserError) Error() string {
	return fmt.Sprintf("more than one user exists for login %s", e.UserLogin)
}

func (u *UserDao) Upsert(user *model.User) (*model.User, error) {
	err := base.GenericDao.Upsert(user, u.ColName)
	return user, err
}

func (u *UserDao) Delete(message *model.User) error {
	return base.GenericDao.Delete(message, u.ColName)
}

func (u *UserDao) GetByLogin(login string) (*model.User, error) {
	query := bson.M{"login": login}
	users, err := u.search(query, nil, 0, 0)
	if err != nil {
		return nil, err
	} else if len(*users) == 1 {
		return (*users)[0], err
	} else if len(*users) == 0 {
		return nil, &UserNotFoundError{login}
	} else {
		return nil, &MoreThanOneUserError{login}
	}
}

func (u *UserDao) GetAll() (*model.Users, error) {
	query := bson.M{}
	return u.search(query, nil, 0, 0)
}

func (u *UserDao) GetAllPotentialFriends(likeLogin string) (*model.Users, error) {
	var query = bson.M{}

	if likeLogin != "" {
		query = bson.M{"login": bson.RegEx{".*" + likeLogin + ".*", ""}}
	}

	selector := bson.M{"login": 1}
	return u.search(query, selector, 0, 0)
}

func (u *UserDao) AddFriend(userLogin string, friend string) error {
	query := bson.M{"login": userLogin}

	var update = bson.M{
		"$addToSet": bson.M{
			"friends": friend,
		},
	}

	return base.GenericDao.Update(u.ColName, query, update)
}

func (u *UserDao) AddTalk(userLogin string, talk *model.Talk) error {
	query := bson.M{"login": userLogin}

	var update = bson.M{
		"$addToSet": bson.M{
			"talks": talk,
		},
	}

	return base.GenericDao.Update(u.ColName, query, update)
}

func (u *UserDao) search(
	q interface{}, selector interface{}, skip int, limit int) (*model.Users, error) {

	results := &model.Users{}

	query := func(c *mgo.Collection) error {
		if limit > 0 {
			if selector != nil {
				return c.Find(q).Skip(skip).Limit(limit).Select(selector).All(results)
			} else {
				return c.Find(q).Skip(skip).Limit(limit).All(results)
			}
		} else {
			if selector != nil {
				return c.Find(q).Skip(skip).Select(selector).All(results)
			} else {
				return c.Find(q).Skip(skip).All(results)
			}
		}
	}

	err := base.GenericDao.SearchCollection(u.ColName, query)

	return results, err
}

/*func (u *UserDao) UseSuggestion(userLogin string) error {

	query := bson.M{
		"login": userLogin,
	}

	update := bson.M{
		"$inc": bson.M{
			"suggestions": -1,
		},
	}

	return base.GenericDao.Update(u.colName, query, update)
}*/

func (u *UserDao) AddNewMessageInTalk(talk *model.Talk, message *model.Message) {
	for _, userLogin := range talk.Users {
		query := u.getUserTalkQuery(userLogin, talk.Id)

		var update = bson.M{}

		if message.Author == userLogin {
			update = bson.M{
				"$set": bson.M{
					"talks.$.last_message": message,
				},
			}
		} else {
			update = bson.M{
				"$set": bson.M{
					"talks.$.last_message": message,
					"talks.$.has_unread":   true,
				},
				"$inc": bson.M{
					"talks.$.ciphered_num": 1,
				},
			}
		}

		err := base.GenericDao.Update(u.ColName, query, update)

		if err != nil {
			log.WithFields(log.Fields{
				"query":  query,
				"update": update,
				"err":    err.Error(),
			}).Error("Error while update message in AddNewMessageInTalk")
		}
	}
}

func (u *UserDao) DecipherMessageInTalk(talk *model.Talk, message *model.Message) {
	for _, userLogin := range talk.Users {
		isAuthor := (message.Author == userLogin)

		query := u.getUserTalkQuery(userLogin, talk.Id)

		var update = bson.M{}

		if u.isLastMessageInTalk(talk, message) {
			if isAuthor {
				update = bson.M{
					"$set": bson.M{
						"talks.$.last_message": message,
						"talks.$.has_unread":   true,
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"talks.$.last_message": message,
					},
					"$inc": bson.M{
						"talks.$.ciphered_num": -1,
					},
				}
			}
		} else {
			if isAuthor {
				update = bson.M{
					"$set": bson.M{
						"talks.$.has_unread": true,
					},
				}
			} else {
				update = bson.M{
					"$inc": bson.M{
						"talks.$.ciphered_num": -1,
					},
				}
			}
		}

		err := base.GenericDao.Update(u.ColName, query, update)

		if err != nil {
			log.WithFields(log.Fields{
				"query":  query,
				"update": update,
				"err":    err.Error(),
			}).Error("Error while update message in DecipherMessageInTalk")
		}
	}
}

func (u *UserDao) MarkTalkAsReaded(userLogin string, talkId bson.ObjectId) error {
	query := u.getUserTalkQuery(userLogin, talkId)

	update := bson.M{
		"$set": bson.M{
			"talks.$.has_unread": false,
		},
	}

	return base.GenericDao.Update(u.ColName, query, update)
}

func (u *UserDao) isLastMessageInTalk(talk *model.Talk, message *model.Message) bool {
	if talk.LastMessage == nil {
		return false
	} else {
		return talk.LastMessage.Id == message.Id
	}
}

func (u *UserDao) getUserTalkQuery(login string, talkId bson.ObjectId) bson.M {
	query := bson.M{
		"login":     login,
		"talks._id": talkId,
	}

	return query
}

func (u *UserDao) GetUserTalk(login string, talkId bson.ObjectId) (*model.Talk, error) {
	query := bson.M{
		"login": login,
	}

	selector := bson.M{"talks": bson.M{"$elemMatch": bson.M{"_id": talkId}}}

	users, err := u.search(query, selector, 0, 0)

	if err != nil {
		return nil, err
	} else if len(*users) == 1 {
		return (*users)[0].Talks[0], err
	} else if len(*users) == 0 {
		return nil, &UserTalkNotFoundError{login, talkId}
	} else {
		return nil, &MoreThanOneUserError{login}
	}
}

func (u *UserDao) ClearCollection() {
	base.GenericDao.ClearCollection(u.ColName)
}

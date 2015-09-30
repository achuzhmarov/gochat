package dal

import (
	"gochat/dal/base"
	"gochat/model"

	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var TokenDaoInstance = TokenDao{"tokens"}

type TokenDao struct {
	ColName string
}

type TokenNotFoundError struct {
	TokenValue string
}

func (e *TokenNotFoundError) Error() string {
	return fmt.Sprintf("token %s not found", e.TokenValue)
}

func (t *TokenDao) Upsert(token *model.Token) (*model.Token, error) {
	err := base.GenericDao.Upsert(token, t.ColName)
	return token, err
}

func (t *TokenDao) GetByValue(tokenValue string) (*model.Token, error) {
	query := bson.M{"value": tokenValue}

	tokens, err := t.search(query, 0, 0)
	if err != nil {
		return nil, err
	} else if len(*tokens) == 1 {
		return (*tokens)[0], err
	} else if len(*tokens) == 0 {
		return nil, &TokenNotFoundError{tokenValue}
	} else {
		return (*tokens)[0], err
	}
}

func (t *TokenDao) search(
	q interface{}, skip int, limit int) (*model.Tokens, error) {

	results := &model.Tokens{}

	query := func(c *mgo.Collection) error {
		if limit > 0 {
			return c.Find(q).Skip(skip).Limit(limit).All(results)
		} else {
			return c.Find(q).Skip(skip).All(results)
		}
	}

	err := base.GenericDao.SearchCollection(t.ColName, query)

	return results, err
}

func (t *TokenDao) ClearCollection() {
	base.GenericDao.ClearCollection(t.ColName)
}

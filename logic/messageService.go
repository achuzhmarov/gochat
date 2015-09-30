package logic

import (
	"gochat/dal"
	"gochat/model"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var MessageService = messageService{dal.MessageDaoInstance, dal.UserDaoInstance}

type messageService struct {
	messageDao dal.MessageDao
	userDao    dal.UserDao
}

func (s *messageService) GetMessagesByTalk(userLogin string, talkIdString string) (*model.Messages, error) {
	talkId := bson.ObjectIdHex(talkIdString)
	result, err := s.messageDao.GetAllByTalk(talkId)

	if err != nil {
		return nil, err
	}

	markErr := s.userDao.MarkTalkAsReaded(userLogin, talkId)
	if markErr != nil {
		log.WithFields(log.Fields{
			"user":   userLogin,
			"talkId": talkId,
			"err":    markErr.Error(),
		}).Error("Error in markTalkAsReaded")
	}

	return result, err
}

func (s *messageService) GetUnreadMessagesByTalk(userLogin string, talkIdString string) (*model.Messages, error) {
	talkId := bson.ObjectIdHex(talkIdString)

	talk, err := s.userDao.GetUserTalk(userLogin, talkId)
	if err != nil {
		return nil, err
	}

	if talk.HasUnread {
		return s.GetMessagesByTalk(userLogin, talkIdString)
	} else {
		return nil, nil
	}
}

func (s *messageService) CreateMessage(userLogin string, message *model.Message) (*model.Message, error) {
	talk, err := s.userDao.GetUserTalk(userLogin, message.TalkId)
	if err != nil {
		return nil, err
	}

	result, err := s.messageDao.Upsert(message)

	//error logging inside, not critical
	s.userDao.AddNewMessageInTalk(talk, message)

	return result, err
}

func (s *messageService) DecipherMessage(userLogin string, message *model.Message) (*model.Message, error) {
	messageDB, err := s.messageDao.GetById(message.Id)
	if err != nil {
		return nil, err
	}

	talk, err := s.userDao.GetUserTalk(userLogin, messageDB.TalkId)
	if err != nil {
		return nil, err
	}

	messageDB.Words = message.Words
	messageDB.CheckDeciphered()

	result, err := s.messageDao.Upsert(messageDB)

	if messageDB.Deciphered {
		//error logging inside, not critical
		s.userDao.DecipherMessageInTalk(talk, messageDB)
	}

	return result, err
}

/*func (s *messageService) CreateMessage(author string, talkId bson.ObjectId, text string,
	cipherType cipher.CipherType) (*model.Message, error) {

	words := getWords(text)
	message := model.CreateMessage(talkId, author, false, cipherType, words)

	return s.messageDao.Upsert(message)
}*/

/*func (s *messageService) FailMessage(message *model.Message) (*model.Message, error) {
	for _, word := range message.Words {
		if word.WordType == wordType.New {
			word.WordType = wordType.Failed
		}
	}

	message.Deciphered = true

	return s.messageDao.Upsert(message)
}*/

/*func (s *messageService) DecipherByGuess(message *model.Message, guessText string) (*model.Message, error) {
	guesses := strings.Split(guessText, " ")

	for _, guess := range guesses {
		for _, word := range message.Words {
			if strings.ToUpper(word.Text) == strings.ToUpper(guess) {
				word.WordType = wordType.Success
			}
		}
	}

	message.CheckDeciphered()

	return s.messageDao.Upsert(message)
}*/

/*func (s *messageService) DecipherBySuggestion(message *model.Message, wordIndex int) (*model.Message, error) {
	message.Words[wordIndex].WordType = wordType.Success
	message.CheckDeciphered()

	return s.messageDao.Upsert(message)
}

func getWords(text string) *model.Words {
	//TODO - implement
	return &model.Words{}
}*/

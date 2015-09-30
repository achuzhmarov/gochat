package logic

import (
	"gochat/dal"
	"gochat/model"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var UserService = userService{dal.UserDaoInstance}

type userService struct {
	userDao dal.UserDao
}

type NoSuggestions struct {
	UserLogin string
}

func (e *NoSuggestions) Error() string {
	return fmt.Sprintf("User %s has 0 suggestions", e.UserLogin)
}

type NegativeSuggestionsCount struct {
	Count int
}

func (e *NegativeSuggestionsCount) Error() string {
	return fmt.Sprintf("SuggestionsCount must be positive. Get %d", e.Count)
}

type UserAlreadyExists struct {
	UserLogin string
}

func (e *UserAlreadyExists) Error() string {
	return fmt.Sprintf("User %s already exists", e.UserLogin)
}

func (s *userService) RegisterNewUser(user *model.UserAuth) (*model.User, error) {
	_, err := s.userDao.GetByLogin(user.Login)
	_, ok := err.(*dal.UserNotFoundError)

	if ok {
		newUser := model.CreateUser(user.Login, getHash(user.Password), "")
		user, err := s.userDao.Upsert(newUser)

		//friend with tigra
		if user.Login != "tigra" {
			_, friendErr := s.MakeFriends("tigra", user.Login)

			if friendErr != nil {
				log.WithFields(log.Fields{
					"userLogin":   "tigra",
					"friendLogin": user.Login,
					"err":         friendErr.Error(),
				}).Error("Error while makeFriends with tigra")
			}
		}

		return user, err
	} else {
		return nil, &UserAlreadyExists{user.Login}
	}
}

func getHash(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword)
}

func (s *userService) GetUserByLogin(userLogin string) (*model.User, error) {
	return s.userDao.GetByLogin(userLogin)
}

func (s *userService) GetNewFriends(userLogin string, searchString string) ([]string, error) {
	user, err := s.userDao.GetByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	allFriends, err := s.userDao.GetAllPotentialFriends(searchString)

	resultFriends := []string{}

	for _, friend := range *allFriends {
		if !contains(user.Friends, friend.Login) && userLogin != friend.Login {
			resultFriends = append(resultFriends, friend.Login)
		}
	}

	return resultFriends, nil
}

func contains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (s *userService) MakeFriends(userLogin string, friendLogin string) (*model.Talk, error) {
	//TODO - what to do if error in process?

	err := s.userDao.AddFriend(userLogin, friendLogin)
	if err != nil {
		return nil, err
	}

	err = s.userDao.AddFriend(friendLogin, userLogin)
	if err != nil {
		log.WithFields(log.Fields{
			"userLogin":   userLogin,
			"friendLogin": friendLogin,
			"err":         err.Error(),
		}).Error("Error while makeFriends in addUserToFriend")
	}

	talkId := bson.NewObjectId()
	talk := model.CreateTalk(talkId, []string{userLogin, friendLogin})

	err = s.userDao.AddTalk(userLogin, talk)
	if err != nil {
		log.WithFields(log.Fields{
			"userLogin":   userLogin,
			"friendLogin": friendLogin,
			"err":         err.Error(),
		}).Error("Error while makeFriends in addTalkToUser")
	}

	err = s.userDao.AddTalk(friendLogin, talk)
	if err != nil {
		log.WithFields(log.Fields{
			"userLogin":   userLogin,
			"friendLogin": friendLogin,
			"err":         err.Error(),
		}).Error("Error while makeFriends in addTalkToFriend")
	}

	return talk, nil
}

func (s *userService) GetUnreadTalks(userLogin string) (model.Talks, error) {
	user, err := s.userDao.GetByLogin(userLogin)
	if err != nil {
		return nil, err
	}

	var talks model.Talks

	for _, talk := range user.Talks {
		if talk.HasUnread {
			talks = append(talks, talk)
		}
	}

	return talks, nil
}

func (s *userService) UseSuggestions(userLogin string, count int) (*model.User, error) {
	if count < 1 {
		return nil, &NegativeSuggestionsCount{count}
	}

	user, err := s.userDao.GetByLogin(userLogin)

	if err != nil {
		return nil, err
	}

	if user.Suggestions >= count {
		user.Suggestions -= count
	} else {
		log.WithFields(log.Fields{
			"suggestions": user.Suggestions,
			"count":       count,
		}).Warn("Not enough suggestions")

		user.Suggestions = 0
	}

	return s.userDao.Upsert(user)
}

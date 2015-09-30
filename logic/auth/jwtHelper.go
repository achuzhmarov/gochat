package auth

import (
	"gochat/dal"
	"gochat/logic"
	"gochat/model"
	"gochat/settings"

	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = time.Second * 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(userLogin string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userLogin
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) Authenticate(requestUser *model.UserAuth) bool {
	userDb, err := logic.UserService.GetUserByLogin(requestUser.Login)

	if err != nil {
		//TODO - add logging
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(requestUser.Password)) == nil
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) time.Time {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		return tm.Add(expireOffset)
	}

	return time.Now().Add(expireOffset)
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	expireAt := backend.getTokenRemainingValidity(token.Claims["exp"])
	tokenDb := model.CreateToken(tokenString, expireAt)
	_, err := dal.TokenDaoInstance.Upsert(tokenDb)
	return err
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(tokenString string) bool {
	token, err := dal.TokenDaoInstance.GetByValue(tokenString)

	if token != nil {
		return true
	}

	if err != nil {
		//TODO - log error
	}

	return false
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(settings.Get().PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(settings.Get().PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
